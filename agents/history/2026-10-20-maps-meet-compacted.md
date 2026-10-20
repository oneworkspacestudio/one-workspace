# Maps and Meet integration — in progress

User authorized two native standalone apps, shared design/MyAccounts, Drive Office creation, installation, org repos and releases. Free/open source foundations only; no paid services or GitHub Actions. New app roots: I:/Projects/my-maps and I:/Projects/my-meet. Suite branch codex/maps-meet is isolated from concurrent Clips/navigation work in the primary checkout; reconcile the newer primary commits before release.

Shared apphost now has explicit opt-in geolocation, safe HTTPS image origins and same-origin websocket CSP. Added Service for embedded engine children under the existing Windows kill-on-close job boundary with an explicit scrubbed environment. My Maps uses Leaflet, user-triggered cached/rate-limited Photon and FOSSGIS routing, OSM attribution and provider configuration; no tile prefetch. Meet embeds pinned LiveKit 1.13.10 Windows x64 with SHA256 validation, separate signal route, per-room host admission and server-side permission checks before websocket proxying. All private state uses host issuer/subject identities and NTFS ACLs. Both AUI/TS frontends and Go development EXEs compile. Runtime and external-media checks remain outstanding; do not claim parity or a release yet.

Maps provider policies: OSM tiles forbid offline bulk downloads; Photon is best effort; FOSSGIS maximum one request/second. Remote LiveKit requires reachable media UDP/TCP (HTTP Cloudflare tunnel alone cannot carry media). No live traffic, Street View, proprietary reviews, Meet telephony, cloud transcription or unlimited conferencing claimed.

Build tools/caches: I:/Projects/tooling/maps-meet-20261020 and existing four-apps/env.ps1. GOTMPDIR is now task-owned on I: for large builds. No owner install changed yet.

## Release preparation

Merged the primary checkout's 7104646 (Workspace52 navigation and Clips3) before
integration. Published apphost dependency 0578d36ce297 is pinned in both new Go
modules. Suite54 adds maps/meet ports, routes, identity/session lifecycle, private
storage, icons, source pins, licences, app-update groups and contracts. Default
catalogue has 24 apps. Meet's derived TCP/UDP ports are included in install
preflight, with a clear lower-base-port error. README overview/table includes all
apps and retains existing hero artwork. New-app exclusion/update selection uses
the existing default-app contract.

Maps release source c25322c; Meet e92087e; Drive 1319b08b. Windows standalone
packages built locally, no Actions. Drive3 only promotes existing Office creation
to the main New menu; Photos bytes and pinned editor/Studio vendor assets remain
unchanged, verified against the public Drive2 release. Its untracked preview files
are preserved. Source/licence/checksum assets are staged for org publication.

Chrome synthetic preview: Maps OSM tiles and Photon results render on desktop
and mobile. Meet has real LiveKit connections, two synthetic media publishers,
host admission and cross-client chat. Physical media quality, public NAT paths,
capacity and full Google parity remain unverified. The 120 engine module licences
and missing npm protobuf licence texts are included. Maps stale reverse-lookup
selection, exact pin coordinates, Settings overlay restoration and mobile tab
shrink were corrected during visual review.

Next: finish public app releases, run the complete pinned org build, publish
Workspace54 and app update manifests, then apply the normal backed-up owner
update. Recheck current canonical Git and installed state first because the
Workspace52 task may still be publishing/updating. Remote Meet forwarding choice
is pending; never infer approval. Clean only task-owned synthetic preview
processes/data on C:, not the owner's rollback or other active work.

## Final Workspace54 candidate

Merged canonical 1bd7da5 to preserve Clips4 restart-safe migrations. The owner's
normal update to53 has now completed independently; recheck its live state
before applying54. Final pins are Maps2/Meet3 with minimum manager54. Both have
an integrated compact app picker; Meet binds its signalling grants to the
authenticated app session and restores video after camera re-enable. Drive3
main-menu creation produced real DOCX/XLSX/PPTX/PDF files, all reporting saved
in an isolated synthetic folder. Public release uploads/build and owner update
receipts remain to be finalized. Do not claim public-network media acceptance.

## Installed54; final release superseded by55

Normal backed-up owner update53-to54 completed; 24 apps and23 running services
healthy, other20 app versions/packages and all identity/storage fields preserved.
Maps2 and Meet3 are public with managed-update manifests. Drive3 remains uploading.
Public /maps/ and /meet/ return the expected unauthenticated login redirect.
Chrome control could not load its request-header policy; managed browser SSO was
not rechecked. Two synthetic Meet3 clients did pass media, lobby, chat and camera
re-enable checks before installation. Synthetic C fixture was removed (56MB);
normal updater cleanup freed about4.2GB and retained its verified newest rollback.

Concurrent53 publication exposed an inherited release-ordering hazard: choose()
used publication time, allowing an older version published late to displace a
newer version.55 filters valid published version tags, sorts numerically and
refuses manager downgrades regardless of publication time. No new tests added or
run for this follow-up. Build/package integrity and installation checks remain.
54 draft418837047 is held; its source tag stays immutable. Owner54 is healthy.
Final55 packaging/publication and normal owner manager update are pending.

## Installed 55; large release transfers pending

The normal 54-to-55 update completed at 2026-10-20T21:01Z. All 24 app versions and
packages, identity and storage locations were retained; all 23 service entries
were healthy. The newest rollback manifest was verified. Workspace.exe SHA256
is `52d7e5feab9c16c46d6fcc1de0ecc0498c37754b3a8007d353d24d7ff575e049`.
Do not rerun the 54/55 update scripts against this newer installation.

Final 55 source/tag is `de39246412774ee035d1b5a6cdb298f2930378d3`; main also merges
the independent Clips 53 installation evidence. Final bundle contains 24 apps
and 119 verified payloads. Maps 2 and Meet 3 are public. Drive 3 and Workspace 55
remain draft while large uploads retry HTTP 408 failures. The failed concurrent
Workspace 53 publisher has stopped; its draft is not an installed-version target.
The finishing publisher now uploads sequentially and verifies every remote
asset size/SHA256 before publication. It must not wait indefinitely for failed 53.
No Actions or paid providers were enabled. Remaining work: finish transfers,
verify public hashes against installed packages, and record publication receipts.
