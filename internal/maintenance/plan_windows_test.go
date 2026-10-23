package maintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

// Opt-in acceptance: a real protected app update through an expired HTTP plan.
// Never accepts an owner installation or an arbitrary path from the environment.
func TestNativeExpiredUpdateWindow(t *testing.T) {
	if os.Getenv("WORKSPACE_NATIVE_EXPIRED_UPDATE") != "1" {
		t.Skip("requires named disposable installation")
	}
	data := `C:\Projects\.cache\update-expiry-20261023\fixture\Private`
	if !strings.Contains(strings.ToLower(os.Getenv("APPDATA")), "update-expiry-20261023") {
		t.Fatal("fixture AppData is not isolated")
	}
	before, secrets, err := suite.Load(data)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	t.Cleanup(func() {
		stop, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := suite.Stop(stop, before); err != nil {
			t.Error(err)
		}
	})
	plan, err := distribution.CheckUpdates(ctx, before.Catalog)
	if err != nil {
		t.Fatal(err)
	}
	wantVersion := ""
	for _, ch := range plan.Choices {
		if ch.ID == "drive" && ch.Available {
			wantVersion = ch.Version
		}
	}
	if wantVersion == "" {
		t.Fatal("fixture requires an available Drive update")
	}
	plan.CheckedAt = time.Now().Add(-65 * time.Minute)
	marker := filepath.Join(before.Files, "Documents", "expired-update-preserve.txt")
	if err := os.WriteFile(marker, []byte("preserve this user document"), 0600); err != nil {
		t.Fatal(err)
	}
	s := &Server{Data: data, Token: "native-fixture-capability", Context: ctx, Initial: before, Plan: &plan}
	server := httptest.NewServer(s.Handler())
	defer server.Close()
	s.Origin = server.URL
	r, _ := http.NewRequestWithContext(ctx, "POST", server.URL+"/api/apply", strings.NewReader(`{"selected":["drive"]}`))
	r.Header.Set("Origin", server.URL)
	r.Header.Set("X-Workspace-Update", s.Token)
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("expired plan rejected: HTTP %d", response.StatusCode)
	}
	acceptedStatus := response.StatusCode
	for {
		s.mu.Lock()
		busy, problem := s.Busy, s.Problem
		s.mu.Unlock()
		if !busy {
			if problem != "" {
				t.Fatal(problem)
			}
			break
		}
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(time.Second):
		}
	}
	after, afterSecrets, err := suite.Load(data)
	if err != nil {
		t.Fatal(err)
	}
	drive, _ := after.Catalog.Find("drive")
	if drive.Version != wantVersion || !reflect.DeepEqual(secrets, afterSecrets) || before.Owner != after.Owner || before.Programs != after.Programs || before.Files != after.Files {
		t.Fatal("expected version or preserved identity/storage differs")
	}
	for _, app := range before.Catalog.Apps {
		if app.ID != "drive" && app.ID != "docs" && app.ID != "photos" {
			got, _ := after.Catalog.Find(app.ID)
			if !reflect.DeepEqual(app, got) {
				t.Fatalf("unselected app changed: %s", app.ID)
			}
		}
	}
	content, err := os.ReadFile(marker)
	if err != nil || string(content) != "preserve this user document" {
		t.Fatal("user document not preserved")
	}
	response, err = http.Get(fmt.Sprintf("http://localhost:%d/healthz", after.Port))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var health struct {
		Status string            `json:"status"`
		Apps   map[string]string `json:"apps"`
	}
	if json.NewDecoder(response.Body).Decode(&health) != nil || health.Status != "ok" || len(health.Apps) == 0 {
		t.Fatal("health check failed")
	}
	for id, status := range health.Apps {
		if status != "running" {
			t.Fatalf("%s: %s", id, status)
		}
	}
	backups, err := suite.UpdateBackups(after)
	if err != nil || len(backups) == 0 {
		t.Fatal("no recovery backup", err)
	}
	if path := os.Getenv("WORKSPACE_NATIVE_EXPIRED_RECEIPT"); path != "" {
		body, _ := json.MarshalIndent(map[string]any{"expiredMinutes": 65, "httpStatus": acceptedStatus, "updatedGroup": "drive", "updatedVersion": drive.Version, "ownerAndCredentialsPreserved": true, "documentPreserved": true, "unselectedAppsPreserved": true, "backupRecorded": true, "services": health.Apps, "publicReleaseRecheck": true, "realAppDownloadAndInstall": true}, "", "  ")
		if err := os.WriteFile(path, body, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
