# Workspace Backups

Available in Preview 28. Windows builds and package integrity checks passed;
runtime installation, scheduled backup and restore acceptance remain pending.

## Set up

Sign in as the Workspace owner and open **Backups** in the app grid or the
Workspace apps switcher. It is an internal app; there is no separate executable
or account. Scheduling is **off by default** on new and upgraded installations.

1. Choose a dedicated backup folder, preferably on a different physical drive,
   such as `E:\Workspace backups`. **Browse server folders** lists drives and
   ordinary folders on the Workspace server directly in your browser, including
   on a phone or remote PC. Choose a parent and enter a new folder name in the
   destination field if needed. **Check folder** still validates the destination.
   The browser picker never shows files or follows junctions. Large folder
   listings are bounded; you can enter a full server path directly.
2. Use **Check folder**. Private account data needs Windows ACLs, so the destination
   must be local NTFS storage. FAT, exFAT, mapped network shares, junctions, drive
   roots and folders overlapping the installation are refused. A same-volume
   destination works, but does not protect against that drive failing.
3. Choose an interval of 1–168 hours (24 hours initially), check **Enable scheduled
   backups**, and save. You can leave scheduling off and use **Back up now**.

The PC and Workspace must be running. A missed scheduled run is picked up within
30 seconds of startup. Runs do not overlap updates, recovery or each other.
An unavailable/full drive produces a visible error; failures retry after
30 minutes while Workspace remains running. Restarting Workspace may retry a
missed run sooner. Copies created by the Backups app are retained until the owner removes them; no automatic
pruning is performed.

## What happens during a backup

The page checks storage and free space before pausing apps. The app processes and
their descendants are stopped, along with this installation's shared service.
Workspace's portal and Cloudflare connector stay running so the owner can see the
current phase and copied bytes. App requests receive a temporary unavailable
response until the run finishes. Native connections such as SMTP may be interrupted
and must retry. Save current work before a manual run. Unsaved browser/editor
content is not server data, and external editors must not change files during copying.

Workspace copies accounts, databases, app content, installed programs and ordinary
files. It verifies copied file hashes, alternate streams and supported Windows
metadata using the existing update snapshot format. Only verified snapshots enter
the recovery list. Apps are restarted and checked after success or failure. A
restart error is reported separately from the presence of a verified snapshot.
The page shows the last verified backup, next run, history and remaining phases.

Allow at least twice the source size plus 512 MB free on the destination for a
snapshot and recovery archive. Large libraries cause a longer interruption; this
is a full copy, not an incremental backup. Logs, runtime scratch, previous backups
and maintenance state are excluded to prevent recursion and copying changing files.
Reparse points, EFS files and directory alternate streams are refused; no incomplete
copy is presented as a verified snapshot. Partial directories can remain after a
failed run; keep them separate from verified backups until reviewed.

## Storage layout

The single installation folder contains:

```text
Workspace/
  Apps/                       Verified program payloads and licenses
  Private/                    Owner-protected settings and coordination files
    accounts/                 MyAccounts identity and settings
    drive/                    Drive settings, database and cache
    photos/                   Photos working directory
    gitea/                    Git repositories and configuration
    music/                    Music settings and library index
    mail/                     Mail engine, messages and hosting configuration
    calendar/                 Calendar private state
    contacts/                 Contacts private state
    shared/                   Calendar/Contacts shared service and database
    notes/                    Notes, attachments and history
    my-tube/                  Private videos and metadata
    my-vault/                 Vault encrypted content and configuration
    runtime/<app>/temp/       Per-app runtime scratch (excluded from backup)
    music-runtime-temp/       Existing owned Music extraction scratch
  Files/
    Documents/                Ordinary documents and other user files
    Photos/                   Ordinary pictures
    Music/                    Ordinary music files
```

Only selected apps create their private folders. Common settings and credentials
remain in the protected Private root. Workspace helpers/shared services inherit
protected Workspace scratch; process apps have their own TEMP/TMP and working
folders. Installer downloads are staged on the selected program drive and removed
afterwards. Windows shortcuts, autostart registration, trust and OS-managed caches
remain Windows resources. Existing three-folder installations keep their paths;
this change does not move existing user data automatically.

Private content such as mailboxes, notes, private videos and password vaults must
remain access controlled. It is not exposed as shared Drive files. Documents,
photos and music that are already ordinary shared files remain in Files. Standalone
apps launched outside Workspace retain their own standalone storage settings.

## Recover

Open **Backups → Open recovery tools → Open update window on this PC**, then choose
a verified snapshot. Recovery runs locally on the Windows server, stops the suite
and makes a separate archive of its current state before restoring. Newer ordinary
files and the recovery archive are preserved by the existing rollback procedure.
See [Updates and recovery](UPDATES.md).

Keep the entire external snapshot, including `backup.json` and `receipt.json`.
The receipt records its manifest checksum and installation paths. The local recovery
list also records it under `Private/updates/backups.json`. Do not edit manifests or
move their members independently. If the original disk or local recovery index is
lost, retain the external copy and request administrator-assisted recovery; this
version does not provide a bare-disk restore wizard or automatically import receipts.

These are file snapshots, not bootable Windows images or portable account exports.
Workspace secrets use Windows machine DPAPI, and restore requires the original
Windows security context. A replacement Windows installation cannot automatically
decrypt those secrets. My Vault's independent encrypted export/recovery process
is still needed for cross-machine vault recovery.

## Advanced installation

**Apps & folders → Advanced: choose your apps** lists the bundled apps, all checked
initially. MyAccounts stays required. Checking an app checks its dependencies;
unchecking a dependency also unchecks dependent apps. The server validates the
selection again against the downloaded bundle before writing installation data.
Omitted standalone executables and unneeded shared engines are not extracted or
launched. Licenses and provenance remain available. Updates preserve omitted apps.

Docs is an entry to an editor embedded in Drive: unchecking Docs removes that tile;
it cannot remove editor assets already compiled into Drive. Photo Editor's omitted
standalone Studio.exe is not extracted, while Drive's embedded editing remains part
of Drive. Backups stays available as an internal owner app with scheduling off.

Normal installation checks GitHub for the newest published suite. Use `--offline`
to deliberately use the installer’s embedded versions. A newer installer refuses
an older published manager instead of silently losing its features.
