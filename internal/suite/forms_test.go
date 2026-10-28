package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"strings"
	"testing"
)

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
