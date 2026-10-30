# My Tickets

My Tickets joins General and is selected by default. It uses MyAccounts and opens
at `/tickets/` on the Workspace origin, including through the remote tunnel.
Existing installations gain it through the normal 106 upgrade while explicit
app exclusions, existing apps, credentials and user files remain preserved.

The native Windows x64 app embeds the Vikunja2.7.1 UI/API and project/issue
extensions in one executable. Its databases, encryption keys and attachment data
stay in private application storage. Managed updates retain attachment blobs and
back up the application metadata needed for rollback. Independent app updates
and recovery are supported; the new installer integration5 requires manager105; use106 or newer with Drive.

## Acceptance

Real fresh106 setup and normal104-to-106 upgrade passed on this PC with synthetic
NTFS data. MyAccounts sign-in opened Tickets without another password. Desktop
and mobile project/issue flow passed. An independent synthetic app update and
normal verified rollback retained issue data, user files and unrelated binaries.
The complete ordinary published-pin catalog contains34 apps and307 owned packages;
publication rejects candidate builds and requires published component hashes.

Native Chrome/API checks cover scoped tokens, private project permissions,
16 concurrent unique issue keys, attachment lifecycle and encrypted backup
restored to a fresh NTFS directory with both databases and exact file bytes.
10,000-row search took159-170ms on this developer host. Backend lint, frontend
checks and five Workspace integration packages pass. The clean-VM requirement
was explicitly waived; no physical-phone, reboot or elevated-service claim.

## Preview limits

This evaluation preview has a protected shared-filesystem Drive adapter; a remote
Drive API backend is not implemented. Webhook senders must retry failed requests;
there is no durable inbound retry queue. Scheduled due-date automation is absent.
Git resync covers the latest50 commits. Live owner Gitea delivery and Windows SCM
service/reboot are unverified. Full feature statuses and sanitized screenshots
are included in the app release's ACCEPTANCE.md and acceptance evidence archives.

[My Tickets source and releases](https://github.com/oneworkspacestudio/my-tickets)

## Drive startup correction and local delivery

Use Workspace106 or newer. Preview105 is withdrawn: Drive could create an empty
attachment tree before the native SDK ownership marker, preventing first startup.
Preview106 prepares ownership first and repairs only that exact empty layout;
files, foreign markers and reparse points remain protected. Updates explicitly
wait for Tickets readiness.

Fresh106 setup and normal104-to-106 upgrade were repeated with Drive and Tickets
co-installed, using synthetic NTFS data. Both passed MyAccounts sign-in, issue
creation and responsive UI; independent app update and verified rollback also
passed. The app EXE is unchanged from the verified preview1 release.

The authorized primary PC now runs Workspace106/Tickets1. The protected update
retained34 app contracts,282 unchanged installed payload hashes, settings and
credentials; all33 supervised services are running. The remote Tickets route
responds with the expected sign-in redirect. No owner identity or data was used
as a fixture. Final owner-profile UI inspection was blocked by the browser
control URL policy; isolated Chrome flow tests passed. Physical phone, elevated
service and reboot checks remain unverified.

Workspace106 is published and promoted with329 verified assets and a source-bound signed Android10 companion; anonymous APK and exact public channel checks passed. GitHub Actions remained disabled.
