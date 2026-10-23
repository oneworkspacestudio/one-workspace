# Workspace repository

Use repository-history before substantive changes and maintain agents/history.
Never test against the existing live MyAccounts/FileBrowser installations.
Test new suite installations in fresh NTFS private-data folders; program payloads
may use another drive. Preserve component license/provenance and user files.
Do not claim a production or clean-host/service/reboot test not performed.
The bundle is an evaluation build while upstream acceptance gates remain open.

## README and app branding

Keep the main README friendly to people without technical knowledge. Explain
what people can do, how to get started and where to get help in plain language.
Use short, accurate app descriptions. Put implementation details, commands,
compatibility requirements and testing notes in linked guides. Keep numbered
preview announcements, version histories and release-status commentary in release
notes or docs/RELEASE-STATUS.md, not below the app table or elsewhere in README.md.
Do not imply that an unfinished feature is available or that untested use is proven.

List every app, including apps sharing a repository, plus the Workspace hub and
Backups. Keep docs/branding/apps.json, the README app table and the logo overview
in sync when apps are added, removed or renamed. Run
scripts/render-branding-overview.cjs to regenerate the image and README sections;
verify that all images and links work and review the result at desktop and phone
widths before publishing. Use a new content-addressed overview image when artwork
changes so GitHub does not keep showing an older collection.

Use the shared SVG sources in internal/branding/icons. Blue is the dominant
colour; mint and pale neutrals are secondary accents. Keep recognisable symbols,
consistent visual weight, transparent margins and legibility at small sizes.
Regenerate PNG/ICO exports from the SVGs with sync-branding.cjs --local-only;
do not hand-edit generated artwork or silently copy into sibling repositories.

## Build and preview cleanup

Use I:/Projects/tooling for large build artifacts, downloads and cleanup archives.
When a build needs NTFS, use one named C:/Projects/.cache build folder per app
and remove its disposable contents after the task. Preserve source changes first;
never treat the live installation or unknown private data as a disposable cache.

Track each test installation's program, data and file roots and all child PIDs.
Stop task-owned preview services and tunnels when verification finishes. Verify
that they have stopped; do not leave older preview suites running alongside the
user's active installation. Check C: free space before large builds/install tests.

Before recursive cleanup, verify absolute paths and active process references.
Do not follow junctions into other folders. Archive unknown test state or older
rollback data before removal, verify the archive, and encrypt private archives
when the destination cannot preserve Windows ACLs. Keep the current installation
and newest usable rollback backup. Do not clear unrelated projects or global temp.

## Verify deployed branding

For an installed-app branding request, publishing artwork or README changes is
not completion. Rebuild the manager or affected app binaries, install with the
normal backed-up update flow and verify the running deployment serves the new
content-addressed icon bytes. Distinguish source artwork from deployed artwork
in the completion report.

## Shared app navigation

Never inject a floating Workspace launcher into app content. It can cover phone
bottom tabs and actions. Use existing app header Apps menus/triggers and the shared
picker event; retain keyboard navigation. Scrolling, route changes, hidden headers
and apps without an explicit trigger must never create an overlay button.
