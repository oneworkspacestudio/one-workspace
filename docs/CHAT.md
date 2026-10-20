# Chat

Chat is a One Workspace fork of Tinode with its native messaging engine and web
client. Open the Chat tile to reuse your MyAccounts browser session.

## What is included

- Direct conversations, group topics, message history and attachments.
- One Workspace colours and a navy/mint/silver Chat icon.
- A locally bundled client: React, translations and icon font do not need a CDN.
- A single Chat.exe starts the gateway, Tinode engine and private PostgreSQL
  runtime without Docker, a Windows service or a separate database installation.

Workspace registers a separate OIDC client and exact `/chat/auth/callback`. Chat
validates the identity token, state, nonce and PKCE before issuing single-use
engine tickets. Identity links use issuer plus subject; email is never the key.
No second Chat password is created. Logout closes the Chat browser session and
its active connections; Workspace sign-out also ends the current Chat session before clearing its cookies.
Workspace43 gives Chat a separate browser hostname and updates its existing
MyAccounts callback without changing its client credentials or account links.
Old `/chat` links redirect to the new address, retaining the conversation link.
Host-only cookies and per-account browser caches keep Chat separate from other
apps. Directly opened active attachments are sandboxed.

Managed Cloudflare installations automatically allocate an unused first-level
`chat` hostname (or a numbered alternative), retain it in private settings and
repair its owned tunnel route on startup. Unrelated DNS and tunnel routes are
preserved. Local installations use the other loopback hostname, so Workspace
and Chat do not share a cookie host. This does not configure calling or TURN.

## Finding people and opening conversations

Workspace35+ connects the MyAccounts directory to Chat automatically, checking
new accounts every15 seconds. Search by display name, username, verified email,
MyAccounts ID or native Chat ID. People can be found before they first open Chat.
Chat6 refreshes the open Find query every5 seconds while visible and online, so a
new account appears after directory synchronization without retyping the query.
The Find field searches as you type and also has an explicit Search button.

Open By ID to paste a native user/group/channel ID or tinode:///id/ link, or copy
your own ID and QR code. Invalid entries show guidance rather than blanking Chat.
Disabled members disappear from search and their gateway sessions are revoked.
The private synchronization credential is never sent to the browser. Account
links use issuer+subject, preserving existing messages and IDs across updates.
Only verified email is matched; names/usernames also work without email validation.
The current limit is10,000 members and50 matches per query.

## Storage and backup

Private/chat holds PostgreSQL files, ordinary attachment files, protected secrets
and regenerable runtime software. Workspace's stopped-state backups include this
folder. PostgreSQL WAL crash recovery is retained; do not copy its running data
folder as a standalone backup or open the same folder from two processes. Moving
protected credentials to another Windows machine requires a planned restore.

Chat is optional in Advanced installation. Existing installations can add it
through the installer’s app-selection flow; ordinary updates preserve deliberately
omitted apps. Future Chat releases are offered from oneworkspacestudio/chat.

## Candidate5/6/7/8 changes and scope

Chat5 adds a full-height branded layout, app switcher, synchronized custom lists,
account/conversation drafts and static offline shell. Restored sends use one
durable client ID. Chat7 also preserves this ID through the actual manual Retry
menu: a real dropped-acknowledgement test leaves one recipient copy after reload.
The prior menu recreated the message and could duplicate an already delivered
send; this defect was reproduced and disclosed for Chat6/Workspace44. Older
ambiguous queued messages remain pending with guidance rather than receiving
an invented ID. The original reply/document data is retained. A deleted
message's retry returns Gone. Uploaded bytes and conversation membership are
checked before media access or attachment linking. The account and security
surfaces route centrally managed name/password controls to MyAccounts.

Chat6 refreshes an open directory search after new-account synchronization. New
group messages retain the original server-authored recipient audience: every
original recipient must acknowledge delivery/read. Removing an unread member
cannot manufacture all-read, and retry after membership changes retains one
message. Read implies delivery and text labels accompany icons. Unknown legacy
or removed-recipient counters are conservatively Sent; historical receipt
retention after deletion/rejoin and full Message Info/privacy remain incomplete.

