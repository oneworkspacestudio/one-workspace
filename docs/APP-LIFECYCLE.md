# Creating and maintaining a One Workspace app

These rules apply to new apps, changes to existing apps, and delivery repairs.
They are release requirements. A standalone executable, local shortcut or one-PC
configuration change does not complete a managed One Workspace delivery.

## Define the integration before implementation

1. Add the app to `internal/catalog/managed-apps.json`. Choose its stable ID,
   name, organization repository, route, update group, authentication model and
   runtime/storage profile. Declare whether it is a default or optional app in
   category/setup policy. Do not confuse an optional app with a missing app.
2. Reserve public and private ports. Define executable, arguments, working
   directory, origin/base-path handling, health/readiness and dependency order.
3. Put mutable data, settings, keys and logs under its dedicated private-data
   directory. Ordinary user files belong under Files. Neither belongs in the
   replaceable program directory. Use the standard backup/restore mechanism.
4. Define first-use activation and migration from earlier/standalone installs.
   Preserve existing accounts, data, files, credentials, paths and recovery copies.
   Security-sensitive features remain inactive until the owner enables them.

## Complete every managed surface

- Register the full release catalog and ordinary organization rebuild, not just
  a local test catalog. Pin published component versions, sizes and hashes.
- Keep every app-owned runtime payload in its package namespace (`<package>-...`).
  `catalog.InstallPackageIDs` carries that namespace through both setup selection
  and new-app updates. Test selecting and omitting the app; an executable-only
  selection is incomplete even when the full release zip contains its sidecars.
  A changed selection/runtime contract must advance the bundle installer-integration
  guard so older controllers request the new installer before changing anything.
- Register rebased ports, runtime supervision, private/temp folders, health,
  storage validation, recovery scope and package ownership for sidecars too.
- Connect MyAccounts with exact OIDC callbacks, owner binding, authorization,
  logout/session expiry and the shared app-opening flow. Open WebSocket/desktop
  streams must be revoked too; deleting an HTTP cookie alone is insufficient. Any independent login
  is an explicit product decision, never a workaround for missing integration.
- Provide correct proxy/base-path/WebSocket behavior. Public paths require a
  narrow allowlist and native authorization. A new route is not blanket public
  access to APIs or private files.
- Mark an existing header Apps control as a shared-picker trigger (the native
  AUI `apps` link variant is recognized). Keep a usable hub link as fallback;
  never add a floating launcher. Foundation cookie cleanup derives from the
  managed registry and must pass the sign-out coverage test.
- Include the app in setup selection, Workspace tiles, shared picker, categories,
  search, branding and docs. A tile must open the installed service from another
  browser/device; a browser-local `localhost` link is not remote integration.
- Register update discovery, managed update metadata, compatibility/minimum
  manager version and every owned payload. Independent app updates must retain
  unrelated applications and the manager.
- Ship runtime dependencies, licence notices and corresponding source where
  required. Verify/lock downloaded and installed payloads. Do not fetch executable
  dependencies dynamically during ordinary startup.

## Verify the supported user paths

Use fresh synthetic accounts and dedicated test roots. Never use the owner's
identity, files, projects or installation as a development fixture.

Required acceptance for a new app:

1. Build the ordinary complete catalog and verify it against the managed registry.
   Test that omitting the new app causes release validation to fail.
2. Install using the actual installer on a fresh supported filesystem. Check
   setup selection, tile/picker discovery, launch, first sign-in, permissions,
   health, static assets, required sidecars and the primary user workflow.
   Exercise opt-in native setup in the supervised worker's restricted environment,
   not only an interactive shell. Declare required Windows command paths and
   executable extensions explicitly. An idle health check does not verify host setup.
3. Upgrade an earlier catalog/manager through the normal backed-up path. Check
   new default-app discovery and optional-app selection. Verify old app payloads,
   settings, identity, exclusions and files are retained as intended.
4. Exercise independent app update and recovery/rollback where supported. Test
   failure cases such as missing runtime files and incompatible managers.
5. Check desktop and mobile presentation. Include direct navigation, expired or
   split sessions and WebSockets when used. State hardware/reboot tests actually
   performed; browser emulation does not establish a physical-phone test.
6. After authorized owner deployment, verify the installed binary and served
   assets against the tested release, plus the normal completion/recovery receipt.

## Publish and promote

Commit/push complete source, create an immutable release with verified assets and
managed update metadata, then promote the public update channel. Verify anonymous
download/channel selection using the same path used by ordinary installers.
Follow `RELEASING.md`, including the signed Android companion for Workspace
releases. Never replace published assets or bypass a guard to claim delivery.

Use the ordinary `scripts/rebuild.ps1` / `scripts/build-org.ps1` path. With
`managedSeed` set, the organization lock pins a verified public baseline and the
managed component archive; the builder preserves its unrelated app contracts and
payloads. Never seed a release from private installed-machine files. Other pin
fields support the legacy full-component builder and do not override a managed
seed build. Update the appropriate public pins deliberately during maintenance.

Keep the managed inventory, component pins, catalog policy, branding and update
registrations synchronized. Run their tests on maintenance releases as well as
initial delivery. Document capability gaps and security restrictions accurately.

## Completion and incident reporting

Report these states separately: source implemented; managed integration verified;
release/channel available; installed and verified on each target machine. Include
versions and any blocked target. Do not collapse them into “done.”

A local workaround is not a distribution repair. If an app is missing, audit its
full path from registry and package to setup, catalog, navigation, runtime and
updates. Reproduce on a fresh install and an upgrade, fix the shared source and
release path, then remove or migrate task-created workarounds safely. Record the
cause, reusable prevention rules, evidence and remaining work in `agents/history`.
