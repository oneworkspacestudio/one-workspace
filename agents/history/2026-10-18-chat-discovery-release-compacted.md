# Chat discovery integration / preview35

Owner requested tester accounts, discovery/ID/messaging fixes, push and releases.
Used repository-history and an isolated protected NTFS two-app fixture, never the
owner's live suite as a development fixture. Chat repo handoff contains details.
Four synthetic example.test members; Alice/Bob independent browser contexts.

Manager now passes a per-run Chat directory token and polls authoritative
MyAccounts users every15seconds. Full snapshots only after all pages succeed;
issuer/subject IDs and verified email only, explicit disabled flags, limit10,000.
Private Chat endpoint never exposes manager credentials to browsers. New members
are searchable before their first login. User data/keys and other component pins
remain unchanged. Chat4 fixes QR constructor, native ID guidance, search, media
base path, WAL readiness, unloaded-contact presence and private packet logging.
Go suite snapshot validation/verification tests and Chat gateway/ID checks passed.
Browser OIDC, discovery, new-account sync, native/QR IDs, direct/group messaging,
history reload, exact300KB attachment bytes, anonymous/sandbox checks and disabled
account/session revocation passed in the fixture. No production/recovery claim.

TargetChatv0.1.0-preview.4 source32d9bf3bf73dd5e577bdd81ca28856c4bea6b636,
ZIP SHA256e3b54436732020da35f72f23e5996a8f3420884ec7a69a5320885edce7b83d3d.
TargetWorkspacev0.3.0-preview.35. Artifacts tooling/chat-build-preview4 and
workspace-chat-release35. Publication/final package verification not yet confirmed
at this edit; append results. Keep older immutable tags/assets and two unrelated
older untracked histories. Stop task-owned fixture/processes and clean named C
build/fixture roots after recording final checks. Owner data/backups untouched.

## Public delivery and cleanup confirmed

Published Chatv0.1.0-preview.4 from32d9bf3bf73dd5e577bdd81ca28856c4bea6b636
with12 verified assets and Workspacev0.3.0-preview.35 from
2e1cb9cc1b02ca0f8af08141179647b87aaa9f98 with19 verified assets. Anonymous
metadata/download HEADs and every GitHub size/SHA256 digest match receipts.
Bundle SHA25621a48b53d07588cab4ee0e91885f4f4767afa68260d69d6b4668b90074c23d24;
installer SHA256708fcc094eb744ee13a7ca934cd2d9505bfa170c5a4021f1d80b1ab669bd262b.
All14 entries/63 packages verified; only manager, Chat, Chat dependency notice and
the two suite provenance/status notices changed. Other app executable pins retained.

The exact release4 EXE passed the browser matrix again. Official packaged updater
upgraded a stopped fixture with verified original34/Chat3 executable bytes to35/4,
created/verified a backup, restarted successfully, retained synthetic messages and
the two-app selection. Post-upgrade browser matrix passed. Additional checks:
Workspace login then Chat tile skips another password form; disabling a member
closes their existing WebSocket; Chat logout closes its socket and rejects access.
Public installer update detection resolves exactly Workspace35 and Chat4 as current.
Git source/archive differences checked were checkout CRLF normalization only.

Artifacts/receipts/scripts remain on I under chat-build-preview4,
workspace-chat-release35 and chat-acceptance-20261018. Verified and removed named C
browser build cache. Stopped the fixture through its own manager and verified no
fixture executable descendants remain; removed its dedicated synthetic-data root
and temporary backups after retaining non-secret receipts/screenshots on I. No
fixture/tunnel remains. Owner live installation/data/backups were not modified;
it still needs the normal Workspace+Chat update. Read-only sizing showed its full
verified backup needs about6GB free on its local NTFS backup volume, more than
available during this task. Do not bypass its normal backup checks. No Actions.
These completion notes follow the immutable source tags; no binaries retagged.
