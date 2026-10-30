package suite

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// AppUpdateScope is an independent recovery copy, never a delta depending on a
// previous backup. Paths come from the manager's storage registry, not a release.
// Schema 2 prevents older controllers treating this as a full workspace snapshot.
type AppUpdateScope struct {
	Apps     []string        `json:"apps"`
	Packages []string        `json:"packages"`
	After    catalog.Catalog `json:"after"`
}

type updateBackupPath struct{ area, path string }

// Shared engines, identity and storage layout changes need Workspace-wide
// metadata recovery. New layouts must be reviewed here before using an app scope.
func appUpdatePrivatePaths(id string) []string {
	switch id {
	case "my-tube", "notes":
		return []string{id, id + "-settings.json"}
	case "chat", "slicer", "tasks", "pdf", "monitor", "search", "maps", "meet", "create", "my3d", "cinema":
		return []string{id}
	case "authenticator", "shield", "finance", "remote", "tickets":
		return []string{id}
	case "forms":
		return []string{"forms", "forms-drive"}
	case "hosting":
		return []string{"my-hosting"}
	case "vault":
		return []string{"my-vault"}
	case "myai":
		return []string{"myai", "myai-drive", "myai-projects"}
	case "clips":
		return []string{"clips", "videos-clips"}
	case "social":
		return []string{"social", "photos-social"}
	}
	return nil
}

func appOwnsUpdatePackage(a catalog.App, id string) bool {
	return a.Package != "" && (id == a.Package || strings.HasPrefix(id, a.Package+"-"))
}

func sameUpdateAppContract(a, b catalog.App) bool {
	a.Version, a.Source, b.Version, b.Source = "", "", "", ""
	return reflect.DeepEqual(a, b)
}

// Only existing, independently owned apps with unchanged integration contracts
// qualify. Other supported changes use Workspace-wide metadata recovery.
func planAppUpdateScope(c Config, target catalog.Catalog) *AppUpdateScope {
	if c.Catalog.Selection == nil || c.Version != target.Version || c.Catalog.Version != target.Version || len(c.Catalog.Apps) != len(target.Apps) || len(c.Catalog.Packages) != len(target.Packages) {
		return nil
	}
	s := &AppUpdateScope{After: target}
	affected := map[string]bool{}
	for _, a := range c.Catalog.Apps {
		b, ok := target.Find(a.ID)
		if !ok || !sameUpdateAppContract(a, b) {
			return nil
		}
		if !reflect.DeepEqual(a, b) {
			affected[a.ID] = true
		}
	}
	for id, p := range c.Catalog.Packages {
		next, ok := target.Packages[id]
		if !ok || p.File != next.File {
			return nil
		}
		if p == next {
			continue
		}
		s.Packages = append(s.Packages, id)
		owned := false
		for _, a := range c.Catalog.Apps {
			if appOwnsUpdatePackage(a, id) {
				affected[a.ID], owned = true, true
			}
		}
		if !owned {
			return nil
		}
	}
	if len(s.Packages) == 0 {
		return nil
	}
	for id := range affected {
		s.Apps = append(s.Apps, id)
	}
	sort.Strings(s.Apps)
	sort.Strings(s.Packages)
	if validateAppUpdateScope(c, s) != nil {
		return nil
	}
	for _, id := range s.Apps {
		if id != "my-tube" && id != "notes" {
			continue
		}
		var settings struct {
			DataDir  string `json:"data_dir"`
			MediaDir string `json:"media_dir"`
		}
		if err := platform.LoadJSON(filepath.Join(c.Data, id+"-settings.json"), &settings); err != nil && !os.IsNotExist(err) {
			return nil
		}
		for _, path := range []string{settings.DataDir, settings.MediaDir} {
			if path == "" {
				continue
			}
			rel, err := filepath.Rel(filepath.Join(c.Data, id), filepath.Clean(path))
			if !filepath.IsAbs(path) || err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
				return nil
			}
		}
	}
	return s
}

func updateScopePaths(c Config, s *AppUpdateScope) []updateBackupPath {
	paths := []updateBackupPath{{"private", "workspace.json"}, {"private", "secrets.dpapi"}, {"private", "Workspace.exe"}}
	for _, id := range s.Apps {
		for _, path := range appUpdatePrivatePaths(id) {
			paths = append(paths, updateBackupPath{"private", path})
		}
	}
	for _, id := range s.Packages {
		paths = append(paths, updateBackupPath{"programs", c.Catalog.Packages[id].File})
	}
	return paths
}

func scopeContainsPath(c Config, s *AppUpdateScope, area, path string) bool {
	path = strings.ToLower(filepath.ToSlash(path))
	for _, p := range updateScopePaths(c, s) {
		root := strings.ToLower(filepath.ToSlash(p.path))
		if area == p.area && (path == root || (scopeTreePath(p) && strings.HasPrefix(path, root+"/"))) {
			return true
		}
	}
	return false
}

