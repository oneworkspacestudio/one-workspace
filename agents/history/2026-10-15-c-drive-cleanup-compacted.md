# C: duplicate test/build cleanup

User explicitly requested cleanup of repeated test installations on C:. This
supersedes the earlier boundary preserving those obsolete previews, but does not
authorize publication or changing the active installation's data/storage roots.
Active installation remains C:/tmp/workspace. The preceding installed-flow audit
and its unresolved Mail/integration findings remain separate; no app fixes here.

## Completed safeguards and process cleanup

Inventory covered the task-related children of C:/Projects/.cache, two empty
C:/tmp installer checks, old active-install rollback snapshots and old global
Music extraction directories. Unrelated C:/Projects projects, legacy F: apps,
system caches and active Apps/Files/Private state are excluded.

Stopped obsolete local Workspace suites on19390/19590 using their authenticated
stop command, then verified their children exited. Stopped separate old PublicMail
and its preview tunnel, Calendar/Contacts demo processes on18998/18999. Five stale
WorkspaceSetup windows were later closed because one held a test log open. The
current manager/children and public tunnel remain running. No cache-referencing
Windows services or scheduled tasks were found.

Eight global _MEI directories identified as Music payloads were checked against
loaded modules and removed only after the old Music processes exited. Current
Music remains in its dedicated private runtime-temp root; unrelated contents in
that root were preserved. Manager19 already bounds its own extraction cache.

## Archive policy and current operation

Task artifacts outside Git:
I:/Projects/tooling/workspace-cleanup-20261015.
Archive-And-Clean.ps1 inventories and archives each fixed allowed target, refuses
active references/reparse ancestors, skips following child links, verifies an
AES encryption/decryption SHA256 round trip, seals recovery metadata with Windows
CurrentUser DPAPI, checks source file length/mtime stability, then removes only
the verified source. Unknown test/source state is retained encrypted on exFAT I:;
no plaintext private archive is written there. Recovery helper/README supplied.
Recovery requires the same Windows account/profile and protected NTFS output.
Archives preserve file contents/empty directories; links are recorded for manual
recreation, not traversed; Windows ACLs need appropriate reapplication on restore.

First archive was initially slow due small output writes; an unfinished archive
was removed after stopping its own worker. Buffered output was added. First full
Office build archive verified/removed6.76GiB. Next Mail attempt refused a locked
installer log without deleting source; stale installer windows were closed, that
unfinished archive removed, and operation resumed. Later output suppresses ZIP
entry flushes to the external drive while flushing the completed archive explicitly.
completed.json retains every finalized archive/removal across resumed runs.

Initial C: free space approximately0.53GiB. Interim verification reached20.58GiB
free with11 archives totaling15.01GiB of original files removed. Cleanup remains
in progress at this checkpoint; native exec session43612 is processing Music build
(20,000/28,761 files at last output). Do not claim final completion yet.

21 loose test logs/downloads (~9.9MiB) were moved within the validated cache root
to workspace-cleanup-loose-20261015 and appended to inventory.json. The running
worker loaded the old list, so run the helper once more after it exits to archive
that one staged folder. Removed prior sources will be skipped. All folders are
fixed-boundary targets; do not blindly reuse the cleanup helper for new state.

## Remaining completion

Wait for worker43612, resolving only confirmed old-test locks if encountered.
It should archive remaining cache folders, two empty installer checks and older
rollback snapshots. Keep Private/backups/install-apps-20261015T160228.753290200.
Then rerun once for the staged loose files, verify completed.json and all paired
archives/DPAPI keys, actual C: free space, no remaining preview references, current
local/public health, preserved Apps/Files/Private and retained newest rollback.
Update this handoff with final measured results and precise remaining exceptions.

AGENTS.md now instructs future work to keep large artifacts on I:, use bounded
NTFS build caches when necessary, stop task-owned previews/tunnels at completion,
verify process/path ownership before cleanup and encrypt archived private state
when destination ACLs are unavailable. No commit/push/release is authorized here.

## Final completion — supersedes the in-progress checkpoints above

Worker43612 completed successfully. Reran the helper with explicit OnlyTarget
for the staged loose folder; no rescan of new/unknown state. All52 intended targets
were archived/verified/removed:37 original cache folders, one staged loose-file
folder, two empty installer checks and12 older rollback snapshots.180,172 original
files total18.18GiB; separate eight unused Music extractions removed. All52 paired
DPAPI recovery metadata files decrypt and match their archive/source records;
sources are absent. Every archive had already passed a plaintext SHA256 round
trip before deletion. The detailed report/records are outside Git in the task root.

Final C: free25,846,857,728 bytes =24.07GiB, up from approximately0.53GiB at start.
The measured free-space delta also reflects stopped-process/system activity and
must not be confused with the18.18GiB archived-file total. C:/Projects/.cache is
empty; zero process executable/argument references remain. Keep newest full backup
Private/backups/install-apps-20261015T160228.753290200; it is the only backup folder
remaining. Current Apps/Files/Private state was retained and not migrated.

Both localhost8990 and the public Workspace health endpoint report statusok with
all child services running, including the current public tunnel. No restart of
the active manager was needed. No product implementation/tests or publication
changes were made; only AGENTS cleanup guidance and this handoff changed in the
repository. The prior audit handoff remains untracked and untouched. No required
cleanup work remains. The earlier browser end-to-end audit is still incomplete
for its separately documented tool-policy blocker; cleanup does not establish
UI editing/saving/delivery acceptance or fix those Mail/integration defects.
