# Mail19 incoming recovery and settings persistence

Continuation of Mail setup repair, publication/install previously authorized.
Mail19 source f3967576751c75723778d94bbdbe4577459fbb33 is public; every uploaded
asset size/SHA256 verified. Public release: workspace-mail-v0.1.0-preview.19.
Promoted into updates/releases-v1.json and organization lock, retaining all other
repository entries and Workspace75. No Actions or automated tests run in this
follow-up. JS syntax and Windows compilation passed. Mail19 adds generation-safe
status polling, preserves open panels, connection-only retry, native gateway
lookup/targeted SSDP, structured blocker/next-check help, and verified missing
outer-rule recovery. Compatible unowned routes stay read-only.

Affected real installation is workspace-2, now Workspace75 and Mail18. It has
completed post-domain DNS setup; latest actual status says SMTP25 belongs to a
different PC. Asked user whether original receiver should be shared or replaced;
decision pending. Do not silently redirect another mail service. One publicSMTP25
cannot serve two independent internal hosts just by adding a new email subdomain.
This PC C:/tmp installation is stopped; do not start it as an unrelated fixture.
Remote owner updater75 is available and can install Mail19 using normal snapshots.
Deployment and UI observation pending. Artifacts I:/Projects/tooling/mail-network-20261023.

## Mail20 correction and shared-receiver decision
Mail19 managed-update check was correctly refused because it lacked
workspace-update.json. No installation was attempted. Mail20 adds schema1 /
integration2 metadata for all three owned packages, published and size/SHA256
verified. Immutable19 remains intact. Updated the component lock and refreshed
the public channel using scripts/publish-update-index.py. Windows compilation
and JavaScript syntax checks passed; no automated tests run.

The owner confirmed the other PC should remain the main incoming receiver.
Preserve that SMTP route. Multiple PCs should access the main Mail service and
its identities/mailboxes; independent stores require a gateway and routing.
Central URL/connection mode is being clarified. Do not equate identities across
separate MyAccounts issuers. Mail20 remote owner deployment remains pending.
