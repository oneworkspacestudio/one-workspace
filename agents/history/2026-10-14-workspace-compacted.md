# Workspace suite installer

User requested a new repository and one Windows setup EXE installing MyAccounts,
custom Drive/Docs and Photos, extensible for future applications. Default name
One Workspace, private fahimc/one-workspace; optional name/visibility question pending.
Source now I:/Projects/one-workspace (new Git repository). Existing MyAccounts/Drive source
and live integration previews must stay intact. Components are process-managed,
not merged into a new identity/file implementation. Docs is built into Drive.

Architecture in implementation: offline embedded payload + small protected runtime,
manifest-driven app catalog/dependencies/routes/OIDC/env; same-origin portal/Drive/
Photos; optional manual/Cloudflare public connection; local NTFS private data and
separate program/file locations; owner normal-realm admin plus protected server
recovery credential; optional email remains off by default. New ordinary users
must not become global built-in admins. Future native app packages registered by
local owner CLI with SHA256 and explicit manifest.

Reusable source components: MyAccounts released preview-20261014.2 in H repo dist/
releases; photo-enabled FileBrowser/Photos are local uncommitted development builds
under I:/Projects/tooling/releases/photos-preview. Record provenance honestly;
do not claim published Photos source or push the companion repo implicitly.

Build tool Go I:/Projects/tooling/go1.27.2/go/bin/go.exe, caches off C. F is the only
non-C NTFS data drive and has roughly100MiB free: keep large programs/bundles on I,
use small fresh F test data and clean only owned fixtures after tests. H/I are exFAT.
Current live public integration runtime F:/MyAccounts/integration-preview must not
be touched.

Implementation now includes the offline bundle installer, four-step browser wizard,
DPAPI/ACL bootstrap, ordinary-realm owner, dedicated OIDC portal/Drive clients,
process job/supervisor, gateway, CLI app import, Cloudflare provisioner and tests,
working extension example, README and development/rebuild workflows. The user
explicitly corrected the name to **one-workspace**; product is One Workspace and
private GitHub repository https://github.com/fahimc/one-workspace has been created
(renamed from the just-created m8e-workspace). An unrelated public fahimc/workspace
from 2023 was discovered and preserved. No source push/commit/release yet.

Fresh first wizard install succeeded using I: programs/files and fresh F: NTFS data.
Real browser portal login and Drive OIDC callback succeeded. Real Docs typing and
autosave produced a DOCX whose XML contained the typed text. Fixed two findings:
enable Accounts auto-signin on suite clients, and proxy root /sdkjs + /fonts for
the embedded converter worker. Photos UI must receive its own loopback Host; only
UI/static paths go to Photos, all authenticated APIs go directly to Drive to retain
Secure cookies. Added regression tests. First private fixture was stopped/cleaned;
its generated files remain only under ignored dist. Final renamed artifact build
and fresh smoke/app-extension/browser verification still pending as of this edit.

Source package script snapshots existing FileBrowser + pinned editor source without
committing that separate repo or copying agents/history/credentials. Packaging is
running; component notice files include editor AGPL and cloudflared licenses.
Before final push: finish fresh artifact smoke, verify Photos and SSO reuse, verify
post-install extension route, finalize components.lock digest, acceptance evidence,
screenshots/product name; publish source/notices with installer/seed release; inspect
CI. Preserve existing live apps. Avoid modifying Go sources during a bundle build
(one interrupted build failed because module imports were renamed mid-compile).

## Final implementation and verification

Final installer 0.1.0-preview.1 SHA256
`229fd58e71e9f082e87009f0e82ec85a197381c049ae15307c2f910165129fa9`.
Manager SHA256 `34f9edd0e46949992c17325389ca9a82772725ee928dc644041e86dd49c1b0f4`.
Exact final EXE installed into a fresh fixture, no SMTP, programs/files I: and
private data F: NTFS. Native tests/vet pass. Stop/restart, exited-child containment,
manifest import, anonymous API/extension denial and authenticated Hello page pass.
Real browser final-artifact portal login reused the Accounts session for Drive
without a second password; Docs typing/autosave was verified in DOCX XML; Photos
listed/viewed a synthetic image. Owner organization users panel is accessible.
Enable both enableAutoSignin AND enableSigninSession on cloned clients: the latter
was essential for app-to-app SSO and has a regression test. Embedded ZIP is read
with ReaderAt directly, avoiding a full ~300MB copy before wizard display.

PowerShell localhost readiness timed out on this host despite healthy services;
the smoke script now uses 127.0.0.1 and a five-second timeout, which passed. The
installer/actual browser origins remain localhost for OIDC/cookie consistency.

User explicitly prohibited GitHub Actions. Repository Actions is disabled through
the GitHub API; remove workflow templates and use scripts/build.ps1, rebuild.ps1
and smoke.ps1 locally. The first push was rejected for missing OAuth workflow
scope; unpublished root was amended before successful source publication. Do not
request expanded workflow credentials or enable Actions on a later continuation.

Component seed release components-20261014 is published privately, digest pinned
in components.lock.json. It contains core EXEs, available Drive/Photos/editor
source, MyAccounts source and source manifest with additional engine source links.
Seed manager is ignored by rebuild.ps1; that script builds the current manager.
The final installer payload changes when manager source changes, independently of
the pinned core component seed. README screenshots use One Workspace branding.

The earlier I:/Projects/workspace directory remains separate; a different ongoing
build was observed there. Preserve it. Current canonical source is one-workspace.
Fresh test fixtures are disposable and must be stopped/cleaned after release;
never touch F:/MyAccounts/install or F:/MyAccounts/integration-preview.

Remaining publication step at this edit: commit final fixes/evidence, publish
v0.1.0-preview.1 EXE/catalog/checksums/acceptance, verify remote digest, record it.
Cloudflare fake API tests pass; live provisioning, clean host, reboot/service,
code-signing and cross-machine migration are not acceptance claims.

## Published handoff

Published https://github.com/fahimc/one-workspace/releases/tag/v0.1.0-preview.1
at tested source c5592ab91fe6e959a5515300502d75e2de8d314a, also retained on
release/0.1.0-preview.1. GitHub's server-computed EXE digest exactly matches
`229fd58e71e9f082e87009f0e82ec85a197381c049ae15307c2f910165129fa9`;
size311577088 bytes. Manager/catalog/checksums/acceptance are uploaded alongside it.
Component seed assets and their server-side digests were verified as well.

Concurrent Git/Music source commit7bd4ac2 arrived during push. Preserved it by
merging into main (0238ce0), while keeping the tested 0.1 installer tag at c5592ab.
Main includes the separate 0.2 preview docs/timeout changes; do not confuse its
additional apps with the four core apps in the published 0.1 EXE. Merged Go tests
and vet passed locally. Embedded-payload smoke test accepts additional app entries
so future extended bundles are not rejected for exceeding four apps.

All owned F: private smoke folders for this task were stopped and removed; disk
space restored. Removed only the test-created Start Menu shortcut after checking
its target. Synthetic documents/images remain under ignored dist fixtures on I:.
Existing live MyAccounts/Drive and the separate Git/Music test installation remain
outside this task's cleanup. No Actions runs; repository Actions remains disabled.
