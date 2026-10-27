# My AI native app gateway (in progress)

User requests My AI know and use all installed Workspace apps, including Hosting,
Git and Drive. Started clean 0a3689d, manager92/MyAI11 installed on both machines.
Uncommitted new myai_apps*.go plus portal/runtime wiring implement a current-
browser/account-scoped native operation registry at loopback integration endpoint.
Companion MyAI worktree my-ai has discovery, descriptions, reviewed calls and UI.

Gateway authenticates MyAI service capability AND live Workspace browser subject,
uses only cached native cookies from that browser, verifies destination issuer/
subject and native CSRF, and keeps native roles. No arbitrary HTTP proxy, minted
app sessions, credential disclosure, admin recovery adoption or public MCP app
authority. Apps without registered actions are discoverable/openable, not falsely
advertised as fully integrated. Special identity adapters: Drive OIDC username;
Notes identity endpoint; Tube read-only protected local ID mapping; Mail/DAV native
subject; Music OIDC extra fields plus native Music/JWT CSRF. Photos uses Drive API.

Git actions are scoped existing MyAI managed repositories. Hosting deployment
packages immutable pushed Git blobs, refuses links/credentials and oversized
source, requires native site authorization and reviewed exact revision; durable
per-user receipt reconciles uncertain upload outcomes instead of duplicating jobs.
No executable trust enabled. Native Hosting reports actual job state and URL.

Initial fake-protocol tests passed for browser isolation/expiry/identity/CSRF,
unregistered parameters/routes, human-review requirement, redaction and digest
stability. Native APIs/examples still being verified. Real fixture deployment,
Drive/Git roundtrip/model workflow, full regression/vet and final review remain.
No version bumps, release/channel promotion/install for this feature yet. Do not
claim all app operations tested or deploy over immutable existing releases.
Artifacts I:/Projects/tooling/my-ai-workspace-apps-20261026; synthetic NTFS private
root only C:/Projects/.cache/my-ai-20261024. No new fixture started yet. Preserve
all owner state/normal backup and coordinate existing authorized peer before any
shared channel/apply work. See previous project-picker handoff for delivered11.

## Updated resume point

Fast-forwarded to0926501 (Tasks4 promotion); all external pins retained. New
myai_apps_schema.go validates registered write shapes before review, app discovery
is required by MyAI, and empty read bodies normalize to no HTTP body. New
myai_apps_deployments.go separates native upload receipts from immutable Git
archive preparation: server-derived per-turn action IDs distinguish retries from
later redeployments; pending source pointers survive unconfirmed responses;
history reconciliation requires matching source/environment AND creation time;
aliases prevent duplicate retries; definite native rejections remain distinct
from ambiguous409/5xx/network outcomes. Native Hosting1 itself is unchanged.
New tests prove accepted/repeated/later deployments, lost responses across action
IDs, old-history refusal and later retry of definite rejection. Targeted suite
tests all pass. Native Tube identity SQLite URL test passes on Windows.

Full suite exposed11 pre-existing failures; immutable baseline0926501 archive
reproduced exactly those failures. Test-only fixes in ecosystem/mailboxes/seamless,
update/reinstall/update_cleanup fixtures now use current role probes and91/92
metadata recovery. Old-manager rejection remains explicitly tested. Current user
libraries and earlier full-data recovery copy remain protected. Full suite passed
after fixture corrections; rerun after latest connector source additions. No
production update protection or Mail implementation changed.

Real fixture p12 exists: C:/Projects/.cache/my-ai-20261024/p12 private NTFS,
programs under I:/Projects/tooling/my-ai-workspace-apps-20261026/fixture-programs.
Portal39600, isolated Caddy39680 provisioned normally from pinned official ZIP.
Seven native services run; actual SSO, discovery and app connections checked.
No successful real AI mutation/deployment yet. Three qwen4b attempts exposed
contract/read-body/model-argument issues, now fixed in source/tests; none created
a task. Latest dev candidate build is running. Follow MyAI companion handoff's
resume steps, current PID receipt, and do not use owner installations as fixtures.
No feature version bump/commit/publication/channel/install yet; versions11/92.

