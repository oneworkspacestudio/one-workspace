# Notes editor crash hotfix bundle

User explicitly requested fixing the Notes commands-null crash in Workspace,
pushing and releasing. Notes source4c6e422 fixes destroyed-editor effects/callbacks,
keeps separate per-note undo history, adds a regression and bumps preview4. The
regression reproduces the exact reported error against published Notes preview3,
and passes against4. All6 browser workflows, Go tests/vet, frontend/type build and
portable restart/attachments/backup/restore/prefix checks pass in disposable state.

Suite source uses isolated workspace-notes-fix worktree from public c0a0df5 so
primary release artifacts and unrelated edits are preserved. Only Notes pin/hash
and build defaults change. Bundle version0.3.0-preview.21 avoids reusing the prior
prepared preview20 number. All other org pins stay unchanged, including the old
Music executable and explicit rejected-hash guard for its withheld newer build.
No Defender protection change, GitHub Actions or live installation update.

Notes preview4 is published first with exact source/licenses/checksums; build-org
fetches that published component and verifies the entire pinned org set. Rebuild
the full installer/manager/catalog/bundle and publish a new prerelease with hashes
and source provenance. Do not rewrite old tags/assets. Existing live stores,
identity configuration and active suite processes are untouched. Source/payloads
on I:, temporary isolated NTFS fixtures under the named Notes test cache only.

Workspace Go tests pass. Final bundle/catalog/download checks and release receipt
will be appended after publication. Do not infer clean-host/service/reboot or new
public provider/integration acceptance from these tests; existing release limits
in docs/RELEASE-STATUS.md remain.

## Published Workspace hotfix

Published v0.3.0-preview.21 from sourcec70d9e793ba802ad33b2ba246822442fd71e26b2:
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.21

All11 uploaded asset sizes/digests matched local SHA256 before publishing.
Installer public HEAD returned200. All43 inner package files match the bundle
catalog, all11 app records/arguments pass organization checks, and the Notes
package equals the verified standalone4 EXE. Exact Workspace source, full
installer/manager/bundle/catalog, component provenance, notices and checksums
are attached. Builds are local; Actions remains disabled in both repositories.

Installer SHA2568738e67d1efa40f564812eeb83097db3119ca28a4915523b35386801c2d1fe05
Bundle SHA2569110139d5d22abead18c28c77044d03b4a763cecf7630764fffcd3dd1b58e5ce
Notes ZIP SHA256e7ec8c75801d5ba0a8cc9c459bcb523783f671451519fcd44953bc3b1a81584c

A newer source-notice commit5d4221f appeared during work; it was preserved by
normal rebase before the source/build/tag were finalized. Primary checkout and
release20 work were not reset/overwritten. Later history-only commits do not
alter immutable release tags or their matching sources. Existing live Workspace
and Notes processes/stores remain unchanged; installing the fix is a separate
normal owner-controlled upgrade. Disposable NTFS Notes fixtures removed after
process/path/reparse checks. No test preview/tunnel remains from this task.

The online selector currently orders publication timestamps, not semantic tags.
Do not later publish an older prepared Notes3 bundle after this hotfix; re-pin
Notes4 or leave it unpublished, to retain this fix in normal online installs.
Anonymous installer download verification is recorded below after completion.

Anonymous final installer check passed (exit0): downloaded and verified the full
573MiB public bundle, selected Workspace0.3.0-preview.21 and reported Notes
0.1.0-preview.4 among all11 apps. No GH_TOKEN/GITHUB_TOKEN or CLI login used.
Mail's independent organization check retained verified preview16. Download
scratch files were removed by Bundle.Close. This proves the real online setup
resolution/download/hash path, not a new GUI installation or live data upgrade.
