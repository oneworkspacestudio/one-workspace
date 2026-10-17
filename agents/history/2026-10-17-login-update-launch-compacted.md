# App sign-in and Windows update restart repair

## Request and diagnosis

User reported a second Drive login and an update/automatic rollback that both
failed to launch Private/Workspace.exe with Access is denied. The live suite is
the protected `C:/tmp/workspace` installation. It was still preview24, predating
the preview27 automatic native app handoff. Its attempted preview28 update had
restored24 and completed the rollback, but could not restart. ACL inspection
permitted owner execution; manual `version` and normal launch worked. There was
no relevant Defender detection. The idle old maintenance process was closed
through its local capability API, and24 was restarted while repair work continued.

Windows host jobs can forbid CREATE_BREAKAWAY_FROM_JOB. The old unconditional
flag reproduces ERROR_ACCESS_DENIED in an isolated real Windows nested job even
with correct file permissions. The fixed launcher checks current job limits,
requests explicit breakaway only when permitted, and otherwise respects inherited
host policy. The manager itself remains outside its separate app kill-on-close
job. No security protection, ACL, credential, DNS or autostart change is needed.

## Source repair

- `internal/platform/maintenance_windows.go`: job-aware maintenance/runtime start.
- `internal/suite/update.go`: launch preflight before app downtime; distinguish
  completed data rollback with failed startup from unfinished recovery. Completed
  rollback now gives restart guidance instead of an inapplicable Recover action.
- Windows regression tests reproduce the old failure and exercise restricted and
  silent-breakaway hosts. Updater tests cover safe preflight abort and guidance.
- New seamless tests cover every native OIDC app, existing sessions, cookie/token
  isolation, Drive/Docs/Photos/Studio sharing, canonical routes and server errors.
- `docs/UPDATES.md` explains the launch failure and completed rollback recovery.

Full Go tests passed for final source on isolated NTFS; vet passed before the
final small silent-breakaway flag refinement. No Actions. Existing untracked
installed-preview22 and runtime-restart histories are preserved.

## Local artifacts and acceptance state

Local hotfix version **0.3.0-preview.28+restart.1**, derived from the verified public
preview28 bundle with only its manager/catalog replaced. It is not a GitHub
release and must not be represented as official preview28 bytes. Later public29
will still compare newer. Artifacts/receipts/scripts are outside Git under
`tooling/workspace-login-update-repair`; one named NTFS cache is
`C:/Projects/.cache/workspace-login-update-repair`.

Final manager SHA256 `6a5958a9c04adad46e489dd7a043ccfd252bdfbdbdf8055b295789130306a5b3`.
Final maintenance helper SHA256 `88157aa3afcbea33e9be9ad565a7c0e7087b10f76ed22de5616ddf95b96e32b3`.
Final local bundle SHA256 `52912ce26d96177f1b1139c0cdf3d524319e1dadd57c284d7b91c1f0b0ccf962`.

Fresh published24 fixture used synthetic accounts, a saved note and ordinary file.
Earlier candidate upgraded successfully under a deliberately restrictive Windows
host job, retaining credentials semantically, owner/storage bindings and content.
An initial harness assertion wrongly compared randomized DPAPI ciphertext; it was
corrected to compare decrypted logical credentials without exposing values. The
stopped first fixture was removed after exact-path/reparse/process checks.

The Python HTTP sign-in harness stopped at Drive callback: Domain=localhost
cookies are handled differently by its cookie jar. Provider auto-authorization
itself redirected correctly. This is not browser acceptance. The failed harness
receipts are retained, and no cookie/state/security check in an app was weakened.
The second disposable fixture is entirely under the named cache fixture2 root;
only its synthetic owner password was reset for diagnostic resume. Live user
credentials were not reset. Chrome's existing real Workspace tab is signed in
on24 and available for user-requested read-only flow verification after repair.

## Earlier checkpoint (completed below)

Exact final native upgrade/rollback acceptance is running via native-final.py,
with task-owned restrictive host processes and protected disposable data. Wait
for its receipt before deployment. Planned deployment uses the normal verified
backup/update pipeline with the fixed helper, retaining user files, credentials,
storage/domain and optional app selection. Then verify existing Chrome sessions
and local/public health, preserve exact source/artifacts, and stop/clean only the
owned fixture/cache. No commit, push, tag or release is authorized by this repair
request. No full email delivery, editing/playback, reboot or service acceptance.

## Final repair and deployment completed

Exact restart.1 native acceptance passed: restrictive-host upgrade and real
rollback restored the original manager, existing client credentials, saved note,
ordinary document, storage/identity bindings and app subset. Every fixture process
stopped; its disposable Private/Apps/Files tree was removed after path, reparse and
process checks. No live test account or password was created/reset.

