# Chat parity — active integration, not published

The October19 brief authorizes upgrading existing Chat, Workspace integration,
local Windows builds, push/releases, no Actions. All consumer families remain
mandatory; no full parity or production E2EE claim. No subagents authorized.
Use repository-history first. Detailed state and exact resume point are in Chat
agents/history/2026-10-19-whatsapp-parity-compacted.md; read that before continuing.

Baselines: Workspace0db607a/main/public42 (source5ebc588); Chatc5c48b02/master/
public4 (source32d9bf3). Both fetched October20 and0/0 remote divergence. Current
branches codex/chat-parity and codex/whatsapp-parity. Changes uncommitted; no new
push/tag/release/owner update yet. Preserve unrelated untracked installed-history
files. Owner42 with20 apps is untouched, never used as a fixture.

Workspace implementation: dedicated Chat hostname (local other loopback host;
managed CF first-level hostname/owned route preservation and tunnel lock), saved
address continuity, exact existing OIDC client callback update without credential
rotation, legacy/chat redirects and OAuth restart without transplanting codes.
Private per-user app catalog with owner-only monitor and absolute destinations;
Chat bundles actual icons and AUI switcher, no wildcard CORS/shared cookies/scripts.
Global logout is private verified issuer/subject capability backchannel, not a
Chat cookie read from portal's wrong hostname. Open sockets cancelled. Actual
synthetic browser test passed: Alice logged out, Bob stays in; client heartbeat
catches401 before a rejected WS can reach onConnect, no silent provider re-login.
Managed public CF/DNS operation is implemented but not yet verified remotely.

New distribution app contract schema2/integration3 with minimumManager and caps,
strict old-schema rejection, checks at decode and SelectUpdates, CLI flags and
unit tests. Recognized caps currently app-update-requirements-v1/chat-state-v2;
do not advertise origin/E2EE capabilities before final integration proof.
Important bypass: old42 can still choose a Chat5 package from a new full-suite
catalog even when standalone manifest2 is rejected. New Chat managedAdmission
runs before directories/secrets/extraction/DB, requires manager-supplied
isolated-chat-origin-v1 plus home/catalog endpoints. Runtime command now injects
that capability outside App.Env. Never place manager capability in catalog Env.
Actual native executable rejected an old launch shape without creating/opening
its sentinel data directory. Still need actual published42 manager offered new
manifest, full-bundle app-only rollback before migration, and final paired update.

Full suite tests passed using a compiled suite.test.exe built on I and executed
with GOTMPDIR named C/go-tests. Go1.27 t.TempDir uses GOTMPDIR; changing TEMP only
while GOTMPDIR=I produces expected ACL/exFAT failures. Quote '-test.timeout=180s'
in PowerShell. New Chat seamless test asserts no shared-cookie probe, not a skip.
Distribution and update-manifest tests passed. Re-run after final changes.

Chat source has real atomic UUID receipts/message+links+sequence, gateway-bound
WS/LP, media membership and anti-forged-link checks, scoped caches/drafts,
static-only worker, durable offline replay and ack clearing, real custom lists,
layout/token changes. Native packaged synthetic tests passed restart/IDs/bytes,
role denial/removal, changed-body/retry/delete semantics, list CRUD across sessions,
actual group UI creation/invitation/messages, five viewports in two themes,
offline reload/replay once, global logout and no uncaught errors. Latest source
adds neutral header/plain default surface/product empty state, managed admission
and release clean-source/matching-client build receipts. Final prototype re-run
is currently underway; see tooling foundation-*-final.log. SourceDirty dev build
is not a released artifact. No public5/43 tags or org component pin changes yet.

OpenMLS native/browser bidirectional proof passed tamper/replay/AAD and enforced
recovered-snapshot send fencing. Not wired to Chat. Default E2EE, signed device
trust/recovery, encrypted media/search/outbox, social families, SFU/TURN E2EE calls,
native clients/push/fresh-machine restore remain mandatory and incomplete. Native
signing/Mac/device question is unanswered; do not use it to omit browser work.

