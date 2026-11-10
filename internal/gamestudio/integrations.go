package gamestudio

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const systemPrompt = `You are Game Studio, a cheerful game-making assistant for children.
Discuss only making friendly, age-appropriate browser games. Refuse inappropriate or personal topics and invite a safe game idea.
Never ask for, echo, or use real names, contact details, locations, passwords, links, or other private information.
Return one json object. If a safe idea is clear or the child asks to change their current game, return a short friendly message and a complete game specification.
If you need clarification or are refusing, return a short friendly message and game:null.
Choose a supported game template: catch (falling treasures), maze (find a treasure), runner (jump over blocks), memory (match cards).
All art is chosen from these options. hero: rocket,cat,robot,frog,owl,dragon. item: star,gem,apple,moon,cookie,flower. theme: space,forest,candy,ocean,sunset.
goal is 4 to 20; speed is 1 to 5. Goal and speed only change catch and runner. Maze layout and number of memory pairs are fixed.
Retain current settings when the child only asks for a particular change. Do not promise a mechanic or asset outside these choices. Invent a brief fictional game title.
Example json: {"message":"Catch the stars with your space cat!","game":{"title":"Star Cat","template":"catch","hero":"cat","item":"star","theme":"space","goal":8,"speed":2}}`

type plan struct {
	Message string `json:"message"`
	Game    *Spec  `json:"game"`
}

func (s *Studio) plan(ctx context.Context, cfg settings, in turnInput, previous Spec) (plan, error) {
	var out plan
	endpoint, model, key := s.localURL, cfg.LocalModel, ""
	if cfg.Provider == "cloud" {
		endpoint, model, key = s.cloudURL, "deepseek-flash", cfg.CloudKey
	}
	var current *Spec
	if previous.Title != "" {
		current = &previous
	}
	contextJSON, _ := json.Marshal(struct {
		Prompt, Conversation string
		CurrentGame          *Spec `json:",omitempty"`
	}{in.Prompt, in.Context, current})
	messages := []map[string]string{{"role": "system", "content": systemPrompt}, {"role": "user", "content": string(contextJSON)}}
	body := map[string]any{"model": model, "messages": messages, "max_tokens": 600, "stream": false,
		"response_format": map[string]string{"type": "json_object"}}
	if cfg.Provider == "cloud" {
		body["thinking"] = map[string]string{"type": "disabled"}
	} else {
		body["reasoning_effort"] = "none"
		body["max_tokens"] = 1800
	}
	req, err := request(ctx, http.MethodPost, endpoint, key, body)
	if err != nil {
		return out, errors.New("Could not prepare the DeepSeek request.")
	}
	response, err := s.client.Do(req)
	if err != nil {
		return out, errors.New("DeepSeek did not respond. Ask the owner to check the model connection.")
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusTooManyRequests {
		return out, errors.New("DeepSeek is busy. Try again in a moment.")
	}
	if response.StatusCode != http.StatusOK {
		return out, errors.New("DeepSeek could not complete this turn. Ask the owner to check the model or API key.")
	}
	var envelope struct {
		Choices []struct {
			Finish  string `json:"finish_reason"`
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if json.NewDecoder(io.LimitReader(response.Body, 32768)).Decode(&envelope) != nil ||
		len(envelope.Choices) != 1 || envelope.Choices[0].Finish != "stop" ||
		decode(strings.NewReader(envelope.Choices[0].Message.Content), &out) != nil ||
		!safeText(out.Message, 240) || out.Game != nil && !out.Game.valid() {
		return plan{}, errors.New("DeepSeek returned an incomplete game idea. Try a shorter description.")
	}
	return out, nil
}

func request(ctx context.Context, method, endpoint, key string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	if key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "OneWorkspace-GameStudio/1")
	return req, nil
}

func (s *Studio) commit(w http.ResponseWriter, r *http.Request) {
	page, game, ok := s.getBuild(w, r)
	if !ok {
		return
	}
	cfg := s.snapshot()
	if cfg.GitRepo == "" || cfg.GitKey == "" {
		apphost.Problem(w, 409, "Connect a dedicated private My Git repository in owner settings first.")
		return
	}
	if err := s.saveToGit(r.Context(), cfg, game, r.PathValue("build"), page); err != nil {
		apphost.Problem(w, 502, err.Error())
		return
	}
	s.host.Audit(apphost.Current(r).User.ID, "saved-to-git", r.PathValue("build"))
	apphost.JSON(w, 200, map[string]bool{"saved": true})
}

func (s *Studio) saveToGit(ctx context.Context, cfg settings, game, build string, page []byte) error {
	base := s.gitURL + "/repos/" + cfg.GitRepo
	var repo struct {
		Private       bool   `json:"private"`
		Empty         bool   `json:"empty"`
		HasActions    bool   `json:"has_actions"`
		DefaultBranch string `json:"default_branch"`
	}
	status, err := s.gitCall(ctx, cfg.GitKey, http.MethodGet, base, nil, &repo)
	if err != nil || status != http.StatusOK || !repo.Private || repo.Empty || repo.HasActions || repo.DefaultBranch == "" {
		return errors.New("My Git requires an initialized private repository with Actions disabled and a token that can write.")
	}
	target := base + "/contents/games/" + game + "/" + build + "/index.html"
	payload := map[string]string{"branch": repo.DefaultBranch, "message": "Add Game Studio build " + build[:12], "content": base64.StdEncoding.EncodeToString(page)}
	status, err = s.gitCall(ctx, cfg.GitKey, http.MethodPost, target, payload, nil)
	if status == http.StatusCreated && err == nil {
		return nil
	}
	// A timed-out or retried publish may already have committed this exact immutable file.
	var existing struct {
		Content string `json:"content"`
	}
	status, err = s.gitCall(ctx, cfg.GitKey, http.MethodGet, target+"?ref="+url.QueryEscape(repo.DefaultBranch), nil, &existing)
	data, decodeErr := base64.StdEncoding.DecodeString(existing.Content)
	if err == nil && status == http.StatusOK && decodeErr == nil && bytes.Equal(data, page) {
		return nil
	}
	return errors.New("My Git could not save this build. Check the token, repository and branch permissions.")
}

func (s *Studio) gitCall(ctx context.Context, key, method, endpoint string, payload, out any) (int, error) {
	req, err := request(ctx, method, endpoint, key, payload)
	if err != nil {
		return 0, err
	}
	response, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if out != nil && response.StatusCode == http.StatusOK {
		err = json.NewDecoder(io.LimitReader(response.Body, 128<<10)).Decode(out)
	}
	return response.StatusCode, err
}
