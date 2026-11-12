# Podcast co-host and intro track release

My Podcasts preview.7 (2257f7d) keeps production notes out of spoken dialogue,
inserts separate first-person introductions for Alex and Jamie, and labels
AI narration on episode pages and transcript exports. Account users can
upload two MP3/WAV intro tracks. The worker pins one randomly at enqueue,
places it after the cold open, and derives later chapter offsets from
samples. The upload and selection tables are an additive app migration.

My Voice preview.10 (5c4aeab) adds the pinned male Marius and female Alba
Pocket presets and advertises their cast through the existing loopback
grant. Older installs can add the second preset via owner-approved Pocket
reinstallation. Existing credentials, routes, namespaces and the fallback
single-preset engines remain compatible. No manager or Android build.

Go/vet, eight Podcasts browser checks, native Pocket and MOSS PCM, three
signed-out scheduled WAVs, and a native intro placement/remastering check
passed with synthetic account/research inputs. The protected installer
verified both updates, all 164 payloads, managed pairing, track upload,
and rollback on NTFS. Supplied tracks remain private and are omitted
from public releases. Owner installation and private upload are pending.

Channel commit 2cf7689 promotes both immutable releases and pins. The normal
anonymous Workspace preview.123 updater offered and applied Podcasts 6→7
and Voice 9→10, verified all 164 files, and preserved DeepSeek settings and
cadence. A signed-in browser uploaded two synthetic MP3s through the managed
route. The protected rollback restored both previous versions and kept
manager and Accounts binaries unchanged.
