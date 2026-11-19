# Sunny & Moon for Android

[Download Sunny & Moon 1.0.0](https://github.com/oneworkspacestudio/one-workspace/releases/download/android-v0.1.0-preview.22/SunnyMoon-Android-1.0.0.apk).

Open the APK on an Android 9 or newer phone and allow installation from the
browser or Files app when Android asks. Open Sunny & Moon from the home screen.
The game fills the screen in portrait and landscape. Swipe from the screen edge
to reveal Android navigation. Back closes the open settings or earlier check-in
dialog, then leaves the game.

Sunny becomes sleepy as bedtime approaches and wiggles when a toilet break is
due. Tap the toilet button after a visit, or use the earlier check-in button to
record a visit from school. Settings start at 9 pm bedtime, 7 am wake-up and a
toilet reminder one hour before bed. The morning brings back a happy Sunny.

The app works offline. Settings and check-ins are saved privately on this phone
and survive closing the app and APK updates. This edition has its own routine,
separate from the routine saved in a Workspace account. Uninstalling the app or
clearing its Android storage deletes its saved settings and check-ins. Bars
follow the clock. Children can go whenever their body tells them to.

## Build and release

The `:sunnymoon` Android module packages the canonical game from
`internal/sunnymoon/web` and substitutes its routine request with
`src/main/web/device-routine.mjs`. The server edition keeps its authenticated
routine API. No runtime downloads, network permissions or server connection
are required for this APK. The local adapter validates saved routines, clears
expired check-ins and uses revisions to avoid overwriting a newer save. Android
confirms each save on disk through private preferences before displaying success.
Earlier development WebView saves migrate automatically on first open.

Build using `scripts/build-sunny-moon-android.ps1`, adding `-Release` for permanent
signing. Use JDK 17+, SDK 35 and the repository Gradle wrapper. Supply the existing
four `WORKSPACE_ANDROID_*` signing variables through the local environment.
Keys stay outside source. The APK records the committed Android and game trees.

Publish with `scripts/publish-release.ps1 -SunnyMoon` and an Android version tag.
It builds and validates both signed APKs, includes their source and licences,
checks public downloads and promotes the public channel. The shared Android
source ZIP contains everything needed to rebuild Sunny & Moon.
