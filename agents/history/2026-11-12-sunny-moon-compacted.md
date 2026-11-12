# Sunny & Moon

Sunny & Moon preview.1 is an optional General app in One Workspace. It is
implemented in the suite repository's `codex/sunny-moon-release` worktree.
The deployment target is the managed Workspace channel. No Netlify site was
created. No owner installation was requested.

The screen has an animated face, bladder/bowel/tiredness clock cues, a one-tap
toilet check-in, optional earlier visit time, a separate evening check-in, and
a sleep opportunity estimate. Defaults are 9 pm bed, 8 pm toilet reminder and
7 am morning reset. The morning is always a fresh start. The app uses device
time and reminders are displayed in the open app. No push or medical data
integration. Workspace MyAccounts sessions isolate accounts, and the private
SQLite store keeps settings and two expiring check-in timestamps per account.

Managed ID `sunnymoon`, route `/sunnymoon`, base port 9047 (manager +57),
private directory `sunnymoon`, and installer integration 17. The browser UI is
embedded in SunnyMoon.exe. Owned packages include the compatibility contract,
guide and DM Sans licence/source. Standard Workspace setup, shared navigation,
storage validation, app recovery and update discovery all use this ID.

Working artifact root `I:/Projects/tooling/sunny-and-moon-20261112`. Use only
isolated NTFS test data under `C:/Projects/.cache/sunny-and-moon-20261112`.
The ordinary build uses the pinned public seed, builds the full app inventory including the concurrent Pet integration
and carries the latest previously promoted component pins. Publish only via
the signed Android publisher after installer and update acceptance.
