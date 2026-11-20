# Vuno

Vuno is the new name for the song creation app previously called My Music.
The managed route, saved library and settings carry forward when you update.

Vuno creates songs from a brief or custom lyrics, plays and imports tracks,
and arranges clips into a stereo WAV. It is an optional app in Setup under
**Creative & Development**. On an existing Workspace, select **Vuno** in
Updates when applying the compatible Workspace release. Open the **Vuno**
tile and continue with MyAccounts. The Workspace owner can use the app.

## Connect a generation engine

Run a local [ACE-Step 1.5 API server](https://github.com/ace-step/ACE-Step-1.5/blob/v0.1.8/docs/en/INSTALL.md)
on the Workspace host. Vuno includes the client. The GPU runtime and
weights are a separate installation. Use a loopback HTTP listener such as
http://127.0.0.1:8001. Open **Engine** to check the connection, change the
loopback origin or provide the engine's API key.

To enable **Draft lyrics** and **Refine brief**, save a DeepSeek API key in
**Engine**. These actions send the music brief to DeepSeek. ACE-Step generates
the audio locally. Keys are stored in the app's protected settings.

Use **Library** to import audio from an earlier standalone installation. The
standalone files are under **%APPDATA%/One Workspace/My Music/Audio**. Import
creates a private copy and leaves the original in place. You can download
tracks or export a mix to your Workspace Files folder.

Use **Studio** to place tracks, set levels, preview the mix, and export a stereo
WAV. Export before closing the page because arrangement edits are held in the
current session. Workspace backups include the managed audio library and engine
settings. Restoring protected API keys to a different Windows account or PC
may require entering them again.

## Standalone Windows launch

Double-click **Vuno.exe** to open a separate loopback-local app window. Its
library stays separate from the managed Workspace instance. Keep the companion
console open while you work. Set **MYMUSIC_ACE_URL**, **ACESTEP_API_KEY**, or
**DEEPSEEK_API_KEY** before launch to configure the standalone instance.
