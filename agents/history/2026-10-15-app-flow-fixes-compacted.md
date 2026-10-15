# App flow fixes — local installation acceptance

## Scope and state
User requested fixes for the installed-app audit. Changes are local, uncommitted
and unreleased. Preserve earlier AGENTS/cleanup/audit history changes. The active
user-authorized suite was updated in place; ordinary files and private stores
were preserved. One original rollback was retained. Do not touch legacy apps.

## Implementation
- portal.go: safe relative return paths in server-held PKCE state; app gateway
  requires a Workspace session for private browser requests. Native Git/Drive
  Authorization requests still undergo upstream validation. Deliberate public
  Drive shares retain only share-specific cookies; private APIs stay protected.
- navigation.go/navigation/: authenticated catalog, consistent app switcher,
  home/sign-out links injected into top-level proxied HTML with CSP nonce and
  bounded HTML buffering; compressed HTML requested as identity upstream.
- Logout clears app cookie scopes including Drive domain cookies, invalidates
  portal session and POSTs an ID-token-bound logout to the same-origin issuer.
  Exact post-logout redirect is reconciled on the existing portal OIDC client.
  Native access tokens/public shares have separate lifecycles.
- Native upstream redirects preserve mounted app prefixes and external origin;
  external identity redirects are untouched. Navigation refreshes an existing
  exact-scope Music worker on upgrade, without creating new registrations.
- Signed displayName takes precedence over name; immutable subject unchanged.

## Verification
Targeted Go regressions passed. Browser verified Notes return path/query after
Workspace SSO, catalog navigation and friendly portal/Drive identity. Cross-app
fixes and limits are documented in component histories dated today. Mail sign-in
is blocked by Chrome Safe Browsing classification; approved Google review was
submitted successfully, with no delisting decision. Never bypass that warning
or use an alternate session/credential to complete the blocked Mail test.

## Resume
Full local report: tooling/workspace-fixes-20261015/FIXES-AND-VERIFICATION.md
(outside Git, under the projects tooling root). Rebuild/publish only with current
user authorization. Independent users/devices, clean-host/service/reboot,
large interrupted uploads and complete SMTP delivery remain acceptance gates.
Notes adapters, Mail Contacts/Drive pickers and Calendar invitations were feature
plans outside the reproduced defect list; they are not implemented here.

## Final local acceptance
Music final client passed normal SSO, friendly greeting, Recently added fixture
and full suite switcher. Server-authoritative navigation replaced offline HTML
fallback; static PWA assets remain cacheable. Browser logout closed Workspace
app access: authenticated Photos before, Photos/Music gates after; seven anonymous
private API/media routes401. Chrome blocked final Accounts /api/logout with
ERR_BLOCKED_BY_CLIENT; provider browser session termination is unverified.
No warning bypass. UI wording explains this separate step. Google Mail-related
classification review submitted successfully; no delisting decision/live Mail
acceptance. Source histories/output retained, named C build caches removed after
path/reparse/process checks; active installation/rollback retained. Nothing
committed, pushed or released. Detailed local report records limits/open features.
