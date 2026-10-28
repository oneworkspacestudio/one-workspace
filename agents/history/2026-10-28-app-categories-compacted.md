# App categories, default selection and both-PC delivery

User requests General/Administration/Development sections, fresh installs selecting
General+Administration and other apps optional, then push/release/install on both
owner PCs. No Actions or subagents. Work in I:/Projects/one-workspace-android,
codex/android-workspace. Clean initial884a67b fast-forwarded to public51aa2f3;
read 2026-10-28-my-shield-integration-compacted.md. Both PCs began on Workspace95 with
29 healthy services; primary C:/tmp/workspace/Private, remote I:/Workspace/Private
at workspace-2.m8e.co.uk. Shield3, MyAI15, Maps10 and all other pins must be retained.
Remote has Workspace95; its previously open maintenance child appears stale. Close/reopen the idle updater before applying96; preserve safety checks.

## Source implementation

- internal/catalog/categories.go centralizes 3 category names/policy. Admin IDs:
  accounts/backups/monitor/shield/push/updates/app-status. Development: gitea, hosting and myai. Remaining apps General. DefaultInstall true for General+Admin.
- Categories are derived, NOT added to strict serialized App/catalog contracts,
  preserving older-manager update compatibility. API adds Category/CategoryID.
- Portal groups tiles into category sections; owner administration controls stay
  inside Administration. Existing owner checks/4-field catalogTile retained.
  Home template has fallback Tiles rendering for legacy renderer callers.
- Shared app picker groups headings while retaining scrollable body/fixed close.
  Android grid groups matching API categories with fallback IDs; search hides
  empty headings. Android incremented5→6 for this source change.
- Fresh nil selections now choose category defaults, not legacy Photos boolean.
  Explicit selections still honored. Removed legacy config Photos filtering.
  WantedApps(previous) preserves installed/opted-in apps and explicit exclusions,
  adding only new default IDs. Used by online installer reconciliation and update
  new-default detection; optional future additions no longer install silently.
- Setup API offers category/defaultInstall; app checkbox sections/legends preselect
  defaults, Development unchecked. Existing dependency toggle logic retained.
- catalogcheck also rejects default apps depending on optional apps. No tests
  added/run per developer instruction. Candidate manager builds and --help starts
  successfully (includes template initialization). Full build pending.
- Build scripts default to Workspace96. README/release notes describe policy.
  Need update Android README/download links to suite96/APK6 before release commit.

## Remaining

Review generated UI visually using only synthetic public metadata/static preview
or post-install real owner verification; never live owner data as development
fixtures. Browser-screenshot skill read from .agents/skills/browser-screenshot.
Computer-use skill/guidance/confirmations read in preparation for browser-only
remote updates. Native Computer Use forbids terminal automation; use authenticated
remote Workspace maintenance UI rather than Chrome Remote Desktop PowerShell.
No browser control started yet; CUA inventory first. If owner login is expired,
ask user to sign in privately while continuing independent build/primary work.

Build full org96 locally with existing caches on I and Go1.27.3. Baseline wrapper:
I:/Projects/tooling/my-shield-integration-20261028/build95.ps1 (do NOT rerun old95
deployers; use as reference only). Component cache I:/Projects/tooling/workspace-org-components;
old release95 artifacts under tooling/my-shield-integration-20261028/release-workspace95.
New tooling root I:/Projects/tooling/workspace-categories-20261028. Keep caches and
artifacts on I; C had23GB free but no new NTFS private fixture is required yet.

Commit/push exact release source, use REQUIRED scripts/publish-release.ps1 with
permanent signing key outside Git for suite96+Android6. Private signing loader
reference tooling/android-servers-20261027/publish-android4.ps1; never print key
values. Every suite release must contain signed APK/source/notices/provenance.
Current publisher old gh may fail uploading a draft (see Shield handoff); if so
upload by exact draft ID with streamed transfers, preserving all verification,
then resume the required publisher. No published assets may be overwritten.
Refresh/promote public channel and verify anonymously, then normal Workspace-only
protected updates on both hosts. Preserve all app contracts/bytes/settings/secrets,
verify health, complete transaction and latest usable recovery. Current native
primary95 deployment reference script Deploy-ManagedShield95.ps1 can inform a new
96 preservation audit, but its old hashes/version checks must not be rerun.

Task is in progress. No release96 or live update has been completed yet; source
commits c0d1043 and 0b2b16a were pushed to public main.

## Build and visual review in progress

First source implementation committed c0d1043. Candidate manager compiles/starts.
Full org build caught My AI depending on optional Git; My AI now Development,
keeping its real dependency contract and default selection coherent. Default
dependency guard retained. Desktop/phone synthetic previews use only public
metadata at loopback19180, all assets/task artifacts on I. Added category jump
links to avoid long scrolling on phones. Preview PID83236 must be stopped.
Secondary maintenance UI opens authenticated; refreshed plan shows its actual
MyAI14/Chat10 selections, so preserve installed contracts per machine rather than
assuming identical pins. No installation or public release96 yet.

## Build storage fallback

I: reads stalled during the org build and headless preview. Task-owned build
PID37428 and preview83120 stopped; the unfinished build must not be published.
Clean checkout of pushed source now C:/Projects/.cache/workspace-categories-20261028/repo.
C had20GB free before full staging. Remove this named scratch after delivery.
Rebuild manager/setup locally with official Go1.27.3 checksum verified from go.dev.
Reuse only the immutable published95 bundle after its GitHub SHA256, every embedded
package size/SHA256, catalog and exact unchanged organization lock are verified.
Replace manager, suite version and current release-status notice; all app programs
and other notices stay byte-identical. Required signed publisher remains mandatory.
Desktop synthetic grouped layout reviewed. Phone fallback screenshot was clipped
and not conclusive; added explicit narrow-viewport styles. Additional preview was
blocked by automatic approval review with no detailed reason. Browser control is
now disabled in Work settings; user asked to enable while independent work proceeds.
Secondary requires owner browser update when access resumes. Primary95 transaction
is complete but Workspace is currently stopped (localhost8990 refused), so verify
startup after normal protected installation. No production data fixtures used.
