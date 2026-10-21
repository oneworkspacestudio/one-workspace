# Installer address choice and explicit reuse

User requested editable automatic Workspace/Accounts addresses and deliberate
override of occupied names. Started from clean main a24dddb (Workspace58 port
repair); source I checkout was older, so used the current user C checkout.

Automatic suggestions remain defaults. Browser can edit both, checks read-only
hosting/review, sees exact current web-record destinations, and confirms replacement.
Edits/reconnect clear approval; polling retains edits. Credentials stay server
side. Review digest binds zone/account, both normalized addresses and all sorted
provider records/metadata. Install validates early and Provision rechecks before
mutating each name. Automatic addresses use distinct first-level zone hosts for
standard TLS. Existing advanced certificate arrangements retain nested support.

Override is scoped to exact A/AAAA/CNAME names; non-web DNS blocks takeover. No
existing tunnels, ingress routes or unrelated DNS are removed. Save full originals
in protected cloudflare-dns-before.json before writes; on provisioning failure
remove only new-tunnel CNAMEs and restore deleted record values/metadata. Provider
IDs can change on restoration. Concurrent unknown changes are preserved and an
incomplete-recovery error points to the snapshot. CF has no conditional cross-name
DNS transaction; crash recovery needs provider inspection and the saved snapshot.

Go suite/setup/distribution/branding and vet pass, plus browser fixture editing,
polling, consent reset, summary/submission and mobile overflow checks. Fake CF
tests cover refusal, absent/stale review, multi-record replacement, mail-record
protection, pagination, mid-setup changes, rollback, lost POST response and failed
recovery. No live DNS or owner data used for development tests.

Version59 build outputs: I:/Projects/tooling/workspace-domains-20261021. Named NTFS
Go/cache area C:/Projects/.cache/workspace-domains-20261021. Build from committed
source with regenerated blue Windows resources; preserve all25 app payloads from
verified public58. Native wizard, publication and local installer copy receipts
are to be recorded after building. No owner DNS takeover is authorized by merely
testing this feature. Keep existing owner installation and newest recovery copy.
