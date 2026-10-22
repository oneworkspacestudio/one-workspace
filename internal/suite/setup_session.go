package suite

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type SetupSession struct {
	Schema    int       `json:"schema"`
	Data      string    `json:"data"`
	Programs  string    `json:"programs"`
	Files     string    `json:"files"`
	State     string    `json:"state"`
	Phase     string    `json:"phase"`
	Error     string    `json:"error,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
	Progress  Progress  `json:"progress"`
}
type setupCheckpoint struct {
	Request Request
	Config  Config
	Phase   string
	Offline bool
}

func setupSessionDir(data string) string {
	h := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(data))))
	return filepath.Join(filepath.Dir(data), ".workspace-install-"+hex.EncodeToString(h[:8]))
}
func SetupSessionDirectory(data string) string { return setupSessionDir(data) }
func prepareSetupSession(r Request) error {
	dir := setupSessionDir(r.Data)
	if e := platform.RejectReparse(dir); e != nil && !os.IsNotExist(e) {
		return e
	}
	if e := os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	return platform.ProtectDir(dir)
}
func saveSetupCheckpoint(r Request, c Config, phase string) error {
	return platform.SaveSecret(filepath.Join(setupSessionDir(r.Data), "continue.dpapi"), setupCheckpoint{r, c, phase, r.Offline})
}
func SaveSetupSession(r Request, p Progress) error {
	if e := prepareSetupSession(r); e != nil {
		return e
	}
	return platform.SaveJSON(filepath.Join(setupSessionDir(r.Data), "status.json"), SetupSession{1, r.Data, r.Programs, r.Files, p.State, p.Phase, p.Error, time.Now().UTC(), p})
}
func ReadSetupSession(r Request) (SetupSession, bool, error) {
	var status SetupSession
	if e := r.ResolveStorage(); e != nil {
		return status, false, e
	}
	for _, root := range []string{r.Data, r.Programs, r.Files, setupSessionDir(r.Data)} {
		for p := filepath.Clean(root); ; p = filepath.Dir(p) {
			if e := platform.RejectReparse(p); e != nil && !os.IsNotExist(e) {
				return status, false, e
			}
			if filepath.Dir(p) == p {
				break
			}
		}
	}
	if _, e := platform.CheckDataVolume(r.Data); e != nil {
		return status, false, e
	}
	e := platform.LoadJSON(filepath.Join(setupSessionDir(r.Data), "status.json"), &status)
	if os.IsNotExist(e) {
		for _, path := range []string{setupLockPath(r), filepath.Join(r.Data, "install.lock")} {
			if _, err := os.Stat(path); err != nil {
				continue
			}
			lock, err := platform.Lock(path)
			if err != nil {
				return SetupSession{Schema: 1, Data: r.Data, Programs: r.Programs, Files: r.Files, State: "installing", Phase: "Another installer is working in this folder", UpdatedAt: time.Now().UTC()}, true, nil
			}
			lock.Close()
		}
		return status, false, nil
	}
	if e != nil {
		return status, false, e
	}
	if status.Schema != 1 || !sameSetupPaths(r, status) {
		return status, false, errors.New("Saved setup belongs to different folders")
	}
	guard, e := platform.Lock(setupLockPath(r))
	active := e != nil
	if guard != nil {
		guard.Close()
	}
	if !active && status.State == "installing" {
		status.State = "interrupted"
		status.Progress.State = "interrupted"
		status.Error = "The installer stopped before recording completion. Continue to verify installed files and finish the remaining steps."
	}
	return status, active, nil
}
func sameSetupPaths(r Request, s SetupSession) bool {
	return strings.EqualFold(filepath.Clean(r.Data), filepath.Clean(s.Data)) && strings.EqualFold(filepath.Clean(r.Programs), filepath.Clean(s.Programs)) && strings.EqualFold(filepath.Clean(r.Files), filepath.Clean(s.Files))
}

func ContinueSetup(ctx context.Context, r Request, z *zip.Reader, cat catalog.Catalog, report func(string)) (Config, error) {
	if e := r.ResolveStorage(); e != nil {
		return Config{}, e
	}
	if _, e := InspectSetupStorage(r); e != nil {
		return Config{}, e
	}
	if e := platform.RejectReparse(setupSessionDir(r.Data)); e != nil {
		return Config{}, e
	}
	lock, e := platform.Lock(setupLockPath(r))
	if e != nil {
		return Config{}, errors.New("Installation is still running. Reconnect to its progress instead")
	}
	defer lock.Close()
	var cp setupCheckpoint
	if e = platform.LoadSecret(filepath.Join(setupSessionDir(r.Data), "continue.dpapi"), &cp); e != nil {
		return Config{}, errors.New("This older installation has no resumable checkpoint. Open an installed Workspace, or use the protected reinstall recovery option")
	}
	cp.Request.Offline = cp.Offline
	if !sameSetupPaths(cp.Request, SetupSession{Data: r.Data, Programs: r.Programs, Files: r.Files}) {
		return Config{}, errors.New("Saved installation folders do not match")
	}
	if c, _, err := Load(r.Data); err == nil && c.Owner != "" {
		return c, nil
	}
	// A crash during external DNS changes needs a fresh explicit review; do not
	// repeat an uncertain external operation or silently create another tunnel.
	if cp.Phase == "domain" {
		return Config{}, errors.New("The domain connection was interrupted. Its records were retained; review the connection through recovery before continuing")
	}
	return installFreshFrom(ctx, cp.Request, z, cat, report, &cp)
}

func SetupHintPath() string {
	base, e := os.UserConfigDir()
	if e != nil {
		return ""
	}
	return filepath.Join(base, "OneWorkspace", "Setup", "last.json")
}
func RememberSetup(r Request) error {
	p := SetupHintPath()
	if p == "" {
		return nil
	}
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		return e
	}
	if e := platform.ProtectDir(filepath.Dir(p)); e != nil {
		return e
	}
	return platform.SaveJSON(p, SetupSession{Schema: 1, Data: r.Data, Programs: r.Programs, Files: r.Files})
}
func LastSetup() (Request, error) {
	var s SetupSession
	e := platform.LoadJSON(SetupHintPath(), &s)
	return Request{Data: s.Data, Programs: s.Programs, Files: s.Files}, e
}
