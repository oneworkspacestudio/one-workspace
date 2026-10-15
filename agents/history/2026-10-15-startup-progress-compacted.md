# Startup investigation and progress UI

User reported the installer stuck on Starting for ten minutes and requested
current-install investigation plus actual progress/current/remaining feedback.
No new GitHub publication authorized. Preserve all prior uncommitted hosting
and one-folder work; see 2026-10-14-mail-compacted.md for that context.

## Actual installation and findings

User installed at C:/tmp/workspace (Apps, Private, Files), base8990, Accounts8991,
Drive8992, Photos8993, Mail8996/engine8997/localSMTP8998. Public workspace.m8e.co.uk
and accounts.m8e.co.uk, dedicated tunnel f6a5abcb-d770-4355-b49d-9edb75fdd076.
The first launcher had already exited. launcher.log only said health check timed
out; Accounts log showed it had started, but no downstream app logs existed.
Original logging cannot prove whether the local or public Accounts health gate
failed. A transient first-domain/tunnel readiness delay is an inference, not an
established root cause. The definite UI bug: POST/launch returned after process
spawn and neither monitored its exit nor polled actual startup readiness.

A scoped retry of this install (not the separate live F:/MyAccounts/FileBrowser)
made public/local Accounts healthy and started all apps. The original failure
was not a still-running installation. Mail configured locally but stopped its
Internet hosting plan at preflight because existing m8e.co.uk mail/DNS settings
collide; no takeover was attempted. Preserve this real existing service.

## Implementation

- New suite/progress.go snapshots named steps/status/detail, completed counts,
  timestamps, state, error and PID. Installation checklist is held by setup UI;
  startup is persisted under Private/startup-progress.json. Percent is completed
  checks, not elapsed-time prediction. Current, completed and pending steps show.
- Runtime tracks ports, connector process, each app, public sign-in and portal.
  WaitHTTPDetailed reports observed HTTP result, catches child exit immediately,
  and retains the failed step on timeout. Public Accounts wait5min with feedback.
  Mail native hosting progress/warnings are surfaced; portal can run while Mail
  Internet delivery requires attention. OIDC discovery HTTP client timeout45s.
- Setup monitors cmd.Wait, supports retry, hides Open until ready and corrects
  installed-vs-running wording. CSS [hidden] rule now applies to nav too, removing
  old Back/Install controls after installation. ProcessRunning verifies recorded
  PID/executable to prevent stale ready snapshots from claiming a dead launcher.
- Setup --resume DATA opens owner-only local progress for an installed workspace.
  Manager no-argument browser readiness allowance increased to10min; connector
  start/exit lifecycle logged. No secrets or log contents returned in progress.

## Recovery and evidence

Normal stop, manager-only replacement and restart were applied to this actual
C:/tmp/workspace installation. Prior manager executables and workspace.json were
saved under Private/backups/startup-progress-fix-20261015. Accounts, app binaries,
secrets and user files were preserved. Config/catalog manager hash+size/version
updated to match deployed manager (0.3.0-progress-preview.4). Mail remains
0.1.0-hosting-preview.2. No other installations were stopped/changed.

The updated runtime completed8/8 startup checks. UI displayed all apps healthy,
with Mail's DNS collision as a warning and a working Open Workspace link. Local
and public health endpoints returned200 during recovery. No mail was sent and
no DNS migration was performed. No automated tests added/run. Go builds, Node
syntax and diff whitespace checks passed; current-install recovery was observed.

Local ignored artifacts: dist/WorkspaceSetup-progress-preview.exe and
Workspace-progress-preview.exe; checksums in SHA256SUMS-progress-preview.txt.
Current status monitor metadata/logs: C:/Projects/.cache/workspace-mail/progress-setup*.
Runtime PID is in C:/tmp/workspace/Private/running.json. New status tab retained.
No GitHub Actions, commit, push or release. Follow up separately on mail-domain
migration only with explicit user direction; do not replace existing root MX.

Final status screen was captured from the repaired installation in
docs/screenshots/startup-progress-ready.png. It shows8/8 with the Mail DNS warning;
Back/Install controls are hidden. The live progress tab remains open. Replacing
the status UI process did not stop the recovered runtime; public health remained200.