## Native acceptance update

The latest qwen3.5:4b fixture completed a reviewed Tasks creation, independently
visible in native Tasks Inbox, and a reviewed Notes creation, independently
visible in native Notes. Exactly one of each synthetic acceptance item exists.
Full MyAI Go tests including bundled search and the suite's full Go tests passed.
Drive revealed an actual trailing-slash registry route bug: connector constructed
/drive//api/users and treated native307 as a connection issue. Normalize route
prefixes for both native calls and connection probes; added a direct HTTP regression
for routes with and without trailing slash and already-prefixed/native paths.
Rebuilding candidates/targeted suite tests for this correction. Native Drive and
Git-to-Hosting deployment acceptance still unfinished; no publication/install yet.


## Completed native acceptance and release preparation

Real qwen3.5:4b Code created General automatically, wrote index.html and pushed
it to private native Git. Registry HEAD and remote Git file bytes match Drive.
My AI then reviewed/created a static Hosting site, reviewed the exact pushed
revision, submitted it and polled its actual job to published. The returned local
preview page was opened and rendered the expected heading. Native Tasks, Notes
and Drive mutations are independently visible in their apps. Acceptance receipt
is docs/workspace-apps-acceptance.json in MyAI and artifact native-acceptance.json.
All tests used only the named synthetic fixture; no live owner app fixtures.

Installed app names now route from the actual catalogue and preload their native
action descriptions, covering future apps without a fixed one-scenario list.
A small local model still made invalid/irrelevant calls before recovering; do not
claim perfect model instruction-following or certify every registered app. Native
body validation, permissions and reviews remain authoritative. Mixed public web
research/private app saves should use separate turns; legacy full Codex and cloud
providers are not live-certified. Mail is read-only, unsupported apps open-only.

Final full Go tests in both repos passed, including bundled SearXNG runtime16.1s.
Both go vet ./... passed; AUI/TypeScript/Vite passed; npm production audit0. Go
vulnerability scan found no reachable vulnerable call paths (module advisory
GO-2026-5932 is reported for the existing crypto dependency; do not call it zero
advisories). Cleanup's bytes-freed assertion now uses actual stale-stage bytes:
hard-linked backup removal can correctly free0. Production updater unchanged.

Final versions prepared: MyAI12 requires manager93 in compatibility and managed
metadata. Candidate binaries built. Final fixture restart/version/persistence
check, clean source commits, immutable packaging, publication/channel promotion,
both protected owner installs and fixture cleanup remain. Preserve current main
b76ccf3 Maps5/Tasks4 promotion and unrelated pins. Authorized coordination peer
reports no overlapping update or manager93 reservation; secondary needs a fresh
transaction check before applying.


## Final privacy correction / delivery versions

Final review caught raw prior-user fallback in public research: a short follow-up
could refer to a dynamically recognised app request not matched by fixed filters.
When private Workspace context exists, the researcher now receives no raw message
history; only its separately verified public ResearchContext may resolve a follow-up.
A regression ensures a private app request followed by "tell me more" never starts
public research. Full app tests including bundled runtime13.6s and both vets passed.

App12 source e96a8fb and immutable public assets were already published, but never
promoted or installed. Explicitly withhold12; do not alter/delete its tag or assets.
Corrected app13 requires manager93. Manager93 remains unpublished/uncommitted.
Core native acceptance is unchanged; final13 candidate health and persisted data
verified after synthetic restart. Phone viewport override did not change the Chrome
extension viewport (still1794px), so mobile-specific rendering is NOT certified.
Desktop app cards render correctly. Continue app13 packaging/publication, suite93
source/package, channel generation (exclude12), protected installs and cleanup.


## Final14 direct native tools accepted

The local4B model's occasional invented dispatcher/app names motivated direct
workspace_APP_ACTION functions with flat typed fields, generated from registered
native schemas for relevant/discovered apps. They adapt to exactly the same
registered gateway, native permissions, privacy flags and conversation reviews.
Only current browser scopes receive them;48-action cap retains provider headroom.
Unit tests cover typed fields, path/query/body separation, unknown authority and
absence in public scopes. Full app tests plus final targeted cases and both vets
passed. Frontend assets unchanged from validated AUI build.

