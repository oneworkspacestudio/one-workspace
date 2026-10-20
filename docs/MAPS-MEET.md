# My Maps and My Meet

Workspace53 adds both apps to the default installation and update catalogue.
Use the Apps page to open them; existing MyAccounts sessions start their normal
authenticated SSO flow automatically. Neither app uses a second managed password.
Advanced installation can exclude them. Standalone EXEs provide their own local
owner setup when no Workspace identity provider is configured.

## My Maps

`Maps.exe` uses Leaflet and OpenStreetMap. It includes search, dropped pins,
private saved lists and notes, share links, measurement, location following,
driving/cycling/walking routes with alternatives and up to eight stops, and
GeoJSON/GPX downloads. State is under the managed `maps` private data folder.

The default community tile, Photon and FOSSGIS services are free and have no
availability guarantee. Requests disclose the requested search/route to those
providers. There is no bulk tile download. Operators can configure compatible
HTTPS providers or host open source services themselves. No traffic, Street View,
satellite, transit timetable, review corpus, offline package or voice navigation
is included. [App documentation](https://github.com/oneworkspacestudio/my-maps).

## My Meet

`Meet.exe` embeds the pinned, checksum-verified LiveKit Windows binary. It extracts
into protected `meet/engine` data and is supervised as a child process; one EXE
launch starts both the app and engine. Meetings, participant permissions and chat
stay in the managed `meet` data folder and participate in normal backups.

Includes camera/mic preview and selection, meetings and dated links, an account
and host admission lobby, video/audio, screen sharing, pinning, host removal/mute,
room locking, end for all, chat, hands, reactions and browser screen recording.
Recordings download as ordinary WebM files; they can be uploaded to Drive.

The web UI and authenticated signalling use `/meet/` on the usual Workspace URL.
Media needs a reachable UDP or TCP route to the PC. The default UDP port is
10017 and TCP is 10018; alternate Workspace base ports change these. Meet Settings
shows the actual ports. Configure router forwarding/firewall rules and
`MYMEET_EXTERNAL_IP=true` before claiming external calling. Restrictive networks
may need a separately hosted TURN service. No paid service is enabled.

There is no anonymous guest entry, telephony, breakout rooms, polls, cloud
recording, captions/transcription or background effects. The 32-participant room
limit is a configured ceiling, not a capacity certification. Media is encrypted
in transit but handled by the self-hosted server, not end-to-end encrypted.
[App documentation](https://github.com/oneworkspacestudio/my-meet).

## Drive Office files

Drive → **New** now offers **Word document**, **Excel spreadsheet**,
**PowerPoint presentation** and **PDF document** directly. The existing creation
and editor paths preserve Drive permissions and store ordinary files in the
selected source folder, with the existing autosave/history features.

## Evidence

Local TypeScript/AUI/Vite and Windows builds completed. Isolated Chrome previews
used synthetic identities and media: Meet host admission, two clients receiving
media and shared meeting chat worked. Maps displayed OSM tiles and Photon results.
These checks do not establish physical-device quality, public NAT traversal,
large meeting capacity or full feature parity. Release and installation receipts
are separate from this scoped preview evidence.
