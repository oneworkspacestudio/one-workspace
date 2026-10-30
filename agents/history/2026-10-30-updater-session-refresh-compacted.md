# Stale updater after manager replacement

User reported missing app rows/new apps after installing Workspace106, then
requested a product fix so other users do not encounter the same problem.

## Observed cause

The secondary owner installation was106 with31 installed app registrations and30
healthy supervised services. Its maintenance process remained the older updater
from103. The marker reused remote protocol4, so the portal accepted that process
despite its different executable. The local opener similarly reused any live
window. The old loader rejected the new catalog and fell back to a verified103
backup configuration. Its check returned installer106 required and choices=null;
the frontend cleared the table without rendering an explanation.

Closed only the confirmed idle controller through its authenticated close API,
reopened through the installed manager and checked releases. Controller106 then
returned32 groups, including Remote5, Tickets1 and Finance1 available to install.
No app update/reinstall was submitted as part of that diagnostic repair. All30
existing services remained healthy. No credentials, task data or browser sessions
were used as test fixtures. Do not copy capability-bearing updater logs to evidence.

## Source fix, not yet published

Branch codex/updater-session-refresh, based on1bb6c67. A marker records SHA256 of
the actual running updater executable; reuse requires matching installed manager
package digest as well as remote protocol. Both local and remote launch paths
retire mismatched idle controllers through authenticated state/close calls and
wait for exit before replacement. Busy, ambiguous, unresponsive or rejected-close
controllers are preserved. Close still rechecks busy under the mutation gate.

The state endpoint recognises a changed installed binary using a minimal stable
header before attempting to load its catalog or a backup configuration. Thus an
unreadable future catalog reports the installed version and a handoff, not the
older backup version. Stale controllers reject further update mutations. Idle
local pages request the installed manager once to open the replacement page;
remote portal requests automatically reconnect via the manager. Busy work waits.
An empty release plan now renders an explicit app-list/installer explanation.

## Validation

Windows Go1.27.3, GOWORK=off, cached pinned module dependencies. PASS focused
suite/maintenance tests: TestUpdateController*, TestMaintenanceCapability*,
TestReinstallRequires*, TestUpgradedController*. Tests use synthetic temp headers
and HTTP servers, covering same-protocol different binaries, legacy markers,
future catalog schema, busy/missing/malformed status, close races and stale
mutation rejection. PASS go vet ./internal/suite ./internal/maintenance.
PASS node scripts/update-controller-handoff-test.cjs: actual refresh function
opens the replacement once locally; busy and remote refreshes do not launch it.
PASS git diff --check. No native end-to-end upgrade, browser rendering, reboot,
new installer acceptance or publication is claimed for this source patch.

Build cache is I:/Projects/tooling/updater-session-20261030/gocache. Named NTFS
test temp C:/Projects/.cache/updater-session-20261030 was removed only after it
was empty. Diagnostic helpers remain in tooling and contain no owner secrets.

Rollout needs a new immutable manager/setup release through the repository's
signed-Android publisher and normal protected updater, followed by native
old-to-new controller acceptance. Do not replace published106 assets or claim
this patch is installed on owner systems. The current owner repair only restarted
the existing106 updater. Existing public106 bytes remain unchanged.
