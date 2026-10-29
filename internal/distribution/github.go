// Package distribution resolves published Windows bundles from GitHub.
package distribution

import (
	"archive/zip"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

// Set by the local release build; no credential is ever embedded.
var Repository = "oneworkspacestudio/one-workspace"

const bundleName = "workspace-bundle.zip"
const maxBundleSize int64 = 2 << 30

type Asset struct {
	Name   string `json:"name"`
	URL    string `json:"browser_download_url"`
	APIURL string `json:"url"`
	State  string `json:"state"`
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}
type Release struct {
	Tag       string    `json:"tag_name"`
	Draft     bool      `json:"draft"`
	Published time.Time `json:"published_at"`
	Assets    []Asset   `json:"assets"`
}
type Bundle struct {
	ZIP       *zip.Reader
	Catalog   catalog.Catalog
	Overrides map[string]string
	close     func()
}

func (b *Bundle) Close() {
	if b != nil && b.close != nil {
		b.close()
		b.close = nil
	}
}

func githubToken(ctx context.Context) string {
	for _, name := range []string{"GH_TOKEN", "GITHUB_TOKEN"} {
		if token := strings.TrimSpace(os.Getenv(name)); token != "" {
			return token
		}
	}
	return ""
}

func client() *http.Client {
	return &http.Client{Timeout: 10 * time.Minute, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > 6 || req.URL.Scheme != "https" {
			return errors.New("download redirect rejected")
		}
		host := req.URL.Hostname()
		if host != "github.com" && host != "api.github.com" && host != "release-assets.githubusercontent.com" && host != "objects.githubusercontent.com" {
			return errors.New("download host rejected")
		}
		// GitHub asset redirects use signed URLs; do not forward the API credential.
		if host != "api.github.com" {
			req.Header.Del("Authorization")
		}
		return nil
	}}
}

func request(ctx context.Context, endpoint, token string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "OneWorkspaceSetup")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Cache-Control", "no-cache")
	if token != "" && req.URL.Hostname() == "api.github.com" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req, nil
}

