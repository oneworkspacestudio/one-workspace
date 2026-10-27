# One Workspace for Android

An Android launcher for the apps installed on your own One Workspace server.
Enter the **HTTPS Workspace address**, sign in with MyAccounts, and open an app.
The launcher caches the app list; each app retains its existing server and data.
Chat uses the existing mobile conversation UI and its emoji/sticker features.

## Android integration

- Shared WebView cookies for Workspace/MyAccounts app sign-in, trusted server
  origins, a saved launcher, native file selection and Android back navigation.
- Existing Workspace servers can supply their normal secure launch redirects;
  new managers also publish `/.well-known/one-workspace.json`.
- Chat device credentials encrypted with Android Keystore, separately revocable
  on the server. Native session restoration avoids another Chat login while the
  linked-device grant remains valid. The phone never stores your password.
- Android message channels, conversation notifications, unread badges, private
  lock-screen content by default, reply and mark-as-read actions, and a tap into
  the corresponding conversation. Actions are bound to the linked account.
- My Push, your self hosted notification service, carries opaque wake signals.
  No Firebase Android SDK, Google project or service-account key is needed.
- A visible Android remote-messaging foreground service keeps the private
  WebSocket connection open, reconnects with backoff, recovers pending signals
  and acknowledges receipt. It can be disabled in native Settings.
- WorkManager provides bounded authorized Chat sync and a periodic fallback of
  **at least 15 minutes**, subject to Android power and network policy.

Background delivery needs the matching Chat server, My Push and new Workspace
manager integration. Workspace configures the service connection automatically.
Older servers still open their web Chat. Force-stop, denied notification
permission, offline servers and manufacturer battery restrictions can prevent
or delay alerts. The app does not bypass those controls or guarantee delivery
while Android suspends network access.

This is a development implementation, not full WhatsApp feature parity. Existing
Chat E2EE/calling limitations remain. Native Telecom/CallStyle incoming calls,
background camera/microphone calls, Play distribution and physical-device/Doze
acceptance are not claimed. No automated or device tests were requested/run.

## Build locally

Requires JDK 17+, Android SDK platform 35, and Gradle 8.11.1. Gradle wrapper files
are included. Android 9/API28 or newer is supported; target SDK is35.

```powershell
./scripts/build-android.ps1 -Sdk "$env:LOCALAPPDATA/Android/Sdk"
```

The default output/cache location is `I:/Projects/tooling/workspace-android`.
If a filesystem cannot handle Gradle's many small files efficiently, use a named
temporary NTFS cache and remove that cache after copying the output elsewhere.
Never use a live Workspace installation as a mobile development fixture.

Debug APKs use a development signing key. They are not Play/production releases.
Release builds fail unless all four environment values are supplied:
`WORKSPACE_ANDROID_KEYSTORE`, `WORKSPACE_ANDROID_STORE_PASSWORD`,
`WORKSPACE_ANDROID_KEY_ALIAS`, `WORKSPACE_ANDROID_KEY_PASSWORD`.
Keep that permanent signing key outside Git and reuse it for app updates.

## Chat server setup

The companion source is `workspace/native.go` and `workspace/native_push.go` in
the One Workspace Chat repository. It must be built/deployed before enabling
native background access. Existing web Chat continues to work without My Push.
See its `docs/workspace/ANDROID.md` for the protected server configuration and
the limits of this development build. Service and device credentials never belong
in source, logs or public release assets.

After signing in, the phone links to the current Chat identity for up to30 days.
Open Chat again to renew expired access through normal sign-in. The server checks
the current MyAccounts directory, retains membership checks, and revokes access
when the account is removed or Workspace explicitly revokes its Chat sessions.
Chat's **Linked phones** screen can remove individual devices. Native settings
can disconnect the current phone. Restoring older authentication data can restore
older device records too; remove stale devices after a deliberate full restore.

## Source and licences

The Android source is part of the Apache-2.0 One Workspace repository. It uses
AndroidX, Material Components and OkHttp under their upstream
licences. A notices file and Apache-2.0 text are bundled in the APK assets.
Chat retains its own Tinode GPL-3.0 and web client Apache-2.0 licensing.
