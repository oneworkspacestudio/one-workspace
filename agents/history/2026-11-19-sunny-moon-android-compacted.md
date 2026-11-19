# Sunny & Moon Android

The user requested an installable Android APK of the existing game. The
`:sunnymoon` Android module packages the canonical Sunny & Moon preview.4 web
assets with a local routine adapter. The managed server edition is untouched.
The offline Android edition has an independent per-phone routine, private
WebView storage, no network permissions and permanent release signing under
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
`sunny-moon-android-20261119`. Native acceptance and publication remain pending.
No owner installation or physical-phone acceptance was requested.
