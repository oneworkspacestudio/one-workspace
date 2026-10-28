package distribution

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

type updateTransport map[string][]byte

func (tr updateTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.String() == publicIndexURL {
		index := releaseIndex{Schema: 1, PublishedAt: time.Now().UTC(), Repositories: map[string][]Release{}}
		for u, data := range tr {
			if strings.HasPrefix(u, "https://api.github.com/repos/") && strings.HasSuffix(u, "/releases?per_page=100") {
				var rs []Release
				if json.Unmarshal(data, &rs) == nil {
					index.Repositories[strings.TrimSuffix(strings.TrimPrefix(u, "https://api.github.com/repos/"), "/releases?per_page=100")] = rs
				}
			}
		}
		b, _ := json.Marshal(index)
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(b)), Request: r}, nil
	}
	body, ok := tr[r.URL.String()]
	if !ok && r.URL.Host == "github.com" {
		for _, data := range tr {
			var rs []Release
			if json.Unmarshal(data, &rs) != nil {
				continue
			}
			for _, release := range rs {
				for _, a := range release.Assets {
					if a.URL == r.URL.String() {
						body, ok = tr[a.APIURL]
					}
				}
			}
		}
	}
	status := 200
	if !ok {
		status = 404
		body = []byte(`{"message":"fixture not found"}`)
	}
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(body)), Request: r}, nil
}
func fixtureAsset(repo, tag, name, id string, body []byte) Asset {
	sum := sha256.Sum256(body)
	return Asset{Name: name, Size: int64(len(body)), Digest: "sha256:" + hex.EncodeToString(sum[:]), State: "uploaded", APIURL: "https://api.github.com/repos/" + repo + "/releases/assets/" + id, URL: "https://github.com/" + repo + "/releases/download/" + tag + "/" + name}
}
func fixturePackage(file string, body []byte) catalog.Package {
	sum := sha256.Sum256(body)
	return catalog.Package{File: file, Size: int64(len(body)), SHA256: hex.EncodeToString(sum[:])}
}
func updateCatalogFixture() catalog.Catalog {
	return catalog.Catalog{Schema: 1, Version: "0.3.0-preview.22", InstallerIntegration: 2, Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("old-manager")), "notes": fixturePackage("MyNotes.exe", []byte("old-notes")), "notes-LICENSE": fixturePackage("licenses/Notes-LICENSE", []byte("license"))}, Apps: []catalog.App{{ID: "notes", Name: "Notes", Kind: "process", Package: "notes", Executable: "MyNotes.exe", Port: 25112, Route: "/notes", Version: "0.1.0-preview.4", Source: "https://github.com/oneworkspacestudio/notes/releases/tag/v0.1.0-preview.4"}}}
}
func setupUpdateTransport(t *testing.T, withMarker bool) (catalog.Catalog, updateTransport, AppUpdateManifest) {
	t.Helper()
	old := updateCatalogFixture()
	latest := updateCatalogFixture()
	latest.Version = "0.3.0-preview.23"
	latest.Packages["manager"] = fixturePackage("Workspace.exe", []byte("new-manager"))
	raw, _ := json.Marshal(latest)
	suiteTag := "v" + latest.Version
	cat := fixtureAsset(Repository, suiteTag, "catalog.json", "1", raw)
	manager := fixtureAsset(Repository, suiteTag, "Workspace.exe", "2", []byte("new-manager"))
	bundle := fixtureAsset(Repository, suiteTag, "workspace-bundle.zip", "3", []byte("bundle"))
	suiteRelease := Release{Tag: suiteTag, Published: time.Unix(20, 0), Assets: []Asset{cat, manager, bundle}}
	oldRelease := Release{Tag: "v" + old.Version, Published: time.Unix(10, 0)}
	repo := "oneworkspacestudio/notes"
	tag := "v0.1.0-preview.5"
	app := old.Apps[0]
	app.Version = "0.1.0-preview.5"
	app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
	exe := fixtureAsset(repo, tag, "MyNotes.exe", "10", []byte("new-notes"))
	license := fixtureAsset(repo, tag, "LICENSE", "11", []byte("license"))
	manifest := AppUpdateManifest{Schema: 1, Integration: 2, Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{"notes": {Package: fixturePackage("MyNotes.exe", []byte("new-notes")), Asset: "MyNotes.exe"}, "notes-LICENSE": {Package: old.Packages["notes-LICENSE"], Asset: "LICENSE"}}}
	markerRaw, _ := json.Marshal(manifest)
	marker := fixtureAsset(repo, tag, "workspace-update.json", "12", markerRaw)
	assets := []Asset{exe, license}
	if withMarker {
		assets = append(assets, marker)
	}
	appRelease := Release{Tag: tag, Published: time.Unix(30, 0), Assets: assets}
	tr := updateTransport{cat.APIURL: raw, manager.APIURL: []byte("new-manager"), exe.APIURL: []byte("new-notes"), license.APIURL: []byte("license"), marker.APIURL: markerRaw}
	tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], _ = json.Marshal([]Release{suiteRelease, oldRelease})
	tr["https://api.github.com/repos/"+repo+"/releases?per_page=100"], _ = json.Marshal([]Release{appRelease, {Tag: "v0.1.0-preview.4", Published: time.Unix(5, 0)}})
	original := http.DefaultTransport
	http.DefaultTransport = tr
	t.Cleanup(func() { http.DefaultTransport = original })
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	return old, tr, manifest
}
func TestDiscoverSelectAndStageIndependentAppAndWorkspace(t *testing.T) {
	current, _, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	target, sources, e := SelectUpdates(current, plan, []string{"workspace", "notes"})
	if e != nil {
		t.Fatal(e)
	}
	if target.Version != "0.3.0-preview.23" || target.Apps[0].Version != "0.1.0-preview.5" || target.Apps[0].Port != 25112 {
		t.Fatal(target)
	}
	if current.Version != "0.3.0-preview.22" || current.Packages["notes"].SHA256 == target.Packages["notes"].SHA256 {
		t.Fatal("installed catalog mutated")
	}
	if len(sources) != 2 {
		t.Fatal("unchanged license downloaded", sources)
	}
	dir := t.TempDir()
	if e = StageUpdateFiles(context.Background(), target, sources, dir, func(string) {}); e != nil {
		t.Fatal(e)
	}
	if sum, _, e := digestFixture(filepath.Join(dir, "packages", "MyNotes.exe")); e != nil || sum != target.Packages["notes"].SHA256 {
		t.Fatal(sum, e)
	}
}

