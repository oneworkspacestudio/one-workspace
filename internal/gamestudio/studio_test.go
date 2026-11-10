package gamestudio

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func TestModelOutputIsRestrictedToTemplates(t *testing.T) {
	valid := Spec{"Star Cat", "catch", "cat", "star", "space", 4, 1}
	for _, tc := range []struct {
		name  string
		spec  Spec
		valid bool
	}{
		{"expected", valid, true},
		{"untrusted template", Spec{"Star Cat", "../secret", "cat", "star", "space", 4, 1}, false},
		{"invalid bounds", Spec{"Star Cat", "catch", "cat", "star", "space", 200, 1}, false},
		{"contact details", Spec{"name@example.test", "catch", "cat", "star", "space", 4, 1}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.spec.valid() != tc.valid {
				t.Fatal("unexpected validation result")
			}
		})
	}
	page, err := render(Spec{"<script>alert(1)</script>", "catch", "cat", "star", "space", 4, 1})
	if err != nil || strings.Contains(string(page), "<script>alert(1)</script>") {
		t.Fatal("game title was not escaped", err)
	}
}

// Run only when an external Playwright test owns the fixture lifetime.
func TestBrowserFixture(t *testing.T) {
	root := os.Getenv("GAME_STUDIO_TEST_FIXTURE")
	if root == "" {
		t.Skip("browser fixture is opt-in")
	}
	var mu sync.Mutex
	files := map[string]string{}
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/model":
			var in struct{ Messages []struct{ Content string } }
			if json.NewDecoder(r.Body).Decode(&in) != nil || len(in.Messages) != 2 {
				w.WriteHeader(400)
				return
			}
			spec := Spec{"Cosmic Frog", "catch", "frog", "star", "space", 4, 1}
			prompt := in.Messages[1].Content
			switch {
			case strings.Contains(prompt, "maze"):
				spec = Spec{"Forest Maze", "maze", "frog", "gem", "forest", 4, 1}
			case strings.Contains(prompt, "runner"):
				spec = Spec{"Dragon Dash", "runner", "dragon", "cookie", "candy", 4, 1}
			case strings.Contains(prompt, "matching"):
				spec = Spec{"Ocean Pairs", "memory", "owl", "moon", "ocean", 4, 1}
			case strings.Contains(prompt, "blue"):
				spec.Theme = "ocean"
			}
			content, _ := json.Marshal(plan{"Your game is ready!", &spec})
			apphost.JSON(w, 200, map[string]any{"choices": []any{map[string]any{"finish_reason": "stop", "message": map[string]string{"content": string(content)}}}})
		case r.Header.Get("Authorization") != "Bearer synthetic-token":
			w.WriteHeader(401)
		case r.URL.Path == "/repos/test/builds":
			apphost.JSON(w, 200, map[string]any{"private": true, "empty": false, "default_branch": "main"})
		case strings.HasPrefix(r.URL.Path, "/repos/test/builds/contents/games/"):
			mu.Lock()
			defer mu.Unlock()
			if r.Method == http.MethodPost {
				if files[r.URL.Path] != "" {
					w.WriteHeader(409)
					return
				}
				var in map[string]string
				if json.NewDecoder(r.Body).Decode(&in) != nil || in["branch"] != "main" {
					w.WriteHeader(400)
					return
				}
				data, err := base64.StdEncoding.DecodeString(in["content"])
				if err != nil || !strings.Contains(string(data), "Cosmic Frog") {
					w.WriteHeader(400)
					return
				}
				files[r.URL.Path] = in["content"]
				w.WriteHeader(201)
				return
			}
			apphost.JSON(w, 200, map[string]string{"content": files[r.URL.Path]})
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	origin := "http://" + listener.Addr().String()
	h, err := apphost.New(apphost.Config{ID: "gamestudio", Name: "Game Studio", Version: Version, Origin: origin + "/gamestudio", Listen: listener.Addr().String(), Data: filepath.Join(root, "data")}, UI())
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	s, err := New(h)
	if err != nil {
		t.Fatal(err)
	}
	s.localURL = upstream.URL + "/model"
	s.gitURL = upstream.URL
	if err = h.Bootstrap("owner@example.test", "SyntheticOwner123"); err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: h.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go server.Serve(listener)
	defer server.Close()
	if err = os.WriteFile(filepath.Join(root, "origin.txt"), []byte(origin), 0600); err != nil {
		t.Fatal(err)
	}
	fmt.Println("Synthetic Studio fixture ready")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for {
		if _, err := os.Stat(filepath.Join(root, "stop")); err == nil {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatal("external browser fixture timed out")
		case <-time.After(100 * time.Millisecond):
		}
	}
}
