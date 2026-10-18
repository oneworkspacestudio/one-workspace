# New default apps in updates

## Request and diagnosis

User updated to preview.36 but My Hosting was absent. Read-only local inspection
confirmed current manager/catalog36 and no Hosting entry. Distribution explicitly
skipped absent app groups, scoped unknown packages away, and never appended apps.
Staged update preflight also rejected every new app as a callback migration.
Installer selection only sent checked IDs from its embedded app list, so a newer
online bundle's new apps were implicitly omitted. No durable opt-out policy existed.

## Implementation

- Installed catalog now records `selection.excluded`. Fresh setup records actual
  exclusions against the displayed embedded catalog before fetching a newer bundle.
  New bundle defaults are included; dependency exclusions propagate.
- Legacy migration preserves missing pre-Hosting apps as omitted; Hosting/future
  IDs are new defaults. Historical explicit omissions cannot otherwise be recovered
  after old managers advanced catalog.version. Do not infer opt-outs from version.
- Shared catalog package-selection helper keeps fresh and upgrade package sets
  aligned (Hosting executable plus three notices; standalone contract unchanged).
- Online/offline update plans discover missing defaults. Choosing Workspace includes
  them; independently selecting a new app includes an available manager update.
  A blocked new integration prevents a misleading partial Workspace success.
- Forward-compatible tagged organization apps may form groups from the verified
  suite catalog; standalone manifests still cannot introduce arbitrary integrations.
- Existing ports, clients, data bindings and extensions remain preserved. New ports
  use the installation base; full target ports are checked before backup/apply.
- New OIDC clients are provisioned after verified backup/program installation and
  before startup. Identity/storage failures use the existing transaction rollback.
- Additional-app CLI also preserves opt-outs and existing custom ports. UI describes
  new apps and locks their selection when bundled into the selected manager update.
- A legacy-controller startup bridge runs only with a matching `testing` journal
  token and missing selection metadata. It fetches the exact installed manager's
  catalog/digest, adds defaults within the original transaction, and refuses any
  change to existing app definitions/packages. Normal startup never installs apps.
  Corrected offline installers do not need this network-dependent legacy bridge.

## Build and deployment results

Final native manager, full bundle and embedded setup built successfully as
`0.3.0-preview.36.1-local`, using pinned organization assets and I: build storage.
Anonymous GitHub API build lookups hit rate limits; an external PowerShell wrapper
used existing authenticated gh metadata, retaining digest checks. A superseded
intermediate build was stopped; do not use its initial manager output. Final files
are in dist; external receipts/logs under tooling/workspace-app-install-fix.
Read-only real-config plan:14→15 apps, only Hosting added on9007, four Hosting
packages, no existing app-definition or port changes. No new tests added/run.
No current-task commit/push/release authorized or performed. Preserve unrelated
untracked October15/16 histories.

Actual local repair completed through `WorkspaceSetup --update ... --apply --apps
workspace,hosting --bundle ...`. Verified full archive before downtime, backed up
and verified existing data/programs, provisioned Hosting client, saved policy and
restarted17/17 startup steps. Transaction complete; installed15 apps; Hosting
running9007. All14 supervised health entries (including Cloudflare) running.
Workspace local health200/ok, direct Hosting health200, public tunneled Workspace
health200/ok. Existing app definitions/ports were unchanged in the read-only plan.
This is an actual owner-install repair, not a synthetic account or file test.
No automated tests, browser owner-login flow or legacy-controller bridge test run.
Do not claim those are accepted; fresh/offline/legacy opt-out regression checks
remain follow-up verification if requested. Final local build is not published.

Verified backup retained: `update-20261018T223640-6f4736b00e3b` under the existing
protected backups root. No existing backups/user data deleted. Build sources and
receipts remain on I:; owner installation stays on C: per user instruction.
The old idle updater was closed through its authenticated local API (Origin is
required on POST; first omitted-Origin request was correctly rejected).

Separate observed startup warning: Mail HTTP server runs, but incoming TCP25
external reachability/UPnP718 mapping needs attention. No router changes made for
this Hosting/updater repair; do not infer Internet email works from healthz.
Standalone public Hosting review demo was left as previously documented pending
the user's keep/stop preference; its19080 Caddy port differs from managed default9080.

Next: source ready for review/publication on explicit authorization; no Actions.
If publishing, choose a new preview greater than36.1-local, build from final source,
and verify the legacy-controller online handoff plus opt-out/fresh-installer cases
in isolated NTFS fixtures if testing is authorized. Never use live MyAccounts or
Drive as test fixtures. Maintain strict standalone managed-package boundaries.
