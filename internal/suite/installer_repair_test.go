package suite

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
)

func TestInstallerRepairReportsMissingFailedAppsAndConnector(t *testing.T) {
	c := Config{Tunnel: true, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "accounts", Kind: "process"}, {ID: "tasks", Kind: "process"}, {ID: "docs", Kind: "builtin"}}}}
	response := func(body string) *http.Response {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}
	}
	r, e := readInstallerHealth(response(`{"status":"ok","apps":{"accounts":"running","tasks":"running","cloudflare":"running"}}`), c, "next")
	if e != nil || len(r.NeedsAttention) != 0 {
		t.Fatal(r, e)
	}
	r, e = readInstallerHealth(response(`{"status":"ok","apps":{"accounts":"running","tasks":"failed"}}`), c, "next")
	if e != nil || strings.Join(r.NeedsAttention, ",") != "cloudflare,tasks" {
		t.Fatal(r, e)
	}
	for _, bad := range []string{`{"status":"ok","apps":{}}`, `{"status":"failed","apps":{"accounts":"running"}}`, `invalid`} {
		if _, e = readInstallerHealth(response(bad), c, "next"); e == nil {
			t.Fatal("false healthy result", bad)
		}
	}
}

func TestInstallerProgressSurvivesSplitWritesAndIgnoresNonProtocolOutput(t *testing.T) {
	var phase, stage string
	p := &installerRepairOutput{ctx: operationprogress.Observe(context.Background(), func(e operationprogress.Event) { stage = e.Stage }), report: func(s string) { phase = s }}
	p.Write([]byte("not protocol\n{\"phase\":\"Back up"))
	p.Write([]byte(" your workspace\"}\n{\"progress\":{\"stage\":\"restart\",\"detail\":\"Checking apps\"}}\n"))
	if phase != "Back up your workspace" || stage != "restart" {
		t.Fatal(phase, stage)
	}
}

func TestInstallerRepairCannotRestartAnIncompleteSetup(t *testing.T) {
	if _, e := ReinstallInstalled(context.Background(), t.TempDir(), nil, catalog.Catalog{}, true, nil); e == nil {
		t.Fatal("incomplete data could reset accounts")
	}
}
