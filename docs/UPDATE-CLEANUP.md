# Automatic update cleanup

Workspace preview.39 cleans up after the updated apps restart successfully.
The updater shows **Clean up old update files**, then reports how many older
copies were removed and approximately how much space was freed. The result and
any warnings remain visible when the update window is reopened.

## What is kept

- The latest fully verified automatic update backup, for rolling back the update.
- Scheduled and manual backups, plus all interrupted-state recovery archives.
- Unknown legacy snapshots. The last completed update journal can identify its
  automatic backup; other old copies cannot safely be distinguished from manual
  or scheduled backups and remain intact.
- Ordinary app data, user files, logs, source files and unregistered app runtime
  caches. An arbitrary folder is never deleted based only on its name.

## What is removed

After checking every file in the retained recovery copy, older automatic update
backups are verified and removed. Their history entries are removed too. Only
backups belonging to this exact installation are eligible. Unknown added files,
modified copies, junctions, active executables or ownership errors prevent cleanup.

Downloaded update packages and their staging folders are removed when the normal
update operation returns. This release labels its staging folders so a later
successful update can remove leftovers from a crash. The current operation's
folder and unlabelled legacy staging folders are preserved by that sweep.

## Failures and interrupted cleanup

A failed update or rollback never triggers backup retention cleanup. Cleanup
starts only after the successful update has been committed to its journal.
It never deletes the final recovery copy to make room for an update.

Before removing an older backup, the updater marks it unavailable for rollback.
If Windows locks a file or cleanup is interrupted, its remaining files and
ownership marker are kept so the next successful update can retry. A partial
copy is never offered as a healthy restore choice. Cleanup warnings do not undo
an otherwise healthy update; inspect the **Update cleanup** details.

Scheduled backups remain independent of this retention policy. Reserve enough
space for the update download, a new full recovery copy and working space before
starting. Keeping one recovery copy can still require substantial space when your
mailboxes and ordinary files are large. Existing backups are not purged before
the new update and its backup have passed verification.
