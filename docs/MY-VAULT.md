# My Vault integration

One Workspace preview 26 includes My Vault 0.2.0 RC2 from the public organization
release, with its own app tile/logo, light interface, optional MyAccounts session
reuse and independent recovery guides. Fresh installations include it by default.

My Vault is a native Windows Vaultwarden companion from
[oneworkspacestudio/my-vault](https://github.com/oneworkspacestudio/my-vault).
It preserves the established vault protocol/client encryption and provides
independent local management, encrypted full backups, portable-export guidance
and restoration without the original Workspace, MyAccounts or server.

The installed MyVault.exe runs one supervised loopback-only engine, with state in
`Private/my-vault`. The fixed gateway port is Workspace base+14; the following
port serves dedicated local HTTPS. Trust the app's unique localhost certificate
explicitly through its local owner tools before browser use. Its catalog has
no required MyAccounts dependency. The updated manager configures an optional,
dedicated upstream OIDC client without sharing the portal cookie or encryption
keys. Independent email/master-password sign-in stays available. The process lock prevents
another standalone/service instance from adopting that same root concurrently.

The `/vault` entry is a launcher. It never proxies cryptographic HTML through the
shared Workspace origin or injects Workspace navigation into that client. A
dedicated configured HTTPS vault hostname is linked directly. Until remote hosting
is configured, the entry explains setup and provides the local computer's direct
loopback link. Configure a trusted HTTPS proxy/tunnel to that fixed listener using
the app's owner tools. Saving a hostname does not create DNS/certificates.

On startup the updated manager registers only its dedicated `workspace-vault`
application, using the vault's exact OIDC callback, standard signed ID tokens and
PKCE. It stores the optional client in protected `my-vault/myvault-sso.json`.
Create the independent vault account first and verify its matching MyAccounts
email. Workspace launches with a sign-in hint; the unmodified client reuses its
vault session before starting OIDC. Locked clients still need their encryption
unlock. Local owner opt-out persists across starts, and identity setup failure
does not stop the independent vault. No public SSO signup or unverified-email
linking is enabled. The browser's independent-login route works with MyAccounts
stopped; recovery archives intentionally exclude the OIDC client file.

If Workspace or MyAccounts is unavailable, stop any remaining supervised vault
and run the installed MyVault.exe directly. Its adjacent workspace-location.json
selects the same private root. Vault email/master password/MFA remain independent.
Do not delete/recreate the account or change encryption to compensate for a
Workspace sign-in failure.

For server loss, use the independent age-encrypted backup or a password-protected
portable compatible-client export on a fresh computer/account. Preserve separate
master/export/backup passwords, recovery identity and MFA codes off the server.
Offline clients must have synced and remain logged in, and unlocking must be
tested before relying on the cache. Mobile acceptance is deferred by owner.

App updates are a separate verified group; launcher/license notices move together.
Workspace's full stopped-state backup includes the private vault. Portable
independent backup/export/recovery remains necessary and is not replaced by an
installation rollback tied to the original Windows computer.

Read the app's shipped acceptance matrix and export-format limits before relying
on attachment, organisation, passkey, TOTP, SSH or other client migration paths.
No unsupported client or unperformed recovery/platform test is claimed.

## Add My Vault to an existing installation

Ordinary updates preserve your selected app set. To explicitly add the new app,
download the verified complete bundle, stop Workspace and use the installed
manager's existing command:

```powershell
Workspace.exe --data YOUR_PRIVATE_FOLDER install-apps --bundle workspace-bundle.zip
```

Existing component versions must match the bundle; update them first if the
command reports a mismatch. The command retains account state, files and existing
app data. Start Workspace again, then open the My Vault entry. Its independent
account and certificate setup are separate owner actions. Do not reinstall over
or delete existing private data to add a new app.
