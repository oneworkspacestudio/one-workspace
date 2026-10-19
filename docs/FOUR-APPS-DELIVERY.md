# Tasks, PDF, Monitor and Search

Four original Apache-2.0 Windows x64 apps join the existing Workspace installer,
launcher, MyAccounts, update groups and backup system in preview.37. End users
need a modern browser and protected NTFS private-data folders; no build tools.

| App | Download and source | Main capability |
| --- | --- | --- |
| My Tasks | [preview.1](https://github.com/oneworkspacestudio/my-tasks/releases/tag/v0.1.0-preview.1) | Private/shared lists, recurrence, reminders, board and ICS |
| My PDF | [preview.1](https://github.com/oneworkspacestudio/my-pdf/releases/tag/v0.1.0-preview.1) | Viewer, queued tools, annotations, true raster redaction and Drive copies |
| My Monitor | [preview.1](https://github.com/oneworkspacestudio/my-monitor/releases/tag/v0.1.0-preview.1) | Owner-only app health, metrics, incidents and optional public status |
| My Search | [preview.1](https://github.com/oneworkspacestudio/my-search/releases/tag/v0.1.0-preview.1) | Permission-aware Drive, Notes, Tasks, Calendar and Contacts search |

Notes [preview.5](https://github.com/oneworkspacestudio/notes/releases/tag/v0.1.0-preview.5)
adds its authenticated current-user identity endpoint. Existing Notes users must
update that independent app group to enable its Search connector.

## Installation and updates

New installs include all four by default; Advanced setup can exclude them.
Manager updates add new defaults while preserving explicit exclusions, existing
ports, user files and identities. Older catalogs retain their inferred omissions.
Existing app upgrades remain independent selections. A compatible older controller
completes additions inside its verified backup/recovery transaction, never on an
ordinary launch. Dedicated per-app private data is covered by Workspace backups;
standalone apps support encrypted consistent backup/new-folder restore. Stopped
suite remove-app/install-apps commands preserve data and manage selection.

## Recorded local verification

All builds/tests ran locally without GitHub Actions, against synthetic accounts
and disposable NTFS data. The owner's installation and files were preserved.

- Each app: TypeScript/AUI, frontend contracts, Go tests/vet, pinned-module native
  build, desktop/mobile Chrome and same-window private-draft clearing. Eight final
  standalone browser checks pass.
- Four encrypted backups restored to new protected folders with matching state
  and invalid old sessions; wrong-password/tamper tests reject archives.
- Native fresh/upgraded suites: one MyAccounts login opens all four plus Drive
  and Notes. Real Drive upload -> PDF import -> save-copy preserves ordinary-file
  bytes. Monitor discovers live services, public status stays off. Search uses
  actual per-user sources, escapes snippets and refreshes moved/deleted records.
- Second-user isolation, viewer-sharing/edit denial/revocation, non-owner Monitor
  denial and clearing all four sessions on logout passed.
- PDF malformed/encrypted/signed fixtures and independent pypdf extraction of
  raster redaction passed; original bytes unchanged. Signed fixture is self-signed,
  without a trusted-signature claim.
- Updater regressions cover new defaults, exclusions, custom ports, new clients
  and data-preserving removal. Native verified updates back up and restart apps.

## Explicit preview limits

Unsigned previews; no second physical clean host, real reboot, Windows-service or
independent security-audit certification. Go source scans found no reachable or
imported-package vulnerabilities in the four apps; module-only findings remain.
New app runtime npm audits reported zero findings. Notes retains the inherited
moderate sprintf-js precision denial-of-service advisory with no upstream fix.

Search does not index Mail bodies, Chat messages or Vault secrets. Tasks has no
invented central mail-notification API or full CalDAV VTODO. PDF annotation and
redaction flatten output and invalidate signatures; drawn signatures are not
cryptographic signatures. No universal Office conversion/OCR is advertised.
Workers have Windows resource limits, not an OS security sandbox. Monitor blocks
private/reserved custom targets and rebinding; domain expiry/email routing are not
implemented. Public status is off by default. See each app's complete security,
API, operations, architecture and acceptance documents.

The original Workspace showcase remains; the app-mark overview and aligned README
table include all four original icons.
