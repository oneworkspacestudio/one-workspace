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


## Public78 checkpoint
Source34682f3ce99c48e92dcd30311365d81a7ce8d6a4 is tagged/published78.
Committed-source portal/safe-return tests passed; manager and full Windows setup
compiled, all147 payload size/hash contracts and seven icon frames verified.
All26 apps retained. Baseline public77 catalog; unchanged bytes reused from the
verified public76 bundle, with manager/Mail24/notices replaced. Every uploaded
asset matched local size/SHA256. Actions disabled. Named C mail-signin78-build
scratch was removed after exact-path/reparse/active-build checks. Source/assets
stay on I. Raw-channel promotion and normal deployment are next. Another normal
77 operation was observed on the secondary while78 built; never interrupt/replay.
Main currently77 + Mail24; final78 browser verification still pending.


## Deployment running
Raw channel promotionf001b1c is pushed; unauthenticated index lists78 and Mail24.
Both normal updater plans selected78, with Mail24 already current. Main and
secondary accepted workspace-only transactions after prior77 operations completed.
Do not replay apply on connection loss. Main control script reads its protected
loopback window marker; no capability is logged. Secondary owner update page was
replaced by another task's current tab; use the fresh current page, not an old
OAuth callback. Main/secondary concurrent browser sign-in verification is next
only after both normal updates finish. Preserve each newest verified rollback.


## Completed real deployment and Chrome verification
Both PCs completed normal77-to78 manager-only updates, Mail24 preserved, with all
apps healthy. Main completed in9m15s and secondary15m01s. Each retained its newest
verified rollback and removed one older automatic copy (~4413.2MB /3391.1MB).
Main installed manager78 and Mail24 hashes match the verified release artifacts.
MainC free space was17.03GiB; named build scratch and obsolete task-controller logs
are removed. No unknown/recovery/user data was manually deleted.

In Chrome, opened two normal portal authorization entries together on each host.
Main Workspace and Mail both opened signed in as the real owner. On secondary,
Workspace opened as the owner and Mail automatically handed off to the main
Inbox, signed in and displaying retained verification messages. No password
prompt or Login-state error appeared. This uses existing owner/provider sessions;
it is not a new-password, arbitrary-provider or clean-host acceptance claim.
A callback without credentials stayed rejected with the new plain-language
recovery page; Continue sign-in returned to the authenticated secondary Workspace.
All code/PKCE/token validation constraints remain as documented above.

After both78 updates, sent a labelled Gmail message, opened its exact content in
Mail via the shared client entry, replied from Mail and verified receipt in Gmail.
No test message body, OAuth artifact or credential is copied into this history.
Earlier24 Advanced-panel checks survived a live checking-to-responding transition.
Shared main setting persisted through76,77,78. Main remains the only Internet
receiver; no independent identity or mailbox stores were merged. No passwords
changed. No remaining task-owned preview service or C build scratch is required.
Public78/Mail24 plus promoted channel are delivered. Physical reboot, clean-host,
arbitrary-router/provider, mobile and production acceptance remain unperformed.
If a future issue occurs, capture its exact normal URL/error and relevant version;
do not reset passwords or replace another receiver's rule speculatively.
