# Update-page tunnel reconnection

User screenshot shows75 /updates dumping Cloudflare1033 HTML while an operation
is running. User explicitly requests fix, push, release and install on BOTH this
PC (workspace.m8e.co.uk, C:/tmp/workspace/Private) and workspace-2 (different PC).
At task start the remote Chrome Updates tab had an existing user-started
multi-app operation. Waited for that and later Workspace76/Mail24 operations
before starting the separate Workspace77-only update. Never replayed their POSTs.

New source clone at I:/Projects/tooling/update-reconnect-20261023/source, branch
codex/update-reconnect, basefe17dd7 containing public76 sourcea66d81d and channel
promotion. Primary I:/Projects/one-workspace is stale56 with uncommitted histories;
leave it untouched. No agents or Actions.

request_page.go supplies bounded API requests, generic connection errors for
network/timeouts/5xx/429, safe bounded plain/JSON error messages, validates state,
and never surfaces proxy HTML. Dedicated connection status retains last progress.
Controls start disabled and remain disabled for disconnect/in-flight action;
action generation drops stale polls. Unknown POST acknowledgement only polls
state, never retries a mutation. Real409 errors remain across successful polls.
Remote controller revision2 retires only idle older helper via existing busy
checks, ensuring refreshed pages receive new UI without killing active updates.

Browser acceptance scripts/update-reconnect-acceptance.mjs passes synthetic502,
200HTML, malformed JSON state, request timeout/recovery, lost POST/no replay,
real409 persistence and desktop/phone overflow checks. Exact packaged native
page/capability acceptance is opt-in TestNativeReconnectPage, fresh NTFS task
root only. Go vet passes; full tests show five Mail-role probe fixture failures
also reproduced on clean76 baseline. Maintenance tests and relevant update tests
pass; do not claim a passing whole suite. No live owner state as fixtures.

Target77 manager/full installer from verified76 bundle; change only manager and
release-status notice, retain all26 app contracts and all other payload bytes.
Artifacts I:/Projects/tooling/update-reconnect-20261023; only disposable C scratch
C:/Projects/.cache/update-reconnect-20261023. Preserve source/evidence and stop
tracked fixture PIDs before bounded no-reparse cleanup. Never touch owner files,
unknown caches or recovery archives. Publication and both installations completed
as recorded below; task-owned C scratch removed after process/path checks.


Exact77 Windows manager acceptance passed: real native protected state denies
anonymous requests, accepts its synthetic capability, advertises revision2, and
its actually served page passes the browser outage/timeout/no-replay checks.
Initial fixture omitted its Windows owner marker and exited before serving;
fixed harness to call normal ProtectDir. No product workaround. Test process
stopped automatically. Native receipt records the exact executable digest.


Published immutable77 fromf7f74562f048e05e4cb3ee6bc35e987e7223f1e0. All24 public
asset sizes/SHA256 match. Public index promotiondc0c61d preserves concurrent
Mail24 promotion17c5a44; unauthenticated default-channel selection verified after
the raw CDN's300-second cache expired. Do not move the source release tag.

## Completed delivery

Release: https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.77
Exact tested manager SHA256:
5608ea0a594aa63346a3a3a991b644df568d4cb7551eac4dcde1b1f186142b1b.
Public manager downloaded and digest verified. No Actions used.

Both PCs updated76 ->77 through the normal verified-backup updater. Local Go
receipt confirms all app contracts, other package bytes, app selection, owner
configuration and decrypted identity/credentials retained; transaction complete,
both manager copies match the tested binary, newest recovery copy retained,
all25 services running. Actual owner Updates page runs protocol2 helper with
matching digest and new reconnect status; normal release check succeeds.

Remote normal UI selected only Workspace, accepted one install request, downloaded
the verified package, then stopped its tunnel during backup/restart. Returned77
after roughly15 minutes with cleanup complete and newest76 recovery copy kept.
Refreshed page has the new reconnect element, empty error and successful release
check; before/after table comparison changes only Workspace. Remote /healthz
reports all25 services running. Remote binary was validated by the normal update
flow, not manually hashed from this PC. Existing owner session remains usable.

Sanitized publication/local/remote receipts are under docs/evidence/ with
update-reconnect-*20261024.json names. Actual installed-page proof is
docs/screenshots/update-installed77-{local,remote}.png; synthetic outage images
remain separate. No private capability, cookie, secret or owner content included.
Durable build/runtime evidence remains I:/Projects/tooling/update-reconnect-20261023.
No task fixture processes remain. C:/Projects/.cache/update-reconnect-20261023
removed; normal owner runtimes/helpers intentionally remain running. No remaining
delivery work for this fix; five pre-existing Mail fixture failures remain as above.
Concurrent main34682f3 adds portal sign-in recovery after77; preserve it when
pushing this final documentation-only delivery record. Immutable77 stays atf7f7456.
