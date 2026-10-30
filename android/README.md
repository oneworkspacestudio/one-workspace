# One Workspace for Android

One Android app opens the web versions of the apps on your remote One Workspace
server. You do not need to install individual Android apps for Drive, Mail, Chat
or the other Workspace apps. The server continues to host the apps and their data.
On first launch, enter your **remote HTTPS Workspace address**. Leave
**Use Chat from this Workspace** enabled to find Chat automatically, or turn it
off and enter your own **HTTPS Chat address**. Sign in with MyAccounts and open
an app. Your server choices and valid sign-in sessions are remembered when you
close and reopen Android; you do not need to configure the addresses again.
The launcher caches the app list; each app retains its existing server and data.
Chat uses the existing mobile conversation UI and its emoji/sticker features.
Chat appears alongside the other apps in the launcher, with no dedicated Chat
tab. Settings are available from the launcher header. Android Back returns from
an app to the launcher; native Chat notifications and replies remain supported.
This simpler launcher layout requires Android preview 4 or newer.

## Choose or change servers

- Workspace example: `https://workspace.example.com` (the server address, without
  a sign-in page or callback).
- Chat example: `https://chat.example.com/chat`. Include the path shown when you
  open Chat in your browser. A bare Chat hostname is also checked for `/chat`.
- Open **Settings → Servers** to change either address later. Addresses are checked
  before they are saved; a failed check keeps the previous configuration.
- Switching servers asks for confirmation and removes this phone's previous
  sign-in, saved launcher, queued replies and notification access. Server data
  stays on the original server. Saving the same servers keeps the valid session.
- Chat in the Apps launcher and background notifications both use
  the selected Chat server. Push is discovered from that server after sign-in;
  users do not need to enter a separate push address or key.

Use reachable servers with valid HTTPS certificates. For seamless single sign-in,
configure Workspace and Chat with the same MyAccounts service. A Chat server
using another MyAccounts service can ask you to sign in to that account once.
The saved launcher remains available offline; connecting or signing in for the
first time needs access to the selected servers.

The separate Chat address and Servers editor require Android preview 3 or newer.
[Download Android preview 11](https://github.com/oneworkspacestudio/one-workspace/releases/download/v0.3.0-preview.109/OneWorkspace-Android-0.1.0-preview.11.apk).

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

The signed preview APK accompanies every new Workspace release. Workspace
preview 96 includes Android preview 6, rebuilt with the existing signing key
and the same launcher behaviour as preview 4. Use the Android download above.
Use Chat preview 11 and My Push preview 2 or newer. This is an evaluation
implementation; full WhatsApp feature parity is not claimed. Existing
Chat E2EE/calling limitations remain. Native Telecom/CallStyle incoming calls,
background camera/microphone calls, Play distribution and physical-device/Doze
acceptance are not claimed. Automated server authorization/queue tests and an
isolated HTTPS two-user Chat-to-WebSocket/queue/ACK/revocation flow pass. Physical
notification/reply, Doze and reboot acceptance remain pending; installation alone
does not demonstrate these behaviours.

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
Every One Workspace release must include this signed APK. Maintainers use the
[local release publisher](../docs/RELEASING.md), which builds and checks Android
before publishing; no GitHub Actions are used.

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
