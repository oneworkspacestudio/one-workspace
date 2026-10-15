# My Tube package integration

User requested the supplied My Tube SVG in its README and adding My Tube to the
One Workspace README and installed package. The supplied artwork is copied without
path or gradient changes to internal/branding/icons/my_tube.svg. Both the existing
illustrated hero and current-logo overview are retained; overview adds My Tube.

The component lock pins My Tube v0.1.0-preview.1 and the published ZIP SHA256.
Build-org verifies the public release, digest and bytes, includes MyTube.exe and
its three license notices, and produces twelve app entries. All previous pins
remain, including Notes preview4 and the older approved Music executable. Never
include the withheld Music hash or use GitHub Actions for Workspace builds.

Manager provisions a confidential OIDC client through the existing generic path,
uses base port+13, and generates managed settings outside the dedicated private
My Tube data directory. The explicit Workspace owner subject is its initial admin.
Drive import/export targets the direct loopback /drive mount to avoid the portal's
separate browser session requirement. Existing settings are preserved. Logout
clears My Tube cookies at root and /my-tube/. Optional FFmpeg helpers are separate.

Verification: go test ./... and go vet ./... passed. PowerShell scripts parse;
both SVG copies match the supplied file; regenerated app overview was inspected.
Full local package build and fresh NTFS installation verification follow. Named
build/test staging: tooling/workspace-my-tube22, with private disposable data under
the dedicated C: cache. Do not touch the current live installation.

Fresh-package verification exposed a pre-existing Calendar/Contacts ambiguity when
several managed workspaces are registered: their CLI checked discovery before the
environment-selected root. The catalog now supplies --shared-data explicitly,
and its launch-contract check requires it. Component pins and binaries are unchanged.
The first test installer succeeded and My Tube started; the suite then stopped at
Calendar. Its separately managed shared service was explicitly stopped as part of
cleanup. Rebuild and repeat the fresh install before publishing preview22.

## Published and verified

Complete release v0.3.0-preview.22 is published:
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.22
Source tag is c4d953c3d9cc14088a9ab132f18540646317eaa0. All thirteen public
asset sizes and GitHub SHA256 digests match the local release. Every bundled
package file matches the twelve-app catalog. Notes preview4 and the previous
Music hash are retained; withheld Music is excluded. No Workspace Actions build.

The rebuilt single-file offline installer passed a second fresh NTFS-data
installation, all managed suite startup steps, real MyAccounts code/PKCE/nonce
callbacks for portal and My Tube, explicit owner admin, private managed settings,
shared-origin web shell, launcher/logo and scoped logout. This is native/API
integration evidence on the existing Windows host, not clean-host, browser UI,
reboot, public-hosting or full app acceptance. Existing preview limits remain.

The released installer also anonymously downloaded and verified the newest public
bundle and reported all twelve apps including My Tube preview1. Both public
READMEs reference the supplied logo and hosted SVG bytes match the original.
Git archive applies Windows line endings; its SVG differs only by newline format.

All task-owned native app, shared-service and engine processes were stopped and
verified absent. Disposable private data, test program/file roots and installation
markers were removed after path/process/reparse checks. The current live suite
was not upgraded. Build, sanitized acceptance and publication receipts remain
under tooling/workspace-my-tube22. Published tags/assets are immutable; this
handoff follow-up is documentation only.
