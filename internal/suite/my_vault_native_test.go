package suite

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Explicitly enabled only for a disposable, marked synthetic identity fixture.
// The test contains no credentials and never chooses the live installation.
func TestNativeVaultIdentityFixture(t *testing.T) {
	root := os.Getenv("MYVAULT_NATIVE_SSO_FIXTURE")
	if root == "" {
		t.Skip("requires disposable native MyAccounts fixture")
	}
	if _, e := os.Stat(filepath.Join(root, ".synthetic-vault-test")); e != nil {
		t.Fatal("fixture marker missing")
	}
	var fixture struct {
		Config  Config  `json:"config"`
		Secrets Secrets `json:"secrets"`
	}
	b, e := os.ReadFile(filepath.Join(root, "fixture.json"))
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(b, &fixture); e != nil {
		t.Fatal("invalid fixture configuration")
	}
	if !strings.EqualFold(filepath.Clean(fixture.Config.Data), filepath.Clean(root)) {
		t.Fatal("fixture data root does not match")
	}
	rt := &Runtime{Config: fixture.Config, Secrets: fixture.Secrets}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if e = rt.prepareVaultIdentity(ctx); e != nil {
		t.Fatal(e)
	}
}
