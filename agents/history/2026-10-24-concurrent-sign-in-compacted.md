# Concurrent portal sign-in recovery — October24

## Scope and finding
User reports accounts-2 sign-in repeatedly fails during shared Mail deployment.
A real callback visibly returned “Login state did not match” after a restart and
several administration/app tabs. Portal used one workspace_flow cookie for all
pending requests: another tab can overwrite it. Fresh single-flow sign-in had
succeeded, so do not reset passwords. No authentication artifact is saved here.
Mail history in 2026-10-23-mail-network-recovery-compacted.md gives network, source,
release, two-PC role and real Gmail delivery context. User authorizes fixes,
local builds, checks, push/release/install; no GitHub Actions.

## Implementation
portal_flow.go hashes each random state into a separate host-only flow cookie
name, restricted to /auth/callback, short-lived, HttpOnly, Secure for HTTPS and
SameSite=Lax. Constant-time matching and atomic single-use map consumption remain.
Correct legacy in-flight requests remain bound; new requests never overwrite
legacy state. Correctly bound restart/expired/replayed callbacks discard the code
and redirect to a fresh normal login, preserving a known safe return route only.
Unbound/forged/missing-cookie callbacks remain400 with a safe retry page, no raw
OAuth data, no external resources, no-store and no-referrer. Existing token
exchange, PKCE verifier, nonce, issuer/audience/signature checks are unchanged.
No credentials are changed or merged. No app or portal authorization is weakened.

## Verification and publication
New pure in-memory tests passed for overlapping tabs, per-flow cookie protections,
PKCE/nonce initiation, invalid/missing/other-tab/oversized callbacks, single use,
expiry/restart/replay recovery and matching legacy binding. No live account or
Drive database was used. Broader portal checks and committed Windows build next.
Choose Workspace78: unrelated77 tag already exists at f7f7456. Base source includes
that safe update-UI reconnect fix; preserve it. Build78 with verified76 full bundle,
replace manager/Mail24/notices and verify every payload/icon; keep all26 apps.
Publish locally, promote raw channel, then normal backed-up update both PCs.
Do not move an existing tag or replace published assets. Physical fresh-host,
service/reboot and arbitrary provider/router acceptance remain unperformed.

## Current real installations
Main C:/tmp/workspace: Workspace76 + Mail24 healthy, post-install Gmail received;
earlier Gmail round trip passed. Secondary I:/Workspace (remote PC) completed76,
saved main shared URL persisted and normal Mail entry automatically opened main
Inbox. Its24-only update reached retained-backup cleanup; observe completion,
never replay. Additional78 manager updates will be needed after publication.
Named C scratch for76 was removed and obsolete temporary75 controller logs cleared.
Use I tooling for large artifacts and only named C NTFS scratch for full linking.