func TestInstalledAppAbsentFromPublishedSuiteKeepsIndependentBaseline(t *testing.T) {
	for _, newer := range []bool{false, true} {
		t.Run(fmt.Sprintf("newerRelease=%v", newer), func(t *testing.T) {
			current, tr, _ := setupUpdateTransport(t, true)
			endpoint := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
			var releases []Release
			if err := json.Unmarshal(tr[endpoint], &releases); err != nil {
				t.Fatal(err)
			}
			var published catalog.Catalog
			asset, _ := releaseAsset(releases[0], "catalog.json")
			if err := json.Unmarshal(tr[asset.APIURL], &published); err != nil {
				t.Fatal(err)
			}
			published.Apps = nil
			delete(published.Packages, "notes")
			delete(published.Packages, "notes-LICENSE")
			raw, _ := json.Marshal(published)
			tr[asset.APIURL] = raw
			for i, a := range releases[0].Assets {
				if a.Name == "catalog.json" {
					releases[0].Assets[i] = fixtureAsset(Repository, releases[0].Tag, "catalog.json", "1", raw)
				}
			}
			tr[endpoint], _ = json.Marshal(releases)
			if !newer {
				tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"], _ = json.Marshal([]Release{{Tag: "v0.1.0-preview.4", Published: time.Unix(5, 0)}})
			}
			plan, err := CheckUpdates(context.Background(), current)
			if err != nil {
				t.Fatal(err)
			}
			for _, choice := range plan.Choices {
				if choice.ID != "notes" {
					continue
				}
				want := "0.1.0-preview.4"
				if newer {
					want = "0.1.0-preview.5"
				}
				if choice.Version != want || choice.Available != newer || choice.Reason != "" || len(choice.Apps) != 1 || len(choice.Packages) != 2 {
					t.Fatalf("unexpected update choice: %+v", choice)
				}
				if newer {
					target, _, err := SelectUpdates(current, plan, []string{"notes"})
					if err != nil || len(target.Apps) != 1 || target.Apps[0].Version != want {
						t.Fatalf("independent update lost app: %+v, %v", target, err)
					}
				}
				return
			}
			t.Fatal("installed app missing from update plan")
		})
	}
}
func digestFixture(path string) (string, int64, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", 0, e
	}
	defer f.Close()
	h := sha256.New()
	n, e := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), n, e
}
func TestUncontractedNewReleaseIsVisibleButNotInstalled(t *testing.T) {
	current, _, _ := setupUpdateTransport(t, false)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	var notes UpdateChoice
	for _, ch := range plan.Choices {
		if ch.ID == "notes" {
			notes = ch
		}
	}
	if notes.Available || !strings.Contains(notes.Reason, "v0.1.0-preview.5") {
		t.Fatal(notes)
	}
	if _, _, e = SelectUpdates(current, plan, []string{"notes"}); e == nil {
		t.Fatal("uncontracted executable selected")
	}
}
func TestIndependentContractCannotReplaceOtherAppsOrEscapeArchive(t *testing.T) {
	current, tr, m := setupUpdateTransport(t, true)
	var releases []Release
	json.Unmarshal(tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"], &releases)
	r := releases[0]
	g := updateGroup{}
	for _, candidate := range updateGroups {
		if candidate.id == "notes" {
			g = candidate
		}
	}
	for _, mutate := range []func(*AppUpdateManifest){func(m *AppUpdateManifest) { m.Repository = "someone/notes" }, func(m *AppUpdateManifest) { m.Apps[0].ID = "accounts" }, func(m *AppUpdateManifest) {
		p := m.Packages["notes"]
		p.File = "MyAccounts.exe"
		m.Packages["notes"] = p
	}, func(m *AppUpdateManifest) {
		p := m.Packages["notes"]
		p.Member = "../outside.exe"
		m.Packages["notes"] = p
	}, func(m *AppUpdateManifest) { delete(m.Packages, "notes-LICENSE") }} {
		raw, _ := json.Marshal(m)
		var copy AppUpdateManifest
		json.Unmarshal(raw, &copy)
		mutate(&copy)
		raw, _ = json.Marshal(copy)
		if _, e := manifestChoice(g, current, "oneworkspacestudio/notes", r, raw); e == nil {
			t.Fatal("unsafe contract accepted", string(raw))
		}
	}
}
func TestStagedDownloadChecksumAndArchiveMemberChecks(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	target, sources, e := SelectUpdates(current, plan, []string{"notes"})
	if e != nil {
		t.Fatal(e)
	}
	d := sources["notes"]
	tr[d.Asset.APIURL] = []byte("tampered")
	if e = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {}); e == nil {
		t.Fatal("tampered executable staged")
	}
	for _, duplicate := range []bool{false, true} {
		var buf bytes.Buffer
		z := zip.NewWriter(&buf)
		name := "nested/MyNotes.exe"
		f, _ := z.Create(name)
		f.Write([]byte("new-notes"))
		if duplicate {
			f, _ = z.Create(name)
			f.Write([]byte("new-notes"))
		}
		z.Close()
		a := fixtureAsset("oneworkspacestudio/notes", "v0.1.0-preview.5", "portable.zip", "99", buf.Bytes())
		tr[a.APIURL] = buf.Bytes()
		tr[a.URL] = buf.Bytes()
		d.Asset = a
		d.Member = name
		sources["notes"] = d
		e = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {})
		if duplicate && e == nil {
			t.Fatal("duplicate member accepted")
		}
		if !duplicate && e != nil {
			t.Fatal(e)
		}
	}
}
func TestWorkspaceVersionAndSelectionGuards(t *testing.T) {
	for _, tc := range []struct {
		next, old string
		want      bool
	}{{"0.3.0-preview.23", "0.3.0-preview.20-fixes", true}, {"0.3.0-preview.22", "0.3.0-preview.23", false}, {"0.4.0", "0.3.0-preview.23", true}, {"0.3.0-preview.23", "development", false}} {
		if got := newerWorkspace(tc.next, tc.old); got != tc.want {
			t.Fatal(tc, got)
		}
	}
	current, _, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	for _, selection := range [][]string{{"notes", "notes"}, {"unknown"}, {}} {
		if _, _, e = SelectUpdates(current, plan, selection); e == nil {
			t.Fatal("invalid selection accepted", selection)
		}
	}
	changed := current
	changed.Version = "0.3.0-preview.24"
	if _, _, e = SelectUpdates(changed, plan, []string{"workspace"}); e == nil {
		t.Fatal("stale plan could downgrade installation")
	}
}
func TestPublisherContractRoundTripAndCustomExtensionOwnership(t *testing.T) {
	current, tr, m := setupUpdateTransport(t, true)
	current.Packages["notes-addon"] = fixturePackage("CustomNotes.exe", []byte("custom extension"))
	current.Apps = append(current.Apps, catalog.App{ID: "notes-addon", Name: "Custom notes helper", Version: "1.0.0", Kind: "process", Package: "notes-addon", Executable: "CustomNotes.exe", Port: 25119, Route: "/custom-notes"})
	generated, e := MakeAppUpdateManifest(current, "notes", "0.1.0-preview.5", "oneworkspacestudio/notes", "v0.1.0-preview.5", m.Packages)
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(generated)
	var releases []Release
	json.Unmarshal(tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"], &releases)
	var notesGroup updateGroup
	for _, group := range updateGroups {
		if group.id == "notes" {
			notesGroup = group
		}
	}
	ch, e := manifestChoice(notesGroup, current, generated.Repository, releases[0], raw)
	if e != nil {
		t.Fatal(e)
	}
	ch.Available = true
	plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{ch}}
	target, _, e := SelectUpdates(current, plan, []string{"notes"})
	if e != nil {
		t.Fatal(e)
	}
	if target.Packages["notes-addon"] != current.Packages["notes-addon"] || len(target.Apps) != 2 {
		t.Fatal("custom extension was replaced")
	}
}

