# Shield private-history update discovery, 30 October 2026

User confirmed workspace-2 still loads old index-CzNDSlaq.js and its Settings
.map crash; My Shield4's fixed build is index-B3l6dnYm.js. Shield4 was published
in Workspace97 at2026-10-28T21:54:32Z; public latest107 catalog also pins Shield4.
The public channel intentionally has an empty my-shield release history.

Found discovery defect in CheckUpdates: different app sources with missing release
timestamps retain the installed contract. Private histories are omitted, so both
timestamps are zero and Shield3 is retained even with verified suite Shield4.
Fixed by allowing a strictly newer canonical same-repository app version from
the verified suite catalog when the public index omits history. Older pins still
retain installed app. No private release metadata is exposed. Regression exercises
public discovery and app-only selection (one payload, retained manager), plus
rejection of a lower suite app version. Targeted regression passed. Whole
internal/distribution package check passed (13.620s).

Source changes: internal/distribution/updates.go, private_suite_test.go.
This correction is NOT yet in a published updater or installed on either PC.
Needs complete test check, commit/push and next unclaimed suite release via mandatory
signed-Android publisher; preserve concurrent latest main/app pins. Do not overwrite
97/107. Need supported remote protected app-only Shield4 update after corrected
updater installed; old updater can falsely say up-to-date. The manifest401 remains
separate: browser link lacks use-credentials, Workspace protects app paths; do not
open app APIs publicly to silence it. Existing primary Shield4 was verified earlier.
Native Chrome automation now observes via @oai/sky but input was stopped for inability
to determine current browser URL; no remote update submitted.

## Release108 execution in progress

User explicitly requested push, release and installation on both authorised PCs.
Release source e0faae9 (includes e0f696c discovery fix) pushed; immutable tag
v0.3.0-preview.108 pushed to reserve this source. No release108 published yet.
Clean build checkout I:/Projects/tooling/shield-discovery-20261030/source.
Ordinary scripts/rebuild.ps1 uses unchanged managedSeed103/Remote5/Tickets1 pins.
All34 catalog apps validated. Manager108 and signed Android10 prebuild complete;
offline bundle1292354447 bytes complete, setup EXE build currently running.
Build exec session18482, log build108-cached.log, PASS marker build108-status.txt.
Original redundant seed download stopped (only exact task Python90348); cached
published103 ZIP verified SHA cc61e4bf...fe76a4 and supplied through -SeedArchive.
Partial download remains outside release108 output, not used/published.

All current filtered distribution/catalog/maintenance/suite update/controller tests
passed. Before-fix overlay regression failed specifically available=false/want=true;
fixed regression passes and full distribution package passed earlier. Logs in task
folder, no owner data used as fixtures. Signed Android prebuild log android108.log,
exec87546 (likely completed); required publisher rebuilds incrementally again.

Next: await builder PASS; run prepare108.py, then publish108.ps1 (mandatory signed
Android publication), then verify108.py. Do not bypass publisher/exported-package
hash checks or overwrite published assets. Helpers prepared in task folder; no
publication started yet. Standard publisher exports per-package assets; large upload.
After verified publication/channel proof public108-verification.json, run
Deploy-Primary108.ps1 (guarded manager-only protected update, hash/config/all-secret
comparisons in memory, recovery and health checks). Source helper parse checked.
Fresh observation overrides earlier107 handoff: META-GAMING current
C:/tmp/workspace/Private originworkspace.m8e.co.uk is manager106, binary106 matches
catalog, Shield4,34 apps, transactioncomplete. Do not assume107 installed.
C had4.8GiB free before final setup link; only C/F are NTFS, F effectively full.
No unrelated files/backups cleaned. Existing C my-shield-release4 source clones
275MiB clean but retained; no need to delete until actually necessary.

