//go:build linux

package platform

import secure "github.com/oneworkspacestudio/one-workspace/internal/portablesecure"

// Native apphost storage uses the existing shared Linux ownership/permission
// and per-installation encryption implementation. Windows DPAPI is unchanged.
type DataVolume = secure.DataVolume

var CheckDataVolume = secure.CheckDataVolume
var ProtectDir = secure.ProtectDir
var RejectReparse = secure.RejectReparse
var Lock = secure.Lock
var OpenBrowser = secure.OpenBrowser
var protectSecret = secure.SealSecret
var unprotectSecret = secure.OpenSecret
