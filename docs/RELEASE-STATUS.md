# October 20 release status

## Preview46: mobile Chat and visible Workspace dismissal

Pins Chat9 from the organisation's release, with compact phone navigation,
real conversation times, readable rows/composer targets, and a Workspace picker
whose header/Close stays outside its scrolling app grid. The shared Workspace
picker in other apps now also scrolls only its catalog. Only existing working
Chat destinations are offered; calls/Updates/Communities remain incomplete.

Exact Chat9 packaged mobile/browser/security/offline/retry/QR/navigation/restart
checks passed using synthetic accounts. Official rendered Android/iOS design
examples were reviewed; current reference-account pixel parity and physical-phone
acceptance are not claimed. Full application E2EE and the broader parity brief
remain unfinished. Installer/download/installed checks are recorded separately.
## Preview45: Chat identity, navigation and safe retry

Chat7 repairs the actual Retry menu so it uses the original queued publication
and durable ID. A real server ACK was deliberately dropped after recipient
delivery: Chat6's manual Retry produced two copies. The corrected packaged7
menu produces one copy after reload and retains the original ID. Controller
checks cover retained reply/document headers, ambiguous legacy drafts, failed
acknowledgements and the subscription race with automatic replay.

The first unpublished45 fresh test then exposed a By ID startup race: QR could
encode an unavailable initial ID. The corrected Chat8 waits for identity and
refreshes QR; actual3s delayed-login/QR-pixel tests pass. The UI suite also found
a late-data-end callback crash during navigation. Failed8 candidates were held.
Owned/guarded callback cleanup, deferred-pin fencing, stale subscribe handling
and pending-leave re-entry now pass ordered real-frame network delay,13 controller
units and the complete packaged native/browser/security/UI/restart regressions.
Final45 was built with the released8 bytes. Actual anonymous installer/bundle
downloads passed size/SHA checks; default-online fresh setup, browser/security/UI,
supported rollback and paired42→45 upgrade passed. Existing browser SSO, client
credentials, IDs/history/document bytes and explicit opt-outs were preserved.
No8/45 immutable assets were rewritten. [Exact evidence](CHAT-PUBLIC45-ACCEPTANCE.json).

The operational installation was then updated from42 to45 using those published
bytes. Its20 selected apps comprise18 process apps and2 built-in routes; all18
app processes and Cloudflare are running. Public Workspace, Accounts and dedicated
Chat health endpoints return200. Existing locations/owner/issuer and the latest
verified recovery backup remain. No owner conversation or sign-in flow was tested.

This is still an **incomplete parity candidate**. Application E2EE, signed
device/recovery, complete rich/social/calling and native Android/iOS work remain
unfinished. No production/full-parity, physical-phone or independent-audit claim.

## Preview44: Chat discovery, receipts and identity binding

The candidate combines Chat6's visible/online directory-query refresh and
original group-recipient receipt checks with MyAccounts' exact original callback
binding. Actual synthetic native tests cover account creation while Find is open,
three-recipient acknowledgement transitions, denied forged audiences and retry
after membership changes. Provider migration preserved existing tokens, browser
SSO and Chat identities/history/document bytes; malformed callbacks, PKCE failure,
replay and concurrent redemption were tested against the real provider.

The manager places its owned exact Chat tunnel route before wildcard fallbacks,
repairs only its owned Host override and preserves unrelated rules/options.
Unit/full-suite tests passed. The public44 installer and bundle were downloaded
and SHA-verified. Exact local44 fresh/paired42 upgrade preserved opt-outs, existing
browser SSO, IDs/history and document bytes; old42 refused incompatible6 before
storage, then automatically restored its backup. Actual public Cloudflare routing
remains unverified. The later manual Retry defect is disclosed on both6/44 release
pages; immutable published assets were not rewritten. Use the corrective release.

