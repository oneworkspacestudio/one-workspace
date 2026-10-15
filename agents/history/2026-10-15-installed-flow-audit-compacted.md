# Installed suite flow audit — incomplete browser acceptance

User asked to open/use every installed app, test cross-app flows and highlight
issues. This task authorizes testing, not fixes/commits/pushes/releases. Target
is the user-authorized C:/tmp/workspace installation19, not separate legacy live
MyAccounts/FileBrowser installations. Previous organization-apps handoff remains
authoritative for package provenance and earlier release/startup evidence.

Current Git: Workspace main tracks studio/main, clean before this handoff. Mail
checkout behind studio/main3; its frontend hash matches the installed Mail15
served file exactly. Drive's unrelated untracked preview files preserved.

## Browser blocker

Selected existing Chrome binding could not name a session/create a tab: browser
request-header policy unavailable; OpenAI policy service request returned403.
Read browser troubleshooting and retried once after independent work; still
blocked. User asked asynchronously to reconnect extension/restart Codex. No
browser warning bypass, alternative automation or authenticated UI acceptance.
Do not claim apps were opened/used successfully in this audit. Resume with the
same selected browser after restoring tools; read CUA documentation on resumption.

## Evidence and findings

Local artifacts outside Git:
I:/Projects/tooling/workspace-flow-audit-20261015/REPORT.md plus sanitized route,
asset, anonymous-API and synthetic Mail compose results/script. No user content
changed or messages sent. No screenshots because browser unavailable.

HTTP: all child services running.13 routes local+public26 responses consistent;
all24 initial static assets200. Anonymous Drive/Music401, Git identity403,
native Mail/Notes401. Portal Mail/Notes APIs redirect to login rather than401.
No claim of exhaustive authorization/security tests or actual SMTP delivery.
C: only about0.50GB free; ordinary files and private databases reside there.
No new failed write induced; substantial storage reliability risk remains.

Four synthetic checks execute exact installed Mail compose JavaScript with a
minimal DOM and stubbed JMAP (not browser E2E): forwarding drops attachments;
HTML-only forwarding omits original body; reopened drafts retain attachments
internally but render no attachment rows/removal controls; CC-only send is
blocked by To-only validation. Source app.js compose lines53–57. Installed/source
frontend SHA256487869d25abe615a4cca5fbc4e138b6a692109965cde256c6cc541a8ff9d1b7b.

Portal source+HTTP confirm requested app/deep link lost on sign-in: gate302 to
/auth/login without return path; callback always redirects /. Full login was
not run. Preserve allowlisted relative return paths in any future fix.

Existing integration gaps documented separately from reproduced defects: Notes
has no Drive picker/direct Calendar API (ICS only); Mail lacks Contacts search
and Drive attachment picker; installed Mail keeps private engine while Calendar/
Contacts share another; Calendar doesn't deliver invitations; no consistent app
switcher/return-home; portal logout leaves app sessions alive (documented).
Do not auto-migrate real Mail data to resolve the backend difference.

## Resume

Restore CUA and perform labelled synthetic browser journeys for all11 entries:
sign-in/deep links/logout, Drive uploads/resume/duplicates, Office save/history/
collaboration, Photos/Studio shared file edits, Notes persistence/ICS, Calendar/
Contacts CRUD, Mail drafts/forward/internal delivery, Music indexing/playback,
Git private test repository/commit. Capture actual visual proof. Full remaining
checklist and reproduction steps are in the artifact report. No audit findings
were fixed, no source implementation changed, and nothing was published.

## Interactive retry completed after cleanup

User renewed request to test all apps and list fixes. Chrome3 was unavailable;
documented browser discovery found reconnected Chrome4. Native browser controls
now worked, so the old tool-policy block is superseded for this pass. Current
installed manager19 unchanged. User-authorized target remains C:/tmp/workspace,
not the legacy MyAccounts/FileBrowser installs. Prior cleanup removed old previews;
no new installations were created for these tests.

