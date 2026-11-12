package pet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

type Settings struct {
	Mode       string `json:"mode"`
	Model      string `json:"model"`
	ServiceKey string `json:"serviceKey"`
	Owner      string `json:"owner"`
}

var modelName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_./:-]{0,99}$`)
var serviceKey = regexp.MustCompile(`^[!-~]{32,200}$`)

type brain struct {
	localURL  string
	sharedURL string
	client    *http.Client
}

func newBrain(port string) brain {
	b := brain{
		localURL: "http://127.0.0.1:11434/api/chat",
		client: &http.Client{
			Timeout:       65 * time.Second,
			Transport:     &http.Transport{MaxIdleConns: 2, IdleConnTimeout: 30 * time.Second},
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
	if n, err := strconv.Atoi(port); err == nil && n > 0 && n <= 65535 {
		b.sharedURL = "http://127.0.0.1:" + strconv.Itoa(n) + "/_internal/chat"
	}
	return b
}

func (p *Habitat) configure(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Mode     string `json:"mode"`
		Model    string `json:"model"`
		Key      string `json:"key"`
		ClearKey bool   `json:"clearKey"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	owner := apphost.Current(r).User.ID
	next := Settings{Mode: in.Mode, Model: strings.TrimSpace(in.Model), ServiceKey: strings.TrimSpace(in.Key), Owner: owner}
	if in.ClearKey {
		next.ServiceKey = ""
	} else if next.ServiceKey == "" && p.settings.Owner == owner {
		next.ServiceKey = p.settings.ServiceKey
	}
	local := next.Mode == "local" || next.Mode == "both"
	shared := next.Mode == "service" || next.Mode == "both"
	if next.Mode != "instincts" && !local && !shared ||
		local && !modelName.MatchString(next.Model) ||
		next.ServiceKey != "" && !serviceKey.MatchString(next.ServiceKey) ||
		shared && (next.ServiceKey == "" || p.brain.sharedURL == "") {
		apphost.Problem(w, 400, "Enter an installed Ollama model or pair Pet with the installed AI Service for this mode.")
		return
	}
	if err := platform.SaveSecret(p.settingsPath(), next); err != nil {
		apphost.Problem(w, 500, "Decision settings could not be saved.")
		return
	}
	if p.cancel != nil {
		p.cancel()
	}
	p.generation++
	p.settings, p.pending, p.modelStatus = next, Decision{}, ""
	p.nextTurn, p.sharedNext = time.Time{}, false
	apphost.JSON(w, 200, map[string]bool{"saved": true})
}

// schedule captures only simulation data. An outstanding decision never blocks a tick.
func (p *Habitat) schedule(ctx context.Context, workers *sync.WaitGroup, now time.Time) {
	cfg := p.settings
	if p.world == nil || p.cancel != nil || cfg.Mode == "instincts" || cfg.Owner == "" ||
		now.Before(p.nextTurn) || p.world.urgent() != "" || p.world.Action == "sleep" ||
		p.world.Action == "eat" || p.pending.Action != "" {
		return
	}
	if p.host.Config.Issuer != "" && (p.host.Config.Owner == "" ||
		cfg.Owner != apphost.Identity(p.host.Config.Issuer, p.host.Config.Owner)) {
		return
	}
	shared := cfg.Mode == "service" || cfg.Mode == "both" && p.sharedNext
	p.sharedNext = !p.sharedNext
	p.nextTurn = now.Add(90 * time.Second)
	turn, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	gen := p.generation
	world := *p.world
	workers.Add(1)
	go func() {
		defer workers.Done()
		defer cancel()
		idea, err := p.brain.decide(turn, cfg, world, shared)
		p.mu.Lock()
		defer p.mu.Unlock()
		p.cancel = nil
		if gen != p.generation || turn.Err() != nil {
			return
		}
		if err != nil {
			p.modelStatus = "Model unavailable · instincts active"
			return
		}
		p.pending = idea
		p.modelStatus = idea.Source + " suggested the next activity"
	}()
}

func (b brain) decide(ctx context.Context, cfg Settings, w World, shared bool) (Decision, error) {
	available := []string{}
	for _, action := range []string{"explore", "play", "draw", "build"} {
		if w.allows(action) {
			available = append(available, action)
		}
	}
	snapshot, err := json.Marshal(map[string]any{
		"actions": available, "currentAction": w.Action, "mood": w.Mood, "minuteOfDay": w.Minute,
		"weather": w.Weather, "energy": int(w.Energy), "fullness": int(w.Fullness),
		"happiness": int(w.Happiness), "curiosity": int(w.Curiosity), "food": len(w.Food), "creations": len(w.Creations),
	})
	if err != nil {
		return Decision{}, err
	}
	body := map[string]any{
		"model": cfg.Model, "stream": false, "format": "json", "keep_alive": "0s",
		"options": map[string]any{"num_ctx": 2048, "num_predict": 256, "temperature": 0.55},
		"messages": []map[string]string{
			{"role": "system", "content": `You choose the next activity for a curious woodland pet. Use only an action listed in actions and return exactly {"action":"explore","style":0}. Pick an integer style 0–3 for art or structures. Prefer variety and play when happiness is low. Do not output prose or extra fields.`},
			{"role": "user", "content": string(snapshot)},
		},
	}
	url, source := b.localURL, "Local model"
	if shared {
		url, source = b.sharedURL, "AI Service"
		delete(body, "model")
		delete(body, "keep_alive")
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return Decision{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return Decision{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if shared {
		req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)
		req.Header.Set("X-AI-App", "pet")
		req.Header.Set("X-AI-User", cfg.Owner)
	}
	res, err := b.client.Do(req)
	if err != nil {
		return Decision{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Decision{}, errors.New("inference endpoint unavailable")
	}
	var reply struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 128<<10)).Decode(&reply); err != nil {
		return Decision{}, err
	}
	var idea Decision
	d := json.NewDecoder(strings.NewReader(reply.Message.Content))
	d.DisallowUnknownFields()
	if len(reply.Message.Content) > 512 || d.Decode(&idea) != nil || d.Decode(new(any)) != io.EOF ||
		!w.allows(idea.Action) || idea.Style < 0 || idea.Style > 3 {
		return Decision{}, errors.New("unsupported habitat decision")
	}
	idea.Source = source
	return idea, nil
}
