package suite

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// UpdateControllerDigest identifies the actual running executable, not the
// installed configuration, which can change while the controller stays open.
func UpdateControllerDigest(path string) (string, error) {
	sum, _, err := fileDigest(path)
	return sum, err
}

// Read only the stable package identity, even when an older controller cannot
// interpret a newer installation's complete catalog.
func UpdateControllerCurrent(data, digest string) bool {
	var header struct {
		Catalog struct {
			Packages map[string]struct {
				SHA256 string `json:"sha256"`
			} `json:"packages"`
		} `json:"catalog"`
	}
	if digest == "" || platform.LoadJSON(filepath.Join(data, "workspace.json"), &header) != nil {
		return false
	}
	return header.Catalog.Packages["manager"].SHA256 == digest
}

func controllerMatches(c Config, window updateWindow) bool {
	return window.RemoteProtocol == RemoteUpdateProtocol && window.BinarySHA256 != "" && window.BinarySHA256 == c.Catalog.Packages["manager"].SHA256
}

func closeIdleUpdateController(ctx context.Context, target *url.URL) error {
	client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	origin := "http://" + target.Host
	request, _ := http.NewRequestWithContext(ctx, "GET", origin+"/api/state", nil)
	request.Header.Set("X-Workspace-Update", target.Fragment)
	response, err := client.Do(request)
	if err != nil {
		return errors.New("The existing updater is not responding. Retry after it finishes.")
	}
	var state struct {
		Busy *bool `json:"busy"`
	}
	err = json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(&state)
	response.Body.Close()
	if err != nil || response.StatusCode != http.StatusOK || state.Busy == nil || *state.Busy {
		return errors.New("An update operation may still be running. Wait for it to finish, then reopen Updates.")
	}
	request, _ = http.NewRequestWithContext(ctx, "POST", origin+"/api/close", nil)
	request.Header.Set("X-Workspace-Update", target.Fragment)
	request.Header.Set("Origin", origin)
	response, err = client.Do(request)
	if err != nil {
		return errors.New("Could not restart the idle updater. Retry shortly.")
	}
	response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return errors.New("The updater became busy. Retry after it finishes.")
	}
	return nil
}

func retireUpdateController(ctx context.Context, window updateWindow, target *url.URL) error {
	if err := closeIdleUpdateController(ctx, target); err != nil {
		return err
	}
	for i := 0; i < 40; i++ {
		if !platform.ProcessRunning(window.PID, window.Executable) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return errors.New("The previous updater is closing. Reopen Updates to continue.")
}

// Start through the installed manager so it retires the idle old helper and
// copies the verified installed executable before opening a new local page.
func ReopenInstalledUpdates(data string) error {
	command := exec.Command(filepath.Join(data, "Workspace.exe"), "--data", data, "updates")
	if err := platform.StartDetached(command); err != nil {
		return err
	}
	go command.Wait()
	return nil
}
