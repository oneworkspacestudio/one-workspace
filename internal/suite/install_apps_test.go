package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"reflect"
	"testing"
)

func TestReenableAnyVerifiedAppAndDependencies(t *testing.T) {
	c := catalog.Catalog{Apps: []catalog.App{{ID: "accounts"}, {ID: "drive", Depends: []string{"accounts"}}, {ID: "photos", Depends: []string{"drive"}}, {ID: "music"}}}
	s := &catalog.AppSelection{Excluded: []string{"photos", "drive", "music"}}
	if err := reenableRequestedApps(c, s, []string{"photos", "photos"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s.Excluded, []string{"music"}) {
		t.Fatal("unrelated opt-out changed", s.Excluded)
	}
	if err := reenableRequestedApps(c, s, []string{"music", "unknown"}); err == nil {
		t.Fatal("unknown app accepted")
	}
	if !reflect.DeepEqual(s.Excluded, []string{"music"}) {
		t.Fatal("failed request changed selection")
	}
}

func TestOrganizationAppSources(t *testing.T) {
	for _, raw := range []string{"https://github.com/oneworkspacestudio/drive", "https://github.com/oneworkspacestudio/Photo-Editor/releases/tag/studio-v0.1.0"} {
		if !organizationSource(raw) {
			t.Fatal(raw)
		}
	}
	for _, raw := range []string{"https://github.com/fahimc/filebrowser", "https://github.com/oneworkspacestudio.evil/drive", "http://github.com/oneworkspacestudio/mail", "https://github.com@evil.example/oneworkspacestudio/mail", "https://github.com/oneworkspacestudio/mail?download=x"} {
		if organizationSource(raw) {
			t.Fatal(raw)
		}
	}
}
