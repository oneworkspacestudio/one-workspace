# Installer status and continuation

User left an installation on another machine; webpage returned to first screen.
They request clear results, recognition of existing/in-progress apps and continue.
Actual remote failure/logs unavailable; do not claim to have diagnosed that host.
Use named synthetic NTFS only; preserve owner66 and its newest rollback. Current
task fixes/replaces the installer; no owner runtime update is needed to test it.
Prior push/release authorization continues. No agents or GitHub Actions used.

Started clean on codex/slicer-flow-release, then branch codex/installer-resume.
Latest verified baseline is public Workspace66 (all25 apps/125 payloads), not the
older58 in this chat. Current owner66 includes independent Chat9/Mail17.

Found reload stripped capability hash without retaining it; active poll did not
restore progress screen. Fix session-scoped capability retention, restore active
screen, repeated reconnect with explicit lost-connection details. Durable status
records phases, errors, timestamps and progress outside the incomplete Data tree.
Protected user-config hint lets another installer process discover saved roots.
Parent install lock governs active vs interrupted, preventing duplicate writers.

DPAPI+ACL-protected request/config checkpoints bind original paths, ports/catalog
and credentials. Continue verifies/reuses matching package files and retains bad
partial bytes before replacing only the recorded package. Same recovery admin,
owner and existing client secrets reused after interrupted identity creation.
Fresh identity behavior stays unchanged. Completed configuration is adopted for
startup/retry without reinstall. Unknown older checkpoints show last saved phase
and protected reinstall recovery instead of pretending they can be resumed.
Uncertain external domain changes fail closed for explicit recovery review.
Fresh install deadline increased from15m to2h for large/slow downloads.

Go suite/setup/catalog tests passed, including lock vs interruption, encrypted
credentials, package reuse/partial retention, capability and origin protection.
First native Windows candidate passed second-process attachment, refresh after
hash removal, killing original installer mid-extraction, continuation with original
password, completed discovery after process exit, startup and real Accounts login.
Screenshots reviewed; one summary text was corrected to say interrupted rather
than still continuing. All-package tests/vet/final exact-build acceptance pending.

Artifacts I:/Projects/tooling/installer-resume-20261022. C named cache
C:/Projects/.cache/installer-resume-20261022; fake APPDATA/LOCALAPPDATA and existing
fake Start Menu Programs are mandatory for native tests. No real shortcut writes.
Go1.27.2 caches/build temp on I; tests require GOTMPDIR on C NTFS. Git per-process
GIT_CONFIG_GLOBAL=NUL/GIT_CONFIG_NOSYSTEM=1 plus exact safe.directory.

Next: finish tests and final exact native acceptance; export immutable source67,
retain all25 app contracts/bytes from verified66 while replacing only manager and
suite release notice, verify all payloads/icons, publish67 assets/hashes and copy
installer for the user. Preserve earlier tags/assets and owner66. Record final
receipt and stop/remove only named fixture roots after process/reparse checks.


## Reconciled concurrent delivery

Another authorized Slicer task completed public67 (source5feafa8) and owner67
while this installer change was uncommitted. Its committed app changes were
preserved. Final installer release is68, retaining all25 public67 app contracts
and payloads including Slicer4. Candidate67 in this task is only a local fixture,
not the unrelated public67 release. Do not move or replace that public tag.

All-package Go tests, vet and setup JS syntax passed. Extended native continuation
after identity exposed a duplicate private-manager CopyFile; fixed by hashing and
reusing the recorded manager, or retaining partial bytes before replacement.
Final68 exact-package tests must cover this path before publication.


## Final delivery complete

Public Workspace68 source f9e9c96054507c16f40faf6cf4f674bc41c10375,
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.68.
All uploaded assets match their tested local sizes/SHA256;125 bundled payloads
verified and all25 application contracts/bytes retained from public67. Seven
Windows icon frames verified. Exact public installer copied/hash-checked as
Downloads/WorkspaceSetup-0.3.0-preview.68-windows-amd64.exe.

Final exact Windows68 tests passed: second-process attachment; reload after
capability hash removal; killing installer during extraction; checked continuation;
replaying a checkpoint after account/client creation with unchanged owner ID and
original password; completed discovery after installer process exit; app file
inventory; launch, real MyAccounts sign-in, phone390 layout and zero page errors.
Screenshots reviewed. Explicit recovery-review mode remains stable while polling
is paused; interrupted UI offers both continuation and recovery options.
All Go packages/vet and embedded setup JavaScript syntax passed. No external DNS,
clean-host/reboot, actual remote-machine crash diagnosis or Actions claim.

Owner67/Slicer4 was not updated or used as a functional fixture by this task;
all24 services observed healthy afterward. Other authorized concurrent67 delivery
was retained. User on the other machine must run68, choose the same Workspace
folder and Check installation / continue setup. Older failed installs without
checkpoints use protected recovery; do not claim automatic restoration of their
unknown in-memory state. Keep existing owner data and latest rollback untouched.

Evidence docs/evidence/installer-continuation-20261022.json and installer-resume
screenshots. All named native services/browsers stopped; bounded synthetic roots
removed after process/reparse checks. Durable artifacts and delivery scripts on
I:/Projects/tooling/installer-resume-20261022. Published tags/assets are immutable;
the later documentation commit must not rebuild or replace68. No work remains.