func setupCompatibilityUpdate(t *testing.T, tracked bool) (catalog.Catalog, updateTransport, AppUpdateManifest, Release) {
	t.Helper()
	repo, tag := "oneworkspacestudio/my-social", "v0.1.0-preview.6"
	current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.82", InstallerIntegration: 3, Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("manager")), "social": fixturePackage("MySocial.exe", []byte("old social"))}, Apps: []catalog.App{{ID: "social", Name: "My Social", Kind: "process", Package: "social", Executable: "MySocial.exe", Port: 29013, Route: "/social", Version: "0.1.0-preview.5", Source: "https://github.com/" + repo + "/releases/tag/v0.1.0-preview.5"}}}
	document := []byte(`{"schema":1}`)
	if tracked {
		current.Packages["social-contract"] = fixturePackage("docs/MySocial/compatibility.json", []byte("old document"))
	}
	app := current.Apps[0]
	app.Version = "0.1.0-preview.6"
	app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
	exe := fixtureAsset(repo, tag, "MySocial.exe", "20", []byte("new social"))
	doc := fixtureAsset(repo, tag, "compatibility.json", "21", document)
	m := AppUpdateManifest{Schema: 1, Integration: 2, Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{"social": {Package: fixturePackage("MySocial.exe", []byte("new social")), Asset: exe.Name}, "social-contract": {Package: fixturePackage("docs/MySocial/compatibility.json", document), Asset: doc.Name}}}
	raw, _ := json.Marshal(m)
	marker := fixtureAsset(repo, tag, "workspace-update.json", "22", raw)
	release := Release{Tag: tag, Published: time.Unix(30, 0), Assets: []Asset{exe, doc, marker}}
	suiteTag := "v" + current.Version
	catalogRaw, _ := json.Marshal(current)
	cat := fixtureAsset(Repository, suiteTag, "catalog.json", "1", catalogRaw)
	bundle := fixtureAsset(Repository, suiteTag, "workspace-bundle.zip", "2", []byte("bundle"))
	tr := updateTransport{cat.APIURL: catalogRaw, exe.APIURL: []byte("new social"), doc.APIURL: document, marker.APIURL: raw}
	tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], _ = json.Marshal([]Release{{Tag: suiteTag, Published: time.Unix(20, 0), Assets: []Asset{cat, bundle}}})
	tr["https://api.github.com/repos/"+repo+"/releases?per_page=100"], _ = json.Marshal([]Release{release, {Tag: "v0.1.0-preview.5", Published: time.Unix(5, 0)}})
	original := http.DefaultTransport
	http.DefaultTransport = tr
	t.Cleanup(func() { http.DefaultTransport = original })
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	return current, tr, m, release
}

