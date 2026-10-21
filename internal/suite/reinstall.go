package suite

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type SetupStorageReview struct {
	Kind    string              `json:"kind"` // new, repair, restart
	Review  string              `json:"review"`
	Message string              `json:"message"`
	Volume  platform.DataVolume `json:"volume"`
}

func setupLockPath(r Request) string {
	h := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(r.Data))))
	return filepath.Join(filepath.Dir(r.Data), ".workspace-setup-"+hex.EncodeToString(h[:8])+".lock")
}

// Read-only inspection. No permissions, running service, or existing files are
// changed merely by selecting a folder in the wizard.
func InspectSetupStorage(r Request) (SetupStorageReview, error) {
	var result SetupStorageReview
	if err := r.ResolveStorage(); err != nil {
		return result, err
	}
	for _, root := range []string{r.Data, r.Programs, r.Files} {
		if !filepath.IsAbs(root) || filepath.Dir(filepath.Clean(root)) == filepath.Clean(root) {
			return result, errors.New("Choose a dedicated Workspace folder")
		}
		for p := filepath.Clean(root); ; p = filepath.Dir(p) {
			if _, err := os.Lstat(p); err == nil {
				if err = platform.RejectReparse(p); err != nil {
					return result, err
				}
			} else if !os.IsNotExist(err) {
				return result, err
			}
			if filepath.Dir(p) == p {
				break
			}
		}
	}
	// Validate the separation independently of account/domain form fields.
	paths := []string{r.Data, r.Programs, r.Files}
	for i, a := range paths {
		for j, b := range paths {
			if i != j && (strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) || strings.HasPrefix(strings.ToLower(filepath.Clean(a)), strings.ToLower(filepath.Clean(b))+string(os.PathSeparator))) {
				return result, errors.New("Workspace apps, private settings and files must use separate folders")
			}
		}
	}
	volume, err := platform.CheckDataVolume(r.Data)
	result.Volume = volume
	if err != nil {
		return result, err
	}
	h := sha256.New()
	fmt.Fprintln(h, filepath.Clean(r.Data), filepath.Clean(r.Programs), filepath.Clean(r.Files))
	occupied := false
	for _, root := range []string{r.Data, r.Programs, r.Files} {
		st, e := os.Lstat(root)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return result, e
		}
		if !st.IsDir() {
			return result, errors.New("Choose folders for Workspace apps, settings and files")
		}
		if root == r.Files {
			continue
		}
		entries, e := os.ReadDir(root)
		if e != nil {
			return result, e
		}
		if len(entries) > 0 {
			occupied = true
		}
	}
	result.Kind = "new"
	if occupied {
		result.Kind = "restart"
		if c, secrets, e := Load(r.Data); e == nil && c.Owner != "" && secrets.AdminPassword != "" {
			if !strings.EqualFold(c.Data, filepath.Clean(r.Data)) || !strings.EqualFold(c.Programs, filepath.Clean(r.Programs)) || !strings.EqualFold(c.Files, filepath.Clean(r.Files)) {
				return result, errors.New("This Workspace uses different storage folders. Open its original installer location or Updates and backups.")
			}
			result.Kind = "repair"
			for _, name := range []string{"workspace.json", "secrets.dpapi"} {
				sum, _, e := fileDigest(filepath.Join(r.Data, name))
				if e != nil {
					return result, e
				}
				fmt.Fprintln(h, name, sum)
			}
		} else {
			// An unfinished setup can contain user data too. Bind the decision to
			// the reviewed tree and refuse junctions before moving anything.
			for _, root := range []string{r.Data, r.Programs} {
				if _, e := os.Stat(root); os.IsNotExist(e) {
					continue
				}
				err = filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
					if e != nil {
						return e
					}
					if e = platform.RejectReparse(p); e != nil {
						return e
					}
					if d.Name() == "runtime.lock" || d.Name() == "install.lock" {
						return nil
					}
					st, e := d.Info()
					if e != nil {
						return e
					}
					if !d.IsDir() && !st.Mode().IsRegular() {
						return errors.New("This folder contains a special file that needs a Windows backup")
					}
					stamp := st.ModTime().UnixNano()
					if d.IsDir() {
						stamp = 0
					}
					fmt.Fprintln(h, p, st.Size(), stamp, st.Mode())
					return nil
				})
				if err != nil {
					return result, err
				}
			}
		}
	}
	result.Review = hex.EncodeToString(h.Sum(nil))
	switch result.Kind {
	case "repair":
		result.Message = "Workspace is already here. Reinstall its program files while keeping your accounts, app data, files and connection settings. A recovery copy is saved first."
	case "restart":
		result.Message = "A previous setup is incomplete. Reinstall here with a new sign-in setup. Your Files folder stays in place; previous accounts, app data and settings are saved in a protected recovery folder."
	default:
		result.Message = "Ready to install Workspace in this folder."
	}
	return result, nil
}

