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
