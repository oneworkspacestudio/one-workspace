# My Music in One Workspace

My Music preview.2 accompanies Workspace preview.142. The legacy Music player
keeps its separate id and route. My Music is optional in Creative & Development
with id mymusic, route /mymusic, catalog port 9049 (base + 59) and a MyAccounts
callback at /mymusic/auth/callback. The shared apphost protects library audio,
generation, backups and engine settings for the Workspace owner.

The same executable supports standalone launch with no arguments. Its prior
library remains in the user config folder. To move tracks into the managed app,
import copies in Library. Workspace stores its own audio and encrypted engine
settings under its app-specific private directory. The shared picker is
available from the header Apps link.

ACE-Step 1.5 uses an independently installed loopback API. Generation jobs are
polled so a long model run does not hold a reverse-proxy request open. DeepSeek
is opt-in through Engine settings. The release owns a single EXE, compatibility
contract, guide, licence and notices. Model weights and Python/CUDA are
dependencies of the external engine installation.

The ordinary rebuild uses immutable preview.141 bundle bytes and introduces
My Music with installer integration 23. The registry tracks the first
integration containing the app so previous bundle revisions can still be
validated. An older installer must hand off before applying integration 23.

Preview.1's standalone paths passed synthetic API, playback, import, search,
Studio preview/export and narrow layout acceptance. The managed release is
tested with synthetic accounts in NTFS roots under the named task cache.

## Verified release and installation

The ordinary preview.142 rebuild validated 47 apps and 939 packages. The
Windows setup and app EXE used for acceptance match the public release.
Catalog, distribution, suite and apphost tests passed with NTFS test roots.
Publisher tests and update-index tests passed. Synthetic setup covered My
Music selected and omitted, the real MyAccounts owner flow and shared picker,
CSRF and expired-session denial, audio imports/playback, Studio WAV export
and narrow/desktop layouts. A native app-only fixture update retained the
manager, account and audio. The preview.141 installer rejected integration
23. The current installer upgraded the same baseline with identity and files
preserved and restored it through supported rollback. The large offline
bundle needed NTFS staging for timely full verification on this host.

The permanently signed Android companion shipped with the 958 verified
release assets. The immutable suite tag points to source ad5dc8e. Channel
promotion initially raced an unrelated My Voice/Podcasts channel push.
Merged the published source over the newer main and regenerated the supported
index. Channel commit 3c5d2ad preserved those releases. An anonymous
workspace-package-mymusic.bin download matched the staged and installed
SHA-256 (834063354ecb9419af8e99d2bb031f4acb674c94288e05abd7d46ef93f4df190).
The ordinary installer check sees preview.142 and My Music preview.2 up to date.

The authorized installation on the current Windows PC completed through the
normal backed-up update, with the same payloads and served JavaScript. The
Workspace origin and native app health both passed. The existing Music player
hash and locations were preserved. The new recovery copy was verified. Older
backups that cleanup could not verify were retained by the installer.
No physical handset or reboot test was performed.

The ACE-Step v0.1.8 CUDA dependencies and model checkpoints were prepared
separately. PyTorch confirmed CUDA. The earlier synthetic generation QA used
an API fixture. The initial live engine readiness check timed out and was
retried with lazy initialization. DeepSeek needs the owner's optional API key.

The lazy ACE-Step retry reached HTTP readiness and accepted a generation
request, then the Python backend exited while loading the DiT checkpoint.
The client received a connection reset and no WAV was produced. No ACE-Step
process was left listening on the production port. Keep generation acceptance
open. Preserve the downloaded engine, pinned checkpoints and diagnostics for a
follow-up. Low free space on the private NTFS drive precluded staging a second
full Windows CUDA layout without jeopardizing the existing installation.

## Vuno rename and interface redesign

The user renamed the song generation app to Vuno to distinguish it from
the existing Music player and requested the reference charcoal palette.
Vuno preview.3 adds Home, a split Create view, searchable and sortable
Library, persistent playback controls, and a waveform arrangement Studio.
The layout, icon and wording are Vuno originals. The app ID mymusic,
managed executable path, routing, sessions, OIDC callback, package namespace,
standalone library and managed backup root retain their existing contracts.
The standalone public download is named Vuno.exe.

