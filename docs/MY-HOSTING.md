# My Hosting integration

My Hosting is an original native Windows application from
[oneworkspacestudio/my-hosting](https://github.com/oneworkspacestudio/my-hosting).
The preview installer registers this optional app using its pinned organization
release and compatibility manifest. Native runtime/security and fresh MyAccounts
installation checks gate publication. See the app acceptance record for scope.

The existing catalog registers app ID `hosting`, route `/hosting`, loopback port
base+17 (9007 by default), health `/healthz`, callback `/hosting/auth/callback` and
private state `${DATA}/my-hosting`. The normal app picker, owner selection,
lifecycle supervisor, backups and update groups apply. OIDC client credentials
stay in the manager's protected configuration/environment; its configured owner
subject is passed explicitly. Hosting never grants owner status to the first
unrecognised internet sign-in.

Workspace starts the normal MyAccounts code/PKCE flow when Hosting needs a
session. An existing provider session skips another password form. Global logout
revokes Hosting's server-side session and removes its cookies, including for log
streams. The only anonymous mounted route is an exact POST to
`/hosting/hooks/github/<project-id>`; Hosting verifies HMAC before processing it.

`MyHosting.exe` embeds its React dashboard and SQLite control plane. Caddy and
project Node runtimes are additional managed processes. Caddy has one owned
instance, disabled plaintext admin and private mutual-TLS loopback control; hosted websites use separate origins from
Workspace. Never publish user site HTML underneath the authenticated Workspace
origin. Node is optional and provisioned only when an approved project needs it.
Official runtime ZIPs are SHA256 pinned; settings also accept matching offline
ZIPs. Private state/builds need persistent Windows ACLs on local NTFS storage.

The initial scope is trusted administrators/projects, not hostile multi-tenant
hosting. Executable projects use Windows low-integrity restricted tokens,
per-project ACLs and Job Objects. Builds fail closed under LocalSystem or when
restrictions cannot be applied. The supported startup path is Workspace running under the signed-in Windows
user. SCM/Session0/unattended-before-login/ServerCore operation has not been
accepted. Do not install it as a SYSTEM service or bypass process restrictions.

New installations include Hosting only after the organization pin points to a
published accepted release. Existing installations add it through the supported
additional-app installation flow; the updater must not silently add an app the
owner previously unchecked. Update manifests own Hosting's executable and all
its notice packages together. Uninstall retains project/private data by default.

The recommended Cloudflare setup creates exact first-level site addresses such
as `my-site-sites.example.com` and `deploy-ID--my-site-sites.example.com`. Workspace
passes its connected zone/tunnel and protected API token only to Hosting's control
process, never to project builds. Hosting preserves unrelated DNS/routes and
refuses conflicts. It checks fresh public DNS and the expected deployment header
over verified HTTPS before marking a managed tunnel deployment published.
No inbound router ports, zone-wide wildcard ingress or paid edge certificates
are required for this flat naming pattern. The owner's zone is m8e.co.uk; other
installations use their own connected zone.

Optional direct HTTPS supports nested `my-site.sites.example.com` addresses and
requires DNS/public TCP80/443 reachability. A tunnel for nested names needs edge
certificates covering that extra level. Hosting shows DNS/certificate
status and recovery guidance rather than reporting readiness from DNS alone.


Native acceptance includes real static publication/previews, GitHub Express,
Next/Nuxt builds, WebSocket frames, Node crash restart, cancellation descendant
shutdown, runtime repair, retention, canonical redirect and encrypted recovery.
A fresh synthetic Workspace installation reuses its MyAccounts sign-in and global
logout revokes the path-scoped Hosting session. Other-device/reboot/SCM and
cross-machine recovery acceptance remain pending. The app is an original
Apache-2.0 repository, not a cosmetic upstream fork. Builds/releases are local.
