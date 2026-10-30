package suite

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// UpdateMetadataPolicy records exactly which program paths may be replaced.
// Data libraries are preserved in place, not copied into update recovery. Full
// data protection belongs to the separate manual/scheduled Backups operation.
type UpdateMetadataPolicy struct {
	Version  int             `json:"version"`
	Programs []string        `json:"programs"`
	Target   catalog.Catalog `json:"target"`
	After    catalog.Catalog `json:"after"`
}

// Version 1 is a storage contract, not a filename/size heuristic. Never exclude
// a whole DB store (for example Vault's engine or Mail's store). DB/WAL/SHM,
// credentials, keys, schema files and settings remain in the recovery copy.
// Keep this policy stable for old recovery records; new layouts need a new version.
var preservedUpdateDataV1 = []string{
	"accounts/uploads", "accounts/logs", "drive/cache", "drive/.filebrowser-history", "photos/cache", "photos/logs",
	"gitea/repositories", "gitea/data/attachments", "gitea/data/lfs", "gitea/data/packages", "gitea/data/repo-archive", "gitea/data/avatars", "gitea/data/repo-avatars", "gitea/log",
	"music/cache", "music/logs", "mail/runtime", "mail/blobs", "mail/logs", "shared/bin", "shared/blobs", "shared/logs",
	"my-tube/media", "my-tube/helpers", "my-tube/uploads", "my-tube/backups", "my-tube/logs",
	"notes/attachments", "notes/uploads", "notes/backups", "notes/logs",
	"chat/runtime", "chat/attachments", "chat/logs", "chat/backups",
	"my-vault/runtime", "my-vault/engine/attachments", "my-vault/engine/sends", "my-vault/engine/icon_cache", "my-vault/backups",
	"my-hosting/projects", "my-hosting/runtimes", "my-hosting/cache", "my-hosting/backups", "my-hosting/logs",
	"slicer/engine", "slicer/models", "slicer/slices", "slicer/projects", "slicer/logs",
	"cinema/runtime", "cinema/jellyfin/cache", "cinema/jellyfin/logs", "cinema/transcodes",
	"myai-drive", "myai-projects/projects", "myai-projects/git", "myai/models", "myai/search-runtime", "myai/agent-runtime", "myai/cache",
	"videos-clips", "photos-social", "clips/uploads", "clips/cache", "clips/outputs", "social/uploads", "social/cache", "social/outputs",
	"maps/tiles", "maps/offline", "maps/cache", "maps/downloads", "meet/recordings", "meet/engine",
	"create/uploads", "create/outputs", "create/assets", "create/cache", "my3d/models", "my3d/projects", "my3d/outputs", "my3d/cache",
	"pdf/uploads", "pdf/outputs", "pdf/cache", "search/cache", "monitor/logs",
	"tickets/core/attachments", "tickets/tmp",
}

func preservedUpdateData(path string) bool {
	path = strings.ToLower(filepath.ToSlash(path))
	for _, root := range preservedUpdateDataV1 {
		if path == root || strings.HasPrefix(path, root+"/") {
			return true
		}
	}
	return false
}

func containsPreservedUpdateData(path string) bool {
	path = strings.ToLower(filepath.ToSlash(path))
	if path == "." || path == "" {
		return true
	}
	for _, root := range preservedUpdateDataV1 {
		if strings.HasPrefix(root, path+"/") {
			return true
		}
	}
	return false
}

