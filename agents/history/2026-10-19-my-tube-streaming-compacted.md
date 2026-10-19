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

Workspace41 code/source commit17a37ac1f98c514593805e390b9b44b5df5435cc is pushed
to main and codex/four-apps. Local native19-app bundle has91 verified package
hashes; My Tube packages match the public component manifest and other app
payloads/pins are unchanged. Refresh both the release-status notice and embedded
organization-components.json when incrementally assembling verified payloads.

Isolated native baseline40 Accounts/Drive/MyTube1 installation opened My Tube
through one Workspace sign-in and seeded a synthetic original. The released40
controller applied the exact Tube2 package; actual41 installer then updated the
manager. Release binary hashes, identity/locations and healthy startup matched.
Native cleanup removed one older experimental rollback (~853MB), retaining the
latest. Managed playback after the final update reused verified helpers and
four qualities, preserved the original checksum/ID and decoded360p at simulated
1.2Mbps. No personal runtime used as a functional fixture.

Stopped task-native processes/browsers, checked process/path/reparse boundaries
and removed1,616,551,814 bytes from the exact named C task root. About768MB of
interrupted task-only build temporaries were also removed earlier. Owner health
still200 and installation/credentials/tunnel/data untouched. Evidence/binaries
remain under I:/Projects/tooling/tube-streaming.

## Published and complete

Published [Workspace v0.3.0-preview.41](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.41)
from17a37ac1f98c514593805e390b9b44b5df5435cc, with24 assets. Every remote asset
size/SHA256 matches local staged files; public catalog downloaded/hashed and
anonymous installer/bundle downloads checked. Tag source matches the build
provenance. Both repositories' main branches contain the implementation and
Actions remain disabled. My Tube component has ten verified public assets.

No pending release work. Publication does not deploy to the owner's running
installation: use Workspace Updates for Workspace/My Tube, then keep the suite
running while legacy originals prepare. Existing owner settings omit processing
flags, so new defaults apply on update. Preserve explicit opt-outs, owner data,
latest owner rollback and three unrelated untracked earlier histories. Do not
retag or replace published release assets.