These are tested increments, not full WhatsApp parity. Application E2EE, complete
device/recovery, rich/social/calling and native-mobile families remain incomplete.
Historical group receipts after subscriber deletion/rejoin still need a durable
ledger; unknown original-recipient state is conservatively Sent. No Actions.

## Preview43 candidate: Chat integration and update compatibility

Adds a dedicated Chat browser origin, preserved MyAccounts client/identity links,
legacy entry redirects, a private installed-app catalog and subject-scoped logout
backchannel. Chat5 adds branded responsive messaging, custom lists, durable drafts,
offline reconnect/idempotency and stricter conversation/media authorization.
The manager advertises its tested origin capability; app manifests can require a
minimum manager, while Chat rejects incompatible managed startup before storage.

Synthetic native browser checks passed SSO, directory/ID lookup, actual direct and
group messages, attachment bytes and denied access, restored outbox replay once,
custom-list synchronization and global logout. This is an incomplete candidate
for the WhatsApp parity brief. Application E2EE, full social/calling/native-client
features and fresh-machine key recovery remain unfinished. The MLS interoperability
proof is separate from the shipped messenger. See [Chat scope](CHAT.md).
Actual public42→43 upgrade preserved selected apps, IDs/history and actual UI
document filename/bytes. Default-online fresh installation passed baseline,
security and UI after directory-ready assertion. Its initial immediate query
missed a new account before synchronization, prompting the Chat6 refresh fix.
The old42 manifest refusal and automatic pre-migration rollback were tested.
Built locally; no GitHub Actions.

## Preview42: My Slicer

Adds My Slicer 0.1.0-preview.1 from the organization release with a matching icon,
shared MyAccounts sign-in/logout, Drive ordinary-file copies, `/slicer/` routing,
private data, backups and default-app installer/updater selection. Native A1/A1 mini
slicing uses a separately downloaded, hash-pinned official engine. See
[My Slicer integration and exact acceptance limits](MY-SLICER.md).
No GitHub Actions are used. Existing preview and component limitations remain.

## Preview 41: automatic adaptive My Tube streaming

Pins My Tube 0.1.0-preview.2. Windows x 64 automatically downloads checksum-pinned
media tools from the original publisher, then queues new and existing completed
originals for bounded 360/480/720/1080p HLS copies. Originals stay unchanged and
playable. Auto starts low and adapts conservatively. Settings shows byte/percentage
progress; failed downloads retry without blocking app startup. Explicit processing
and download opt-outs are preserved. New managed setup enables the feature.

Seven isolated native upgrade/playback checks passed, including real helper
download, measured renditions, unchanged source/ranges, private-media denial and
Chrome decoding on a simulated 1.2Mbps connection. Source 21.38Mbps produced
approximately 0.59/0.95/1.70/3.23Mbps copies. This does not certify every video,
remote network, clean host/reboot/service or an independent security audit.
See [streaming details](MY-TUBE-STREAMING.md). Built locally; no GitHub Actions.

## Preview40: reliable My PDF saves to Drive

Pins My PDF 0.1.0-preview.3 in both the complete installer and update catalog.
Managed PDF uploads, image conversions and browser/native processed PDF copies
save as ordinary Drive files while the app remains open. Per-file save status,
numbered duplicates, durable retries and custom destinations preserve originals.
Existing library documents are not silently exported; use Save to Drive.

Ten isolated managed browser checks, a final native completion smoke, all PDF
Go tests/vet and eleven frontend checks passed. Actual Drive downloads were
compared, including duplicate naming, native rotation, browser copies and failed
or interrupted transfers. See [PDF saving](PDF-DRIVE-SAVE.md) for exact scope.
No live account data was used. Built locally without GitHub Actions. Existing
unsigned preview, clean-host/reboot/service, independent audit and component
limits remain. Schema2 app-data downgrade requires a compatible backup restore.

## Preview37: Tasks, PDF, Monitor and Search

