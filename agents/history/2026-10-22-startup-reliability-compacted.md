# Chat / Calendar / Contacts startup reliability

User screenshot shows three failed startup checks, not failed package extraction.
Actual other-machine logs unavailable; do not claim its specific cause confirmed.
Started clean from087fe01 on codex/startup-reliability; public68 baseline retained.
Repository-history loaded; no agents or Actions. No owner/live identity fixtures.

Isolated public68 Accounts/Chat10/Calendar3/Contacts3 all started successfully on
this developer PC. PE imports prove shared Stalwart requires MSVCP140.dll,
VCRUNTIME140.dll and VCRUNTIME140_1.dll; bundle lacked MSVC prerequisites. This is
a confirmed packaging gap, not proof of the screenshot's remote cause.

Added Microsoft-signed app-local release CRT packages, private engine deployment
with catalog hashes/reparse checks and unfamiliar-byte retention. Official
Microsoft VC_redist.x64.exe14.51.36247.0 downloaded; Authenticode valid. Extracted
embedded CABs without a system installation; no System32 copying. Provenance
receipt required before delivery. App-local copies require future security updates.
Integration3 prevents older setup versions silently omitting the required packages.

Shared manager stdout/stderr previously discarded; now protected shared/startup.log,
early-exit observation and bounded Windows memory retries. Startup failures now
show fixed categories, guidance and log locations, keeping raw provider errors
private. Chat health timeout3min. Setup offers credential-free status report.

Artifacts I:/Projects/tooling/startup-reliability-20261022; NTFS fixtures/tests
C:/Projects/.cache/startup-reliability-20261022. Synthetic public68 repro manager
and shared service stopped; retain roots until final cleanup/process checks.
Fake APPDATA/LOCALAPPDATA + precreated Start Menu Programs mandatory for GUI tests.

Next: finish targeted/all-package tests, final native acceptance incl app-local DLL
module paths and startup report/retry, all25 selected; build immutable69+ only after
checking GitHub releases for version collisions. Preserve all25 app contracts and
payloads, add CRT/provenance, replace manager+release notice only. Publish verified
assets and copy installer to Downloads. Owner67 stays untouched. No clean-host,
reboot, remote-host or external DNS acceptance claim.


Wider import review includes Mail17's private Stalwart engine, which extracts to
Data/mail/runtime. Manager supplies the same pinned app-local CRT there and keeps
packages for Mail-only selection. Shared and private Mail storage remain separate.
Runtime packages belong to the manager repair group, so protected repair of an
existing68 installation adds prerequisites without app upgrades/downgrades. Tests
cover this scoping, all three consumer selections and status report privacy.
All Go packages/vet and setup JS syntax passed after these changes. First local69
candidate93f86f2 only covered the shared engine; unpublished artifacts are retained
as candidate-first69, never public assets. Final69 must use the next source commit.
