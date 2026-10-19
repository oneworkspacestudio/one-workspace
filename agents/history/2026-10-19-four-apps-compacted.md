# Four native apps — work in progress

## Authorization and storage

User requests implementation of the Four Apps master brief, local automated and
browser tests, GitHub pushes and actual releases. No GitHub Actions (standing
user restriction overrides brief). Four original public repositories created in
oneworkspacestudio: my-tasks, my-pdf, my-monitor, my-search; Actions disabled.
No current-task commits/pushes/releases yet. Preserve unrelated Oct15/16 histories.
Use I: for source/build/download/cache; protected disposable NTFS fixtures only
under a named C:/Projects/.cache root. Preserve live installation/newest rollback.
C free about16.6GiB at start; about16.6GiB still during implementation.

## Source and architecture

Actual suite source I:/Projects/one-workspace; branch codex/four-apps. Inherits the
unpublished new-default-app updater repair from the Oct18 handoff. Suite Go module
canonicalized to github.com/oneworkspacestudio/one-workspace, including imports and
ldflags. New public pkg/apphost reuses Apache suite platform/Hosting conventions:
issuer+subject identity, code+PKCE/state/nonce OIDC, sessions+CSRF, standalone local
bootstrap disabled in managed mode, NTFS/reparse checks, SQLite WAL, immutable
blobs, bounded migrations, structured logs, encrypted chunked Argon2id/AES-GCM
consistent backup/new-root restore (restored browser sessions removed). A same-EXE
worker reuses suspended-start Windows Job with512MiB/process limit1/CPU60s and
parent90s context cancellation. App worker joins before Host.Close.

ADRs under docs/adr/001..004 and docs/FOUR-APPS-DELIVERY.md. Source directories
I:/Projects/my-{tasks,pdf,monitor,search}. Each is its own git main, remote org repo.
Apps currently require suite v0.0.0 placeholder. External go.work under
I:/Projects/tooling/four-apps overrides it; must pin real public immutable suite
pseudo-version, remove dev replacements and independently build before release.
Plan: push suite foundation to feature branch, derive pseudo-version, release
apps, pin actual published hashes in organization-components.lock.json, then
final suite main/release. Never force or retag. No private release assets.

## Implemented app source

- Tasks backend: lists/private/shared owner/editor/viewer, task/subtask fields,
  optimistic versions, bulk/drag positions, recurrence separate rules+template,
  idempotent next occurrence/undo, IANA DST wall time, exception dates, bounded
  calendar rules (rrule-go1.8.2), persisted in-app reminders/retry/dedupe, saved
  views/search/ICS escaping/folding. New /api/time resolves datetime-local by IANA
  zone, rejects spring gap, chooses earlier autumn fold by default. Frontend AUI
  editor, views/board, sharing, reminders, saved views, Markdown. Recent polish
  adds date pickers, quick Enter add, hides sharing management behind button.
- PDF backend: per-user immutable uploads/docs/annotations, jobs/cancel/status,
  same-EXE pdfcpu0.16.2 worker inspect/merge/collect/remove/rotate/optimise/blank/
  properties. Parent paths only after ACL; password in RAM only; crash reset jobs.
  Frontend PDF.js6.4.299 local worker/cmaps/fonts/wasm/iccs, pdf-lib1.17.1, page/
  zoom/find/text/PNG, images→PDF, mark geometry/ink/comments/highlight/text,
  image stamp support being added, annotations and true redaction export brand-new
  raster PDFs (no original objects, explicit loss-of-text/accessibility warning).
  Browser Drive picker uses /drive/api/users?username=self, resources, actual
  download?source=&file= and POST resources override=false; compare OIDC subject.
  Save copy only; no unsupported overwrite/version, OCR or Office claims.
- Monitor: one service, loopback-capability central registry, owner-only APIs,
  checks HTTP/TCP/discovered process/Windows service, memory/disk/CPU counters,
  samples30d, incidents/ack/quiet/maintenance, TLS14-day warning, encrypted webhook,
  safe aggregate status off by default. SSRF checks ALL DNS answers and pins dial
  IP; rejects loopback/private/linklocal/metadata/CGNAT except EXACT registry URL.
