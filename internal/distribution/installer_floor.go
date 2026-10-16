package distribution

import (
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
)

// A newer setup must not download an older manager and silently lose its features.
// This can occur when trying an unreleased local installer against public releases.
func CheckInstallerBundle(embedded, incoming catalog.Catalog) error {
	if newerWorkspace(embedded.Version, incoming.Version) {
		return fmt.Errorf("the newest published bundle (%s) is older than this installer (%s). Wait for the matching release, or deliberately run this installer with --offline to use its bundled apps", incoming.Version, embedded.Version)
	}
	return nil
}
