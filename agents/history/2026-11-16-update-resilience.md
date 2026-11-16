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
