# My Clips

Workspace 0.3.0-preview.51 adds My Clips 0.1.0-preview.2 as the22nd ecosystem app.
Source and standalone Windows/PWA/Android packages:
https://github.com/oneworkspacestudio/my-clips

Features include vertical For You/Following/Friends feeds, camera capture,
multi-clip editing, filters/glow, timed text, CC0 soundtrack defaults, authorized
My Music import, voiceover, scheduling, private drafts, audience controls,
Duet/Stitch, profiles/follows, likes/comments, mutual-follow messaging,
notifications and small WebRTC live rooms.

The original recommendation system learns from watch completion, skips,
hashtags and follows. It does not reproduce TikTok's proprietary model.
No commercial music catalog, face AR, automatic captions, Shop/coins/gifts,
production-scale live service or full TikTok parity is claimed.

## Integration

Default loopback port9014; managed route `/clips/`; MyAccounts OIDC and the
shared app picker. Installer selection and normal Workspace updates can add the
app; existing opt-outs, identities, paths and unrelated app versions persist.
Standalone app updates require `scoped-clips-storage-v1` on Workspace51+.

Originals live in `Data/videos-clips/My Photos/My Clips/users/<identity-hash>`.
The native **My Clips Originals** source is outside broad Drive/Files roots,
disabled by default and read-only in each account's own physical scope.
Its loopback authenticated guard verifies containment and actual permissions.
The app fails closed if that source is unsafe or unavailable. Derivatives remain
outside native original scopes and every playback request rechecks audience.

Videos are transcoded into 360p/540p/720p HLS with two-second segments, measured
bandwidth, bounded buffers and MP4 byte-range fallback. A single render worker
uses bounded FFmpeg processes. Operator tools can be selected with
`MYCLIPS_FFMPEG` / `MYCLIPS_FFPROBE`, or the app obtains its pinned verified
upstream helper. FFmpeg is not redistributed in the suite payload.

Native Android package `studio.oneworkspace.clips` connects to the Workspace
HTTPS origin, handles camera/microphone and system file selection, and optionally
checks activity approximately every15minutes. Browser Web Push requires HTTPS
and supported permission/install conditions. Physical-device and external TURN
acceptance remain separate gates.

## Validation

Local Go suite tests/vet and default-app update regression include Clips.
Synthetic native sign-in and Drive tests check byte identity plus other-account,
traversal, case-variant and broad-source denial. Browser tests exercise rendering,
privacy, HLS/ranges, throttled playback, drafts, capture, social actions, live video
and a390px layout. See the app's public acceptance record and this release's
machine-readable receipts. Never use the owner's installation as a fixture.

Workspace53 and Clips4 place Workspace inside the mobile navigation and desktop
sidebar. Shared navigation enhances data-workspace-apps-trigger links into the
picker, retaining their href fallback. The floating launcher is hidden only while
an integrated trigger is visible; resizing and SPA changes restore it as needed.

Clips4 also fixes reopening initialized databases: only pending migrations run.
Clips2/3 rejected their own completed schema on restart. Workspace52 was held
before publication;53 includes the correction. Existing data is preserved.