Adds four original Apache-2.0 native Windows applications from pinned organization
releases, plus Notes preview.5's authenticated current-user identity endpoint.
The installer and updater add new default apps, preserve explicit opt-outs and
existing ports, and register app sign-in, private data, backup and health checks.
The launcher, original overview artwork, app marks and README table include all
four applications. No GitHub Actions are used.

Local synthetic NTFS fixtures passed native standalone desktop/mobile flows,
same-window account switching, encrypted new-folder recovery, real MyAccounts
single sign-on, Drive/PDF import and save-copy, Monitor registry privacy, Search
federation/deletion/account isolation, viewer sharing/revocation and suite logout.
Source tests/vet and pinned independent release builds passed. See
[delivery evidence and limits](FOUR-APPS-DELIVERY.md). This is an unsigned preview;
no second-machine, actual reboot, Windows-service or independent security audit
certification is claimed. Existing component limitations below remain.

## Preview36: My Hosting and free site publishing

Adds original Apache-2.0 My Hosting as an optional app from its organization
release, with logo, app launcher, OIDC registration, data/backup/update contract
and license/provenance packages. Existing component pins are retained. Free
Cloudflare site publishing uses first-level flat hostnames, preserves unrelated
DNS/ingress, verifies public HTTPS and avoids router/paid-certificate setup.
Vault/Hosting provider writes share a named mutex. Provider logout uses the
supported OIDC navigation and revokes the browser's path-scoped Hosting session.

Full Hosting native suite passed192.843s including actual Next/Nuxt, static/Node,
WebSocket echo, crash recovery, strict descendant cancellation, runtime repair,
retention, canonical routing and encrypted new-folder recovery. Fresh selected
Accounts/Hosting installer + actual MyAccounts browser SSO/logout and full suite
Go regressions pass. No owner account data was used as a fixture.

Preview for trusted administrators/projects only. SCM/ServerCore/reboot,
unattended pre-login startup, other-machine recovery and hostile multi-tenancy
are not accepted. Supported background path is the interactive Windows user's
startup. See the packaged My Hosting acceptance/operations docs. Existing app
limitations below remain. No GitHub Actions are used.

## Preview35: Chat member search and verified browser flows

Pins Chat0.1.0-preview.4 and synchronizes the MyAccounts member directory every
15 seconds through a private loopback credential. New accounts are discoverable
before first Chat login. Native issuer+subject identity links remain unchanged.
By ID no longer crashes; native IDs, copy/QR links and invalid-input guidance work.
Attachment download paths and PostgreSQL recovery readiness are fixed. Raw engine
authentication/message packets are no longer logged by the Workspace build.

Separate synthetic MyAccounts/Chat Windows fixtures passed OIDC, people search
(name, username, verified email, account ID), native/QR IDs, direct/group messages
both ways, reload history, attachment exact bytes/anonymous denial, and forbidden
account search/session revocation. Focused Go/Node regressions pass. No owner
accounts or messages were used as development fixtures. All other app executable
pins remain unchanged. Existing app and Vault limitations below still apply.

Builds and releases are local; no GitHub Actions. Clean-host, full network/expiry/
global-signout, service/reboot and database-restore acceptance remain pending.


## Preview 34: Chat

The preview33 installer draft was withheld after the first actual installation
exposed Tinode's required empty push-notification configuration. Preview34 pins
Chat preview3 with that generated startup setting fixed for fresh installations.
The existing thirteen apps were restored while the corrected app was prepared.
No account data or user files were removed.

Adds the Tinode-based Chat app to the organization installer and shared launcher.
It uses One Workspace styling and a matching icon, automatically registered
MyAccounts OIDC, independent session-protected transport/attachments and a private
PostgreSQL runtime started by one Windows executable. Browser dependencies are
bundled locally. Chat is optional during advanced installation and has its own
organization release/update contract. Existing app executable pins are retained.

