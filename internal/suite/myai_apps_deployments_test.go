package suite

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestAIConnectorDeploymentReceipts(t *testing.T) {
	for _, mode := range []string{"accepted", "lost-response", "old-result-only", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			archive := []byte("Immutable fixture source")
			hash := sha256.Sum256(archive)
			sourceHash := hex.EncodeToString(hash[:])
			var mu sync.Mutex
			posts := 0
			jobs := []map[string]any{{"id": strings.Repeat("e", 32), "sourceHash": sourceHash, "environment": "preview", "created": time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339Nano), "state": "published"}}
			site := strings.Repeat("a", 32)
			native := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				c, e := r.Cookie("myhosting_session")
				if e != nil || c.Value != "fixture-native" || r.Header.Get("Authorization") != "" {
					t.Error("Wrong native authority")
					http.Error(w, "denied", 403)
					return
				}
				base := "/api/projects/" + site + "/deployments"
				if r.Method == "POST" && r.URL.Path == base {
					posts++
					body, _ := io.ReadAll(r.Body)
					if !bytes.Equal(body, archive) || r.Header.Get("X-CSRF-Token") != "fixture-csrf" || r.Header.Get("Content-Type") != "application/zip" {
						t.Error("Source/CSRF changed")
					}
					if mode == "rejected" && posts == 1 {
						w.WriteHeader(400)
						w.Write([]byte(`{"error":"Queue not ready"}`))
						return
					}
					if mode != "old-result-only" {
						jobs = append(jobs, map[string]any{"id": fmt.Sprintf("%032x", posts), "sourceHash": sourceHash, "environment": "preview", "created": time.Now().UTC().Format(time.RFC3339Nano), "state": "queued"})
					}
					if mode == "lost-response" || mode == "old-result-only" {
						w.WriteHeader(500)
						w.Write([]byte(`{"error":"Response unconfirmed"}`))
						return
					}
					w.WriteHeader(202)
					_ = json.NewEncoder(w).Encode(jobs[len(jobs)-1])
					return
				}
				if r.URL.Path == base {
					_ = json.NewEncoder(w).Encode(map[string]any{"deployments": jobs})
					return
				}
				for _, job := range jobs {
					if r.URL.Path == base+"/"+job["id"].(string) {
						_ = json.NewEncoder(w).Encode(job)
						return
					}
				}
				http.NotFound(w, r)
			}))
			defer native.Close()
			u, _ := url.Parse(native.URL)
			_, port, _ := net.SplitHostPort(u.Host)
			n, _ := strconv.Atoi(port)
			app := catalog.App{ID: "hosting", Name: "My Hosting", Port: n, Route: "/hosting", StripPrefix: true}
			p := &Portal{rt: &Runtime{Config: Config{Data: t.TempDir(), Origin: "https://workspace.example.test"}}}
			connection := aiAppConnection{Cookies: []*http.Cookie{{Name: "myhosting_session", Value: "fixture-native"}}}
			call := aiAppCall{App: "hosting", Operation: "deploy_project", Parameters: map[string]string{"site": site, "project": strings.Repeat("b", 48)}, RequestID: strings.Repeat("1", 64)}
			submit := func(c aiAppCall) (any, error) {
				return p.aiSubmitDeployment(context.Background(), strings.Repeat("a", 64), app, connection, "fixture-csrf", c, archive, strings.Repeat("c", 40), "preview", "https://workspace.example.test/git/private/project")
			}
			_, e := submit(call)
			if (mode == "accepted") != (e == nil) {
				t.Fatal("Wrong initial upload result", e)
			}
			retry := call
			if mode == "lost-response" {
				retry.RequestID = strings.Repeat("2", 64)
			}
			_, again := submit(retry)
			mu.Lock()
			count := posts
			mu.Unlock()
			if count != 1 {
				t.Fatal("One action submitted twice", count)
			}
			switch mode {
			case "accepted":
				if again != nil {
					t.Fatal(again)
				}
				call.RequestID = strings.Repeat("2", 64)
				if _, e = submit(call); e != nil {
					t.Fatal(e)
				}
				mu.Lock()
				count = posts
				mu.Unlock()
				if count != 2 {
					t.Fatal("A later deliberate redeployment was suppressed", count)
				}
			case "lost-response":
				if again != nil {
					t.Fatal("Existing accepted job was not reconciled", again)
				}
				if _, e = submit(retry); e != nil {
					t.Fatal("Recovered action alias was not retained", e)
				}
				call.RequestID = strings.Repeat("3", 64)
				if _, e = submit(call); e == nil {
					t.Fatal("Fixture lost response should remain unconfirmed on a new action")
				}
				mu.Lock()
				count = posts
				mu.Unlock()
				if count != 2 {
					t.Fatal("A new explicit action was not distinct", count)
				}
			case "old-result-only":
				if again == nil {
					t.Fatal("An old deployment was mistaken for the current upload")
				}
				call.RequestID = strings.Repeat("2", 64)
				if _, e = submit(call); e == nil {
					t.Fatal("A new action bypassed the unresolved upload")
				}
				mu.Lock()
				count = posts
				mu.Unlock()
				if count != 1 {
					t.Fatal("Unresolved source uploaded again", count)
				}
			case "rejected":
				if again == nil {
					t.Fatal("Rejected action retried without a new user turn")
				}
				call.RequestID = strings.Repeat("2", 64)
				if _, e = submit(call); e != nil {
					t.Fatal("A later reviewed retry of a definite rejection was blocked", e)
				}
				mu.Lock()
				count = posts
				mu.Unlock()
				if count != 2 {
					t.Fatal(count)
				}
			}
		})
	}
}
