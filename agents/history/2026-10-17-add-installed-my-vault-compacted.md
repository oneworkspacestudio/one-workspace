# My Vault added to the existing suite

## Cause and deployment

The owner again reported My Vault missing from the latest installed version.
Read-only inspection confirmed the active suite was preview30 with12 selected
apps and no Vault binary/private root. Public preview31 already includes all13
apps and the pinned org MyVault v0.2.0-rc.2. Ordinary updates preserve the selected
app catalog, so publishing a complete bundle had not added this app locally.

With the user's requested repair, used the supported `install-apps` action of
the verified public31 manager/bundle. Existing process packages matched exactly;
no existing app upgrade or schema migration was performed. Closed the idle
maintenance window with busy/pending guards, stopped the suite gracefully, then
stopped its persistent shared engine through its authenticated owner endpoint.
InstallApps created its normal config/secrets/replaced-program backup. Restarted
the installed manager hidden. Installed version is now preview31 with13 apps,
including catalog ID `vault`; private storage is `my-vault`.

Private maintenance receipts verify all prior configuration/storage/identity
bindings, decrypted logical suite credentials and ordinary file digests retained;
only Vault was added to the app set. Exact Vault binary SHA256:
`bf8ec3ec6fd5a73ab35946d926efd00cd3b1ffe537b38a1b2c7e4f6e787916fe`.
An initial receipt assertion used the wrong app ID `my-vault`; installation had
already succeeded and restart ran in finally. Corrected the external script and
reran only read-only preservation verification against the original receipt;
all checks passed. Do not rerun the installer script as an idempotent operation.

## Runtime and browser evidence

Startup reached ready100; local and public health endpoints return200/ok with
Vault and all process services running. Existing suite tunnel restarted through
the supervisor; unrelated tunnels were untouched. Dedicated MyAccounts Vault
identity settings were prepared by normal startup. Using earlier owner approval
for the genuine local Vault certificate, its `trust-local` command succeeded;
Windows current-user Root contains the exact generated localhost leaf.

Chrome opened the actual secure local Vault login page without a warning, and
the owner's existing signed-in public Workspace page visibly lists My Vault.
Those two Chrome tabs were intentionally left as deliverables. The IAB provider
reported an authority error for the new localhost certificate; no browser warning
was bypassed. Windows trust and Chrome success do not establish IAB trust.
The local portal process had already cached roots before trust and still showed
its guidance page; its local link works in Chrome. Public /vault intentionally
shows local-access/dedicated-origin guidance until a separate public Vault
hostname is configured. Do not claim remotely accessible vault encryption UI.

The first real Vault account is not enrolled: the login page reports signups
disabled/invitation required. No actual master password or account was created
or reset. Owner enrollment and first master-password setup remain required;
do not confuse healthy deployment with end-to-end encrypted-item acceptance.
The old browser18935 synthetic signup fixture is not this installation.

## Boundaries and follow-up

No tests against live apps, source implementation change, GitHub push or release
was performed for this deployment. Public31 already supplied the required app.
No backups/user data removed; existing credentials, domain, folders and autostart
retained. Sanitized task scripts are outside Git under local
tooling/workspace-add-vault; sensitive receipts/logs stay in protected Private
maintenance storage. Preserve older untracked histories.

Future UX work: expose a supported Add available apps action so older selected
catalogs can opt into newly bundled apps without needing an operator CLI.
