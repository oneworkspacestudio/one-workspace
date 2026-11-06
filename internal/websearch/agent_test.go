package websearch

import (
	"context"
	"encoding/base64"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestModelActionsDoNotExecuteCode(t *testing.T) {
	for _, raw := range []string{`__import__('os').system('calc')`, `{"name":"shell","arguments":{"action":"wait"}}`, `{"action":"evaluate","text":"document.cookie"}`, `{"action":"download","url":"https://example.com"}`} {
		if _, e := ParseAction(raw); e == nil {
			t.Fatalf("unsafe action accepted: %s", raw)
		}
	}
	a, e := ParseAction(`<tool_call>{"name":"computer_use","arguments":{"action":"left_click","coordinate":[40,90]}}</tool_call>`)
	if e != nil || a.Action != "left_click" || a.Coordinate[1] != 90 {
		t.Fatal(a, e)
	}
}
func TestReservedEgress(t *testing.T) {
	for _, s := range []string{"100.64.0.1", "192.0.2.1", "198.18.0.1", "64:ff9b::a00:1"} {
		if !inReserved(netip.MustParseAddr(s)) {
			t.Fatal("reserved address accepted", s)
		}
	}
	for _, host := range []string{"localhost", "127.0.0.1", "::1", "169.254.169.254", "10.0.0.1", "192.168.1.1"} {
		if _, e := publicIP(context.Background(), host); e == nil {
			t.Fatal("private address accepted", host)
		}
	}
}
func TestFaraNativeDisplayContract(t *testing.T) {
	a := browserAction(Action{Action: "left_click", Coordinate: []float64{500, 500}})
	if a.Coordinate[0] != 640 || a.Coordinate[1] != 400 {
		t.Fatal("normalized center did not reach browser center", a)
	}
	a = browserAction(Action{Action: "scroll", Pixels: -1000})
	if a.Direction != "down" || a.Pixels != 800 {
		t.Fatal("Fara signed scroll contract", a)
	}
}
func TestPortableRealBrowserAndModel(t *testing.T) {
	if os.Getenv("WEBSEARCH_NATIVE_ACCEPTANCE") != "1" {
		t.Skip("explicit native runtime acceptance only")
	}
	if e := OwnChildren(); e != nil {
		t.Fatal(e)
	}
	root := os.Getenv("WEBSEARCH_TEST_ROOT")
	if root == "" {
		t.Fatal("dedicated NTFS fixture required")
	}
	data, e := os.MkdirTemp(root, "agent-")
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	runtimeRoot := os.Getenv("WEBSEARCH_RUNTIME_ROOT")
	modelRoot := os.Getenv("WEBSEARCH_MODELS_ROOT")
	if modelRoot == "" {
		modelRoot = filepath.Join(runtimeRoot, "models")
	}
	a, e := StartAgent(ctx, runtimeRoot, modelRoot, data)
	if e != nil {
		t.Fatal(e)
	}
	defer a.Close()
	a.observe = func(action Action, source Source) {
		t.Logf("%s %s %v %s", source.URL, action.Action, action.Coordinate, action.Answer)
		if shot, e := a.browser.Screenshot(ctx); e == nil {
			if png, e := base64.StdEncoding.DecodeString(shot); e == nil {
				_ = os.WriteFile(filepath.Join(root, "synthetic-browser-evidence.png"), png, 0600)
			}
		}
	}
	a.debug = func(raw string) {
		if len(raw) > 600 {
			raw = raw[:600]
		}
		t.Logf("Synthetic model action: %s", raw)
	}
	screen, e := a.browser.Screenshot(ctx)
	if e != nil || len(screen) < 1000 {
		t.Fatal("real browser screenshot", e)
	}
	query := os.Getenv("WEBSEARCH_ACCEPTANCE_QUERY")
	expected := os.Getenv("WEBSEARCH_ACCEPTANCE_EXPECTED")
	if query == "" {
		query = "Open https://example.com and quote the prominent English notice at the top of the page, with its source URL."
		expected = "This domain is for use in documentation examples"
	}
	result, e := a.Search(ctx, query, 20)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("Actual research result: %+v", result)
	evidence := ""
	for _, source := range result.Sources {
		evidence += source.Text + "\n"
	}
	if expected == "" || result.Status != "completed" || !strings.Contains(strings.ToLower(evidence), strings.ToLower(expected)) {
		t.Fatal("did not report the independently verified live page text")
	}
	if len(result.Sources) == 0 {
		t.Fatal("missing browser source metadata")
	}
}
