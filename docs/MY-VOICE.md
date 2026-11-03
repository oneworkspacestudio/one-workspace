# My Voice

My Voice lives in Creative at `/voice`. Select it during setup, or update to
Workspace111 and explicitly select the new optional My Voice row in Updates.
Accounts SSO and the shared Apps picker use the normal Workspace session.
Existing apps, paths and independent app versions remain selected as before.

Speech and private voice cloning run on the Windows host CPU. In Settings,
review and download a model once: Pocket18/7 provides lightweight English preset
speech; MOSS provides cloning and its explicitly listed languages. Downloads use
fixed public sources and verified hashes. No Hugging Face account is required.
Weights are not bundled. The app shows download failures and retry state.

Create a voice only from a recording you are authorised to use and confirm consent.
References, prepared states, projects, history and generated audio stay in private
per-account storage. Shared voices require an explicit grant; revoking it blocks
future use. Delete references and prepared states from My Voices. WAV, MP3 and
batch ZIP exports are available; scoped Drive import/export checks your access.

Installation, runtime, recovery, API and native CPU measurements are documented
in the [My Voice repository](https://github.com/oneworkspacestudio/my-voice) and
[preview2 release](https://github.com/oneworkspacestudio/my-voice/releases/tag/v0.1.0-preview.2).
MOSS was slower than real time on the tested CPU; older low-memory PCs and
non-English pronunciation need separate review. This preview does not claim full
ElevenLabs feature or voice-quality parity. Model downloads remain an explicit
first-use action on each PC.
