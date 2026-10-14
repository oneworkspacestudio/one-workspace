# Public connection setup

Use two stable HTTPS hostnames. Workspace is the application origin; MyAccounts
must have its own origin because discovery and callback registrations use its
issuer URL. Do not expose the individual Drive or Photos loopback ports.

With base port 8990:

| External hostname | Local target |
| --- | --- |
| workspace.example.com | http://localhost:8990 |
| login.example.com | http://localhost:8991 |

Drive's exact callback is `https://workspace.example.com/drive/api/auth/oidc/callback`.
The portal's is `https://workspace.example.com/auth/callback`. Setup creates both.
Photos uses authenticated Drive APIs under the same origin. The gateway broadens
the Drive session cookie's path to `/`, retaining HttpOnly, SameSite and Secure.
It sends photo API calls directly to Drive; it does not route them through the
desktop Photos proxy, which has local-only cookie handling.

Cloudflare setup uses the documented [remotely managed Tunnel API](https://developers.cloudflare.com/tunnel/get-started/).
Use a token scoped to your account and DNS zone with Cloudflare Tunnel Edit,
DNS Edit and Zone Read. Enter its account/zone IDs in the wizard and authorize
creation. Both hostnames must be unused subdomains in that zone. The installer
refuses existing records; it never takes over a domain already serving an app.
Use a different hostname from an existing MyAccounts installation when testing
a fresh, separate suite.

Failed provisioning attempts roll back only the tunnel/records they created.
If rollback fails, setup identifies the tunnel to remove manually. Resource IDs
are saved in private data after successful provisioning. The manager retains
only the tunnel connector token in DPAPI, not the setup API token. Revoking the
setup token afterwards does not stop an existing connector; rotate/revoke the
connector token through Cloudflare as appropriate.

The connector uses its own config file and does not edit `~/.cloudflared/config.yml`
or another application's tunnel. Existing proxy mode creates no Cloudflare
resources. Public startup waits for the configured Accounts issuer to become
reachable before starting OIDC clients. Existing proxy mode requires those routes
to be configured by the operator beforehand. Changing origins later requires
coordinated Accounts issuer and client callback updates; merely changing DNS or
one JSON field is insufficient.
