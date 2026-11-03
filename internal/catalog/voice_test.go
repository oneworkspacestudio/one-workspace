package catalog

import "testing"

func TestVoiceSelectedAndOmitted(t *testing.T) {
	c := Catalog{Packages: map[string]Package{"voice": {File: "MyVoice.exe"}, "voice-contract": {File: "docs/MyVoice/compatibility.json"}, "voice-manifest": {File: "docs/MyVoice/manifest.json"}, "unrelated": {File: "Other.exe"}}, Apps: []App{{ID: "voice", Package: "voice"}, {ID: "other", Package: "unrelated"}}}
	for _, selected := range []bool{false, true} {
		ids := InstallPackageIDs(c, map[string]bool{"voice": selected, "other": true})
		if ids["voice"] != selected || ids["voice-contract"] != selected || ids["voice-manifest"] != selected || !ids["unrelated"] {
			t.Fatal("incomplete selected Voice payload", ids)
		}
	}
	if CategoryFor("voice").ID != "development" || DefaultInstall("voice") {
		t.Fatal("Voice is an optional creative app")
	}
	full := Catalog{}
	for _, app := range ManagedApps() {
		full.Apps = append(full.Apps, App{ID: app.ID, Name: app.Name, Route: app.Route, OIDC: app.OIDC})
	}
	if e := full.ValidateManagedInventory(); e != nil {
		t.Fatal(e)
	}
	for i, a := range full.Apps {
		if a.ID == "voice" {
			full.Apps = append(full.Apps[:i], full.Apps[i+1:]...)
			break
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("public catalog can silently omit Voice")
	}
}
