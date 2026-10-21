# App startup recovery / My Slicer

User reports Windows paging-file error launching MySlicer and requests continued
installation/startup plus retry of failed apps. Screenshot is startup after the
verified package installation. Preserve atomic package/update integrity and all
app data/versions; normal backed-up local delivery and GitHub publication remain
authorized from this conversation. No Actions or live owner data test fixtures.

## Findings and implementation

Error is Win32 process-creation commitment exhaustion, before Slicer code runs.
Installed Slicer EXE is about16MB; engine is separately installed and on-demand.
Owner Slicer is currently healthy; sampled32GB RAM/~20GB pagefile has ample free
commit. Historical resource exhaustion is not reproducible from current readings;
no evidence of corrupt Slicer or eager engine launch, no pagefile settings changed.

Start failures previously aborted the manager and killed every other app. New
startAvailable continues independent apps, marks failed/dependency-blocked steps,
and retains Accounts as essential sign-in. Windows memory creation errors receive
two bounded delays/retries; persistent errors get plain close-unused-apps/retry
guidance, raw diagnostic only protected app log. Each app has its own kill-on-close
job, cleanup before retry, active-process guard and normal supervisor. Restart
health is verified. Retries serialize with backup/shutdown and preserve healthy
processes. Shared preparation failures can be retried. Backup restart progress is
refreshed. Also normalizes Windows path spelling before temp ancestor walk (a
forward-slash synthetic fixture exposed the old infinite ancestor-loop edge case).

Setup keeps Open Workspace available with optional failures, provides Retry apps
that need attention and shows retry progress. Owner App status and retry page is
linked from hub and app-unavailable guidance. Owner+origin checks protect browser
retry; setup capability+origin protect local wizard; loopback and encrypted local
control capability protect setup-to-manager retry. Failed/blocked checks are
terminal attempts, labelled steps checked; failures remain visible, not success.

## Verification and delivery state

All Go packages passed; vet and embedded JS syntax checks passed. New regressions:
wrapped Win321455 transient/persistent/non-memory/cancel; independent continuation,
dependency blocking/recovery, Accounts fatality, permissions/origin/maintenance
lock and progress error reset. Opt-in real Windows child fixture passed: failed
Slicer then independent process, recovery without healthy PID changes, repeat
retry and timeout termination of primary+real descendant. Fresh named NTFS only.
Initial native test stalled on preexisting slash-normalization bug; scoped fixture
was stopped, normalization fixed, repeated native test passed in about4seconds.

Next: build62 from committed source and verified public61 bundle, preserve all25
app contracts/payloads. Real native installer+MyAccounts browser partial-start and
retry verification at desktop/mobile, publish immutable new62 with all asset
hashes, normal protected owner61->62 manager-only update, verify data contracts,
latest rollback and all services, cleanup named task cache. Artifacts/logs on
I:/Projects/tooling/app-startup-20261021; only C:/Projects/.cache/app-startup-20261021
is disposable after task-process/reparse checks. Preserve owner/live data and
latest verified recovery. No physical printing, external mail or live DNS tests.

## Shortcut recovery follow-up before publication

Final review found built-in Docs could retain a blocked status after Drive
recovered. Built-ins now derive dependency readiness without fake process entries;
a targeted regression proves recovery through a built-in dependency. Keep the
public62 source tag immutable; its release420043453 remains an unpublished draft
with20 small assets (large uploads stopped). Prepare final63 from new source.
Normal owner61->62 backup/update was already in progress and must finish safely;
then apply63 normally. Do not interrupt or tamper with an active transaction.
Final source/release/local claims must refer to63 and reconcile both receipts.

## Final verified delivery (supersedes pending next steps)

Source b018b00 and immutable v0.3.0-preview.63 are public, all 22 release asset
sizes/digests verified. Normal protected owner 61->62->63 manager-only update passed:
all 25 app contracts/packages retained including Mail 17 and newer installed Chat 9,
identity/settings/paths/selections unchanged, all 24 service entries running.
Newest recovery manifest verified. Exact tested local artifacts were applied
while large GitHub uploads continued; public-download status is claimed only
after final publication verification. No source tag or published asset replaced.

Actual final installer installed fresh Accounts/Drive/Docs/Slicer/Tasks. Deliberately
unavailable Slicer retained a failed row while Tasks/portal stayed ready. Real
MyAccounts login, owner retry, setup retry with Open Workspace still available,
Slicer recovery, unchanged healthy PIDs, phone/desktop screenshots and zero browser
errors passed. First browser run caught malformed punctuation encoding in setup;
it was corrected and installer rebuilt before the 62 tag. Manager and
all app payloads remained byte-identical across this setup-only correction.
Final 63 includes the separate tested built-in shortcut dependency fix; 62 draft/tag remain preserved and unpublished.

Synthetic managers, descendants and browsers stopped; exact named C task cache
removed after no active-reference/reparse checks. Current owner installation and
newest recovery retained. Installer 63 SHA-verified in Downloads. Evidence:
docs/evidence/app-startup-recovery-20261021.json and startup-retry screenshots.
Historical real Windows memory pressure cannot be recreated from current free
memory; bounded retry/degraded operation tested via Win32 1455 injection and real
process failures, with no OS pagefile tuning or physical printing. Native owner
Slicer currently healthy. No clean-host/reboot or Internet-mail delivery claim.

Final public portal health confirms all 24 services healthy; Actions remains
disabled. Removed 5,828,934,177 bytes / 1,103 known synthetic cache files after
process/reparse checks. Preserved owner63 and its newest verified recovery.