Built locally without GitHub Actions. Chat and Workspace Go builds and the themed
browser bundle compile. No new automated/multi-user runtime, messaging, attachment,
reconnection, service/reboot or database restore acceptance is claimed. Chat is a
preview; its recorded acceptance scope and other apps' existing limits remain.

## Preview 32: automatic Vault hosting and sign-in

Workspace provisions a separate permanent HTTPS hostname for My Vault through
its existing managed Cloudflare connection. It selects an unused name, preserves
existing DNS/tunnel routes and retains the chosen address on updates. Local-only
and manually configured hosting remain available; origin changes retain the
existing encrypted database and keys. Managed upstream SSO permits verified
members to initialize their vault and uses the native canonical SSO identifier.
The extra independent email form is reserved for explicit recovery sign-in.

Builds and packaging are local, without GitHub Actions. Historical acceptance
below remains labelled historical. Current release build/asset verification and
remote deployment observations are recorded separately; no new automated tests,
clean-host, recovery or browser credential-submission acceptance is implied.
Vault remains a prerelease with its outstanding production gates.


## Preview 31: finished updates stop showing progress bars

The updater now enforces hidden elements even when other CSS sets their display
style. Both progress bars appear only during running operations. Completed and
failed operations show a static result, completed/stopped duration, steps and
activity details. Duplicate phase/result text is removed; the visible detail
retains its accessible status announcement. No transaction, backup, download or
startup behavior changes.

The completed screen was visually reviewed in a browser with synthetic data;
running and failed snapshots were inspected as well. Running states retain their
bars and busy controls; terminal snapshots have no visible progress bars. These
are isolated UI previews, not a replay of a real update or rollback. The Go build,
browser JavaScript syntax and release packaging hashes are checked. No new test
suite, live app exercise, clean-host or reboot acceptance is claimed. All13 app
entries and component binaries remain pinned to the same organization releases.
Earlier preview limitations remain. Built locally without GitHub Actions.

## Preview 30: visible update progress

The Windows-local updater shows its current phase, completed and remaining
stages, actual transfer bytes/speed, elapsed time and approximate current-stage
remaining time. Backup copying and verification, package installation and native
app startup report their own progress. An inactive transfer or lost browser
connection has explicit feedback. A transfer reaching100% does not declare the
update finished; completion follows successful transaction and startup checks.
The local capability, origin and operation guards remain in place.

Download labels identify the actual release and asset, rather than an arbitrary
internal package ID. Matching published manager/license/notice/provenance assets
can be downloaded directly when both catalog size and SHA256 agree, with the
verified complete bundle retained as fallback. This avoids fetching an entire
bundle solely for a small changed release notice in supporting updaters.

This version is built locally from committed source; browser JavaScript syntax
and release package hashes are checked. No new tests, full progress UI review,
real update/rollback, clean-host, reboot or service acceptance is claimed for
these progress changes. Earlier Preview29 checks are historical evidence for
that repair, not a full acceptance run of this version. All13 component app
entries and their organization binaries remain unchanged from Preview29.
Prior preview limits, including Vault acceptance and the withheld newer Music
binary, remain. GitHub Actions is disabled and unused.

## Preview 29: app sign-in and Windows update restart repair

The manager now respects Windows host job restrictions when starting maintenance
or the restored runtime, fixing an Access is denied failure from unconditional
job breakaway. Updates check executable launch before stopping apps. A completed
rollback whose startup fails gives restart guidance instead of an unavailable
Recover action. No Windows permissions or antivirus settings are weakened.

Connected app entry reuses existing MyAccounts sessions. Music's scoped worker
integration prevents an older offline shell from replacing login/callback pages;
media caches, other workers and the approved Music executable are retained.
All thirteen component app entries and their verified organization versions are
unchanged from Preview 28.

Full Go tests, vet and the Node worker simulation passed. An isolated Windows
upgrade and real rollback under a restrictive host job preserved synthetic data,
credentials and storage bindings. The local repair used the verified backup/update
pipeline and preserved ordinary files and logical credentials. Chrome then opened
all twelve entries in that existing installation and internal Backups without
entering credentials or clicking sign-in, including Mail and Music. Vault was not
installed there. See SIGN-IN.md for the exact scope.

