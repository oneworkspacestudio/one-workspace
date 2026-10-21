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
