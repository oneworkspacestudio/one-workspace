# Portable Web Search service (in progress)

User requests building/testing/installing a One Workspace web-research API:
queued requests, real browser with persistent profile/session, small local model
controlling screenshots/input. Latest steering explicitly prohibits Docker and
removes the model-name constraint. No subagents authorized or used.

Actual current suite source worktree is I:/Projects/workspace-web-search,
codex/web-search-service, based on a78bc7b from workspace-ai-service. The user's
C workspace suite checkout and I:/Projects/one-workspace are older; preserve their
WIP. Read AI Service/remote AI Service handoffs. Owner installation is
C:/tmp/workspace/{Apps,Private}, port8990, public workspace.m8e.co.uk. Owner state
was inspected for paths only, never used as a functional test fixture.

New internal/websearch + cmd/websearch: native apphost OIDC/owner/CSRF; hashed
server credentials; SQLite bounded FIFO, per-grant isolation/cancellation,
logout/revocation checks; screenshot-only CDP actions (no DOM/JS); persistent
dedicated Chrome profile; public-HTTPS proxy with private/reserved IP filtering
and pinned DNS; local Fara1.5-4B Q4_K_M via portable llama.cpp Vulkan. Windows Job
Objects should own all descendants. This is a portable process runtime with the
browser renderer sandbox, not an OS-level container/VM. Profile shared only with
explicitly trusted owner-issued server credentials.

Model weights are private app model-library data, installed only through explicit
owner /api/model/install; immutable revision/size/SHA256 verification. Executable
runtime payloads ship in websearch-* catalog namespace. No ordinary-startup
download and no external daemon/system service/Python/WSL/Docker prerequisite.
Model binaries3,383,083,424bytes; source/download/cache/build artifacts stay on I:.

Suite surfaces added: managed registry, optional Administration category,
port9043/manager+53, runtime working dir, recovery/update scope, SSO/logout cookie,
native server-only ingress routes, branding icon/directory, complete build
component, installer integration13 (minimum120), publisher guard13. UI is AUI
compiled with pinned parser/React/esbuild; npm audit0. See docs/WEB_SEARCH.md.
Managed integration has NOT yet passed full acceptance. No public release or
owner installation/activation yet. Do not report delivery complete.

Artifacts I:/Projects/tooling/web-search-20261106; synthetic NTFS root
C:/Projects/.cache/web-search-20261106 only. Runtime downloads pin CfT157.0.8103.38
and llama.cpp b12167 Vulkan. Use direct Go1.27.3 at
I:/Projects/tooling/gomodcache/golang.org/toolchain@v0.0.1-go1.27.3.windows-amd64/bin/go.exe;
Go1.27.2 entrypoint invokes auto-downloaded1.27.3. GOCACHE go-build-cache and
GOMODCACHE gomodcache under tooling. GOTMPDIR should be the named C task root:
I scratch-directory metadata writes made tests/builds very slow. C free space
recovered to about7GB after stopping orphan task model processes.

Native testing exposed: Google returns CAPTCHA; use Bing default and direct
requested HTTPS URLs. Initial Fara output lacked valid action JSON; added strict
JSON-schema output and disabled thinking. Actual successful model-controlled
research is still pending. Repeated test used same profile and found leaked
process/profile locks; added graceful Browser.close, wait and native test job
ownership, plus unique fixture profiles. Need verify forced stop and profile
restart. Task-owned Chrome/model orphans were stopped by exact artifact executable
paths. No owner browser process was touched.

Docker Desktop was initially started before user prohibited Docker. No service
container was built or run. A base Python image was pulled; DockerCli -Shutdown
was issued. Need verify task-started Desktop is stopped and remove only that
unneeded task image if appropriate, without affecting other Docker state.
Docker source engine draft was removed from this worktree.

