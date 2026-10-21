package suite

import (
	"context"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestFirstRunMailNetworkProgress(t *testing.T) {
	for _, tt := range []struct{ network, want string }{{"reachable", "done"}, {"attention", "warning"}, {"unverified", "warning"}, {"", "warning"}} {
		t.Run(tt.network, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]any{"configured": true, "hostingBusy": false, "inboundState": tt.network, "inboundMessage": "Router 192.168.1.1 TCP25 SMTP ISP diagnostics"})
			}))
			defer server.Close()
			u, _ := url.Parse(server.URL)
			port, _ := strconv.Atoi(u.Port())
			rt := &Runtime{progress: NewProgress("starting", []ProgressStep{{ID: "mail", Label: "Mail"}}, "")}
			rt.mailHostingProgress(context.Background(), catalog.App{ID: "mail", Port: port})
			step := rt.progress.Snapshot().Steps[0]
			if step.State != tt.want {
				t.Fatalf("state=%s want %s", step.State, tt.want)
			}
			if tt.want == "warning" && !strings.Contains(step.Detail, "Mail") {
				t.Fatal("missing recovery guidance")
			}
			for _, jargon := range []string{"192.168.1.1", "TCP25", "SMTP", "ISP"} {
				if strings.Contains(step.Detail, jargon) {
					t.Fatalf("technical diagnostic exposed in startup: %s", step.Detail)
				}
			}
		})
	}
}

func TestMailStartupKeepsDiagnosticsOutOfUserMessage(t *testing.T) {
	for _, state := range []map[string]any{
		{"identityError": "Mail's ordinary User role is missing or changed. Restore the standard User role before repairing sign-in"},
		{"hostingError": "Mail's ordinary User role is missing or changed. Restore the standard User role before repairing sign-in"},
		{"configured": true, "hostingError": "x:Role/set: invalidProperties"},
		{"hostingError": "Existing DNS conflict", "hostingIssue": map[string]string{"code": "existing_mail_dns"}},
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _ = json.NewEncoder(w).Encode(state) }))
		u, _ := url.Parse(server.URL)
		port, _ := strconv.Atoi(u.Port())
		rt := &Runtime{progress: NewProgress("starting", []ProgressStep{{ID: "mail", Label: "Mail"}}, "")}
		rt.mailHostingProgress(context.Background(), catalog.App{ID: "mail", Port: port})
		step := rt.progress.Snapshot().Steps[0]
		server.Close()
		if step.State != "warning" || !strings.Contains(step.Detail, "Open Mail") || strings.Contains(step.Detail, "User role") || strings.Contains(step.Detail, "x:Role") {
			t.Fatalf("unsafe or unhelpful startup warning: %+v", step)
		}
	}
}
