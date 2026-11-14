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
