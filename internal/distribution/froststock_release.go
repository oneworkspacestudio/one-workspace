package distribution

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

const froststockRepository = "oneworkspacestudio/froststock"

var froststockLaunch = catalog.App{
	ID: "froststock", Name: "Froststock", Kind: "process", Package: "froststock",
	Executable: "Froststock.exe", Port: 9048, Route: "/froststock", StripPrefix: true,
	Health: "/healthz", Callback: "/froststock/auth/callback", OIDC: true,
	Depends: []string{"accounts", "drive"},
	Args:    []string{"--data", "${DATA}/froststock", "--listen", "127.0.0.1:${PORT}", "--origin", "${ORIGIN}/froststock", "--no-browser"},
}

var froststockPackagePaths = map[string]string{
	"froststock":          "Froststock.exe",
	"froststock-contract": "docs/Froststock/compatibility.json",
	"froststock-guide":    "docs/Froststock/guide.md",
	"froststock-license":  "licenses/Froststock-LICENSE.txt",
	"froststock-notices":  "licenses/Froststock-NOTICE.txt",
}

func froststockReleaseChoice(ctx context.Context, h *http.Client, token string, g updateGroup, current catalog.Catalog) UpdateChoice {
	previous, installed := current.Find(g.id)
	choice := UpdateChoice{ID: g.id, Name: g.name, Current: previous.Version, Version: previous.Version, Source: previous.Source, Optional: !installed}
	if !installed {
		choice.Current = "Not installed"
		if _, ok := current.Find("drive"); !ok {
			choice.Reason = "Install My Drive to use Froststock."
			return choice
		}
	}
	releases, err := releasesFor(ctx, h, token, froststockRepository)
	if err != nil {
		choice.Reason = err.Error()
		return choice
	}
	var latest Release
	for _, r := range releases {
		if r.Draft || r.Published.IsZero() || !strings.HasPrefix(r.Tag, "v") {
			continue
		}
		version := strings.TrimPrefix(r.Tag, "v")
		if !managerRequirementVersion.MatchString(version) || (installed && !newerWorkspace(version, previous.Version)) {
			continue
		}
		if latest.Tag == "" || newerWorkspace(version, strings.TrimPrefix(latest.Tag, "v")) {
			latest = r
		}
	}
	if latest.Tag == "" {
		return choice
	}
	choice.Version = strings.TrimPrefix(latest.Tag, "v")
	choice.ReviewURL = "https://github.com/" + froststockRepository + "/releases/tag/" + latest.Tag
	asset, ok := releaseAsset(latest, "workspace-update.json")
	if !ok {
		choice.Reason = "This Froststock release has no Workspace install manifest."
		return choice
	}
	raw, err := fetchUpdateBytes(ctx, h, token, froststockRepository, latest, asset, 1<<20)
	if err == nil {
		var next UpdateChoice
		next, err = froststockManifestChoice(g, current, latest, raw)
		if err == nil {
			next.Current = choice.Current
			next.Available = true
			return next
		}
	}
	choice.Reason = err.Error()
	return choice
}

func froststockManifestChoice(g updateGroup, current catalog.Catalog, r Release, data []byte) (UpdateChoice, error) {
	previous, installed := current.Find(g.id)
	// Only this registered external app can introduce an application from its
	// own repository. Seed the existing manifest validator with the exact
	// package namespace and launch contract the Workspace manager supports.
	baseline := current
	baseline.Apps = slices.Clone(current.Apps)
	baseline.Packages = make(map[string]catalog.Package, len(current.Packages)+len(froststockPackagePaths))
	for id, pkg := range current.Packages {
		baseline.Packages[id] = pkg
	}
	if !installed {
		baseline.Apps = append(baseline.Apps, froststockLaunch)
		baseline.Selection = current.EffectiveSelection()
		baseline.Selection.Excluded = slices.DeleteFunc(baseline.Selection.Excluded, func(id string) bool { return id == g.id })
	}
	for id, path := range froststockPackagePaths {
		if p, ok := baseline.Packages[id]; ok && p.File != path {
			return UpdateChoice{}, errors.New("Froststock has a different installed package path")
		} else if !ok {
			baseline.Packages[id] = catalog.Package{File: path}
		}
	}
	next, err := manifestChoice(g, baseline, froststockRepository, r, data)
	if err != nil {
		return next, err
	}
	if len(next.Apps) != 1 || len(next.Packages) != len(froststockPackagePaths) ||
		r.Tag != "v"+next.Apps[0].Version || (installed && !newerWorkspace(next.Apps[0].Version, previous.Version)) {
		return next, errors.New("invalid or non-advancing Froststock release")
	}
	if next.MinimumManager == "" {
		return next, errors.New("Froststock release is missing its minimum Workspace version")
	}
	launch := next.Apps[0]
	launch.Version, launch.Source, launch.Description = "", "", ""
	if !reflect.DeepEqual(launch, froststockLaunch) {
		return next, errors.New("Froststock release changes its managed launch contract")
	}
	if !installed {
		next.NewApps, next.Optional = []string{g.id}, true
	}
	return next, nil
}
