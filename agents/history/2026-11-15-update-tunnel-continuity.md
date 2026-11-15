# Cloudflare continuity during Workspace updates

During the user's preview.133 install, the updater stopped the running manager
for a consistent backup. The manager owned its Cloudflare child in a Windows
job, so stopping the manager stopped the only connector. Cloudflare reported
1033, then 502 once the connector returned while the local portal was starting.
The PC's startup recovery loop had logged that it was waiting on the updater.
The proxied DNS record for the host belongs to its dedicated Windows Workspace
tunnel. No DNS or Umbrel routing change was needed. Preview.133 completed,
reported all apps running, and returned the expected public sign-in redirect.

Preview.134 introduces installer integration 20 so an older updater hands off
to the new installer before stopping its runtime. The new updater starts an
independent replica from a verified local copy and checks its loopback readiness
endpoint before stopping Workspace. Active updates and automatic rollback share
the replica. The normal origin returns a generic retryable maintenance response
while the update has stopped the apps and while the runtime starts. The new
runtime verifies that its own connector is registered before it declares ready.

Suite, catalog, distribution, maintenance, and publisher tests passed. The
standard Windows build verified 939 packages from the published preview.133
seed. Fresh setup passed on isolated NTFS data, and the real second connector
registered without disturbing the active public hostname.

The backed-up live upgrade from preview.133 to preview.134 completed. Public
polling recorded 68 healthy responses, 231 Workspace maintenance responses,
one 502 during the origin handoff, and no 1033. The retained recovery copy
verified, the temporary connector closed, and the live portal, Froststock,
and its primary Cloudflare connector are healthy. Test fixture roots were
removed. The Windows installer is prepared for publication. The standard
publisher stopped before a draft because this PC has no Android release signer.
