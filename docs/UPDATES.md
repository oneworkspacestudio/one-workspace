# Update an existing Workspace

Download the newest **WorkspaceSetup-windows-amd64.exe** from the
[organization releases](https://github.com/oneworkspacestudio/one-workspace/releases).
On its first screen, expand **Already have One Workspace?**, enter your existing
Workspace parent folder and open **Updates and recovery**. For an installation
with separate storage locations, use the command below with its Private folder.
This bootstraps installations made before the update window existed.

After the first update, the Workspace owner can choose **Updates and backups**
under Workspace administration from any device at the normal Workspace address.
Check, install, reinstall, recover and restore run on the server with progress in
your browser. An idle older controller is replaced automatically. An active
operation is never interrupted to open another window.

The private update controller stays on loopback; its capability is never sent to
the remote browser. Every remote action requires owner access and a matching
Workspace origin. A protected, HttpOnly owner continuation lasts at most two
hours so an existing update page can reconnect after a manager restart; signing
out revokes all update continuations. If it expires, sign in again to resume
viewing progress. Updates may briefly disconnect the tunnel while services
restart. Keep the page open; do not repeat an installation request during that
period. Your PC must be running and reachable. Backup folder paths always refer
to the server, including when administration is opened from a phone.

To repair program files or install a required newer installer, choose **Check for
updates**, then **Back up and reinstall Workspace**. The update window downloads
and verifies the installer from that exact release, runs it automatically, saves
a recovery copy and restarts the suite. Existing accounts, selected apps, settings
and files stay in place, including independently updated apps. No new password is
required. Reinstall results list the restarted services; use **Open App status and
retry** if any service needs attention. Keep the window open until checks finish.

This also handles future releases that need an installer format the old update
window cannot read. The old window hands execution to the verified new installer
instead of trying to install an unfamiliar application catalog. If prompted after
completion, close the old updater and reopen **Updates and backups** in Workspace.
Older update windows without this action must be bootstrapped once using the latest
downloaded installer.

1. Schedule a maintenance window. Close editors and pause external writers such
   as mail delivery, sync jobs and independently configured app services until
   startup checks finish. Recovery retains newer state but does not merge concurrent
   mailbox/database writes automatically.
2. Choose **Check for updates**. Review installed and available versions.
3. Select Workspace, individual apps, or all available updates. Drive/Docs/Photos
   update together; Calendar/Contacts and their shared engine update together.
4. Optionally choose a larger local NTFS backup folder. The default is Private/backups.
5. Choose **Back up and install selected updates**. Wait for **Update complete**
   after verification and startup checks, then choose **Open Workspace**.

You can leave the update window open. If its release check is over 30 minutes old,
install and reinstall refresh it automatically before changing anything. They
continue only when the selected program files still match what you reviewed.
If releases or the installation changed, review the refreshed choices and try
again. Failed requests stay visible while the page checks status.

Downloads and every extracted executable/helper/notice are SHA256 and size
verified before downtime. Installed program hashes are checked too. A stale plan
is refused if another operation changed the installed catalog. The updater only
uses published organization assets with valid GitHub download URLs and digests.
The withheld Music executable remains blocked.

The updater retains your Programs, Private and Files locations, owner, identities,
OIDC clients, hosting settings, disabled optional apps and custom extensions.
It refuses an automatic change to existing data paths, storage environment
bindings, private Mail backend or OIDC callbacks. Those need an explicit migration.

### New apps are included

Updating **One Workspace** also installs newly introduced default apps, including
My Hosting. They appear as **Not installed** in the release list and are included
with the Workspace selection. New apps receive their own private storage and a
MyAccounts client automatically, using your existing owner account and sign-in.
The update is complete only after the new apps pass the normal startup checks.
Their installation and identity changes are covered by the same verified backup
and rollback transaction as the manager update.

Advanced setup remembers the apps you uncheck. Future updates respect those
exclusions and exclude apps whose dependencies you omitted. If an installer
downloads a newer bundle than the app list it originally displayed, new defaults
in that bundle are included; an absent checkbox is not treated as an opt-out.
Both online updates and explicit offline bundles follow the same policy.

Old installations did not record exclusions. The compatibility migration preserves
missing apps from the pre-Hosting app list as omitted; My Hosting and future app
IDs are considered new defaults. This deliberately avoids re-enabling an older
optional app that may have been unchecked. Reinstalling is not required to add
Hosting. Use the corrected installer/updater once to repair an installation made
by an older updater, even if it already reports the current Workspace version.
When an older controller downloads a corrected manager, the manager also finishes
adding missing defaults during that update's startup checks. It uses the exact
matching verified release and the original transaction's backup. Normal launches
do not download or install apps. Explicit offline repairs should use the corrected
installer with the complete bundle.

## Follow update progress

The local updater shows the current stage and a checklist of completed and
remaining work: download, package checks, restart preflight, backup, backup
verification, installation and app startup. Downloads name the actual release
asset, including the complete application bundle when a changed release notice
comes from that bundle. Internal package IDs are not download descriptions.

The current-stage bar shows measured bytes or completed files/packages/checks.
Download and backup stages show speed and an approximate remaining time after
work begins. These estimates apply to that stage; later stages need additional
time. Overall elapsed time and activity details remain visible. Slow operations
show how long progress has been unchanged; a disconnected browser reconnects
without declaring the update failed or enabling another operation.

A transfer reaching100% means its bytes have arrived. **Update complete** appears
only after the transaction and app startup checks succeed. Both progress bars
disappear when an operation finishes or fails. Completed updates show a static
result and **Completed in** duration. Failures retain the error and recovery
guidance instead of presenting a completed update.

## Backups and rollback

The updater stops the selected installation's supervisor, apps and persistent
shared service, then verifies that its ports are released. It creates a complete
backup of program files, private databases/configuration/credentials/media and
ordinary files. Empty directories, modification times, ordinary-file permissions,
file attributes and regular-file NTFS alternate data streams are preserved.
Backup contents and alternate streams are hashed and verified before replacement.
Backups are protected on a local NTFS volume, including when programs/files use exFAT.

Allow space for **two full copies** of the backed-up installation plus working
space. Downloads stage under Private/updates, so that volume also needs staging
space. Existing backups, diagnostic logs and disposable runtime extraction caches
are excluded from each new snapshot. After a healthy update, automatic cleanup
keeps the latest verified update recovery copy and removes older identified
automatic update backups. Scheduled/manual snapshots, recovery archives and
unknown legacy copies remain. Backups are never purged before a new update to
make room. A failed backup prevents any program replacement. See
[automatic update cleanup](UPDATE-CLEANUP.md) for the retention and safety rules.
Custom app storage outside Programs, Private and Files is not part of this
snapshot; back up those external locations separately before maintenance. Current
installed versions and source URLs are recorded in the protected installation
catalog; bundled source-pin notices remain provenance for the original package.

Every app must pass startup before the update completes. A failed update restores
the verified previous state. Before restoring, the newer state is separately
archived so changes remain recoverable. Failed or interrupted states are labelled
as recovery archives and cannot be selected as healthy full-installation backups.
Extra ordinary files are retained when an earlier snapshot is restored.

### Windows could not restart Workspace

The updater checks it can launch the installed manager before stopping apps.
It respects Windows host process-group limits when starting maintenance or the
replacement manager. A terminal or installer may forbid process breakaway;
requesting that option unconditionally can produce an "Access is denied" error
even when the executable has the correct file permissions.

If an older updater reports **data restored, but startup needs attention**, the
restore completed but the apps are stopped. Open the existing Private folder's
Workspace.exe to start the restored version. Use a corrected installer/updater
for the next attempt. **Recover installation** is for an unfinished transaction;
it cannot fix a Windows launch restriction after a completed rollback. Keep both
verified backup and recovery archives. If launch still fails, check the executable
permissions and Windows Security's protection history; do not disable protection.

To manually undo an update, choose a recorded backup and **Archive current state
and restore backup**. This creates a verified snapshot of the current installation
first. You can restore that newer snapshot again if you need post-update edits.
The protected backup directory also contains the files and app databases for
individual recovery. The newest automatic update rollback is retained; use the
separate Backups app for long-term snapshots that update cleanup does not retire.

A durable transaction journal blocks normal startup during an unfinished update.
Manager replacement is ordered before app replacement; recovery retains a verified
controller until all data and programs are restored, including for legacy launchers.
If the process or PC stops, reopen the newest installer’s update window and choose
**Recover installation**. Recovery resumes even if workspace.json was interrupted.

These are application backups, not Windows volume images. Reparse points/junctions,
EFS-encrypted files and directory alternate data streams are refused rather than
silently copied without their special semantics. Use an appropriate Windows backup
for those storage layouts. Directly managed external services must stay stopped
during maintenance. Clean-host, real power-loss/reboot and Windows service acceptance
remain separate preview gates.

## Command line

Run the **new installer**, rather than a running manager being replaced:

```powershell
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private"
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --check-updates
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --apply
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --apply --apps workspace,notes
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --apply --backup-root "E:\WorkspaceBackups"
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --recover
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --rollback BACKUP_ID
```

Update groups: workspace, accounts, drive, gitea, music, mail, studio, calendar,
notes, my-tube, vault, chat, hosting. New organization apps with a compatible
suite contract are discovered from the verified suite catalog. Older installations must update Workspace first or select it
together with their app updates. No password or recovery credential is needed in
the command line. An administrative Windows user still needs access to the
installation’s protected private storage.

For a deliberately offline update, download an authentic organization
workspace-bundle.zip, compare its SHA256 with the release checksum file, then use:

```powershell
WorkspaceSetup-windows-amd64.exe --update "C:\Workspace\Private" --apply --bundle "D:\Downloads\workspace-bundle.zip"
```

The offline option verifies the catalog and every packaged file. Authenticity of
a locally supplied archive is the local operator's responsibility. Ordinary
online updates verify the organization release digest through GitHub.

## Independent app releases

Workspace checks the complete suite release and each installed app repository.
A newer compatible suite pin can update an individual app without updating the
others. Independent releases can update sooner by publishing the explicit
**workspace-update.json** contract described in [Publishing app updates](UPDATE-PUBLISHING.md).

A newer release without that contract is shown for review and left unchanged.
This avoids assuming that an arbitrary standalone executable supports existing
managed data, arguments, helper binaries and shared authentication. It becomes
automatically available when a compatible Workspace bundle includes it, or the
publisher supplies the contract in a new release. A partial/failed repository
check is reported instead of silently installing an older fallback.

## Optional scheduled backups (Preview 28)

The internal owner **Backups** app adds full verified snapshots on a configurable
interval, switched off by default. See [Backups](BACKUPS.md) for setup, downtime,
storage and recovery limits. Its snapshots appear in the same local recovery list.

## Release ordering

Workspace55 chooses the highest published Workspace version that supplies a bundle.
Publication dates only break ties between equivalent version labels. An older
release published later is not offered as a manager downgrade.


## Public update list and GitHub limits

Public installations read one publisher-maintained update list and download
verified files from public release URLs. A GitHub account or token is not needed.
The update window shows when the publisher refreshed the list. If it cannot be
read, retry when your connection is available; no program files are changed.
Authenticated/private repository checks may still show GitHub's API retry time.
Publishers must follow [release promotion](../updates/README.md) after releases.

## When the connection drops during an update

Keep the Updates page open while Workspace and its tunnel restart. A short
reconnection notice replaces gateway error pages; the last received progress
remains visible. Controls stay disabled until the page receives current server
status. A lost response does not resend the install, reinstall or restore request.
If the server PC or its network stays offline, restore its connection first;
the browser cannot restart an unreachable PC. Refresh an older open Updates tab
after installing the new manager to load the current controller.
