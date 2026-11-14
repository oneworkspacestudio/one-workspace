# Froststock and My Hosting listener collision

On a Workspace with base port 9022, Froststock's saved port was 9080.
My Hosting's Caddy also listened on IPv4 port 9080, so Froststock could not
bind. The Workspace proxy then reached Caddy, which returned the
unknown-hostname page.

Managed Froststock now binds to IPv6 loopback. Workspace uses that same
address for preflight, startup and retry health checks, the proxy, session
probes, logout and the monitor registry. The saved port, public route,
OIDC callback and My Drive layout are preserved. The preview.131 build
uses the verified preview.130 bundle as its baseline to retain other apps.

The Windows regression test holds the IPv4 port and checks that the
authenticated portal route reaches the IPv6 listener. The native
Froststock executable also passed health on [::1]:9080 while the real
Caddy answered separately on 127.0.0.1:9080. The fixture used fresh NTFS
private and Drive paths and was stopped after the check.

Suite, catalog, distribution, Froststock and apphost tests and vet passed.
Full installer, signed publication, update-channel promotion and live
installation still need verification.
