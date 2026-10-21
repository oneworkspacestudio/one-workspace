# Mail installation failure and readable startup messages

User screenshot: Start Mail warns ordinary User role missing after fresh install.
Pinned engine native fixture reproduced exact failure: configureIdentity ran in
recovery mode before engine defaults create Role records. Mail17 now starts
normally after saving explicit listeners, initializes native defaults, connects
MyAccounts, verifies and saves. Retry found/fixed listener name read-only patch
failure too. No new permission allowlist or administrator role for mailbox users.

Mail source isolated from studio/main ca2bf240 into I:/Projects/mail-setup-permissions;
original Mail checkout has unrelated Linux edits left untouched. Mail source
3eeb4c4, tag workspace-mail-v0.1.0-preview.17 planned. Shared blue Mail artwork/resources
sync and obsolete brand-test assertion corrected. Go all tests/vet and JS checks
pass. Native fresh + interrupted-before/after-defaults tests cover token mailbox
access, isolation, denied administration/wrong audience, and restart stability.

Suite60 starts from clean8968aeb (public59 delivery). Startup maps technical Mail
identity/setup failures to plain actions and preserves details in protected
logs/mail-setup.log. Existing-domain decisions stay visible. Source pins Mail17;
other24 app contracts/packages are preserved from public59, including Monitor3.
Mailbox identities/custom policy remain intact. No existing live data fixture,
external email or GitHub Actions. Publication uses own organization fork, not
upstream Stalwart. No issue/PR/contact to upstream maintainers.

Artifacts I:/Projects/tooling/mail-setup-20261021. Go/private fixture root
C:/Projects/.cache/mail-setup-20261021. Build from committed sources. Native59
broken fresh setup ->60 protected upgrade/SSO UI, owner deployment, publication
and cleanup to be recorded after verification. Preserve current owner data,
all unselected app versions (including installed Chat9) and verified recovery.

Full native browser check found a first-use restart race too. Mail now waits for
authenticated engine registry readiness before returning setup success or serving
configured startup. Immediate access, retry and restart Go tests pass; updated
Mail17 source 7e17754 is pinned before regenerating the unreleased60 bundle.

## Consumer-facing incoming-mail status follow-up

Native59-to60 protected upgrade and real MyAccounts browser inbox passed at
desktop/mobile widths after engine readiness was added. Mail17 source7e17754c
is public. Normal backed-up owner58-to60 installation passed: all24 services
healthy, Mail configured with no identity error, encrypted Mail config unchanged,
all24 other app contracts/packages including Chat9 retained, recovery verified.

Operational startup also exposed an existing separate incoming-network warning:
the detailed router/IP/TCP25/SMTP/ISP diagnostics were still copied directly into
the consumer checklist. Version61 keeps the genuine warning, states Mail is ready
on this computer and asks the user to open Mail to finish the incoming connection
check. Detailed diagnostics stay in protected logs/Mail setup; no false Internet
delivery claim. Reachable/unknown/checking messages are also plain language.
Unit regressions reject network jargon in the startup detail. Suite/setup/update
tests and vet pass. Mail17 and all25 app bytes/contracts remain unchanged from60.
Use new61 tag/artifacts;60 tags/assets stay immutable. Final61 build/publication,
manager-only normal local update, proof and cleanup are to be recorded.

## Final verified delivery

Mail17 (7e17754c) and Workspace60 (f1e5e90) are public. Final Workspace61
(c31eaf3) is public with unchanged Mail17 and all25 app contracts/payloads.
Every published asset size and GitHub SHA256 matches the local artifact.
Source tags/assets remain immutable; final documentation commits are separate.

Native59 broken setup ->60 protected update and real MyAccounts browser inbox
passed with no browser errors at desktop/mobile widths. Signed-token native
engine tests prove immediate access after setup, resumed setup, restart, user
isolation, no administrator role for mail users, wrong-audience denial and custom
policy preservation. Earlier login-driver failures were corrected to handle both
provider SSO and login forms; an actual engine-readiness race was then fixed.

Normal backed-up owner58->60/Mail17 and manager-only60->61 updates completed.
All24 service entries healthy; Mail configured, no identity error, encrypted Mail
configuration unchanged, identity/paths/ports/selections retained, all25 app
contracts/packages retained from60 including Chat9. Latest recovery manifest
verified. Incoming Internet receiving still has a genuine connection warning;
startup now shows a plain next step instead of router/IP/protocol diagnostics.
No external email sent, no live DNS mutation as a test, and no owner data fixture.

Verified final installer in Downloads/WorkspaceSetup-0.3.0-preview.61-windows-amd64.exe.
Evidence and synthetic screenshots: docs/evidence/mail-setup-deployment-20261021.json,
docs/screenshots/mail-setup-fixed-*.png. Native services/browser stopped. Named
C task cache removed with explicit contained file cleanup after no reparse/no
active-reference checks. Durable artifacts remain on I. Preserve current owner61
and its newest verified recovery copy. No clean-host/reboot or delivery proof claim.
