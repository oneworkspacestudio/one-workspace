# Dedicated Vault hosting and RC3/preview32 publication

## Scope

Owner explicitly requested My Vault behave like the remotely hosted suite apps,
with all pending changes pushed and released without GitHub Actions. Targets are
My Vault v0.2.0-rc.3 and One Workspace v0.3.0-preview.32. This supersedes the prior
local-only publication boundary; previous local handoff/evidence remains historical.

## Implementation and invariants

Added `vault_hosting.go`: managed Cloudflare startup chooses an unused first-level
Vault hostname on this installation's existing tunnel, persists an intent, appends
only its gateway route before the catch-all, creates its own proxied CNAME and
saves the dedicated origin in existing Vault settings. Other records/routes are
preserved. A pending intent resumes on restart. After local rollback, an existing
address can be adopted only if both its DNS target and exact Vault tunnel/gateway
route match this installation. No unrelated occupied hostname is overwritten.

Ready addresses persist across updates without repeated provisioning. Manual
public origins and owner-pending origin changes are respected. New Cloudflare
installations keep the hosting grant in protected DPAPI Secrets.HostingToken;
existing domain-hosted installations can reuse their already-retained grant.
No credentials enter browser state, logs, source or release payloads. Private
resource metadata is validated against the connected zone/account.

Startup exposes provisioning/remote HTTPS progress. A remote reachability failure
retains local independent access and produces a clear attention landing instead
of an unconditional redirect. The manager configures the dedicated native OIDC
callback after origin selection. Component RC3 contains the preceding canonical
SSO/JIT fixes and limits local recovery to independent login when public OIDC
browser binding belongs to another origin. No encryption/engine dependency change.

Updated README, hosting/sign-in/security guides, build defaults and exact RC3 org
pin. Existing README images/logos are retained. Ordinary app selection, folders,
identity credentials and native encrypted state must be preserved on deployment.
The owner still chooses/unlocks their encryption master password themselves.

## Build and publication plan/evidence boundary

Remote main refs matched local bases before source changes; Actions confirmed
disabled in both repos. Read-only live Cloudflare preflight confirmed the existing
managed connection and an unused preferred Vault hostname. The before snapshot of
DNS/ingress is protected in private maintenance storage, not in Git/artifacts.

Component source is committed as732e782a8ae26f073c3d11701368a5adfaecea44.
RC3 is built from its exact frozen source; original engine and client JS unchanged.
Windows ZIP SHA2569ff7dca32503a8354044b3054a100f1498704c0c64ccdceaea35a9c14bf27bdf;
MyVault.exe SHA25615ece2eeffe454224ba9ec24c7b4c4ac283d24441f1fe4e81706007984f12ebf.
No new automated tests were added/run. Earlier RC2 recovery/client evidence stays
labelled historical, and productionEligible remains false. Build/packaging,
publication, actual remote deployment and cleanup receipts will be appended after
completion; planning or an upload attempt is not a passed release.

Artifacts/scripts outside Git: local tooling/my-vault/release-rc3-build and
tooling/workspace-vault-release32. Preserve real backups and older untracked
histories. Use only the exact named disposable NTFS workspace-release32 build root
when needed, then preserve source/artifacts and check paths/processes before cleanup.

## Completed publication and owner deployment

Published locally built prereleases (Actions remain disabled):
- [Workspace preview32](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.32), source fbb4294cd2baaf07584d5a58ceec93ae9217d80b, 18 assets.
- [Vault RC3](https://github.com/oneworkspacestudio/my-vault/releases/tag/v0.2.0-rc.3), source732e782a8ae26f073c3d11701368a5adfaecea44, 26 assets.
All asset sizes and GitHub SHA256 digests match packaging receipts. Anonymous
release metadata and main Windows downloads return200. Installer embeds the exact
verified13-app/59-package bundle. Workspace.exe SHA256
cd10c28aaf0c721fc890b1a8061ddd3400d1a3ea4a49e1f821a244ffbd647382;
setup SHA256c5067bf4a2223d168fb1489107183f2c03a117bc9d4d3f5ffb56cc4f3843352f.

Initial installer compile exhausted C temporary storage; retried successfully with
GOTMPDIR/TEMP/TMP on I. Preserved source/artifacts and safely removed only the
named NTFS workspace-release32 build root (652MiB). A read-only inactive software
cache plan was made but no cache deletion executed. Existing user data/backups
were retained; latest C free about5GiB. Do not attribute recovered space to an
unverified cause or delete unrelated caches/snapshots.

Owner deployment used coordinated stopped-state replacement with a newly created
and verified encrypted native Vault backup on I, machine-DPAPI backup key only in
protected private maintenance, cold Accounts/settings backup and changed-program
copies. Prior verified full-suite snapshot remains retained. This was not a new
native full-suite update transaction/clean-install acceptance. First attempt
rolled back because local upstream DNS still cached the old pre-provisioning
address; authoritative Google/Cloudflare DNS already had the new proxied record.
After cache expiry a separate attempt made fresh backups and completed. Do not
blindly rerun one-use deployment scripts.

Installed preview32/RC3 runs with local, public Workspace and dedicated public
Vault health200/ok; Workspace launcher303 points to its saved permanent origin.
Suite logical credentials/settings,13 apps and local certificate are unchanged.
Protected before/after Cloudflare comparison confirms all90 earlier DNS records,
all3 earlier ingress entries and unrelated tunnel configuration unchanged; only
one owned proxied CNAME and Vault route were added. Address survives updates.

Normal owner Chrome launch from authenticated Workspace reached remote native
#/lock with no additional account login/security warning. Owner master-password
unlock is still required; it was neither read nor entered. Screenshot retained
only in private maintenance; browser left as owner handoff. No encrypted-item,
clean-host, service/reboot, full recovery or new automated acceptance claimed.
Existing Mail network warning remains unrelated. Original18935 fixture stays off.

Public README retains logos/showcase and describes managed remote hosting. A
post-publication README cleanup links preview32 explicitly and removes a malformed
historical paragraph. Follow-up documentation commits do not retag release source
or replace immutable published assets. Completion receipts/scripts remain under
local tooling/workspace-vault-release32; encrypted private archives stay outside
release stage/Git. No credentials, personal screenshots or full DNS data in Git.