Remote workspace-2 still needs corrected updater then Shield4 through supported
maintenance UI; no remote change submitted. Async request asks user bring Chrome
Updates URL/address bar to foreground. Native @oai/sky works for observations but
previous Ctrl+T was blocked by browser-URL verification while Codex occluded Chrome.
Try documented activate_window on freshly returned Chrome first once independent
release work is finished, then observe/one-action loop. Do not operate Codex UI,
security dialogs, credentials or remote terminal. Existing Remote app setup modal
is unrelated and must not be changed. MiniPC excluded. No claim of both delivered.

## Build complete; draft transfer active

Build108 passed:34 apps,307 packages, manager17538048 bytes, bundle1292354447 bytes,
setup1306242560 bytes. Signed Android10 SHA ffa2e02f922eae4bed45c978cf8dd068654a7ca0e24c68b78667692dfec7c400.
Mandatory publisher passed signing/package checks and exported all307 package
assets, but gh CLI draft lookup returned release-not-found. Draft429393129 is
confirmed target e0faae924422a7124ea6244356cb59565e6375ff, tag108, not published.
No assets were uploaded by failed gh call. Task-owned transfer108.py verifies
checksums and streams to that exact draft ID; it never publishes. Exec session29982
currently transfers assets (log transfer108.log), then automatically resumes the
mandatory publish108.ps1 (log publish108-resume.log), then verify108.py. Do not start
a competing upload or overwrite assets. Last19:12 UTC draft23 assets, bundle upload
active. Release/helper paths all I:/Projects/tooling/shield-discovery-20261030/.

After successful guard/channel verification, proof public108-verification.json
contains exact asset SHA256s and enables guarded Deploy-Primary108.ps1. All input
hashes, config/origin, transaction, app contracts, unselected payloads and secret
values are checked; manager-only update. Helper uses large buffered SHA streams.
Do not apply until proof exists. Native tests/controller checks passed and baseline
overlay regression rejected the old discovery behavior. Primary current106 binary
matches106 config, originworkspace.m8e.co.uk,34apps/Shield4; earlier107 receipt does
not describe current state. C free later4.8GiB (temporary link used some), no cleanup
performed. Only C/F NTFS, F nearly full; do not put private recovery on exFAT I/E/H.

Remote still unchanged by this task. Native Chrome activation-first via fresh
@oai/sky returned window may restore URL confidence; only attempt after independent
publication/current-PC work completes in case tool forces turn termination again.
User foreground Updates-page question remains pending. Never automate unrelated
Remote setup/PIN/security dialog, Codex UI or a remote terminal. GUI may need old idle
updater Close updater + refresh before manager update; inspect busy state first.
Then install manager108, reconnect corrected controller, select only Shield4 and
verify Settings render. Existing apps/credentials/recovery must be retained.

## Publication and primary deployment completed

Release v0.3.0-preview.108 published through mandatory signed-Android publisher,
tag e0faae9, release429393129, channel57d62f6. All333 immutable asset sizes/digests
verified against checksums and anonymous public channel; signed Android10 downloaded
anonymously and verified. Private Shield history deliberately remains empty, now
covered by suite catalog fallback. Actions disabled. Concurrent Remote7 changes
preserved by fetch/rebase before regenerating channel; no immutable assets replaced.

Primary supported manager-only protected update106 ->108 completed. All owner
configuration and decrypted-secret comparisons passed in memory; all unselected app
contracts/package bytes retained; exact public manager digest and recovery manifest
verified. Shield remains4/running. Final health32 running, Cinema failed readiness
HTTP503 after restart; do not claim all services healthy. Native transaction complete,
latest verified1840.1MB recovery retained; normal cleanup removed one older backup.
Evidence docs/evidence/shield-discovery108/. No owner identity/database fixture used.

Remote workspace-2 remains pending. Activate freshly selected native Chrome before
observing Updates; use supported manager108 then Shield4 flow, preserve unrelated
apps. Do not interact with unrelated Remote setup/security modal. MiniPC excluded.
Manifest credential-fetch401 is separate and still unresolved; app4 fixes Settings
map crash,108 fixes global update discovery. AI stays off/Monitor-only.