Chat8 waits for the actual user ID before enabling By ID actions and creating
its QR. The QR refreshes when identity arrives; stale asynchronous work is ignored.
Owned conversation callbacks are detached on navigation so late data/pin work
cannot crash Settings or update a different conversation. Re-entry coordinates
with an outstanding leave, and stale subscribe completion cannot select an old
chat. Real delayed-login QR-pixel and ordered delayed-frame navigation tests pass.

Chat5/6/7/8 require Workspace43. The app-update manifest rejects older managers;
the executable also checks the manager capability before touching Chat storage,
protecting complete-bundle update paths. Use the paired Workspace update first.

This candidate is **not full WhatsApp parity and does not send end-to-end
encrypted messages**. Its separately tested OpenMLS native/browser proof is not
part of Chat.exe. Linked-device trust, encrypted media/search/recovery, complete
social features, calls and native clients remain unfinished. See the
[feature matrix and actual test scope](https://github.com/oneworkspacestudio/chat/blob/master/docs/whatsapp-parity/IMPLEMENTATION-STATUS.md).

The package remains a preview. Separate synthetic Windows accounts passed browser
OIDC, search, ID/QR lookup, two-way direct/group messaging, exact attachment bytes,
history reload and disabled-user checks. Database startup waits for WAL recovery.
Isolated local cross-origin logout/socket revocation passed. Full network/expiry,
clean-host, Windows service/reboot, fresh-machine database/key restore and mobile
acceptance remain open. [Recorded Chat checks](https://github.com/oneworkspacestudio/chat/blob/master/docs/workspace/ACCEPTANCE.md). Calls, TURN/SFU hosting and external push services are not
configured by this installer. Tinode is not an end-to-end encrypted messenger by
default. Keep the server and its backups private.

## Public45 integration verification

Actual published installer/bundle downloads passed default-online fresh setup,
supported rollback/reapply and the paired42→45 upgrade. Synthetic browser tests
cover SSO, preserved IDs/history/document bytes, delayed-login QR, ordered late
data navigation, lost-ACK Retry, offline/logout and seven viewports in both themes.
The operational installation now runs45/Chat8, with its permanent dedicated
Chat address healthy through Cloudflare. Owner conversations and owner browser
sign-in were not test fixtures. [Recorded boundaries and hashes](CHAT-PUBLIC45-ACCEPTANCE.json).

Source and component notices: https://github.com/oneworkspacestudio/chat

## Phone layout in Chat9

Chat9 uses a compact branded header, search and filters, a separate compose action,
readable working navigation tabs and actual conversation timestamps. Workspace
has a fixed close action above its scrolling app grid. See the
[mobile reference and verification scope](https://github.com/oneworkspacestudio/chat/blob/master/docs/whatsapp-parity/MOBILE-LAYOUT.md).

## Offline continuity in Chat10

Workspace49 includes Chat10. After an online visit, the cached app and history
can reload at the same address during an outage. Chrome network-disabled checks
cover composing, durable queueing across reloads and exact-ID synchronization
when the network returns. Local emoji and six Workspace stickers are included.

Recent existing one-to-one conversations automatically prepare signed device
permissions and public Nostr discovery. With Internet still available, browsers
can reconnect over WebRTC while the main server is down, without a separately
installed signaling service. Queued, directly delivered and server-saved states
are displayed separately. The desktop rail uses shared design-system buttons.

Permissions last at most two hours and public relays/direct routes may be
unavailable. No TURN service is included. First-time visits, new sign-ins,
groups/files/calls and mobile background delivery still need server support.
Acceptance used real Chrome profiles and synthetic accounts on one Windows PC;
separate-network and physical Android/iOS checks remain open. Normal server
messages are not application E2EE.
