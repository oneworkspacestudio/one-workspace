# Automatic readable Mail addresses

User requests admin email on setup and mailboxes for new MyAccounts users.
See stalwart-mail matching handoff for native implementation and security rules.
Owner gets admin@the configured Mail domain; other users get normalized username
addresses with collision suffixes. Stable sub/native mailbox IDs are retained;
existing messages, aliases and IdP recovery/contact email stay intact.

runtime.go generates a fresh supervisor-only mailbox provisioning token per run,
passes it to Mail via environment and starts maintainMailboxes after apps start.
New mailboxes.go polls the existing MyAccounts administrative API every15s, scoped
to organization myaccounts, paginates100 users, skips disabled/deleted users and
reuses its cookie session. It posts batches directly to local Mail with token and
Origin. Native Mail requires matching local Host/peer and rejects forwarded
requests. No MyAccounts EXE change or broad new app permission was required.
Private mailbox-sync.json and logs/mailbox-sync.log report state without tokens
or user lists. Poll retries transient startup/directory errors automatically.

Targeted mock integration passes: owner present, disabled user excluded, adding
a new MyAccounts user reaches provisioning without any Mail login. Mail's actual
native-engine fixture passes account creation, alias routing before sign-in,
collision safety, existing owner account/alias preservation and stable repeat
assignments. No actual IdP test user or external email created in this task.

Installed suite/installer0.3.0-mailbox-preview.13 bundles Mail0.1.0-mailbox-preview.10,
retains Drive drive-upload-preview.3/Photos photos-upload-preview.2 and prior
Accounts/Cloudflared. Normal verified manager stop and stopped Mail snapshot,
changed EXEs/config backups under Private/backups/mailbox-addresses-20261015-103639
(runtime/log excluded; restore untested). Startup ready8/8 and private sync ready
with1 user; actual public MyAccounts SSO/Compose From shows admin@inbox.m8e.co.uk.
Existing Sent test remains. Native Mail has no identity/provisioning error.

Local ignored dist/WorkspaceSetup-mailbox-preview.exe, Workspace-mailbox-preview.exe,
SHA256SUMS-mailbox-preview.txt. Packaging source C:/Projects/.cache/photos-upload-suite;
Mail source C:/Projects/.cache/mail-mailbox-fix. No commit/push/release/Actions.
Preserve all prior work and separate live MyAccounts/FileBrowser installations.
Continue from running.json for current runtime. Remaining previous deliverability
work includes DKIM and Gmail Spam placement; readable addresses do not establish
inbox placement or external inbound reachability.