func TestIndependentUpdateWithOptionalUntrackedCompatibilityDocument(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		t.Run(fmt.Sprint("tracked=", tracked), func(t *testing.T) {
			current, _, _, _ := setupCompatibilityUpdate(t, tracked)
			plan, e := CheckUpdates(context.Background(), current)
			if e != nil {
				t.Fatal(e)
			}
			var choice UpdateChoice
			for _, c := range plan.Choices {
				if c.ID == "social" {
					choice = c
				}
			}
			if !choice.Available || choice.Reason != "" || choice.ReviewURL != "" {
				t.Fatal("ordinary update was blocked", choice)
			}
			target, sources, e := SelectUpdates(current, plan, []string{"social"})
			if e != nil {
				t.Fatal(e)
			}
			if target.Apps[0].Version != "0.1.0-preview.6" || target.Version != current.Version || target.Apps[0].Port != current.Apps[0].Port {
				t.Fatal("wrong update target")
			}
			if _, ok := target.Packages["social-contract"]; ok != tracked {
				t.Fatal("optional document registration changed")
			}
			if _, ok := sources["social-contract"]; ok != tracked {
				t.Fatal("optional document download changed")
			}
			if target.Packages["manager"] != current.Packages["manager"] {
				t.Fatal("manager package changed")
			}
			if e = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {}); e != nil {
				t.Fatal(e)
			}
		})
	}
}

