# Game Studio

Game Studio is an optional Creative & Development app in Workspace. Source is
in the dedicated `codex/game-studio` worktree based on `cb9aacf`, separate
from the dirty Workspace checkout and the unborn outer aggregate. App version
0.1.0-preview.1, local suite candidate 0.3.0-preview.123, installer integration 15.

The owner approves one MyAccounts subject and enables either local DeepSeek R1
via Ollama or opt-in DeepSeek Cloud. The model supplies validated settings for
catch, maze, runner, and memory. No model-provided code runs. Static HTML builds
are immutable and open in a private sandboxed tab. Text chat history stays in
the active browser page. Owner publication commits the selected build to a
private, initialized My Git repo with Actions disabled and uploads the same
HTML to My Hosting using the matching owner's session. The same game/build
identities support publication retries. See `docs/GAME_STUDIO.md`.

Managed source covers catalog/setup, package namespace, origin and rebasing
(+55), MyAccounts callback/shared launcher/logout, private storage, backup,
updates/recovery, branding and ordinary bundle integration. Go tests, vet and
publisher checks passed. Synthetic desktop and 390px mobile browser coverage
included four playable templates, new tabs, revisions, Git retry, Hosting tab
and access denial/revocation. Model, Git and Hosting were test endpoints.
Windows CRLF script hash, AUI UTF-8, and duplicate HTTP listen regressions were
found and fixed. The exact final managed installer passed fresh setup,
restart, optional omission, old-installer guard, 122-to-123 upgrade with
Accounts retained, app-only update and rollback. MyAccounts sign-in, the
shared picker, anonymous denial and app sign-out were verified in the
supported managed paths.

Artifacts and evidence are under `I:/Projects/tooling/game-studio-20261109`.
`managed-verified` is the final ordinary-build output. Bundle contains 44
apps and 923 packages. `bundle-proof.json` records the five Game Studio
payload hashes and `managed-acceptance.json` records the native results.
Earlier candidate folders are diagnostic only. No public release, push,
owner installation, real provider/hosting/account activation, physical phone
or reboot acceptance occurred. A published suite must use the signed Android
publisher and immutable public assets.

## Release preparation

Public Workspace preview.122 was current when publication started. Rebased onto
the public My Commute preview.4 pin and channel updates. The ordinary preview.123
release build uses those public pins and Android companion preview.20. Publish
through the signed Android entry point and check the anonymous update path.
