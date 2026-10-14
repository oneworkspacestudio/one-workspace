# Git and Music extension (in progress)

User requested new Gitea and swingmx/swingmusic forks, MyAccounts login, shared
Music folder visible in FileBrowser, standalone Music and Workspace hosting.
Forks created at fahimc/gitea and fahimc/swingmusic; sources I:/Projects/gitea and
I:/Projects/swingmusic. Root Workspace creation task is concurrently implementing
and publishing this same workspace repository, renamed fahimc/one-workspace.
Preserve its work and use a separate build snapshot for the extended installer.
No live MyAccounts/FileBrowser installation was changed here.

Workspace changes: Config maps gitea/music to base+4/+5; fresh setup creates Music;
ConfigureIdentity registers catalog OIDC apps after owner creation; child runtime
injects Music folder/base path/exact owner sub; /git and /music native auth passes
through gateway with prefix stripping and exact route redirects. Generic add-ons
retain portal gating. Build.ps1 optional paired -Gitea/-Music embeds both verified
EXEs and manifests; other builds remain base edition. docs/GIT-MUSIC.md documents.

Gitea launcher independent module integrations/workspace: embeds official pinned
Gitea28.1.0 and MinGit2.56.0.2; verifies extracted runtime each start, locks down
NTFS data ACL, migrates only fresh database, seeds native OIDC source, captures
random recovery admin to machine DPAPI. No Gitea core rebuild claimed. Existing
or interrupted DB bootstrap without marker refused. Stock OIDC handles users,
no email auto-linking/admin group mapping; SSH/mail off; HTTP Git tokens work.

Music: Authlib1.6.5 code/state/nonce/PKCE/JWKS, stable SQLite issuer/sub map, scoped
HttpOnly secure HTTPS JWT bounded by provider expiry, ordinary roles except exact
owner. Disables first-user local/guest/pairing in OIDC mode, fixes Music root,
15-second serial change scans, Waitress on Windows, normal standalone auth kept.
Pinned v3 webclient patch fixes API/audio/worker paths under /music; package bundles
FFmpeg7.1. Public OSS optional LicenseManager None settings bug fixed, package
__init__ added for Python3.12 resources. No premium gating bypass.

Verified so far: Workspace suite/catalog tests; real signed local OIDC fixture
(wrong audience/state rejected, PKCE verifier, ordinary/owner separation, stable
links and media auth); Vite client and PyInstaller build; native Music health;
Gitea fresh SQLite/bootstrap/Git startup with dummy issuer (provider unavailable
as expected in that isolated bootstrap). Full fresh suite integration pending.
New programs use C:/Projects/.cache for fast isolated tests, I contains durable
fork outputs. Extended suite snapshot C:/Projects/.cache/workspace-gitmusic-build,
build 0.2.0-gitmusic-preview.1. Finish full installer/SSO/playback/FileBrowser checks,
publish fork source, and reconcile Workspace remote/publication without overwriting
the other task's release. No new production DNS/tunnel created here.

## Final verification and publication

Gitea ad6c9da and Swing Music52ce4d51 pushed/remote confirmed. Both have published
workspace-20261014 standalone EXE preview releases. Extended installer
0.2.0-gitmusic-preview.2 built, SHA256
c0ee4622875c7294418ad1bb04394f07ccde363b14468618e680ce485ff6d012.
Durable artifacts I:/Projects/workspace/dist/gitmusic. Exact installer fresh test
uses I:/Projects/workspace/dist/gitmusic-final-programs and -files, protected
C:/Projects/.cache/workspace-gitmusic-final/private. Loopback19390..19395. All five
processes healthy. Real browser portal/Drive/Music/Gitea callbacks passed; normal
Drive user uploaded MP3 into ordinary shared Music file, scanner indexed it, media
returned identical bytes, actual player produced successful206 range stream and
browser audio played. Prefix stayed /music; no browser page errors. Anonymous
media/APIs denied. Gitea external DB user confirmed non-admin. Snapshots source
docs/screenshots/git-music-*.png. Synthetic credentials/browser state private only.

Final installer leaves live pre-existing MyAccounts/FileBrowser and DNS unchanged.
New global empty authorized_keys created by initial Gitea default was removed only
after verifying new creation time/empty content; final launcher explicitly disables
authorized-key generation and sets SSH root under private data, verified no global
file after corrected starts. One initial Music build retry failed because its own
standalone test EXE held the output file; stopped only task-owned process trees.

Full Workspace Go tests and signed OIDC fixture passed; native final installer
and full real browser flow passed on this host, not a clean Windows/reboot/public
tunnel acceptance. Startup timeout for Gitea allows3min for cold runtime extraction.
Final clean stop/restart and release download health are checked before handoff.
Root I workspace .git is a portable skeleton while the other creation task uses
a separate source Git mirror; this task publishes only its remaining docs/runtime
changes from C:/Projects/.cache/workspace-gitmusic-git (clone of remote096e43d).
The remote initial commit already contains the earlier catalog/OIDC/proxy changes.
Do not overwrite the other task's release or git index.


Final cold install on I/C passed the full browser script with no page errors.
Clean stop verified all ports19390..19395 closed, then restart retained accounts,
library and cookies; actual player still streamed successful range206 afterrestart.
Live final preview remains19390 with synthetic owner@example.test; protected owner
password/browser state only in C:/Projects/.cache/workspace-gitmusic-final. Source
commit7bd4ac2 published; combined v0.2.0-gitmusic-preview.2 draft upload is the last
publication step. Removed only task-owned old C prototype programs after checking
no process used them, reclaiming ~1GB. Durable installer/artifacts stay on I.


## Published handoff

Completed: https://github.com/fahimc/one-workspace/releases/tag/v0.2.0-gitmusic-preview.2
(published preview,8 assets). GitHub's uploaded installer SHA256 exactly matches
local c0ee4622875c7294418ad1bb04394f07ccde363b14468618e680ce485ff6d012.
Source/binary code target7bd4ac2, plus later documentation-only handoff update.
Fork standalone releases workspace-20261014 published with matching SHA256 asset
digests and licenses/checksums. Final verification.json records actual player,
clean all-listener stop and persisted restart. Current localhost19390 preview
remains healthy; it is an isolated synthetic demonstration, not the existing
login.m8e.co.uk deployment. No work or approvals remain for this requested preview.
New installs choose their own owner in the wizard; this preview uses private
DPAPI synthetic owner credentials. Source CI not configured in base repo; all
local Go tests and signed OIDC/native browser checks passed. gh release edit is
unavailable in this CLI, so publication used authenticated GitHub PATCH draft=false.