- Search: browser federation through actual source session APIs for Drive metadata
  (Docs/Photos types), Notes, Calendar, Contacts, Tasks. User identity match before
  snippets; no admin crawler/foreign DB access; Vault excluded, Mail/Chat excluded.
  Bounded responses/concurrency, source statuses/partial results, opt-in per-user
  recents in own SQLite. Notes/Tasks pagination; Drive has no offset contract and
  must offer refine/open-source rather than repeated pagination. Calendar contacts
  use collection APIs; links currently calendar date / contacts app, not record
  selection. Highlight/filters/suggestions still need polish and honest limits.

All frontends use actual suite tokens/Hosting adapter and .aui generation via
@codedia/parser0.2.1, React19.3.1/Vite8.3.4/TS. Original matching SVGs added to
app public/docs assets and suite branding; transparent light/dark same readable
palette. PNG exports/screenshots/docs still pending. AUI compiler owns JSX;
adapter/integration React.createElement only. Do not rerun external scaffolding
scripts blindly: several are one-shot and would duplicate/revert later edits.

## Suite integration done so far

New routes /tasks,/pdf,/monitor,/search, ports base+18..21. Runtime passes owner,
home and Monitor-only registry capability; generic OIDC provisioning handles new
clients. Seamless probes require matching issuer/subject for new apps. Captured
path-scoped sessions revoke on root logout; browser cookies cleared. Registry
GET /_integration/monitor requires exact loopback+derived capability, exposes no
secrets/user file paths. Public /monitor/status only exactGET (app opt-in needed).
/_workspace/apps keeps old tile fields, adds ID/Route/Version/issuer/subject.
Search discovery must account for launcher URLs /_workspace/open/{id}, not just
literal /app/ routes. Conditional global home search and Ctrl/Cmd+K app switcher
search added. CLI remove-app supports four independent apps stopped, preserves
private data/clients/backups, records opt-out so updates don't reinstall them.
Build.ps1/build-org.ps1 accept/embed/download new app artifacts+notices+contracts;
new update groups APPENDED after old groups (tests formerly index notes group8).
Organization lock pins not yet added, because releases are not published.

Notes actual source I:/Projects/my-notes, shallow org notes main. Added authenticated
GET /api/identity returning own issuer+subject, plus test; needs new Notes preview
release (version currently0.1.0-preview.4) and suite lock update. No other Notes
changes. Do not use workspace-notes-fix (that is an old suite checkout).

## Verification to date

Suite full go test ./... passed after module migration/updater repair. Initial
apphost auth/local/CSRF/backup/restore/migrations/Host tests and RSA fake-provider
OIDC PKCE/nonce/audience/issuer/expiry/replay tests PASS. Later suite/apphost tests
passed; adding update groups at front caused one indexed-fixture test failure,
fixed by appending groups; distribution tests pass again. Re-run final full gates.
Tasks backend tests PASS ownership/sharing/revocation/version conflicts/DST fall
recurrence/undo/dedupe/reminder retry/ICS. New time gap/fold tests not run yet.
PDF tests PASS actual bounded TestMain --worker operations/malformed PDF/context
cancel/source unchanged, per-user docs/annotation/job ACL and queued native job.
Encrypted/signed/redaction text-extraction and full Drive browser gates pending.
Monitor tests PASS private/mixed-DNS/metadata SSRF, exact registry trust, status
off, encryption/tamper, incident threshold/dedupe/recovery. More maintenance/
quiet/TLS/restart tests pending. Search four node contract tests PASS discovery/
vault exclusion, wrong-account denial, dedupe/text strings/pagination, revoked or
offline source. Need browser injection, actual integrated permission/deletion.
All four AUI/TypeScript/Vite frontends compiled once successfully. Several polish
edits since; Tasks latest build passed; others need final rebuild. Logs under
I:/Projects/tooling/four-apps. No clean-host/reboot acceptance claim.

