# Workspace updater and installer resilience

The live portal returned HTTP 200 and the protected transaction was complete
when this investigation began. The last repair had installed preview.137 while
keeping Froststock preview.4. Its browser still had an older maintenance helper
open. Reopening the verified installed controller restored the Updates window
without restarting the server. The current release check offered Froststock
preview.6 as a separate app update.

Installer integration 21 exercises the staged manager on the exact target
catalog before taking the live server down. Upgrading through the installer
also advances newer installed apps published in the Workspace bundle, retaining
independent pins, newer app versions, optional exclusions, and saved ports.

The portal can still read progress from a busy older helper during handoff.
Starting a helper is serialized until it registers, and closing a helper blocks
new mutations. The local page can retry a failed reopen request. Recovery uses
the same public maintenance response as the normal update path. The temporary
updater link and its capability are no longer written to the persistent log.

The owner approved Windows-only preview.138. The build carries the verified
preview.137 app inventory and payloads forward.

[Windows preview.138](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.138)
shipped from `1838c46`. The publisher verified 949 assets and promoted the
channel in `5f7af5f`. The anonymous preview.137 updater selected the required
preview.138 installer and matched its SHA-256. Go and publisher regression
tests, a fresh Windows setup, a native updater handoff, and app-only rollback
and reapplication with the exact preview.138 installer passed. Temporary test
roots were removed.

The published installer, matched against its public SHA-256, upgraded the live
server from preview.137 and Froststock preview.4 to preview.138 and Froststock
preview.6 through the normal backed-up repair entry point. Owner identity,
storage roots, and the saved Froststock port were preserved. The transaction
completed, the replacement Updates helper passed a public release check, and
the public portal reported all 46 installed services running. Public polling
recorded 78 healthy responses, 365 maintenance responses, three brief 502
responses at origin handoffs, and no Cloudflare 1033.