Resume: finish native model loop and security tests; fix any actual failures;
complete exact runtime licence texts and reproducible component packaging;
full suite tests/vet/build; synthetic supported fresh/upgrade/app-update/recovery
paths and desktop/mobile browser QA. Reconcile latest public suite/channel before
choosing an immutable release number. Publication must use mandatory signed
Android scripts/publish-release.ps1; Actions disabled. Only after passing supported
paths, install through normal backed-up owner update and verify served bytes,
model setup and health. Never edit live owner SQL to bypass login or use its data
as fixtures. Record each delivery stage separately.

## Continuation checkpoint (13:07 UTC)

Complete42-app candidate bundle built at tooling/web-search-20261106/suite120:
workspace-bundle.zip1,673,116,566bytes. Initial setup link failed C scratch disk
space; direct ordinary cmd/setup rebuild with GOTMPDIR on I succeeded, resulting
WorkspaceSetup-windows-amd64.exe1,687,039,488bytes. This is a CANDIDATE from earlier
source, not the final current source. Must rebuild after source freezes and tie
acceptance to exact final bytes. Public119 remains another task's draft;118 is
the last verified primary deployment. Do not mutate that draft.120 is intended
new immutable suite version; Android17/default scripts/publisher/check-org guard13
are prepared. Permanent signing loader can be reused privately from
I:/Projects/tooling/wonderpal-providers-20261106/publish119.ps1; it decrypts the
existing user DPAPI signing configuration. Never print signing values.

All websearch security unit tests plus catalog/suite/distribution tests passed.
Publisher tests4pass. npm audit0. Browser code must be readable by Chromium's
AppContainer; app grants read/execute ONLY to immutable browser code, and if
Programs is exFAT stages verified browser payloads into private NTFS
browser-code-cache (explicit appcontainer RX, no cookies/keys/profile sharing).
Native runtime C cache is C:/Projects/.cache/web-search-20261106/runtime; model
test weights C:/Projects/.cache/web-search-20261106/models. These3.4GB fixture
weights must be removed only after verified no task processes reference them,
before another model-bearing native install or owner deployment needs C space.
Original verified model artifacts remain on I, so these are disposable test copies.

Fara server loads C weights in9seconds. Tiny32KB Go reads of I weights took several
minutes, fixed with bulk4MB I/O. Initial Example Domain heading assertion was
STALE: the live page changed to a multilingual notice; screenshot proved the
model was accurately quoting its prominent English notice. Do not call this a
model failure. Native test default now verifies that actual notice.

More demanding live Python-version query exposed real integration/model issues:
Fara-native coordinates normalize1000x1000 and signed scroll positive=up, negative
=down. Adapter now correctly rescales, bounds scroll and constrains JSON numbers.
Full user instruction was an ineffective initial search string; start on Bing
homepage and let the model choose concise search keywords. Native model sometimes
anchored to an unsupported3.13.2 fact while still on Bing. Primary-source fact
provenance and finish gates were added; high-DPI1920x1200 screenshots improve
small text reading. Do NOT claim current-version factual acceptance yet.
Ground truth independently browsed Python official downloads:3.15.2,Oct25,2026.
Native model test is running in exec session86606 (same Python query); inspect
results. It logs only synthetic action text. Synthetic screenshot at
C:/Projects/.cache/web-search-20261106/synthetic-browser-evidence.png is safe to
view and revealed the actual live pages. No production screenshots/reasoning are
persisted; debug hooks are test-only.

Source now includes supported install-model CLI (explicit --data, optional offline
--source of pinned GGUFs), model reinstall repairs invalid existing files, owner
console blocks activation until model exists, auth loss clears private UI state.
UI esbuild uses automatic JSX runtime (needed for generated AUI imports).
Standalone native console UI QA is running via websearch-console-acceptance.cjs
with fresh console-r1 synthetic data; check desktop/mobile screenshots and result.
Compiled standalone WebSearch.exe at task artifact root is recent source; final
ordinary release still needs rebuild.

