# AI Fieldnotes

Workspace preview.145 introduces the optional Creative & Development app
`fieldnotes` (service preview.1). It is a supervised native apphost service with
embedded AUI reader, MyAccounts owner sign-in, a dedicated private database,
encrypted provider settings and standard recovery ownership. The published
preview.144 seed preserves other app contracts and payloads. Installer
integration 24 covers the new app's package, licenses and launch contract.

Search uses DeepSeek Search or the optional Workspace Web Search server API.
Source selection precedes DeepSeek generation. Drafts are validated against
numbered source IDs and persisted with citations. The app contains original
starter articles on loop engineering, a proposed five-stage maturity framework
and WCAG-based interface evaluation. Requests run on demand. Keys and provider
responses containing credentials do not appear in the reader or exports.

Android package `uk.co.m8e.fieldnotes` version 1.0.0 reuses the bundled reader,
ships with no Internet permission and imports bounded JSON packs through the
Android file picker. Imported reading state is private and commits before UI
acknowledgement. The signed publisher's `-Fieldnotes` flag pairs the APK with the
required signed Workspace companion and verifies its provenance and public
channel entry. The web-to-Android handoff uses explicit file transfer.

Use only synthetic NTFS private folders during acceptance. The baseline
preview.144 bundle is pinned by public SHA-256. Do not install on an owner
machine absent a separate request. Finish release verification and record
any outstanding acceptance here before reporting managed delivery.