func reinstall(ctx context.Context, r Request, z *zip.Reader, cat catalog.Catalog, report func(string)) (Config, error) {
	review, err := InspectSetupStorage(r)
	if err != nil {
		return Config{}, err
	}
	if r.ReinstallReview != "" && r.ReinstallReview != review.Review {
		return Config{}, errors.New("This folder changed. Review the reinstall option again before continuing.")
	}
	if review.Kind == "new" {
		return Config{}, errors.New("This folder no longer needs reinstalling. Check it again and choose Install Workspace.")
	}
	if report == nil {
		report = func(string) {}
	}
	if review.Kind == "restart" {
		if err = r.Validate(); err != nil {
			return Config{}, err
		}
		// Validate the complete replacement payload before retaining the old
		// setup. The fresh installer independently checks its selected bytes.
		selected, e := SelectInstallApps(cat, r)
		if e != nil {
			return Config{}, e
		}
		if e = catalog.Verify(z, selected); e != nil {
			return Config{}, e
		}
		backup, e := retainIncompleteSetup(r, review)
		if e != nil {
			return Config{}, e
		}
		report("Previous setup saved in " + backup + ". Your Files folder is retained.")
		r.Reinstall = false
		r.ReinstallReview = ""
		return installFresh(ctx, r, z, cat, report)
	}
	c, s, err := Load(r.Data)
	if err != nil {
		return Config{}, err
	}
	// Recreate absent roots only after the user reviewed their exact locations.
	// Existing files are never cleared or replaced here.
	for _, root := range []string{c.Programs, c.Files} {
		if err = os.MkdirAll(root, 0700); err != nil {
			return Config{}, err
		}
	}
	if err = updateRoots(c); err != nil {
		return Config{}, err
	}
	guard, err := lockUpdate(c)
	if err != nil {
		return Config{}, err
	}
	defer guard.Close()
	if UpdatePending(c) {
		return Config{}, errors.New("An interrupted update has a recovery copy. Open Updates and backups and choose Recover before reinstalling.")
	}
	var overrides map[string]string
	if !r.Offline {
		bundle, e := distribution.Latest(ctx, report, filepath.Dir(c.Programs))
		if e != nil {
			return Config{}, e
		}
		defer bundle.Close()
		if e = distribution.CheckInstallerBundle(cat, bundle.Catalog); e != nil {
			return Config{}, e
		}
		z, cat, overrides = bundle.ZIP, bundle.Catalog, bundle.Overrides
	}
	target, err := repairCatalog(c.Catalog, cat)
	if err != nil {
		return Config{}, err
	}
	stage, err := newUpdateStage(c)
	if err != nil {
		return Config{}, err
	}
	defer removeUpdateStage(c, stage, report)
	report("Verifying replacement program files before reinstalling")
	fromBundle := catalog.Catalog{Schema: 1, Version: target.Version, Packages: map[string]catalog.Package{}}
	for id, p := range target.Packages {
		if cat.Packages[id] == p {
			fromBundle.Packages[id] = p
		} else {
			// A newer independently installed app must never be downgraded. Keep
			// its verified bytes; report a missing matching installer clearly.
			path := filepath.Join(c.Programs, filepath.FromSlash(p.File))
			held, e := platform.LockedBinary(path, p.SHA256)
			if e != nil {
				return Config{}, fmt.Errorf("%s needs its matching app installer before Workspace can reinstall it", p.File)
			}
			held.Close()
			if e = platform.CopyFile(path, filepath.Join(stage, filepath.FromSlash(p.File))); e != nil {
				return Config{}, e
			}
		}
	}
	if err = catalog.Extract(z, fromBundle, stage, overrides); err != nil {
		return Config{}, err
	}
	// Recheck the reviewed identity and paths after all potentially slow reads.
	current, err := InspectSetupStorage(r)
	if err != nil || current.Review != review.Review {
		return Config{}, errors.New("Workspace settings changed. Review this folder again before reinstalling.")
	}
	ops := nativeUpdateOperations
	ops.repair = true
	if err = applyStagedWith(ctx, c, s, target, stage, "", report, ops); err != nil {
		return Config{}, err
	}
	returnConfig, _, err := Load(c.Data)
	if err == nil {
		err = platform.SaveJSON(filepath.Join(c.Programs, "workspace-location.json"), map[string]string{"data": c.Data})
	}
	if err == nil {
		err = platform.Autostart(filepath.Join(c.Data, "Workspace.exe"), c.Data, c.Autostart)
	}
	if err == nil && !r.NoShortcuts {
		err = platform.Shortcut(filepath.Join(c.Data, "Workspace.exe"), c.Data)
	}
	return returnConfig, err
}

