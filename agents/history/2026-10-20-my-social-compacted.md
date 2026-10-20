# My Social integration

User authorizes creating My Social from the supplied full Instagram-style PWA
brief, local implementation, GitHub source pushes and versioned releases of
My Social and Workspace after local tests and GitHub Actions succeed. This task
explicitly requests Actions; earlier tasks' no-Actions preferences are superseded
for these two repositories. Owner installation is not a test fixture.

Sources: existing I:/Projects/one-workspace plus new I:/Projects/my-social.
Build artifacts I:/Projects/tooling/my-social-20261020; private synthetic data
under C:/Projects/.cache/my-social-20261020 (NTFS). I: is exFAT: set GOTMPDIR to
the named C test temp for Go tests needing protected state. No owner data used.

Shared apphost adds explicit opt-in MediaCapture policy and self/blob video CSP.
Existing apps retain camera/microphone denial. Shared capability regression plus
all suite Go tests/vet passed locally. Initial commit publishes this dependency
and a Windows source verification workflow; remaining integration remains local
until the application and isolated install checks pass.

Discovery: My Photos is a permission-aware ordinary-file view of Quantum. The
suite's default My Drive source currently shares Files. New adapter reserves
Files/My Photos/My Social, installs Quantum denyAll and fails closed if any
specific ancestor/descendant user/group allow overrides it. The local capability
guard rechecks native Quantum policy; Social mediates owner originals and
audience-checked derivatives. Never expose that root directly. Chat uses Tinode
and per-app tickets; Social currently offers an explicit launcher, no invented
DM/RTC API. Full parity and physical-device/clean-host checks remain open.

Managed acceptance passed: fresh three-app installation, one MyAccounts sign-in,
private publication, exact original bytes, same immutable Drive subject, Quantum
raw download/Photos media/folder denial and shared picker. Guard auth now reuses
its token in memory for two minutes; authorization rules are fetched every time
(regression tested). Full suite Go tests/vet pass. Social preview1 targeta7dfda3
passed Windows CI39410458193 and packaged EXEf8de0bce repeated Chrome workflow.
Social ZIPa78f08d7 is pinned; existing20 component app versions remain unchanged.
Final full21-app installer, CI, release and normal owner update are next.

Security correction: Windows native logical ACLs preserve casing. An isolated
real Drive probe returned404 for exact denied path but200 and identical marker
for lower/upper-case paths. Earlier exact-path acceptance was insufficient.
Do not use preview1/Workspace47 storage for private media. Published tags/assets
remain immutable; release notes mark predecessors superseded. The local owner
root contained only its marker, no media, and uploads were paused by renaming
that marker. The new adapter uses Data/photos-social/My Photos/My Social outside
all shared sources, default-disabled My Social Originals native source, and
read-only per-subject scopes; physical native scope containment provides the
boundary. Social2 requires Workspace48 and scoped-social-storage-v2.
Existing legacy media is refused, never silently moved. Managed tests and final
release48/publication/local update are pending this correction handoff.

Scoped adapter acceptance now passes on a fresh native installer: one MyAccounts
sign-in, Social upload and publication, original exact bytes via Social and
native own My Photos scope, cross-account traversal/case denial, and shared
Drive alias denial. CreateUser applies defaults to incoming source permissions;
bootstrap creates without a private scope, then PATCHes read-only scope. Scope
provisioning is serialized to avoid concurrent first-login create races.
Social source0744a86 passes CI39429278014 and packaged Chrome checks.

Delivery complete: Social v0.1.0-preview.2 target0744a86c826a978c76cbfe521b00336dbe187dc0
(release418427788,9 verified assets, anonymous public ZIP verified) and Workspace
v0.3.0-preview.48 target1e589609b86b96ca2aa2a4e61b0d5430893b21ce
(release418432236,21 verified assets) are public. Source CI39429278014 and
39430126063 passed, as did packaged Chrome and native managed acceptance.
A real published47-to48 backed-up upgrade preserved identity/issuer/paths,
ordinary bytes, Accounts/Drive hashes and18 app opt-outs, then passed own native
original identity and cross-account/case/traversal denial. Existing20 other
component payloads/versions remain unchanged.

The authorized real local installation now runs48/Social2, public-digest-matched
manager1c459041 and Social3f938d8f, with21 apps and health verified. Normal updater
created a verified rollback; preserve the newest recovery copy. Owner data was
never a development fixture. Existing operator FFmpeg/FFprobe are configured; no
helper binaries redistributed. Synthetic C fixtures were stopped and removed.
Source clones in the user workspace are synchronized. Full Instagram parity,
advanced video/audio, Social messaging/live, physical Android/iOS push/capture,
clean-host restore and independent security review remain documented gates.
The predecessor1/47 tags/assets stay immutable and their notes are superseded.
Do not revive shared-folder ACL storage. Existing nonempty legacy media still
requires a separate explicit migration; this local owner had no legacy media.
