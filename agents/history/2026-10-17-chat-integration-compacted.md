# Chat integration / preview33

Owner requested a Tinode fork in oneworkspacestudio with Workspace styling, matching
icon, MyAccounts login and installer integration. Public fork oneworkspacestudio/chat
is local I:/Projects/workspace-chat, default branch master; Actions disabled.
App source4affea9378c18b01feebfd4d22a828683267350a. Component targetv0.1.0-preview.1;
ZIP SHA25638a2350a0124276128840234be04179376b4f3a217a49f93dc0b203d59fb007c.

Added Chat.exe to org-pinned build, full14-entry catalog, icon/README/guide, advanced
selection, per-app working folder, base port+16 (Vault local TLS uses+15), separate
OIDC callback /chat/auth/callback, seamless probe/login and sign-out cookie cleanup.
Chat uses /chat under the existing Workspace tunnel, so no new DNS record needed.
Its update group owns only chat executable/license packages. New apps still require
installer selection; updates preserve intentionally omitted apps. MyAccounts client
registration uses existing generic installer flow; other apps/source unchanged.

App details/history in chat repo: Tinode engine, vendor webapp, matching SVG/CSS,
local React/i18n/font, state/nonce/PKCE, native issuer-subject links and one-time
engine tickets, HttpOnly sessions and private loopback PostgreSQL runtime. Single
Chat.exe supervises children with Windows jobs. Private/chat contains database,
ordinary attachments, machine-DPAPI credentials and regenerable software. Existing
backup includes this folder only after processes stop; no restore acceptance yet.

Workspace and Chat Go compilation and browser production build pass (webpack size
warnings only). No automated or multi-user runtime tests requested/run. Do not
claim messaging/upload/reconnect/lifecycle/recovery/clean-host acceptance. Owner
live suite was not used as a development fixture; no preview runtime started.

Build/download/output under tooling/chat-build. Initial I frontend extraction slow;
named C:/Projects/.cache/chat-build held only public build/npm software. Preserved
outputs/generated source on I, verified hashes/no reparse/process use and removed
that exact cache (374MiB logical software). C free observed~6.3GiB afterward; no
claim about unrelated free-space changes. No owner data/backups deleted.

Component publication, full preview33 package/install and public download results
will be appended after completion. Source history is not evidence of successful
upload/publication. Existing unrelated older untracked histories must be retained.
