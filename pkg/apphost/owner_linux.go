//go:build linux

package apphost

import secure "github.com/oneworkspacestudio/one-workspace/internal/portablesecure"

func validateOwner(root string) error { return secure.ValidateOwner(root) }
