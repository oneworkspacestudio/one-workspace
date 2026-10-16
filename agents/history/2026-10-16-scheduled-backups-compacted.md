# Optional backups, storage layout and installation selection

## Scope and Git boundary

User requested an internal Backups app with an optional interval, switched off
initially, plus per-app storage on the selected Workspace drive and Advanced
installer app checkboxes. Implemented in `I:/Projects/one-workspace` on top of
`27ef326`; no new commit, push, release or live installation upgrade is authorized
by this feature request. Preserve the unrelated untracked installed-preview22
handoff. Public latest remains Preview 27. Build defaults now target unpublished
Preview 28; README distinguishes local development from published downloads.

## Implementation

- Owner-only `/backups` tile in home and the app switcher. Matching SVG and setup
  page include destination check/server-native picker, 1–168-hour interval (24
  initially), disabled scheduling, manual backup, copy progress/phases, last/next
  run, errors and verified history. POSTs require the configured origin and an
  owner session. Backups is internal, not a catalog executable or new repository.
- `backups.go` persists settings in protected `backup-settings.json`, status under
  excluded `maintenance/backups/`, and receipts in existing update backup history
  plus each external snapshot. Scheduler runs while Workspace is running, checks
  every 30 seconds and delays failed attempts 30 minutes. No automatic deletion.
- Reuses the update snapshot format: local ACL-capable NTFS, non-overlapping
  destination, no junctions/reparse points, EFS or directory ADS; copies Programs,
  Private and Files, verifies hashes and regular-file metadata/streams. Runtime
  scratch, logs, maintenance and prior snapshots are excluded. Capacity checks
  allow two copies plus512MiB. Recovery uses the existing local update/rollback UI.
- Backups serialize with updater/recovery via update.lock and manager operations.
  Preflight reads configuration and verifies installed programs before downtime.
  Process apps and supervisors use a separate job/context from the tunnel/portal.
  Snapshot pauses block app proxy requests with503; the owner portal stays alive.
  Whole-job termination waits for descendants before copying; owned shared service
  and listener checks follow. App restart/health checks run after success/failure.
  Shared initialization resets only after old app goroutines have ended.
- `selectedApps` supports explicit checked IDs; legacy Photos flag remains for older
  headless requests. UI recursively checks dependencies/unchecks dependents;
  MyAccounts stays required. Server validates against the downloaded catalog before
  creating installation state. Extraction uses only selected packages and required
  helpers, plus notices/licenses. Drive configuration/proxies/client are conditional.
  Coupled suite and standalone updates are scoped to installed app/package subsets.
  Docs editor assets are compiled into Drive; unchecking Docs removes its tile.
- Managed apps have private working directories and per-app TEMP/TMP on the chosen
  data drive; manager/shared helpers inherit protected Workspace scratch. Existing
  owned Music scratch is retained. Installer download staging uses the selected
  program drive. Existing storage bindings are not migrated; private content remains
  protected and is not exposed in shared Drive. Standalone apps retain their own
  settings. `docs/BACKUPS.md` records the layout and recovery limits.
- New local installer refuses a published bundle older than itself instead of
  downloading an older manager and losing these features. Until Preview28 is
  published, local candidate installation requires the explicit `--offline` option.

## Limits and remaining checks

No automated tests or runtime backup/restore/install acceptance were requested or
run for this task. Windows manager and installer compilation passed. Do not claim
end-to-end, clean-host, reboot or production acceptance.
No live account data, active suite, DNS, trust, autostart or GitHub release changed.

This is a full snapshot with application downtime, not incremental backup or a
bootable image. Unsaved work is not copied. External editors can make a copy fail
and trigger retry. Windows machine DPAPI still needs the original Windows security
context. Existing recovery requires a surviving installation/index; there is no
bare-disk restore wizard or automatic external receipt import. Keep independent
My Vault encrypted exports for cross-machine recovery. Existing component preview
limitations and organization pins remain unchanged.

Next acceptance, if authorized: fresh isolated NTFS fixtures only; check minimal
and full selections, dependency toggles, no omitted executables after updates,
owner/origin gates, default-off settings, scheduling/restart, unavailable/full/
overlapping destinations, copy interruption, shared-service pause/restart and
actual data-preserving rollback from a scheduled snapshot. Never use the active
installation. No task-owned app or preview services were started.

## Final build and storage incident

The initial full org build completed; a later build with final adjustments failed
at the Windows linker with FlushViewOfFile / I/O device error while writing on I:.
Source changes were preserved and hash compared in the named NTFS fallback
`C:/Projects/.cache/workspace-backups-build`. Windows reported the I: exFAT volume
as Warning while its USB physical disk reported Healthy/Online. System storage
reset events also existed, but no particular disk failure was diagnosed. No disk
repair or changes to user data were performed.

Final source was copied to the isolated NTFS root and manager/installer compiled
there. Existing compressed org components were retained without recompression;
only the manager and catalog were replaced. Every one of59 bundled package file
sizes and SHA256 digests matched the final catalog (13 app entries), including
the rebuilt manager. No GitHub Actions or app processes were used. These are
packaging checks, not runtime or restore acceptance. Build defaults are28; public
Preview27 assets remain untouched.

Static review also fixed initial configuration when Drive is added later through
`install-apps`, guarded against older bundle downgrades, and made online cache-file
rotation tolerable during preflight size estimation (the stopped tree is sized
again before copying). Backups cache the verified update controller before pause
so a fresh installation can use existing recovery tools without a prior update.

Final artifacts, source snapshot and build/package receipts are verified under
`tooling/workspace-backups/final-preview28`. Installer674937856bytes, SHA256
`4750688f8580e1eb8ff3a18582d214a3a40731552e8df022b8dd72823881973a`.
All final artifacts were flushed to I: and read back with matching hashes. The four
core dist outputs, checksums and embedded payload were aligned to this final build.
Other old release metadata in dist was not regenerated: reconcile source/tag,
build-info, limitations and new guide before any authorized publication.

A second small source archive/receipt is retained under the calling workspace's
`builds/backups-preview28` on C: because of the reported I/O incident. Removed only
the eight exact public component extraction stages created by the two builds and
the isolated NTFS fallback root after source/artifact preservation, canonical path,
reparse and active-process checks. No fixture/app process was started. C: free space
returned to approximately18.1GB. Existing component caches, live apps, private state
and rollback backups are untouched. Follow-up history/receipt bytes do not change
the compiled Go inputs or the recorded installer checksum.

## Publication authorized

User now explicitly requested pushing and releasing. Publish the Workspace changes
as v0.3.0-preview.28 with the complete13-app bundle, using unchanged component pins.
No separate app binaries/source changed. No Actions, runtime tests or live upgrade.
README and backup/update guides now identify28; release status explicitly records
pending runtime backup/restore acceptance. Refresh the bundled release-status
notice, installer, exact committed-source ZIP, provenance and checksums before
publication. Leave the unrelated untracked installed-preview22 handoff untouched.
