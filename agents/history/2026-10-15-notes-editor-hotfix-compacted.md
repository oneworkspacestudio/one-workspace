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
