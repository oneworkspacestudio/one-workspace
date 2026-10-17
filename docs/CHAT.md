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
its active connections; Workspace sign-out also clears Chat cookies.

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

The initial package is a preview. Multi-user messaging, attachments, reconnection,
Windows service/reboot, database restore and mobile clients need acceptance before
a production claim. Calls, TURN/SFU hosting and external push services are not
configured by this installer. Tinode is not an end-to-end encrypted messenger by
default. Keep the server and its backups private.

Source and component notices: https://github.com/oneworkspacestudio/chat
