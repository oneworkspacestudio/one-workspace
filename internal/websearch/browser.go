package websearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Browser exposes only rendered screenshots and physical browser input. There
// is deliberately no DOM, accessibility, selector, JavaScript or shell tool.
type Browser struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	seq      int
	process  *exec.Cmd
	cancel   context.CancelFunc
	proxy    *http.Server
	target   string
	endpoint string
	done     chan struct{}
}
type pageInfo struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
	WS    string `json:"webSocketDebuggerUrl"`
}

func freePort() (int, error) {
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		return 0, e
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func StartBrowser(ctx context.Context, executable, profile string, proxy string) (*Browser, error) {
	if e := os.MkdirAll(profile, 0700); e != nil {
		return nil, e
	}
	// A unique profile prevents connection to any owner's existing browser.
	_ = os.Remove(filepath.Join(profile, "DevToolsActivePort"))
	child, cancel := context.WithCancel(ctx)
	args := []string{"--headless=new", "--user-data-dir=" + profile, "--remote-debugging-address=127.0.0.1", "--remote-debugging-port=0", "--window-size=1280,800", "--disk-cache-size=134217728", "--media-cache-size=16777216", "--no-first-run", "--no-default-browser-check", "--disable-extensions", "--disable-sync", "--disable-background-networking", "--disable-features=Translate,OptimizationHints,MediaRouter", "--force-webrtc-ip-handling-policy=disable_non_proxied_udp", "--disable-component-update", "--disable-quic", "about:blank"}
	if os.Getenv("WEBSEARCH_NATIVE_ACCEPTANCE") == "1" {
		args = append(args, "--enable-logging=stderr", "--v=1")
	}
	if proxy != "" {
		args = append(args, "--proxy-server="+proxy, "--proxy-bypass-list=<-loopback>")
	}
	cmd := exec.CommandContext(child, executable, args...)
	cmd.Env = childEnvironment()
	hideProcess(cmd)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if os.Getenv("WEBSEARCH_NATIVE_ACCEPTANCE") == "1" {
		if log, e := os.Create(filepath.Join(profile, "startup.log")); e == nil {
			cmd.Stderr = log
			defer log.Close()
		}
	}
	if e := cmd.Start(); e != nil {
		cancel()
		return nil, e
	}
	b := &Browser{process: cmd, cancel: cancel, done: make(chan struct{})}
	go func() { _ = cmd.Wait(); close(b.done) }()
	fail := func(e error) (*Browser, error) { b.Close(); return nil, e }
	var endpoint string
	for n := 0; n < 600; n++ {
		raw, e := os.ReadFile(filepath.Join(profile, "DevToolsActivePort"))
		if e == nil {
			p := strings.Split(string(raw), "\n")
			if len(p) > 1 {
				endpoint = "http://127.0.0.1:" + strings.TrimSpace(p[0])
				break
			}
		}
		select {
		case <-b.done:
			return fail(errors.New("portable browser exited before readiness"))
		case <-child.Done():
			return fail(child.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
	if endpoint == "" {
		return fail(errors.New("portable browser did not become ready"))
	}
	b.endpoint = endpoint
	request, _ := http.NewRequestWithContext(ctx, "GET", endpoint+"/json/list", nil)
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{Proxy: nil}}
	res, e := client.Do(request)
	if e != nil {
		return fail(e)
	}
	var pages []pageInfo
	e = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&pages)
	res.Body.Close()
	if e != nil {
		return fail(e)
	}
	for _, p := range pages {
		if p.Type == "page" {
			b.target = p.ID
			conn, _, err := websocket.DefaultDialer.DialContext(ctx, p.WS, nil)
			if err != nil {
				return fail(err)
			}
			b.conn = conn
			break
		}
	}
	if b.conn == nil {
		return fail(errors.New("browser page unavailable"))
	}
	if e = b.Call(ctx, "Page.enable", nil, nil); e != nil {
		return fail(e)
	}
	_ = b.Call(ctx, "Browser.setDownloadBehavior", map[string]any{"behavior": "deny"}, nil)
	if e = b.Call(ctx, "Emulation.setDeviceMetricsOverride", map[string]any{"width": 1280, "height": 800, "deviceScaleFactor": 1.5, "mobile": false}, nil); e != nil {
		return fail(e)
	}
	return b, nil
}

func (b *Browser) Close() {
	b.mu.Lock()
	conn := b.conn
	b.mu.Unlock()
	graceful := false
	if conn != nil {
		ctx, stop := context.WithTimeout(context.Background(), 2*time.Second)
		e := b.Call(ctx, "Browser.close", nil, nil)
		if e != nil && b.endpoint != "" {
			if b.Reconnect(ctx) == nil {
				e = b.Call(ctx, "Browser.close", nil, nil)
			}
		}
		graceful = e == nil
		stop()
		_ = conn.Close()
	}
	// Let Chromium flush cookies and profile state before a forced fallback.
	if graceful && b.done != nil {
		select {
		case <-b.done:
		case <-time.After(3 * time.Second):
		}
	}
	if b.cancel != nil {
		b.cancel()
	}
	if b.done != nil {
		select {
		case <-b.done:
		case <-time.After(3 * time.Second):
		}
	}
	if b.proxy != nil {
		_ = b.proxy.Close()
	}
}
func (b *Browser) Call(ctx context.Context, method string, params any, result any) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	conn := b.conn
	if conn == nil {
		return errors.New("browser connection unavailable")
	}
	if e := ctx.Err(); e != nil {
		return e
	}
	deadline := time.Now().Add(60 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	_ = conn.SetWriteDeadline(deadline)
	_ = conn.SetReadDeadline(deadline)
	stopCancellation := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stopCancellation()
	b.seq++
	id := b.seq
	if e := conn.WriteJSON(map[string]any{"id": id, "method": method, "params": params}); e != nil {
		return e
	}
	for {
		var v struct {
			ID     int             `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if e := conn.ReadJSON(&v); e != nil {
			return e
		}
		if v.ID != id {
			continue
		}
		if v.Error != nil {
			return errors.New("browser control: " + v.Error.Message)
		}
		if result != nil {
			return json.Unmarshal(v.Result, result)
		}
		return nil
	}
}
func (b *Browser) Screenshot(ctx context.Context) (string, error) {
	var v struct {
		Data string `json:"data"`
	}
	e := b.Call(ctx, "Page.captureScreenshot", map[string]any{"format": "png"}, &v)
	return v.Data, e
}
func (b *Browser) Location(ctx context.Context) (Source, error) {
	b.mu.Lock()
	target := b.target
	b.mu.Unlock()
	var v struct {
		Target pageInfo `json:"targetInfo"`
	}
	e := b.Call(ctx, "Target.getTargetInfo", map[string]string{"targetId": target}, &v)
	return Source{URL: v.Target.URL, Title: v.Target.Title}, e
}
func (b *Browser) Navigate(ctx context.Context, raw string) error {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil {
		return errors.New("only public HTTPS pages are supported")
	}
	return b.Call(ctx, "Page.navigate", map[string]string{"url": u.String()}, nil)
}
func (b *Browser) Key(ctx context.Context, keys []string) error {
	modifiers := 0
	key := ""
	code := ""
	vk := 0
	for _, k := range keys {
		switch strings.ToLower(k) {
		case "ctrl", "control":
			modifiers |= 2
		case "alt":
			modifiers |= 1
		case "shift":
			modifiers |= 8
		case "enter", "return":
			key = "Enter"
			code = "Enter"
			vk = 13
		case "tab":
			key = "Tab"
			vk = 9
		case "escape", "esc":
			key = "Escape"
			vk = 27
		case "backspace":
			key = "Backspace"
			vk = 8
		case "delete":
			key = "Delete"
			vk = 46
		case "arrowleft", "left":
			key = "ArrowLeft"
			vk = 37
		case "arrowright", "right":
			key = "ArrowRight"
			vk = 39
		case "home":
			key = "Home"
			vk = 36
		case "end":
			key = "End"
			vk = 35
		case "arrowdown", "down":
			key = "ArrowDown"
			vk = 40
		case "arrowup", "up":
			key = "ArrowUp"
			vk = 38
		case "pageup":
			key = "PageUp"
			vk = 33
		case "pagedown":
			key = "PageDown"
			vk = 34
		case "a", "c", "l":
			key = strings.ToLower(k)
			code = "Key" + strings.ToUpper(k)
			vk = int(strings.ToUpper(k)[0])
		default:
			return fmt.Errorf("unsupported browser key")
		}
	}
	if key == "" {
		return errors.New("missing browser key")
	}
	// Content input only: address-bar, developer-tools and OS/clipboard shortcuts
	// must never let the pilot escape the checked HTTPS navigation primitive.
	contentKey := key == "ArrowLeft" || key == "ArrowRight" || key == "ArrowUp" || key == "ArrowDown" || key == "Home" || key == "End" || key == "PageUp" || key == "PageDown" || key == "Backspace" || key == "Delete"
	if modifiers != 0 && !(modifiers == 2 && (key == "a" || contentKey)) && !(modifiers == 8 && (key == "Tab" || contentKey)) && !(modifiers == 10 && contentKey) {
		return errors.New("browser chrome shortcuts are unavailable; use visit_url for HTTPS navigation")
	}
	p := map[string]any{"type": "rawKeyDown", "key": key, "code": code, "modifiers": modifiers, "windowsVirtualKeyCode": vk}
	if e := b.Call(ctx, "Input.dispatchKeyEvent", p, nil); e != nil {
		return e
	}
	p["type"] = "keyUp"
	return b.Call(ctx, "Input.dispatchKeyEvent", p, nil)
}
func (b *Browser) Act(ctx context.Context, a Action) error {
	switch a.Action {
	case "visit_url":
		return b.Navigate(ctx, a.URL)
	case "web_search":
		return b.Navigate(ctx, "https://www.bing.com/search?q="+url.QueryEscape(a.Query))
	case "history_back":
		var history struct {
			Current int `json:"currentIndex"`
			Entries []struct {
				ID  int    `json:"id"`
				URL string `json:"url"`
			} `json:"entries"`
		}
		if e := b.Call(ctx, "Page.getNavigationHistory", nil, &history); e != nil {
			return e
		}
		if history.Current < 1 {
			return nil
		}
		entry := history.Entries[history.Current-1]
		if !strings.HasPrefix(entry.URL, "https://") {
			return errors.New("history destination is not public HTTPS")
		}
		return b.Call(ctx, "Page.navigateToHistoryEntry", map[string]int{"entryId": entry.ID}, nil)
	case "type":
		if len(a.Text) > 8000 {
			return errors.New("text limit exceeded")
		}
		return b.Call(ctx, "Input.insertText", map[string]string{"text": a.Text}, nil)
	case "key":
		return b.Key(ctx, a.Keys)
	case "left_click", "double_click", "mouse_move":
		if len(a.Coordinate) != 2 || a.Coordinate[0] < 0 || a.Coordinate[0] >= 1280 || a.Coordinate[1] < 0 || a.Coordinate[1] >= 800 {
			return errors.New("coordinates outside viewport")
		}
		p := map[string]any{"x": a.Coordinate[0], "y": a.Coordinate[1], "type": "mouseMoved"}
		if a.Action == "mouse_move" {
			return b.Call(ctx, "Input.dispatchMouseEvent", p, nil)
		}
		p["button"] = "left"
		p["clickCount"] = 1
		if a.Action == "double_click" {
			p["clickCount"] = 2
		}
		p["type"] = "mousePressed"
		if e := b.Call(ctx, "Input.dispatchMouseEvent", p, nil); e != nil {
			return e
		}
		p["type"] = "mouseReleased"
		return b.Call(ctx, "Input.dispatchMouseEvent", p, nil)
	case "scroll":
		pixels := a.Pixels
		if pixels < 0 || pixels > 1600 {
			return errors.New("scroll exceeds limit")
		}
		if pixels == 0 {
			pixels = 600
		}
		if a.Direction == "up" {
			pixels = -pixels
		} else if a.Direction != "down" && a.Direction != "" {
			return errors.New("unsupported scroll direction")
		}
		return b.Call(ctx, "Input.dispatchMouseEvent", map[string]any{"type": "mouseWheel", "x": 1000, "y": 650, "deltaY": pixels, "deltaX": 0}, nil)
	case "wait":
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			return nil
		}
	default:
		return errors.New("unsupported browser action")
	}
}
