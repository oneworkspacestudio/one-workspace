# Photos upload recovery and duplicate naming

User reported a Photos upload error containing raw Cloudflare502 HTML. The
screenshot was from07:55UTC during earlier recovery; user retried and confirmed
uploads now work. Later reported409 and explicitly requested `(xx)` name suffixes.
See filebrowser-drive's matching history for code/build details and constraints.

Drive/Photos frontend fix normalises gateway errors and adds Photos in-memory
retry queue, confirmed-chunk state, XHR progress/Saving, destination retention
and success summary. Duplicate uploads choose `photo (1).png`, `photo (2).png`
automatically, checking names before bytes and handling racing409 without
overwrite. Unknown final replies may create numbered copies on retry; documented.

Actual installation C:/tmp/workspace was updated locally: manager/installer
0.3.0-photos-upload-preview.10, Drive/Photos photos-upload-preview.2. Existing
Mail0.1.0-subdomain-preview.8 retained (inbox.m8e.co.uk); root mail DNS unchanged.
Normal control stop with exact path checks; changed binaries/workspace.json
backed up under Private/backups/photos-upload-fix-20261015. Updated catalog
hashes/sizes/versions; ordinary files and identity/Mail data unchanged.

Build source isolated from concurrent Drive sidebar/design-system edits. Those
unrelated source edits were preserved and excluded from this deployed build.
Suite staging C:/Projects/.cache/photos-upload-suite uses the previously deployed
Mail source and new paired app inputs. Full Office/Studio embeds retained.
Build checks: Vite production and Go compilation plus JS syntax/whitespace.
No tests added/run or agent uploads performed. User's successful upload was on
the preceding runtime; new duplicate/resume behavior has no functional acceptance
claim. No commit/push/release/Actions authorized.

Local artifacts: dist/WorkspaceSetup-photos-upload-preview.exe,
Workspace-photos-upload-preview.exe, SHA256SUMS-photos-upload-preview.txt.
Read actual process PID from Private/running.json. Refresh Photos after update.

Deployment readback:8/8 startup, all apps running, Photos version.2 and new
retry/numbered-name text served from its executable. Mail domain remains inbox
with no hosting/identity error. This is configuration/asset-readback evidence,
not functional duplicate-upload acceptance. No user files uploaded by agent.
