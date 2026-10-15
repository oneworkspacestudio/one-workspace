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

Published release: v0.3.0-preview.16 with Mail0.1.0-preview.13. Mail release published
first at https://github.com/oneworkspacestudio/mail/releases/tag/workspace-mail-v0.1.0-preview.13
Mail SHA256 4ab76b887b0d33aa9b2a6b062b8ccb4f00dc96a87f5fc282c2d81ec4b889a44e.
Installed dedicated Workspace manager and Mail updated with backups, normal stop,
verified child exit, checksums/config update and hidden restart; retained existing
core+Mail app set and ordinary files, without adding Git/Music to that install.
Both forwarding hops now verified and external SMTP reachable; authorized Gmail
receipt checked. Anonymous compiled download and lease-duration delivery checks completed successfully.
## Published artifacts and additional checks
Public release v0.3.0-preview.16 at
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.16
Published program/catalog/bundle assets all matched local size and GitHub SHA256.
Bundle SHA256 7ca90638fa1be194ffe2382756f17568f727f345248ea9e95fcf4e5b7a0af906.
The compiled installer downloaded all500MiB anonymously, verified every package,
then resolved organization Mail0.1.0-preview.13. Exit0; log outside Git at
C:/Projects/.cache/workspace-public-release/dist/public-download-verification-v16.txt.
A separate live anonymous test used the older bundle, downloaded newer published
Mail, verified the complete catalog with the override, and extracted the new EXE
into a fresh temporary directory. Exit0. No existing installation used for that
extraction. Native fresh Mail provisioning fixture also passed with ephemeral
loopback ports and NTFS temp data, without sending messages or changing live Mail.

Rebuild seed now pins preview.16 and current Mail. Optional manual manifest helper
accepts an explicit version; generated current-version manifest/hash/provenance
checked. Current Mail docs/notices now link organization source. A small external
Mail-NOTICE.md release asset corrects the older embedded notice's former source
URL; immutable binary/bundle assets were retained and release notes explain it.
Drive/Photos source remains the matching preview.15 source asset linked in notes.

Installed manager and Mail match release/bundle hashes. Local ignored artifacts
are version-suffixed preview.16 to avoid the existing locked canonical installer.
No unrelated preview or live MyAccounts/FileBrowser process was stopped.

Final running-suite verification: startup ready8/8, current Mail incoming state
reachable, no hosting/identity errors. Autonomous router observation exceeded the
original30-minute lease with no manual changes. The final authorized Gmail test
was ingested and visibly received at14:55:20 BST. Evidence is stored outside Git
under C:/Projects/artifacts/mail-repair-20261015/. Actions remain disabled in both
organization repositories. Both corresponding source trees are published; future
builds use the corrected public organization Mail notice and manual helper links.

Mail follow-up: Mail0.1.0-preview.14 distinguishes failed SMTP from an unavailable
or inconclusive third-party probe in the banner. Networking/engine unchanged.
Live anonymous organization resolution against published suite16 successfully
downloaded Mail14, reverified the entire catalog and extracted its override;
exit0, log outside Git dist/mail14-overlay-verification.txt in release staging.
This proves the online installer can pick a newer Mail independently of its
bundled13. Deliberate offline16 still uses13; no suite asset was rewritten.
Installed Mail updated to14 with verified stop/backup/hash/config/restart. A
further authorized Gmail delivery completed at15:00:13 BST. Final Chrome sign-in
is currently stopped on a browser security page, requiring the user's inspection;
do not bypass it. Prior signed-in receipt screenshot proves the lease-window test.

User screenshot confirms a Google Safe Browsing phishing interstitial. Mail15
corrects identity cues (own One Workspace logo, independent-app notice, actual
MyAccounts host). No security bypass or ownership grants were made. Public Google
status reports no unsafe content, which does not override Chrome. Search Console
property is not verified for the current user, so a Google report-error request
with public callback path and remediation details is prepared. User approval is
required before submitting; no report sent yet. See Mail router-recovery handoff.
Installed/released Mail15, five asset digests checked; normal public entry page
visually confirms own branding/version15. Sign-in callback warning remains until
Google resolves classification. No Go networking/auth changes in14/15. Installer
resolver14 real upgrade test passed; same supported manifest integration1 for15.