Workspace preview.143 seeds from the immutable preview.142 archive. A
historical inventory check accepts the old preview.142 display name. New
catalogs use Vuno. Synthetic browser checks passed at desktop and mobile
widths, including WAV import/export and mocked ACE-Step requests. The
prior live ACE-Step model loading issue is outside this UI revision.

Fresh managed setup and offline upgrade passed on isolated NTFS data. The
upgrade retained the imported audio and MyAccounts session, showed Vuno in
the app picker, and rollback restored preview.142 and its prior branding.

Published Workspace preview.143 from cd91faadbabed59cd271b8eea81c0ba58d7b6396
through the signed Android publisher. The publisher verified 958 assets and
promoted channel commit ef82afd. Anonymous checks confirmed the live index and
catalog list Vuno preview.3, and both Vuno.exe and the managed mymusic package
match SHA-256 c18ff0c6e802b38f58baaff22aa5aa34dfb2e297fe7cb36616ad72b1ac07bcb5.
The owner installation remained on preview.142 during the UI-only revision.


## Live Vuno installation

The owner requested installation and live verification for future managed app
updates. Use the normal backed-up Workspace updater on the active target and
verify the deployment. Keep functional test data isolated from owner data.

The public channel offered a selective Workspace and Vuno update. The current
Windows PC finished that update to Workspace preview.143 and Vuno preview.3,
with a verified recovery copy. All five Vuno payloads matched the public
catalog. The separate Music executable and existing install settings were
preserved. The running Workspace and Vuno native health checks passed. The
served HTML, CSS and JS matched the release source, OIDC was configured, and
the library rejected anonymous requests. Browser checks against the exact
installed binary passed import, playback, and stereo Studio WAV export using
isolated data.

The retry against the real local ACE-Step server reached inference submission.
During DiT model loading the Windows volume fell below 2 GB free, and the test
stopped the dedicated backend to preserve room for Workspace. No audio was
generated. The test backend was stopped and Vuno health passed again after
cleanup. Keep live song generation open until the host can load the model.


## ACE-Step runtime acceptance

After reclaiming space on the Windows volume, installed the prepared official
ACE-Step 1.5 v0.1.8 runtime at its pinned upstream commit dce6214. It serves the
Turbo model at the existing Vuno loopback origin 127.0.0.1:8001. The working
profile uses eager initialization, CUDA residency for the audio models and
runs with the optional 5 Hz language model disabled. CPU offload triggered
paging and the disk reserve during synthesis on this host. The 5 Hz model also
triggered the reserve during startup on the original configuration.

Browser tests used the exact installed Vuno binary with isolated NTFS data and
the live engine. Both the instrumental job and the song with supplied lyrics
created, saved and played non-silent 10-second stereo WAVs at 48 kHz. The
prompt-only vocal job requires the language model and failed as expected under
this profile. The engine stayed healthy after that check. The per-user Start
menu entry ACE-Step for Vuno launches on demand and was exercised while an
engine was already running, which preserved the active listener. It leaves the
existing Vuno installation and owner library intact.


## Vuno lyric timing

Vuno preview.4 estimates duration from sung line length, written section
boundaries and tempo. With no explicit tempo it uses 105 BPM for auto-length
tracks with supplied lyrics. Manual length remains available and warns when
it is shorter than the estimate. The ACE-Step request sends the supplied
lyrics unchanged, with the selected vocal language and explicit formatting,
caption enhancement and language enhancement disabled.

A synthetic six-section nursery rhyme established the baseline on the live
ACE-Step Turbo engine. A 30-second clip omitted most of the lyrics. A 45-second
clip using just the first two verses also skipped opening lines. A separate
180-second run at 105 BPM with English and formatting off included all six
sections in order according to local audio transcription. Wording and timing
can still vary across renders. Workspace preview.144 seeds from the immutable
preview.143 bundle and preserves Vuno's existing installation contract.