## Current processes and browser

Only CURRENT-task preview started: MyTasks.exe at http://localhost:29108,
PID135036 (recheck receipt/process!), executable tooling/four-apps/build/MyTasks.exe,
data C:/Projects/.cache/workspace-four-apps-20261018/previews/tasks. Receipt
tooling/four-apps/tasks-preview.json, stdout/err adjacent. Bootstrap synthetic
owner@example.test, existing disposable password available in external launch
command; do not record credentials here. Started with --no-browser, hidden.
Browser CUA handle fourTasks, in-app tab17, signed in to this synthetic account.
Browser create/edit/complete/Completed/undo task passed. Last view Inbox with one
synthetic task. Screenshot revealed overly large controls and sharing clutter;
polish.py addressed CSS widths and management/editor layout. Rebuild EXE to a NEW
filename while old one running, restart only proven preview PID before deleting
old output/data. Owner live apps were NOT used as test fixtures.
Browser guidance restored; selected iab2. cua.rewriteDocumentation only if resumed
after summary; native computer control disabled. Current CUA supports screenshot,
AX, Playwright read-only evaluate, file chooser (read upload docs first). User brief
explicitly asks Playwright automated tests; standalone test harness acceptable for
synthetic fixtures. No browser warning bypass or new credentials entry through UI.

## Precise outstanding work / risks

1. Finish PDF image mark validation and native metadata: Mark.image added, PNG
   magic/dimensions check; compile/gofmt/test. Current pdfcpu properties-remove
   only removes custom properties/XMP, not all standard Info: fix and test instead
   of advertising full metadata removal. Browser metadata removal must delete old
   unreferenced Info/Metadata objects too (pdf-lib saves unreferenced objects).
   PDF.js6 uses loadingTask.destroy (not doc.destroy); isEvalSupported removed.
   CSP now PDF-only wasm-unsafe-eval + font blob/data; no JS unsafe-eval.
   Redaction/annotations raster export coordinates protected against rotated-view
   marks; image marks need actual UI test. Bound aggregate raster pixels/pages and
   add cancellable browser export; thumbnails/split ZIP still need implementation.
2. Monitor maintenance currently suppresses alerts but creates down incidents and
   counts downtime; fix to maintenance state/exclude uptime. Quiet period currently
   records only successful alerts; record attempts too, show delivery status.
   Service API currently mgr.Connect/OpenService requests excessive privileges;
   replace read-only SCM/service query. Add configurable timeout/history graph.
3. Task recurrence existing-rule edits silently ignored: either implement explicit
   rule/exception update semantics or reject unsupported changes clearly. Recheck
   editable membership inside write transactions, not only before BEGIN. Browser
   board currently completed column empty unless completed/all view; proper full
   list board query needed. datetime fold choice UI, list rename/delete controls,
   bulk move destination field and tags filter improvements desirable.
4. Search highlight/type/date/tag/owner capability filters, suggestions, stale
   result race/account change; expose unsupported sources metadata-only, no fake
   mail/chat APIs. Notes identity new preview pin. Actual federation fixtures.
5. Add suite new-default/opt-out/legacy/new-app OIDC/independent update/removal and
   registry/logout tests. Implement uninstall UI if needed (CLI currently only).
6. PNG/icon exports using existing sharp; docs/README/screenshots/licenses/SBOM/
   API/recovery/acceptance/performance/security. Bundled Node modules path from
   load_workspace_dependencies is in primary runtime; contains sharp likely.
7. Independent clean Go builds after pinning public foundation commit, package
   exact Windows executables + compatibility + update contracts + checksums,
   virus/vulnerability review, isolated full suite fresh/upgrade/backup/restore/
   seamless browser tests; actual GitHub publication only after gates truthful.
8. Clean only task-owned NTFS fixtures after verified stops; C free check; keep
   live suite and its most recent valid backup. Prior Hosting public review at19007
   is unrelated retained demo (see Oct18 storage history), do not strand its routes.

