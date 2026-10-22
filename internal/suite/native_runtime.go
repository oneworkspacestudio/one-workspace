package suite

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Stalwart's Windows binary imports the release MSVC runtime. Windows does not
// install these libraries on every PC. Deploy the verified app-local set beside
// the private engine, rather than relying on a developer machine's System32.
var nativeRuntimeDLLs = []string{"concrt140.dll", "msvcp140.dll", "msvcp140_1.dll", "msvcp140_2.dll", "msvcp140_atomic_wait.dll", "msvcp140_codecvt_ids.dll", "vcruntime140.dll", "vcruntime140_1.dll", "vcruntime140_threads.dll"}

func prepareSharedRuntime(c Config, enginePath string) error {
	// Older catalogs remain launchable with their existing system prerequisites.
	_, supplied := c.Catalog.Packages["msvc-msvcp140"]
	if !supplied {
		return nil
	}
	target := filepath.Dir(enginePath)
	if err := platform.RejectReparse(target); err != nil {
		return err
	}
	for _, name := range nativeRuntimeDLLs {
		id := "msvc-" + name[:len(name)-4]
		p, ok := c.Catalog.Packages[id]
		if !ok || p.File != "runtime/msvc/"+name {
			return errors.New("the Microsoft runtime package is incomplete; repair Workspace with the latest installer")
		}
		source, err := platform.Inside(c.Programs, p.File)
		if err != nil {
			return err
		}
		held, err := platform.LockedBinary(source, p.SHA256)
		if err != nil {
			return fmt.Errorf("verify Microsoft runtime %s: %w", name, err)
		}
		dest := filepath.Join(target, name)
		if _, err = os.Lstat(dest); err == nil {
			err = platform.RejectReparse(dest)
			if err == nil {
				var b []byte
				b, err = os.ReadFile(dest)
				if err == nil {
					sum := sha256.Sum256(b)
					if hex.EncodeToString(sum[:]) == p.SHA256 {
						held.Close()
						continue
					}
				}
			}
			// Never overwrite a loaded or unfamiliar runtime during normal startup.
			held.Close()
			return fmt.Errorf("Microsoft runtime %s needs a protected repair; existing files were retained", name)
		} else if !os.IsNotExist(err) {
			held.Close()
			return err
		}
		err = platform.CopyFile(source, dest)
		held.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
