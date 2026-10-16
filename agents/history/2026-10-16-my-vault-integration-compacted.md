# My Vault integration — Windows preview

The owner supplied the My Vault brief and requested a new org repository,
implementation, installer integration and README/logo. Mobile work was explicitly
deferred. Do not change the running C:/tmp/workspace suite or use live accounts,
mailboxes, files or identity state as fixtures. No GitHub Actions or subagents.
Preserve the unrelated untracked installed-preview22 handoff.

New component source: I:/Projects/my-vault, public
https://github.com/oneworkspacestudio/my-vault. Published preview1 from1ff33c0,
native MyVault.exe66,552,832bytes; org ZIP SHA
d0912e0f7780f4dbb4d4525480147f7b1079d1f45627d704a3aee7318064da1d.
Pinned unmodified Vaultwarden1.37.4 application source, dependency-only security
lock updates rustls0.23.46; native MSVC static CRT, bundled OpenSSL/SQLite. Web
vault2026.7.0 crypto JS unchanged. GPL/AGPL notices and corresponding sources
provided. The SDK source archive briefly included upstream's separately licensed
directory: release was temporarily drafted, that source asset/checksums corrected
to exclude44 restricted paths, then republished. Binary/Windows ZIP unchanged.
Keep that source-packaging correction transparent; do not suppress audit findings.

My Vault has independent email/master password/MFA, no MyAccounts dependency,
OIDC secret or portal cookie. Its engine/gateway are loopback-only; native local
HTTPS uses a unique localhost-only leaf, CA=false, explicitly trusted by the owner
for their Windows user. Startup never changes trust. Native local HTTPS is gateway
port+1. Owner tools use separate local capability/Host/Origin checks. A pending
question asks whether the owner permits trusting a test certificate; no answer
has arrived and Windows trust must remain unchanged until approved.

Official Bitwarden CLI2026.9.1 passed login/sync/decryption, secure note, folder/
search/generation/TOTP/attachment, encrypted portable export into a different
account, offline cache unlock with server stopped, and age full-backup recovery
into a fresh data root/clean client profile with exact attachment bytes. Repeated
those recovery checks with the final security-patched executable. Same physical
Windows host; no clean-host claim. Cross-user denial, management/admin/enrollment,
TLS/rate limits/forwarded IP, duplicate ownership, wrong secret/corrupt archives
also exercised. Component acceptance.json remains productionEligible=false for
browser trust/registration/login, desktop/extension, clean host/service/reboot,
MFA/passkey/shared collections and independent security review. Mobile deferred.

Workspace changes: catalog app13 `vault`, base+14 gateway/base+15 local TLS reserved
in preflight; app args run --data ${DATA}/my-vault --listen127.0.0.1:${PORT}
--no-browser, env MYVAULT_WORKSPACE_ORIGIN=${ORIGIN}, no dependencies/OIDC. Logo,
README and branding overview updated, recovery/security/client guides and license
files installed, separate verified updater group/owned notices and org lock pin.

The /vault route is an anonymous launcher only: no cryptographic client proxy or
shared navigation injection. Canonical same-origin comparison includes default443.
Localhost HTTPS origins show setup instead of sending remote visitors to their
own localhost. Configured public dedicated HTTPS origins redirect directly.
Existing updater preserves selected apps. Existing installs explicitly add the
new entry with stopped Workspace install-apps --bundle; identical existing
component hashes required. See docs/MY-VAULT.md, no destructive reinstall advice.

Fixed installer Autostart(false): only remove a matching registration for that
same exe/data pair; a new portable/test install must not disable another install.
All Go tests/vet pass, including new origin/session/port integration coverage.

Full org build verified13app definitions. Final manager was refreshed after the
last localhost-launch correction, retaining all component ZIP entries byte for
byte with archive/zip.Writer.Copy; catalog manager digest and embedded setup were
rebuilt from final source. Final native install uses actual packaged installer in
I:/Projects/tooling/my-vault/workspace-suite/Apps and a disposable C:NTFS private
root under the named my-vault-build cache. APPDATA/LOCALAPPDATA/TEMP isolated,
autostart registration fingerprint checked unchanged. Installer completed; startup
health/anonymous vault launcher/installed guides checks are running (exec60797).
Test script always stops owned supervisor/children; verify no leftovers afterward.

