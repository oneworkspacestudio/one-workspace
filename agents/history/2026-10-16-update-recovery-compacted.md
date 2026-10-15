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
