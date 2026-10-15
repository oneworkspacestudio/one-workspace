# Shared Stalwart for Calendar, Mail and Contacts

Calendar and Contacts source and standalone Windows programs: https://github.com/oneworkspacestudio/calendar.

Build locally using scripts/build.ps1 with -Calendar, -Contacts, -Services (WorkspaceServices.exe), -Stalwart (official pinned stalwart.exe), and an updated -Mail. Calendar can be included without Mail. No GitHub Actions are used.

The service is installed once under `${DATA}/shared`; app programs and preferences remain separate. The manager verifies the engine and manager package hashes, registers one shared MyAccounts client with the three exact callback paths, then starts the shared backend before the first dependent app. The shared client uses immutable `sub` mapping. Calendar uses CalDAV, Mail JMAP and Contacts CardDAV with each user's own token.

Default listeners: private HTTP only; adding Mail does not open SMTP, IMAP or POP3, change DNS, or reset storage. Legacy private Mail data is not automatically migrated/adopted. A dedicated administrator migration is required. The old published bundle remains compatible; the new mode is selected only when the two shared packages are present.

New shared-mode Mail has no recovery administrator credential, does not own/start/stop Stalwart, and cannot edit its backend through the legacy hosting wizard. Stop individual apps freely. Use the standalone shared manager's explicit stop/backup/upgrade commands for maintenance. Uninstalling an app only removes its shared registration; shared data and other apps remain.

Existing live installations were not changed. See Calendar's technical guide and verification report for native install-order, OIDC, ACL, recovery and app tests. Windows logon/service reboot and a clean machine are separate acceptance gates; this remains an unsigned preview.

## Hosting-compatible organization bundle

The complete organization build marks the published hosting-capable Mail with
`WORKSPACE_MAIL_BACKEND=private`. Calendar/Contacts share their own service, while
Mail retains its existing SMTP hosting, mailbox provisioning and private store.
This intentionally runs two isolated engines. Merely adding Calendar never
adopts legacy Mail storage, changes its OIDC audience or disables Internet mail.
Fully shared Mail is still available through the separate Calendar preview and
requires a separately planned migration. See [organization apps](ORG-APPS.md).
