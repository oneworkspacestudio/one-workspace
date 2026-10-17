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

## Preview 30 publication complete

Published https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.30
from source/tag `287e07100ae37c248f73801498e2279cd75cdb50`. Source is pushed to org
main; README and build defaults identify30. Exact committed source was compiled
on isolated NTFS; combined browser JavaScript syntax passed. No new tests, full
progress UI or real update/rollback acceptance is claimed. GitHub Actions stayed
disabled and was not used. Historical29 evidence remains labelled separately.

All13 app definitions and component pins are unchanged from29. All59 package
size/SHA256 checks passed; only manager and Workspace release-status notice
changed. The complete installer embeds the exact verified bundle. All18 uploaded
asset sizes/digests match local files; anonymous discovery selects30, installer
HEAD returns200, and public README identifies30. Prior releases are untouched.

Artifacts, exact source archive, packaging/upload/verification/cleanup scripts
and receipts are retained under local `tooling/workspace-progress-release30`.
Ignored dist and embedded payload match the final public artifacts. Installer
675001856bytes, SHA256
`fd7468851f3539f3acddc8467c05a136503b5b1ac8d6b47af192748a5a95d723`.
Manager SHA256 `562b1a1b5c07f7103e47c93cf4bd05801ff1ec68b39a403cd761977322ae5454`.
Bundle SHA256 `7897181bcabef9e2b2a3f9d7b2bb962a7727b5d9ead9f84cf6de0d4b43428992`.
The exact named NTFS build root was removed after source/artifact preservation
and canonical/reparse/process checks. No unknown data or real backup was removed.

After publication, the idle earlier local progress helper was closed through its
capability API with busy/pending guards. The public30 stage/Workspace.exe opened
the local maintenance window and ran the normal read-only release check. It
completed with no error, offering30 to installed29 and selecting direct
Workspace.exe plus KNOWN-LIMITATIONS.md for the two changed packages:11834313
bytes instead of the665687763byte complete bundle. This is actual release-plan
discovery, not an update/backup/rollback execution or browser UI acceptance.

The published helper remains intentionally running for the owner's next update;
its PID/receipt is in local-window.json. Capability fragments/logs remain in
protected Private storage and must not be printed or uploaded. Installed runtime
stays29 until the owner applies the offered update. No app downtime, user data,
identity, credentials, DNS, autostart or app-selection change occurred. The source
history follow-up records completion without retagging or replacing any asset.

## Preview31 terminal-state display correction

User screenshot showed an indeterminate bar after a completed public30 update.
Root cause: author `progress{display:block}` overrides the browser's normal
[hidden] rule. In addition, the overall stage bar deliberately stayed visible
on completion, and phase/detail repeated the same result. This is a30 UI bug;
read-only live state confirms installed30, not busy, no pending transaction.

Added author `[hidden]{display:none!important}`; both progress bars now render
only while status is running. Terminal states retain steps/history but show a
static result and Completed in/Stopped after duration. Legacy phase hides when
structured progress is present, and visible detail retains status semantics.
This changes presentation only, not transaction/backup/startup logic.

Browser-screenshot skill used on isolated mock-state HTML generated from the
actual page constants. Completed screenshot visibly has no bars/duplicate status;
running and failed DOM snapshots respectively retain/hide bars as intended.
No live app, account, update or rollback was exercised; no test suite was added
or run. Synthetic screenshot/source/server receipt are outside Git under local
tooling/workspace-terminal-release31. Temporary review tab was closed; review
server must be stopped after capture. NTFS build root is the exact named
C:/Projects/.cache/workspace-release31; preserve source/artifacts before cleanup.

User's publication expectation applies to this correction: build/publish31
locally without Actions, keep30 untouched, verify all packages/assets and update
README/build defaults. No unattended live upgrade is authorized. Preserve real
backups despite limited C: space; never prune user snapshots to build a release.
