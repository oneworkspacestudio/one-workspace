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
Chat preferences use their own browser-storage namespace; directly opened
attachments are sandboxed on the shared Workspace origin.

## Finding people and opening conversations

Workspace35+ connects the MyAccounts directory to Chat automatically, checking
new accounts every15 seconds. Search by display name, username, verified email,
MyAccounts ID or native Chat ID. People can be found before they first open Chat.
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

## Scope

The package remains a preview. Separate synthetic Windows accounts passed browser
OIDC, search, ID/QR lookup, two-way direct/group messaging, exact attachment bytes,
history reload and disabled-user checks. Database startup waits for WAL recovery.
Full network/expiry/global-signout, clean-host, Windows service/reboot, database
restore and mobile acceptance remain open. [Recorded Chat checks](https://github.com/oneworkspacestudio/chat/blob/master/docs/workspace/ACCEPTANCE.md). Calls, TURN/SFU hosting and external push services are not
configured by this installer. Tinode is not an end-to-end encrypted messenger by
default. Keep the server and its backups private.

Source and component notices: https://github.com/oneworkspacestudio/chat
