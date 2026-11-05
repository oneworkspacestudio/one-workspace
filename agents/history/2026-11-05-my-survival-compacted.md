# My Survival integration candidate

Based on `1f505e0` in an isolated detached worktree. Work is local and uncommitted;
no public release, channel promotion, GitHub Actions or owner installation. User
asked to create an offline survival app below 1 GB. The app lives in the sibling
`my-survival` repository; preserve unrelated work in other suite checkouts.

## Contract and changes

Stable ID `survival`, name My Survival, route `/survival`, General default app,
OIDC callback `/survival/auth/callback`, update/recovery namespace `survival`.
One `MySurvival.exe`, apphost health endpoint, no sidecars; private `${DATA}/survival`.
Port 9041 / base+41 deliberately avoids My Commute's observed 9040 reservation.
Proposed installer integration 9 and minimum manager preview.115 need reconciliation
with newer parallel suite work before publication.

Registry, runtime/storage/update ownership, callback and logout inventories,
branding, catalog completeness and ordinary build plumbing now include survival.
Static handbook is deliberately public: GET/HEAD allowlist in `survival.go`, no
query/traversal/API/auth bypass. Private backups retain suite gate plus apphost
per-user authorization and CSRF. Portal HTML injection is skipped to retain exact
offline hashes; app requests shared navigation separately when online.

`scripts/survival_component.py` validates identity, compatibility, hash, namespace
and <=150 MiB package / <=20 MiB offline payload. Normal builds require an immutable
public source pin. Explicit `--survival-candidate` / `-MySurvivalCandidate` allows
local candidate packaging. No fake pin added to organization-release-lock.json.
Publisher expects complete integration-9 inventory and rejects omitted Survival.
Generated branding assets include only intended Survival and suite overview files;
incidental Finance/Authenticator icon regeneration was restored to baseline.

## Evidence and remaining gates

Catalog/suite/branding Go tests passed, including public-handbook sign-in exception,
storage, package ownership and missing-app rejection. Python component boundary
and four publisher tests passed. Frontend/native app verification is recorded in
the app's `docs/ACCEPTANCE.md` and compact history.
The actual final candidate package passed component ownership, hash and size
validation (11 owned files); the normal public-release path correctly rejected
its candidate metadata. App browser tests (six), Go vet and native offline/private
backup/logout acceptance passed. Local preview services were stopped.

Full organization candidate construction, supported managed fresh install, upgrade,
app-only update/recovery, MyAccounts end-to-end SSO and physical-device verification
are not claimed. Before releasing, reconcile current registry/integration version,
pin reviewed component/source hashes, run all lifecycle gates on synthetic isolated
NTFS accounts, then use mandatory signed-Android publication when authorized. Do
not use owner installs as fixtures or replace immutable release assets.

## Authorized installation continuation

User now requested installation in Workspace. Reconciled with public main dba9313
on `codex/my-survival-install`, preserving News/Commute/Finance6 and startup fixes.
Primary target only; no second-PC deployment requested. Waited for live Commute116
update to finish. New guard11 / minimum manager117; Android14 prepared. Corrected
canonical9041 to base+51, matching generic upgrade rebasing; prior base+41 is obsolete.

Expanded checks found missing distribution update-group/capability/package-owner
entries. Added them and a SelectBundleUpdates regression for default discovery,
exclusion, owned payloads, identity and ports. Catalog/distribution/suite/branding/
setupui tests pass. Native app pins public apphost116 with GOWORK disabled.

First complete candidate117:40apps/520packages, every size/hash verified. Isolated
fresh setup, restart, optional MyAccounts backup sign-in, shared picker, encrypted
backup preservation, offline cache and omission passed. Upgrade caught the missing
package-owner rule before owner mutation. Corrected candidate117b normal rebuild
is running. The harness also corrected two assumptions: public guidance does not
force backup sign-in, and shared navigation changes the Apps link to a button role.
No public source/release or owner update by Survival yet. See task tooling handoff
for current sessions, fixtures and prepared publication/deployment scripts.

Corrected candidate117b now passes all six native managed checks, including the
116→117 default-discovery upgrade, independent update and rollback preserving the
encrypted plan, account payload and Files sentinel. Final changed-package vet
passes. Candidate evidence is under docs/evidence/my-survival. App source60be07b
is public and immutable app1 release has13 verified assets; anonymous app ZIP
matches, Actions disabled. Full suite release and primary installation pending.
