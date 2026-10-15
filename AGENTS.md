# Workspace repository

Use repository-history before substantive changes and maintain agents/history.
Never test against the existing live MyAccounts/FileBrowser installations.
Test new suite installations in fresh NTFS private-data folders; program payloads
may use another drive. Preserve component license/provenance and user files.
Do not claim a production or clean-host/service/reboot test not performed.
The bundle is an evaluation build while upstream acceptance gates remain open.

## Build and preview cleanup

Use I:/Projects/tooling for large build artifacts, downloads and cleanup archives.
When a build needs NTFS, use one named C:/Projects/.cache build folder per app
and remove its disposable contents after the task. Preserve source changes first;
never treat the live installation or unknown private data as a disposable cache.

Track each test installation's program, data and file roots and all child PIDs.
Stop task-owned preview services and tunnels when verification finishes. Verify
that they have stopped; do not leave older preview suites running alongside the
user's active installation. Check C: free space before large builds/install tests.

Before recursive cleanup, verify absolute paths and active process references.
Do not follow junctions into other folders. Archive unknown test state or older
rollback data before removal, verify the archive, and encrypt private archives
when the destination cannot preserve Windows ACLs. Keep the current installation
and newest usable rollback backup. Do not clear unrelated projects or global temp.
