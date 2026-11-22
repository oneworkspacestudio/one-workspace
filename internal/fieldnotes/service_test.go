package fieldnotes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

type fixtureTransport func(*http.Request) (*http.Response, error)

func (f fixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestResearchPersistsCitedBriefingAndProtectsOwnerLibrary(t *testing.T) {
	const origin = "http://localhost:29760"
	host, err := apphost.New(apphost.Config{
		ID: "fieldnotes", Name: "AI Fieldnotes", Version: Version, Data: t.TempDir(),
		Listen: "127.0.0.1:29760", Origin: origin,
	}, UI())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { host.Close() })
	if err := host.Bootstrap("reader@example.invalid", "synthetic-reading-passphrase"); err != nil {
		t.Fatal(err)
	}
	app, err := Register(host)
	if err != nil {
		t.Fatal(err)
	}
	s := app.(*service)
	var calls int
	s.client.Transport = fixtureTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Authorization") != "Bearer fixture-deepseek-token" || r.URL.Hostname() != "api.deepseek.com" {
			t.Fatal("provider credential or destination is incorrect")
		}
		var body string
		switch r.URL.Path {
		case "/v1/tools/search":
			body = `{"results":[{"title":"Published engineering paper","url":"https://example.org/paper","content":"Fixture evidence long enough to meet the candidate selection threshold."},{"title":"Injected URL","url":"javascript:alert(1)","content":"Untrusted provider result with an unsafe scheme."}]}`
		case "/chat/completions":
			input, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(input), "Fixture evidence") || strings.Contains(string(input), "fixture-deepseek-token") {
				t.Fatal("draft input missing evidence or contains private credential")
			}
			content := `{"title":"A grounded briefing","summary":"A private test draft","sections":[{"heading":"Evidence","body":"A source-backed section.","refs":[1]},{"heading":"UI implications","body":"Review the proposed changes.","refs":[]},{"heading":"Pilot","body":"Measure review time.","refs":[]}]}`
			encoded, _ := json.Marshal(content)
			body = `{"choices":[{"finish_reason":"stop","message":{"content":` + string(encoded) + `}}]}`
		default:
			t.Fatal("unexpected production API call")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	handler := host.Handler()
	var cookies []*http.Cookie
	var csrf string
	request := func(method, path, body string, status int) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, origin+path, strings.NewReader(body))
		r.RemoteAddr = "127.0.0.1:19001"
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", origin)
		r.Header.Set("X-CSRF-Token", csrf)
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != status {
			t.Fatalf("%s %s: got %d, want %d (%s)", method, path, w.Code, status, w.Body.String())
		}
		return w
	}
	request("GET", "/api/fieldnotes", "", 401)
	login := request("POST", "/api/login", `{"email":"reader@example.invalid","password":"synthetic-reading-passphrase"}`, 200)
	cookies = login.Result().Cookies()
	var session apphost.Session
	if err := json.Unmarshal(request("GET", "/api/me", "", 200).Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	request("PUT", "/api/providers", `{"provider":"deepseek","key":"fixture-deepseek-token"}`, 403)
	csrf = session.CSRF
	request("PUT", "/api/providers", `{"provider":"deepseek","key":"fixture-deepseek-token"}`, 200)
	var searched map[string]string
	search := request("POST", "/api/research", `{"query":"Loop engineering","provider":"deepseek"}`, 202)
	if err := json.Unmarshal(search.Body.Bytes(), &searched); err != nil {
		t.Fatal(err)
	}
	request("POST", "/api/research", `{"query":"Duplicate","provider":"deepseek"}`, 409)
	if err := s.work(context.Background()); err != nil {
		t.Fatal(err)
	}
	prefix := `{"search":"` + searched["id"] + `","indices":`
	request("POST", "/api/research", prefix+`[1]}`, 400)
	want := request("POST", "/api/research", prefix+`[0]}`, 202)
	var generated map[string]string
	if err := json.Unmarshal(want.Body.Bytes(), &generated); err != nil {
		t.Fatal(err)
	}
	if err := s.work(context.Background()); err != nil {
		t.Fatal(err)
	}
	request("PUT", "/api/saved", `{"id":"`+generated["id"]+`","saved":true}`, 200)
	state := request("GET", "/api/fieldnotes", "", 200).Body.String()
	if calls != 2 || !strings.Contains(state, "A grounded briefing") || !strings.Contains(state, "https://example.org/paper") ||
		strings.Contains(state, "fixture-deepseek-token") || strings.Contains(state, "javascript:") {
		t.Fatal("private library is missing the brief or included invalid provider data")
	}
	if _, articles, _, err := s.readLibrary(context.Background(), "other-owner"); err != nil || len(articles) != 0 {
		t.Fatal("research crossed account boundaries")
	}
	request("POST", "/auth/logout", "", 200)
	request("GET", "/api/fieldnotes", "", 401)
}
