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

## Authorized release and installation follow-up

User now explicitly requested release and install. Source committed/rebased onto
org main74a42be, preserving MyAI9 promotion. Final code sourcea4de985 includes a
preview91 floor for schema3 recovery controllers, caching the staged compatible
manager during90->91, and actual changed-asset download previews (JSON d.asset,
unchanged coupled downloads removed). All source pushed; no Actions/tests used.
Earlier local/untracked/nonpublication statements above are historical.

Public immutable releases:
- Workspace91, tagv0.3.0-preview.91, sourcea4de985cf1f1564a7fa75756b1e1419a46e55d44.
  Full Windows installer1,093,954,048bytes, bundle1,080,601,520bytes, manager16.6MB.
  All28 apps/166 payload hashes checked. Offline defaults overlay public90 with
  Drive5/Photos6, Tube5 and MyAI9; normal updates retain independent app choices.
- Drive/Photos6, sourceee090cbd389191975ffc8f24c6dbb364c923f9e3. Drive executable
  unchanged from5; Photos/frontend rebuilt with Go1.27.3/Vite. Photos EXE SHA256
  9ffc67cd8d133658c7866a727b9c9750aec62c0aee22f5474754bc548bb0e759.

Every published asset size/digest matches local. Promotionsa1e9f83/38176a1 add
Photos6 then Workspace91 and concurrently public MyAI10; old releases and
withholdings unchanged. Anonymous binaries/catalog/manifest/raw channel checks
pass. Anonymous normal updater offers Workspace91 as16,652,166bytes in3 direct
assets, without downloading the full suite bundle. MyAI10 is available but not
selected in this rollout; Ai chat is holding installs under human-authorized
coordination until both PCs finish.

Main PC normal verified offline --apps workspace,drive apply completed91 and
Photos6. Snapshot schema3:782,648,768bytes logical,643,203,551bytes reused,
139,445,217bytes newly copied. Recorded Files and all registered preserved paths
are excluded. Independently audited original protected backup vs current config,
credentials in memory, app selections/locations/identities and all other app
contracts/packages; unchanged. MyAI8 retained. Private Photos.log exists. All27
services run locally and through public Chrome healthz. Initial receipt helper's
anonymous urllib public health request returned403 after the successful update;
Chrome public health was OK. Completion receipt was reconstructed with the full
before/after audit from the verified protected backup, not by rerunning apply.

Secondary currently downloading/verifying the public91 installer via its normal
remote Reinstall flow. This bootstraps the new updater without old90's8GB data
copy; repair retains its installed app versions. MyAI10/Drive checkboxes were
unchecked for this manager-only step. After91, apply only Drive/Photos6, retain
MyAI9 and every other row. Still pending; do not claim both-PC completion yet.

Artifacts/logs/receipts on I:/Projects/tooling/update-efficiency-20261026. Build
retries removed their own embedded old payload before manager compilation; final
build/cache/scratch on I:. Owned empty C build scratch removed after process/path
checks; no live/backup folders deleted. Full old data snapshots remain protected.
No native rollback, clean-host/reboot/Linux or production acceptance claimed.

### Patch92 required by installed selection

Secondary manager91 repair completed;27 public services running, MyAI9/other
versions retained. Reloaded updater shows new metadata-only descriptions. Before
applying Photos6, its real plan exposed a download regression: a component equal
to the suite pin skipped its matching standalone manifest, so Photos selected the
1.08GB suite ZIP instead of7.8MB Photos.exe. No such apply was submitted; all
checkboxes cleared. Both PCs remain healthy91, main Photos6/secondary Photos4.

Preparing92: resolve matching pinned-tag managed manifests too, replacing only
download mappings whose package path/size/SHA256 equal the original suite target.
Do not replace the suite app contract for that case. Legacy/no-match fallback is
still the verified bundle. Protocol4 refreshes idle91 controllers so the resolver
change reaches remote owners. Metadata policy remains v1/floor91, preserving
91 recovery compatibility.91 public assets/tags remain immutable;92 must be a
new full release. Still no tests/rollback fixture authorized/run.

Ai peer may carry out its independently human-authorized MyAI10-only operations
during this build window using installed91 controllers. Asked it to report both
completion before92 applies and avoid channel/pin changes. Preserve actual dynamic
versions in92; do not assume8/9 after peer updates. No local/remote operation from
this task is active while92 builds. Need publish/promote92, manager-only updates
on both, then secondary Photos-only plan must show the standalone7.8MB file.

### Workspace92 publication and primary installation

Frozen release source b45a6d6a60913a82c735346a5ab0acf7a5f981bb is pushed/tagged
v0.3.0-preview.92. Matching-pin download mappings are used only for changed
package bytes, validated against exact original target path/size/hash; unchanged
components do not fetch extra manifests. Windows manager/setup compile and all
166 payload hashes pass. Public full installer1,093,956,608bytes; every published
asset matches size/SHA256. Actions disabled, no tests/rollback fixtures run.
Public release https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.92
Channel promotion7653baf adds only92; every previous release and withholding is
unchanged. Anonymous Photos/manifest/manager92/catalog/channel checks pass.
Normal anonymous manager-only plan downloads16,646,073bytes in2 direct assets,
never the suite ZIP.

Secondary91 repair completed with27 services healthy. The independently
human-authorized Ai task subsequently installed MyAI10 on BOTH PCs and confirmed
normal completion before our92 applies. Earlier8/9 retained/pending notes above
are historical. Its public exe hash is
b0cad46690fcecc023ee7ef998f5833baf78a24bbd732525175bba2ba46a334b.

Primary normal online manager-only92 apply completed. Independent protected
before/after audit confirms owner, credentials (decrypted only in memory), paths,
app selections, all app contracts/other package bytes preserved; Photos6/MyAI10
retained. Schema3 recovery203,867,792bytes logical,47,561,183 reused; Files and
registered preserved libraries excluded. Cleanup freed~52.4MB from1 older
eligible update copy; latest recovery and full data snapshots preserved. All27
services healthy locally and via Chrome public URL. Private Photos.log confirmed.
Secondary manager-only92 apply is now running with16.6MB selected; Photos remains
unchecked until idle92 controller produces its standalone7.8MB plan.

### Both-PC rollout completed

Secondary normal manager92 apply completed and refreshed controller4. Other
installed versions retained, including MyAI10. Its Photos-only plan then selected
exactly1 direct Photos.exe asset,7,780,352bytes (previously1.08GB suite ZIP).
Normal Photos6 apply completed in3m17s. Recovery238.1MB checked,104.5MB reused,
133.5MB newly copied; existing user libraries stayed in place. Cleanup freed~154MB
from1 older eligible update copy, keeping newest recovery and old full snapshots.
All27 services healthy via secondary public Chrome healthz after completion;
primary all27 public/local health and independent protected snapshot audit pass.
No second apply was sent during temporary progress loss/restart502s. Source92
fix and public receipts are immutable; no further binary release needed for this
handoff. Artifacts/receipts/proof on I:/Projects/tooling/update-efficiency-20261026.
C scratch removed;~15.9GB free after completion. No automated tests, actual native
rollback, Linux/clean-host/reboot or production acceptance is claimed.
