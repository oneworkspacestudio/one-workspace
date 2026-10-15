# Organization apps

`scripts/rebuild.ps1` now calls `scripts/build-org.ps1`, which packages every
currently released organization application. Every app's catalog source must
belong to `https://github.com/oneworkspacestudio/`. Download URLs, byte sizes and
GitHub SHA256 digests are verified, including cached downloads. The package
contains `licenses/components/organization-components.json` with exact provenance.

| App | Organization repository | Packaged version |
| --- | --- | --- |
| MyAccounts | [myaccounts](https://github.com/oneworkspacestudio/myaccounts) | preview-20261014.2 |
| Drive / Office | [drive](https://github.com/oneworkspacestudio/drive) | drive-upload-preview.3 |
| Photos | [drive](https://github.com/oneworkspacestudio/drive) | photos-upload-preview.2 |
| Photo Editor | [Photo-Editor](https://github.com/oneworkspacestudio/Photo-Editor) | studio-v0.1.0 |
| Git | [gitea](https://github.com/oneworkspacestudio/gitea) | workspace-20261014 |
| Music | [Music](https://github.com/oneworkspacestudio/Music) | workspace-20261014 |
| Mail | [mail](https://github.com/oneworkspacestudio/mail) | 0.1.0-preview.15; independently resolved online |
| Calendar / Contacts | [calendar](https://github.com/oneworkspacestudio/calendar) | 0.1.0-preview.1 |
| Notes | [notes](https://github.com/oneworkspacestudio/notes) | Awaiting its first release |

Drive/Photos use the newer upload fixes already published in the organization's
Workspace16 bundle. Their catalog sources name that exact organization release;
the matching source snapshot is linked from Workspace15. The older executable
in the Drive repository is deliberately not substituted for the upload repairs.
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

## Notes release boundary

Notes is being built in a separate task. An empty repository/no published binary
cannot be installed. It is not part of the default lock until its actual flags, OIDC, health and package
contract have been checked. Do not label a tile or placeholder as installed Notes.
Its first release can use the documented `add-app` manifest contract without
reinstalling existing apps. Future bundle pins must include its exact org release.