These are targeted update and existing-session checks. Fresh Mail/Vault login,
email delivery, every editing/media flow, scheduled backup/restore, clean-host,
service/reboot and independent security acceptance are not established. Earlier
component limitations and the withheld newer Music binary remain unchanged.
The Windows installer is built locally without GitHub Actions.

The supplied eleven app SVGs are applied to app headers, favicons, Windows icons
and repository README logos. Component executables were built locally; GitHub
Actions is disabled for all nine organization repositories.

My Tube is now included as a twelfth app entry, pinned to its published
0.1.0-preview.1 release. Its separately supplied logo is used in the README,
launcher and switcher. This suite update is built locally; the previously
verified My Tube component binary is reused without modification. Its optional
FFmpeg/FFprobe helpers remain separate and are not shipped in the suite.

## Music uses the previous published version

Defender detected `Trojan:Win32/FakeFolder!pz` while staging the newly built
Music executable. This is not confirmed to be an incorrect detection. The exact
77,280,038-byte file was submitted to Microsoft Security Intelligence with the
owner's approval. Microsoft accepted the submission on October 15, 2026; its
determination is pending.

SHA256: `5f0d121cd9733775b327ac3b39b46b32cf5d9cbdef1b014880a89d6421239837`

The new Music binary is not released. The complete Workspace bundle uses the
previous published Music version (`workspace-20261014`), with its exact source,
SHA256 and size recorded in package provenance. Its internal icon/profile/cache
updates from October 15 are therefore not included. Workspace's Music launcher
and shared app switcher use the new supplied logo. The organization build script
rejects the withheld binary's exact hash.
Do not disable Defender, add exclusions or rename/repackage the executable to
avoid the detection.

After Microsoft's determination, review its findings, update security
intelligence normally and validate the exact intended release binary before
publishing an appropriate Music update and changing the bundled pin.

## Component previews

The other app previews include matching source archives, license notices,
build provenance and SHA256 sums. Their notes distinguish previous sampled
functional checks from the new branding builds. A requested manual Defender
scan failed to start on the build machine; antivirus clearance is not claimed.
The component versions selected for the next bundle are recorded in the lock
file. Existing installers continue to select the latest published complete
Workspace bundle; they do not assemble unreleased source commits.

Mail's Google Safe Browsing review remains pending, and provider browser logout
was blocked during earlier checks. These changes do not assert a new successful
Mail delivery, clean-host, service/reboot or full identity acceptance run.

## Current limitations

- Mail Contacts search and Drive attachment selection are future integrations.
- Calendar invitation delivery and migration of existing private Mail data into
  shared Calendar/Contacts storage are not implemented.
- Notes' direct Drive picker/Calendar adapter, Joplin synchronization and E2EE
  are not implemented. Its own login can drop query parameters.
- Independent users/devices, simultaneous Office merges, complex Office/PDF,
  large interrupted uploads and clean-host/service/reboot acceptance remain open.
- The Windows executables are unsigned previews. This bundle does not establish
  antivirus clearance; the requested manual scan failed to start on the builder.

Back up private configuration, databases, history and ordinary files before
upgrading. An existing workspace is not silently migrated by the new installer.

## Update and recovery window

Preview23 adds local owner maintenance for existing installations: select compatible
Workspace/app releases, verify downloads, snapshot private/program/file storage,
health-check startup and restore on failure. A durable journal and cached controller
support interrupted recovery and legacy launchers. Newer app releases without a
managed-update contract are shown for review. See UPDATES.md and UPDATE-PUBLISHING.md.
Backups remain until explicitly removed. Existing storage, identity, private Mail
and withheld-Music constraints remain. This does not claim clean-host, real reboot,
power-loss or Windows service acceptance.

