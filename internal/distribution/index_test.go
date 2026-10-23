package distribution

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPublicOnlineSetupUsesIndexAndVerifiedBundle(t *testing.T) {
	_, tr, _ := setupUpdateTransport(t, true)
	uri := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
	var rs []Release
	json.Unmarshal(tr[uri], &rs)
	var archive bytes.Buffer
	z := zip.NewWriter(&archive)
	for name, body := range map[string][]byte{"catalog.json": tr[rs[0].Assets[0].APIURL], "Workspace.exe": []byte("new-manager"), "MyNotes.exe": []byte("old-notes"), "licenses/Notes-LICENSE": []byte("license")} {
		f, _ := z.Create(name)
		f.Write(body)
	}
	z.Close()
	a := fixtureAsset(Repository, rs[0].Tag, bundleName, "3", archive.Bytes())
	rs[0].Assets[2] = a
	tr[uri], _ = json.Marshal(rs)
	tr[a.APIURL] = archive.Bytes()
	http.DefaultTransport = mailTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.github.com" {
			t.Fatal("online setup used API while quota exhausted")
		}
		return tr.RoundTrip(r)
	})
	b, err := Latest(context.Background(), nil, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	if b.Catalog.Version != "0.3.0-preview.23" {
		t.Fatal(b.Catalog.Version)
	}
}

// Simulate the exact exhausted shared-IP API allowance while exercising the
// complete public discovery/staging/installer path. Any API call is a failure.
func TestPublicUpdatesAndInstallerIgnoreExhaustedAPILimit(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	uri := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
	var rs []Release
	json.Unmarshal(tr[uri], &rs)
	body := []byte("verified native installer fixture")
	a := fixtureAsset(Repository, rs[0].Tag, installerName, "99", body)
	rs[0].Assets = append(rs[0].Assets, a)
	tr[uri], _ = json.Marshal(rs)
	tr[a.APIURL] = body
	apiCalls, indexCalls, publicDownloads := 0, 0, 0
	http.DefaultTransport = mailTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host == "api.github.com" {
			apiCalls++
			return &http.Response{StatusCode: 403, Header: http.Header{"X-Ratelimit-Remaining": {"0"}}, Body: io.NopCloser(strings.NewReader(`{"message":"API rate limit exceeded"}`))}, nil
		}
		if r.Header.Get("Authorization") != "" {
			t.Fatal("public credential leak")
		}
		if r.URL.String() == publicIndexURL {
			indexCalls++
		} else {
			publicDownloads++
		}
		return tr.RoundTrip(r)
	})
	for n := 0; n < 4; n++ {
		plan, err := CheckUpdates(context.Background(), current)
		if err != nil {
			t.Fatal(err)
		}
		if plan.IndexPublishedAt.IsZero() || plan.Installer == nil {
			t.Fatal("missing channel/installer")
		}
		target, downloads, err := SelectUpdates(current, plan, []string{"workspace", "notes"})
		if err != nil {
			t.Fatal(err)
		}
		if err = StageUpdateFiles(context.Background(), target, downloads, t.TempDir(), func(string) {}); err != nil {
			t.Fatal(err)
		}
		path, err := DownloadInstaller(context.Background(), *plan.Installer, t.TempDir(), nil)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := os.ReadFile(path)
		if string(got) != string(body) {
			t.Fatal("installer changed")
		}
	}
	if apiCalls != 0 || indexCalls != 4 || publicDownloads == 0 {
		t.Fatal(apiCalls, indexCalls, publicDownloads)
	}
}

func TestInvalidOrUnavailableIndexFailsClosed(t *testing.T) {
	_, tr, _ := setupUpdateTransport(t, true)
	for _, mode := range []string{"http403", "malformed", "future", "foreign", "missing-repo", "oversize"} {
		t.Run(mode, func(t *testing.T) {
			index := releaseIndex{Schema: 1, PublishedAt: time.Now(), Repositories: map[string][]Release{Repository: {}}}
			status := 200
			switch mode {
			case "http403":
				status = 403
			case "future":
				index.PublishedAt = time.Now().Add(time.Hour)
			case "foreign":
				index.Repositories["evil/repository"] = nil
			case "missing-repo":
				delete(index.Repositories, Repository)
				index.Repositories[mailRepository] = nil
			}
			data, _ := json.Marshal(index)
			if mode == "malformed" {
				data = []byte(`{"schema":1`)
			}
			if mode == "oversize" {
				data = []byte(strings.Repeat(" ", maxIndexSize+1))
			}
			http.DefaultTransport = mailTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != publicIndexURL {
					t.Fatal("failed index fell back to API", r.URL)
				}
				return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(string(data))), Header: http.Header{}}, nil
			})
			if _, err := CheckUpdates(context.Background(), updateCatalogFixture()); err == nil {
				t.Fatal("invalid channel accepted")
			}
		})
	}
	http.DefaultTransport = tr
}

func TestAuthenticatedFallbackRateLimitGuidance(t *testing.T) {
	_, _, _ = setupUpdateTransport(t, true)
	t.Setenv("GH_TOKEN", "fixture-private-token")
	http.DefaultTransport = mailTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host != "api.github.com" || r.Header.Get("Authorization") != "Bearer fixture-private-token" {
			t.Fatal("wrong private request", r.URL)
		}
		return &http.Response{StatusCode: 403, Header: http.Header{"X-Ratelimit-Remaining": {"0"}, "Retry-After": {"120"}}, Body: io.NopCloser(strings.NewReader("rate limited"))}, nil
	})
	_, err := CheckUpdates(context.Background(), updateCatalogFixture())
	if err == nil || !strings.Contains(err.Error(), "Retry in 120 seconds") {
		t.Fatal(err)
	}
}
