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
