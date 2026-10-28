# Managed My Shield delivery, 28 October 2026

My Shield preview 3 is installed through the normal protected Workspace update on
the current PC and the confirmed remote desktop. The mini PC was not changed.
Both run Workspace preview 95 and the same hash-verified My Shield executable.
Each reports all 29 services running after the update.

| Check | Current PC | Remote desktop |
|---|---|---|
| Previous manager | preview94 | preview93 |
| Existing app contracts and executable bytes retained | 29 apps | 28 apps |
| Existing owner, settings and credentials retained | Pass | Pass |
| Non-Workspace package records retained | Pass; 139 total catalogue records checked | Pass; 129 total catalogue records checked |
| Verified recovery manifest | Pass | Pass |
| Native OneAccount enabled; separate local setup disabled | Pass | Pass |
| Actual owner tile and browser SSO | Pending expired owner-session sign-in | Pass; no separate account form |
| Shared Apps menu in actual deployment | Covered by native acceptance; owner browser check pending | Pass |
| Deployed central icon bytes | Exact source match | Exact source match |
| Old standalone process stopped; exact Start shortcut archived | Pass | Pass |

The remote preview 93 portal controller stopped at its storage-profile preflight before
making changes. The published preview 95 installer then completed the normal Workspace-only
update. It and its embedded bundle were checked against published SHA256 values;
no storage or recovery guard was relaxed. The expected shared MSVC notice changed
with the Workspace group. Every existing app contract and all package records
owned by other apps remained intact. A verifier initially classified that notice
as an app record; correcting it to match the product's explicit ownership rule
allowed the full comparison to pass before retiring the old launcher.

Earlier standalone data is retained separately; no local identity or telemetry
history was silently imported into OneAccount. Primary encrypted standalone
recovery remains retained. The remote standalone installation originally had no
older snapshot; the managed update now has its own verified suite recovery.

The remote owner opened the actual Workspace tile and reached My Shield through
the existing OneAccount session. Safe Models-view and catalogue screenshots omit
owner names and event payloads. The current PC's expired browser session still
requires its normal OneAccount sign-in; no password was requested in chat and no
owner session was fabricated. Its exact app binary, managed bootstrap, private
client, owner continuity and all services were verified.
Other browsers, including the remote desktop's own Chrome session, still use
their normal OneAccount login. The obsolete standalone setup page was replaced
with the normal Workspace sign-in; no separate Shield credentials are needed.

Both old shortcuts are recoverable under the protected legacy app folder. The
remote installer downloads were removed only after exact hash/process checks;
task scripts and the sanitized receipt were moved off Desktop into protected
deployment staging. Generated native fixtures were stopped and removed after
path, reparse-point and process checks. Live account and Drive data were never
used as development fixtures.

This remains a Monitor-only evaluation. AI is disabled; neither model qualifies.
Full model, MFA step-up, Drive export, clean Windows VM and reboot acceptance
remain incomplete. This delivery does not change those gates.

Published suite: [preview 95](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.95).
All 28 assets, the permanent-key signed Android preview 5 companion and the anonymous public
channel were verified. Private Shield release metadata and withheld releases are
omitted from the public index.