Reusable managed acceptance scripts added: websearch-managed-acceptance.cjs and
websearch-update-fixture.py. They exercise actual setup optional inventory, fresh
Accounts/Drive/WebSearch OIDC, explicit model install in supervised worker, real
local model+browser query through native remote API, credentials, restart, app-only
update/rollback, omission and118-to120 selection/old guard. Not run yet. They use
SEARCH_QA_* env paths; Playwright package I:/Projects/my-ai/web/package.json.
Ensure no standalone model test is running and delete its disposable C models
before fresh managed model download to avoid filling C. Program roots should
stay on I, test data under the one named C cache. Every fixture writes task PIDs.

All material source is uncommitted in dedicated worktree. No public release, push,
channel promotion, owner installation or activation has happened. Keep this
explicit. Continue native factual acceptance, final code review/tests/vet, ordinary
complete rebuild, managed acceptance, mandatory signed Android publisher/channel,
then protected owner install selecting workspace,websearch ONLY (retain MyAI25
and every unrelated app). Owner model can be prepared by supported install-model
CLI, leaving security-sensitive activation off until owner chooses Turn on.

## Source freeze and live acceptance

Live Python retrieval now PASSES: Fara navigates from Bing to official downloads,
scrolls/read pages and terminates; returned primary-source screenshot OCR contains
independently verified3.15.2. It is a RETRIEVAL API: result.sources[].text is actual
Windows OCR from browser screenshots. Pilot-generated factual summaries are not
returned as evidence (they proved unreliable). Calling LLM synthesizes the answer.
Full native Fara tool prompt/conversation flow is embedded under MIT, normalized
1000x1000 coordinates and signed scrolling; no DOM/JS extraction. Windows inbox
English OCR is available (en-GB/en-US), explicit prerequisite; helper embedded in
the executable, no external runtime download. Children receive only OS runtime
environment, not OIDC/publication secrets; model API key is child-only env, absent
from process arguments. Larger CDP wait respects job cancellation by closing its
connection. Queue remains bounded; results24h/1000, eight source pages, browser
disk cache128MB/media16MB. Service code is being frozen for final acceptance.

Native console QA passed actual1280x900 and390x844, real synthetic login, grant
creation/revocation, logout clears private view, no JS errors or overflow. Screens
under tooling/web-search-20261106/console-r1 visually reviewed. Added Windows EXE
icon from the shared SVG via pinned rsrc. Four package suites pass, all-package
Go tests/vet session89341 pending; final ordinary complete build session32201
pending (GOTMPDIR on I for large embedded payloads). Follow with exact native
managed installer acceptance after the final build. The earlier suite120 candidate
contains older service code; DO NOT publish/install that candidate.

Public119 finished in another task and primary is NOW119/MyAI26 (confirmed via
latest remote handoff and actual read-only config). Origin/main advanced only
channel/evidence/history to e1b79b0. Fetch done. Commit task changes and rebase on
origin/main before publication, preserving those three unrelated delivery files.
Owner install must select workspace,websearch only and retain MyAI26. Private
mandatory publisher wrapper is prepared at tooling/web-search-20261106/publish120.ps1
using the existing DPAPI signing loader; not executed. C fixture model copies
were removed after process-reference checks, originals retained verified on I;
this frees space for fresh managed first-use model setup. No public120 release,
channel promotion or owner Web Search installation has occurred.

## Managed acceptance completed

Source was rebased on public119 and pushed through c30e410. Browser keyboard
input now excludes address-bar, developer-tools, clipboard and OS shortcuts;
only page input, Ctrl+A and Shift+Tab are allowed. The ordinary complete120
build was repeated after this change. All-package Go tests and vet pass, plus
the new shortcut regression, four publisher tests and npm audit (zero findings).
The earlier all-package attempt used exFAT temporary data and overlapped a bundle
write: its filesystem/embed failures were corrected by isolated NTFS test temp
and sequential packaging; they were not suppressed or treated as passing.

Exact final120 native managed acceptance PASSED: optional setup inventory;
fresh Accounts/Drive/WebSearch; MyAccounts owner binding, shared picker and
anonymous denial; explicit supervised verified model download; real Fara/browser
retrieval through manager server API; credential ownership/revocation; restart;
independent app update and rollback; omission;119-to120 optional selection;
old integration12 controller rejects13 before mutation. Desktop/mobile screens
were visually reviewed. Evidence is docs/evidence/web-search/managed-acceptance.json.
Model pins and browser profile were checked after rollback, then disposable C
model fixture copies were deleted to preserve free space for later test phases.
All search-owned runtime payloads are additionally verified in both fresh and
upgraded installs. This cleanup changes no owner installation. Physical phone,
reboot, owner activation and owner inference remain unclaimed.

