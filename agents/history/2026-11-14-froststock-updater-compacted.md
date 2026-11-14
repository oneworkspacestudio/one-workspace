# Froststock installer and updater

Workspace preview.131 shipped Sunny & Moon preview.4 while the Froststock
listener change was under review. Preview.132 builds on the released preview.131
catalog and ships Froststock preview.2 with the compatible manager.

Apps published by Workspace take their update pins and payloads from the verified
suite catalog. Suite tags no longer appear as standalone release warnings for
Froststock, Game Studio, Web Search, Pet, or Sunny & Moon. Independent repositories
retain their compatibility checks and review links.

Installer integration 19 hands older controllers off to the verified installer
before they can install Froststock without the loopback-capable manager. The
backed-up reinstall carries an older, installed Froststock forward while preserving
a newer pin or an existing exclusion. Other selected apps retain their ordinary
repair behavior.

The standard Windows build verified the SHA-256 of the published preview.131
seed and all retained payloads. The catalog audit passed for 47 apps and 939
packages. The installer passed a fresh NTFS setup with Froststock and a separate
IPv4 listener on the same port. The live backed-up update moved preview.131 and
Froststock preview.1 to preview.132 and Froststock preview.2. The journal is
complete, Froststock and Caddy are both listening on port 9080, and the external
Froststock route requires sign-in. The public update check shows Froststock,
Game Studio and Web Search as current and still offers the separate Sunny & Moon
update. Synthetic services and roots were removed.

The public publisher stopped in its required Android signing step because this
machine has no configured release signer. No preview.132 release was published
or promoted.
