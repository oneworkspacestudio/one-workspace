# Mail permission and hosting recovery

Continuation of startup-progress recovery for actual C:/tmp/workspace install.
User reported Mail403 and unclear DNS conflict guidance. No new publication
authorized. Preserve prior uncommitted hosting/storage/progress work.

See stalwart-mail/agents/history/2026-10-15-mail-recovery-compacted.md for causes
and Mail implementation. Native Authentication singleton had no default user
role links; wrapper schema2 migration restores missing ordinary User linkage.
Existing root MX m8e.co.uk -> m8e.co.uk belongs to a pre-existing route; exact
conflict now appears in Mail settings with keep/migrate guidance. No DNS changed.

Suite mail_progress.go now surfaces Mail identityError separately from DNS
hosting errors. This keeps the portal available while an owner can recover Mail.
Current source local build target0.3.0-mail-recovery-preview.6 bundles Mail
0.1.0-recovery-preview.5, baseline Accounts/Drive/Photos/Cloudflared; optional
Git/Music omitted as in the previous Mail preview. Do not claim a new release
of every optional component.

Actual install repair uses normal control stop and exact process-path checks,
then executable replacement and catalog hashes/version updates. Backups under
Private/backups/mail-access-fix-20261015 cover changed EXEs/config only. User
data, identity credentials and existing mail DNS were preserved. Native role
readback confirmed normal User grant and no administrator grant. Recovered
runtime showed all apps running and8/8 startup with Mail DNS warning. Actual
signed-in mailbox access and external send/receive remain unclaimed; user must
sign in again after restart. No test suites added/run, no mail sent.

Local artifacts: dist/WorkspaceSetup-mail-recovery-preview.exe,
Workspace-mail-recovery-preview.exe, SHA256SUMS-mail-recovery-preview.txt.
Read runtime PID from Private/running.json, never reuse a stale PID. Remaining
decision: whether current root mail route is used or user wants a migration;
asynchronous clarification was requested. No commit/push/release/Actions.

## Later user steering: keep root mail, use an available subdomain

User says root mail may belong to another project and wants availability checks
plus a recommended new subdomain. Mail now checks mail/inbox/post/mail2, preserves
occupied records and provides a review/apply flow. Actual mail.m8e.co.uk is an
existing preview CNAME, so inbox.m8e.co.uk was chosen and applied. Native domain
and account IDs retained; Mail DPAPI config is authoritative for its configured
email domain, so the suite's original root-domain environment does not reset it.
Webmail URL remains https://workspace.m8e.co.uk/mail/. Root MX is unchanged.

New inbox-subdomain A/MX/SPF/DMARC have receipts; native DKIM and ACME DNS01
completed. External SMTP delivery is still unverified and must not be inferred
from web tunnel or process health. No mail sent. Protected stopped Mail data
snapshot and pre-change binaries/config under Private/backups/mail-subdomain-20261015.

Final local suite/installer version0.3.0-subdomain-preview.9 bundles Mail
0.1.0-subdomain-preview.8. Artifact filenames remain the recovery-preview names
above. No GitHub publication or Actions. See the Mail history for the protected
domain-change journal and constraints on moving a previously published domain.

Final deployed Mail status has no hosting or identity error and all suite apps
running. Public portal health also observed. New ready screen captured in
docs/screenshots/mail-subdomain-ready.png; earlier mail-recovery-guidance.png
shows the conflict/recovery screen before selection. No signed-in mailbox or
external-delivery acceptance was performed. Updated installer remains local.
