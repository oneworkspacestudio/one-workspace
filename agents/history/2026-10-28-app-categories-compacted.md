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

Concurrent Maps promotion df63e4e was retained by rebase. Its Maps11 pin is the
only baseline95 lock difference: separately verify immutable public Maps11 ZIP
SHA681d1c46cb12432016b9afcb794b8117648772c039b00d4f6f3e0b1ffbdd6e29, compatibility,
program/notices/contract, and update that package group in the installer catalog.
All other lock entries and95 app packages must match. Existing-PC category updates
remain Workspace-only, preserving each installation's actual current app versions.
Maps delivery is separately in progress; check the update lock before applying.

Installer compile caught two legacy single-byte middle dots in the new selection
legend. Converted only those invalid bytes to UTF-8; other changed text sources
pass UTF-8 decoding. No release from the failed build was published.

## Release96 published; primary delivered; secondary access blocked

Immutable Workspace v0.3.0-preview.96 source6950425e329aa0aec026bb9e34ff2a0c885146bf.
Required scripts/publish-release.ps1 built signed Android6 and verified all26
uploaded assets plus anonymous APK. Source, notices, provenance/checksums included.
No Actions. APK SHA b475f48edee8bf488bbbfcde4365ca3390c4e80028d2d6125c8e2d5ed194d254;
manager SHA bc11fb6bb68b93844ce0b47a627816f42f30f472a91719ef8769829e7661a10a.
Public channel promotion initially conflicted with Maps delivery handoff cb29f96;
rebased only the generated-index commit and pushed685cb12, retaining both. Exact
anonymous raw-main channel contains96/all26 matching assets. Published assets/tag
remain immutable. Evidence under docs/evidence/app-categories-20261028/.

Normal protected primary Workspace-only95→96 update completed.30 existing app
contracts and process executable hashes, owner settings, decrypted credentials/
clients and selection retained (comparison only in memory).29 services healthy,
transaction complete. Maps11 retained. Deployed public navigation.js bytes match
release96 source. No new owner UI/SSO, phone, fresh installation or reboot tests.
The no-proxy loopback path is needed for operational HTTP checks in this host.

SECONDARY NOT INSTALLED. Public health has29 running services; its navigation.js
still differs from96. Last confirmed manager95/Maps11 (Maps deliverycb29f96).
Browser control became disabled in Work settings; requested user enable via async
question, no answer yet. Do not claim both delivered. No remote apply/reinstall for
this task has been submitted. Need user enable browser control, then create/use a
dedicated authenticated https://workspace-2.m8e.co.uk/updates tab. Do not take over
the Maps task's deliverable tab. Refresh/check current controller, uncheck EVERY
unrelated update (MyAI14/Chat10 etc), select ONLY Workspace96, submit once, follow
reconnect/progress through verified backup/startup/complete. Preserve installed
app versions/settings. Close idle old updater after success so next launch uses96.
Verify manager96,29 running services, category navigation and completed recovery.
If normal owner session expired, ask user to sign in privately while preserving
all authorization already given to install on both. No full bundle rebuild needed.

I: read stalls forced an isolated named C build scratch; only task-owned duplicate
build archives were removed during publication.16 residual task headless browser
processes were identified by their exact preview profile prefix and stopped;
loopback19180 verified closed. Final C scratch cleanup is recorded below after
preserving source/evidence and confirming no active process references. User's
installations, latest verified recovery, other projects/caches and signing keys
must remain intact. Desktop grouping reviewed with synthetic public metadata;
phone capture inconclusive and additional preview blocked by automatic approval
review (no detailed reason), so physical/mobile visual acceptance is not claimed.

### Cleanup and final resume point

Primary success is final; do not repeat that update. Full named C scratch deletion
was rejected by automatic approval review with no detailed reason. It was not
retried through another delete mechanism. A narrower cleanup verified the two
published installer/bundle copies against GitHub digests and removed only those
files; together with earlier intermediates, about6.28GiB of owned duplicate build
files were removed. C had11.8GiB free afterward. The named scratch still contains
about2.43GiB of public source/build caches and small receipts; preserve/report it
until manual cleanup or an approved policy change. Permanent signing keys and
owner data/recovery were untouched. No task preview still listens on19180.

