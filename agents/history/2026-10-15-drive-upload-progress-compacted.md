# Drive upload progress repair

User reported Drive uploads at100% offering Resume and random mid-batch pauses.
See filebrowser-drive's matching history for diagnosis and source changes.
Implemented acknowledgement-based completion, separate transfer/saving deadlines,
visible saving/pause explanations and safe pause/resume scheduling. Screenshot
does not prove a specific original network failure. No upload tests added/run.

Local installer/manager0.3.0-drive-upload-preview.11 contains Drive
drive-upload-preview.3, unchanged Photos photos-upload-preview.2 and Mail
0.1.0-subdomain-preview.8. build.ps1 now optionally accepts PhotosVersion so a
Drive-only update doesn't mislabel the unchanged companion. Default still falls
back to DriveVersion for existing callers. Build isolated from concurrent Drive
styling work; preserves prior Mail setup/progress and Photos duplicate fixes.

Compiled Vite frontend/checker, JS syntax and Go Windows binaries; no functional
upload acceptance claimed. No commit, push, release or Actions authorized.
Ignored artifacts dist/WorkspaceSetup-drive-upload-preview.exe,
Workspace-drive-upload-preview.exe, SHA256SUMS-drive-upload-preview.txt.
Packaging source C:/Projects/.cache/photos-upload-suite; core source separately
C:/Projects/.cache/photos-upload-build.

Actual C:/tmp/workspace updated with normal verified manager stop; prior changed
executables/config saved under Private/backups/drive-upload-fix-20261015-100853
(not a full data backup). Drive and Apps/Private manager copies replaced; catalog
hashes/sizes/versions updated. Photos/Mail/Accounts binaries unchanged, existing
files/identity and inbox mail/DNS preserved. Restart hidden explicit run, initial
PID58868; read current running.json for fresh identity. Refresh Drive after any
current uploads finish; check destination before retrying uncertain old rows.

Post-update startup ready8/8, fresh runtime executable path verified; Drive,
Photos and Mail read-only health requests200. Accounts health also200 with its
configured Host (raw-IP Host request403); public-accounts gate passed. No uploads made.