func scopeTreePath(p updateBackupPath) bool {
	return p.area == "private" && p.path != "workspace.json" && p.path != "secrets.dpapi" && p.path != "Workspace.exe" && !strings.HasSuffix(p.path, "-settings.json")
}

func scopeDirectoryAllowed(c Config, s *AppUpdateScope, area, path string) bool {
	for _, p := range updateScopePaths(c, s) {
		if scopeTreePath(p) && area == p.area && (strings.EqualFold(path, p.path) || strings.HasPrefix(strings.ToLower(path), strings.ToLower(p.path)+"/")) {
			return true
		}
	}
	return false
}

var updateDataReference = regexp.MustCompile(`\$\{DATA\}[/\\]([^\s"';]+)`)
var updateAbsoluteDrive = regexp.MustCompile(`[A-Za-z]:[/\\][^/\\]`)

func validateAppUpdateScope(c Config, s *AppUpdateScope) error {
	if s == nil || c.Catalog.Selection == nil || len(s.Apps) == 0 || len(s.Packages) == 0 || c.Catalog.Version != c.Version || c.Version != s.After.Version || c.Catalog.Schema != s.After.Schema || c.Catalog.InstallerIntegration != s.After.InstallerIntegration || !reflect.DeepEqual(c.Catalog.Selection, s.After.Selection) || len(c.Catalog.Apps) != len(s.After.Apps) || len(c.Catalog.Packages) != len(s.After.Packages) {
		return errors.New("invalid app recovery scope")
	}
	if err := c.Catalog.Validate(); err != nil {
		return err
	}
	if err := s.After.Validate(); err != nil {
		return err
	}
	apps, packages := map[string]bool{}, map[string]bool{}
	for _, id := range s.Apps {
		if apps[id] || len(appUpdatePrivatePaths(id)) == 0 {
			return errors.New("app recovery scope requires a known independent storage layout")
		}
		apps[id] = true
		if _, ok := c.Catalog.Find(id); !ok {
			return errors.New("app recovery scope cannot add or remove apps")
		}
	}
	for _, id := range s.Packages {
		p, ok := c.Catalog.Packages[id]
		next, nextOK := s.After.Packages[id]
		if packages[id] || !ok || !nextOK || p.File != next.File || id == "manager" {
			return errors.New("invalid app recovery package")
		}
		packages[id] = true
		owned := false
		for _, app := range c.Catalog.Apps {
			if appOwnsUpdatePackage(app, id) {
				if !apps[app.ID] {
					return errors.New("app recovery package has an unselected owner")
				}
				owned = true
			}
		}
		if !owned {
			return errors.New("app recovery package ownership is unknown")
		}
	}
	for id, p := range c.Catalog.Packages {
		q, ok := s.After.Packages[id]
		if !ok || (!packages[id] && p != q) {
			return errors.New("app recovery scope changes an unrelated package")
		}
	}
	for _, a := range c.Catalog.Apps {
		b, ok := s.After.Find(a.ID)
		if !ok || !sameUpdateAppContract(a, b) || (!apps[a.ID] && !reflect.DeepEqual(a, b)) {
			return errors.New("app recovery scope changes an unrelated app or its integration")
		}
		if !apps[a.ID] {
			continue
		}
		if a.Kind != "process" || a.Package != a.ID {
			return errors.New("app recovery requires an independent app package")
		}
		values := append([]string{}, a.Args...)
		for _, value := range a.Env {
			values = append(values, value)
		}
		for _, value := range values {
			_, assigned, assignment := strings.Cut(value, "=")
			if filepath.IsAbs(value) || updateAbsoluteDrive.MatchString(value) || (assignment && filepath.IsAbs(assigned)) {
				return errors.New("custom absolute storage settings require a full recovery copy")
			}
			refs := updateDataReference.FindAllStringSubmatch(value, -1)
			if strings.Contains(value, "${FILES}") || (strings.Contains(value, "${DATA}") && len(refs) == 0) {
				return errors.New("custom or shared storage requires a full recovery copy")
			}
			for _, ref := range refs {
				if _, err := platform.Inside(c.Data, filepath.ToSlash(ref[1])); err != nil || !scopeContainsPath(c, s, "private", ref[1]) {
					return errors.New("storage binding is outside the app recovery scope")
				}
			}
		}
	}
	return nil
}

