package websearch

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type BrowserTab struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

func (b *Browser) pages(ctx context.Context) ([]pageInfo, error) {
	req, e := http.NewRequestWithContext(ctx, "GET", b.endpoint+"/json/list", nil)
	if e != nil {
		return nil, e
	}
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{Proxy: nil}}
	res, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("browser unavailable")
	}
	var pages []pageInfo
	e = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&pages)
	return pages, e
}

// Reconnect keeps the actual page (including an unfinished sign-in form) alive
// when cancelling an LLM request has interrupted its CDP connection.
func (b *Browser) Reconnect(ctx context.Context) error {
	b.mu.Lock()
	target := b.target
	b.mu.Unlock()
	return b.selectTab(ctx, target)
}
func (b *Browser) selectTab(ctx context.Context, id string) error {
	pages, e := b.pages(ctx)
	if e != nil {
		return e
	}
	var chosen *pageInfo
	for i := range pages {
		if pages[i].Type == "page" && pages[i].ID == id {
			chosen = &pages[i]
			break
		}
	}
	if chosen == nil {
		return errors.New("browser tab unavailable")
	}
	conn, _, e := websocket.DefaultDialer.DialContext(ctx, chosen.WS, nil)
	if e != nil {
		return e
	}
	conn.SetReadLimit(16 << 20)
	b.mu.Lock()
	old := b.conn
	b.conn = conn
	b.target = chosen.ID
	b.mu.Unlock()
	if old != nil {
		_ = old.Close()
	}
	if e = b.Call(ctx, "Page.enable", nil, nil); e != nil {
		return e
	}
	if e = b.Call(ctx, "Target.activateTarget", map[string]string{"targetId": chosen.ID}, nil); e != nil {
		return e
	}
	return b.Call(ctx, "Emulation.setDeviceMetricsOverride", map[string]any{"width": 1280, "height": 800, "deviceScaleFactor": 1.5, "mobile": false}, nil)
}
func (b *Browser) Tabs(ctx context.Context) ([]BrowserTab, error) {
	pages, e := b.pages(ctx)
	if e != nil {
		return nil, e
	}
	b.mu.Lock()
	target := b.target
	b.mu.Unlock()
	tabs := []BrowserTab{}
	for _, p := range pages {
		if p.Type == "page" {
			tabs = append(tabs, BrowserTab{ID: p.ID, Title: limitText(p.Title, 256), URL: limitText(p.URL, 4096), Active: p.ID == target})
		}
	}
	return tabs, nil
}
func limitText(s string, n int) string {
	if len(s) > n {
		return strings.ToValidUTF8(s[:n], "")
	}
	return s
}

func (b *Browser) DesktopFrame(ctx context.Context) (string, Source, []BrowserTab, error) {
	var shot struct {
		Data string `json:"data"`
	}
	if e := b.Call(ctx, "Page.captureScreenshot", map[string]any{"format": "jpeg", "quality": 75, "captureBeyondViewport": false}, &shot); e != nil {
		return "", Source{}, nil, e
	}
	loc, e := b.Location(ctx)
	if e != nil {
		return "", loc, nil, e
	}
	tabs, e := b.Tabs(ctx)
	loc.URL = limitText(loc.URL, 4096)
	loc.Title = limitText(loc.Title, 256)
	return "data:image/jpeg;base64," + shot.Data, loc, tabs, e
}

type DesktopInput struct {
	Type    string   `json:"type"`
	X       float64  `json:"x"`
	Y       float64  `json:"y"`
	Buttons int      `json:"buttons"`
	Count   int      `json:"count"`
	DeltaX  float64  `json:"deltaX"`
	DeltaY  float64  `json:"deltaY"`
	Text    string   `json:"text"`
	Keys    []string `json:"keys"`
	URL     string   `json:"url"`
	Tab     string   `json:"tab"`
}

