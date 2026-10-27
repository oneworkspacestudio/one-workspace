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
