# Mail send regression and outbound acceptance

User explicitly requested end-to-end checking and one test to their Gmail.
See stalwart-mail's matching history for the parser fault, implementation and
evidence. Native submission+implicit Sent update had been wrongly rejected as
Unexpected Stalwart response although delivery could succeed.

Mail0.1.0-send-preview.9 now accepts correlated implicit responses and separates
queued-with-Sent-update-warning from failed submission.16 regression cases pass,
JS syntax and Windows builds pass. Browser test used actual C:/tmp/workspace,
MyAccounts sign-in, compose/send, Sent and the recipient's Gmail. One synthetic
WM-20261015-A1 message sent once. Native Google MX250 and delivery.completed at
09:18:10UTC2026-10-15; matching test found in Gmail at10:18BST. It landed in Spam.
SPF and DMARC PASS; no DKIM verdict/signature observed. No inbox-placement,
inbound or all-provider guarantee. Recipient's other mail was not altered.
Proof screenshot outside Git: C:/Projects/artifacts/mail-send-20261015/gmail-received.jpg.

Suite/installer0.3.0-mail-send-preview.12 built from existing isolated
C:/Projects/.cache/photos-upload-suite, preserving Drive drive-upload-preview.3,
Photos photos-upload-preview.2, Accounts/Cloudflared and previous hosting fixes.
Installed Mail and manager updated after normal verified stop, with changed EXEs
and workspace.json backed up under Private/backups/mail-send-fix-20261015-101623
(not a data backup). Catalog hashes/sizes/versions updated; no DNS, identity or
mailbox migration. Fresh initial manager PID16164; use running.json for current.
Startup ready8/8, Mail configured with no hosting/identity error.

Ignored local artifacts: dist/WorkspaceSetup-mail-send-preview.exe,
Workspace-mail-send-preview.exe, SHA256SUMS-mail-send-preview.txt. No current
commit/push/release/Actions authorized. Next: evaluate unsigned-DKIM/IP reputation
for delivery quality; do not represent health or Sent as proof of inbox placement.
