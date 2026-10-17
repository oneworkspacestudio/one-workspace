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
