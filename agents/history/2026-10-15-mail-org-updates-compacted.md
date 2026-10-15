# Organization Mail resolution and incoming delivery repair

User requested the running Mail repaired, Gmail delivery checked, fixes published
in the organization, and Workspace installing the organization version of Mail.
No GitHub Actions. Read Mail's router-recovery handoff for the network diagnosis.

Online setup now independently resolves the latest published non-draft Mail
release at oneworkspacestudio/mail, including previews, after verifying the latest
Workspace bundle. The Mail release must provide Mail.exe and workspace-mail.json
(schema1/integration1, matching product version/tag/executable). Asset URLs and
API paths must stay in the exact repository; GitHub SHA256/size checks cover the
manifest and EXE. Missing/invalid/incompatible latest Mail fails closed; offline
mode deliberately uses the embedded bundle. Existing installations are not
silently auto-upgraded. Public repo resolution needs no token or gh executable.

A downloaded Mail replaces its bundle package via an owned temporary override,
with updated catalog version/hash/size/exact organization release source. Verify
and Extract both rehash overrides; a missing/corrupt override cannot fall back to
the old bundled file. Archive path/link/duplicate checks remain in force. If the
verified bundled executable already matches the organization release, reuse it.

Tests: go test ./... and go vet on NTFS release staging, JS checks in Mail. Mock
HTTP tests cover newer/identical Mail, corrupt hashes, drafts, wrong repository,
missing manifests and incompatible integration. Catalog tests cover override
extraction, tampering and missing-file failure. Installer progress includes the
independent Mail check under app verification. Reused existing Drive/Photos and
all other component binaries, retained extended Git/Music payload and licenses;
no unrelated Drive work included.

Release plan: v0.3.0-preview.16 with Mail0.1.0-preview.13. Mail release published
first at https://github.com/oneworkspacestudio/mail/releases/tag/workspace-mail-v0.1.0-preview.13
Mail SHA256 4ab76b887b0d33aa9b2a6b062b8ccb4f00dc96a87f5fc282c2d81ec4b889a44e.
Installed dedicated Workspace manager and Mail updated with backups, normal stop,
verified child exit, checksums/config update and hidden restart; retained existing
core+Mail app set and ordinary files, without adding Git/Music to that install.
Both forwarding hops now verified and external SMTP reachable; authorized Gmail
receipt checked. Final anonymous downloads and lease-duration checks pending.