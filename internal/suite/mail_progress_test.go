package suite

import (
	"context"
	"encoding/json"
	"github.com/fahimc/one-workspace/internal/catalog"
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
				json.NewEncoder(w).Encode(map[string]any{"configured": true, "hostingBusy": false, "inboundState": tt.network, "inboundMessage": "Network needs guidance"})
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
		})
	}
}
