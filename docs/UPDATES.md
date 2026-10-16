# Update an existing Workspace

Download the newest **WorkspaceSetup-windows-amd64.exe** from the
[organization releases](https://github.com/oneworkspacestudio/one-workspace/releases).
On its first screen, expand **Already have One Workspace?**, enter your existing
Workspace parent folder and open **Updates and recovery**. For an installation
with separate storage locations, use the command below with its Private folder.
This bootstraps installations made before the update window existed.

After the first update, the Workspace owner can choose **Updates and backups**
under Workspace administration. The update window opens on the Windows server PC.
Remote visitors cannot receive its local maintenance capability or run updates.

1. Schedule a maintenance window. Close editors and pause external writers such
   as mail delivery, sync jobs and independently configured app services until
   startup checks finish. Recovery retains newer state but does not merge concurrent
   mailbox/database writes automatically.
2. Choose **Check for updates**. Review installed and available versions.
3. Select Workspace, individual apps, or all available updates. Drive/Docs/Photos
   update together; Calendar/Contacts and their shared engine update together.
4. Optionally choose a larger local NTFS backup folder. The default is Private/backups.
5. Choose **Back up and install selected updates**. Wait for verification and
   startup checks, then open Workspace and sign in again.

Downloads and every extracted executable/helper/notice are SHA256 and size
verified before downtime. Installed program hashes are checked too. A stale plan
is refused if another operation changed the installed catalog. The updater only
uses published organization assets with valid GitHub download URLs and digests.
The withheld Music executable remains blocked.

The updater retains your Programs, Private and Files locations, owner, identities,
OIDC clients, hosting settings, disabled optional apps and custom extensions.
It refuses an automatic change to existing data paths, storage environment
bindings, private Mail backend or OIDC callbacks. Those need an explicit migration.

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
are excluded from each new snapshot and retained in place. Backups are never
silently purged to make room. A failed backup prevents any program replacement.
Custom app storage outside Programs, Private and Files is not part of this
snapshot; back up those external locations separately before maintenance. Current
installed versions and source URLs are recorded in the protected installation
catalog; bundled source-pin notices remain provenance for the original package.

Every app must pass startup before the update completes. A failed update restores
the verified previous state. Before restoring, the newer state is separately
archived so changes remain recoverable. Failed or interrupted states are labelled
as recovery archives and cannot be selected as healthy full-installation backups.
Extra ordinary files are retained when an earlier snapshot is restored.

To manually undo an update, choose a recorded backup and **Archive current state
and restore backup**. This creates a verified snapshot of the current installation
first. You can restore that newer snapshot again if you need post-update edits.
The protected backup directory also contains the files and app databases for
individual recovery. Keep backups until you have checked your work in every app.

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
notes, my-tube. Older installations must update Workspace first or select it
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
