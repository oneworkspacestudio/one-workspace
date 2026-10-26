# Automatic update cleanup

Workspace preview.39 cleans up after the updated apps restart successfully.
The updater shows **Clean up old update files**, then reports how many older
copies were removed and approximately how much space was freed. The result and
any warnings remain visible when the update window is reopened.

## What is kept

- The latest fully verified automatic update backup, for rolling back the update.
- Older full data snapshots when the retained update recovery contains only
  programs, databases and settings; these copies provide different coverage.
- Scheduled and manual backups, plus all interrupted-state recovery archives.
- Unknown legacy snapshots. The last completed update journal can identify its
  automatic backup; other old copies cannot safely be distinguished from manual
  or scheduled backups and remain intact.
- Ordinary app data, user files, logs, source files and unregistered app runtime
  caches. An arbitrary folder is never deleted based only on its name.

## What is removed

After checking every file in the retained recovery copy, older automatic update
metadata/app recovery copies are verified and removed. Older full data snapshots
remain when the retained recovery excludes libraries. History entries for removed
copies are removed too. Only
backups belonging to this exact installation are eligible. Unknown added files,
modified copies, junctions, active executables or ownership errors prevent cleanup.

Downloaded update packages and their staging folders are removed when the normal
update operation returns. This release labels its staging folders so a later
successful update can remove leftovers from a crash. The current operation's
folder and unlabelled legacy staging folders are preserved by that sweep.

## Failures and interrupted cleanup

When an older updater installs this manager, the new manager waits until all
apps are ready, the older updater records completion and its maintenance lock
is released. It can then finish cleanup under the same lock. Reopening or
restarting does not repeat a completed cleanup. An older controller may leave
additional unlabelled history that cannot be classified safely; installing with
the newest installer provides the complete cleanup workflow and progress UI.

A failed update or rollback never triggers backup retention cleanup. Cleanup
starts only after the successful update has been committed to its journal.
It never deletes the final recovery copy to make room for an update.

Before removing an older backup, the updater marks it unavailable for rollback.
If Windows locks a file or cleanup is interrupted, its remaining files and
ownership marker are kept so the next successful update can retry. A partial
copy is never offered as a healthy restore choice. Cleanup warnings do not undo
an otherwise healthy update; inspect the **Update cleanup** details.

Scheduled backups remain independent of this retention policy. Reserve enough
space for the update download, the selected program/database/settings recovery
and a possible recovery archive before starting. Large databases such as mail
stores can still require substantial space. Ordinary files, registered libraries,
models and cached engines stay in place during update recovery. Full Backups app
snapshots still require room for data. Existing backups are not purged before
the new update and its backup have passed verification.

## Program and metadata recovery

Compatible independent app updates retain affected databases/settings and changed
programs. Workspace/shared updates retain broader metadata, still leaving bulk
files in place. Checked unchanged backup files can be reused. Neither update
recovery type replaces a full data backup, so cleanup preserves recorded full
copies. Recovery archives and manual/scheduled
copies remain excluded from automatic retirement. Cleanup counts approximate
physical disk space freed, including when backup files share unchanged bytes.
See [updates and rollback](UPDATES.md#backups-and-rollback) for scope and limits.
