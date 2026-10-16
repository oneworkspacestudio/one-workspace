# Sign in across Workspace

Automatic app handoff is included from **One Workspace 0.3.0 preview 27**.
Preview 26 and older releases do not include this manager change.

Sign in to MyAccounts once. Opening Workspace, an app tile, the Workspace app
switcher or a bookmarked app page then reuses that browser session. Each app
still creates and verifies its own session using its registered OIDC client.
No Workspace password, portal session or provider cookie is sent to an app.

| App | Account session |
| --- | --- |
| Workspace | Starts MyAccounts directly; reuses its own active session |
| Drive | Native Drive OIDC session |
| Docs, Photos, Photo Editor | Drive session, including when Photos is opened first |
| Mail | Native Mail OIDC session |
| Calendar and Contacts | Separate native sessions using the shared registered client |
| Notes | Native Notes OIDC session |
| My Tube | Native My Tube OIDC session |
| Music | Native Music OIDC session |
| Git | Native Gitea OAuth session; API token authentication remains separate |
| MyAccounts | Its existing provider browser session |
| My Vault | Independent HTTPS origin; existing vault session or optional native SSO, followed by encryption unlock |

Workspace remembers the requested app path and query during its handoff.
Existing installations also have their own registered clients’ session reuse
flags repaired at startup. Client secrets, callbacks, roles and consent
requirements are preserved. A missing or invalid client produces a startup
configuration error rather than a substitute identity.

## When sign-in is still needed

A new browser, expired MyAccounts session, account switch or provider-required
verification can require account sign-in. An app’s permissions are still enforced.
Vault enrollment and master-password unlock are separate from account sign-in.
For local Vault access, its unique HTTPS certificate must be trusted through the
owner tools. The local launcher skips its instruction page only after Windows
trusts that certificate; Workspace cannot suppress browser certificate warnings.

## If an app does not open

Automatic sign-in makes one attempt. If it cannot finish, or you explicitly
signed out of that app, Workspace offers **Try MyAccounts again** and a return
to the hub. Reopening its tile also retries. It does not repeatedly redirect.
Check that MyAccounts is available and the app’s client/callback configuration
matches the installed Workspace address. Server/setup/permission errors are not
silently treated as an expired session.

If Chrome shows a dangerous-site warning or blocks a request, do not disable
browser protection. The site owner needs to resolve its classification through
the browser vendor. The prior Mail/Accounts browser classification remains a
separate verification limitation until the browser permits the normal flow.

## Sign out

An app’s own sign-out ends that app’s session. Automatic sign-in stays paused
until you deliberately reopen it or retry. **Sign out of Workspace** closes
Workspace access, clears known app cookies and requests provider sign-out.
It does not revoke independent native-client tokens or public sharing links.

## Scope

The handoff belongs to the Workspace manager. Standalone executables retain
their native sign-in behavior. This is still a preview, not a production,
clean-computer, service or reboot acceptance claim.

## Recorded local review (16 October 2026)

A fresh 13-app installation used disposable accounts and private NTFS state.
The rebuilt manager was reviewed against the exact preview 26 component set.

| Review | Result |
| --- | --- |
| Workspace restart with MyAccounts already signed in | Hub opens without credentials or another sign-in button |
| Drive, Docs, Photo Editor, Photos | Authenticated content opens; Photos can establish the Drive session first |
| Calendar, Contacts, Notes, My Tube, Music, Git | Authenticated content opens without another sign-in click |
| MyAccounts | Existing account dashboard opens directly |
| Notes sign-out, refresh, retry | Stays signed out on refresh; deliberate retry restores access without credentials |
| Mail | Handoff implemented; browser stopped navigation with ERR_BLOCKED_BY_CLIENT, so full browser acceptance remains open |
| My Vault | Fresh fixture certificate was untrusted; final code retains trust instructions until Windows validates it, so browser acceptance remains open |

No live installation or real account was used. This review concerns browser
sign-in; it does not establish mail delivery, editing, playback, service/reboot
or every account/MFA policy. No automated test suite was added or run in this
change. The candidate manager compiles successfully. The final Vault certificate-trust gate was compile-checked after the browser review. Preview services are stopped
after review; publication and deployment are separate steps.
