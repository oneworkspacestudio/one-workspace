# My Hosting and approved C: cleanup

Owner requested implementation of the MyHosting master brief. New original app
is I:/Projects/my-hosting, public org repo oneworkspacestudio/my-hosting (not a fork
after documented upstream fit assessment). Actions disabled. First unfinished
source checkpoint36d469f is pushed. No Hosting release yet; acceptance gates open.
Suite branch codex/my-hosting has uncommitted app/OIDC/logout/webhook/catalog,
build and icon integration. Do not disturb two older untracked installed histories.
Suite README/catalog org-lock/published package pin remain unfinished.

## Storage steering and cleanup completed

Owner explicitly authorized deleting experimental backups rather than archiving,
then directed leaving the live suite on C: for now. This supersedes the ordinary
newest-rollback retention policy for this cleanup only. Do not reformat or repair
E:/I:, create a VHD, or migrate the live installation. Both drives are exFAT,
Health Warning / Full Repair Needed; fsutil confirms dirty. Preserve source and
keep large public build/download caches under I:/Projects/tooling; private bounded
test state requires protected NTFS. No VHD was created; earlier elevation cancelled.

Stopped only task MyHosting fixture (PID124180, loopback19007); its synthetic
browser profile was already closed. Checked exact absolute paths, no reparse
ancestors/children, process/service/scheduled-task references and held the normal
update.lock while deleting. Removed all16 experimental rollback folders, obsolete
C:/tmp/workspace-rollback-20261015T224617Z and synthetic MyHosting/Vault test roots.
19 folders held22.48GiB. Removed another24 old maintenance backup paths/source ZIP,
323.88MiB. Recovery index updates/backups.json is now[]; transaction remains complete.
First index replacement failed because PowerShell coerced null to empty string;
corrected to same-volume atomic File.Move with replacement and verified[] afterward.

C: free0.36GiB ->23.20GiB. Live Apps/Files/Private preserved, local/public /healthz
both statusok. No live services stopped. Windows pagefile25.88GiB and hibernation
12.76GiB unchanged. User Downloads/software runtime archives and global caches
were excluded; they are not obsolete Workspace test backups. Tiny DPAPI recovery
keys for encrypted Vault backups already on I remain protected in maintenance.
Receipts/scripts outside Git: tooling/my-hosting-work/c-{storage-inventory,
cleanup-20261018,maintenance-cleanup-20261018}.json. No new release for cleanup.

## Resume

Continue app acceptance using fresh bounded NTFS fixtures, never owner data.
App Go/AUI dashboard compile; restricted worker host/cross-project/ACL-negative
tests pass after project SID on restricted token/default DACL + private desktop.
Synthetic browser first-run owner/wizard/navigation/theme/mobile checks passed;
screenshots/receipt on I. Initial fixture now intentionally removed, not running.
Subsequent source fixes need recompile/retest. See app history and implementation
document for remaining runtime/deployment/restore/SSO/service/install gates. Do
not label production-ready or publish binaries before critical gates pass.


## Late 18 October corrective handoff

The owner selected free Cloudflare Tunnel with flat first-level hostnames:
my-site-sites.example.com and deploy-ID--my-site-sites.example.com. No router
changes/paid features. Hosting derives the connected root zone from zoneName in
the Cloudflare receipt, with old portal-origin fallback. Local-only empty domain
stays empty; it no longer becomes invalid sites. Configured Hosting-only installs
receive the domain even without Mail selection.

Both repositories have Actions disabled. App latest pushed checkpoint e4827e0;
subsequent accepted implementation is awaiting publication. Suite integration is
still uncommitted on codex/my-hosting. Full suite Go packages pass with GOTMPDIR
on protected C:; Go1.27 t.TempDir uses GOTMPDIR, so using I: correctly failed ACL
checks before the environment was corrected. No protection checks were weakened.

Fresh synthetic installation using actual suite.Install and published org
MyAccounts plus local Hosting/manager builds passed. scripts/native-hosting-
acceptance.go requires a fresh named root and isolates APPDATA. Current fixture
C:/Projects/.cache/myhosting-suite-acceptance-20261018; receipt external at
I:/Projects/tooling/my-hosting-work/suite-acceptance-process.json. Never test with
owner C:/tmp/workspace. Hosting opens through the app tile without another login.
Supported provider GET logout replaces a cross-origin POST MyAccounts rejects;
current Portal session remembers its Hosting path-cookie so root logout revokes
native session. Actual browser logout returned signed-out home/provider login;
remaining Hosting database sessions0. Added unit tests. The shared Cloudflare
Windows named mutex serializes Vault/Hosting read-merge-write operations without
changing their separate public routes. Sidebar footer reserves the app switcher.

App full native suite passes192.843s with actual Next/Nuxt builds, WebSocket echo,
Node crash restart and strict descendant-exit cancellation. Native runtime repair,
retention, canonical routing, snapshot retry, OIDC/security and backup tests pass.
Real static free-tunnel production/preview URLs verified. Operational limits are
explicit: trusted sources, no accepted SCM/ServerCore/reboot/other-machine test.

Next: local app ZIP/notice inventory, real GitHub Express deployment, publish
accepted app preview first; then add verified actual asset pin to org lock,
wire build-org/check-org-catalog to require Hosting, build/test full installer and
publish paired suite preview. Do not invent pins or mark stable. Stop/remove only
tracked test fixtures after acceptance, preserve public demo intentionally or
remove exact owned provider routes before teardown. Preserve unrelated histories.


App v0.1.0-preview.1 is now actually published at the organization release,
source a0375b832feb8acbadba22cfb71456501a67f1de, public ZIP SHA256
392273b43dc97c468e395ab0010276864b6ac70a719dcb2c18d7b09b81ebd129.
REST GitHub asset digest verified before draft publication (old gh view omits
asset.digest; old gh lacks release edit, so used typed REST PATCH draft=false).
Final full app native acceptance passed205.683s after GitHub PAX metadata and
runtime restart ACL fixes; actual GitHub Express public publication passed.

organization-components.lock.json now pins that published exact asset. build-org
requires it, verifies compatibility, copies its executable/notices/docs/inventory;
check-org-catalog requires the fifteenth app and exact OIDC/launch contract.
Source remains preview, no Actions. Next build full bundle, verify all members,
install selected Accounts/Hosting from the actual final ZIP, verify session and
manual uninstall preserving data, generate per-app update metadata, publish paired
suite release. This is not yet a suite release/publication claim.
