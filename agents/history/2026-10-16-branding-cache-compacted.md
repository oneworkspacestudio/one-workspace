# Workspace branding cache correction

User reported old app icons and Workspace's own logo on the live domain after
the earlier icon refresh. The active preview22 serves current supplied SVG bytes,
but reused unversioned filenames and Cloudflare applies a four-hour asset cache.
Browser/favicon caches can therefore retain old artwork at those same URLs.

The shared branding package now maps content hashes to SVG filenames and serves
only known revisions. App tiles/switcher, portal header/favicon, setup and update
windows use the versioned URLs. Legacy URLs remain with no-cache; versioned URLs
have immutable cache headers and strong ETag validation. No artwork or component
pin changes. Cache-policy, content integrity, stale revision and rendered portal
coverage accompany the change. Build locally with Actions disabled.

Current task authorizes correcting the live site; deploy the manager through the
verified backup/update path, retaining accounts/files/settings and existing live
rollback folders. Do not use live MyAccounts/FileBrowser as test fixtures. Retain
the unrelated untracked installed-preview22 handoff. Prepared version is preview24;
verify build, deployment and public site before recording completion.

## Live correction verified

Code commit `d1f86ca23fc7a7eb706131b4936afd23f4390564` passed all Go tests and vet.
The full locally built installer renders the versioned hub logo/favicon and
serves their exact SVG with immutable caching. Bundle verification confirms only
the manager and Workspace provenance/status notices differ from preview23; all
app definitions and component executables remain pinned and unchanged.

Used the released preview23 updater's explicit Workspace-only bundle path to
upgrade the authorized live installation at C:/tmp/workspace from preview22 to preview24. A full
verified protected snapshot was taken first and is retained as
Private/backups/update-20261016T002023-5ee930f34925. The older rollback folder and
all earlier backups remain. The transaction completed only after all 12 apps were
healthy. Post-update fingerprints confirm unchanged identity state, installation
settings, app definitions and every one of 21 ordinary shared Files. No credentials,
raw private configuration or user file names were copied into source or release.

Public Cloudflare requests confirmed current hashed header/favicon references and
exact content/cache policy for all 13 supplied SVGs. Reviewed a fresh browser capture
of the public sign-in page. The prior source-only publication did not deploy23;
the actual live installation is now preview24, including the update/recovery
feature from23. This is a branding/deployment check, not new full app-flow,
clean-host, reboot/power-loss or Windows-service acceptance.

Owned build/browser C: cache was removed after canonical-path, reparse and process
checks. The live installation remains running; no test suite was created or left
running. Build, bundle/installer checks, sanitized deployment receipts and public
branding screenshot are retained outside Git under tooling/workspace-branding24.

## Publication complete

Pushed the implementation and published
[v0.3.0-preview.24](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.24)
from `d1f86ca23fc7a7eb706131b4936afd23f4390564`. All 15 GitHub asset sizes and
SHA256 digests matched the local staging files before publication. Anonymous
requests then verified the public catalog bytes/version and installer availability.
No Actions build was used; all component pins and prior releases are unchanged.
The packaged manager is the same binary verified in the live installation.

No work remains for this fix. The live site runs preview24 and uses versioned
current app/hub artwork. Future SVG edits automatically generate fresh URLs;
retain that shared branding path for tiles, switcher and static page references.
Do not remove the protected live rollback snapshots as build-cache cleanup.
