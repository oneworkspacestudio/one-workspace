# My Shield managed integration

User explicitly rejects the standalone launch/account flow and requires a normal
Workspace tile, automatic managed startup and OneAccount integration on the current
PC and workspace-2 remote desktop, not the mini PC. Original push/release/both-PC
installation authorization persists. No subagents or Actions.

Work isolated in I:/Projects/workspace-my-shield, codex/my-shield from e67c301,
current public main including secondary Maps10 delivery. My Shield private repo
sourceabfa542 prepares app3, Go1.27.3, minimum Workspace95, explicit managed update
metadata and managed-first docs. Exact release3 executable SHA21286cb47074409b5bbc5dfc791cd4f151a93ed404423e383fd894b75bfdd43b.
Its private source stays private; public Workspace bundle distributes executable/
notices. Channel publisher now checks repository visibility and omits private
release metadata, yielding an empty app repository row in the public index.

Suite95 wires shield foundation SSO, Owner subject/home env, base+32 port9022,
/shield callback, owner-only catalogue, cookie/session signout, working directory,
backup and independent update ownership, new default builder entry and original
central icon. Managed runtime forces MYSHIELD_LAB=0. Private store never excluded
from verified metadata recovery. No actuator/model/security-policy changes.

Catalogue/distribution/suite/apphost tests pass with NTFS GOTMPDIR. Initial test run
used exFAT build temp and correctly failed protected-storage checks; no product
workaround. Package-wide gofmt-only unrelated changes were restored in this clean
owned checkout. App3 Go tests/vet and web checks pass; govulncheck reports no
reachable vulnerabilities (three advisories only in unused required modules).

Fresh C:/Projects/.cache/my-shield-integration-20261028/Fixture uses only generated
owner/member and verified public MyAccounts binary. Source/build/evidence on
I:/Projects/tooling/my-shield-integration-20261028. Initial readiness harness used
an unsuitable account-version endpoint and was replaced by health checks; current
fixture healthy. Avoid interrupting its parent exec session: Ctrl-C also stops
its background manager children. Manager now launched independently, PID recorded
in fixture-processes.json. Use normal --data ... stop and verify owned children.

Native exact app3 binary acceptance PASSES: owner signs in once in real MyAccounts,
Shield tile opens without another form, shared Apps picker works, 390px no overflow,
member has no tile and native owner API rejects member, anonymous401 and local
owner creation403. Workspace logout requires new sign-in and fixture SQLite read
confirms zero remaining owner sessions while member session stays separate. Only
safe result counts and model/catalog screenshots are published; no event payloads,
credentials or tokens. Native receipt docs/evidence/shield/ identifies exact hashes.
No live owner data/identity/Drive as fixtures. Monitor-only, neither model qualifies.

Private app3 is published with eight digest/size-verified assets at its private
v0.1.0-preview.3 release. Shipping suite manager bytes equal the native-tested
candidate (SHA997d7d73ba68f3396362dbc4398d9aacfb11767c0e717e7e16d610db15b9cbb9).

Pinned full org bundle95 built; every catalogue launch contract passes. Scoped
suite vet passes. Fresh native protected94-to95 Workspace-only upgrade passes:
Shield added automatically, owner/Accounts bytes retained, opt-outs preserved,
verified recovery and transaction complete. Receipt in docs/evidence/shield/.

Pending: signed Android5 through required
scripts/publish-release.ps1, public channel and anonymous verification, normal both
owner-PC updates preserving all other app bytes/identities, actual owner SSO/tiles,
stop legacy standalone previews only after managed health. Preserve their private
data and primary encrypted rollback; do not silently merge local identities. No
managed live deployment is claimed yet. C free about10.46GiB before full packaging;
large artifacts/caches I, only the named protected test data on C. Preserve newest
usable owner recovery and unrelated worktrees.
