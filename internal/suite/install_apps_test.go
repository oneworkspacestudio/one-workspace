package suite

import "testing"

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
