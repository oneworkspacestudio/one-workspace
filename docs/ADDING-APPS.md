# Adding applications

An extension is a trusted Windows x64 executable that serves HTTP on a loopback
port. Bundle its web assets into that executable. Multi-file installers, MSI
packages, arbitrary shell commands and automatic migrations are not supported.

Build the included example and generate its real package digest:

```powershell
go build -o dist/Hello.exe ./examples/hello
./scripts/hello-manifest.ps1 -Package dist/Hello.exe -Port 9004
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' stop
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' add-app `
  --manifest dist/hello.json --package dist/Hello.exe
& 'F:\Workspace-Data\Workspace.exe' --data 'F:\Workspace-Data' run
```

Example contract (replace size/hash with actual bytes; the helper does this):

```json
{
  "schema": 1,
  "package": {
    "file": "app.exe",
    "size": 123456,
    "sha256": "REPLACE_WITH_64_LOWERCASE_HEX_CHARACTERS"
  },
  "app": {
    "id": "hello",
    "name": "Hello",
    "description": "An app added after installation",
    "version": "1.0.0",
    "kind": "process",
    "port": 9004,
    "route": "/hello",
    "stripPrefix": true,
    "health": "/healthz",
    "depends": ["accounts"],
    "args": ["--port", "${PORT}"],
    "env": {"EXAMPLE_DATA_DIR": "${DATA}/hello"}
  }
}
```

IDs are unique lowercase letters/digits/hyphens, starting with a letter. Versions
cannot contain path separators. Routes must be unique, nonoverlapping and outside
the suite's reserved paths. Ports are distinct unprivileged loopback ports.
Arguments are passed directly to the process, without a shell. Supported
substitutions are `${DATA}`, `${FILES}`, `${PORT}`, `${ORIGIN}`, `${ISSUER}` and
`${ROUTE}`. Avoid secrets in manifest fields: ordinary app environment entries are
visible in configuration. The launcher verifies package size and SHA256, stores
it in `apps/ID/VERSION/app.exe`, and verifies/locks it again on every start.

Dependencies determine startup order. `/healthz` in the example is checked on the
app's own loopback port before dependent apps start. A running process is restarted
with backoff after an exit. A later unhealthy or hung process is not automatically
replaced merely because its HTTP endpoint stops responding.

Every extension proxy route requires a valid Workspace session. Implement your
own authorization for app actions and data; a manifest does not create app roles.
Use route-aware links or `stripPrefix: true` with relative assets. Applications
that assume they own `/api` or root paths need adapting; those paths belong to Drive.

For app-native OIDC, add `"oidc": true` and an exact callback such as
`"callback": "/hello/auth/callback"`. The local owner command starts Accounts
temporarily, creates one normal-organization OIDC client, then saves credentials
in DPAPI. The process receives:

```
WORKSPACE_OIDC_ISSUER
WORKSPACE_OIDC_CLIENT_ID
WORKSPACE_OIDC_CLIENT_SECRET
WORKSPACE_OIDC_REDIRECT_URI
```

Your app must implement authorization code flow, state, nonce and PKCE, validate
issuer/audience/signature/expiry, and apply its own permissions. OIDC secrets are
intentionally shared with the trusted app process. Browser-only public clients
must not expose these confidential credentials; use a server-side adapter.

Add-app refuses replacing an existing ID. Updates/removal and data migrations
require an explicit maintenance procedure in this preview. Stop the suite and
back up protected configuration before manual maintenance. Never edit only an
EXE's bytes and update its digest without reviewing the corresponding app source.
