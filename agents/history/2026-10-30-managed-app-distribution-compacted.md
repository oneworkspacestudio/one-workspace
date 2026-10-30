# Managed-app distribution repair

User challenged the local-only My Remote shortcut workaround and authorized fixing
the shared distribution plus durable app creation/maintenance rules. The app had
been published standalone but omitted from Workspace catalog/build/runtime/SSO/
navigation/update registrations. Standalone installation was not managed delivery.

## Implemented

- My Remote3 public component: source5f98486629de81fcad256827ecfac61bcf6717d3,
  https://github.com/oneworkspacestudio/my-remote/releases/tag/v0.1.0-preview.3.
  All11 asset sizes/digests and anonymous checksum download verified. Requires104.
- Full managed catalog, port base+35, storage/recovery, runtime supervision and
  sidecar hash locks, MyAccounts owner sign-in, narrow support/public gateway,
  shared picker/category, branding and independent update registration.
- App starts idle without host credentials; only explicit owner/PIN setup enables
  desktop access. TightVNC/Node publisher signatures, private broker, lease and
  approved guest rules retained. Unconfigured legacy launchers have conservative
  backed-up migration; enrolled standalone state requires migration review.
- docs/APP-LIFECYCLE.md is mandatory from AGENTS.md. managed-apps.json inventories
  all33 apps. Build/publication reject omissions; registry tests cover runtime,
  SSO and update coverage. scripts/rebuild.ps1 uses the verified published103 seed
  plus pinned Remote3, preserving unrelated public baseline contracts/payloads.
- Actual fresh installer testing caught executable-only package selection: Remote
  sidecars were dropped. Fixed InstallPackageIDs generically for each selected
  app's package namespace; deselection still omits runtime and retains notices.
- Bundle installerIntegration4 forces older controllers to offer the verified104
  installer before mutation. Update plans retain this guard in installed metadata.
  App-update schema2/integration3 is a separate format and remains unchanged.
- Shared picker recognizes the explicit native AUI Apps-link variant; no floating
  launcher is created. Foundation cookie cleanup derives from the registry, with
  coverage for every foundation app. Remote broker/host ports are included in
  fresh-install group collision checks.
- Android10 category fallback and download links include Remote; signed publication
  must still use scripts/publish-release.ps1. Actions remain disabled.

## Verification and release checkpoint

Catalog/distribution/suite tests pass, including missing/tampered runtime rejection,
new-app omission, sidecar selection/deselection, independent Remote update and
injected-health recovery retaining unrelated app data. Six app bridge/security
checks and exact-component idle/auth/CSRF/invalid-PIN/mobile acceptance pass.
Branding35 app/hub/backup entries render at desktop/phone without missing images
or overflow. Tests use synthetic NTFS data only; no personal identity/files used.

The first large setup link exhausted C: temporary space; task GOTMPDIR now uses
the build drive. Large disposable payload cache objects can be removed only from
this task's cache after builds finish. Never clear shared caches or owner data.

Current resume: finish the104 integration4 candidate build, rerun supported fresh
setup and103 upgrade (including old-controller refusal before mutation), verify
MyAccounts OIDC/picker/public boundaries/logout/desktop-phone UI. Then commit/push,
prepare sources/notices/evidence, publish104 through the required signed-Android
publisher, verify anonymous channel/installer, and apply normal protected primary
update with old app contracts/bytes/credentials preserved. No owner deployment
has been performed at this source checkpoint. Second-PC install remains pending
under the user's earlier instruction; browser control was unavailable. Do not
claim actual desktop activation, second-PC, physical phone, reboot, pre-login,
secure desktop, audio, files or individual monitor acceptance.

## Final component pin

Remote4 source0dee4e0a52cbab34627d63f8f0d92b19e10cced7 adds validated
host reactivation before committing enabled state. Nine app tests and exact
component idle acceptance pass. Workspace104 must pin this component, not the
earlier immutable3 checkpoint. The first full supported fresh/103-upgrade run
passed with the104 manager, shared header picker, owner/guest boundaries and
server session revocation; repeat against the final4 component bundle before
publication. Runtime source/provenance and desktop acceptance limits stay explicit.

## Owner stream revocation before publication

HTTP/native session deletion does not revoke already-open desktop WebSockets.
Remote5 sourcef2f7efbe289c547603b04d8af7dbde93bb074986 adds a loopback,
origin/client-secret authenticated revoke-only channel. Workspace owner logout
requires confirmed acknowledgement while Remote runs, then deletes sessions.
Ordinary members cannot trigger owner revocation. Broker clears owner leases and
terminates owner sockets immediately; guest grants retain separate stop/expiry.
Eleven Node tests, Go manager-boundary tests and exact component idle acceptance
pass. Root catalog/distribution/suite tests pass with confirmed/failing revocation
and member isolation. Final coupled104/Remote5 fresh/upgrade acceptance and public
publication/primary install remain the next gates. Earlier full104/Remote4 gate
passed; it is not the final publication component.

## Final release candidate accepted

The exact104/Remote5 installer passed fresh setup inventory/default selection,
fresh installation, old103-controller refusal before mutation, normal protected
103-to104 upgrade with completed backup/transaction, owner/selection retention,
real MyAccounts OIDC, header picker, owner API/CSRF checks, guest shell/socket
boundaries, manager-only endpoint rejection for browser credentials, confirmed
Workspace logout/native token revocation and desktop/phone no errors/overflow.
Evidence is docs/evidence/remote104/. Synthetic services were stopped afterward.
No real host enrollment was performed. Source is ready for required signed-APK
publication and protected primary deployment; those remain pending at this commit.
