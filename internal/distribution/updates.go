package distribution

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// AppUpdateManifest is a publisher's explicit Windows managed-integration
// contract. It travels with the app release, verified by GitHub's asset digest.
// Arbitrary release executables are never guessed into an installation.
type AppUpdateManifest struct {
	Schema               int                      `json:"schema"`
	Integration          int                      `json:"integration"`
	Repository           string                   `json:"repository"`
	Tag                  string                   `json:"tag"`
	Apps                 []catalog.App            `json:"apps"`
	Packages             map[string]UpdatePackage `json:"packages"`
	MinimumManager       string                   `json:"minimumManager,omitempty"`
	RequiredCapabilities []string                 `json:"requiredCapabilities,omitempty"`
}
type UpdatePackage struct {
	catalog.Package
	Asset  string `json:"asset"`
	Member string `json:"member,omitempty"`
}
type UpdateDownload struct {
	Asset      Asset  `json:"asset"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Member     string `json:"member,omitempty"`
}
type UpdateChoice struct {
	ID                   string                     `json:"id"`
	Name                 string                     `json:"name"`
	Current              string                     `json:"current"`
	Version              string                     `json:"version"`
	Source               string                     `json:"source"`
	Available            bool                       `json:"available"`
	Reason               string                     `json:"reason,omitempty"`
	ReviewURL            string                     `json:"reviewUrl,omitempty"`
	NewApps              []string                   `json:"newApps,omitempty"`
	Apps                 []catalog.App              `json:"apps,omitempty"`
	Packages             map[string]catalog.Package `json:"packages,omitempty"`
	Downloads            map[string]UpdateDownload  `json:"downloads,omitempty"`
	MinimumManager       string                     `json:"minimumManager,omitempty"`
	RequiredCapabilities []string                   `json:"requiredCapabilities,omitempty"`
}
type UpdatePlan struct {
	IndexPublishedAt  time.Time         `json:"indexPublishedAt,omitempty"`
	Installer         *InstallerRelease `json:"installer,omitempty"`
	CheckedAt         time.Time         `json:"checkedAt"`
	BaseCatalogSHA256 string            `json:"baseCatalogSha256"`
	Choices           []UpdateChoice    `json:"choices"`
}

// Fresh bounds how long a reviewed release can be used without another check.
func (p UpdatePlan) Fresh(now time.Time) bool {
	return !p.CheckedAt.IsZero() && !p.CheckedAt.After(now.Add(time.Minute)) && now.Sub(p.CheckedAt) <= 30*time.Minute
}

func catalogIdentity(c catalog.Catalog) string {
	b, _ := json.Marshal(c)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

type updateGroup struct {
	id, name, repo string
	apps           []string
}

var updateGroups = []updateGroup{
	{"finance", "My Finance", "my-finance", []string{"finance"}},
	{"shield", "My Shield", "my-shield", []string{"shield"}},
	{"forms", "My Forms", "my-forms", []string{"forms"}},
	{"push", "My Push", "my-push", []string{"push"}},
	{"maps", "My Maps", "my-maps", []string{"maps"}},
	{"cinema", "My Cinema", "my-cinema", []string{"cinema"}},
	{"myai", "My AI", "my-ai", []string{"myai"}},
	{"my3d", "My 3D", "my-3d", []string{"my3d"}},
	{"create", "My Create", "my-create", []string{"create"}},
	{"meet", "My Meet", "my-meet", []string{"meet"}},

	{"workspace", "One Workspace", "one-workspace", nil},
	{"accounts", "MyAccounts", "myaccounts", []string{"accounts"}},
	{"drive", "Drive, Docs and Photos", "drive", []string{"drive", "docs", "photos"}},
	{"gitea", "Git", "gitea", []string{"gitea"}},
	{"music", "Music", "Music", []string{"music"}},
	{"mail", "Mail", "mail", []string{"mail"}},
	{"studio", "Photo Editor", "Photo-Editor", []string{"studio"}},
	{"calendar", "Calendar and Contacts", "calendar", []string{"calendar", "contacts"}},
	{"notes", "Notes", "notes", []string{"notes"}},
	{"my-tube", "My Tube", "my-tube", []string{"my-tube"}},
	{"vault", "My Vault", "my-vault", []string{"vault"}},
	{"chat", "Chat", "chat", []string{"chat"}},
	{"hosting", "My Hosting", "my-hosting", []string{"hosting"}},
	{"tasks", "My Tasks", "my-tasks", []string{"tasks"}},
	{"pdf", "My PDF", "my-pdf", []string{"pdf"}},
	{"monitor", "My Monitor", "my-monitor", []string{"monitor"}},
	{"search", "My Search", "my-search", []string{"search"}},
	{"slicer", "My Slicer", "my-slicer", []string{"slicer"}},
	{"clips", "My Clips", "my-clips", []string{"clips"}},
	{"social", "My Social", "my-social", []string{"social"}},
}

const withheldMusicSHA = "5f0d121cd9733775b327ac3b39b46b32cf5d9cbdef1b014880a89d6421239837"

func releasesFor(ctx context.Context, h *http.Client, token, repo string) ([]Release, error) {
	if index, ok := ctx.Value(indexContextKey{}).(*releaseIndex); ok {
		rs, exists := index.Repositories[repo]
		if !exists {
			return nil, fmt.Errorf("Published update list is missing %s; contact the Workspace publisher", repo)
		}
		copy := append([]Release(nil), rs...)
		sort.SliceStable(copy, func(i, j int) bool { return copy[i].Published.After(copy[j].Published) })
		return copy, nil
	}
	req, e := request(ctx, "https://api.github.com/repos/"+repo+"/releases?per_page=100", token)
	if e != nil {
		return nil, e
	}
	res, e := h.Do(req)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, releaseHTTPError(repo, res)
	}
	var rs []Release
	e = json.NewDecoder(io.LimitReader(res.Body, 4<<20)).Decode(&rs)
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].Published.After(rs[j].Published) })
	return rs, e
}
func releaseAsset(r Release, name string) (Asset, bool) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, true
		}
	}
	return Asset{}, false
}
func fetchUpdateBytes(ctx context.Context, h *http.Client, token, repo string, r Release, a Asset, max int64) ([]byte, error) {
	if !validNamedAsset(a, repo, r.Tag, a.Name, max) {
		return nil, errors.New("invalid update release asset")
	}
	req, e := request(ctx, assetEndpoint(a, token), token)
	if e != nil {
		return nil, e
	}
	req.Header.Set("Accept", "application/octet-stream")
	res, e := h.Do(req)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("asset download failed (HTTP %d)", res.StatusCode)
	}
	b, e := io.ReadAll(io.LimitReader(res.Body, a.Size+1))
	if e != nil {
		return nil, e
	}
	sum := sha256.Sum256(b)
	if int64(len(b)) != a.Size || "sha256:"+hex.EncodeToString(sum[:]) != a.Digest {
		return nil, errors.New("update metadata failed checksum verification")
	}
	return b, nil
}
func sourceTag(source, repo string) string {
	return strings.TrimPrefix(source, "https://github.com/"+repo+"/releases/tag/")
}
func releaseTime(rs []Release, tag string) time.Time {
	for _, r := range rs {
		if r.Tag == tag && !r.Draft {
			return r.Published
		}
	}
	return time.Time{}
}
func newerWorkspace(candidate, current string) bool {
	re := regexp.MustCompile(`[0-9]+`)
	a, b := re.FindAllString(candidate, -1), re.FindAllString(current, -1)
	if len(a) < 3 || len(b) < 3 {
		return false
	}
	for i := 0; i < len(a) || i < len(b); i++ {
		if i == 3 {
			cs, os := !strings.Contains(candidate, "-"), !strings.Contains(current, "-")
			if cs != os {
				return cs
			}
		}
		x, y := "0", "0"
		if i < len(a) {
			x = strings.TrimLeft(a[i], "0")
		}
		if i < len(b) {
			y = strings.TrimLeft(b[i], "0")
		}
		if len(x) != len(y) {
			return len(x) > len(y)
		}
		if x != y {
			return x > y
		}
	}
	return false
}
func UpdateControllerSupported(version string) bool {
	return newerWorkspace(version, "0.3.0-preview.22")
}

// Schema3 recovery must never cache a legacy controller that cannot read it.
func MetadataRecoverySupported(version string) bool {
	return version == "0.3.0-preview.91" || newerWorkspace(version, "0.3.0-preview.91")
}
func groupPackages(g updateGroup, c catalog.Catalog) map[string]catalog.Package {
	out := map[string]catalog.Package{}
	fixedGroup := false
	for _, existing := range updateGroups {
		if existing.id == g.id {
			fixedGroup = true
			break
		}
	}
	for id, p := range c.Packages {
		if updatePackageOwner(id) == g.id {
			out[id] = p
		}
		// Forward-compatible suite additions use the publisher's app package
		// namespace. Existing fixed coupled groups retain their ownership rules.
		if !fixedGroup && updatePackageOwner(id) == "" {
			for _, aid := range g.apps {
				if a, ok := c.Find(aid); ok && (id == a.Package || strings.HasPrefix(id, aid+"-")) {
					out[id] = p
				}
			}
		}
	}
	return out
}

func suiteUpdateGroups(latest catalog.Catalog) []updateGroup {
	groups := append([]updateGroup(nil), updateGroups...)
	known := map[string]bool{}
	for _, g := range groups {
		for _, id := range g.apps {
			known[id] = true
		}
	}
	for _, a := range latest.Apps {
		if known[a.ID] {
			continue
		}
		u, err := url.Parse(a.Source)
		if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
			continue
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) != 5 || parts[0] != "oneworkspacestudio" || parts[2] != "releases" || parts[3] != "tag" || parts[4] == "" {
			continue
		}
		groups = append(groups, updateGroup{a.ID, a.Name, parts[1], []string{a.ID}})
	}
	return groups
}
func updatePackageOwner(id string) string {
	if id == "finance" || strings.HasPrefix(id, "finance-") {
		return "finance"
	}
	if id == "shield" || strings.HasPrefix(id, "shield-") {
		return "shield"
	}
	if id == "forms" || strings.HasPrefix(id, "forms-") {
		return "forms"
	}
	if id == "push" || strings.HasPrefix(id, "push-") {
		return "push"
	}
	if strings.HasPrefix(id, "msvc-") {
		return "workspace"
	}
	for _, app := range []string{"cinema", "myai", "my3d", "create", "maps", "meet"} {
		if id == app || strings.HasPrefix(id, app+"-") {
			return app
		}
	}
	if id == "clips" || strings.HasPrefix(id, "clips-") {
		return "clips"
	}
	if id == "social" || strings.HasPrefix(id, "social-") {
		return "social"
	}
	if id == "slicer" || strings.HasPrefix(id, "slicer-") {
		return "slicer"
	}
	if id == "tasks" || id == "tasks-LICENSE" || id == "tasks-NOTICE" || id == "tasks-THIRD-PARTY-LICENSES.txt" {
		return "tasks"
	}
	if id == "pdf" || id == "pdf-LICENSE" || id == "pdf-NOTICE" || id == "pdf-THIRD-PARTY-LICENSES.txt" {
		return "pdf"
	}
	if id == "monitor" || id == "monitor-LICENSE" || id == "monitor-NOTICE" || id == "monitor-THIRD-PARTY-LICENSES.txt" {
		return "monitor"
	}
	if id == "search" || id == "search-LICENSE" || id == "search-NOTICE" || id == "search-THIRD-PARTY-LICENSES.txt" {
		return "search"
	}

	if id == "hosting" || strings.HasPrefix(id, "hosting-") {
		return "hosting"
	}
	if id == "vault" || strings.HasPrefix(id, "vault-") {
		return "vault"
	}
	fixed := map[string]string{"chat": "chat", "chat-LICENSE": "chat", "chat-NOTICE": "chat", "chat-THIRD-PARTY-LICENSES.txt": "chat", "manager": "workspace", "cloudflared": "workspace", "license": "workspace", "notice": "workspace", "accounts": "accounts", "drive": "drive", "photos": "drive", "gitea": "gitea", "music": "music", "mail": "mail", "mail-license": "mail", "mail-notice": "mail", "studio": "studio", "studio-license": "studio", "services-shared": "calendar", "stalwart-shared": "calendar", "calendar": "calendar", "contacts": "calendar", "calendar-license": "calendar", "notes": "notes", "notes-LICENSE": "notes", "notes-NOTICE": "notes", "notes-THIRD-PARTY-LICENSES.txt": "notes", "my-tube": "my-tube", "my-tube-LICENSE": "my-tube", "my-tube-NOTICE": "my-tube", "my-tube-THIRD-PARTY-LICENSES.txt": "my-tube"}
	if owner := fixed[id]; owner != "" {
		return owner
	}
	for prefix, owner := range map[string]string{"notice-Calendar-": "calendar", "notice-MyAccounts-": "accounts", "notice-FileBrowser-": "drive", "notice-Editor-": "drive", "notice-Vue-": "drive", "notice-Gitea-": "gitea", "notice-SwingMusic-": "music", "notice-Imageio-ffmpeg-": "music", "notice-FFmpeg-": "music", "notice-Photo-Editor-": "studio"} {
		if strings.HasPrefix(id, prefix) {
			return owner
		}
	}
	if strings.HasPrefix(id, "notice-") {
		return "workspace"
	}
	return ""
}
func groupApps(g updateGroup, c catalog.Catalog) []catalog.App {
	var out []catalog.App
	for _, id := range g.apps {
		if a, ok := c.Find(id); ok {
			out = append(out, a)
		}
	}
	return out
}
func normalizedApps(apps []catalog.App) []catalog.App {
	out := append([]catalog.App(nil), apps...)
	for i := range out {
		out[i].Port = 0
		// Component-only releases can point at identical native app bytes.
		// A different release link alone is not an app update.
		out[i].Source = ""
	}
	return out
}

// Only the verified Workspace catalog can introduce app IDs. Standalone app
// manifests remain limited to their existing managed package contract.
func newDefaultIDs(current, latest catalog.Catalog) map[string]bool {
	out := map[string]bool{}
	for _, a := range latest.WantedApps(current) {
		if _, installed := current.Find(a.ID); !installed {
			out[a.ID] = true
		}
	}
	return out
}
func markNewDefaults(current, latest catalog.Catalog, ch UpdateChoice) UpdateChoice {
	allowed := newDefaultIDs(current, latest)
	for _, a := range ch.Apps {
		if allowed[a.ID] {
			ch.NewApps = append(ch.NewApps, a.ID)
		}
	}
	return ch
}
func scopeUpdateChoice(current catalog.Catalog, ch UpdateChoice) UpdateChoice {
	apps := []catalog.App{}
	selected := map[string]bool{}
	for _, a := range current.Apps {
		selected[a.ID] = true
	}
	for _, id := range ch.NewApps {
		if !current.EffectiveSelection().Excludes(id) {
			selected[id] = true
		}
	}
	for _, app := range ch.Apps {
		if selected[app.ID] {
			apps = append(apps, app)
		}
	}
	ch.Apps = apps
	keep := catalog.InstallPackageIDs(catalog.Catalog{Apps: apps, Packages: ch.Packages}, selected)
	hasNew := len(ch.NewApps) > 0
	for id := range ch.Packages {
		prerequisite := false
		if ch.ID == "workspace" && strings.HasPrefix(id, "msvc-") {
			for _, app := range current.Apps {
				if app.ID == "calendar" || app.ID == "contacts" || app.ID == "mail" {
					prerequisite = true
				}
			}
		}
		if _, ok := current.Packages[id]; !ok && !(hasNew && keep[id]) && !prerequisite {
			delete(ch.Packages, id)
			delete(ch.Downloads, id)
		}
	}
	// Keep the full target contract, but advertise only downloads SelectUpdates
	// will actually stage. A Photos-only change must not count an unchanged
	// Drive/Office executable in the selected download size.
	for id := range ch.Downloads {
		if old, ok := current.Packages[id]; ok && old == ch.Packages[id] {
			delete(ch.Downloads, id)
		}
	}
	return ch
}
func choiceChanged(g updateGroup, current catalog.Catalog, ch UpdateChoice) bool {
	if g.id == "workspace" {
		return current.Version != ch.Version || current.Packages["manager"] != ch.Packages["manager"]
	}
	return !reflect.DeepEqual(groupPackages(g, current), ch.Packages) || !reflect.DeepEqual(normalizedApps(groupApps(g, current)), normalizedApps(ch.Apps))
}

// CheckUpdates downloads only small catalogs/manifests. An incomplete repository
// check is reported and cannot silently substitute an older candidate.
func CheckUpdates(ctx context.Context, current catalog.Catalog) (UpdatePlan, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	totalGroups := 0
	for _, g := range updateGroups {
		if g.id == "workspace" || len(groupApps(g, current)) > 0 {
			totalGroups++
		}
	}
	operationprogress.Report(ctx, "check", "Checking One Workspace releases", 0, int64(totalGroups), "repositories")
	h, token := client(), githubToken(ctx)
	ctx, indexErr := withReleaseIndex(ctx, token)
	if indexErr != nil {
		return UpdatePlan{}, indexErr
	}
	plan := UpdatePlan{IndexPublishedAt: indexPublishedAt(ctx), CheckedAt: time.Now().UTC(), BaseCatalogSHA256: catalogIdentity(current)}
	rs, e := releasesFor(ctx, h, token, Repository)
	if e != nil {
		return plan, e
	}
	suiteRelease, bundle, e := choose(rs)
	if e != nil {
		return plan, e
	}
	ca, ok := releaseAsset(suiteRelease, "catalog.json")
	if !ok {
		return plan, errors.New("published Workspace catalog missing")
	}
	raw, e := fetchUpdateBytes(ctx, h, token, Repository, suiteRelease, ca, 1<<20)
	if e != nil {
		return plan, e
	}
	latest, e := catalog.Read(raw)
	if e != nil {
		if requiresInstaller(raw) {
			plan.Installer = installerRelease(suiteRelease)
			if plan.Installer != nil {
				var header struct {
					Version string `json:"version"`
				}
				if json.Unmarshal(raw, &header) == nil && header.Version == plan.Installer.Version {
					plan.Installer.Required = true
					return plan, nil
				}
			}
		}
		return plan, e
	}
	plan.Installer = installerRelease(suiteRelease)
	if latest.Version != strings.TrimPrefix(suiteRelease.Tag, "v") || !validAsset(bundle, Repository, suiteRelease.Tag) {
		return plan, errors.New("invalid Workspace release version or bundle")
	}
	totalGroups = 1
	newIDs := newDefaultIDs(current, latest)
	for _, g := range suiteUpdateGroups(latest) {
		if g.id != "workspace" {
			included := len(groupApps(g, current)) > 0
			for _, id := range g.apps {
				included = included || newIDs[id]
			}
			if included {
				totalGroups++
			}
		}
	}
	checked := 0
	for _, g := range suiteUpdateGroups(latest) {
		oldApps := groupApps(g, current)
		included := g.id == "workspace" || len(oldApps) > 0
		for _, id := range g.apps {
			included = included || newIDs[id]
		}
		if !included {
			continue
		}
		operationprogress.Report(ctx, "check", "Checking "+g.name+" releases", int64(checked), int64(totalGroups), "repositories")
		checked++
		ch := UpdateChoice{ID: g.id, Name: g.name, Current: current.Version, Version: latest.Version, Source: "https://github.com/" + Repository + "/releases/tag/" + suiteRelease.Tag, Apps: groupApps(g, latest), Packages: groupPackages(g, latest), Downloads: map[string]UpdateDownload{}}
		repo := "oneworkspacestudio/" + g.repo
		var appReleases []Release
		ch = markNewDefaults(current, latest, ch)
		if g.id != "workspace" && len(oldApps) == 0 {
			if len(ch.Apps) == 0 {
				continue
			}
			ch.Current = "Not installed"
			ch.Version, ch.Source = ch.Apps[0].Version, ch.Apps[0].Source
		} else if g.id != "workspace" {
			ch.Current = oldApps[0].Version
			// A locally integrated app can be newer than the public suite catalog.
			// Its absence from that catalog is not an update or an uninstall.
			// Keep its installed pin as the baseline for independent releases.
			if len(ch.Apps) == 0 {
				ch.Apps = oldApps
				ch.Packages = groupPackages(g, current)
				ch.Version = ch.Current
				ch.Source = oldApps[0].Source
			}
			appReleases, e = releasesFor(ctx, h, token, repo)
			if e != nil {
				ch.Reason = e.Error()
				plan.Choices = append(plan.Choices, ch)
				continue
			}
			if len(ch.Apps) > 0 {
				ch.Version = ch.Apps[0].Version
				ch.Source = ch.Apps[0].Source
				oldTime := releaseTime(appReleases, sourceTag(oldApps[0].Source, repo))
				newTime := releaseTime(appReleases, sourceTag(ch.Source, repo))
				if ch.Source != oldApps[0].Source && (oldTime.IsZero() || !newTime.After(oldTime)) {
					if len(ch.NewApps) > 0 {
						ch.Reason = "New apps require a compatible suite pin for the newer installed app; no downgrade was selected."
						plan.Choices = append(plan.Choices, ch)
						continue
					}
					ch.Apps = oldApps
					ch.Packages = groupPackages(g, current)
					ch.Version = ch.Current
					ch.Source = oldApps[0].Source
				}
			}
		} else if !newerWorkspace(latest.Version, current.Version) {
			ch.Packages = groupPackages(g, current)
			ch.Version = current.Version
		}
		for id, p := range ch.Packages {
			d := UpdateDownload{Asset: bundle, Repository: Repository, Tag: suiteRelease.Tag, Member: p.File}
			// A small notice must not force a full application-bundle download.
			// Use a published direct asset only when its digest and length match
			// this exact catalog package; otherwise retain the verified bundle.
			direct := map[string]string{"manager": "Workspace.exe", "license": "LICENSE.txt", "notice": "NOTICE.txt", "notice-Workspace-release-status.md": "KNOWN-LIMITATIONS.md", "notice-organization-components.json": "organization-components.json"}
			if name := direct[id]; name != "" {
				if a, ok := releaseAsset(suiteRelease, name); ok && validNamedAsset(a, Repository, suiteRelease.Tag, name, 2<<30) && a.Digest == "sha256:"+p.SHA256 && a.Size == p.Size {
					d.Asset = a
					d.Member = ""
				}
			}
			ch.Downloads[id] = d
		}
		// An explicit component contract also supplies smaller download assets
		// when its tag equals the suite pin. In that case retain the suite target
		// and bind only assets with identical package path, size and checksum.
		for _, r := range appReleases {
			if len(ch.NewApps) > 0 {
				break
			} // new integration comes from the suite contract
			if r.Draft || r.Published.IsZero() {
				continue
			}
			pin := sourceTag(ch.Source, repo)
			baseline := releaseTime(appReleases, pin)
			samePin := r.Tag == pin
			if baseline.IsZero() || (!r.Published.After(baseline) && !samePin) {
				break
			}
			if samePin {
				changed := false
				for id, pkg := range ch.Packages {
					changed = changed || current.Packages[id] != pkg
				}
				if !changed {
					break // No component bytes need a download location.
				}
			}
			marker, ok := releaseAsset(r, "workspace-update.json")
			if !ok {
				if samePin {
					break // Legacy suite components may have no standalone contract.
				}
				ch.Reason = "Newer release " + r.Tag + " has no managed-update compatibility metadata. It will be available automatically when included in a compatible Workspace bundle."
				ch.ReviewURL = "https://github.com/" + repo + "/releases/tag/" + r.Tag
				break
			}
			data, err := fetchUpdateBytes(ctx, h, token, repo, r, marker, 1<<20)
			if err != nil {
				ch.Reason = err.Error()
				break
			}
			independent, err := manifestChoice(g, current, repo, r, data)
			if err != nil {
				ch.Reason = err.Error()
				break
			}
			if samePin {
				for id, download := range independent.Downloads {
					if expected, ok := ch.Packages[id]; ok && expected == independent.Packages[id] {
						ch.Downloads[id] = download
					}
				}
				break
			}
			independent.Current = ch.Current
			ch = independent
			break
		}
		ch = scopeUpdateChoice(current, ch)
		ch.Available = choiceChanged(g, current, ch)
		for _, p := range ch.Packages {
			if p.SHA256 == withheldMusicSHA {
				ch.Available = false
				ch.Reason = "This Music build is withheld for Microsoft analysis and cannot be installed."
			}
		}
		plan.Choices = append(plan.Choices, ch)
	}
	operationprogress.Report(ctx, "check", "Published release checks complete", int64(totalGroups), int64(totalGroups), "repositories")
	return plan, nil
}

func manifestChoice(g updateGroup, current catalog.Catalog, repo string, r Release, data []byte) (UpdateChoice, error) {
	var m AppUpdateManifest
	d := json.NewDecoder(strings.NewReader(string(data)))
	d.DisallowUnknownFields()
	if e := d.Decode(&m); e != nil {
		return UpdateChoice{}, e
	}
	var trailing any
	if d.Decode(&trailing) != io.EOF {
		return UpdateChoice{}, errors.New("trailing app update metadata")
	}
	if m.Repository != repo || m.Tag != r.Tag {
		return UpdateChoice{}, errors.New("unsupported app update contract")
	}
	if err := checkManifestRequirements(m, current.Version); err != nil {
		return UpdateChoice{}, err
	}
	allowedApps := map[string]bool{}
	for _, id := range g.apps {
		allowedApps[id] = true
	}
	allowedPackages := groupPackages(g, current)
	ch := UpdateChoice{ID: g.id, Name: g.name, Source: "https://github.com/" + repo + "/releases/tag/" + r.Tag, Apps: m.Apps, Packages: map[string]catalog.Package{}, Downloads: map[string]UpdateDownload{}}
	ch.MinimumManager = m.MinimumManager
	ch.RequiredCapabilities = append([]string(nil), m.RequiredCapabilities...)
	seen := map[string]bool{}
	for _, a := range m.Apps {
		if !allowedApps[a.ID] || seen[a.ID] || a.Source != ch.Source {
			return ch, errors.New("app update crosses repository boundaries")
		}
		seen[a.ID] = true
		ch.Version = a.Version
	}
	for _, id := range g.apps {
		if _, installed := current.Find(id); installed && !seen[id] {
			return ch, errors.New("app update is missing a coupled app")
		}
	}
	for id, p := range m.Packages {
		old, ok := allowedPackages[id]
		if !ok {
			// Older component updates did not register the optional integration
			// document. Match scopeUpdateChoice: leave it uninstalled rather than
			// blocking the executable update. Never ignore runtime packages or
			// another application's files.
			if untrackedCompatibilityDocument(g, current, id, p.File) {
				continue
			}
			omitted := false
			for _, a := range m.Apps {
				if a.ID == id && a.Package == id && a.Executable == p.File && a.Kind == "process" {
					if _, installed := current.Find(a.ID); !installed {
						omitted = true
					}
				}
			}
			if omitted && updatePackageOwner(id) == g.id {
				continue
			}
		}
		if !ok || old.File != p.File {
			return ch, errors.New("app update cannot replace another app's package or move its files")
		}
		a, ok := releaseAsset(r, p.Asset)
		if !ok || !validNamedAsset(a, repo, r.Tag, p.Asset, 2<<30) {
			return ch, errors.New("app update asset missing or invalid")
		}
		if p.Member == "" && (a.Digest != "sha256:"+p.SHA256 || a.Size != p.Size) {
			return ch, errors.New("direct app asset checksum differs from its contract")
		}
		if p.Member != "" {
			if _, e := platform.Inside("root", p.Member); e != nil {
				return ch, e
			}
		}
		ch.Packages[id] = p.Package
		ch.Downloads[id] = UpdateDownload{Asset: a, Repository: repo, Tag: r.Tag, Member: p.Member}
	}
	for id := range allowedPackages {
		if _, ok := ch.Packages[id]; !ok {
			return ch, fmt.Errorf("app update is missing package %s", id)
		}
	}
	target := current
	target.Packages = map[string]catalog.Package{}
	for id, p := range current.Packages {
		target.Packages[id] = p
	}
	for id, p := range ch.Packages {
		target.Packages[id] = p
	}
	target.Apps = append([]catalog.App(nil), current.Apps...)
	for i, a := range target.Apps {
		for _, n := range ch.Apps {
			if a.ID == n.ID {
				n.Port = a.Port
				target.Apps[i] = n
			}
		}
	}
	if e := target.Validate(); e != nil {
		return ch, e
	}
	return scopeUpdateChoice(current, ch), nil
}

func untrackedCompatibilityDocument(g updateGroup, current catalog.Catalog, id, file string) bool {
	if updatePackageOwner(id) != g.id {
		return false
	}
	for _, appID := range g.apps {
		a, installed := current.Find(appID)
		if installed && a.Kind == "process" && strings.HasSuffix(a.Executable, ".exe") &&
			id == a.ID+"-contract" && file == "docs/"+strings.TrimSuffix(a.Executable, ".exe")+"/compatibility.json" {
			return true
		}
		// Default-app installation may omit these optional Cinema guides. Ignore
		// only exact known document pairs when untracked. Tracked guides must
		// still be supplied, and can never move into another app's files.
		if installed && a.ID == "myai" && a.Kind == "process" && a.Executable == "MyAI.exe" {
			for _, name := range []string{"SETUP.md", "SUPPORT.md", "ADR-001-architecture.md", "SECURITY-CHECKS.md"} {
				if id == "myai-doc-"+name && file == "docs/MyAI/"+name {
					return true
				}
			}
		}
		if installed && a.ID == "cinema" && a.Kind == "process" && a.Executable == "MyCinema.exe" {
			for _, name := range []string{"SETUP.md", "SUPPORT.md", "ADR-001-architecture.md", "SECURITY-CHECKS.md"} {
				if id == "cinema-doc-"+name && file == "docs/MyCinema/"+name {
					return true
				}
			}
			if id == "cinema-provenance" && file == "docs/MyCinema/SOURCE-PROVENANCE.json" {
				return true
			}
		}
	}
	return false
}

// SelectUpdates retains unselected apps, custom extensions, locations and ports.
func SelectUpdates(current catalog.Catalog, plan UpdatePlan, selected []string) (catalog.Catalog, map[string]UpdateDownload, error) {
	if plan.BaseCatalogSHA256 != catalogIdentity(current) {
		return current, nil, errors.New("installation changed since the release check; check for updates again")
	}
	target := current
	target.Selection = current.EffectiveSelection()
	target.Packages = map[string]catalog.Package{}
	for id, p := range current.Packages {
		target.Packages[id] = p
	}
	target.Apps = append([]catalog.App(nil), current.Apps...)
	// A normal Workspace update includes every newly introduced default app.
	// Existing app updates remain independently selectable.
	selected = append([]string(nil), selected...)
	includeDefaults := false
	for _, id := range selected {
		if id == "workspace" {
			includeDefaults = true
		}
		for _, ch := range plan.Choices {
			if ch.ID == id && len(ch.NewApps) > 0 {
				includeDefaults = true
			}
		}
	}
	if includeDefaults {
		for _, ch := range plan.Choices {
			if len(ch.NewApps) > 0 && !ch.Available {
				return target, nil, fmt.Errorf("cannot include new apps from %s: %s", ch.Name, ch.Reason)
			}
			if !ch.Available || (ch.ID != "workspace" && len(ch.NewApps) == 0) {
				continue
			}
			present := false
			for _, id := range selected {
				if id == ch.ID {
					present = true
				}
			}
			if !present {
				selected = append(selected, ch.ID)
			}
		}
	}
	sources := map[string]UpdateDownload{}
	seen := map[string]bool{}
	for _, id := range selected {
		if seen[id] {
			return target, nil, errors.New("duplicate update selection")
		}
		seen[id] = true
		var ch *UpdateChoice
		for i := range plan.Choices {
			if plan.Choices[i].ID == id {
				ch = &plan.Choices[i]
				break
			}
		}
		if ch == nil || !ch.Available {
			return target, nil, fmt.Errorf("update %s is not available", id)
		}
		if ch.MinimumManager != "" {
			if err := checkRequirements(current.Version, ch.MinimumManager, ch.RequiredCapabilities); err != nil {
				return target, nil, err
			}
		}
		if id == "workspace" {
			target.Version = ch.Version
		}
		for pid, p := range ch.Packages {
			if p.SHA256 == withheldMusicSHA {
				return target, nil, errors.New("withheld Music executable")
			}
			target.Packages[pid] = p
			if current.Packages[pid] != p {
				sources[pid] = ch.Downloads[pid]
			}
		}
		for _, n := range ch.Apps {
			found := false
			for i, a := range target.Apps {
				if a.ID == n.ID {
					n.Port = a.Port
					target.Apps[i] = n
					found = true
					break
				}
			}
			if !found {
				allowed := false
				for _, added := range ch.NewApps {
					allowed = allowed || added == n.ID
				}
				if !allowed || target.Selection.Excludes(n.ID) {
					return target, nil, errors.New("new app is outside the selected default-app contract")
				}
				if !strings.HasPrefix(n.Source, "https://github.com/oneworkspacestudio/") {
					return target, nil, errors.New("new app must use an organization release")
				}
				if n.Kind == "process" {
					accounts, ok := current.Find("accounts")
					if !ok {
						return target, nil, errors.New("MyAccounts is required to install new apps")
					}
					n.Port += accounts.Port - 8991
				}
				target.Apps = append(target.Apps, n)
			}
		}
	}
	if len(seen) == 0 {
		return target, nil, errors.New("select at least one update")
	}
	return target, sources, target.Validate()
}

func OpenUpdateBundle(path string) (*zip.ReadCloser, catalog.Catalog, error) {
	z, e := zip.OpenReader(path)
	if e != nil {
		return nil, catalog.Catalog{}, e
	}
	c, e := bundleCatalog(&z.Reader)
	if e == nil {
		e = catalog.Verify(&z.Reader, c)
	}
	if e == nil {
		for _, a := range c.Apps {
			if !strings.HasPrefix(a.Source, "https://github.com/oneworkspacestudio/") {
				e = errors.New("update bundle contains a non-organization app")
				break
			}
		}
	}
	if e == nil {
		for _, p := range c.Packages {
			if p.SHA256 == withheldMusicSHA {
				e = errors.New("update bundle contains the withheld Music executable")
				break
			}
		}
	}
	if e != nil {
		z.Close()
		return nil, catalog.Catalog{}, e
	}
	return z, c, nil
}
func SelectBundleUpdates(current, incoming catalog.Catalog, selected []string) (catalog.Catalog, error) {
	if incoming.Version != current.Version && !newerWorkspace(incoming.Version, current.Version) {
		return current, errors.New("use rollback to restore an older Workspace version")
	}
	plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current)}
	for _, g := range suiteUpdateGroups(incoming) {
		ch := UpdateChoice{ID: g.id, Name: g.name, Version: incoming.Version, Packages: groupPackages(g, incoming), Apps: groupApps(g, incoming)}
		ch = markNewDefaults(current, incoming, ch)
		if g.id != "workspace" && len(groupApps(g, current)) == 0 && len(ch.NewApps) == 0 {
			continue
		}
		if g.id != "workspace" {
			if len(ch.Apps) == 0 {
				continue
			}
			ch.Version = ch.Apps[0].Version
		}
		ch = scopeUpdateChoice(current, ch)
		ch.Available = choiceChanged(g, current, ch)
		plan.Choices = append(plan.Choices, ch)
	}
	if len(selected) == 0 {
		for _, ch := range plan.Choices {
			if ch.Available {
				selected = append(selected, ch.ID)
			}
		}
	}
	c, _, e := SelectUpdates(current, plan, selected)
	return c, e
}

// StartupDefaultPlan bridges older controllers which upgraded the manager but
// discarded new app entries. Bind additions to that exact manager release, not
// whichever release happens to be newest while the transaction is restarting.
func StartupDefaultPlan(ctx context.Context, current catalog.Catalog) (UpdatePlan, error) {
	plan := UpdatePlan{CheckedAt: time.Now().UTC(), BaseCatalogSHA256: catalogIdentity(current)}
	h, token := client(), githubToken(ctx)
	ctx, indexErr := withReleaseIndex(ctx, token)
	if indexErr != nil {
		return plan, indexErr
	}
	plan.IndexPublishedAt = indexPublishedAt(ctx)
	rs, err := releasesFor(ctx, h, token, Repository)
	if err != nil {
		return plan, err
	}
	var release Release
	for _, r := range rs {
		if r.Tag == "v"+current.Version && !r.Draft && !r.Published.IsZero() {
			release = r
			break
		}
	}
	bundle, hasBundle := releaseAsset(release, bundleName)
	metadata, hasMetadata := releaseAsset(release, "catalog.json")
	if !hasBundle || !hasMetadata || !validAsset(bundle, Repository, release.Tag) {
		return plan, errors.New("matching Workspace release is unavailable; use the newest installer to finish adding apps")
	}
	raw, err := fetchUpdateBytes(ctx, h, token, Repository, release, metadata, 1<<20)
	if err != nil {
		return plan, err
	}
	latest, err := catalog.Read(raw)
	if err != nil {
		return plan, err
	}
	if latest.Version != current.Version || latest.Packages["manager"] != current.Packages["manager"] {
		return plan, errors.New("new-app catalog does not match the verified installed manager")
	}
	for _, g := range suiteUpdateGroups(latest) {
		ch := UpdateChoice{ID: g.id, Name: g.name, Apps: groupApps(g, latest), Packages: groupPackages(g, latest), Downloads: map[string]UpdateDownload{}}
		ch = markNewDefaults(current, latest, ch)
		if len(ch.NewApps) == 0 {
			continue
		}
		ch = scopeUpdateChoice(current, ch)
		ch.Available = true
		for id, p := range ch.Packages {
			ch.Downloads[id] = UpdateDownload{Asset: bundle, Repository: Repository, Tag: release.Tag, Member: p.File}
		}
		plan.Choices = append(plan.Choices, ch)
	}
	return plan, nil
}

// StageUpdateFiles verifies both download archives and their individual members.
func StageUpdateFiles(ctx context.Context, target catalog.Catalog, sources map[string]UpdateDownload, dir string, report func(string)) error {
	if e := os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	h, token := client(), githubToken(ctx)
	// Download each actual release asset once, in a deterministic order. Package
	// map order must never label a complete bundle as an internal notice file.
	ids := make([]string, 0, len(sources))
	for id := range sources {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	downloads := map[string]string{}
	assets := 0
	for _, id := range ids {
		d := sources[id]
		if !validNamedAsset(d.Asset, d.Repository, d.Tag, d.Asset.Name, 2<<30) {
			return errors.New("invalid staged asset")
		}
		if _, ok := downloads[d.Asset.APIURL]; !ok {
			downloads[d.Asset.APIURL] = ""
			assets++
		}
	}
	number := 0
	for _, id := range ids {
		d := sources[id]
		if downloads[d.Asset.APIURL] != "" {
			continue
		}
		number++
		label := fmt.Sprintf("Downloading %s %s â€” file %d of %d (%s)", strings.TrimPrefix(d.Repository, "oneworkspacestudio/"), strings.TrimPrefix(d.Tag, "v"), number, assets, d.Asset.Name)
		report(label)
		archivePath := filepath.Join(dir, "download-"+platform.Token())
		transferCtx := operationprogress.Observe(ctx, func(e operationprogress.Event) {
			operationprogress.Report(ctx, e.Stage, label, e.Done, e.Total, e.Unit)
		})
		if e := downloadMailAsset(transferCtx, h, token, d.Asset, archivePath); e != nil {
			return e
		}
		downloads[d.Asset.APIURL] = archivePath
	}
	report("Preparing and verifying the downloaded update packages")
	for index, id := range ids {
		d := sources[id]
		p := target.Packages[id]
		archivePath := downloads[d.Asset.APIURL]
		operationprogress.Report(ctx, "stage", "Verifying downloaded update packages", int64(index), int64(len(ids)), "packages")
		path, e := platform.Inside(filepath.Join(dir, "packages"), p.File)
		if e != nil {
			return e
		}
		if d.Member == "" {
			if e = platform.CopyFile(archivePath, path); e != nil {
				return e
			}
		} else {
			z, e := zip.OpenReader(archivePath)
			if e != nil {
				return e
			}
			var member *zip.File
			for _, f := range z.File {
				if f.Name == d.Member {
					if member != nil || f.Mode()&os.ModeSymlink != 0 {
						z.Close()
						return errors.New("ambiguous or linked archive member")
					}
					member = f
				}
			}
			if member == nil || member.UncompressedSize64 != uint64(p.Size) {
				z.Close()
				return errors.New("update archive member missing or wrong size")
			}
			in, e := member.Open()
			if e != nil {
				z.Close()
				return e
			}
			e = os.MkdirAll(filepath.Dir(path), 0700)
			if e == nil {
				var out *os.File
				out, e = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
				if e == nil {
					_, e = io.Copy(out, io.LimitReader(in, p.Size+1))
					if e == nil {
						e = out.Sync()
					}
					out.Close()
				}
			}
			in.Close()
			z.Close()
			if e != nil {
				return e
			}
		}
		held, e := platform.LockedBinary(path, p.SHA256)
		if e != nil {
			return e
		}
		st, e := held.Stat()
		held.Close()
		if e != nil || st.Size() != p.Size {
			return errors.New("staged package size mismatch")
		}
		operationprogress.Report(ctx, "stage", "Verifying downloaded update packages", int64(index+1), int64(len(ids)), "packages")
	}
	return nil
}
