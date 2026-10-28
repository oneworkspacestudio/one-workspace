package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormsUpgradeAddsOnlyRecoverySource(t *testing.T) {
	var user map[string]any
	_ = json.Unmarshal([]byte(`{"username":"workspace-recovery","loginMethod":"password","permissions":{"admin":true},"scopes":[{"name":"My Drive","scope":"/original","permissions":{"view":true,"delete":false}}]}`), &user)
	scopes := user["scopes"].([]any)
	original, _ := json.Marshal(scopes[0])
	changed, e := formsRecoveryScope(user)
	if e != nil || !changed {
		t.Fatal(changed, e)
	}
	scopes = user["scopes"].([]any)
	preserved, _ := json.Marshal(scopes[0])
	if len(scopes) != 2 || string(original) != string(preserved) {
		t.Fatal("Existing scope changed")
	}
	added := scopes[1].(map[string]any)
	p := added["permissions"].(map[string]any)
	if added["name"] != formsSource || added["scope"] != "/" || p["view"] != true || p["create"] != true || p["delete"] != false || p["modify"] != false {
		t.Fatal("Unexpected service grant")
	}
	if changed, e = formsRecoveryScope(user); e != nil || changed {
		t.Fatal("Migration not idempotent")
	}
	user["username"] = "ordinary-user"
	if _, e = formsRecoveryScope(user); e == nil {
		t.Fatal("Granted an ordinary account")
	}
}

func TestFormsUpdateBarrierAllowsOnlyAuthenticatedReconciliation(t *testing.T) {
	p := &Portal{rt: &Runtime{Secrets: Secrets{Clients: map[string]Client{"forms": {Secret: "fixture-only"}}}}}
	for _, tt := range []struct {
		body, remote, token string
		want                bool
	}{
		{`{"operation":"list","path":"/"}`, "127.0.0.1:1234", "fixture-only", true},
		{`{"operation":"read","path":"/` + strings.Repeat("a", 64) + `/` + strings.Repeat("b", 48) + `/events/000000001.json"}`, "127.0.0.1:1234", "fixture-only", true},
		{`{"operation":"write","path":"/"}`, "127.0.0.1:1234", "fixture-only", false},
		{`{"operation":"bind","path":"/"}`, "127.0.0.1:1234", "fixture-only", false},
		{`{"operation":"list","path":"/"}`, "192.0.2.1:1234", "fixture-only", false},
		{`{"operation":"list","path":"/"}`, "127.0.0.1:1234", "wrong", false},
		{`{"operation":"list","path":"/../private"}`, "127.0.0.1:1234", "fixture-only", false},
	} {
		r := httptest.NewRequest("POST", "/_integration/forms/drive", strings.NewReader(tt.body))
		r.RemoteAddr = tt.remote
		r.Header.Set("Authorization", "Bearer "+formsCapability(tt.token))
		if got := p.formsReconcileRequest(r); got != tt.want {
			t.Fatalf("%s: got %v", tt.body, got)
		}
	}
}

func TestFormsPublicGatewayBoundaries(t *testing.T) {
	id := strings.Repeat("a", 48)
	for _, tt := range []struct {
		method, path string
		want         bool
	}{{"GET", "/forms/", true}, {"GET", "/forms/api/v1/public/" + id, true}, {"POST", "/forms/api/v1/public/" + id + "/responses", true}, {"GET", "/forms/api/v1/public/" + id + "/responses", false}, {"GET", "/forms/api/v1/forms/" + id, false}, {"GET", "/forms/api/v1/public/../forms/" + id, false}, {"GET", "/forms/api/v1/public/" + id + "/receipt/" + id, true}, {"POST", "/forms/auth/callback", false}} {
		if got := publicFormsRequest(httptest.NewRequest(tt.method, tt.path, nil)); got != tt.want {
			t.Fatalf("%s %s = %v", tt.method, tt.path, got)
		}
	}
}
func TestFormsScopedStoragePaths(t *testing.T) {
	owner, id := strings.Repeat("a", 64), strings.Repeat("b", 48)
	for _, p := range []string{"/", "/" + owner, "/" + owner + "/" + id + "/events/000000001.json", "/" + owner + "/" + id + "/uploads/" + id + ".json"} {
		if !formsPath.MatchString(p) {
			t.Fatal("valid rejected", p)
		}
	}
	for _, p := range []string{"/../drive", "/" + owner + "/" + id + "/secrets/a.json", "C:/Users", "//evil", "/" + owner + "/../../"} {
		if formsPath.MatchString(p) {
			t.Fatal("escape accepted", p)
		}
	}
}
func TestFormsManagedRuntime(t *testing.T) {
	c := Config{Data: t.TempDir(), Programs: t.TempDir(), Port: 19300, Origin: "http://localhost:19300", Issuer: "http://localhost:19301", Owner: "owner"}
	s := Secrets{Clients: map[string]Client{"forms": {ID: "forms-client", Secret: "forms-secret"}}}
	rt := Runtime{Config: c, Secrets: s}
	a := catalog.App{ID: "forms", OIDC: true, Executable: "MyForms.exe", Port: 19333, Callback: "/forms/auth/callback"}
	cmd := rt.command(a)
	env := strings.Join(cmd.Env, "\n")
	for _, want := range []string{"WORKSPACE_OWNER_SUBJECT=owner", "MYFORMS_DRIVE_URL=http://127.0.0.1:19300/_integration/forms/drive", "MYFORMS_DRIVE_TOKEN=" + formsCapability("forms-secret")} {
		if !strings.Contains(env, want) {
			t.Fatal("missing", want)
		}
	}
	if len(appUpdatePrivatePaths("forms")) != 2 || !foundationApp("forms") {
		t.Fatal("missing recovery/identity registration")
	}
	if !strings.HasSuffix(cmd.Dir, "forms") {
		t.Fatal("working directory", cmd.Dir)
	}
}
