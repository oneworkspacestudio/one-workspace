# My Social

My Social is the optional native social PWA in Workspace preview47 and later.
It uses MyAccounts, the shared Workspace switcher/theme and the existing HTTPS
proxy. Default app route `/social`, listener9013 (allocated per installation),
callback `/social/auth/callback`. Dependencies: Accounts and Drive/My Photos.

Capture and compose photos, publish carousels and Stories, render bounded Reels
with operator-supplied FFmpeg, follow people, comment/save, and use Creator Studio.
Activity includes SSE and optional VAPID Web Push. This is an evaluation preview;
full Instagram parity, Social DMs/calls/Live, commerce, advanced editing, physical
device push and production acceptance remain open. The app's
[acceptance matrix](https://github.com/oneworkspacestudio/my-social/blob/main/docs/ACCEPTANCE.md)
is authoritative for feature status.

All originals reside under the configured Files/My Photos/My Social root, keyed
by immutable issuer/subject. The manager reserves it in native Quantum ACLs and
rejects conflicting explicit allow rules, including ancestor rules. Each media
operation verifies a loopback capability guard; only authentication is briefly
reused, never the authorization decision. Social serves owner originals and
audience-checked renditions. No public/static directory mapping is created.
Existing My Photos files remain private unless explicitly imported and published.
The Social picker displays only that user's protected originals; ordinary Photos
imports use the user's own current Drive session and scope.

Runtime app state is in Private/social. The normal coordinated stopped backup
covers both Private and Files; standalone state-only snapshots do not cover the
external media root. Removing the app preserves originals, identity and the
protective Quantum policy. Do not remove that policy to browse private originals.
Story expiry is server enforced; retention is conservative and owner managed.

Install/update/uninstall uses existing package size/SHA verification and reversible
maintenance. OIDC handoff/logout and shared app navigation include Social. Source:
https://github.com/oneworkspacestudio/my-social. Existing apps' pins remain unchanged.
