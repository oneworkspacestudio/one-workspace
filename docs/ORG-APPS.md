# Organization apps

`scripts/rebuild.ps1` calls `scripts/build-org.ps1`, which packages the pinned
organization applications. Every app's catalog source must
belong to `https://github.com/oneworkspacestudio/`. Download URLs, byte sizes and
GitHub SHA256 digests are verified, including cached downloads. The package
contains `licenses/components/organization-components.json` with exact provenance.

The October 15 bundle selects the updated individual app previews and the
previous published Music version. The new Music executable is withheld pending
Microsoft analysis; the build script rejects that exact binary. Current feature
and acceptance limits are documented in [release status](RELEASE-STATUS.md).

The catalog declares installer integration2. Older setup executables reject this
field before installing and tell the user to download the current installer.
They cannot safely register shared-service clients or reassign Calendar ports.
Always download the current full installer when updating from setup17 or earlier.

The preview 34 bundle pins these organization components:

| App | Organization repository | Packaged version |
| --- | --- | --- |
| MyAccounts | [myaccounts](https://github.com/oneworkspacestudio/myaccounts) | preview-20261015.1 |
| Drive / Office | [drive](https://github.com/oneworkspacestudio/drive) | drive-v0.5.0-preview.2 |
| Photos | [drive](https://github.com/oneworkspacestudio/drive) | drive-v0.5.0-preview.2 |
| Photo Editor | [Photo-Editor](https://github.com/oneworkspacestudio/Photo-Editor) | studio-v0.1.1 |
| Git | [gitea](https://github.com/oneworkspacestudio/gitea) | workspace-20261015 |
| Music | [Music](https://github.com/oneworkspacestudio/Music) | workspace-20261014 |
| Mail | [mail](https://github.com/oneworkspacestudio/mail) | 0.1.0-preview.16; independently resolved online |
| Calendar / Contacts | [calendar](https://github.com/oneworkspacestudio/calendar) | 0.1.0-preview.2 |
| Notes | [notes](https://github.com/oneworkspacestudio/notes) | 0.1.0-preview.4 |
| My Tube | [my-tube](https://github.com/oneworkspacestudio/my-tube) | 0.1.0-preview.1 |
| My Vault | [my-vault](https://github.com/oneworkspacestudio/my-vault) | 0.2.0-rc.3 |
| Chat | [chat](https://github.com/oneworkspacestudio/chat) | 0.1.0-preview.3 |

Drive/Photos use their own current pinned organization release, retaining the
upload fixes and new artwork. Preview 34 adds Chat; all existing component executable pins are retained.
Office and the integrated Photo Editor run inside Drive. Studio.exe is installed
as well for standalone local editing. The integrated tile uses Drive's existing
MyAccounts authentication and file permission checks.

The one installer contains the released applications and requires no Docker.
Cloudflared and the pinned native Stalwart engine remain third-party dependencies;
the organization app wrappers and their source are not represented as owning
those upstream projects. Component licenses/notices remain in the package.

## Existing installations

After stopping Workspace, use the new manager from outside its installed folder:

```powershell
& ./Workspace.exe --data 'C:/path/to/workspace/Private' install-apps --bundle ./workspace-bundle.zip
```

This adds missing apps, registers their MyAccounts clients, updates the manager
and organization provenance, and backs up changed program/configuration files.
It refuses replacing an existing process executable whose bytes differ; that
requires a separate explicit upgrade. Existing app data, account databases and
ordinary files are retained. Restart with the protected installed manager.
Private Mail never silently switches to the shared Calendar/Contacts store.

## Notes

Notes became available during this integration. The bundle uses preview4, which
retains the Workspace callback fixes and includes the editor crash hotfix. Its actual
CLI, OIDC environment and `/notes` proxy contract were checked against source.
Workspace generates a separate protected `notes-settings.json`, gives the
Workspace owner Notes administration, and keeps the SQLite data in `Private/notes`.
Existing settings are preserved. Notes remains in managed MyAccounts mode; no
local-login fallback is enabled. Full Joplin sync/E2EE is outside this app's
current preview; see the Notes repository's documented feature boundaries.

## My Tube

MyTube.exe is a verified component of the full offline/online installer payload,
not an optional post-install download. It appears in the home and shared app
switcher at `/my-tube`, with the supplied My Tube artwork. Its backend uses the
selected Workspace port plus13, depends on Accounts/Drive startup, and receives
a dedicated confidential OIDC client through the protected secret store.

The manager creates `Private/my-tube-settings.json` outside the dedicated
`Private/my-tube` data/media directory. The Workspace owner is an explicit My Tube
administrator; existing settings are preserved. Mode stays managed, so no local
owner link or duplicate setup wizard is enabled. Public access/uploads default
off. Drive requests use its direct loopback `/drive` mount with the calling
user's Drive cookie/token, never a shared administrator credential or a shortcut
around the portal's browser session. Workspace logout clears My Tube's cookies.

The package includes My Tube's AGPL license, notice and third-party notices.
FFmpeg/FFprobe are optional separate helpers and are not bundled. Compatible MP4
upload/playback works without them. See the app's own release for the tested
features and remaining preview acceptance gates.

## Music runtime storage

Music's packaged Python runtime extracts into an installation-owned private
`music-runtime-temp` folder. Workspace reclaims only its stale `_MEI` extraction
directories before the next Music start, so forced supervisor shutdowns do not
accumulate copies in the Windows user's global temporary folder. User music,
configuration and unrelated temporary files are preserved.

## Chat

Chat.exe embeds Tinode, its locally bundled web client and a private PostgreSQL
runtime. No Docker or separate database installation is needed. Workspace passes
Private/chat, base port+16, /chat and a dedicated OIDC client/callback. The normal
installer app selection can omit Chat. Ordinary updates do not enable an app the
owner omitted; use installer app selection to add Chat to an existing suite.

The app uses MyAccounts identity tokens and one-minute single-use engine tickets;
issuer/subject determines the stable native identity. Its own HttpOnly session
protects messaging and attachments. The shared launcher starts normal OIDC when
needed, and sign-out clears Chat cookies. Its messages and ordinary attachment
files live in its own protected folder and are included in stopped suite backups.
Chat's GPL/Apache and dependency notices accompany its exact org ZIP pin. See
[Chat setup and limits](CHAT.md); initial runtime acceptance remains pending.

## My Social / preview47

The twenty-first app is [My Social](MY-SOCIAL.md), from the organization-owned
`my-social` source and verified preview1 package. Existing component pins remain
unchanged. The native guarded My Photos adapter and explicit SSO/logout/update
contract require Workspace47 or later. Feature gaps remain explicit.
