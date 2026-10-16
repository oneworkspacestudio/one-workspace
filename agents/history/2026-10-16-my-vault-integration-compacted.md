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
