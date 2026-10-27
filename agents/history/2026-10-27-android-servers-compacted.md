# Android remote server setup

## Request and starting state

User wants Android users to choose their remote Workspace and Chat servers once,
then reopen/use normally. Base is 4cab127 on codex/android-workspace, clean before
this task. See 2026-10-27-my-push-compacted.md for released Workspace94/Chat11/
MyPush2/Android2 and the prior owner installation and phone acceptance boundaries.
This task changes source only; no new publication or installation was requested.

## Implementation

- MainActivity first-run and Settings → Servers offer a Workspace HTTPS address
  and an automatic Chat option, or a separate Chat HTTPS address. Current choices
  populate the editor. No owner URLs, accounts or credentials are hardcoded.
- Validate Workspace discovery and explicit Chat gateway/address before saving.
  Detect its canonical UI base from its same-origin sign-in redirect, including
  bare-hostname input for a /chat deployment. Discover the separate Chat issuer
  for WebView navigation. No TLS bypasses, wildcard origins or cleartext support.
- Encrypt the server choices together with discovery in SecureStore; older
  home/discovery preferences remain readable. Cold launch uses saved choices,
  catalogue and valid sessions without opening setup again.
- Launcher Chat tile, native Chat tab, app-picker Chat links and native device
  linking use the selected Chat server. Push remains discovered by the selected
  Chat server's authenticated native bootstrap, never a manually entered key.
- Validation failure/cancellation keeps the existing configuration. Saving the
  same effective servers retains the account session. Confirmed server switching
  stops background work, requests revocation of the old device grant, clears this
  app's cookies/web storage and atomically replaces encrypted cached account,
  history, queued replies and configuration. Server data is untouched. Remote
  revocation is best effort if the old server is unavailable.
- Async checks/native grants are fenced against cancellation, activity destruction
  and another activity changing the selected server. Disconnect also clears
  saved catalogue/discovery rather than leaving stale account data on disk.
- android/README.md documents setup, switching and same-MyAccounts SSO limits.
  It explicitly distinguishes these changes from the published Android2 APK.

## Verification and delivery boundary

Java/Android compileDebugJavaWithJavac completed successfully, including the final
canonical-base/link refinements (final incremental build:29s).
An initial Gradle invocation failed before compilation because its old temp path
had been removed; the task now uses an existing named temporary folder on I.
No automated tests added or run, no phone/UI/live-account testing, no APK install,
version bump, commit, push or release in this task. Builds/cache remain on I; no
new C test installation or cache was created. Physical acceptance from the prior
task remains pending; do not claim a phone sign-in/notification test here.

## Resume

Review/build the intended Android release with a new version and the permanent
private signing key when publication/installation is requested. Do not overwrite
the immutable published Android2 asset or use owner accounts as test fixtures.
Keep server-selection first-run, same-address persistence, failed validation,
cancelled/successful switching, separate issuer and phone restart/offline UI in
the future explicitly authorized acceptance scope.

## Release authorization

User subsequently requested release. Preparing Android3 (versionCode3) under
`android-v0.1.0-preview.3` in the existing org Workspace repository, with the
permanent signing key and source ZIP/notices/checksums. This Android-only tag is
excluded by the desktop updater's version-tag selector; no desktop bundle rebuild
or owner installation is needed. README now links directly to Android3. Regenerate
and anonymously verify the public channel after publication; Actions remain off.
