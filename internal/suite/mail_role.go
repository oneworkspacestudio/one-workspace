package suite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

var errMainMailService = errors.New("mail accounts are managed on the main Workspace")

// The configured app's loopback endpoint supplies only its service role. Never
// fetch a client-selected remote URL here or share a portal/OIDC credential.
func mailUsesMain(ctx context.Context, app catalog.App) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://127.0.0.1:%d/api/status", app.Port), nil)
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var state struct {
		Main string `json:"mainMailURL"`
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("Mail status is unavailable")
	}
	if err = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&state); err != nil {
		return false, err
	}
	if state.Main == "" {
		return false, nil
	}
	target, err := url.Parse(state.Main)
	return err == nil && target.Scheme == "https" && target.Host != "" && target.User == nil && target.RawQuery == "" && target.Fragment == "" && !strings.Contains(target.Path, "\\"), nil
}
