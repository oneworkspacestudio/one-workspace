# App update efficiency

## Request and baseline

User asked why a single-app update can involve8GB and requested efficiency.
Baseline9ce5506 in the existing I-drive worktree. Loaded repository-history,
AGENTS and the update/recovery/cleanup handoffs. Confirmed in source: app downloads
already select/deduplicate changed release assets, but every apply previously
made/verified a complete Programs+Private+Files snapshot. Tube5 public Windows
ZIP was6,534,069bytes; the secondary combined Tube5/AI8 operation copied8.8GB
of recovery data and took23m51s. Do not conflate those measurements.

## Local implementation

Work is uncommitted on codex/app-update-efficiency. No current-task authorization
to push, release or deploy; neither owner's live installation was changed.

- AppUpdateScope/backup schema2: existing independent known app groups, changed
  owned packages, their registered private roots, workspace.json/secrets/manager
  recovery launcher. Includes known sidecar storage for AI/Clips/Social. Does not
  traverse unrelated app libraries or ordinary Files. Identity/Drive/Git/Music/
  shared engines, app additions/removals, legacy selection migration, changed
  contracts/paths/manager or unknown layouts retain full schema1 snapshots.
- Scope derived from manager registry, never publisher-supplied arbitrary paths.
  Validate unchanged contracts, package ownership/paths, catalog metadata and
  stored scope/receipt; refuse unexpected scopes, junctions, special files,
  custom storage references and missing original package bytes. Notes/Tube
  custom settings paths outside their known roots force full snapshots.
- Scoped failure/manual recovery archives only matching affected state, clears
  only registered scope paths after verification and restores old binaries/data.
  Manual app rollback requires the exact expected installed catalog so it cannot
  silently rewind other later app versions. Full/manual/scheduled paths remain
  full snapshots; older controllers refuse schema2. Remote protocol3 refreshes
  idle legacy controllers, preserving busy operations.
- Large unchanged regular files may be reused via same-volume backup-to-backup
  hard links after live/prior hashes and no-ADS checks. Never link live data.
  Verify through the new name; unsupported/corrupt/cross-volume candidates copy
  normally. No delta chain or prior path dependency. Per-receipt checked/reused
  bytes and progress distinguish data checked from newly copied. This saves
  writes/allocation; large affected files are still read/hashed, and worst-case
  free-space reserve for recovery is retained.
- Cleanup of scoped updates preserves existing full automatic snapshots, since
  an app-only copy cannot replace full coverage. A later healthy full update can
  retire older automatic full/app snapshots. Approximate reclaimed space uses
  volume free-space change, avoiding counting shared links as freed data.
- Updater selection shows distinct asset download MB/file count separately from
  recovery data, and backup choices identify app vs full recovery. UPDATES and
  UPDATE-CLEANUP describe semantics/limits.

## Verification / limits

Windows x64 manager and setup build with Go1.27.3, CGO0; changed inline JS passes
Node syntax checking; git whitespace check clean. No new tests were added and
no test suite or native runtime/update/rollback fixture was run (session says
not to run tests unless the user asks). Do not claim acceptance/performance or
production readiness. Native scoped success/failure/recovery/cleanup/reuse and
old-controller/schema compatibility checks are still needed before deployment.
An attempted Linux cross-build is blocked by existing Windows-only Protect/
Unprotect implementations in platform and calendar/shared/secure. This was not
introduced/fixed by this task; no Linux build/runtime success claim.

Build output/scratch only I:/Projects/tooling/update-efficiency-20261026; no C
fixture/services created, no live backup deleted. No GitHub Actions or source
publication. Ordinary full snapshots remain unchanged outside scoped apply.

## Coordination and resume

Under the previously verified human coordination authorization, told the Ai chat
that this is source-only and no deployment will run. It is independently delivering
MyAI9. Fetch observed origin/main f2f3bf6 (MyAI9 promotion) ahead of local base;
its pin/index changes are not in this dirty branch. Preserve them and rebase onto
current main when a later authorized commit/publication is prepared. Do not
publish the old local channel/index over that promotion or start a simultaneous
installation. Both scope/reuse source files are new/untracked, so include them
with intended tracked changes when committing is authorized.

## Follow-up: separate update recovery from data backups

User approved saving replaced programs plus affected databases/settings, leaving
bulk data to the separate Backups app, conditional on program/data separation.
The design below supersedes the earlier full-copy fallback described above.
Still source-only; no new publication or installation authorization this turn.

- New schema3 `UpdateMetadataPolicy` v1 records replaced old/new program paths,
  original target catalog and expected installed catalog. Independent app scope
  remains optional; Workspace/shared/addition/repair changes now save broader
  private metadata, not full libraries. Snapshot verifier uses stored catalogs,
  not live app settings, so damaged current settings cannot prevent recovery.
- Fixed registered exclusions leave Files and media/project/model libraries,
  Gitea repositories/LFS, Vault attachments, hosting projects, Tube helpers/media,
  Slicer engines/models/slices, Cinema engines, AI document/project stores and
  downloaded engines/caches in place during both snapshot and restore. Never
  exclude whole DB stores: Mail/shared DBs and Vault engine DB/keys remain copied,
  along with WAL/SHM, settings and credentials. Embedded DB content may be large.
- Metadata restore gathers checked paths, skips preserved subtrees and removes
  entries individually; no recursive deletion of an ancestor of preserved user
  files. New files in those libraries survive rollback. Conflicting ancestor
  directory/file replacement is refused. Archives match the original policy and
  retain the original target paths even if a failed apply never saved its catalog.
  Healthy manual rollback requires exact expected catalog. Schema1/2 retain their
  original full-data/app-data semantics and are not silently reinterpreted.
- Startup/update planning reject overlapping Apps/Private/Files and known args or
  env data paths under Apps (including relative paths resolved from app cwd).
  Notes/Tube settings are checked before stop and after stop; unsupported custom
  DB/media roots or unknown app profiles refuse automatic update recovery. New
  program destinations refuse existing unowned files. No blind migration/deletion
  of unknown existing state. Per-app cwd is private, including hosting/slicer.
  Managed children receive WORKSPACE_APP_DATA_DIR in their private folder.
- Read-only primary program inventory found old downloaded Gitea runtime assets
  and an active Photos.log beside executables; no app database/media store was
  identified there. This is not a universal content audit. Companion source fix
  in I:/Projects/filebrowser-drive moves Photos logs/state to injected private
  location, with a standalone user-config fallback and read-only legacy config
  compatibility. Old installed Photos continues writing its old log until an
  authorized deployment; no live files were moved or deleted.
- Separate Backups manual/scheduled full snapshots are unchanged; optional/off
  by default. External vendor model libraries are not implicitly covered. UI/docs
  distinguish metadata recovery, full data backup and download size. Cleanup keeps
  old full data snapshots because metadata recovery cannot replace their coverage;
  only older eligible automatic app/metadata copies are retired.

Windows x64 manager/setup and companion Photos production compilation pass with
Go1.27.3/CGO0. Maintenance JS syntax and Git whitespace checks pass. Builds/scratch
stay on I:. No tests, runtime/update/rollback fixtures or installation were run;
native recovery/compatibility acceptance is still required before release. No
commits, push, tag, release, channel change, Actions or live backup cleanup.
Preserve concurrent public channel work described above when publication is
later authorized. All three new updater Go files remain untracked and required.
