package suite

import (
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"strings"
)

func lockRemotePayloads(c Config, a catalog.App) (func(), error) {
	var held []*os.File
	closeAll := func() {
		for _, f := range held {
			f.Close()
		}
	}
	if a.ID != "remote" {
		return closeAll, nil
	}
	for id, p := range c.Catalog.Packages {
		if !strings.HasPrefix(id, "remote-") {
			continue
		}
		f, err := platform.LockedBinary(filepath.Join(c.Programs, filepath.FromSlash(p.File)), p.SHA256)
		if err != nil {
			closeAll()
			return func() {}, fmt.Errorf("My Remote runtime %s failed integrity verification: %w", id, err)
		}
		held = append(held, f)
	}
	if len(held) == 0 {
		return closeAll, fmt.Errorf("My Remote runtime packages are missing")
	}
	return closeAll, nil
}