Frozen artifacts are tooling/web-search-20261106/final120. Private preparation,
mandatory signed-Android17 publisher, anonymous installer check and protected
primary120 deployment scripts are under that task artifact root. Publication,
public channel promotion and owner119-to120 installation remain the next steps;
do not report these as completed merely from this source/acceptance checkpoint.

## Public120 and primary installation completed

The previous pending-publication checkpoint is superseded. Mandatory local signed
Android publication completed for v0.3.0-preview.120, source7645119, Android17;
all890 immutable assets, anonymous APK and public channel verified. Channel source
promotion commit75f3ae6. The ordinary unauthenticated installer downloaded and
verified public120/WebSearch1 using its supported check-updates path.

Primary119-to120 protected update selected workspace,websearch only. WebSearch1
and all345 owned payloads plus manager match the accepted public release. Existing
app contracts/bytes, identities, credentials, AI Service settings/clients, paths
and recovery copies were retained. Transaction complete; normal cleanup receipt
checked. Installed served HTML/JS/CSS/branding match source, native/public private
APIs deny anonymous calls, and the exact remote route uses native key auth.

Supported install-model CLI imported and verified Fara1.5-4B Q4_K_M and vision
projector (3,383,083,424 bytes) into primary private models. Activation remains OFF,
with no owner API key or job created. Owner opens /websearch/, chooses Turn on,
then creates a named server credential. This is deliberate first-use policy,
not unfinished installation. No owner inference was used as a development fixture.

All synthetic runtimes stopped and private test data removed after evidence
retention. Docker is not a service dependency and no task Docker process remains.
No other PC, physical phone or reboot acceptance is claimed. Public summary is
docs/evidence/web-search/delivery.json; primary private receipts remain in the
task artifact root. Release assets are immutable; final delivery receipts belong
in source, not replacement release uploads. Source worktree is the dedicated
codex/web-search-service branch in I:/Projects/workspace-web-search.

## Remote browser maintenance in progress

User requested remote desktop into the service to resolve browser issues. The
implementation is an owner-only live view of the research browser, not the host
desktop. No Docker/RDP/VNC dependency. Source adds browser-only startup without
model/OCR, same-page handoff from the local pilot, queue pause/interruption,
mouse/drag/keyboard/paste/navigation/tabs/fullscreen, and exact native-session
90-second leases. Logout/role/expiry watchers cancel in-flight frame access;
server credentials cannot use desktop routes. Graceful browser close preserves
profile state. Reconnect retains forms; dead-browser restart is supported.

Unit tests cover CSRF, exact-cookie binding, credentials/roles, expiry and logout
during an in-flight frame, queue interruption/retention and input restrictions.
Actual portable Chrome standalone UI acceptance passed input, tabs, reconnect,
logout and desktop/mobile sizes. Visual review found a stale-tab close race,
fixed by server-side current-tab selection; repeat native acceptance is pending.
QA-only CDP/DOM reads are test oracles; production browser control still has no
DOM, selectors, scripts or host desktop tools. Fixtures are under the single
C:/Projects/.cache/web-search-20261106 root, never primary data. Artifact root
is I:/Projects/tooling/web-search-remote-20261106. The earlier test-root error
was a harness creating temp files inside a required-empty app data folder;
its layout was fixed. No production install was used for testing.

All-package tests/vet session54016 is pending. Next: final native UI and managed
fresh/120-to2 independent upgrade, real local-model handoff/queue resume, rollback,
then ordinary complete121 build/release with Android18 and protected WebSearch-only
primary update retaining manager120 and all unrelated apps. No121 publication,
channel promotion or primary WebSearch2 installation has occurred yet.
