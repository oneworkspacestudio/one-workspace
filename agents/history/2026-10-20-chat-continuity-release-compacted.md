# Chat continuity integration

Current user authorizes Chrome offline verification, pushes and releases without
GitHub Actions. Preserve Workspace48/Social2 and all other components. This
increment targets Chat10 and paired Workspace49; the owner installation is not
a fixture or deployment target in this turn.

Chat source9ec055bf3d742a5044660cad4ca88b45c6a92190 is clean packaged10:
EXEb669639b71301a430245c2d2f10a22f50cad501ae82c81e01f01997baabdebf8;
ZIP2d7b345aacc3d78dec45bf761c62f92b19d027394969108e99b9ec6e438c2c3e.
13 exact-binary real Chrome continuity checks pass: Chrome network disabled,
cached same-URL reload, durable queue across another reload, network restoration,
public discovery reconnection after both clients reload with the whole server
stopped, real text/emoji/sticker delivery, duplicate and revocation checks.
Seven messaging/security/UI/mobile/navigation/retry/QR regression suites pass;
18 controller,10 signaling,5 WebCrypto and13 older controller cases pass.
One tooling batch used .mjs for three older .cjs scripts; those actual .cjs
scripts then ran separately and passed. Native gateway/server/store/rendezvous
tests and complete Workspace tests pass. Four synthetic example.test accounts.
No owner data, mocked RTC success, separate-network or physical-phone acceptance.

Pin Chat10, update build defaults and documentation. Build the complete21-app
bundle by replacing only manager, Chat and current notices in verified48.
Keep Social2 native scoped storage; do not revive preview47 shared-folder ACLs.
Scripts/evidence: I:/Projects/tooling/chat-parity-20261019, final names chat10,
release-chat10 and release-workspace49. Protected C fixture and browser mirror
are disposable after proof capture; preserve owner files/newest backup.

Public Chat tags are immutable. Verify uploaded asset digests before publication,
then verify Workspace49 installer/update integration. The Workspace source
commit uses [skip ci] to honor the current no-Actions instruction; no workflow
dispatch or release-trigger workflow is used.


## Publication complete

Chat10 and Workspace49 are public, with exact uploaded asset sizes/SHA256 and
tag/source mappings verified (15 Chat assets,21 Workspace assets). Chat's
anonymous public ZIP matches the locally tested digest. Workspace49 bundles
21 apps/104 payloads; only manager, Chat and the two release/provenance notices
changed from48. All payload hashes and launch contracts pass. No Actions used;
the Workspace release source commit has zero workflow runs.

Exact-binary Chrome acceptance passed13 outage cases plus seven regression
suites and46 controller/signaling/crypto cases. Fresh49/Chat10 installation
passed MyAccounts, discovery, messaging and permissions. One first-search
15-second timeout occurred during new directory synchronization; after the
sync window the same unmodified browser checks passed. Published48/Chat9 to
49/Chat10 via the normal updater preserved identity/issuer/paths, native Chat
IDs/history/exact attachment bytes and19 explicit app opt-outs. No owner fixture.

Public source commits: Chat9ec055bf3d742a5044660cad4ca88b45c6a92190;
Workspaceb3282bee45343489cd5bcdaec08504d09cbfe498. Tags/assets stay immutable.
Follow-up evidence/history commits change no packaged runtime code.
Owner installation was not updated this turn; its local health remains OK.
Three synthetic suites stopped and all their listeners closed. Named C data,
browser mirror and test temp roots removed after evidence capture, reclaiming
1.79GiB. Unknown data and owner newest recovery copy untouched.
Builds/downloads/evidence remain in I tooling/chat-parity-20261019.

Separate-network/physical Android/iOS/TURN/background delivery remain untested.
Outage grants expire within two hours; normal server messages are not E2EE.
Future release builds must recreate the removed NTFS client mirror.
