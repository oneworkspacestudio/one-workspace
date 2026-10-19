# Smoother My Tube playback

Workspace preview.41 includes My Tube preview.2. On Windows x 64, video
optimisation and automatic media tools are enabled by default. Existing settings
that explicitly disable either feature are respected. Original videos are never
replaced by the smaller streaming copies.

## What happens after updating

1. My Tube prepares verified video tools in its private data folder. The first
   download is about 115MB, directly from the pinned upstream publisher. Settings
   shows actual downloaded bytes and percentage. Allow about 350MB temporary space.
2. Existing completed videos without the current streaming profile join the
   background queue. New uploads join automatically. Keep Workspace running while
   Creator Studio shows processing jobs. Existing originals remain playable.
3. Auto playback starts with the smallest quality and adapts to the connection
   and player size. When a new adaptive copy becomes available, playback position
   is preserved. Smaller sources are not upscaled.

| Quality | Target video rate | Peak video cap |
| --- | --- | --- |
| 360p | 450kbps | 650kbps |
| 480p | 800kbps | 1.1Mbps |
| 720p | 1.5Mbps | 2Mbps |
| 1080p | 3Mbps | 4Mbps |

Audio adds 96/128kbps and transport adds some overhead. The HLS playlist advertises
measured segment bandwidth, including multiplexed audio/transport overhead.
Four-second segments and aligned keyframes support quality switches. Buffers,
worker count, threads, storage and job time remain bounded. Slow connections below
the smallest rendition's rate can still buffer; this is not unlimited bandwidth.

## Failure and configuration

Interrupted or unavailable downloads retry automatically at one-minute intervals.
Originals stay safe and playback stays available. Processing failures are visible
in Creator Studio, where owners can retry. Explicit `transcode:false` or
`auto_helpers:false` disables the corresponding automatic work. Reviewed custom
FFmpeg/FFprobe paths remain supported in advanced administration.

Helpers are separate processes downloaded directly from Gyan's pinned FFmpeg 9.0.2
release, with archive and executable SHA 256/size checks before execution. They are
not bundled in the My Tube or Workspace release ZIP/EXE. Runtime and source links
are in My Tube's docs/media-tools.json; the main apps still launch without Docker,
Node.js or Python. A one-time Internet connection is needed for automatic tools.

## Verification and limits

Real native tests upgraded the previous released executable's synthetic 1080p
original, downloaded tools through the implemented installer, generated all four
qualities, compared original bytes/ranges, denied anonymous private-media access,
decoded Chrome HLS under a simulated 1.2Mbps network and reused caches after restart.
The 21.38Mbps source produced average rates about 0.59/0.95/1.70/3.23Mbps; Auto decoded
360p at the constrained rate. Go tests/vet, including real helpers, and the AUI/
TypeScript production build passed. No live owner data was used as a fixture.
These are synthetic checks, not certification of all codecs, an actual remote
tunnel, clean physical Windows host, reboot/service startup or independent audit.
