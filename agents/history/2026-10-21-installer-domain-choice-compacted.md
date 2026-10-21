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

## Verified delivery

Workspace59 is public from c9397f8. Every published asset size and GitHub SHA256
digest matches the local package. All124 bundle payloads verified; all25 app
contracts and application bytes retained from public58. Seven shared blue PNG
icon frames are embedded in both new Windows executables.

The actual full Windows installer served the new UI in desktop/mobile Chrome.
Synthetic provider responses verified custom names, edits surviving polling,
explicit conflict consent, consent reset after editing, correct summary and
submitted chosen addresses. Native browser had no script errors or mobile
horizontal overflow. Install was intercepted: no live DNS, owner data or owner
installation was changed as a fixture. Native installer and browser stopped.
Verified installer copied to Downloads/WorkspaceSetup-0.3.0-preview.59-windows-amd64.exe.

Privacy-filtered delivery receipt: docs/evidence/installer-domain-delivery-20261021.json.
Task C cache removed after checking containment, no reparse points, no active
process references and empty contents; durable source/build outputs stay on I.
No real provider takeover or clean-host/reboot acceptance claim. Existing owner
Workspace58 and its verified recovery copy are preserved. Release tags/assets
are immutable; future fixes need a new version.
