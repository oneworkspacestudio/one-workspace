# Optional Mail component

One Workspace Mail is an independent Gmail-style client licensed AGPL-3.0-only.
Source: https://github.com/oneworkspacestudio/mail/tree/main/workspace-mail

Mail.exe embeds the **unmodified official Stalwart 0.16.26 Windows engine**,
not a locally rebuilt Rust engine. Stalwart is copyright Stalwart Labs LLC.
Its source is AGPL-3.0-only OR LicenseRef-SEL; the community AGPL option applies.
Exact engine source: https://github.com/stalwartlabs/stalwart/tree/v0.16.26
The app's local build script pins and verifies its official Windows archive.

The Windows platform helpers retain One Workspace's Apache-2.0 license.
Go dependencies include go-oidc and go-jose (Apache-2.0), oauth2 and x/sys
(BSD-3-Clause). See the source module and its NOTICE.md for provenance.
No affiliation with Google or Stalwart Labs is claimed. See Mail-AGPL-3.0-only.txt
for the AGPL license and docs/MAIL.md for deployment and preview limits.
