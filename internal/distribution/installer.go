package distribution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
)

const installerName = "WorkspaceSetup-windows-amd64.exe"

// InstallerRelease remains readable when a future application catalog needs a
// newer installer. It never authorizes the old controller to interpret that catalog.
type InstallerRelease struct {
	Version  string         `json:"version"`
	Source   string         `json:"source"`
	Required bool           `json:"required,omitempty"`
	Download UpdateDownload `json:"download"`
}

func installerRelease(r Release) *InstallerRelease {
	a, ok := releaseAsset(r, installerName)
	if !ok || !validNamedAsset(a, Repository, r.Tag, installerName, maxBundleSize) {
		return nil
	}
	return &InstallerRelease{Version: strings.TrimPrefix(r.Tag, "v"), Source: "https://github.com/" + Repository + "/releases/tag/" + r.Tag, Download: UpdateDownload{Asset: a, Repository: Repository, Tag: r.Tag}}
}

func requiresInstaller(raw []byte) bool {
	var header struct {
		Schema               int    `json:"schema"`
		InstallerIntegration int    `json:"installerIntegration"`
		Version              string `json:"version"`
	}
	return json.Unmarshal(raw, &header) == nil && header.Version != "" && (header.Schema > 1 || header.InstallerIntegration > catalog.InstallerIntegration)
}

func ValidateInstallerPlan(current catalog.Catalog, p UpdatePlan) error {
	i := p.Installer
	if p.BaseCatalogSHA256 != catalogIdentity(current) || !p.Fresh(time.Now()) {
		return errors.New("Workspace or the release check changed. Check for updates again")
	}
	if i == nil || i.Version != strings.TrimPrefix(i.Download.Tag, "v") || i.Source != "https://github.com/"+Repository+"/releases/tag/"+i.Download.Tag || i.Download.Repository != Repository || i.Download.Member != "" || !validNamedAsset(i.Download.Asset, Repository, i.Download.Tag, installerName, maxBundleSize) {
		return errors.New("No verified Workspace installer is available. Check for updates again")
	}
	if i.Version != current.Version && !newerWorkspace(i.Version, current.Version) {
		return errors.New("The published installer is older than this Workspace")
	}
	return nil
}

// DownloadInstaller is called only with a validated release plan and an exclusive
// owner-protected directory. No executable is exposed before its digest is checked.
func DownloadInstaller(ctx context.Context, i InstallerRelease, dir string, report func(string)) (string, error) {
	if i.Download.Repository != Repository || i.Version != strings.TrimPrefix(i.Download.Tag, "v") || !validNamedAsset(i.Download.Asset, Repository, i.Download.Tag, installerName, maxBundleSize) {
		return "", errors.New("Installer release is not verified")
	}
	a := i.Download.Asset
	req, e := request(ctx, a.APIURL, githubToken(ctx))
	if e != nil {
		return "", e
	}
	req.Header.Set("Accept", "application/octet-stream")
	r, e := client().Do(req)
	if e != nil {
		return "", e
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return "", errors.New("Workspace installer download failed; your installation has not been changed")
	}
	path := filepath.Join(dir, installerName)
	f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return "", e
	}
	h := sha256.New()
	counter := &installerProgress{ctx: ctx, total: a.Size, report: report}
	n, e := io.Copy(io.MultiWriter(f, h, counter), io.LimitReader(r.Body, a.Size+1))
	if e == nil {
		e = f.Sync()
	}
	closeErr := f.Close()
	if e != nil || closeErr != nil || n != a.Size || "sha256:"+hex.EncodeToString(h.Sum(nil)) != a.Digest {
		os.Remove(path)
		return "", errors.New("Workspace installer failed size or SHA256 verification; your installation has not been changed")
	}
	return path, nil
}

type installerProgress struct {
	ctx         context.Context
	total, done int64
	last        time.Time
	report      func(string)
}

func (p *installerProgress) Write(b []byte) (int, error) {
	p.done += int64(len(b))
	if time.Since(p.last) > time.Second || p.done == p.total {
		p.last = time.Now()
		operationprogress.Report(p.ctx, "download", "Downloading the verified Workspace installer", p.done, p.total, "bytes")
		if p.report != nil {
			p.report("Downloading the verified Workspace installer")
		}
	}
	return len(b), nil
}