The normal verified update pipeline upgraded the user's24 installation to local
restart.1, retaining all21 ordinary files byte-for-byte and logical credentials.
Chrome confirmed the existing Workspace session reopened automatically. The
requested read-only review then found a real additional Music problem: an older
PWA worker served the cached player shell at the native login URL, producing a
blank page. The same backend URL correctly returned a302 without that cache.
The approved Music executable was retained; the withheld newer binary was unused.

Added `music_navigation.go`, embedded `assets/music-worker-navigation.js` and
Music app-open handoff. The authenticated handoff outside Music's worker scope
retires only its exact /music/ registration; it retains caches and every other
app's data/worker/cookies. The next worker response prepends a scoped live-navigation
handler while retaining original asset/audio caching. Offline shell fallback is
allowed for ordinary Music pages, never for authentication callbacks. Worker
scope stays /music/. Native OIDC callbacks, secrets, state/PKCE checks and standalone
Music are unchanged. Go proxy/handoff tests and the Node worker simulation cover
scope isolation, cache headers, original worker retention, live auth navigation and
offline fallback. `docs/SIGN-IN.md` documents the integration. All final Go tests
and Node worker checks passed. Final vet status is in completion receipts.

The accepted native updater source is byte-identical in the second hotfix, which
was also installed through the normal backed-up update pipeline. The authoritative
active version is now **0.3.0-preview.28+restart.2** in both protected manager and
catalog/programs manager. Final manager SHA256:
`00006902608fd1691b75b584f61075bd7d2c925a2766bd20747d009c84e0387d`.
Maintenance helper SHA256:
`1a31389c14c93e91891459f05e2985e093d73c0f62914add11e85a79dcda5cec`.
Local full bundle SHA256:
`c5240fa881309de593e2bc18e9eb7ae0774e6b17aa6c7e3933b113c316764a51`.
The retained pre-update backup is
`update-20261016T233440-e287e548713c` under the live Private/backups folder (restart.1
state). The preceding24 backup `update-20261016T231955-75eb694e6a0f` also remains,
as do the earlier rollback/recovery archives. No automatic backup pruning.

Both local/public portal and health returned200; every runtime service and
dedicated Cloudflare connector reported running. All21 user files, logical
credentials, storage/domain/owner/ports, autostart false and all12 installed app
entries were retained. Backups scheduling remains off. Public GitHub preview28
has not changed; these are clearly labelled **local, unpublished hotfixes**. No
commit, push, tag, release or Actions run occurred. A later public29 will compare
newer than the local28 hotfix label.

Final Chrome review on restart.2 passed all12 installed entries: Drive, Docs,
Photo Editor, Photos, Calendar, Contacts, Notes, My Tube, Mail, Git, MyAccounts
and Music, plus internal Backups. No credentials were entered and no sign-in
button was required. Music's formerly blank native login URL was separately
revisited and returned automatically to the player. Contacts/Calendar initially
paint their landing/loading view before session hydration; they reached the
authenticated views automatically. The temporary review tab was closed and the
original user Workspace tab left refreshed on restart.2. No file/message edits,
email delivery, audio playback, fresh logged-out Mail/Vault, reboot or service
acceptance is claimed. Vault is not part of this12-entry live catalog; do not
start the old18935 synthetic fixture from ambient browser context.

Sanitized receipts and exact source/artifacts are retained under local
`tooling/workspace-login-update-repair`; restart.1 artifacts are in `final`,
restart.2 in `music-final`. Deployment preservation receipts containing filenames
stay inside protected live `Private/maintenance/{restart-repair,music-sso-repair}`.
Final cleanup removes only the named disposable build cache after source/artifact
hash verification and process/reparse checks. Preserve all live data/backups and
unrelated cloudflared instances. Future publication needs a freshly versioned
complete installer/bundle and explicit authorization; don't relabel these bytes
as the already-published preview28 release.

## Completion confirmation

Final `go test ./... -count=1`, `go vet ./...`, Node worker simulation and Git
diff whitespace checks passed. Exact restart.2 source and three local artifacts
were preserved with flush/read-back SHA256 verification; native test harness Go
sources were retained separately. The exact named NTFS cache was then removed
after confirming no reparse points or active process references. Cleanup receipt
reports about20.1GB free on C:. Live installation, user data, every real backup,
and unrelated tunnels were retained. Final local/public health after cache
removal remains ok/all services running; update journal is complete. Work is
finished locally, source changes uncommitted, public release unchanged.

## Publication authorized

User explicitly requested publication and release. Prepare v0.3.0-preview.29 from
the accepted repair, with the same13 organization app entries and component pins
as public28. Only Workspace is changed; no separate component release is needed.
Build a newly versioned full offline installer, manager, bundle and exact committed
source archive on isolated NTFS; preserve notices/provenance and verify all package
hashes. Publish to the org without Actions. Do not replace28 or modify the live
installation while publishing. README/build defaults and release notes identify29.
