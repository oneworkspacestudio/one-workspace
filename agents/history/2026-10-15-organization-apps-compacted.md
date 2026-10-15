# Complete organization app installation

User requests every organization app included with Workspace and sourced from
oneworkspacestudio. Repository main was clean but behind the published Calendar
integration by two commits; fast-forwarded to 90fb522 before changes. Preserve
legacy hosting-capable Mail and concurrent component development. No Actions.

## Package and behavior

The organization has nine repositories. Notes is an empty repository with no
published executable and is being built in a separate task. It is explicitly
pending, not represented by a fake installed tile. All released applications
are included: Accounts, Drive/Office, Photos, Photo Editor, Git, Music, Mail,
Calendar and Contacts (ten catalog entries, eight supervised app processes).

`scripts/rebuild.ps1` now uses `build-org.ps1` and the new organization component
lock. It verifies published asset URLs, GitHub SHA256/size and cached downloads;
includes licenses and an exact organization-components.json provenance record.
Every app source points at an organization release. Drive/Photos retain the newer
upload repairs already released in org Workspace16; the older EXE in org Drive
must not replace those fixes. Matching source remains linked from org Workspace15.
Studio.exe is packaged for standalone editing; its integrated tile opens the
existing Drive editor with Drive's MyAccounts/file permissions. No new Studio
OIDC adapter is invented. Component source links in active docs now use the org.

Calendar/Contacts share a backend. Published Mail15 stays private via explicit
WORKSPACE_MAIL_BACKEND=private, retaining SMTP hosting, provisioning, identity and
mailboxes. This intentionally uses two isolated engines. Fully shared Mail still
requires a planned migration; do not silently use the separate Calendar Mail EXE.

New `Workspace.exe --data PATH install-apps --bundle ZIP` adds missing suite apps
to a stopped installation, registers OIDC, verifies before extraction, copies
only new/changed packages, updates the manager and saves protected backups. It
refuses changing an installed process binary with different bytes and preserves
extensions, user files and stores. Existing working executable hashes were
unchanged; the manager was explicitly updated.

## Verification and installation

Full Go tests and vet pass. PowerShell syntax checked. Organization catalog
validation checks all ten entries, sources, licenses and exact shared app args;
positive catalog passes, intentionally split args fail. Fresh protected NTFS
synthetic installation with programs/files on another drive was exercised.
It caught a pre-existing PowerShell argument concatenation bug in Calendar and
Contacts, then a restart bug binding the already-owned persistent shared engine.
Both fixed. Corrected isolated startup reached100%, all eight processes healthy;
restart reused the original shared instance. Clean host/reboot remains untested.

Actual authorized installation C:/tmp/workspace updated using normal manager
stop, verified child exit, verified bundle and install-apps command, hidden
restart. Backup is under Private/backups/install-apps-20261015T150549.046787100.
All ten entries installed; startup ready100%, all eight apps and connector running.
Public app entry routes reachable (Calendar/Contacts properly redirect anonymous
requests to sign-in). Mail15 configured, incoming reachable, identity/hosting
errors empty. Shared registration contains Calendar/Contacts only. Existing
Accounts, Drive, Photos and Mail binaries unchanged. Manager SHA256
41226fcdab1862a55869d72ae3494b0d0141c755e0c737dff707ae71fa8b3bd5.

Browser automation was unavailable because its request-header policy service
could not load; do not claim a new rendered-browser or app-by-app sign-in test.
Existing Calendar SDK native tests are prior evidence, not a real MyAccounts
end-to-end result for this task. No additional messages or DNS changes made.

## Build storage and resume

C: was nearly full. Only earlier owned release build outputs were moved to
I:/Projects/tooling/workspace-build-archive-20261015 (old public staging dist
and embedded payload). Current build outputs reside at
I:/Projects/tooling/workspace-org-source-final/dist. Raw ZIP entry reuse staged
the final corrected arguments/licenses/manager without recompressing unchanged
components; complete catalog verification passed. Original compressed pre-final
candidate stays in tooling/workspace-org-release-17. Synthetic NTFS fixture was
normally stopped, its owned service stopped, processes checked and private
fixture removed. Health/progress evidence remains outside Git in
tooling/workspace-org-freshcheck. Leave unrelated preview/live installations alone.

Prepare/publish preview17 with source, complete installer, manager, catalog,
bundle, provenance and checksum manifest. Confirm asset sizes/digests before
publishing, then run the anonymous compiled online update check against it.
Notes remains the only missing application pending its own first release.

## Older-installer compatibility correction

After the first publication check, source review proved setup16 could accept the
schema1 catalog but could not provision shared-service identity or remap the new
Calendar ports. Preview17 was immediately returned to draft; its assets/tag were
not overwritten. The 17 download-only check was canceled as superseded, with no
installation changes. Final publication must be preview18.

Catalog now declares installerIntegration2, supported by current code and rejected
by setup16's strict decoder before installation. Future unknown integrations fail
with newest-installer guidance. Unit test reproduces the exact old strict field
set, verifies rejection, verifies integration2 support and rejects integration3.
Full Go tests/vet pass again. Build/catalog validation requires the guard. Existing
installed apps remain functional; update the manager/catalog to18 before final
release verification. Older saved setup files must be replaced once, as documented.