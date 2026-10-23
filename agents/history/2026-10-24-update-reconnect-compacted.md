# Update-page tunnel reconnection

User screenshot shows75 /updates dumping Cloudflare1033 HTML while an operation
is running. User explicitly requests fix, push, release and install on BOTH this
PC (workspace.m8e.co.uk, C:/tmp/workspace/Private) and workspace-2 (different PC).
The remote Chrome Updates tab has an existing user-started multi-app operation;
never replay its POST or kill it. Remote tunnel currently unavailable. Need wait
for it to return or restore connectivity through authorized available access.

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
unknown caches or recovery archives. Publication and both installations pending.


Exact77 Windows manager acceptance passed: real native protected state denies
anonymous requests, accepts its synthetic capability, advertises revision2, and
its actually served page passes the browser outage/timeout/no-replay checks.
Initial fixture omitted its Windows owner marker and exited before serving;
fixed harness to call normal ProtectDir. No product workaround. Test process
stopped automatically. Native receipt records the exact executable digest.


Published immutable77 fromf7f74562f048e05e4cb3ee6bc35e987e7223f1e0. All public
asset sizes/SHA256 match. Promoting public update index; unauthenticated selection
and both owner deployments follow. Do not move the source release tag.
