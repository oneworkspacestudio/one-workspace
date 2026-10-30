package suite

import (
	"bufio"
	"context"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestMfaRevocationClosesExistingProxyStreams(t *testing.T) {
	for _, upgrade := range []bool{false, true} {
		t.Run(map[bool]string{false: "event-stream", true: "websocket"}[upgrade], func(t *testing.T) {
			var active atomic.Bool
			active.Store(true)
			provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if active.Load() {
					io.WriteString(w, `{"active":true,"sub":"member"}`)
				} else {
					io.WriteString(w, `{"active":false}`)
				}
			}))
			defer provider.Close()
			upstreamClosed := make(chan struct{}, 1)
			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if upgrade {
					c, _, e := http.NewResponseController(w).Hijack()
					if e != nil {
						return
					}
					defer c.Close()
					io.WriteString(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: websocket\r\n\r\n")
					b := make([]byte, 1)
					_, _ = c.Read(b)
					upstreamClosed <- struct{}{}
					return
				}
				w.Header().Set("Content-Type", "text/event-stream")
				io.WriteString(w, "event: ready\ndata: 1\n\n")
				http.NewResponseController(w).Flush()
				<-r.Context().Done()
				upstreamClosed <- struct{}{}
			}))
			defer backend.Close()
			target, _ := url.Parse(backend.URL)
			proxy := httputil.NewSingleHostReverseProxy(target)
			p := &Portal{rt: &Runtime{Config: Config{Issuer: provider.URL}}, oauth: oauth2.Config{ClientID: "client", ClientSecret: "synthetic"}, sessions: map[string]session{"valid": {Subject: "member", AccessToken: "synthetic", Expires: time.Now().Add(time.Hour)}}}
			gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := p.user(r); !ok {
					http.Error(w, "unauthorized", 401)
					return
				}
				ctx, cancel := context.WithCancel(r.Context())
				defer cancel()
				r = r.WithContext(ctx)
				done := make(chan struct{})
				defer close(done)
				go p.watchAccountSession(r, cancel, done, 10*time.Millisecond)
				proxy.ServeHTTP(w, r)
			}))
			defer gateway.Close()
			if upgrade {
				address := strings.TrimPrefix(gateway.URL, "http://")
				c, e := net.DialTimeout("tcp", address, time.Second)
				if e != nil {
					t.Fatal(e)
				}
				defer c.Close()
				c.SetDeadline(time.Now().Add(2 * time.Second))
				io.WriteString(c, "GET / HTTP/1.1\r\nHost: "+address+"\r\nCookie: workspace_session=valid\r\nConnection: Upgrade\r\nUpgrade: websocket\r\n\r\n")
				reader := bufio.NewReader(c)
				line, e := reader.ReadString('\n')
				if e != nil || !strings.Contains(line, "101") {
					t.Fatal("upgrade failed", line, e)
				}
				for {
					line, e = reader.ReadString('\n')
					if e != nil {
						t.Fatal(e)
					}
					if line == "\r\n" {
						break
					}
				}
				active.Store(false)
				if _, e = reader.ReadByte(); e == nil {
					t.Fatal("revoked socket remained readable")
				}
			} else {
				request, _ := http.NewRequest("GET", gateway.URL, nil)
				request.AddCookie(&http.Cookie{Name: "workspace_session", Value: "valid"})
				response, e := http.DefaultClient.Do(request)
				if e != nil {
					t.Fatal(e)
				}
				defer response.Body.Close()
				active.Store(false)
				_, _ = io.ReadAll(response.Body)
			}
			select {
			case <-upstreamClosed:
			case <-time.After(2 * time.Second):
				t.Fatal("revoked upstream remained connected")
			}
		})
	}
}
