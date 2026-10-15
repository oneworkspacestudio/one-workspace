# Mail integration

User requested a new Stalwart fork with Gmail-style mail, MyAccounts and One
Workspace integration; then authorized Cloudflare mail.m8e.co.uk and a setup
wizard. Fork/source: https://github.com/fahimc/stalwart/tree/main/workspace-mail,
local I:/Projects/stalwart-mail. See that repo's matching compact handoff for
engine bootstrap/auth details, public DNS and runtime state.

## Changes

- scripts/build.ps1 optional -Mail payload/catalog, dedicated /mail/auth/callback.
- scripts/mail-manifest.ps1 generates hash-bound addon manifest for existing suites.
- config maps Mail to base+6 and reserves child HTTP +7 / local SMTP +8, with
  upper-port bound check. runtime passes base path and stable owner subject.
- Existing generic authenticated addon routing and automatic tile system handle
  Mail; no portal auth bypass was added. Mail verifies its own MyAccounts session.
- docs/MAIL.md, component notice/AGPL license and main README explain installation,
  existing release boundary, domain wizard and SMTP limitations.

Dedicated client uses existing suite registration; server-held OIDC access token
is used by Mail against each user's own Stalwart mailbox. No email admin linking.
Mail.exe embeds official unmodified Stalwart0.16.26 native binary; source is in
the fork, AGPL option and helper Apache licenses retained. No Docker/runtime install.

## Evidence and state

Fresh isolated NTFS suite runtime at C:/Projects/.cache/workspace-mail/suite-private,
programs/files in sibling suite-programs/suite-files. Preview http://localhost:19590,
Mail /mail, app19596/engine19597/localSMTP19598. Fresh client registration, MyAccounts
SSO, real empty mailbox and compose display observed. No tests suites added/run
under current request; no Internet messages, clean-host, reboot/service or full
mail acceptance claim. Existing live MyAccounts/FileBrowser and previous19390
Git/Music preview untouched. Synthetic owner credentials stay in protected cache.

Local Mail.exe and suite setup compiled without Actions. The local ignored
dist/WorkspaceSetup-mail-preview.exe is an evaluation bundle of baseline
Accounts/Drive/Photos plus Mail; Git/Music optional inputs were not included.
Existing GitHub EXE releases were not replaced; current request publishes sources.

Concurrent branding work was preserved. Its publication moved main through
d333daf and b71ccf7 during this work. Stage only Mail changes/history; reconcile
latest Git state before pushing. User preference: never use GitHub Actions.

Public https://mail.m8e.co.uk hosts a separate unconfigured protected Mail wizard
through a new dedicated tunnel. It is not the synthetic local suite and is not
an active Internet email service. User still needs real client/SMTP/domain values.
No existing MX/root records changed. Details and private launch paths in fork history.


## 2026-10-15: domain-first hosting (local, unpublished)

User clarified that hosting configuration should ask for the domain and automate
technical settings. Current task does not authorize another push or release.

- New `internal/suite/hosting.go` discovers a zone/account from an explicitly
  connected official cloudflared login; finds unused first-level HTTPS hosts
  (hash labels for subdomain inputs); reuses existing cert read-only. The setup
  server's Connect endpoint runs verified bundled `cloudflared tunnel login`
  only when its default credential is absent. Credentials stay server-side.
- Setup normal connection mode is domain + Connect Cloudflare; manual/token
  modes remain Advanced. `?screen=hosting` opens that screen for local review;
  install still validates owner/password/folders. Existing owner/storage steps
  are unchanged. No synthetic accounts are created just to show this screen.
- Installer uses existing dedicated tunnel/identity provisioning; config adds
  hostingDomain/dnsZoneId and DPAPI Mail DNS credential. Runtime passes these
  only to Mail and grants it a three-minute startup allowance. Mail performs
  its own asynchronous first-run setup. `-MailVersion` identifies bundled build.
- docs/HOSTING.md, CONNECTIONS, MAIL and README describe normal/advanced paths,
  credential retention for renewal, existing-service boundaries and network gaps.
- Web HTTPS apps remain loopback. Managed Mail intentionally enables SMTP25;
  firewall/NAT/PTR/ISP access and gateway/relay accounts are NOT automated. Do
  not claim domain-only complete Internet delivery or production readiness.

