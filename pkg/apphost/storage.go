package apphost

import "github.com/oneworkspacestudio/one-workspace/internal/platform"

// SavePrivateSetting atomically writes JSON using Workspace's platform secret store.
func SavePrivateSetting(path string, value any) error {
	return platform.SaveSecret(path, value)
}

// LoadPrivateSetting reads a setting saved by SavePrivateSetting.
func LoadPrivateSetting(path string, value any) error {
	return platform.LoadSecret(path, value)
}

// RejectReparse rejects symbolic links and platform reparse points in a managed path.
func RejectReparse(path string) error {
	return platform.RejectReparse(path)
}
