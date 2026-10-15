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
