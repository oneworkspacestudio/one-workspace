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
