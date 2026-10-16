# Workspace update and recovery implementation

User requested updates for existing Workspace installations and apps without data
loss, followed by pushing and releasing. Use the authoritative one-workspace source;
do not modify the active installation. A concurrent installation handoff confirms
its current preview22 state; that unrelated untracked handoff is preserved.

New local maintenance window is accessible from setup's existing-installation
entry, Setup --update PRIVATE, or owner administration after bootstrap. Local
capability, strict origin/host checks and owner portal gates protect operations.
Workspace and per-app compatible releases can be selected; shared backends update
together. Independent app releases declare workspace-update.json. Missing contracts
are shown, not guessed. cmd/update-manifest builds contracts with owned packages,
helpers and notices. Stale plans, altered installed bytes, data-path/backend/OIDC
callback changes and the withheld Music digest are refused.

Downloads stage before downtime. The stopped installation receives a verified full
backup of Programs, private state and Files on local ACL-protected NTFS storage.
Existing backup/cache/log roots are excluded. Identity clients, paths, owner and
custom/disabled apps are retained. Regular file streams, timestamps/attributes and
ordinary-file permissions are preserved. Junctions/reparse points, EFS and directory
ADS are refused. Disk capacity is checked for backup plus recovery archive.

Journal and cached verified controller prevent mixed legacy startup: update manager
first; guard recovery until data/programs are restored and legacy managers last.
Failed startup restores the original snapshot, archiving current state first.
Manual rollback similarly archives post-update edits. Partial/failed archives are
not offered as complete healthy installations. Recovery can locate its backup even
when workspace.json is unavailable. Application backups are not volume images.

Unit checks cover independent release scoping/selection/staging, tampered downloads
and backups, legacy version labels, stale plans, identity/content preservation,
failed health with changed data, interrupted restore, capabilities/owner/origin,
junctions, storage-binding refusal and Windows regular-file ADS/metadata. Native
preview22 fixture includes real Notes/attachment, private video and documents.
Its first upgrade and rollback passed; final source snapshot validation follows.
A mistakenly overlapping rebuild was stopped without affecting the test fixture.
Use tooling/workspace-updater23 and the named C: cache only; track and stop each
owned preview/shared engine. Preserve active live data and unrelated handoffs.

Build locally; do not use Workspace Actions. Retain all existing component pins,
Notes4 and the approved earlier Music binary. New complete release is planned as
v0.3.0-preview.23. Do not rewrite published tags/assets or older draft previews.

The first mutable development run later stopped during a second restore because
mixed build generations lacked a cached controller. This was a test-artifact
coordination error, not a published release. Verified that the edited note and
ordinary document still existed in the protected post-update archive. All fixture
processes were stopped; the synthetic fixture and its registration were removed
after preservation, path/reparse and process checks. Repeat the complete run from
an immutable source export and do not replace artifacts while their test runs.

## Immutable acceptance completed

Code commit `b90cd366027bc60d161343698552fb174ff96304` passed the full Go test
suite and vet, then was exported before building or running final acceptance.
Native Windows preview22 -> preview23 upgrade, rollback to22, restoration of
post-update edits, and automatic rollback from an invalid Notes executable all
passed. Every stage returned all twelve apps healthy. Real MyAccounts login and
stable Notes/My Tube user IDs, note body, attachment, private video bytes and
ordinary documents were verified. Added ordinary files survived rollback, and
newer note/document edits were recoverable. Five full/recovery archives remained.
All fixture supervisor, app, helper and shared-engine processes were stopped.

The packaged full installer separately opened the existing fixture and its
retained backups, exposed the first-screen existing-installation entry, and
completed an anonymous per-app release metadata check. The release manager SHA256
is `e645093cb3966e5893128df03ac1d0d14dc4a2323b6893a198cd245012eadd63`, identical
to the fixed executable used throughout native acceptance. Every bundled byte,
package size, twelve-app catalog and upstream release asset digest was verified.
No clean-host, real reboot/power-loss or Windows-service acceptance is claimed.

## Published release

Pushed the implementation and immutable source tag, then published
[v0.3.0-preview.23](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.23)
from `b90cd366027bc60d161343698552fb174ff96304`. All 15 assets passed GitHub
size/SHA256 comparison before publication. Local-only builds; Actions remained
disabled. Both installation and update publishing guides accompany the complete
installer, manager, twelve-app bundle/catalog, source/provenance, licenses and
checksum manifest. Existing published assets and component release pins were
preserved. The active installation and its live rollback backup were untouched.

The released full installer then anonymously downloaded and verified the newest
complete public preview23 bundle and reported all twelve app versions correctly.
The task-owned native installation, all rollback snapshots and private fixture
state were encrypted and read back with per-file hash verification before cleanup.
The archive is under local `tooling/workspace-updater23`, with its key wrapped by
current-user Windows DPAPI; it is not a public release asset. Temporary C: cache,
owned native Apps/Files and the single fixture registration were removed only
after exact-path, process and reparse checks. No fixture services remain.

Verification receipts, fixed source export, release staging, screenshot and
publication receipt are retained under local `tooling/workspace-updater23`.
Implementation/publication is complete; this task did not upgrade the active
preview22 installation. Existing users bootstrap through the new installer's
**Already have One Workspace?** entry, then use owner **Updates and backups**.
Future standalone app releases must publish the documented compatibility contract
before they can be installed independently. No app repository release was changed
by this task. Preserve the unrelated untracked installed-preview22 handoff.