func repairCatalog(current, incoming catalog.Catalog) (catalog.Catalog, error) {
	if incoming.Version == current.Version {
		if incoming.Packages["manager"] != current.Packages["manager"] {
			return current, errors.New("This installer does not match the installed Workspace release")
		}
		return current, nil
	}
	// Reinstall the installed selection, rather than silently adding apps.
	scoped := current
	scoped.Selection = current.EffectiveSelection()
	scoped.Selection = &catalog.AppSelection{Excluded: append([]string(nil), scoped.Selection.Excluded...)}
	for _, a := range incoming.Apps {
		if _, ok := current.Find(a.ID); !ok {
			found := false
			for _, id := range scoped.Selection.Excluded {
				if id == a.ID {
					found = true
				}
			}
			if !found {
				scoped.Selection.Excluded = append(scoped.Selection.Excluded, a.ID)
			}
		}
	}
	target, err := distribution.SelectBundleUpdates(scoped, incoming, []string{"workspace"})
	target.Selection = current.Selection
	return target, err
}

func verifyRepairStage(target catalog.Catalog, stage string) error {
	for _, p := range target.Packages {
		path, e := platform.Inside(stage, p.File)
		if e != nil {
			return e
		}
		f, e := platform.LockedBinary(path, p.SHA256)
		if e != nil {
			return e
		}
		st, e := f.Stat()
		f.Close()
		if e != nil || st.Size() != p.Size {
			return errors.New("Replacement program verification failed; no existing files were replaced")
		}
	}
	return nil
}

func retainIncompleteSetup(r Request, review SetupStorageReview) (string, error) {
	// Only stopped, incomplete setups can take this path. A recognised complete
	// installation always uses the journalled repair transaction above.
	var locks []*os.File
	for _, name := range []string{"runtime.lock", "install.lock"} {
		if _, e := os.Stat(r.Data); e == nil {
			f, e := platform.Lock(filepath.Join(r.Data, name))
			if e != nil {
				for _, l := range locks {
					l.Close()
				}
				return "", errors.New("This setup is still running. Close its setup window or stop Workspace, then retry.")
			}
			locks = append(locks, f)
		}
	}
	current, e := InspectSetupStorage(r)
	for _, l := range locks {
		l.Close()
	}
	if e != nil {
		return "", e
	}
	if current.Kind != "restart" || current.Review != review.Review {
		return "", errors.New("This setup changed. Check the folder again before reinstalling.")
	}
	backup := filepath.Join(filepath.Dir(r.Data), "ReinstallRecovery-"+time.Now().UTC().Format("20060102T150405")+"-"+platform.Token()[:8])
	if e = platform.RequireNewDir(backup); e != nil {
		return "", e
	}
	if e = platform.ProtectDir(backup); e != nil {
		return "", e
	}
	type movedDir struct{ from, to string }
	var moved []movedDir
	for _, entry := range []movedDir{{r.Programs, filepath.Join(backup, "Apps")}, {r.Data, filepath.Join(backup, "Private")}} {
		if _, e = os.Stat(entry.from); os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return "", e
		}
		if e = os.Rename(entry.from, entry.to); e != nil {
			for i := len(moved) - 1; i >= 0; i-- {
				_ = os.Rename(moved[i].to, moved[i].from)
			}
			return "", fmt.Errorf("Could not retain the previous setup; it has not been deleted: %w", e)
		}
		moved = append(moved, entry)
	}
	if e = platform.ProtectDir(backup); e != nil {
		return "", fmt.Errorf("Previous setup retained at %s; its Windows permissions need attention: %w", backup, e)
	}
	e = platform.SaveJSON(filepath.Join(backup, "reinstall.json"), map[string]any{"schema": 1, "created": time.Now().UTC(), "programs": r.Programs, "data": r.Data, "files": r.Files, "filesRetained": true})
	return backup, e
}
