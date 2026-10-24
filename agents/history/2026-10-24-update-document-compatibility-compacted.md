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

## Delivered

Public83 source09b3b9fb43cbc78897ce4f494b1df2b88d631eb3:
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.83
All release asset digests/lengths and156 final bundle payloads verified. Anonymous
public manager/catalog verified. Main channel promoted at425a055 and checked
anonymously for83/Capture.2, preserving withheld builds. No Actions.

Native isolated82/four-package Social2 ->83/Social6 update passed, retaining
identity, selection, other apps, ordinary bytes and credentials. Real Accounts/OIDC
login and Social Home plus same-data full restart passed. A final packaging check
corrected the embedded organization provenance notice before publication;
manager/Social executable bytes were unchanged after native acceptance. Final
installer/bundle rebuilt and all payload hashes verified. Evidence distinguishes
that metadata-only correction; it does not claim an additional final-installer run.

Normal owner Workspace/Social offline update completed exit0 with a full verified
recovery copy, preserving every other app contract/package/executable and owner
identity/OIDC credentials/paths/ports/selection. Installed manager83 SHA3349f6c6
and Social6 SHA07959a2e match public assets; all26 services run. Cleanup removed
one superseded automatic copy and preserved newest usable recovery. An older
partial-copy cleanup warning remains pre-existing; unknown files were retained.

Old idle82 update controller retired through its authenticated local API and83
window reopened normally. Close returned connection reset while the server exited;
marker/process retirement confirmed before new launch. New controller hash verified
and real release check shows Social6 Up to date, no reason/review link or extra
native update for Chrome-only Capture.2. Runtime/helper intentionally remain active.
Source proof docs/evidence/update-document-compatibility-20261024.json.

Protected NTFS rerun resolves exFAT fixture failures; only the same five baseline
Mail fixtures fail, all relevant distribution/maintenance/update cases and vet pass.
Native fixtures stopped, C native data and I fixture programs removed with bounded
process/reparse checks. Keep release/evidence/tooling and current owner recovery.
Do not use owner data as fixtures, remove registry guards or regress to four-only
release manifests; support both layouts in the validator. Chrome0.1.2 files unchanged.
