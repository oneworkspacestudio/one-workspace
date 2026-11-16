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

The user approved Windows-only publication after the public channel had
already assigned preview.134 to a separate Froststock release and advanced
through preview.136. The tunnel fix now targets preview.137 and uses the
verified public preview.136 bundle, preserving Froststock preview.6.

A preview.135 update on the private preview.134 installation copied an
integration-19 manager into an installation that still required integration
20. The restart log confirmed that the manager rejected that catalog. Recovery
restored the existing preview.134 programs and data. Online and offline update
selection now reject lower-integration replacement managers before taking the
installation down. Valid independent app updates remain available.

[Windows preview.137](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.137)
shipped from `26d7751`. The publisher verified 949 assets and promoted the
public channel in `b5796f6`. A fresh NTFS install passed with Froststock
preview.6. The anonymous public preview.136 updater required the preview.137
installer and matched its published digest. The restored local preview.134
installation remains healthy and its release check offers compatible Workspace
preview.137 and Froststock preview.6. The signed Android companion stays with
preview.136.