Tooling I:/Projects/tooling/chat-parity-20261019. Named C:/Projects/.cache/chat-
parity-20261019 holds needed NTFS client/Cargo/fixtures only. Protected fixture/
Private plus synthetic credentials outside repo; Apps/Files on I. Baseline public42
accounts+chat only intentionally excludes other apps. Stopped-state baseline42-
cold/Private preserved for real upgrade, do not re-seed/reset identities. Prototype
manual binary/hash replacements are explicitly NOT actual updater acceptance.
Originlocalhost30960, issuer30961, Chat127.0.0.1:30976/chat; fixture-process.json
tracks PID. Helper stops/validates exact fixture, starts hidden, polls127 health.
A running exec session can remain attached to its child after health is ready;
read health/progress rather than waiting on it until fixture shutdown. C free~20GB,
I~3.6TB. End by stopping/verifying and scoped cleanup after preserving source/
evidence; never remove owner or newest rollback/unrelated roots.

Next: finish final test/capture/source/docs/history, stage only intended changes,
publish truthful immutable candidates without Actions, pin released Chat ZIP/hash,
use build-org.ps1/build.ps1/update-manifest tooling, test actual downloads on
fresh/old fixtures and all update gates. Planned app5/Workspace43 only. If a
foundation increment is published, clearly keep the brief's remaining work active;
never mark all parity done from those releases. Continue encryption/rich/social/
calls/native/recovery work afterwards. No new owner data as fixtures.

## October20 correction and release preparation

Supersedes the earlier capability/resume wording: actual local native SSO,
private catalog, global logout, managed account controls and all final suite/
distribution tests passed. Manager source advertises isolated-chat-origin-v1
and injects the exact dedicated callback plus capability outside catalog Env.
Chat managed admission rejects an old manager before creating storage. Actual
old public42 offered-release/bundle rollback remains required. Chat source
commits37102177/4f0fe08e are local; final clean release build is running.
No GitHub release, new component pin or owner update yet. Full brief remains
incomplete, especially application E2EE, social/calls/native clients and recovery.

## Verified old-manager package gate

The actual public42 installer applied a complete compatibility bundle containing
the exact final Chat5 executable115d4f276129c5b58b843401bdca853072c4a84c8f2499718e7105228b8a7ea8,
selecting only Chat. Chat exited with minimum-manager guidance before runtime
extraction (new runtime hash directory absent). The real updater transaction
rolled back automatically. The restored executable matched publicChat4 SHA,
Workspace42/Chat4 versions returned, and synthetic immutable IDs, retained message,
attachment link, browser download filename and exact bytes passed after rollback.
This is not yet the final public43 bundle/download test. Owner data untouched.

Clean Chat build source21fe7325, release5 ZIP SHA
db72538e9e722da41a2d46a74b227b4ed2451e3605312e0f8cd402148d87a933.
All final packaged native/baseline/UI runs passed, including seven viewports in
both themes. Chat master source is pushed; release5 draft upload in progress.
Workspace pin now targets those exact bytes. Final build-org build, public
manifest rejection, downloads/fresh/install/paired upgrade and rollback remain.
Legacy Chat redirects now remove old media transport credentials from the query;
full suite tests passed after this change. Do not claim app E2EE/full parity.

## October20 actual public43 downloads and follow-up

Public43 is published from3c5e951,27 assets all size/SHA-verified. Installer
3b0099182371b59fd3d9afe30786b11807bd5d5b4a1a41b316e55e2b844204ce and bundle
1be166d1a024d53c05c71762aadb4b37dc7c15ca7fd56b4e5b3ac496e7d4543c downloaded.
Actual public42→43 update preserved IDs/history/UI attachment filename and bytes,
selected Accounts+Chat and18 opt-outs. Default-online fresh install passed native
startup plus baseline/security/UI after asserting directory snapshot ready. First
immediate new-account Find failed before synchronization; Chat6 now tests automatic
refresh with a real new account created while Find stays open. Keep this boundary
explicit. Neither public43 nor5 is full parity or application E2EE.

