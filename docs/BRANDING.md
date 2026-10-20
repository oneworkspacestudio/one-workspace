# Suite branding

The shared SVGs in `internal/branding/icons/` are the source for the Workspace
portal, app switcher, README and generated PNG/Windows ICO artwork. The app
directory in `docs/branding/apps.json` includes every catalogue app, the
Workspace hub and the built-in Backups tool.

## Visual style

Blue is the dominant colour across the collection. Use medium blue (`#2563EB`),
light blue (`#60A5FA`) and deep blue (`#0B2A63`) as the main colour family. Mint and
pale neutrals are accents. Preserve recognisable app symbols and transparent
margins; check legibility at 16, 24, 32 and 40 pixels on light and dark backgrounds.
Clips and Backups use the same blue gradient and restrained mint accents as the
rest of the collection. Avoid unrelated black/red, yellow/green or purple themes.

## Keep the collection complete

When adding or renaming an app, update its SVG and `docs/branding/apps.json`.
Apps sharing a repository still need separate entries. The shared directory
validator rejects missing, duplicated and unlisted SVGs.

With Sharp installed in a local tools directory, run from the repository root:

```powershell
$env:SHARP_MODULE = 'path/to/node_modules/sharp'
node scripts/sync-branding.cjs --local-only
node scripts/render-branding-overview.cjs
```

The first command generates every app's transparent PNGs and multi-size ICO.
Existing unqualified PNG paths remain available for older build scripts. The
second command regenerates the complete overview, its content-addressed filename,
the README app table and the light/dark gallery at `docs/branding/index.html`.
The table descriptions and overview order come from the same app directory.

Use `--local-only` for repository artwork work. The legacy positional projects
root intentionally copies assets into supported sibling repositories; use that
mode only as part of an explicitly scoped component update. It does not rebuild
those apps. Do not silently modify sibling checkouts or published release files.

Review the rendered README at desktop and phone widths. Keep its language
friendly to ordinary users; technical details and release history belong in
linked guides and release notes, as recorded in `AGENTS.md`.

Published executables retain their embedded artwork until a later app or suite
build includes these sources. Existing immutable release assets must not be
replaced. GitHub Actions stays disabled; builds and visual checks run locally.
