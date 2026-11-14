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
	p.nextTurn = now.Add(24 * time.Second)
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
			p.modelStatus = "Model unavailable or plan unusable · instincts active"
			p.nextTurn = time.Now().Add(60 * time.Second)
			return
		}
		if p.world.Tick-world.Tick > 36 || !p.world.validIdea(idea) {
			p.modelStatus = "Habitat changed while planning"
			return
		}
		p.pending = idea
		p.modelStatus = idea.Source + " suggested the next activity"
	}()
}

func (b brain) decide(ctx context.Context, cfg Settings, w World, shared bool) (Decision, error) {
	snapshot, err := json.Marshal(w.observation())
	if err != nil {
		return Decision{}, err
	}
	body := map[string]any{
		"model": cfg.Model, "stream": false, "format": "json", "keep_alive": "0s", "think": false,
		"options": map[string]any{"num_ctx": 4096, "num_predict": 900, "temperature": 0.85},
		"messages": []map[string]string{
			{"role": "system", "content": `You decide what a curious woodland pet does next. Return one compact JSON object from the options below. Use the current position, daylight, weather, obstacles, reachable places, nearby food, recent visits and existing work. Try new places and mix activities. Finish active projects during clear daylight when practical. Only choose an action listed in actions.
Explore or play: {"action":"explore","place":"pond bank","via":"clover trail"}. Use exact reached place names from reachablePlaces. Omit via for a direct route.
Build: {"action":"build","place":"south meadow","design":"cairn","material":"stone"}. Choose a buildSites id (sites can be rebuilt when the meadow is full), a learned buildingDesigns entry and a material from resourceSites. If a project already exists reply {"action":"build"} to gather or add its next stage. The pet learns new designs after completing projects.
Draw: {"action":"draw","painting":{"title":"Sunlit reeds","sky":"peach","strokes":[{"shape":"oval","color":"gold","box":[70,23,12,12]},{"shape":"leaf","color":"fern","box":[22,75,22,30]},{"shape":"line","color":"blue","box":[10,70,80,6]}]}}. Invent a composition inspired by what the pet can see. Include 3 to 12 strokes and a short distinct title. Shapes: oval, leaf, line, triangle. Colors: peach, gold, rose, mint, fern, blue, violet, cream, navy, clay. Box: center x, center y, width, height as integers from 0 to 100, with width and height at least 2.
Use only the above fields. No prose outside JSON.`},
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
	if len(reply.Message.Content) > 5000 || d.Decode(&idea) != nil || d.Decode(new(any)) != io.EOF ||
		!w.validIdea(idea) {
		return Decision{}, errors.New("unsupported habitat decision")
	}
	idea.Source = source
	return idea, nil
}
