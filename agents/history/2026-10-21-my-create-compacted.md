# My Create / Workspace57

User authorized implementation, GitHub publication, local installation and release.
New MyCreate.exe uses apphost/OIDC, NTFS private data, embedded AUI/Fabric editor,
SVG/PPTX/PDF/PSD/raster adapters, immutable source/revisions, comparison, packs,
local Windows OCR and browser-user My Drive access. Official Canva boundary and
format limits are in app docs/RESEARCH.md; no companion or catalogue scraping.

Rebased this checkout on eea959c (verified installed blue-branding56). Added
create to port mapping (base+27), sessions, navigation, runtime storage, package
ownership, updater group, organization pin, installer scripts and shared blue
branding. All27 logo directory entries regenerated and visually checked.

App source/release: oneworkspacestudio/my-create, v0.1.0-preview.1, compiled
2687e8a. Separate my-create-templates v1.0.0 contains three original CC0 packs,
no private/user-imported content. GitHub Actions disabled. Public asset hashes
are checked against local release outputs before publication/install.

Tests: app Go/API, frontend build/unit/audit, ten real-format fixture/pack checks,
native executable login/edit/save/reopen/undo/export/mobile/OCR. Suite targeted
Go suites (suite/distribution/branding/apphost) pass. Go1.27 t.TempDir honours
GOTMPDIR, so use C:/Projects/.cache/my-create-20261020/suite-tests for private
fixtures; I is exFAT. Fixed a pre-existing index-based updater test made brittle
by inserting the new group; tests now look up Notes by ID.

Workspace57 bundle has25 apps/124 hash-verified payloads. All24 existing app
contracts, versions and executable bytes are unchanged from verified local56.
Synthetic56-to57 normal backed-up upgrade passed: retained identity/paths/app
selection and portal sign-in; Create provisioned at custom base+27; MyAccounts
owner sign-in, template save, same-subject Drive new-copy write and content-hash
icon URL all verified. Processes stopped, owner never used as a fixture.

Artifacts/build tools: I:/Projects/tooling/my-create-20261020. Synthetic private
roots: C:/Projects/.cache/my-create-20261020. Powershell startup stalled midway;
explicit C Node/Python executables continued successfully. I remains readable.
Workspace57 is public from a401b4e; all release assets verified by GitHub SHA-256.
Normal backed-up owner56-to57 update completed: 25 apps, all 24 service entries
healthy, all 24 previous app versions/packages and owner identity/paths preserved.
Installed manager/Create hashes match public artifacts; rollback manifest verified.
Local and public content-addressed My Create icon bytes match source. Evidence
is in docs/evidence/my-create-deployment-20261021.json and workspace57-publication.json.
Installed My Chat preview9 was retained through independent update selection.
Release tags/assets remain immutable. Preserve the current owner installation and
newest usable rollback backup. No clean-host reboot or production acceptance claim.

Operational note: stale global Git safe-directory entries for unavailable drives
can stall Git when working on exFAT. Per-command GIT_CONFIG_GLOBAL=NUL and
GIT_CONFIG_NOSYSTEM=1 with this exact safe.directory and explicit author identity
worked; no global configuration was changed. Task-owned blocked diagnostic
processes received termination requests; some Windows I/O waits did not exit.
These are separate from the stopped synthetic app services.
