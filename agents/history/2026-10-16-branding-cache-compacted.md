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