func TestUntrackedCompatibilityDocumentDoesNotRelaxPackageBoundaries(t *testing.T) {
	current, _, m, release := setupCompatibilityUpdate(t, false)
	delete(m.Packages, "social-contract")
	g := updateGroup{"social", "My Social", "my-social", []string{"social"}}
	for _, tc := range []struct{ id, file string }{
		{"clips-contract", "docs/MyClips/compatibility.json"},
		{"social-contract", "MySocial.exe"},
		{"social-runtime", "docs/MySocial/compatibility.json"},
		{"social-contract", "docs/Other/compatibility.json"},
		{"social-contract", "docs/MySocial/../compatibility.json"},
	} {
		t.Run(tc.id+"/"+tc.file, func(t *testing.T) {
			copy := m
			copy.Packages = map[string]UpdatePackage{}
			for id, p := range m.Packages {
				copy.Packages[id] = p
			}
			copy.Packages[tc.id] = UpdatePackage{Package: fixturePackage(tc.file, []byte("unexpected")), Asset: "unexpected"}
			raw, _ := json.Marshal(copy)
			if _, e := manifestChoice(g, current, m.Repository, release, raw); e == nil {
				t.Fatal("unregistered runtime/foreign/path-changing package accepted")
			}
		})
	}
	current.Packages["social-contract"] = fixturePackage("docs/MySocial/compatibility.json", []byte("old"))
	raw, _ := json.Marshal(m)
	if _, e := manifestChoice(g, current, m.Repository, release, raw); e == nil {
		t.Fatal("tracked document may not disappear")
	}
	m.Packages["social-contract"] = UpdatePackage{Package: fixturePackage("docs/Other/compatibility.json", []byte("moved")), Asset: "unexpected"}
	raw, _ = json.Marshal(m)
	if _, e := manifestChoice(g, current, m.Repository, release, raw); e == nil {
		t.Fatal("tracked document may not move")
	}
}

func TestComponentReleaseWithIdenticalNativeBytesIsNotAnotherAppUpdate(t *testing.T) {
	current, _, m, release := setupCompatibilityUpdate(t, false)
	current.Apps[0].Version = m.Apps[0].Version
	current.Packages["social"] = m.Packages["social"].Package
	// Only the release link differs, as for a Chrome extension component patch.
	raw, _ := json.Marshal(m)
	g := updateGroup{"social", "My Social", "my-social", []string{"social"}}
	choice, e := manifestChoice(g, current, m.Repository, release, raw)
	if e != nil {
		t.Fatal(e)
	}
	if choiceChanged(g, current, choice) {
		t.Fatal("identical native app offered another update")
	}
	choice.Apps[0].Args = []string{"--changed-startup"}
	if !choiceChanged(g, current, choice) {
		t.Fatal("a real launch-contract change was hidden")
	}
}
