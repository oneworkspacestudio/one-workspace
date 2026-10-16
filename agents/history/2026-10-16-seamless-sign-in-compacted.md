# Workspace seamless app sign-in

## Request and repository state
User requested reviewing every app and removing repeated sign-in screens when
already signed into MyAccounts. Work started from Workspace 296a1eb (preview26
publication handoff). Preserve unrelated untracked installed-preview22 history.
No other component source was changed. No commit, push, release or live upgrade
was requested/performed in this task. No GitHub Actions.

## Implementation
`internal/suite/seamless.go` starts each app’s native OIDC login only when its
own authenticated session is absent. Loopback session probes forward only that
app’s cookies, never portal/provider cookies, user identities or bearer tokens.
Drive supplies Docs/Photos/Studio sessions. Mail/Calendar/Contacts use native
api/session; Notes/MyTube api/me; Music auth/user. MyTube anonymous api/me is
HTTP200 with an empty user ID. Gitea API does not accept browser sessions; probe
its authenticated user/settings page and recognize only its login redirect.
Server/setup/permission failures are not treated as expired authentication.

Canonical /app/ precedes probes: native cookies scoped /app/ are absent at /app.
Initial candidate returned to /contacts and looped; corrected and reviewed.
Guard tied to the current portal session permits one automatic attempt. Explicit
app logout pauses automation; a deliberate app tile/open/retry clears the pause.
Native callbacks retain their state/nonce/PKCE/token verification. Requested app
path/query is retained in a bounded, short-lived, validated same-app cookie.
Workspace native browser app entry goes directly to its OIDC flow, not a hub
landing button. Existing portal session bypasses another authorization. Hub
and switcher use /_workspace/open/{id}; Music worker scope remains /music/.
Startup repairs session reuse flags only for installer-owned workspace clients,
retaining secrets, callbacks, roles and consent. Portal account change/logout
clears app cookies and handoff markers. Standalone app authentication unchanged.

Local Vault entry redirects straight to its independent HTTPS origin only after
a normal Windows TLS handshake validates its local certificate. Otherwise the
existing first-run trust instructions remain reachable. This final guard was
compile-checked after the browser review. Remote
local-only Vault hosting guidance remains. No script/cookie/key is injected into
Vault and its encryption unlock/independent login remain intact.

## Review and limitations
Fresh offline preview26 installation: 13 catalog apps, 11 running process apps.
Rebuilt manager browser review passed Workspace session recovery after restart,
Drive/Docs/Studio/Photos, Calendar/Contacts, Notes, MyTube, published Music, Git,
and MyAccounts. Photos-first established Drive auth. Native Notes sign-out,
refresh pause and deliberate retry passed. No repeated credential entry after
initial synthetic account login. Documentation: docs/SIGN-IN.md + README/technical.

Mail browser navigation stopped with ERR_BLOCKED_BY_CLIENT; no warning bypass,
no full Mail browser acceptance or delivery claim. Fresh Vault certificate
reported ERR_CERT_AUTHORITY_INVALID; no trust change/bypass, no fresh Vault browser
acceptance. Prior Vault RC2 checks remain separate evidence. Published Music’s
old generated-name greeting remains a known component limitation; the withheld
Defender-detected build was not used. Build uses Go1.27.2, buildvcs=false for
exFAT ownership, I: tool caches. No automated tests added/run this task.

Candidate outside Git: tooling workspace-auth-review/Workspace-candidate.exe.
Isolated programs/files on I:, private state under one named NTFS cache directory.
Final process cleanup and source status are recorded at task completion. Do not
adopt this disposable fixture as a live installation or publish its private state.

## Next publication
These are uncommitted source changes, not in public preview26. Rebuild a new
versioned full bundle/installer using the verified org pins, preserve original
preview26 artifacts, and perform normal owner-controlled deployment if requested.
Mail browser block and Vault trust remain explicit verification gates. Do not
claim all apps passed browser acceptance or production acceptance.

## Completion state
Task-owned supervisor/children/shared services stopped; no listeners remained
on the fixture’s 16 ports. Verified canonical roots, no reparse points or active
process references, then removed the disposable C: private/profile/temp tree
and duplicate I: Apps/Files payloads. Protected real installations were untouched.
Receipt/candidate remain outside Git under tooling/workspace-auth-review.
Final executable reports 0.3.0-preview.27-sso-candidate. This is a local candidate,
not a reserved release or installer version. No source commit or external write.
The browser could not close restricted error tabs; ordinary review tabs closed,
and remaining temporary tabs are subject to automatic turn-end cleanup.

## Publication authorized
User subsequently requested pushing all current changes and creating releases.
Publish the Workspace manager changes as v0.3.0-preview.27 with the complete
13-app bundle, using the unchanged org lock/component versions. Individual app
source/binaries did not change, so do not manufacture new component releases.
README, normal build defaults and release/sign-in docs now identify preview27.
No Actions, live installation update or real-account testing. Preserve the
unrelated untracked installed-preview22 handoff. Build, source/tag/checksum and
public-release receipts will be appended after publication is confirmed.