Uncommitted dedicated_tunnel.go helper now prepends only the owned exact route
before wildcards, repairs its Host header preserving options, refuses conflicting
routes and skips occupied candidates. Tests prove input/unrelated-rule preservation,
idempotence and conflicts; full suite PASS (manager-routing-suite.log). No public44
or actual Cloudflare operation yet; v1 origin capability unchanged. Provider callback
fix and Chat6 candidate coordinated in companion histories. Owner still42/healthy,
no owner account/data used as fixture, normal released local update still pending.

Named C task's redundant cold synthetic snapshots removed after path/reparse/
process checks,~1.94GB reclaimed; active fixtures/newest backup retained. Preserve
the four unrelated untracked Workspace histories. Continue all mandatory brief work.

## Preview44 assembly checkpoint

Chat6 and MyAccountspreview-20261020.1 are now PUBLIC,15/14 assets respectively.
Every uploaded asset's server SHA/size and each tag's source commit were verified;
Actions read API confirmed disabled in all three repos before publication.
Chat clean sourceec0885eedf45e1e845333bbbcb5f95c90e8ca963, EXE
5dca548066ede54932b2b70bacb4c9c56ae0818a055e7961108890941467833a,
ZIP2b0e58d8f3b4469f86496a72393dbcdf1e6e5558645ce5104f8ed9085cced498.
Provider clean source87a1c98d83f2ec1751e5d8eac114b94ba30c4660, EXE
c3a309d396774427aeb9fdf0db355b70fb37a107abdd7b7156fe817744ec1d34.
The component lock/default versions now target those bytes; no44 build/release yet.

Exact final Chat6 passed real live new-member query, three-account historical
audience receipt transitions, forged/removed-member/membership retry, baseline,
security,7viewport×2themes and native restart. Historical receipt retention after
subscriber deletion/rejoin needs a durable ledger; full family remains incomplete.
Exact provider passed root/focused tests+vet, fresh independent OIDC and encrypted
SQLite restore, existing synthetic DB/negative PKCE/concurrent claim and retained
Chat identities/history/actual UI filename+document bytes. Source processing/cache
on I was slow; provider built from a clean temporary detached C NTFS worktree of
the same commit, frontend before native. Toolchain copy/fixtures are scoped C
scratch to clean after preserving evidence. Owner42 still untouched/healthy.

Routing helper also now compares exact hostnames case-insensitively to refuse a
foreign mixed-case route. Final full-suite/manifest rerun passed. Source
currently uncommitted44: routing/tests,2version scripts,2component pins,CHAT and
RELEASE-STATUS/UPDATE-PUBLISHING docs plus this handoff. Four other untracked
histories are unrelated and must stay unstaged. Need build-org from committed
exact source,20org pins, new immutable44 assets/source/hash mapping, actual
public42 old-manager offered6 refusal + final44 fresh/upgrade and rollback checks,
then normal owner update/remote health. No user data as fixtures. Entire mandatory
brief remains active; current releases are incomplete and plaintext, not parity.

Published IDs: Chat6 release417921774, provider20 release417922134. Their app
manifests were generated with the actual Workspace tool; Accounts includes both
owned notice files, not only its EXE. Provider inventory1013 declared Go/frontend
build components is explicitly not an audit. Primary prototype is being restored
through actual downloaded public43 setup to its verified pre42 rollback snapshot
update-20261020T002529-b668f2ee318d for an actual old42 offered6 check and future
public44 paired upgrade. Do not re-run native test-client metadata against that
rolled-back provider without recreating its synthetic registration. Fresh31060
remains final Chat6 with public43 manager/old provider for final native checks.
