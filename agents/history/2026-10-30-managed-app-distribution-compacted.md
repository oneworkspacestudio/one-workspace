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

## Published104 and Windows publisher hardening

Workspace104 source01e96882648e9f2b994b12b438783d5c4a776531 is published via
the mandatory signed-Android publisher:338 asset sizes/digests, approved Android10
certificate/source provenance, anonymous APK and public channel verified. Channel
commitd281625 promoted104 plus Remote5.

Windows publication exposed default CP1252 decoding and Win32 argument limits
with300 payloads. Publication used UTF-8 process configuration and short staging;
main now explicitly decodes UTF-8 and batches below24000 UTF-16 units. Three
regression tests cover600 assets/Unicode, no omissions/duplicates and empty/
overlong input. Large exFAT uploads stalled; official portable gh2.103.3 and
hash-verified NTFS copies recovered only the missing draft payloads. The required
publisher then reverified all338 and performed publication/channel promotion.
No public assets/tags were overwritten. Runtime bytes stay at the tested104 tag.

## Delivery completed on the primary PC

The anonymous103 controller requires104 before interpreting the new runtime
contract; the104 installer anonymously downloaded and verified the public bundle
and found Remote5. Required publication verified338 assets and Android10.

Normal protected103-to104 primary deployment completed with installerIntegration4,
33 installed apps and all32 supervised services healthy. All32 previous app
contracts, installation settings and existing identity credentials were retained;
150 unchanged payload digests and every new Remote payload were verified. The
manager binary matches the tested release. The latest recovery backup remains
verified; ordinary updater cleanup retired one older automatic recovery copy.
An initial low-disk preflight correctly refused before mutation and restarted103.
The successful retry followed space recovery from a disposable Go build archive;
no private data, installed programs or recovery copies were manually removed.

Public HTTPS Remote bootstrap uses MyAccounts; actual served JS/CSS match Remote5.
Anonymous desktop access and invalid origins are denied. Managed startup migrated
the unconfigured preview2 shortcuts to the Workspace route, backed up originals,
retired duplicate startup/stop launchers and stopped old listeners while preserving
legacy programs/data. Evidence: docs/evidence/remote104/primary-*.json and
public-distribution.json. Python's default urllib user-agent received an edge403;
an explicit delivery-verification user-agent and PowerShell both passed. No edge
policy or authentication was changed.

Central lifecycle rules are docs/APP-LIFECYCLE.md and the immutable104 release
asset APP-LIFECYCLE.md. User workspace AGENTS.md points to the same rules. Treat
standalone app delivery, managed integration, publication and per-PC installation
as separate evidence-backed stages. No remaining primary distribution work.
Second PC remains pending by user instruction. Desktop hosting is inactive until
owner setup; no actual host enrollment, physical-phone, reboot, pre-login, secure
desktop, audio, files or individual-monitor acceptance is claimed.

## Host setup repair delivered as Remote7

Reproduced screenshot failure in a synthetic NTFS host-configuration fixture and
isolated registry namespace. apphost.Service omits PATHEXT/PATH; Windows PowerShell
reported CantActivateDocumentInPipeline on icacls.exe before private host setup.
Remote7 supplies standard executable extensions and trusted system command paths
to PowerShell without importing interactive PATH.13 Node regressions pass, including
real restricted Windows PowerShell failure/success, plus Go manager-boundary tests.
The original idle acceptance did not exercise opt-in host configuration. Lifecycle
rules now explicitly require that path in the supervised environment.

Remote6 immutable checkpoint added a runtime file and was correctly refused by
the106 planner before mutation.7 moves the helper into existing lifecycle.cjs;
all125 package IDs/paths match5.6 is withheld;7 release assets/digests and anonymous
checksums/channel were verified; normal organization builder pin is7. No suite
release was needed. Do not replace published6 assets.

Primary was already106 (not earlier104) and stopped at deployment preflight.
Old104 installer refused the newer contract; the installed106 protected Updates
controller successfully applied only Remote7. Transaction and cleanup completed,
latest verified recovery kept,33 unrelated app contracts and locations/identity
retained. All33 services healthy, all125 Remote payload hashes match release, public
HTTPS MyAccounts bootstrap reports7. No PIN/host enrollment was performed; the
owner must retry Enable remote access. Second PC remains pending. Evidence is
docs/evidence/remote-host-setup7*.json. No real desktop/reboot/physical-phone claim.

## Remote8 Windows client policy repair

Remote8 source7edcc79 is published as immutable v0.1.0-preview.8. All12 public
assets match sizes/SHA256; anonymous checksums match. The normal builder pin is
advanced from7 to8 without changing package paths. The public channel includes8,
continues withholding6 and excludes private Shield metadata. This channel refresh
also discovers the already-public MyAccounts and Authenticator component releases;
it does not publish a Workspace suite release or install those components.

Windows Restricted blocks the app's unsigned -File scripts before their existing
host checks run. Configure/verify/migration now use the documented process-only
application-hosted policy option; machine/user settings remain unchanged and Group
Policy retains precedence. Owner-only read-only diagnostics expose six allowlisted
enums; setup errors/logs use fixed reason codes without raw PowerShell output.
17 Node checks pass, including full synthetic NTFS configuration under a Restricted
parent and a unique test registry namespace. Exact package idle/auth/CSRF/PIN/mobile
acceptance passes. Payload namespace is unchanged:124 files plus manifest.

Catalog tests pass. The full distribution run initially hit its existing500ms
large-download timing test; the focused rerun passed. Owner deployment and actual
host activation remain pending. Secondary108 was restarted through its installed
launcher after its public tunnel/local manager became unavailable; public Remote7
and MyAccounts sign-in recovered. Its private drive was observed as NTFS. An old
paging-file error in launcher.log is not evidence of today's failure cause.

Remote8 delivery follow-up: primary normal protected app-only update completed on
Workspace108;125 payload hashes match,33 unrelated contracts and existing
locations/identity/credential contents remain. All apps now report running and the
verified recovery copy is retained. DPAPI ciphertext changes on ordinary Save;
compare decrypted contents only in memory for preservation checks, never store or
print them. No credential change occurred.

Second-PC Updates UI and public bootstrap confirm Remote8 installed on108, with
completed update/recovery retained. Owner-only read-only diagnostics confirm
effective Restricted, machine/user Undefined, FullLanguage, Valid host signature,
NTFS. Actual owner host activation remains pending at the open PIN/setup form;
the owner must submit it. App evidence: docs/host-setup8-delivery.json in my-remote.

The final full catalog/distribution run passed catalog but again failed the
existing timing-sensitive TestLargeDownloadExceedsSingleRequestTimeoutSafely
(expected two chunks, got3). Its earlier focused rerun passed. No distribution
source or timing test changed for this channel/pin-only repair; do not report the
full distribution suite as passing.
