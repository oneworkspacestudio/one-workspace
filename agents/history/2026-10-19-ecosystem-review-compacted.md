# Ecosystem review and local fixes

Current user requested a whole-product review and fixes: shared design/navigation,
single sign-in, Drive workflows. This task did not authorize publication. No
commit/push/release, Actions, owner installation change, DNS/tunnel change or live
account fixtures. Starting source main-equivalent15f9c6d on codex/four-apps;
previous preview37 remains the published version. Preserve the two unrelated
untracked October15/16 histories.

## Implemented

- Shared modal app switcher: current app, actual catalog IDs/icons, Search/Drive,
  keyboard/focus/mobile bounds, native app switcher reuse, adaptive appearance.
  Root home follows the browser theme. Scoped Music/Photos/Studio title branding
  retains content titles and legal attribution. MyAccounts and Vault origins
  intentionally receive no injected portal script.
- `/_workspace/connect/{id}` starts each app's native OIDC and returns to an
  installed app entry. Five-minute HttpOnly session-fingerprint-bound return
  cookie, strict same-origin/installed-root validation, logout cleanup. API
  registry advertises appConnections/sharedNavigation so new PDF/Search clients
  fall back on older managers rather than calling missing endpoints.
- PDF connects Drive automatically; first folder and child clicks load listings;
  import closes picker; explicit save-copy preserves originals. Search connects
  sources and returns to query in the same tab (removed external-link behavior).
- Tasks/PDF/Monitor/Search reuse shared switcher; standalone hides empty Apps.
  Notes/Calendar/Contacts share appearance too. Persist only explicit theme
  changes, not incoming storage events: initial browser review caught feedback
  bouncing themes between tabs, now fixed and retested.
- Stable MyAccounts/Backups registry IDs; hide owner-only Monitor launcher for
  members without weakening backend authorization.
- Unconfigured Mail shows guidance before OIDC. Owner-only setup route loads
  DPAPI setup fragment without exposing it to members/registry/logs. Monitor's
  managed probe treats configured:false as incomplete readiness, not healthy HTTP.
- Verified install-apps accepts any requested catalog app + dependencies;
  unrelated exclusions/data/binary checks retained. Actual omitted Photos was
  restored successfully. Legacy nil Request.SelectedApps retains Photos boolean
  contract; full fixture must explicitly select all or pass Photos=true.
- Branded HTML recovery for browser502; APIs retain plain502. Contacts correct
  tab/favicon and Notes capitalization. Added docs/ECOSYSTEM-REVIEW.md + screenshot
  and README link; updated prior four-app PDF test for close-on-import behavior.

## Verification

Fresh NTFS local private fixture; programs/files/builds on I. All19 app entries
plus Backups opened; no1440px page overflow. Music reached player. Mail's pending
domain setup and Vault's certificate handoff are not delivery/unlock acceptance.
Twelve managed browser regressions passed, including real Drive upload/PDF import/
ordinary save-copy byte equality, automatic OIDC return, accessible desktop/mobile
switcher, cross-tab themes, owner Mail setup, registry IDs, readiness, Search return,
dark home. Final scoped titles passed3 native-app browser checks plus preservation
of unrelated content title. Source script scripts/ecosystem-e2e.mjs records the
12-flow suite and accepts only disposable localhost29xxx, synthetic owner and env
package/output/password inputs. It does not save browser credentials.

Targeted Go suite tests passed for connection validation/session binding, app
reinstallation/dependencies, readiness guidance, member visibility, proxy errors,
native sign-in and branding; Monitor readiness unit passed.19 Node contracts/
filters/security checks passed across four apps. AUI generation/TS/Vite and eight
native binaries built (manager + seven app executables). Calendar's tracked built
asset changes are intentional; retain new asset and remove old tracked hash.

External evidence: I:/Projects/tooling/ecosystem-review (build, acceptance,
final-audit, cleanup receipts). The first build raced an in-progress helper edit
and failed Monitor compilation; subsequent final builds passed. A first browser
test started before readiness; final harness polls health. Search reconnect's
new-tab bug and theme feedback were real findings, fixed before final acceptance.

## Cleanup / resume

Stopped fixture manager and its separate shared WorkspaceServices/Stalwart
instance, verified no process references or reparse points, removed exact named
C:/Projects/.cache/workspace-ecosystem-20261019 including synthetic browser state.
779,548,443 bytes removed; C free25,285,738,496 bytes at completion. No task-owned
preview/tunnel remains. Owner C:/tmp/workspace and rollback data preserved.

Changes are uncommitted in one-workspace, my-tasks, my-pdf, my-monitor, my-search,
my-notes (org notes), my-calendar (Calendar+Contacts). No version/catalog/release
pins changed. Publishing later requires current user authorization, local release
builds for changed components, matching suite pins/bundle and usual acceptance.
Do not publish external build receipts or synthetic browser state.

Remaining product gaps are explicit in the review: general Drive pickers for
Notes/Mail/Chat, bidirectional Tasks/Calendar sync, older native theme controls,
Music shared-library vs Drive scope semantics. Vault unlock/HTTPS, private
databases, repository internals and explicit Hosting publishing stay protected.
Do not describe all workflows/themes as fully unified or production certified.

## Authorized release follow-up

User explicitly requested push and release without Actions after the review.
Preparing preview.38 with Tasks/PDF/Monitor/Search preview.2, Notes preview.6,
Calendar/Contacts preview.3. No owner-installation deployment requested. Published
release receipts and final component pins will be recorded after verification.