Release25 source not yet committed/pushed, final bundle not yet published. Local
dist now holds the final13app installer/bundle/catalog; old SHA256SUMS.txt must be
regenerated before release. Verify every package hash/size and source provenance,
record native receipt and honest limitations, commit only task paths, then publish
preview25 with no Actions and compare every GitHub asset hash before publication.
Owned large build artifacts under tooling/my-vault; native Rust target/vendor C:
cache already removed after verified corresponding-source preservation. Remaining
test fixtures/toolchains/cache may be removed only after exact canonical/reparse/
process checks, preserving live installation and its backups. Do not publish
synthetic fixture credentials/session/recovery secrets or private configuration.

## Native suite acceptance passed

The actual packaged offline installer completed a fresh disposable installation.
All13catalog entries were installed, all11process app health checks passed (Docs
and Photo Editor are builtin entries), the My Vault executable exactly matched its
verified org release, the unauthenticated independent launcher showed the correct
dedicated local HTTPS link, and recovery/security guides were installed. Previous
autostart registration was unchanged. Test supervisor and its app processes were
stopped by the harness. Receipt: tooling/my-vault/workspace-suite/acceptance.json.
This verifies installation/health/integration, not renewed full flows in every app,
clean physical host, real reboot/service, browser vault or desktop acceptance.
No live installation or certificate trust was changed.

## Published and public download verified

Published v0.3.0-preview.25 from42363fb72cec752b327b6e6e5f5ce87757d976e1.
All16GitHub assets matched local sizes/SHA256 before publication; no Actions.
Final manager SHA118928429fada318a322fecbce832423dabdf1b68d342b33ffb031ff45e50013.
The released installer then anonymously downloaded and verified the latest public
preview25 bundle and reported all13app versions including vault0.1.0-preview.1.
Receipt tooling/my-vault/public-install-check.json. No live update was performed.

Removed compiler tools/cache and duplicate synthetic C: installations after
canonical/reparse/process checks. Official rustup's twelve relative file symlinks
were individually validated/removed without following them before recursive cache
cleanup. C: free space afterward8,221,650,944bytes. Retained only one stopped
synthetic native-final-recovered root for the unanswered certificate/browser
check; no Windows trust modification. Shared fixture service and Stalwart were
explicitly stopped; no fixture process remained.

README follow-up clarifies vault sign-in/preview status and cache-busts the updated
logo overview. Remaining production acceptance gates stay open. The computer-use
skill excludes automating password-manager desktop apps; desktop/clean-host checks
need manual acceptance, and browser testing awaits explicit local test-certificate
approval. Mobile was deferred, not silently counted as tested.

Local origin was still the earlier fahimc repository. Corrected it to
https://github.com/oneworkspacestudio/one-workspace.git and confirmed README/history
follow-up2bd05bf on the org main branch. Release25 remains pinned to42363fb;
follow-up edits are documentation only. Keep future pulls/pushes on the org.
The initial follow-up also reached the legacy remote; no history was reset or
release created there. Preserve existing legacy history rather than force-pushing.

## Optional Vault identity and local release candidate

Owner requested branded production completion and existing-session sign-in. Local
changes only; no new commit/push/release, public preview25 and org lock unchanged.
Live installation untouched. My Vault remains independently recoverable, with no
required Accounts catalog dependency and a separate HTTPS browser origin.

prepareVaultIdentity registers only dedicated admin/workspace-vault MyAccounts
client using existing local recovery-admin API, standard JWT/code/PKCE and exact
vault /identity/connect/oidc-signin callback. Saves protected optional config, honors
explicit disabled state, preserves unrelated clients and continues independent
startup on failure. Native fixture registration and unit checks passed. /vault
remains anonymous launcher; workspace=1 triggers standard optional SSO. Existing
identity browser session reused successfully, vault master unlock still separate.
No vault keys/passwords or portal cookies cross origins. Guides/README/build docs
updated. Current source tests/vet passed; targeted tests passed after final URL fix.

scripts/build-candidate.go verifies seed catalog/package size/SHA, copies unchanged
compressed components, replaces only manager/vault/guides and rejects extra entries.
RC labels/new output directory required. Existing dist/workspace-bundle.zip remains
published25. Ignored internal/payload/bundle/payload.zip now holds local RC: beware
before future setup builds. Candidate manager version0.3.0-rc.1 and vault0.2.0-rc.1.
Final artifact tooling/my-vault/candidate-suite-final/WorkspaceSetup-windows-amd64.exe
is675388928bytes; exact vault SHA572da2056458614f131eac9f6145795b47e6905d94dd29fea9827bd4f2738ce6.

