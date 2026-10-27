# My Push replaces the Android Firebase dependency

User asks for self hosted Firebase-style notification app and Chat integration.
Local new I:/Projects/my-push implements native Go Centrifuge0.40.1 console/broker.
Read its agents/history/2026-10-27-my-push-compacted.md for complete source,
security, artifacts and remaining acceptance. No publication/install/tests asked
or performed. Prior Android history's required Firebase-project question is
superseded: Android preview2 removes Firebase and uses My Push.

This worktree retains previous Android edits and adds push branding31stentry,
apphost/SSO/navigation/default installer/update group, package/recovery storage,
loopback service-key env and public authenticated phone transport through tunnel.
build.ps1 accepts MyPush/Version/Notices; build-org requires explicit development
package or a real published myPush pin, never a fabricated release. Existing org
pins/channel unchanged. check-org-catalog schema guard corrected from stale2 to
actual current3 and includes push. No full installer payload built/published.

Linux apphost compile support delegates to existing calendar/shared/secure;
SaveSecret/LoadSecret use path-aware wrappers, Windows DPAPI unchanged. My Push
module replaces calendar with sibling source for the development cross-build.
No claim of full suite Linux support. Preserve original unrelated Linux work.

Manager and Chat Windows candidates, My Push Windows/Linux, Android debug2
compiled. Artifacts I:/Projects/tooling/my-push-20261027. Tests/live/physical
acceptance not run. UI sample-only desktop/phone renders inspected. Final source,
hashes, package metadata and cleanup evidence in receipt files there.

Current baseb6f88cc remains older than now-public main93; before release reconcile
latest main and all newer app pins. Preview94 in local My Push manifest is only
provisional minimum; no version reservation or external promotion. Need explicit
publish/install task, real-phone acceptance and permanent Android signing key.

## Authorized release continuation, 27 October

User now authorizes physical-phone testing, push, release and installation. Local
integration rebased onto public a96102b (Maps8), preserving MyAI14/Tasks4 and all
newer pins. Source commit before actual releases contains the apphost contract;
new My Push module will pin it publicly before reproducible Windows/Linux builds.
Linux helpers now live in internal/portablesecure with Apache notices, replacing
the unpublished sibling calendar dependency. Full suite Linux support is separate.

Automated distribution/apphost/suite tests pass on a fresh protected NTFS root.
Cleanup test no longer assumes physical free-space delta is positive for small
resident files; actual removal and newest verified recovery remain asserted.
My Push synthetic isolation/ACK/revocation tests pass. Android debug2 installed
on the explicitly authorized attached phone (no previous package); device UI
acceptance pending unlock. New permanent release key held outside Git in protected
user storage. Synthetic HTTPS fixture uses isolated Quick Tunnel config; never
load owner named-tunnel config. No public release or owner update yet.

My Push1 and Chat11 are now publicly released with all uploaded asset SHA256/size
checks matching local builds, then individually promoted to the public channel.
Chat source15996b0d fixes release mirror text checks to normalize Windows CRLF;
all substantive text and binary-artwork checks remain. Final MyPush1/Chat11
binaries reran the synthetic HTTPS/native/WebSocket/ACK/revoke acceptance. Signed
Android preview2 package/version and v2 APK signature verified with permanent
release certificate. Physical UI remains blocked by locked phone, secondary update
by expired owner browser session. Preparing full Workspace94 with actual org pins;
no owner update started at this point. No Actions.

## Delivered previews and remaining human steps

Public Workspace94 source73055a0 contains Chat11 (15996b0d), My Push2 (23161db)
and signed Android2. All24 Workspace assets match GitHub size/SHA256; the public
channel promotion1efc1e9 was verified anonymously. No Actions. My Push1 was
superseded/withheld before owner installation: its default port collided with
My AI. Push now uses9021 (managed base+31); cmd/catalogcheck validates the complete
catalogue before packaging. Fixed stale Maps floor, MSVC runtime argument/seed,
and optional publisher metadata authentication in the organisation build script.

Distribution/apphost/suite tests and scoped push tests pass. Final binaries pass
16 HTTPS/native protocol checks, including a real server publish acknowledgement,
opaque WebSocket wake, queue/ACK, native sync/reply/session restoration/revocation.
An optimistic direct-fallback send is not a server commit: the test now waits for
the actual publish ACK. The normal protected93→94 fixture upgrade automatically
added Push and preserved26 excluded apps. Live console SSO/branding/navigation
and390px no-overflow checks pass. See docs/ANDROID-DELIVERY.json.

Primary owner installation completed through protected manager then Chat updates;
all28 services running, old credentials/settings and unselected executable bytes
retained (MyAI14/Maps8/Tasks4). Newest verified recovery kept; normal cleanup only.
The phone has the exact published release-signed APK. Appium session deleted and
app force-stopped/absence verified. Physical notification/reply/Doze/reboot tests
remain unperformed because the phone stayed locked. Secondary installation NOT
performed: its MyAccounts browser session requires the user to sign in. Resume
from a fresh secondary Workspace root URL, not a stale OAuth callback/state URL.
No owner account/data was used as a development fixture.

All synthetic services/three Quick Tunnels stopped. Named C test/build tree
removed (905422159 logical bytes); source/release/evidence stays on I and the
permanent Android signing key remains outside Git in protected user storage.
The synthetic credentials/data were removed; recreate a fresh isolated fixture
for physical acceptance after the user unlocks. Do not rerun primary install.
Full-suite Linux/gateway support remains separate; only My Push Linux x64 cross-
build is certified here. Original unrelated Linux worktrees remain untouched.
