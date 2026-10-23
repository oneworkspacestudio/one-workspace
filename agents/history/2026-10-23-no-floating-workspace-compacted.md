# Floating Workspace navigation removed

User requested removal across every app, push, release and local installation.
Published v0.3.0-preview.71.1 from source6c5a37d933b26ffbbd6f364b47b8ff4cac785f78.
The shared navigation no longer creates a floating launcher. Existing header Apps
menus, shared picker events and Ctrl/Cmd+K remain. Dialog scroll fits short phones.
AGENTS.md now prohibits floating app-content navigation. No individual app payload
changes: public71 baseline retained all25 app contracts and134 other payloads.

Verification: Go tests/vet, JS syntax, actual embedded navigation on25 synthetic
catalog routes at phone/desktop widths, short-phone scrolling/Close/focus/keyboard,
and isolated native Windows protected67-to71.1 update. Actual Social reproduced
the old overlay; updated five bottom tabs/Create are clickable, Social header Apps,
Slicer picker and Drive checked. Synthetic route coverage is not25 real-app tests.
No owner functional test, physical print, GitHub Actions or clean-host reboot.
See docs/evidence/shared-navigation-no-floating-20261023.json.

Local first attempt failed on public Accounts OIDC discovery502, then the normal
recovery guard required more space. Normal Recover restored67 successfully.
Older rollback copies were encrypted, integrity-tested and file-verified before
archival; the newest usable full backup remained on protected NTFS. Ordinary NTFS
compression on recovery-only backup files retained logical SHA256 values. No
journal bypass, WOF compression or operating-system setting changes.

A concurrent authorized My3D update installed Workspace72 (source3b13a2977d32),
which includes this same navigation fix with CRLF line endings. Avoid a downgrade
to71.1. Independently verify the installed manager against the public72 manager
and catalog digests, served JS/CSS against that binary, and navigation against the
tested71.1 source after line-ending normalization. The evidence records actual
local services/app count, identity/paths/ports preservation and nonblocking startup
warnings. New app updates from other chats are not this task's test claims.
The approved older ISO installer was moved toI with SHA256 verification.

Work isolated from concurrent My3D and remote owner-update changes. Never reset or
stage unrelated shared work. Public tag/assets immutable. Use repository-history
and reconcile actual local paths/versions; earlier other-host handoffs differ.