Tooling: dot-source I:/Projects/tooling/four-apps/env.ps1, then override GOWORK to
external go.work for new apps. Go I:/Projects/tooling/go1.27.2/go/bin/go.exe.
Set TEMP/TMP and GOTMPDIR to named NTFS fixture for tests; never generic C caches.
Set PYTHONUTF8=1 or explicit encoding='utf-8' on ALL Python Path read/write calls
(default cp1252 caused an interrupted one-shot patch; resumed only its tail).
Git safe.directory exact I:/Projects/<repo>; pass with git -c. No globalwildcard.

## Foundation gate update

Shared owner-marker validation, immediate standalone setup session, GUI error
reporting, object checksum/SQLite restore validation and Windows icon resources
were added. Signed PDF fixture has an actual detached CMS signature independently
verified with OpenSSL (self-signed trust deliberately not claimed); encrypted PDF
correct/wrong-password and signed original preservation tests pass. PDF split ZIP,
1GiB account quota, garbage collection and atomic cancelled-result handling added.
Monitor maintenance/quiet-attempt/read-only service query/timeout implemented.
Search fixed actual Drive terms + repeated scope contract and launcher registry
URLs; identity checked again after source query, late recents rejected, bound
requests refuse account changes. Tasks date gap/fold, transactional ACL checks and
clear unsupported existing-series rule edit error added.
New tests PASS: online/offline manager update adds four defaults, preserves optouts
and custom ports; scoped Monitor registry secrecy/loopback capability, full new-app
server logout/cookie cleanup, stopped removal preserving data and clients.
Full suite tests and vet passed before the final checksum-restore hardening;
repeat final gates. Current app-suite PNG updated with four new marks and a new
content-addressed README URL; original illustrated showcase preserved. Source
foundation publication to feature branch is next; no product release yet.

## Final delivery correction (19 October)

Earlier resume lists above describe intermediate work. Final independent app
sources pin public apphost b1c1f801ccb2 via immutable pseudo-version; no local
replacements/GOWORK for release. My Tasks/PDF/Monitor/Search are original Apache2
repos; actual Notes org repo is notes, native preview.5 adds /api/identity.

Final native standalone binaries passed eight Chrome checks including same-SPA
account switching, private-draft clearing and expected-account rejection. Final
managed native package update passed all eighteen checks in scripts/four-apps-e2e.mjs:
one MyAccounts password, six app SSO flows, actual Drive/PDF import/save-byte-copy,
Monitor registry/status privacy, real Search sources/deletion/XSS strings,
second-user isolation/viewer sharing/edit denial/revocation/non-owner Monitor403,
and app-initiated full Workspace logout/cookie cleanup. Full suite tests/vet pass.
Four encrypted standalone backups/new-folder restores passed. Actual old36 install
updated with new37 controller added four defaults, preserved ports/identity/old
omissions; separate Notes update enables identity connector. Old-controller bridge
requires published exact manager/catalog; postpublication native check remains.

Packaging: PowerShell7 local builds, pinned dependencies, full license texts,
compatibility/build metadata, original matching SVG/PNG/ICO and synthetic UI
screenshots. New runtime npm audits0; Go source scans0 reachable/imported findings,
module-only findings retained. Existing Notes moderate sprintf-js precision DoS
has no upstream fix. No clean-host/reboot/SCM/independent-audit claim; preview only.

Release outputs/evidence remain external on I:/Projects/tooling/four-apps.
Publish four preview.1 apps and Notes preview.5, pin actual published ZIP digests,
build full org19 catalog, push main fast-forward, publish suite preview.37. No
Actions. Root has intentional NoShortcuts fixture isolation, strict19 catalog
check, provenance-after-all-downloads and managed E2E script changes. Preserve two
unrelated untracked older histories. Final release/cleanup outcome follows below.

Cleanup: only named C:/Projects/.cache/workspace-four-apps-20261018 fixtures after
verified process stops. Owner C:/tmp/workspace and newest rollback preserved;
unrelated Hosting acceptance/public-preview19007 preserved. All large builds,
downloads/source/cache/packages on I:. Do not delete unknown owner/user data.
