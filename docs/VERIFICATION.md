# Verification and current limits

This is a Windows x64 evaluation installer. Results below describe this machine,
not a clean Windows VM or a production deployment.

## Exercised

- Single EXE setup wizard through a real browser: synthetic owner, separate I:
  program/file folders and F: NTFS private data, local mode, no SMTP.
- Account/bootstrap creation and dedicated portal/Drive OIDC clients.
- Portal code flow with PKCE, nonce, state and verified ID token; real owner login.
- Authenticated Drive file listing and the created Documents/Photos folders.
- Real Docs editing and DOCX save; extracted `word/document.xml` contained the
  typed text. A root-relative converter worker route was fixed during this test.
- Unit checks for catalog traversal, Windows device/ADS paths, bad hashes,
  executable/package mismatch, dependency cycles, route overlap and duplicate ZIP
  entries; wizard capability/Origin/Host checks; proxy cookie protections; Photos
  UI/API separation; Cloudflare existing-record refusal and failure rollback.

Additional final artifact and extension smoke results are recorded in the release
acceptance JSON, rather than inferred from the source tests alone.

## Limits

- Live Cloudflare resource creation is implemented but was tested with a fake API,
  not a new production tunnel or DNS zone. Existing public previews were preserved.
- No clean-host installation, Windows service, sign-in/reboot, code-signing or
  cross-machine encrypted-data migration acceptance is claimed.
- No in-place suite upgrades, extension replacement/migrations, automatic global
  logout or app-native role provisioning. Installed extensions are trusted native
  code, not sandboxes. A manifest hash is an integrity check, not publisher trust.
- The supervisor checks health at startup and restarts exited child processes;
  it does not kill later hung processes based on health checks. The Cloudflare
  connector's own reconnect behavior is used; a connector process exit is exposed
  as stopped in status and requires restarting the suite.
- Portal sessions are held in memory and end at ID-token expiry or restart. App
  sessions have their own lifecycle; account revocation is not an immediate
  suite-wide session purge.
- The initial Drive source is shared by all suite users. Adjust scopes and source
  permissions for private per-person storage. Drive SSO users are not automatically
  Drive administrators; use the separate protected recovery account when needed.
- Latest Drive/Photos are local development payloads. Existing upstream acceptance
  gates and converter reproducibility limits remain; this bundle does not remove
  those limits. Full Office/Photos format regression suites were not rerun here.
- Data and history consume growing disk space. The wizard requires fresh empty
  program/private-data folders; a failed partial install must be inspected before
  choosing fresh folders. It does not silently delete or overwrite existing data.
- Photos writes its own `Photos.log` beside its EXE, following the component's
  existing behavior. Other app logs are bounded in the protected data folder.