func (in DesktopInput) valid() bool {
	switch in.Type {
	case "pointerDown", "pointerUp", "pointerMove":
		return !math.IsNaN(in.X) && !math.IsNaN(in.Y) && in.X >= 0 && in.X < 1280 && in.Y >= 0 && in.Y < 800 && (in.Buttons == 0 || in.Buttons == 1) && in.Count >= 0 && in.Count <= 2
	case "wheel":
		return !math.IsNaN(in.DeltaX) && !math.IsNaN(in.DeltaY) && math.Abs(in.DeltaX) <= 1600 && math.Abs(in.DeltaY) <= 1600 && in.X >= 0 && in.X < 1280 && in.Y >= 0 && in.Y < 800
	case "text":
		return len(in.Text) > 0 && len(in.Text) <= 8000
	case "key":
		return len(in.Keys) > 0 && len(in.Keys) <= 3
	case "navigate":
		return validDesktopURL(in.URL)
	case "newTab":
		return in.URL == "" || validDesktopURL(in.URL)
	case "closeTab":
		return len(in.Tab) <= 80
	case "selectTab":
		return len(in.Tab) > 0 && len(in.Tab) <= 80
	case "back", "forward", "reload":
		return true
	}
	return false
}
func validDesktopURL(raw string) bool {
	u, e := url.Parse(raw)
	return e == nil && len(raw) <= 4096 && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && (u.Port() == "" || u.Port() == "443")
}
func (b *Browser) DesktopInput(ctx context.Context, in DesktopInput) error {
	if !in.valid() {
		return errors.New("unsupported browser input")
	}
	switch in.Type {
	case "navigate":
		return b.Navigate(ctx, in.URL)
	case "reload":
		return b.Call(ctx, "Page.reload", nil, nil)
	case "text":
		return b.Call(ctx, "Input.insertText", map[string]string{"text": in.Text}, nil)
	case "key":
		return b.Key(ctx, in.Keys)
	case "pointerDown", "pointerUp", "pointerMove":
		kind := "mouseMoved"
		button := "none"
		count := 0
		if in.Type == "pointerDown" {
			kind = "mousePressed"
			button = "left"
			count = 1
		}
		if in.Type == "pointerUp" {
			kind = "mouseReleased"
			button = "left"
			count = 1
		}
		if count > 0 && in.Count > 0 {
			count = in.Count
		}
		return b.Call(ctx, "Input.dispatchMouseEvent", map[string]any{"type": kind, "x": in.X, "y": in.Y, "button": button, "buttons": in.Buttons, "clickCount": count}, nil)
	case "wheel":
		return b.Call(ctx, "Input.dispatchMouseEvent", map[string]any{"type": "mouseWheel", "x": in.X, "y": in.Y, "deltaX": in.DeltaX, "deltaY": in.DeltaY}, nil)
	case "back", "forward":
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
		index := history.Current - 1
		if in.Type == "forward" {
			index = history.Current + 1
		}
		if index < 0 || index >= len(history.Entries) {
			return nil
		}
		entry := history.Entries[index]
		if entry.URL != "about:blank" && !validDesktopURL(entry.URL) {
			return errors.New("HTTPS navigation required")
		}
		return b.Call(ctx, "Page.navigateToHistoryEntry", map[string]int{"entryId": entry.ID}, nil)
	case "newTab":
		address := in.URL
		if address == "" {
			address = "about:blank"
		}
		var created struct {
			ID string `json:"targetId"`
		}
		if e := b.Call(ctx, "Target.createTarget", map[string]any{"url": address}, &created); e != nil {
			return e
		}
		return b.selectTab(ctx, created.ID)
	case "selectTab":
		return b.selectTab(ctx, in.Tab)
	case "closeTab":
		b.mu.Lock()
		current := b.target
		b.mu.Unlock()
		if in.Tab == "" {
			in.Tab = current
		}
		pages, e := b.pages(ctx)
		if e != nil {
			return e
		}
		found := false
		other := ""
		for _, p := range pages {
			if p.Type != "page" {
				continue
			}
			if p.ID == in.Tab {
				found = true
			} else {
				other = p.ID
			}
		}
		if !found {
			return errors.New("browser tab unavailable")
		}
		if in.Tab != current {
			other = current
		}
		if other == "" {
			if e := b.DesktopInput(ctx, DesktopInput{Type: "newTab"}); e != nil {
				return e
			}
		} else {
			if e := b.selectTab(ctx, other); e != nil {
				return e
			}
		}
		return b.Call(ctx, "Target.closeTarget", map[string]string{"targetId": in.Tab}, nil)
	}
	return errors.New("unsupported browser input")
}
