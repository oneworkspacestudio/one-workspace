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
