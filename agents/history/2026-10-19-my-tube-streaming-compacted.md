# My Tube adaptive streaming — 2026-10-19

## Purpose and scope

User requested source optimisation, push and release after buffering diagnosis.
Work in I:/Projects; protected disposable native state only in named
C:/Projects/.cache/workspace-tube-streaming-20261019. Preserve owner installation.
No GitHub Actions. Build locally; both repositories have Actions disabled.

## Implementation

My Tube v0.1.0-preview.2 source da3dc1d automatically prepares adaptive HLS on
Windows x64. Direct upstream pinned FFmpeg archive/executable hash and size
verification precedes execution. Helpers are downloaded into private app data,
not redistributed inside the suite. Background status/retry preserves startup
and original playback. Existing explicit opt-outs/custom helper paths survive.
New and legacy original-only videos queue for bounded 360/480/720/1080p copies.
Measured actual segment bandwidth, aligned keyframes, bounded encoding threads
and conservative hls.js startup reduce slow-network buffering. Originals remain
unchanged; visibility checks also protect playlists/segments.

Workspace41 enables these defaults for new managed configurations, pins the
published component ZIP and adds docs/MY-TUBE-STREAMING.md. Existing settings
remain intact. Other component pins/payloads must match public Workspace40.

## Verification completed

My Tube source Go tests/vet with real reviewed FFmpeg9.0.2, AUI/TS/Vite build,
seven native legacy-upgrade/helper-download/streaming checks and three final
packaged EXE checks passed. A synthetic21.38Mbps source produced average
0.59/0.95/1.70/3.23Mbps copies; actual Chrome playback decoded360p at simulated
1.2Mbps. Exact original hash/range bytes, private-media anonymous denial and
restart reuse passed. Managed settings regression/vet passed. No live owner data
used as functional fixtures; physical clean-host, reboot/SCM, all-codec and
independent audit limitations remain.

## Release preparation / resume point

My Tube2 is public from da3dc1d8f59d9932c09cac71b3528c484c9cdf6f. All ten
asset sizes/SHA256 matched; public ZIP downloaded/hashed. ZIP SHA256 is
bf788ebef31d770224eb7d548d7c771e177c537707330ebcf0ef8d4521c6a5ac.
GitHub REST release listing omitted the newly created draft although GraphQL
listed it; fetching draft by databaseId returned full authenticated metadata.
Verify draft by ID and never create duplicates based only on REST listing.

Next: commit/push Workspace41 source; assemble verified40 bundle with only
manager/MyTube/notices changed, build exact native installer, verify91 package
hashes and component manifest, run isolated selected-app install/update gates,
publish after remote digest verification. External tooling/evidence is under
I:/Projects/tooling/tube-streaming. Update this file with actual release/native
acceptance and cleanup results. Preserve unrelated untracked earlier histories.