Actual final offline installer passed fresh isolated13app installation, all11process
health checks, exact vault binary, optional MyAccounts auto-configuration, anonymous
dedicated HTTPS launcher, installed recovery/sign-in guides and unchanged prior
autostart. Receipt tooling/my-vault/rc-workspace-suite/acceptance.json. Programs,
private test root and files removed after supervisor/shared helpers stopped and
no remaining process references. Local testing uses --offline; online still fetches
the public preview25, not unpublished RC. No clean-host/reboot/all-app-flow claim.

Vault final candidate passed independent encrypted recovery/export/import/offline
CLI and exact attachments with MyAccounts stopped. UI/SSO tested in synthetic browser
fixture only. Manual new-account registration pending human credentials; user saw
no page, IAB visibility enabled and blank password form captured for handoff. Current
preview C:/Projects/.cache/my-vault-build/sso-fixture/my-vault on HTTPS localhost18935,
PID79588 (verify), older same-UI RC dist executable. Synthetic identity19951 stopped.
Previous native-final-recovered root is deleted; approved test cert retained. Never
publish fixture secrets or automate new credential entry. Production eligibility
remains false for manual registration, desktop/extension, clean physical host,
elevated service/reboot, MFA/passkeys/collections and independent security review.
Mobile deferred. Preserve unrelated installed-preview22 handoff.

## RC2 first-run testing and final local installer

Owner requested disposable-password testing instead of manual registration. Vault
native SDK account creation, browser sign-in/create/edit, independent CLI decrypt,
same-root executable replacement and retained original password/item/certificate
passed. Browser Create account final submission remains outside automated UI
acceptance; no owner account/password or live state was used. New screenshot caught
dark dialog/card surfaces with light text: complete pinned upstream light-token
palette now fixes OS-dark first-run screens without client crypto/settings changes.
Public recovery page HTTP example corrected and regression-tested.

Current local manager0.3.0-rc.2 and My Vault0.2.0-rc.2 are bundled in
tooling/my-vault/candidate-suite-rc2/WorkspaceSetup-windows-amd64.exe (675386368bytes,
SHA3ee9447043f721fe45a06a4b50516bda14753b56dac62516bcbc51eea27ff818). Vault SHA
bf8ec3ec6fd5a73ab35946d926efd00cd3b1ffe537b38a1b2c7e4f6e787916fe. Seed packages
verified and unchanged components copied verbatim. Actual offline installer passed
fresh isolated13app install/11process health, exact vault, auto identity settings,
anonymous dedicated HTTPS launcher/guides, autostart unchanged/live install untouched.
Receipt tooling/my-vault/rc2-workspace-suite/acceptance.json and package sidecar
ACCEPTANCE.json. Embedded guides record pre-test acceptance; sidecar has final RC2.
Online launch still fetches public preview25; run this unpublished RC with --offline.

Exact RC2 vault repeated stopped-source independent encrypted restore, attachments,
export/import and offline CLI unlock. Tests/vet pass; no production sign-off claimed.
All owned previews/suite/shared helpers stopped and RC2 temporary private/program/
files/recovered/client-profile roots removed after boundary/reparse/process checks.
Retained one stopped synthetic sso-fixture source for reproducibility; manual handoff
is cancelled. Ignored internal/payload/bundle/payload.zip is now RC2, dist baseline
bundle still public25. No org lock change, commit/push/release or Actions. Use this
state for later publication only after explicit authorization and open production
gates; preserve unrelated installed-preview22 history and live user installation.

## Preview 26 publication authorized

Owner explicitly requested My Vault included in app/README and all current work
pushed/released. Target v0.3.0-preview.26 keeps13apps and pins My Vault v0.2.0-rc.2
from the org, retaining all other pins. Update normal rebuild defaults and shipped
guides/README to match. Build locally with no Actions; publish exact source,
component provenance, checksums, installer/bundle and independent vault contract.
Releases remain previews and must disclose open production gates. Never stage
the unrelated installed-preview22 handoff or change the live installation.
