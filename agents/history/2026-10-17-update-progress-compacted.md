# Live update progress

## Request and cause

User reported a disabled button stuck on `Downloading verified release for
notice-Workspace-release-status.md`. That package came from the complete bundle;
map iteration chose an internal notice ID for the download label. No transfer
callback existed in the managed-update downloader. This was poor feedback rather
than a confirmed hang. Read-only inspection found the live update had completed
on public Preview29, healthy and with no pending transaction.

## Implementation

- Optional context observer in `internal/operationprogress` carries structured
  stage, detail, completed/total and units. Streaming downloads emit every250ms;
  existing CLI reporters and SHA256/size checks remain in place.
- Distribution stages unique real release assets deterministically and labels
  repository/version, file number and actual asset name. Package extraction/hash
  checks have their own stage. Release checks identify each installed app group.
- Suite emits preflight, backup bytes, backup verification counts, installed
  package counts and native startup checks. Automatic recovery retains the
  observer while ignoring cancellation, as the previous recovery path did.
- Local maintenance API exposes operation status, stages, timestamps, bounded
  activity history and current-stage rate/ETA. Estimates stop when progress stalls;
  no full-update duration is invented. Completion follows transaction and native
  health success, never merely a100% transfer. Failures preserve error guidance.
- UI adds overall stage checklist, measured-stage bar, elapsed time, approximate
  stage ETA, activity details, inactivity and reconnect messages. Polls do not
  overlap; control/capability/origin/busy guards remain. Docs explain the scope.

## Local artifact and verification scope

Built on isolated NTFS and Node syntax-checked the combined browser script.
No tests were added or run, no real update/rollback was replayed, and no full
progress UI acceptance is claimed. All source remains uncommitted/unpublished.
No GitHub Actions or release change. Public29 remains unchanged.

Unpublished helper version `0.3.0-preview.29+progress.1` is preserved under local
`tooling/workspace-update-progress/WorkspaceUpdater.exe`; SHA256:
`832c8ef52c85230719c7a476faeea51a549e9edef6b005d5c818b6722de95db0`.
The previous idle maintenance window was closed through its capability API only
after verifying not busy and no pending transaction. New helper opened the local
update window against the same installation; its state endpoint is available.
It is intentionally left running for the user's next update. Capability/logs
stay in protected Private storage; never print/upload the window fragment.

Installed manager, apps, ordinary files, accounts, backups, domains, autostart and
app selection were not changed. Installed manager is public29 now (the user's
update completed during this request); live Vault remains absent from the older
12-app selection. The new window reopens through existing manager discovery while
it runs; a later restart/closed helper uses the published manager until these
changes are released. Future publication needs a fresh full versioned bundle and
installer, explicit authorization, and must not replace29 assets. Preserve older
untracked installed-preview22/runtime-restart histories.

Source snapshot/artifact receipts are retained outside Git in the same tooling
folder. Cleanup removes only `C:/Projects/.cache/workspace-update-progress` after
source preservation and canonical/reparse/process checks; retain the I: helper.

## Publication requested

User requested seamless publication of the completed fix rather than a separate
manual release request. Publish these changes as Preview30 in the org, built
locally without Actions. README/build defaults/status identify30; older assets
must remain untouched. Small Workspace manager/license/notice/provenance packages
now select direct published assets only if catalog size and SHA256 match; the
verified bundle remains fallback. This removes the notice-triggered full-bundle
transfer for supporting updaters. No component app binary changes. No new tests
are requested or run; compile, JavaScript syntax and packaging checks are the
release verification scope. The previous broader Preview29 checks must remain
labelled historical. Publishing is not permission for unattended live upgrades.
