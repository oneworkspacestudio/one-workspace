# One Authenticator managed delivery

User authorized brief implementation, push, release and installation. No Actions
or subagents. Isolated source branch starts f8d91d5 and fast-forwards public108
main e5a9776, preserving latest Remote7 pin/channel/history. Original working tree
has unrelated changes; do not stage it or the unborn outer workspace.

Adds35th managed app authenticator to registry, Administration default selection,
portbase+37, private runtime/storage/update/recovery, branding/overview/native
Android fallback and public signed-out access. Strict static GET/HEAD allowlist
preserves independent PWA bootstrap; generic proxy does not mutate offline shell.
Normal managed builder validates official seed103 and public Remote7/Tickets1,
replaces paired MyAccounts executable, adds auth package; candidate mode is local
only and publisher rejects it. Integration6 checks synchronized across builder,
catalog and publisher. New schema2 component metadata requires manager109 and
oneaccount-mfa-v1 / oneaccount-offline-totp-v1 capabilities.

Portal sessions retain provider access token; introspection binds exact subject
and fails closed each request. Stream watcher revalidates every5sec (5sectimeout),
cancels reverse proxy context including WebSocket/SSE. Synthetic actual upgraded
socket/eventstream tests pass. Full catalog/distribution/suite tests pass with
GOTMPDIR and TEMP on task NTFS. Initial full-suite run put test fixtures on exFAT
and correctly refused protected storage; setting only TEMP was insufficient when
GOTMPDIR remained exFAT. No production bypass introduced.

Paired app source main22120a7; provider926137b main pushed, preview tags reserved.
Both exact components passed tests/builds. Native adversarial pending/replay/
recovery/rate/SSO/token invalidation and encrypted backup/restore pass, as does real
settings UI provisioning independent offline phone PWA with host stopped. Synthetic
camera QR and media shutdown passed. Physical phone/reboot/clean-host remain open.

Normal35-app candidate build and actual managed fresh/upgrade/recovery acceptance
are in progress. Local previously named seed cache was incomplete; builder rejected
its size before trusting it. Download official exact seed with digest verification.
Artifacts in I:/Projects/tooling/one-authenticator-20261030. Only named own C task
scratch is a fixture; owner C:/tmp/workspace is never used for testing.

Next: complete installer tests, freeze/push source and rebuild from public pins;
use mandatory scripts/publish-release.ps1 with existing protected permanent signing
configuration, signed Android11 and all package hashes. Public component drafts
must be hash-verified before publishing. Then public channel/anonymous checks and
supported backed owner update. Preserve existing public108 recovery and credentials;
no manual standalone shortcut is delivery. Other installed PCs/physical phones
are not claimed. Update this handoff with exact releases and per-target receipts.
