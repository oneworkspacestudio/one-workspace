package distribution

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestManagerRequirementsVersionAndCapabilities(t *testing.T) {
	m := AppUpdateManifest{Schema: 1, Integration: 2}
	if err := RequireManager(&m, "0.3.0-preview.43", []string{"app-update-requirements-v1", "chat-state-v2"}); err != nil {
		t.Fatal(err)
	}
	if m.Schema != 2 || m.Integration != 3 {
		t.Fatal("requirements silently retained old contract")
	}
	if err := checkManifestRequirements(m, "0.3.0-preview.42"); err == nil {
		t.Fatal("old installed manager accepted new app")
	}
	for _, version := range []string{"0.3.0-preview.43", "0.3.0-preview.44", "0.3.0"} {
		if err := checkManifestRequirements(m, version); err != nil {
			t.Fatalf("compatible %s: %v", version, err)
		}
	}
	m.RequiredCapabilities = []string{"unimplemented-encryption-capability"}
	if err := checkManifestRequirements(m, "0.3.0-preview.43"); err == nil {
		t.Fatal("unimplemented origin capability advertised")
	}
	m.RequiredCapabilities = []string{"chat-state-v2", "chat-state-v2"}
	if err := checkManifestRequirements(m, "0.3.0-preview.43"); err == nil {
		t.Fatal("duplicate capability accepted")
	}
	m.Schema = 1
	m.Integration = 2
	if err := checkManifestRequirements(m, "0.3.0-preview.43"); err == nil {
		t.Fatal("schema1 concealed requirements")
	}
}

func TestFrozenLegacyDecoderRejectsRequiredManifest(t *testing.T) {
	// Mirrors the published42 strict JSON contract. An actual downloaded old
	// manager offered the published candidate remains an integration release gate.
	type legacyManifest struct {
		Schema      int `json:"schema"`
		Integration int `json:"integration"`
	}
	var old legacyManifest
	d := json.NewDecoder(strings.NewReader(`{"schema":2,"integration":3,"minimumManager":"0.3.0-preview.43"}`))
	d.DisallowUnknownFields()
	if err := d.Decode(&old); err == nil {
		t.Fatal("old strict parser ignored manager requirement")
	}
}

func TestApplyPlanRechecksRequirementBeforePackages(t *testing.T) {
	current := updateCatalogFixture()
	current.Version = "0.3.0-preview.42"
	plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{{ID: "chat", Available: true, MinimumManager: "0.3.0-preview.43", RequiredCapabilities: []string{"chat-state-v2"}}}}
	if _, _, err := SelectUpdates(current, plan, []string{"chat"}); err == nil {
		t.Fatal("stored plan bypassed manager gate")
	}
}
