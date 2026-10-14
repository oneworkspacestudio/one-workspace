# Workspace suite installer

User requested a new repository and one Windows setup EXE installing MyAccounts,
custom Drive/Docs and Photos, extensible for future applications. Default name
One Workspace, private fahimc/one-workspace; optional name/visibility question pending.
Source I:/Projects/workspace (new Git repository). Existing MyAccounts/Drive source
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