func planUpdateMetadata(c Config, target catalog.Catalog, repair bool) (*UpdateMetadataPolicy, error) {
	if !distribution.MetadataRecoverySupported(c.Version) && (!distribution.MetadataRecoverySupported(target.Version) || c.Catalog.Packages["manager"] == target.Packages["manager"]) {
		return nil, errors.New("update One Workspace to preview91 or newer first, or select it together with the app updates, to enable program and metadata recovery")
	}
	if err := validateManagedStorageSettings(c, true); err != nil {
		return nil, err
	}
	files := map[string]bool{}
	for id, old := range c.Catalog.Packages {
		next, ok := target.Packages[id]
		if repair || !ok || old != next {
			files[old.File] = true
			if ok {
				files[next.File] = true
			}
		}
	}
	for id, p := range target.Packages {
		if _, exists := c.Catalog.Packages[id]; !exists {
			files[p.File] = true
		}
	}
	p := &UpdateMetadataPolicy{Version: 1, Target: target, After: target}
	for name := range files {
		p.Programs = append(p.Programs, name)
	}
	sort.Strings(p.Programs)
	owned := map[string]bool{}
	for _, old := range c.Catalog.Packages {
		owned[strings.ToLower(old.File)] = true
	}
	for _, path := range p.Programs {
		if owned[strings.ToLower(path)] {
			continue
		}
		full, err := checkedScopePath(c.Programs, path)
		if err != nil {
			return nil, err
		}
		if _, err = os.Lstat(full); err == nil {
			return nil, errors.New("a new program path already contains an unmanaged file; move it outside the Apps folder before updating")
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return p, validateUpdateMetadata(c, p)
}

func validateUpdateMetadata(c Config, p *UpdateMetadataPolicy) error {
	if p == nil || p.Version != 1 {
		return errors.New("unsupported update recovery storage policy")
	}
	if err := validateProgramDataSeparation(c); err != nil {
		return err
	}
	allowed := map[string]bool{}
	for _, cat := range []catalog.Catalog{c.Catalog, p.Target, p.After} {
		if err := cat.Validate(); err != nil {
			return err
		}
		candidate := c
		candidate.Catalog = cat
		if err := validateProgramDataSeparation(candidate); err != nil {
			return err
		}
		for _, app := range cat.Apps {
			if !knownUpdateStorageProfile(app.ID) {
				return errors.New("this app needs a managed data-storage profile before automatic update recovery can exclude its user files")
			}
		}
		for _, pkg := range cat.Packages {
			allowed[strings.ToLower(pkg.File)] = true
		}
	}
	seen := map[string]bool{}
	for _, path := range p.Programs {
		key := strings.ToLower(path)
		if _, err := platform.Inside(c.Programs, path); err != nil || !allowed[key] || seen[key] {
			return errors.New("update recovery contains an unowned program path")
		}
		seen[key] = true
	}
	return nil
}

func knownUpdateStorageProfile(id string) bool {
	if id == "authenticator" || id == "finance" || id == "remote" || id == "tickets" {
		return true
	}
	switch id {
	case "accounts", "drive", "docs", "photos", "studio", "gitea", "music", "mail", "calendar", "contacts", "notes", "my-tube", "chat", "hosting", "vault", "slicer", "tasks", "pdf", "monitor", "search", "maps", "meet", "create", "my3d", "cinema", "myai", "clips", "social", "push", "shield", "forms":
		return true
	}
	return false
}

func metadataPathAllowed(p *UpdateMetadataPolicy, area, path string, directory bool) bool {
	if p == nil {
		return true
	}
	if area == "files" {
		return false
	}
	if area == "private" {
		return !preservedUpdateData(path)
	}
	path = strings.ToLower(filepath.ToSlash(path))
	for _, file := range p.Programs {
		file = strings.ToLower(file)
		if file == path || (directory && strings.HasPrefix(file, path+"/")) {
			return true
		}
	}
	return false
}

func sameMetadataPolicy(a, b *UpdateMetadataPolicy) bool {
	return a != nil && b != nil && a.Version == b.Version && reflect.DeepEqual(a.Programs, b.Programs) && reflect.DeepEqual(a.Target, b.Target)
}

// Gather and check first, then remove metadata entries individually. Never
// RemoveAll an ancestor of a preserved library: files created after an update
// must survive its rollback too. No preserved tree is traversed or modified.
func clearRecoveryMetadata(root, start string) error {
	var paths []string
	err := filepath.WalkDir(start, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		first := strings.ToLower(strings.SplitN(rel, "/", 2)[0])
		if privateUpdateExclusions[first] || strings.HasPrefix(first, ".apps-stage-") || preservedUpdateData(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.EqualFold(rel, "Workspace.exe") {
			return nil // Keep the recovery guard until the last verified copy.
		}
		if err = platform.RejectReparse(path); err != nil {
			return err
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			return errors.New("unsupported special file in update metadata")
		}
		if path != root && (!d.IsDir() || !containsPreservedUpdateData(rel)) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for i := len(paths) - 1; i >= 0; i-- {
		path := paths[i]
		if err = cleanupPath(path); err != nil {
			return err
		}
		if err = platform.ClearFileReadOnly(path); err != nil {
			return err
		}
		if err = os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
