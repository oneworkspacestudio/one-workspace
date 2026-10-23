package apphost

import (
	"strings"
	"testing"
)

func TestCADWorkerEvalIsScoped(t *testing.T) {
	h := testHost(t)
	h.Config.ID = "my3d"
	for _, path := range []string{"/", "/assets/index-abc.js", "/assets/kernel.worker-abc.js"} {
		w := call(h, "GET", path, nil, nil, "")
		csp := w.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, "'wasm-unsafe-eval'") {
			t.Fatal("WASM missing", csp)
		}
		has := strings.Contains(csp, " 'unsafe-eval'")
		if has != (path == "/assets/kernel.worker-abc.js") {
			t.Fatal("eval escaped worker scope", path, csp)
		}
	}
}