func Latest(ctx context.Context, progress func(string), temporaryRoots ...string) (*Bundle, error) {
	if progress == nil {
		progress = func(string) {}
	}
	progress("Checking GitHub for the latest Workspace apps")
	token := githubToken(ctx)
	httpClient := client()
	ctx, err := withReleaseIndex(ctx, token)
	if err != nil {
		return nil, err
	}
	releases, err := releasesFor(ctx, httpClient, token, Repository)
	if err != nil {
		return nil, err
	}
	release, asset, err := choose(releases)
	if err != nil {
		return nil, err
	}
	if !validAsset(asset, Repository, release.Tag) {
		return nil, errors.New("Latest release has an invalid bundle URL, size or SHA256. Download the newest installer or contact the publisher")
	}
	progress("Downloading latest Workspace " + strings.TrimPrefix(release.Tag, "v"))
	temporaryRoot := ""
	if len(temporaryRoots) > 0 {
		temporaryRoot = temporaryRoots[0]
	}
	dir, err := os.MkdirTemp(temporaryRoot, "workspace-download-")
	if err != nil {
		return nil, err
	}
	success := false
	defer func() {
		if !success {
			os.RemoveAll(dir)
		}
	}()
	path := filepath.Join(dir, bundleName)
	if err := downloadMailAsset(ctx, httpClient, token, asset, path); err != nil {
		return nil, fmt.Errorf("Workspace bundle download: %w", err)
	}
	progress("Verifying the latest application bundle")
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, errors.New("Downloaded bundle is not a ZIP archive")
	}
	cat, err := bundleCatalog(&archive.Reader)
	if err == nil && cat.Version != strings.TrimPrefix(release.Tag, "v") {
		err = errors.New("Release version and bundle version do not match")
	}
	if err == nil {
		err = catalog.Verify(&archive.Reader, cat)
	}
	if err != nil {
		archive.Close()
		return nil, fmt.Errorf("Latest bundle could not be verified: %w. Download the newest installer if its catalog format has changed", err)
	}
	overrides, err := latestMail(ctx, httpClient, token, dir, &cat, progress)
	if err == nil {
		err = catalog.Verify(&archive.Reader, cat, overrides)
	}
	if err != nil {
		archive.Close()
		return nil, fmt.Errorf("Latest organization Mail could not be verified: %w. No older Mail was installed", err)
	}
	success = true
	return &Bundle{ZIP: &archive.Reader, Catalog: cat, Overrides: overrides, close: func() { archive.Close(); os.RemoveAll(dir) }}, nil
}
func choose(releases []Release) (Release, Asset, error) {
	versionTag := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?$`)
	candidates := make([]Release, 0, len(releases))
	for _, release := range releases {
		if !release.Draft && !release.Published.IsZero() && versionTag.MatchString(release.Tag) {
			candidates = append(candidates, release)
		}
	}
	// A delayed older release must not displace a higher Workspace version.
	// Publication time only breaks ties between equivalent version labels.
	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := strings.TrimPrefix(candidates[i].Tag, "v"), strings.TrimPrefix(candidates[j].Tag, "v")
		if newerWorkspace(a, b) {
			return true
		}
		if newerWorkspace(b, a) {
			return false
		}
		return candidates[i].Published.After(candidates[j].Published)
	})
	for _, release := range candidates {
		for _, asset := range release.Assets {
			if asset.Name == bundleName {
				return release, asset, nil
			}
		}
	}
	return Release{}, Asset{}, errors.New("No published Workspace application bundle is available yet. Retry after the publisher finishes the release, or deliberately use --offline")
}
func validAsset(a Asset, repository, tag string) bool {
	return validNamedAsset(a, repository, tag, bundleName, maxBundleSize)
}
func validNamedAsset(a Asset, repository, tag, name string, maxSize int64) bool {
	if a.Name != name || a.State != "uploaded" || a.Size <= 0 || a.Size > maxSize || !strings.HasPrefix(a.Digest, "sha256:") {
		return false
	}
	digest, err := hex.DecodeString(strings.TrimPrefix(a.Digest, "sha256:"))
	if err != nil || len(digest) != 32 {
		return false
	}
	api, err := url.Parse(a.APIURL)
	if err != nil || api.Scheme != "https" || api.Host != "api.github.com" || api.User != nil || api.RawQuery != "" || api.Fragment != "" {
		return false
	}
	prefix := "/repos/" + repository + "/releases/assets/"
	id := strings.TrimPrefix(api.Path, prefix)
	if !strings.HasPrefix(api.Path, prefix) || id == "" || strings.Trim(id, "0123456789") != "" {
		return false
	}
	expected := "https://github.com/" + repository + "/releases/download/" + url.PathEscape(tag) + "/" + name
	return a.URL == expected
}
func bundleCatalog(z *zip.Reader) (catalog.Catalog, error) {
	var found *zip.File
	for _, f := range z.File {
		if f.Name == "catalog.json" {
			if found != nil || f.UncompressedSize64 > 1<<20 {
				return catalog.Catalog{}, errors.New("invalid bundle catalog")
			}
			found = f
		}
	}
	if found == nil {
		return catalog.Catalog{}, errors.New("bundle catalog missing")
	}
	in, err := found.Open()
	if err != nil {
		return catalog.Catalog{}, err
	}
	defer in.Close()
	data, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return catalog.Catalog{}, err
	}
	return catalog.Read(data)
}

type downloadProgress struct {
	total, done int64
	tag         string
	last        time.Time
	report      func(string)
}

func (p *downloadProgress) Write(data []byte) (int, error) {
	p.done += int64(len(data))
	if time.Since(p.last) > time.Second || p.done == p.total {
		p.last = time.Now()
		p.report(fmt.Sprintf("Downloading Workspace %s: %d%% (%d / %d MB)", strings.TrimPrefix(p.tag, "v"), 100*p.done/p.total, p.done>>20, p.total>>20))
	}
	return len(data), nil
}