// Enumerate only affected paths, so unrelated large libraries are not even
// traversed. Missing paths are still in the scope and removed on rollback if
// the failed update created them. Existing ancestors must never be junctions.
func walkScopedBackup(c Config, s *AppUpdateScope, policy *UpdateMetadataPolicy, visit func(area, root, path string, entry fs.DirEntry) error) error {
	if s == nil {
		return walkBackupMetadata(c, policy, visit)
	}
	if err := validateAppUpdateScope(c, s); err != nil {
		return err
	}
	for _, p := range updateScopePaths(c, s) {
		root := c.Data
		if p.area == "programs" {
			root = c.Programs
		}
		start, err := checkedScopePath(root, p.path)
		if err != nil {
			return err
		}
		if _, err := os.Lstat(start); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		if err = filepath.WalkDir(start, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if !metadataPathAllowed(policy, p.area, filepath.ToSlash(rel), d.IsDir()) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if err = platform.RejectReparse(path); err != nil {
				return err
			}
			if !d.IsDir() && !d.Type().IsRegular() {
				return errors.New("unsupported special file in app recovery copy")
			}
			return visit(p.area, root, path, d)
		}); err != nil {
			return err
		}
	}
	return nil
}

func checkedScopePath(root, relative string) (string, error) {
	path, err := platform.Inside(root, relative)
	if err != nil {
		return "", err
	}
	for ancestor := path; ; ancestor = filepath.Dir(ancestor) {
		if _, err := os.Lstat(ancestor); err == nil {
			if err = platform.RejectReparse(ancestor); err != nil {
				return "", err
			}
		} else if !os.IsNotExist(err) {
			return "", err
		}
		if ancestor == root {
			break
		}
	}
	return path, nil
}

func sameRecoveryScope(a, b *AppUpdateScope) bool {
	return a != nil && b != nil && reflect.DeepEqual(a.Apps, b.Apps) && reflect.DeepEqual(a.Packages, b.Packages)
}

func checkScopedRollback(c Config, b UpdateBackup) error {
	if b.Metadata != nil && !reflect.DeepEqual(c.Catalog, b.Metadata.After) {
		return errors.New("this update recovery no longer matches the installed versions; choose a full data backup instead")
	}
	if b.Scope != nil && !reflect.DeepEqual(c.Catalog, b.Scope.After) {
		return errors.New("this app recovery copy no longer matches the installed versions; use a full Workspace backup instead")
	}
	return nil
}

func validateBackupScope(b UpdateBackup, r BackupReceipt) error {
	if b.Schema == 3 {
		if b.Metadata == nil || !r.MetadataOnly {
			return errors.New("update metadata policy and receipt differ")
		}
		if err := validateUpdateMetadata(b.Config, b.Metadata); err != nil {
			return err
		}
		if b.Scope == nil {
			if len(r.Apps) != 0 {
				return errors.New("Workspace update recovery cannot declare an app-only receipt")
			}
			return nil
		}
		if !reflect.DeepEqual(b.Scope.Apps, r.Apps) || !reflect.DeepEqual(b.Scope.After, b.Metadata.After) {
			return errors.New("app recovery scope and receipt differ")
		}
		return validateAppUpdateScope(b.Config, b.Scope)
	}
	if b.Metadata != nil || r.MetadataOnly {
		return errors.New("legacy full-data snapshots cannot use an update metadata policy")
	}
	if b.Schema == 1 {
		if b.Scope != nil || len(r.Apps) != 0 {
			return errors.New("a full Workspace backup cannot declare an app scope")
		}
		return nil
	}
	if b.Schema != 2 || b.Scope == nil || !reflect.DeepEqual(b.Scope.Apps, r.Apps) {
		return errors.New("app recovery manifest and receipt differ")
	}
	return validateAppUpdateScope(b.Config, b.Scope)
}

// Called only after both the original snapshot and the matching current-state
// archive have been verified. Never clear an unselected app, the general Files
// tree, or manager program file. The recovery guard survives until the last copy.
func clearAppUpdateScope(c Config, s *AppUpdateScope, policy *UpdateMetadataPolicy) error {
	for _, p := range updateScopePaths(c, s) {
		if p.area == "private" && strings.EqualFold(p.path, "Workspace.exe") {
			continue
		}
		root := c.Data
		if p.area == "programs" {
			root = c.Programs
		}
		path, err := checkedScopePath(root, p.path)
		if err != nil {
			return err
		}
		if _, err = os.Lstat(path); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		if policy != nil && p.area == "private" {
			if err = clearRecoveryMetadata(c.Data, path); err != nil {
				return err
			}
			continue
		}
		if err = cleanupPath(path); err != nil {
			return err
		}
		if err = filepath.WalkDir(path, func(current string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if err = platform.RejectReparse(current); err != nil {
				return err
			}
			if !d.IsDir() && !d.Type().IsRegular() {
				return errors.New("unsupported file in app recovery scope")
			}
			return platform.ClearFileReadOnly(current)
		}); err != nil {
			return err
		}
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}
