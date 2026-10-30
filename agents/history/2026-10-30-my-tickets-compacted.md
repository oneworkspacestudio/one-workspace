# 2026-10-30 My Tickets managed integration

User authorizes implementation, source push, release and installation on this PC.
VM requirement explicitly waived. No Actions. Owner installation is deployment,
never a test fixture. Keep private state/credentials out of source and evidence.

## Implementation
Tickets is General/default, /tickets, port9026/rebased+36, foundation OIDC,
app-owned updates/recovery/private storage. Native in-process apphost SDK host
embeds upstream Vikunja2.7.1 plus custom UI/API, no helper process/private listener.
SDK438b9cfe adds custom UI, consistent extra DB snapshots and checked recovery.
App source8650cb8bacbef64c7a6a6c4b86bcc12f58dc5b20, public org/my-tickets,
Actions disabled, archived inherited workflows. Pin preview1 ZIP770b37f590caa23c6d1a732156658556fd9f080b25e7819d6fe6a7fc3a11e536,
EXEf112e9635c179ac923458e75de1c6e0a5f5c2b03ca20a7c81fb39d04446e3b54.
Manager105 integration5 protects native storage/identity contract. Normal builder
preserves verified public103 baseline, adds Remote5 and Tickets. Candidate option
accepts local SHA-pinned app only for acceptance; publisher rejects candidate:true.

## Verification
Five Workspace Go packages PASS; complete catalog34 apps307 packages validates;
omission/min-manager/explicit exclusions covered. Real fresh105 installer and
104-to-105 upgrade PASS with synthetic NTFS private data on this PC: MyAccounts
once, Tickets opens without second password, project/issue CRUD, desktop/mobile,
exclusions retained. Independent synthetic Tickets1-to-2 update and normal verified
rollback PASS, existing issue/user file/manager/accounts hashes preserved. Synthetic2
never published. Native browser/API/encrypted two-DB+attachment restore PASS bound
to final EXE hash;16 concurrent unique issue keys;10K search159-170ms. Backend lint0,
frontend type/style/build pass. PE requires Windows system DLLs only.

Real managed tests found SDK correctly strips unverified email but upstream user
factory required it. App now allows absent email only for stable authenticated
external issuer/subject; no fake email or local validation weakening. Actual OIDC
passes. Fixed SQLite pool bounds (upstream ignored config), mounted CORS origin,
and rich-text phone toolbar overflow. Tests are regression coverage.

## Limits and remaining delivery
Evaluation preview: protected shared-filesystem Drive source, remote Drive API
backend absent; no durable webhook queue or scheduled due-date automation; Git
resync latest50 commits. Live owner Gitea, elevated SCM, reboot/physical phone not
claimed. See docs/MY_TICKETS.md and app ACCEPTANCE.md.
App release/normal noncandidate suite105 build, mandatory signed Android publisher,
public channel, backed-up owner104-to-105 install and installed byte/health receipt
must be completed before delivery is claimed. Artifacts/caches on I task tooling;
small NTFS fixtures on C. No owner data used as test fixtures.
