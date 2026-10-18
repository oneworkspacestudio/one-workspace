package distribution

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// MakeAppUpdateManifest lets each app's release pipeline publish a contract
// without rebuilding Workspace. The existing managed launch/storage contract is
// retained; every owned executable/helper/license must be supplied together.
func MakeAppUpdateManifest(base catalog.Catalog, group, version, repo, tag string, files map[string]UpdatePackage) (AppUpdateManifest, error) {
	var result AppUpdateManifest
	var found *updateGroup
	for i := range updateGroups {
		if updateGroups[i].id == group && group != "workspace" {
			found = &updateGroups[i]
			break
		}
	}
	if found == nil || repo != "oneworkspacestudio/"+found.repo || tag == "" {
		return result, errors.New("choose an organization app update group and matching repository")
	}
	if e := base.Validate(); e != nil {
		return result, e
	}
	owned := groupPackages(*found, base)
	if len(owned) != len(files) {
		return result, errors.New("supply all executable, helper and license files in the app update group")
	}
	target := base
	target.Packages = map[string]catalog.Package{}
	for id, p := range base.Packages {
		target.Packages[id] = p
	}
	for id, p := range files {
		old, ok := owned[id]
		if !ok || old.File != p.File {
			return result, fmt.Errorf("package %s is not owned by this app group", id)
		}
		if !regexp.MustCompile(`^[A-Za-z0-9_.+-]{1,180}$`).MatchString(p.Asset) {
			return result, errors.New("use a plain release asset filename")
		}
		if p.Member != "" {
			if _, e := platform.Inside("root", p.Member); e != nil {
				return result, e
			}
		}
		if p.SHA256 == withheldMusicSHA {
			return result, errors.New("the withheld Music build cannot be published as an update")
		}
		target.Packages[id] = p.Package
	}
	apps := groupApps(*found, base)
	if len(apps) == 0 {
		return result, errors.New("app is missing from the base catalog")
	}
	source := "https://github.com/" + repo + "/releases/tag/" + tag
	for i := range apps {
		apps[i].Version = version
		apps[i].Source = source
	}
	target.Apps = append([]catalog.App(nil), base.Apps...)
	for i, a := range target.Apps {
		for _, n := range apps {
			if a.ID == n.ID {
				target.Apps[i] = n
			}
		}
	}
	if e := target.Validate(); e != nil {
		return result, e
	}
	return AppUpdateManifest{Schema: 1, Integration: 2, Repository: repo, Tag: tag, Apps: apps, Packages: files}, nil
}
