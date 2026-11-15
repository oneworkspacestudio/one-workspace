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

The signed publication attempt stopped before creating a draft because the
Android release signer is on another PC. The owner approved a Windows-only
exception for preview.132. The standard publisher restricts that mode to this
tag and verifies the public catalog and update index.

[Preview.132](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.132)
shipped from `e927e4f` with 950 verified Windows release assets. The public
index was promoted by `602b63c` alongside the next release's recognition
changes. An anonymous preview.131 updater offered the required preview.132
installer from the promoted index. Anonymous catalog and Froststock package
downloads matched their published hashes. The installed Updates window also
verified the promoted channel. The previous signed Android APK remains available.
