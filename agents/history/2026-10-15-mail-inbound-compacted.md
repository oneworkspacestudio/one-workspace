# First-run incoming Mail networking

User requests missing Gmail-to-Workspace delivery fixed for first-time installers.
Read stalwart-mail agents/history/2026-10-15-mail-inbound-compacted.md for confirmed
router root cause, security limits and message-delivery state. No new email sent: Gmail retried the original message and native ingest plus actual
Mail UI confirmed receipt at11:17:29BST. A fresh test is unnecessary. No publication
authorized.

internal/suite/mail_progress.go consumes Mail inboundState/inboundMessage. Startup
keeps checking while incoming route/probe runs; done requires reachable TCP25.
attention/unverified/legacy fields produce warning with Mail settings guidance;
other apps remain usable. TCP connection does not establish delivery/spam placement.
mail_progress_test.go uses local mocks for ready/blocked/unverified/old Mail cases;
passes with existing mailbox-directory provisioning test. Mail tests cover router
conflict preservation/renewal and bounded public probe interpretation, plus opt-in
fresh native provisioning fixture. No full clean-host installer/reboot test claimed.

Built suite/setup0.3.0-inbound-preview.14 with Mail0.1.0-inbound-preview.11, unchanged
Drive drive-upload-preview.3/Photos photos-upload-preview.2/Accounts/Cloudflared.
Local ignored dist/WorkspaceSetup-inbound-preview.exe, Workspace-inbound-preview.exe,
SHA256SUMS-inbound-preview.txt. Packaging C:/Projects/.cache/photos-upload-suite.

Actual C:/tmp/workspace updated after normal exact-manager-path verified stop,
changed EXE/config backups Private/backups/mail-inbound-20261015-111320, mailbox/files
preserved. Final initial manager56856; refresh running.json for process identity.
Final deployed build public incoming state reachable, no identity/hosting errors,
startup ready8/8 and mailbox sync ready. Actual Mail UI shows Gmail sender/owner
recipient/receipt timestamp; proof C:/Projects/artifacts/mail-inbound-20261015/
gmail-received.jpg outside Git. No separate live installations tested.
Preserve all prior uncommitted suite/Mail fixes and concurrent Drive styling.
No commit/push/release/Actions. Mail network requirements still apply to ISP/shared
NAT/unsupported-router users; guidance replaces misleading automatic-success claim.
