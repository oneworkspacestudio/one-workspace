# Sunny & Moon Android

The user requested an installable Android APK of the existing game. The
`:sunnymoon` Android module packages the canonical Sunny & Moon preview.4 web
assets with a local routine adapter. The managed server edition is untouched.
The offline Android edition has an independent per-phone routine, private
Android preference storage, no network permissions and permanent release signing under
the existing Workspace Android certificate. Package `uk.co.m8e.sunnymoon`,
version 1.0.0, versionCode 1. Default bedtime remains 21:00, wake 07:00 and
the toilet reminder 60 minutes before bed. Earlier visits, expression cues,
bedtime check-in and morning reset are shared with the server edition.

Native Back closes a dialog before leaving. Immersive viewport supports both
orientations. Web assets use the HTTPS Android asset loader with explicit module
MIME types. No external navigation or runtime downloads are permitted. Source,
AndroidX notices, the font licence and source provenance ship with the release.

`scripts/build-sunny-moon-android.ps1` builds the game. The supported signed
publisher accepts `-SunnyMoon` on Android releases and requires both the exact
companion and game APKs. Each has checked package/version, pinned certificate
and embedded source provenance. Shared web source is in the Android source ZIP.
The companion advances to preview.22 for this Android source change. Ordinary
suite publication retains its existing signing requirements.

Routine timing tests, seven publisher tests and an offline browser acceptance
passed. Browser acceptance checked five viewport sizes, durable settings and
check-ins, evening cue, happy morning reset, revision conflict and failed saves.
Build and acceptance artifacts are in the tooling folder
`sunny-moon-android-20261119`. Final native and public checks are recorded below.
No owner installation or physical-phone acceptance was requested.

Native acceptance found that WebView localStorage could lose a just-saved visit
on immediate force-stop. The phone adapter now confirms a synchronous private
SharedPreferences commit before reporting success. The native bridge is exposed
only inside the bundled, navigation-restricted offline WebView. A single-task
activity keeps one writer. Older development localStorage saves migrate on open.
The browser adapter retains localStorage for isolated web acceptance. Native
retesting passed. An Android-only release also preserves direct packages
on the latest suite bundle when the channel needs to trim historic rows.

## Verified release

Source `f347c04` shipped Sunny & Moon 1.0.0 in the public Android preview.22
release. Channel commit `4502dff` promoted 12 verified immutable assets,
including the required signed Workspace companion preview.22, both APKs'
source/provenance/notices and checksums. The supported publisher anonymously
downloaded both APKs and checked their exact public channel entries.

Sunny APK SHA256 is
`49b3e8a4125e63d61ee461b46d94cc619ca519f26a33825b6382470ab6ff18f9`.
The published bytes match the exact permanently signed native acceptance APK.
On a synthetic Android 15 emulator with airplane mode enabled, fresh launch,
Back closing dialogs, immediate force-stop/reopen with retained settings and
visit, rotation without scrolling, and signed APK replacement passed. Eight
offline browser checks, routine tests, seven publisher tests and six channel
tests passed. No physical-phone or owner-installation claim is made.

The emulator needed NTFS configuration and an isolated copy-on-write data
overlay on the build drive. It was stopped after acceptance and its temporary
SDK/AVD/overlay files removed. The system GitHub CLI 2.5.1 timed out while
listing asset-heavy releases during upload. Resuming the same draft with the
existing current CLI under the tooling folder succeeded without replacing
assets. Keep using a current GitHub CLI for publication.

[Download the APK](https://github.com/oneworkspacestudio/one-workspace/releases/download/android-v0.1.0-preview.22/SunnyMoon-Android-1.0.0.apk).
The same APK is copied to the outer workspace's output/sunny-and-moon folder.
