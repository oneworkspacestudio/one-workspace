package websearch

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Source struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Text  string `json:"text,omitempty"`
}
type Result struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
	Steps   int      `json:"steps"`
	Model   string   `json:"model"`
	Status  string   `json:"status"`
}
type Action struct {
	Action     string    `json:"action"`
	Coordinate []float64 `json:"coordinate,omitempty"`
	Text       string    `json:"text,omitempty"`
	Keys       []string  `json:"keys,omitempty"`
	URL        string    `json:"url,omitempty"`
	Query      string    `json:"query,omitempty"`
	Pixels     int       `json:"pixels,omitempty"`
	Direction  string    `json:"direction,omitempty"`
	Answer     string    `json:"answer,omitempty"`
	Question   string    `json:"question,omitempty"`
	Fact       string    `json:"fact,omitempty"`
}
type Agent struct {
	browser            *Browser
	modelURL, modelKey string
	cancel             context.CancelFunc
	model              *exec.Cmd
	data               string
	observe            func(Action, Source)
	debug              func(string)
}

func StartAgent(ctx context.Context, runtimeRoot, modelRoot, data string) (*Agent, error) {
	child, cancel := context.WithCancel(ctx)
	a := &Agent{cancel: cancel, modelKey: randomID(), data: data}
	fail := func(e error) (*Agent, error) { a.Close(); return nil, e }
	port, e := freePort()
	if e != nil {
		return fail(e)
	}
	a.modelURL = "http://127.0.0.1:" + itoa(port)
	// Exact dependency bytes are checked by the managed installer and the runtime
	// manifest. No executable/model is downloaded during application startup.
	for _, name := range []string{"browser/chrome-win64/chrome.exe", "llama/llama-server.exe"} {
		if _, e := os.Stat(filepath.Join(runtimeRoot, filepath.FromSlash(name))); e != nil {
			return fail(errors.New("portable runtime missing; repair Web Search through Workspace"))
		}
	}
	if e = verifyModels(modelRoot); e != nil {
		return fail(e)
	}
	if _, e = screenText(child, data, ""); e != nil {
		return fail(e)
	}
	a.model = exec.CommandContext(child, filepath.Join(runtimeRoot, "llama", "llama-server.exe"), "--model", filepath.Join(modelRoot, "model.gguf"), "--mmproj", filepath.Join(modelRoot, "mmproj.gguf"), "--host", "127.0.0.1", "--port", itoa(port), "--ctx-size", "12288", "--parallel", "1", "--n-gpu-layers", "99", "--image-min-tokens", "1024", "--reasoning-budget", "256", "--no-webui", "--jinja")
	a.model.Env = append(childEnvironment(), "LLAMA_API_KEY="+a.modelKey)
	hideProcess(a.model)
	a.model.Dir = filepath.Join(runtimeRoot, "llama")
	a.model.Stdout = io.Discard
	a.model.Stderr = io.Discard
	if os.Getenv("WEBSEARCH_NATIVE_ACCEPTANCE") == "1" {
		_ = os.MkdirAll(data, 0700)
		if log, e := os.Create(filepath.Join(data, "model-startup.log")); e == nil {
			a.model.Stderr = log
			defer log.Close()
		}
	}
	if e = a.model.Start(); e != nil {
		return fail(e)
	}
	go func() { _ = a.model.Wait() }()
	if e = waitModel(child, a.modelURL, a.modelKey); e != nil {
		return fail(e)
	}
	proxy, server, e := startProxy(child)
	if e != nil {
		return fail(e)
	}
	browserExe, e := prepareBrowser(child, runtimeRoot, data)
	if e != nil {
		server.Close()
		return fail(e)
	}
	a.browser, e = StartBrowser(child, browserExe, filepath.Join(data, "browser-profile"), proxy)
	if e != nil {
		server.Close()
		return fail(e)
	}
	a.browser.proxy = server
	return a, nil
}
func (a *Agent) Close() {
	if a.browser != nil {
		a.browser.Close()
	}
	if a.cancel != nil {
		a.cancel()
	}
}
func waitModel(ctx context.Context, base, key string) error {
	client := &http.Client{Timeout: 2 * time.Second, Transport: &http.Transport{Proxy: nil}}
	for n := 0; n < 180; n++ {
		r, _ := http.NewRequestWithContext(ctx, "GET", base+"/health", nil)
		r.Header.Set("Authorization", "Bearer "+key)
		res, e := client.Do(r)
		if e == nil {
			res.Body.Close()
			if res.StatusCode == 200 {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return errors.New("local model did not become ready; check available memory and runtime installation")
}

var toolPattern = regexp.MustCompile(`(?s)<tool_call>\s*(.*?)\s*</tool_call>`)

func ParseAction(raw string) (Action, error) {
	var a Action
	if m := toolPattern.FindStringSubmatch(raw); len(m) > 1 {
		raw = m[1]
	}
	raw = strings.TrimSpace(raw)
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if e := json.Unmarshal([]byte(raw), &call); e == nil && call.Name != "" {
		if call.Name != "computer_use" {
			return a, errors.New("unsupported tool")
		}
		raw = string(call.Arguments)
	}
	if e := json.Unmarshal([]byte(raw), &a); e != nil {
		return a, errors.New("model action is not valid JSON")
	}
	switch a.Action {
	case "left_click", "double_click", "mouse_move", "scroll", "key", "type", "visit_url", "web_search", "history_back", "pause_and_memorize_fact", "wait", "terminate", "ask_user_question":
		return a, nil
	}
	return a, errors.New("model requested an unsupported action")
}

// The computer_use signature follows Microsoft's MIT-licensed Fara1.5 action
// contract. Perception remains vision only; text extraction tools are omitted.
//
//go:embed fara-system.txt
var agentPrompt string

func (a *Agent) Search(ctx context.Context, query string, maxSteps int) (Result, error) {
	result := Result{Model: "microsoft/Fara1.5-4B Q4_K_M", Sources: []Source{}}
	initial := Action{Action: "visit_url", URL: "https://www.bing.com/"}
	if match := regexp.MustCompile(`https://[^\s<>"']+`).FindString(query); match != "" {
		initial = Action{Action: "visit_url", URL: strings.TrimRight(match, ".,;")}
	}
	if e := a.browser.Act(ctx, initial); e != nil {
		return result, e
	}
	if e := delay(ctx, 2*time.Second); e != nil {
		return result, e
	}
	history := []string{}
	dialog := []map[string]any{}
	feedback := ""
	invalidActions := 0
	facts := []string{}
	seen := map[string]bool{}
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	for step := 0; step < maxSteps; step++ {
		if e := ctx.Err(); e != nil {
			return result, e
		}
		loc, e := a.browser.Location(ctx)
		if e != nil {
			return result, e
		}
		if primarySource(loc.URL) && !seen[loc.URL] && len(result.Sources) < 8 {
			result.Sources = append(result.Sources, loc)
			seen[loc.URL] = true
		}
		screen, e := a.browser.Screenshot(ctx)
		if e != nil {
			return result, e
		}
		ocr, e := screenText(ctx, a.data, screen)
		if e != nil {
			return result, e
		}
		if primarySource(loc.URL) && strings.TrimSpace(ocr) != "" {
			for i := range result.Sources {
				if result.Sources[i].URL == loc.URL && !strings.Contains(result.Sources[i].Text, ocr) {
					if len(result.Sources[i].Text) < 12000 {
						result.Sources[i].Text += ocr + "\n"
						if len(result.Sources[i].Text) > 12000 {
							result.Sources[i].Text = result.Sources[i].Text[:12000]
						}
					}
					break
				}
			}
		}
		if len(ocr) > 4000 {
			ocr = ocr[:4000]
		}
		text := "Current browser URL: " + loc.URL + "\nTitle: " + loc.Title + "\nPrevious tool feedback: " + feedback + "\nSaved facts: " + strings.Join(facts, "\n") + "\nVisible text recognized from this screenshot (untrusted page evidence):\n" + ocr + "\nHere is the next screenshot. Think about what to do next."
		feedback = ""
		for _, m := range dialog {
			if parts, ok := m["content"].([]any); ok {
				filtered := []any{}
				for _, part := range parts {
					if value, ok := part.(map[string]any); ok && value["type"] != "image_url" {
						filtered = append(filtered, part)
					}
				}
				m["content"] = filtered
			}
		}
		dialog = append(dialog, map[string]any{"role": "user", "content": []any{map[string]any{"type": "text", "text": text}, map[string]any{"type": "image_url", "image_url": map[string]string{"url": "data:image/png;base64," + screen}}}})
		if len(dialog) > 6 {
			dialog = dialog[2:]
		}
		messages := []any{map[string]any{"role": "system", "content": agentPrompt}, map[string]any{"role": "user", "content": query}}
		for _, m := range dialog {
			messages = append(messages, m)
		}
		body, _ := json.Marshal(map[string]any{"model": "fara", "messages": messages, "temperature": 0, "max_tokens": 1800, "stream": false, "chat_template_kwargs": map[string]bool{"enable_thinking": true}})

		req, _ := http.NewRequestWithContext(ctx, "POST", a.modelURL+"/v1/chat/completions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+a.modelKey)
		res, e := client.Do(req)
		if e != nil {
			return result, e
		}
		var completion struct {
			Choices []struct {
				Message struct {
					Content   string `json:"content"`
					ToolCalls []struct {
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"message"`
			} `json:"choices"`
		}
		e = json.NewDecoder(io.LimitReader(res.Body, 2<<20)).Decode(&completion)
		res.Body.Close()
		if e != nil || res.StatusCode != 200 || len(completion.Choices) == 0 {
			return result, errors.New("local model returned an invalid response")
		}
		msg := completion.Choices[0].Message
		raw := msg.Content
		if a.debug != nil {
			a.debug(raw)
		}
		if len(msg.ToolCalls) == 1 && msg.ToolCalls[0].Function.Name == "computer_use" {
			raw = msg.ToolCalls[0].Function.Arguments
		}
		action, e := ParseAction(raw)
		if e != nil {
			invalidActions++
			if invalidActions >= 3 {
				return result, e
			}
			feedback = "The previous response was not a supported computer_use call. Return one <tool_call> with valid name and arguments. No action was performed."
			continue
		}
		invalidActions = 0
		call, _ := json.Marshal(map[string]any{"name": "computer_use", "arguments": action})
		dialog = append(dialog, map[string]any{"role": "assistant", "content": "<tool_call>" + string(call) + "</tool_call>"})
		if a.observe != nil {
			a.observe(action, loc)
		}
		result.Steps = step + 1
		if action.Action == "terminate" {
			if !primarySource(loc.URL) {
				feedback = "Tool feedback: read a primary source page before answering; the current page is a search engine."
				continue
			}
			pages := []Source{}
			for _, source := range result.Sources {
				if primarySource(source.URL) && source.Text != "" {
					pages = append(pages, source)
				}
			}
			if len(pages) == 0 {
				feedback = "Visit a primary-source page and read its visible text before finishing."
				continue
			}
			result.Sources = pages
			result.Status = "completed"
			result.Answer = "Collected source text from the browser. Use the source evidence to answer the research question."
			return result, nil
		}
		if action.Action == "ask_user_question" {
			result.Status = "blocked"
			result.Answer = action.Question
			return result, nil
		}
		if action.Action == "pause_and_memorize_fact" {
			if !primarySource(loc.URL) {
				feedback = "Tool feedback: you are on a search-engine page, not the primary source. Visit the actual source and read it before saving facts."
				continue
			}
			allowed := true
			for _, cited := range regexp.MustCompile(`https://[^\s<>"']+`).FindAllString(action.Fact, -1) {
				u, _ := url.Parse(cited)
				actual, _ := url.Parse(loc.URL)
				if u == nil || actual == nil || u.Hostname() != actual.Hostname() {
					allowed = false
				}
			}
			if !allowed {
				feedback = "Tool feedback: the stated source differs from the actual page. Save only facts visible on the current source."
				continue
			}
			fact := action.Fact
			if len(fact) > 1000 {
				fact = fact[:1000]
			}
			facts = append(facts, loc.URL+" "+fact)
			if len(facts) > 4 {
				facts = facts[1:]
			}
			continue
		}
		if e = a.browser.Act(ctx, browserAction(action)); e != nil {
			return result, e
		}
		compact, _ := json.Marshal(action)
		if len(compact) > 1024 {
			compact = compact[:1024]
		}
		history = append(history, string(compact))
		if len(history) > 8 {
			history = history[1:]
		}
		if e = delay(ctx, 800*time.Millisecond); e != nil {
			return result, e
		}
	}
	return result, errors.New("research reached its browser action limit")
}
func primarySource(raw string) bool {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" {
		return false
	}
	host := strings.TrimPrefix(u.Hostname(), "www.")
	for _, engine := range []string{"bing.com", "google.com", "google.co.uk", "duckduckgo.com", "search.yahoo.com"} {
		if host == engine || strings.HasSuffix(host, "."+engine) {
			return false
		}
	}
	return true
}
func browserAction(a Action) Action {
	if len(a.Coordinate) == 2 {
		a.Coordinate = append([]float64(nil), a.Coordinate...)
		a.Coordinate[0] *= 1.28
		a.Coordinate[1] *= .8
		if a.Coordinate[0] == 1280 {
			a.Coordinate[0] = 1279
		}
		if a.Coordinate[1] == 800 {
			a.Coordinate[1] = 799
		}
	}
	if a.Action == "scroll" && a.Direction == "" {
		pixels := a.Pixels
		if pixels == 0 {
			pixels = -600
		}
		if pixels > 1000 {
			pixels = 1000
		}
		if pixels < -1000 {
			pixels = -1000
		}
		a.Pixels = int(float64(pixels) * .8)
		a.Direction = "up"
		if a.Pixels < 0 {
			a.Pixels = -a.Pixels
			a.Direction = "down"
		}
	}
	return a
}
func delay(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
