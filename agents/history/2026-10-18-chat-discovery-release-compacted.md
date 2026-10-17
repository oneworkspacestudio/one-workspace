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
