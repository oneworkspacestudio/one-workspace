# Blue icon family and consumer README

User requested similar blue-dominant icons, a complete Workspace icon overview
and app list, removal of preview-version details below the table, a nontechnical
README, and a durable AGENTS.md rule. Existing source-publication authority
applies. This task updates repository artwork and documentation, not installed
executables or immutable release assets. GitHub Actions remains disabled.

Based on remote main2d6e284 (including concurrent Maps/Meet integration and release
ordering corrections). Work isolated in I:/Projects/one-workspace-readme-branding-20261020
on codex/blue-icons-readme to avoid the pending Clips53 delivery helper. Preserve
the canonical checkout's unrelated untracked histories and background releases.

All26 SVGs now use a common blue family with mint/pale accents. Recognisable
geometry remains; Clips and Backups were redrawn to match. Drive/Music/Hosting,
Tasks/PDF and Vault were adjusted so blue dominates. Every PNG/ICO export was
regenerated locally, including previously missing newer-app exports. No sibling
app checkouts were modified. Historical images remain as prior references.

docs/branding/apps.json is the shared directory for the26-entry README table,
complete content-addressed overview app-suite-f6de3110.png and light/dark gallery.
scripts/app-branding.cjs rejects duplicate/missing/unlisted SVG entries. The
overview renderer regenerates README sections and the gallery from that list.
sync-branding.cjs --local-only generates every icon without copying to sibling
repositories. Existing unqualified PNG paths remain supported for older builds.

README uses the complete current icon collection instead of the old illustrative
hero with obsolete initials. Descriptions, setup, updates and help use plain
language; technical/version histories live in linked guides/release notes.
AGENTS.md now requires consumer-friendly README language, no preview changelog
in the main page, complete synchronized listings, and blue-dominant shared art.
BRANDING.md and ICON-SOURCE.md describe the new palette and owner-authorized
changes, superseding old statements that artwork is byte-for-byte unchanged.

Validation: directory exactly covers all24 catalog apps plus Workspace/Backups;
all26 SVGs render and are primarily blue by coloured-pixel count; PNG sizes/alpha
and seven-size ICO containers verified. Go internal/branding tests pass. README
relative links/images resolve. GitHub GFM markup rendered in Chrome at1440/390px:
one heading,26 app rows, no broken images or page overflow; complete light/dark
gallery includes16/24/32px samples. All screenshots visually inspected. Regenerating
the overview keeps the same hash. No broad app/runtime test or new EXE claimed.

Sanitized review artifacts and scripts: I:/Projects/tooling/blue-branding-20261020.
Task-owned browser/server processes close after capture; remove only the named
C:/Projects/.cache/blue-branding-20261020 cache after checking processes/reparse
points. Publish only this task's docs, icons, generators and AGENTS/history paths.
