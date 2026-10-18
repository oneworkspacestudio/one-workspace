package suite

import (
	"bytes"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
)

func TestHomeUsesCurrentWorkspaceAndAppIcons(t *testing.T) {
	var page bytes.Buffer
	err := home.Execute(&page, map[string]any{
		"LoggedIn": true, "User": map[string]string{"Name": "Example"},
		"Tiles": []catalogTile{{"Notes", "Example", "/notes", branding.Icon("notes")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(page.String(), branding.Icon("workspace")) != 2 || !strings.Contains(page.String(), branding.Icon("notes")) {
		t.Fatal("header, favicon or app tile uses stale artwork URLs")
	}
}
