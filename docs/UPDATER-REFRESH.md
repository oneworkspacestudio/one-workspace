# Update pages after a Workspace upgrade

Workspace checks that an open updater matches the installed manager's executable
SHA256 as well as its communication protocol. An older idle updater closes through
its protected maintenance API before the current manager opens a replacement.
An installation, download, backup or recovery already running is allowed to finish.

A local update page opens the current updater when it notices a completed manager
replacement. Remote owner administration reconnects through the installed manager.
If an updater cannot confirm that the old session is idle, it shows a retry message
and preserves that session. An unavailable app list displays an explanation.

The updater recognises a manager replacement before interpreting a newer app
catalog. It reads the installed version from a stable header rather than labelling
an older recovery configuration as the current installation.

## Preview 107 release scope

Preview107 delivers this fix for Windows. The maintainer explicitly requested
this one release without Android because permanent Android signing is configured
on another PC. No new Android APK accompanies107; the existing signed Android
release remains available. The normal signed Android release policy is unchanged.

Existing application payloads and contracts are retained from the pinned published
organization components. User accounts, app data and libraries follow the normal
protected update and recovery path.

Regression checks cover binary identity even when protocols match, legacy markers,
unreadable future catalogs, busy/ambiguous sessions, close races, stale mutation
rejection and a single local handoff. Native acceptance uses synthetic NTFS data
and public binaries; it never uses owner accounts or files as fixtures. Release
acceptance receipts identify which native checks ran. No physical-phone, reboot,
clean-host or production-certification acceptance is implied.
