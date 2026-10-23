# Remote Workspace administration

User requests remote Backups, Updates and other local-only administration, push,
release and local install. Base public main7fd5cd0 includes71 installer fixes and
the navigation hotfix. Owner on this machine is C:/tmp/workspace/Private version67;
do not apply the separate prior I:/Workspace machine's deployment assumptions.
Dedicated branch/worktree codex/remote-administration. Preserve other worktrees.
Artifacts on I: under tooling/remote-admin-20261023. Target72, no Actions/agents.

Updates portal brokers only a fixed allowlist of update-controller methods to a
verified exact-path running helper at its protected loopback URL. No client-chosen
upstream, cookies or Authorization forwarding; local controller token is injected
server-side and never returned to the remote browser. Owner plus exact Origin on
mutations; strict bounded bodies, redirects disabled, no-store responses. Local
controller remains usable. Idle legacy helpers retire through their authenticated
close endpoint; active operations are preserved. Background launch opens no local
browser. Existing verified-download/operation locks/backups remain unchanged.

A separate HttpOnly/Secure/SameSiteStrict /updates cookie has only a SHA256 hash,
owner binding and two-hour expiry in a private DPAPI file (maximum8 receipts).
It allows this update page to reconnect after manager/session restart, not general
app access. A signed-in non-owner cannot use it. Logout removes all receipts and
expires the cookie. Updated pending-operation middleware permits update/auth/
branding paths while other app access stays paused. Tunnel outages still briefly
interrupt progress; the UI reconnects and never repeats an install automatically.

Backup browser lists only local server drive/folder names with bounded entries
and reparse-ancestor rejection. No file reads or Windows desktop dialog. Selection
still goes through existing NTFS/overlap/destination checks before save/run.
Startup status/retry, Mail settings, accounts and other existing web controls were
already remotely routed. Internal stop tokens, initial desktop install/OS UAC and
device certificate trust are not public administration endpoints.

Compilation, release packaging, publication and protected local installation
pending. No automated tests requested or run; do not claim end-to-end update,
rollback, scheduled backup, clean-host or reboot acceptance from compilation.

During packaging the separate71.1 local transaction entered startup checks and
failed portal discovery with a public Cloudflare502. Added bounded two-minute
discovery retries (each HTTP request still bounded45s) so transient tunnel startup
does not abort immediately. Issuer/signature/owner validation is unchanged. This
is not evidence that a persistent DNS/tunnel error is resolved. Member/anonymous
sign-out does not revoke other browsers' owner continuations; owner logout does.