## Branding cache refresh

Workspace app tiles, the shared switcher, the hub logo and favicon now use
content-addressed SVG filenames. Setup and update windows use those filenames too.
A new SVG receives a new URL, so browser/CDN entries for previous artwork cannot
hide a branding change. Versioned files have immutable caching and ETags; legacy
filenames remain available with revalidation. App component pins are unchanged.

## My Vault Windows preview

My Vault is the thirteenth app in fresh installations. Preview 26 pins the public
organization release 0.2.0 RC2, with branded light screens and optional MyAccounts
OIDC/session reuse. Its separate HTTPS origin and independent master-password
route do not require MyAccounts. Browser sign-in/save/edit/unlock, official SDK
registration, retained-data executable replacement, encrypted export/import,
offline CLI access and full age-encrypted restore/attachment checks passed in
isolated same-host Windows fixtures. Browser Create account submission, clean-host,
desktop/extension, service/reboot and independent security review remain open.
Read the installed
docs/MyVault/acceptance.json and SECURITY-REVIEW.md. This is an unsigned preview
and is not approved for important production credentials. Mobile is deferred.
Existing installations keep their selected app set during update; no new vault
or account is silently enabled. The installer does not update a live suite as
part of building or publishing this release.

## Preview 27: shared browser sign-in

The Workspace manager opens connected apps through their own OIDC clients when
their native session is absent or expired, reusing the existing MyAccounts
browser session. App tiles, the shared switcher and browser bookmarks avoid
another local sign-in button. Cookie path canonicalization, bounded return-page
handling and one-attempt recovery prevent repeat redirects; explicit app logout
pauses automatic sign-in until a deliberate reopen/retry. Startup repairs shared
sign-in flags on this installationâ€™s registered clients without replacing their
secrets, callback addresses, permissions or consent requirements.

All thirteen apps remain bundled from the same verified organization component
pins as preview26. No component executable or encryption implementation changes.
Local Vault launch skips its instruction page only after ordinary Windows TLS
validation succeeds; independent login and master-password unlock remain.

Isolated browser review reached authenticated Workspace, Drive/Docs/Photos/Studio,
Calendar/Contacts, Notes, MyTube, Music, Git and MyAccounts, and checked Notes
logout/pause/retry. Mail navigation was blocked with ERR_BLOCKED_BY_CLIENT; the
fresh Vault certificate was not trusted. These two full browser acceptance gates
remain open. The final Vault TLS gate was compile-checked after that review.
No new automated tests, real account exercise, live upgrade, clean-computer,
service/reboot, mail-delivery or all-features acceptance is claimed. See SIGN-IN.md.

## Preview 28: optional scheduled backups and app selection

The owner-only Backups app is built into the Workspace launcher and app switcher.
Scheduling starts off. Choose a separate local NTFS folder, preferably on another
physical drive, and an interval of1â€“168hours, or run a manual backup. Full snapshots
include private app data, accounts, programs and ordinary files. The portal and
connector stay running while apps are stopped for copying, hash verification and
restart checks. Status, progress, errors and verified history are visible there.

Advanced installation can omit apps and their executable payloads, with required
MyAccounts and automatic dependency selection. Coupled updates retain omitted
apps. Managed apps use private working/scratch folders on the selected drive;
ordinary shared files keep their normal Files location. Existing data is not moved.
All13 component app entries and their published org binaries remain unchanged.

This release was built locally without Actions. Every bundled package size/SHA256
was checked. No new automated tests, installation run, scheduled-backup/restore,
clean-host, service/reboot or live-upgrade acceptance was performed. Backups are
full copies with application downtime and no automatic retention pruning. Windows
DPAPI still needs the original Windows security context; this is not a bare-disk
or cross-machine recovery wizard. Keep independent Vault encrypted exports.
Read BACKUPS.md before enabling the optional schedule. Prior component limitations,
Mail/Vault browser gates and the withheld Music build remain unchanged.
