package distribution

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLatestPublishedBundle(t *testing.T) {
	older := Release{Tag: "v0.2.0", Published: time.Unix(1, 0), Assets: []Asset{{Name: bundleName}}}
	newer := Release{Tag: "v0.3.0-preview.15", Published: time.Unix(2, 0), Assets: []Asset{{Name: bundleName}}}
	draft := Release{Tag: "v99", Published: time.Unix(3, 0), Draft: true, Assets: []Asset{{Name: bundleName}}}
	seed := Release{Tag: "components-seed", Published: time.Unix(4, 0), Assets: []Asset{{Name: bundleName}}}
	r, _, err := choose([]Release{older, draft, newer, seed})
	if err != nil || r.Tag != newer.Tag {
		t.Fatalf("wrong release %v %v", r, err)
	}
	if _, _, err = choose([]Release{draft, seed}); err == nil {
		t.Fatal("accepted unpublished bundle")
	}
}
func TestReleaseAssetScope(t *testing.T) {
	repo := "oneworkspacestudio/one-workspace"
	tag := "v0.3.0-preview.15"
	good := Asset{Name: bundleName, State: "uploaded", Size: 10, Digest: "sha256:" + strings.Repeat("a", 64), APIURL: "https://api.github.com/repos/" + repo + "/releases/assets/123", URL: "https://github.com/" + repo + "/releases/download/" + tag + "/" + bundleName}
	if !validAsset(good, repo, tag) {
		t.Fatal("valid asset rejected")
	}
	for _, mutate := range []func(*Asset){func(a *Asset) { a.APIURL = "https://evil.example/x" }, func(a *Asset) { a.APIURL += "?x" }, func(a *Asset) { a.URL = strings.Replace(a.URL, repo, "other/repo", 1) }, func(a *Asset) { a.Digest = "" }, func(a *Asset) { a.Size = maxBundleSize + 1 }, func(a *Asset) { a.State = "new" }} {
		a := good
		mutate(&a)
		if validAsset(a, repo, tag) {
			b, _ := json.Marshal(a)
			t.Fatalf("unsafe asset accepted %s", b)
		}
	}
}
func TestDownloadRedirectCredentialScope(t *testing.T) {
	for _, host := range []string{"release-assets.githubusercontent.com", "objects.githubusercontent.com", "evil.example"} {
		u, _ := url.Parse("https://" + host + "/asset")
		req := &http.Request{URL: u, Header: http.Header{"Authorization": []string{"Bearer fixture"}}}
		err := client().CheckRedirect(req, nil)
		if host == "evil.example" {
			if err == nil {
				t.Fatal("untrusted redirect accepted")
			}
		} else if err != nil || req.Header.Get("Authorization") != "" {
			t.Fatal("credential forwarded or trusted redirect rejected")
		}
	}
}
