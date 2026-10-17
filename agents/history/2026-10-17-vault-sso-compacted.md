# Vault automatic sign-in repair

See component handoff `my-vault/agents/history/2026-10-17-workspace-sso-compacted.md`
in the sibling repository for implementation and privacy-filtered evidence.

User asked to remove the second Vault login. Managed Vault setup now writes
`provisionAccounts: true` into its dedicated optional SSO configuration. This
allows verified MyAccounts members to initialize their own vault through upstream
OIDC/PKCE, then choose their encryption master password once in the native client.
Standalone identity config defaults to no provisioning; explicit enabled:false
opt-out remains preserved. No user passwords, keys or portal cookies are forwarded.
Updated docs/MY-VAULT.md to describe this first-use handoff.

Component fixes the native canonical SSO organization identifier and duplicate
independent email form. Local owner's unverified-email state was repaired using
an exact match to their authenticated Gmail profile; only email_verified changed
through native admin API. Do not generalize that evidence to new installations.

Compiled manager `0.3.0-preview.31+vault-sso.1` and final Vault
`0.2.0-rc.2+workspace-sso.2`; deployed through two supported verified full-backup
transactions. Both complete; settings/storage/suite credentials/13selected apps
retained. Local/public suite health200/ok after restart. Normal Chrome launch now
reaches the native first-master-password page using the existing MyAccounts
browser session, without another identity login/email form. User credential setup
remains a handoff; no real master password or encrypted-item acceptance performed.
No tests added/run; local builds, JS syntax and all59package hashes checked.

No commit/push/release authorized/performed for this fix. Public31 and component
RC2 are unchanged. Local artifacts/receipts/scripts are outside Git under
tooling/my-vault/seamless-sso; private backups and screenshots remain protected.
Publish the corresponding component and suite pin together before replacing this
local repaired component through public update discovery, which can otherwise
offer the old same-source RC2 bytes. Retain prior deployment/untracked histories.
