# Mail (Stalwart)

[One Workspace Mail source](https://github.com/fahimc/stalwart/tree/main/workspace-mail)
adds a Gmail-style native Windows client and a pinned Stalwart 0.16.26 engine.
It can also run standalone. See that project's README and notices for the mail
setup wizard, licensing, portable runtime and current preview limits.

Automatic domain hosting described below is a local development preview. This
documentation update does not publish its manager changes or a new installer.

## Add to an existing installation

Build Mail.exe locally in the fork, then use this repository's current manager:

```powershell
./scripts/mail-manifest.ps1 -Exe C:/Downloads/Mail.exe -Port 8996 -Output C:/Downloads/mail.json
./Workspace.exe --data C:/WorkspacePrivate stop
./Workspace.exe --data C:/WorkspacePrivate add-app --manifest C:/Downloads/mail.json --package C:/Downloads/Mail.exe
./Workspace.exe --data C:/WorkspacePrivate run
```

Choose three unused consecutive ports. The first is the web app, followed by
the private Stalwart HTTP and local SMTP ports. App packages are trusted native
executables running with the installing user's rights. The manifest records
the exact package hash; the manager holds it against modification while running.

The Mail tile is generated from the catalog. `/mail` is routed through the
authenticated Workspace portal; Mail also verifies its own MyAccounts session.
Its dedicated OIDC client has the exact `/mail/auth/callback` callback. Owner
privileges use the Workspace owner's stable subject, never a contact email.
Private state is `${DATA}/mail` with NTFS ACLs and DPAPI credentials.

## Include in a new installer

Pass `-Mail C:/Build/Mail.exe` to `scripts/build.ps1` along with the existing
Accounts, Drive, Photos and Cloudflared arguments. Git/Music may still be included
using their existing optional arguments. Build locally; do not use Actions.

The installer assigns Mail base port +6, engine +7 and SMTP +8 and registers
the client automatically. With **Host on my domain**, it also supplies the
mail domain and protected DNS connection: Mail configures itself on first launch.
The owner sees hosting progress and any blocker, then signs in. Local and advanced
installations retain the manual wizard. The suite manager reserves the internal
ports before startup. Use `-MailVersion` to identify the bundled Mail build.

## Domain and delivery

See [Automatic hosting](HOSTING.md) for the domain-first path and remaining
SMTP network requirements. The form-based instructions below are advanced
options for existing installations; this change does not migrate their domains.


The wizard accepts a webmail hostname, email domain, Cloudflare Zone ID and
request-only DNS token, with a review step before record creation. Existing
records are never replaced. It can configure an authenticated TLS outgoing
relay without sending a message. Normal Internet receiving requires public SMTP,
TLS, DKIM, MX, SPF/DMARC and reverse DNS; HTTP tunnelling does not provide SMTP.
Public-domain settings must be completed by the installation owner.

This is an evaluation integration. Inbox/compose/search/labels/drafts/attachments
are implemented; it does not claim all Gmail behavior. Do not test against the
live MyAccounts or FileBrowser installations. Use fresh private NTFS state.