Real final14 candidate read the Drive acceptance file correctly in one call,
read Git status, reviewed/submitted a new intentional preview of the same pushed
commit, and polled the native job to published. Opened its actual returned page
and verified the expected heading. It recovered from some invalid native arguments;
do not claim perfect4B behaviour. Native permissions/reviews held throughout.
12/13 remain immutable, never installed and explicitly withheld;14 is the final
promoted app target with manager93. Preserve main7663857 Maps6/Tasks4 changes.
Normal compatibility checks require manager93 first, then app14, each through the
protected updater. Neither owner installation has started yet. SDK/installer
updates and cleanup remain; no further feature expansion needed before delivery.


## Delivered and verified — 2026-10-27

Source changes are pushed to both main branches. Immutable releases:
- MyAI14 / 5b9ec7f14a3c55c2d9186afdb70f407ced039921;
  MyAI.exe817e3c373901b7485ae69f7e10ed42cb9be4ca4ba52a4b0556d54bc2d03400ba.
- Workspace93 / fb07cdd06d96cb5474c76e117b2209a5ac1d8ba4;
  Workspace.exe7456db79f1cf330b36d5d324a44845adf0dadbbde3875e427d95fe3ac0921409.
- Public channel promotion a974397 retains prior entries/other pins and excludes
  immutable superseded app12/13. Both executable downloads and default raw channel
  were verified anonymously. Actions remained disabled. No tags/assets replaced.

Both PCs are installed through normal protected updates: manager93 first, then
MyAI14 (current-version compatibility requires those two stages). Both report all27
services running; actual MyAI health says14. Other app versions including Maps6
and Tasks4 retained. Primary config/identity preserved; all executable bytes match
installed contracts and exact tested manager/MyAI hashes; all18 prior models kept.
Its initial audit stopped on unequal DPAPI ciphertext, after both updates already
completed. Do NOT rerun install-primary.py. Config.Save re-encrypts secrets on save;
audit-primary.py instead compared decrypted settings ONLY IN MEMORY against both
pre-app-update recovery and an older full backup: identical. Protected provider
files and app-owned private SearXNG/Python3.13.18 were verified. No credentials
printed/stored in artifacts. Corrected sanitized receipt is primary-delivery.json.

Primary recovery update-20261026T235839-f2e2d9a024f3; secondary recovery
update-20261027T000935-58884a42ac1f. Secondary updater independently reported all
backup/check/install/restart/cleanup stages complete; public health27running and
MyAI14 checked in Chrome; table versions compared before/after. Its binary hashes
were checked by the normal updater, not a separate remote filesystem audit.
No update remains active. No live owner account/Drive functional fixture used.

Named synthetic p12 and only the named inactive Go temporary directory were
stopped, checked for process references/reparse points and removed. Six test tabs
closed; temporary viewport reset. Artifacts/provenance/screenshots and source stay
on I:/Projects/tooling/my-ai-workspace-apps-20261026. User's second-PC updater tab
is left open as installation evidence. No live library, model, data or manually
selected recovery backup deleted; normal updater cleanup kept usable recovery.

Acceptance and delivery evidence in MyAI docs/workspace-apps-acceptance.json.
Supported actions retain native permissions/reviews. Unregistered apps open-only,
Mail reads only, no Vault/admin/physical printing authority. Every connector,
cloud provider, separate full Codex transport and clean-host/service/reboot are
NOT individually certified. Mixed private app/public research uses separate turns.
The delivered native Chat/Code core workflow is complete; no feature work pending.

Post-install UI check: opened the second PC MyAI through the normal Workspace
SSO entry, then its Tools & skills dialog. The actual installed app catalogue
and native action counts render. This was metadata/navigation verification only;
no live owner record or file was used as a functional fixture. The fresh MyAI
Tools tab and updater are left open for the user; no test or health tab remains.
