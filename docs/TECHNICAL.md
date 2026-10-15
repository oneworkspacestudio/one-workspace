# One Workspace technical guide

[Back to One Workspace](../README.md)

Configuration, administration, hosting, backups, extensions and local builds.
For the included apps and installation steps, start with the main README.

## Find what you need

| Topic | Guide |
| --- | --- |
| Accounts, passwords and app sessions | [Accounts and administration](#accounts-and-administration) |
| Day-to-day commands, backups and removal | [Operate and back up](#operate-and-back-up) |
| Public access and Cloudflare | [Connection setup](CONNECTIONS.md) |
| One-folder setup and automatic hosting preview | [Local hosting preview](HOSTING.md) |
| Mail integration and delivery requirements | [Mail guide](MAIL.md) |
| Add an app or build an integration | [App manifests and examples](ADDING-APPS.md) |
| Package contents | [Existing package coverage](GIT-MUSIC.md) |
| Tests performed and outstanding limits | [Verification](VERIFICATION.md) |
| Build from source | [Development and release](#development-and-release) |
| Component sources and licensing | [Component provenance](COMPONENTS.md), [LICENSE](../LICENSE), [NOTICE](../NOTICE) |

## App documentation

| App | Source or guide |
| --- | --- |
| Drive | [Custom FileBrowser](https://github.com/fahimc/filebrowser) |
| Office | [Office integration and build guide](https://github.com/fahimc/filebrowser/blob/master/OFFICE.md) |
| Photos | [Photos guide](https://github.com/fahimc/filebrowser/blob/master/PHOTOS.md) |
| MyAccounts | [Drive sign-in configuration](https://github.com/fahimc/filebrowser/blob/master/MYACCOUNTS.md) |
| Git | [My Git](https://github.com/fahimc/gitea) |
| Music | [My Music](https://github.com/fahimc/swingmusic) |
| Photo Editor | [Studio integration](https://github.com/fahimc/miniPaint/blob/master/STUDIO.md) |
| Mail | [Mail integration](MAIL.md) |

## Requirements and storage

- Windows x64. Published packages are unsigned evaluation previews.
- Allow at least 2 GB for the installer and programs, plus room for files,
  document history and backups.
- Published installers ask for separate program, private-data and file folders.
  The private-data folder requires a local NTFS volume so Windows can protect
  account settings and credentials. Program/file folders can use another drive.
- The local one-folder preview requires an NTFS parent folder and organises
  `Apps`, `Private` and `Files` beneath it. It checks compatibility before setup.
- Keep existing installations and user data separate from fresh installation
  folders. This preview has no automatic in-place suite upgrade or data migration.
- Initial Drive storage is a shared source, not a private home per person.
  Configure FileBrowser scopes and permissions before inviting other users.
- Account creation and administrator password resets work without an email
  service. Email verification and self-service recovery are optional.
- The installer includes component binaries, editor assets and the optional
  Cloudflare connector. End users do not need Go, Node.js, Python, Docker or a
  separate Office server to install the published package.

## Published packages and local previews

The [releases page](https://github.com/fahimc/one-workspace/releases) is the source
of downloadable installers and their included-app lists. Component versions,
hashes, notices and verification records accompany each package. Source and
screenshots can be newer than the binaries embedded in a published installer.
Photo Editor is a separate integration; the Mail-enabled installer is a local
preview rather than an update to the published EXEs.

The newer one-folder, domain-first Cloudflare and startup-checklist screens are
local development previews. Publishing these guides does not publish those
pending implementation changes or a new installer. See [HOSTING.md](HOSTING.md)
for the preview workflow, screenshots and exact boundaries. The published setup
wizard retains its earlier folder and connection choices.

## Accounts and administration
The wizard creates a Workspace owner in the ordinary `myaccounts` organization.
Sign in to the portal, then use **Manage accounts** to add/edit/remove users and
reset passwords. Email delivery is disabled initially; administrators can reset
passwords without email. Configure an email service separately in MyAccounts if
you want verification emails or self-service recovery.

Apps share an identity. An existing MyAccounts session is reused when opening
Drive's **Sign in with MyAccounts** link; Docs and Photos reuse Drive's session.
Each app retains its own session. Portal sign-out clears the portal session;
sign out of Drive and Accounts as well on shared computers. Immediate global
session revocation across all apps is not implemented in this preview.

## Local and public addresses
Default local addresses are `http://localhost:8990` for Workspace and
`http://localhost:8991` for Accounts. Drive is under `/drive/`, Docs opens the
Documents folder, and standalone Photos is at `/photos-app/` on Workspace.
Web apps bind to loopback. The local automatic Mail hosting preview additionally enables a public SMTP listener on TCP 25; Windows/network access still needs to allow it.

For public setup, enter two stable HTTPS addresses, such as
`https://workspace.example.com` and `https://login.example.com`. Existing proxy
mode must route them to the Workspace base port and the next port respectively,
preserving the external Host and HTTPS forwarding headers. Configure those
routes before launching the suite.

The local development preview offers domain-first Cloudflare setup: enter a
domain, authorize Cloudflare once, then install. It discovers the account and
zone, finds unused first-level web addresses, provisions HTTPS tunnelling and
registers MyAccounts clients. Existing records are preserved. Technical fields
are available in the advanced connection modes.

The official connector login retains its standard Cloudflare credential. The
suite's tunnel token and Mail DNS credential are machine-DPAPI protected; Mail
needs the DNS credential for signing-key publication and certificate renewal.
Created tunnel/record IDs are recorded in `cloudflare-resources.json`.
Cloudflare HTTPS tunnelling does not carry Internet SMTP. See
[automatic hosting and its network requirements](HOSTING.md).
See [public connection setup](CONNECTIONS.md).

## Add another app
Future apps have a manifest describing their package hash, launch arguments,
port, route, dependencies and optional OIDC callback. Stop Workspace, add a trusted
single-EXE app, then start it again:

```powershell
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' stop
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' add-app `
  --manifest 'I:\Downloads\my-app.json' --package 'I:\Downloads\MyApp.exe'
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' run
```

The launcher validates the package and adds a portal tile, supervised process
and protected proxy route. OIDC-aware apps receive a dedicated client registration
and credentials through environment variables. App packages run as your Windows
user: only add apps you trust. [Manifest guide and working example](ADDING-APPS.md).

## Operate and back up
Use the small `Workspace.exe` in the **private data folder** with `--data PATH`:

| Command | Result |
| --- | --- |
| `run` | Start the supervised suite |
| `open` | Open the portal in your browser |
| `status` | Show app process status |
| `apps` | List installed apps and routes |
| `stop` | Stop this suite and its child processes |
| `add-app --manifest FILE --package EXE` | Install an extension while stopped |
| `recovery` | Display the local server/Drive recovery credentials in your console |

Keep recovery credentials private. The owner's password is not stored by the
Workspace manager. Secrets use machine-bound DPAPI plus restricted Windows ACLs;
copying encrypted data to another machine is not a migration method. Back up
Accounts through its supported backup procedure, Drive's settings/history and
your file folder. Close document tabs and stop the suite before taking a full
filesystem backup. Logs rotate at 5 MiB per file under private data.

Windows sign-in startup is optional. It starts for the installing user after
they sign in; this preview is not a Windows service. To uninstall, stop the suite,
remove the `OneWorkspace` value under your HKCU `Software\Microsoft\Windows\CurrentVersion\Run`
key if enabled, and remove the Start Menu shortcut and program folder. Retain
your private-data and file folders until backed up. Remove dedicated Cloudflare
resources separately using their recorded IDs. No automatic in-place upgrade or
destructive uninstaller is supplied.

## Development and release
```powershell
go test ./...
go vet ./...
./scripts/build.ps1 -Accounts 'I:\Components\MyAccounts.exe' `
  -Drive 'I:\Components\FileBrowser.exe' -Photos 'I:\Components\Photos.exe' `
  -Cloudflared 'I:\Components\cloudflared.exe' -ComponentNotices 'I:\Components\notices'
```

The result is `dist/WorkspaceSetup-windows-amd64.exe`; SHA256 sums and the
component catalog are beside it. Component binaries stay out of Git. To rebuild
from the pinned component seed in this private repository, run
`./scripts/rebuild.ps1 -Go 'I:\Tools\Go\bin\go.exe'`. Builds and tests run locally;
GitHub Actions is disabled for this repository.

This is an **evaluation release**, not production acceptance. Published payloads
are identified by their release notes and component catalogs; local previews
can use different builds. [Component provenance](COMPONENTS.md) and
[verification/limitations](VERIFICATION.md) distinguish tests performed here
from outstanding clean-host, service/reboot and live Cloudflare provisioning
checks. Component licenses and source distributions accompany release assets;
the bundled editor retains its upstream AGPL terms and notices.

## Brand assets and product images

The source includes coordinated SVG logos, browser icons and Windows executable icons. The showcase above is generated artwork; the setup capture and linked app guides show actual UI. Published installers retain the assets embedded when they were built.

[Brand assets and rebuild steps](BRANDING.md) · [Showcase generation notes](branding/showcase-generation.md)