Evidence: Go components/manager/installer and Mail compile; both frontend scripts
pass Node syntax checks. No test suites or mail-delivery tests added/run. Official
Cloudflare login credential was used for read-only zone/DNS capability discovery;
no live DNS records, tunnels, mail services or existing suite state were changed.
Current local baseline+Mail hosting installer is ignored dist/WorkspaceSetup-hosting-preview.exe,
with manager Workspace-hosting-preview.exe. Build version 0.3.0-hosting-preview.2,
Mail 0.1.0-hosting-preview.2. Git/Music optional payloads are not included. No Actions.

Local installer UI preview uses a fresh process and token-only localhost URL;
PID/log metadata are in the C:/Projects/.cache/workspace-mail/hosting-setup* files.
It is a screen preview, not a completed automatic public installation. The old
isolated suite at localhost19590 and public unconfigured mail hostname still
run their prior binaries; they were not upgraded/migrated by this task.

Resume: reconcile these uncommitted source changes and matching Stalwart handoff.
If the user requests validation/release, use fresh NTFS state and an authorized
unused domain, exercise official provider login and public DNS/ACME lifecycle,
then inbound/outbound mail/network gates. Current code cannot configure a router,
ISP policy, provider PTR or a paid mail gateway merely from domain ownership.

Final UI preview captured in docs/screenshots/domain-hosting.png; screen rendered
from the compiled installer (no account form submission or provider authorization
performed). Preview tab remains open; current process metadata in the cache.


## 2026-10-15: setup usability corrections (local, unpublished)

User requested a shorter password rule, foreground folder picker, early storage
requirements/validation, removal of No SMTP needed / Windows x64 evaluation,
and then simplified storage to one parent-folder choice.

- Normal UI sends `folder` only. suite.Request.ResolveStorage derives sibling
  Apps / Private / Files under it. Existing three-path headless requests remain
  supported. A selected parent must be dedicated rather than a drive root.
- Storage step explains local NTFS before selection. Read-only POST/api/storage
  checks the real volume and existing installation collisions on selection,
  typed-path change and Continue. Install handler and suite.Install repeat the
  same preflight before folder creation. Existing private/app state is refused.
- platform.CheckDataVolume finds the closest existing ancestor of a future path,
  checks persistent ACL support and rejects network storage. ProtectDir reuses
  it. Errors name the drive/format and suggest NTFS rather than raw hex flags.
  ACL enforcement is preserved. This does not migrate existing data.
- Native picker now has the foreground setup window as owner, an initialization
  callback that activates/focuses it with temporary input-queue attachment, and
  starts at the current path/nearest existing parent. UI disables duplicate
  opens and server serializes picker requests. No permanent topmost setting.
- UI and request validation minimum are 8, maximum128. Fresh install identity
  setup also updates admin/myaccounts organization passwordOptions to AtLeast8
  before creating the owner; recovery admin still uses a generated long secret.
  The bundled MyAccounts executable and existing live accounts are unchanged.
- Removed both requested labels from setup footer. New `?screen=storage` review
  URL shows the one-folder step without creating accounts or installing apps.

Evidence: Go packages + full manager/installer compile; extracted setup script
passes Node syntax check; diff whitespace check clean. New compiled installer
screen rendered/captured in docs/screenshots/one-folder-setup.png, one path field
and upfront NTFS note visible. Native focus behavior, fresh 8-character account
creation and full installation were NOT acceptance-tested in this turn; no test
suites added/run per current request. No existing live MyAccounts/FileBrowser,
DNS, Windows account policy or user files changed. Read-only MyAccounts source
was consulted for organization policy API shapes; no changes made in that repo.

Local artifacts (same prior links, replaced): dist/WorkspaceSetup-hosting-preview.exe,
dist/Workspace-hosting-preview.exe, version0.3.0-setup-preview.3. Mail payload
remains0.1.0-hosting-preview.2. Checksums in dist/SHA256SUMS-hosting-preview.txt.
Local preview PID/log pointer remains C:/Projects/.cache/workspace-mail/hosting-setup*;
only the idle task-owned old installer was stopped after confirming not busy.
The updated storage preview tab remains open. No GitHub push/release or Actions.


Later on 2026-10-15 the user completed a real install at C:/tmp/workspace and
reported stalled startup. See 2026-10-15-startup-progress-compacted.md for the
recovery, updated progress build and current deployed state; it supersedes the
prior screen-only preview state for that installation.
