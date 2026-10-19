# Automatic cleanup of verified updates

User explicitly requested implementing update cleanup, push and release. Build
Windows binaries locally; GitHub Actions stays disabled. Work started from
f4e29d7 (main-equivalent) on codex/four-apps. Current owner live installation was
restored in a separate availability task; leave it running and unchanged here.
Preserve the three unrelated untracked availability/history files.

## Implementation

- Successful committed updates retain their newest fully verified rollback and
  remove older positively identified automatic update backups. Add purpose to
  backup manifest/receipt; Backups app snapshots are separate. Recovery copies,
  unknown legacy backups and arbitrary app runtime caches remain intact.
- The prior completed update journal can identify its legacy automatic backup;
  this proof is carried into the new journal. Other unlabelled snapshots cannot
  be distinguished from scheduled/manual backups and are preserved.
- Fully verify retained copy and deletion candidates; enforce exact installation
  roots/owner, canonical boundaries, checksum/purpose/time, no unrecorded files,
  no reparse points and no executing programs in cleanup trees. Delete individual
  entries after fresh path checks; metadata last. No broad directory-name purge.
- Mark old backup cleanupPending in the protected index before deletion. Partial
  deletion is never advertised as a healthy rollback; next successful update can
  retry with its owner marker. Failed updates/rollbacks never run retention.
- New owned update stage markers enable removal of crash leftovers after a
  later healthy update. Active stage excluded from sweep; normal return removes
  current downloads/staging. Unknown old staging folders remain intact.
- Persist cleanup counts/approximate bytes/warnings; update window has a cleanup
  progress step and result card. Cleanup failures do not roll back healthy apps.
  README/UPDATES and new UPDATE-CLEANUP guide explain retention and limits.

## Verification and release preparation

Compatibility follow-up: Run starts a post-readiness cleanup watcher. It waits
for a completed journal and the older updater's maintenance lock, rechecks the
installed owner/version, and skips an already completed cleanup receipt. This
lets the new manager perform cleanup when installed by an older controller.
The completed journal proves an unlabelled retained legacy copy is automatic;
unknown older copies remain protected. Handoff/idempotency regression passes.
Final native assembly replaces the manager/catalog in the verified full bundle
using Go archive/zip.Copy for unchanged members, then recompiles the installer;
all member sizes/digests are rechecked. Component pins remain unchanged.

Windows cleanup regression cases passed: automatic vs manual/scheduled/legacy/
recovery retention, incomplete commit, corrupted retained copy, locked-file
partial deletion/retry/rollback denial, proven legacy cleanup, unexpected files,
foreign paths, junction target preservation, successful vs failed startup.
The shared junction test helper now converts forward slashes for Windows cmd;
without that it silently skipped tests with Invalid switch Projects. Real junction
creation and both cleanup/older update-junction regressions now pass.
Suite/maintenance/platform/distribution tests and vet passed before final release
packaging. Initial test-fixture failures (read-denying lock, missing stage root)
were fixed; they did not represent production deletion of data.

Next: local preview.39 full19-app release build with unchanged component pins,
isolated native legacy-to-new update/retention check and update-window review,
then verified GitHub assets/publication. Do not label production/clean-host/reboot/
independent security acceptance complete. External scripts/evidence/stages live
under I:/Projects/tooling/update-cleanup-release; NTFS disposable state uses only
C:/Projects/.cache/workspace-cleanup-20261019. Track and stop all fixture PIDs;
clean only that verified non-reparse named root after testing. Current owner
C:/tmp/workspace/{Apps,Private,Files} must remain untouched.
