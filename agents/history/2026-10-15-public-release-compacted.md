# Public releases and latest bundle installation

User requests GitHub push and latest Git versions during install; explicitly
provides oneworkspacestudio organization for public repositories. Public copies:
one-workspace and stalwart. Retain original personal remotes. Disable Actions
before push; all Windows artifacts built locally. No change to existing running
C:/tmp/workspace required for installer-only update.

Normal GUI/headless install resolves newest published non-draft GitHub release
with workspace-bundle.zip (previews included) from oneworkspacestudio/one-workspace,
ordered by published_at. API metadata + whole-asset SHA256/size + every catalog
package verified before writing installation folders. No silent stale fallback.
Bounded HTTPS, restricted GitHub redirects, strip Authorization off asset CDN,
public installs need no GitHub account/CLI. Optional GH_TOKEN/GITHUB_TOKEN only
for explicitly private build configurations; no embedded credentials.

--offline deliberately uses embedded versions; --check-updates downloads/verifies
without installing; --resume monitors existing installation without fetching.
Latest means latest compatible published suite bundle, not arbitrary Git HEADs
or independently chosen upstream apps. Existing installations are not silently
upgraded. Unknown/incompatible catalogs stop with newest-installer guidance.

Download percent/MB reported through install progress; GUI version updates to
installed catalog. build.ps1 emits workspace-bundle.zip and supports -ReleaseRepo.
Tests: release ordering/draft/seed filtering, asset scope/digests, credential
redirect stripping; existing catalog/suite/setup UI tests pass. Real public
release download verification and publication details appended after success.

Preserve all prior uncommitted Mail/Workspace fixes and concurrent Drive design
work. Bundle includes latest local Drive/Photos repair binaries and their exact
source snapshots; Git/Music baseline retained with original component notices.
No separate live installation used for tests; no clean-host/reboot claim.


## Published state

Public source main implementation commit2dedc67; Windows release:
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.15
Six assets uploaded and each GitHub size/digest matched the local SHA256 manifest
before publication: installer, manager, catalog, bundle, Drive/Photos corresponding
source ZIP and SHA256SUMS. Full bundle preserves Git/Music baseline and notices.
Main tracks studio/main; personal origin remains. Actions disabled and confirmed.

Build source C:/Projects/.cache/workspace-public-release; Mail source/build
C:/Projects/.cache/mail-inbound-fix version0.1.0-preview.12. Normal local runtime
was not upgraded/restarted during publication. Versioned ignored local copies
WorkspaceSetup-v0.3.0-preview.15-windows-amd64.exe, Workspace-v0.3.0-preview.15.exe,
catalog/SHA256SUMS/bundle suffixed v0.3.0-preview.15. Canonical old setup EXE was
locked by an existing preview and was retained; use versioned/new GitHub download.
Ignored source payload ZIP updated to the new release. No user preview was stopped.

components.lock.json/rebuild.ps1 now seed from this public verified release,
retain all component versions and avoid the old private core-only seed. Script
syntax checked; no second full rebuild run for this build-only follow-up.

Anonymous compiled --check-updates initially saw no release during GitHub's brief
publication propagation delay, failed closed, then selected the correct release
and started downloading without a token. Completion verification recorded below.