Canonical local report (outside Git):
I:/Projects/tooling/workspace-flow-audit-20261015/run2/ISSUES-AND-COVERAGE.md.
JOURNAL.md in same directory records evidence/test state. Screenshots remain local.

All11 catalog apps covered, with sampled functional journeys passing in9 and2
blocking failures. Calendar event created/saved/reloaded/reopened; Contacts record
created/saved/reloaded; Notes rich content persisted, revisions restored and text
attachment downloaded byte-identical. MyAccounts Users and actual Manage accounts
link target work (read-only admin check; no identities added). Drive QA folder,
small image/text/audio uploads Completed without Resume. Photos upload/numbered
duplicate/search/favourite/folder albums work; favourite appears in Drive pins.
Studio opens own Drive file/saves to original; edited image visible in Photos.

Office ordinary files: DOCX edited/autosaved/reopened; named baseline, subsequent
edit, read-only history preview and successful restore verified on-disk. XLSX
formula1+2 and cached3 verified; PPTX2slides verified; one-page PDF created/rendered/
saved with PDF1.7header. Same-account two tabs Live2devices: first-tab edit appeared
in second without manual reload; final synthetic textRW. No multi-user/device/
simultaneous-merge acceptance. Generated1-second WAV uploaded via Drive Music,
automatically indexed; actual playback currentTime0.8039/duration1/pausedfalse/
readyState4/noerror. Playback tab closed to stop repeat.

## Current issue list

- High: Git initialized private repo creation500. Server log: pre-receive hook
  attempts WSL /bin/bash, execvpe fails; init commit push rejected. Git dashboard
  remains empty afterward. Clone/commit/push not accepted; fix Windows hook/runtime
  selection and validate without WSL dependency. Do not infer clean-host portability.
- High: Mail SSO again stops on Chrome Security error at Accounts OAuth authorize
  path. Tool cannot attach. User asked asynchronously for screenshot; no reply at
  finish. Current classification/cause not established. No bypass, alternate host,
  session reuse or token workaround used; no send/receive/draft/browser acceptance.
  Keep security tab for user inspection when possible; no new Google report sent.
- Medium: Studio Image Rotate90 converts320x200 to320x320 with added TRANSPARENT
  margins (appear white in Studio), rather than200x320. Saved dimensions/alpha
  checked. Whole-image/layer rotation semantics need correction/clear labeling.
- Medium: opening Notes while portal signed out completes normal provider SSO but
  lands at portal /. HTTP/source destination-loss finding now browser-reproduced.
- Medium: internal subject/generated identities shown in Drive/Studio/Git/Music/
  Calendar. Keep immutable authorization identity, correct friendly display only.
- Medium: Contacts retains Calendar shell controls and incorrect login wording.
- Medium: Calendar switcher only offers3apps; One Workspace links to GitHub.
  Correction to earlier source inference: a switcher exists but is incomplete.
- Medium documented gap: portal logout leaves Photos/Drive app session accessible;
  observed signed-out portal and authenticated Photos simultaneously. Provide
  coordinated logout/clear scope; do not present as undisclosed vulnerability.
- Low: new PPTX template layout labels/document language Chinese in English UI.

Four Mail compose defects rerun against byte-identical installed15frontend, SHA
unchanged from prior handoff; all four reproduce in stubbed synthetic logic checks.
Still not browser/delivery proof. Known Notes Drive/Calendar, Mail Contacts/Drive,
Calendar invitation integration gaps remain documented separately.

## State and boundaries

Labelled QA event/contact/note/Drive folder plus root photo+duplicate/favourite and
Music WAV retained for review; small synthetic contents only. No existing user
documents/messages edited or deliberately opened individually (app listings/thumbs
visible); no email sent. No permanent data deletion, source implementation fix,
commit/push/release. Services remain healthy; C:about24.25GiB free. Earlier low-disk
finding is resolved and must not be listed as current. Preserve uncommitted cleanup
AGENTS/history changes. Test report lists untested large/resume uploads, independent
device/users, complex PDF/Office, clean-host/standalone and Mail acceptance gates.
