# Native component updates and optional compatibility documents

User reports My Social update rejection and authorizes fixing/applying updates
without presenting package decisions to users. Work isolated from unrelated dirty
checkouts in I:/Projects/one-workspace-update-docs, branch codex/updater-package-docs
from current mainf17dd86. Read current AGENTS, release82/MyAI history and public
update channel rules. No agents/Actions. Artifacts under I tooling named
social-update-contract-20261024; synthetic NTFS cache same name on C only.

Cause: published73/80/82 catalogs contain social-contract, but installations that
received earlier standalone Social updates can have only four Social packages.
Previous Social capture validation against fresh public80 therefore missed this
supported legacy layout. Strict manifest validation rejected the extra document
before scopeUpdateChoice could discard an untracked package. Do not modify owner
catalog manually, remove security checks or mutate immutable public assets.

Fix skips only an unregistered same-app -contract at the exact installed
executable's docs/<ExecutableStem>/compatibility.json. It is not installed, fetched
or added to the catalog. Tracked documents remain required, verified and cannot
move. Foreign/runtime/unexpected paths remain rejected. Identical native version,
packages and launch contract with only a release-link change no longer offers
a redundant native update for Chrome-only component releases. Pin native Social6
and build Workspace83 with Go1.27.3 from verified82 bundle, retaining other apps.

Regression tests cover both layouts through check/select/staging, rejection of
foreign/runtime/path changes and tracked-document omission, and identical-native
component releases. Next: complete local checks/native packaging and acceptance,
publish83/promote raw update index, normal backed-up Workspace/Social owner update,
verify exact bytes and unchanged identity/other apps. Owner data is never fixture.

Distribution/maintenance/catalog/platform tests pass with the new regressions and
Go vet passes. Initial full-suite run used exFAT GOTMPDIR and hit NTFS guards in
protected-state tests; rerun those packages with all temp roots on named C NTFS.
The five known Mail probe/provision fixture failures from82 remain baseline.
Release will include only manager, release notice and verified native Social6
contract/payload changes; all other public82 apps retained.
