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

## Release107 and installation completed

Later human request explicitly authorised a new release and installation. Android
signing is on another PC; the human explicitly instructed: "release this version
without android". This exception is scoped to107. The repository's permanent
Android signing policy and supported publisher were not weakened. A task-local
107-only adaptation retained source/clean-tree, Actions-disabled, catalog,
immutable draft, full asset digest, anonymous manager and public channel guards.
Release notes and docs/UPDATER-REFRESH.md explicitly describe Windows-only scope.

Frozen release source f780d09eb32b59ae7b6c2e357055dc5444970b3f. Public tag
v0.3.0-preview.107; all325 assets (4,940,468,846bytes) match local sizes/SHA256.
Anonymous public manager download and exact public raw-main channel verified;
promotion67e1eab. No published assets or tags were replaced. Actions stayed disabled.
Manager SHA25601d84f0a4736ceefb0d4c428631db617abd88a7772c9002d398b7ed54e1f2716.

Normal managed build retained34 apps/307 packages from published component pins.
Native acceptance discovered a fixture copy-over-existing mistake, corrected by
using atomic replacement. Actual Windows image locks also required a bounded
access-error retry after confirmed idle exit. That runtime fix is in the frozen
source. Final native public106-to107 helper retirement, matching identity/state
and repeat-open reuse PASS on synthetic NTFS data. No owner fixtures or app
services/tunnels were created by this test; its helper processes were stopped.
All suite/maintenance/distribution/catalog tests and vet PASS. JavaScript idle,
busy, remote and single-reopen checks PASS. Exact final verified bundle was found
byte-for-byte inside the final installer. Public acceptance receipts accompany
the release. No phone, reboot, clean-host or service certification is claimed.

The secondary owner at I:/Workspace/Private had already reached106 with34 apps
before this apply; earlier31-app diagnostics above are superseded. Anonymous
normal plan selected only Workspace107, downloading17,536,512bytes. Protected
106-to107 apply completed. Recovery1879.7MB checked,1324.4MB reused,555.3MB newly
copied; user libraries remained in place. Normal cleanup retired1 older eligible
automatic copy (~240.2MB), retaining the newest verified recovery:
update-20261030T155040-f589b36e3748.

Independent read-only before/after fingerprints confirm existing identity,
issuer/origin, storage, ports, connection/autostart settings, selections and all34
app contracts retained. All276 retained nonmanager/nonrelease-note package files
match their prior SHA256. All33 services report running; journal phase complete.
Opening Updates through the installed107 manager automatically retired the older
helper. Installed manager, maintenance executable and active marker all match the
public107 SHA256. Fresh anonymous release check returns32 groups, no errors, no
pending operation and no stale-controller flag. The update page is left open.
Finance2 is independently available; this manager-only apply preserved Finance1.

Evidence/artifacts: I:/Projects/tooling/updater-session-20261030, including
release107/{build-info.json,native-controller-acceptance.json,
installer-bundle-verification.json}, owner-delivery107.json and owner-install107.log.
Do not publish private updater capability logs or the owner config snapshot.
Only this PC was installed. Other-PC rollout was not requested or claimed.
Large public artifacts and build cache remain on I:. Named C NTFS temp root was
removed only when empty after native helper/build/publisher processes ended.