Permanent I source checkout fast-forwarded to5759958 (release/channel/delivery
records). This follow-up records cleanup only. Remaining delivery requires user
enable browser control in Work settings; the async request has no response. All
other authorized independent work is complete. Secondary has not received96 and
its old navigation bytes were verified. APK6 is published, not phone-installed.
Public release remains https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.96.

### 2026-10-29 category refinement
Renamed Development display label to Creative & Development, preserving stable development ID. Added my3d, studio (Photo Editor), create, slicer and push; moved Push out of Administration. Existing Git/Hosting/MyAI retained. Updated Go category policy, web picker fallback, installer labels, Android grouping and README consistently. This group remains optional on fresh installs; existing selected apps are preserved. Source diff checked; no build, release or installation performed for this follow-up.

### Release101, 29October
Rebased category refinement on newer upstream Finance100/Forms2/Shield4 history. Source aec20de; published immutable v0.3.0-preview.101 with 30 assets and signed Android7 (approved permanent certificate), all asset digests and anonymous exact raw-main channel verified. Publisher used locally; no Actions. Old gh2.5 could not upload draft; used digest-verified official current CLI. Channel generator needed retry after IncompleteRead; channel promotion7bd007a pushed. Catalog tests passed and all32 bundled app contracts/payloads verified against published100 baseline before reusing. Only manager and release notice changed in desktop bundle.
Local100→101 normal protected Workspace-only update is in final recovery cleanup; 31 services healthy, deployed navigation exact source. Existing versions/selections/credentials verification helper still awaits transaction completion. Secondary browser tool continues to report Work browser disabled; no secondary apply. Build/release artifacts on I:/Projects/tooling/workspace-creative-20261029. Removed exact disposable C embedded bundle after matching published digest (1.21GB); owner recovery/data untouched.


Final101 local delivery completed: transaction complete,31 services running, all32 pre-existing app contracts and process executable bytes retained, selections/settings/decrypted identity and client credentials unchanged; manager bytes match published101. Runtime/navigation evidence and publication/delivery receipts under docs/evidence/creative-categories-20261029. Cleanup retained newest verified1.824GB recovery, removed one older automatic backup (~122.5MB). Secondary remains uninstalled by this task due browser-control restriction; no remote apply submitted. No phone/reboot/Linux runtime claim. Source release tag remains immutable aec20de.

### Download recovery follow-up
User screenshot from secondary98 shows download-stage failure at10m48. Code inspection confirms shared staged release downloader has one attempt, client10-minute timeout, and conflates interrupted copies/size mismatch with checksum mismatch. Exact remote cause unverified without logs. Added three bounded attempts, Range resume for incomplete transfers with strict Content-Range validation, safe full restart when ranges ignored or checksum fails, context-aware backoff, final SHA256 and Sync. Added interrupted-response/resume and corrupt-response retry tests. Added plain recovery guidance and GitHub-access check/retry buttons to maintenance failed-download view; same-tab retry preserves submitted selections; after reload requires review. No hash bypass. Initial installer full-bundle downloader remains separate. Full distribution tests exposed stale Slicer-default assertions from category change; updated to expect optional Slicer absent/download-free. Source-only follow-up; not released or installed.

### Release102 and local delivery
Published immutable v0.3.0-preview.102 from a816cec through required local publisher,30 assets including permanent-key-signed Android8. All public asset sizes/SHA256 verified, anonymous APK download passed, exact raw-main channel verified; promotion e290f6b. Baseline101 app contracts/payloads verified before reuse; manager/release notice replaced. Maintenance and distribution tests passed, including simulated truncated/resumed and corrupt/retried transfers. Normal protected local101→102 update completed; all31 services running, all32 old app contracts and process executable bytes retained, selections/settings/decrypted credentials unchanged, manager matches published bytes. Newest verified1.825GB recovery retained; one older update backup removed. Disposable C build payload removed only after digest match (~1.21GB); release artifacts retained on I:/Projects/tooling/workspace-recovery-20261029. Evidence under docs/evidence/download-recovery-20261029. Secondary98 not updated: fresh explicit Chrome connection still rejected by tool with Work-settings-disabled message despite enabled local config. No remote apply submitted. No live-browser recovery-button/phone/reboot/Linux-runtime acceptance claimed. Original initial-installer bundle download remains a separate single-attempt path; staged updater retries are shipped.
