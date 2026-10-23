package distribution

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

type mailTransport func(*http.Request) (*http.Response, error)

func (f mailTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestOrganizationMailRelease(t *testing.T) {
	for _, mode := range []string{"new", "same", "corrupt", "incompatible", "missing-manifest", "foreign-repo", "draft"} {
		t.Run(mode, func(t *testing.T) {
			body := "new Mail executable"
			marker := []byte(`{"schema":1,"integration":1,"version":"0.1.0-preview.13","asset":"Mail.exe"}`)
			if mode == "incompatible" {
				marker = []byte(strings.Replace(string(marker), `"integration":1`, `"integration":2`, 1))
			}
			tag := "workspace-mail-v0.1.0-preview.13"
			asset := func(name, id string, data []byte) Asset {
				return Asset{Name: name, State: "uploaded", Size: int64(len(data)), Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(data)), APIURL: "https://api.github.com/repos/" + mailRepository + "/releases/assets/" + id, URL: "https://github.com/" + mailRepository + "/releases/download/" + tag + "/" + name}
			}
			exe := asset("Mail.exe", "1", []byte(body))
			manifest := asset(mailManifestName, "2", marker)
			if mode == "foreign-repo" {
				exe.APIURL = strings.Replace(exe.APIURL, mailRepository, "someone/mail", 1)
			}
			assets := []Asset{exe, manifest}
			if mode == "missing-manifest" {
				assets = assets[:1]
			}
			r := Release{Tag: tag, Published: time.Now(), Draft: mode == "draft", Assets: assets}
			metadata, _ := json.Marshal([]Release{r})
			downloads := 0
			client := &http.Client{Transport: mailTransport(func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "api.github.com" || !strings.HasPrefix(req.URL.Path, "/repos/"+mailRepository+"/releases") {
					t.Fatal("wrong repository", req.URL)
				}
				data := string(metadata)
				if strings.HasSuffix(req.URL.Path, "/1") {
					downloads++
					data = body
					if mode == "corrupt" {
						data = "bad"
					}
				}
				if strings.HasSuffix(req.URL.Path, "/2") {
					data = string(marker)
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(data)), Header: make(http.Header)}, nil
			})}
			p := catalog.Package{File: "Mail.exe", Size: 3, SHA256: strings.Repeat("a", 64)}
			if mode == "same" {
				p.Size = exe.Size
				p.SHA256 = strings.TrimPrefix(exe.Digest, "sha256:")
			}
			cat := catalog.Catalog{Schema: 1, Version: "0.3.0", Packages: map[string]catalog.Package{"mail": p}, Apps: []catalog.App{{ID: "mail", Name: "Mail", Version: "old", Kind: "process", Package: "mail", Executable: "Mail.exe", Port: 8996, Route: "/mail"}}}
			overrides, err := latestMail(context.Background(), client, "fixture-private-token", t.TempDir(), &cat, func(string) {})
			want := mode == "new" || mode == "same"
			if (err == nil) != want {
				t.Fatalf("error=%v", err)
			}
			if want {
				if cat.Apps[0].Version != "0.1.0-preview.13" || !strings.HasPrefix(cat.Apps[0].Source, "https://github.com/oneworkspacestudio/mail/") {
					t.Fatal("wrong provenance", cat)
				}
				if mode == "new" {
					data, e := os.ReadFile(overrides["Mail.exe"])
					if e != nil || string(data) != body || downloads != 1 {
						t.Fatal("new Mail not used", e)
					}
				} else if downloads != 0 || len(overrides) != 0 {
					t.Fatal("redownloaded identical verified binary")
				}
			}
		})
	}
}
func TestMailManifestStrict(t *testing.T) {
	for _, data := range []string{`{"schema":1,"integration":1,"version":"other","asset":"Mail.exe"}`, `{"schema":1,"integration":1,"version":"1","asset":"Mail.exe","extra":1}`, `{"schema":1,"integration":1,"version":"1","asset":"Mail.exe"} {}`} {
		if _, err := readMailManifest([]byte(data), "workspace-mail-v1"); err == nil {
			t.Fatal("accepted incompatible manifest")
		}
	}
}
