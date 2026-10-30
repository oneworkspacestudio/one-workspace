package suite

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestUpdateControllerReplacementRetriesOnlyTransientAccessFailure(t *testing.T) {
	calls := 0
	err := retryUpdaterReplace(context.Background(), func() error {
		calls++
		if calls < 3 {
			return os.ErrPermission
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Fatalf("transient retry: calls=%d error=%v", calls, err)
	}
	calls = 0
	err = retryUpdaterReplace(context.Background(), func() error { calls++; return os.ErrNotExist })
	if !errors.Is(err, os.ErrNotExist) || calls != 1 {
		t.Fatal("unexpected retry for non-sharing failure")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls = 0
	err = retryUpdaterReplace(ctx, func() error { calls++; return os.ErrPermission })
	if !errors.Is(err, os.ErrPermission) || calls != 1 {
		t.Fatal("retry ignored cancellation")
	}
}

func TestUpdateControllerIdentitySurvivesProtocolReuse(t *testing.T) {
	c := Config{Catalog: catalog.Catalog{Packages: map[string]catalog.Package{"manager": {SHA256: "new-binary"}}}}
	for _, tc := range []struct {
		digest   string
		protocol int
		want     bool
	}{
		{"new-binary", RemoteUpdateProtocol, true},
		{"old-binary", RemoteUpdateProtocol, false},
		{"", RemoteUpdateProtocol, false},
		{"new-binary", RemoteUpdateProtocol - 1, false},
	} {
		if got := controllerMatches(c, updateWindow{BinarySHA256: tc.digest, RemoteProtocol: tc.protocol}); got != tc.want {
			t.Fatalf("%+v: %v", tc, got)
		}
	}
}

func TestUpdateControllerRecognisesUpgradeWithoutParsingNewCatalog(t *testing.T) {
	data := t.TempDir()
	save := func(hash string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(data, "workspace.json"), []byte(`{"schema":999,"catalog":{"schema":999,"packages":{"manager":{"sha256":"`+hash+`"}}}}`), 0600); err != nil {
			t.Fatal(err)
		}
	}
	save("before")
	if !UpdateControllerCurrent(data, "before") {
		t.Fatal("current binary rejected")
	}
	save("after")
	if UpdateControllerCurrent(data, "before") {
		t.Fatal("old controller accepted after upgrade")
	}
	if !UpdateControllerCurrent(data, "after") {
		t.Fatal("new binary rejected with future catalog")
	}
}

func TestUpdateControllerRetirementRequiresConfirmedIdleAndClose(t *testing.T) {
	for _, tc := range []struct {
		name, state       string
		closeCode         int
		wantClose, wantOK bool
	}{
		{"idle", `{"busy":false}`, 202, true, true},
		{"active", `{"busy":true}`, 202, false, false},
		{"missing status", `{}`, 202, false, false},
		{"malformed", `no`, 202, false, false},
		{"became busy", `{"busy":false}`, 409, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			closed := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Workspace-Update") != "synthetic-capability" {
					t.Error("missing capability")
					w.WriteHeader(403)
					return
				}
				switch r.URL.Path {
				case "/api/state":
					w.Write([]byte(tc.state))
				case "/api/close":
					closed = true
					if r.Method != "POST" || r.Header.Get("Origin") != "http://"+r.Host {
						t.Error("missing origin/method")
					}
					w.WriteHeader(tc.closeCode)
				default:
					t.Error("unexpected path")
				}
			}))
			defer server.Close()
			target, _ := url.Parse(server.URL + "/#synthetic-capability")
			err := closeIdleUpdateController(context.Background(), target)
			if (err == nil) != tc.wantOK || closed != tc.wantClose {
				t.Fatalf("closed=%v error=%v", closed, err)
			}
		})
	}
}
