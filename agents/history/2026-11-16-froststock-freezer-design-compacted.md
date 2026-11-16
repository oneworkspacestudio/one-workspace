# Froststock freezer redesign

The owner supplied a mobile reference for a warm neutral, inventory-first
layout. **My freezer** opens to the full inventory with photo cards, date
status, search, date filters and the main scan action. Phone navigation has
Freezer and Scan. Settings opens from the header gear and includes sorting,
history, backup, DeepSeek and the Workspace app picker.

The managed Go service still embeds the same vanilla web app. Existing
inventory, scan queuing, expiry review, crop, recognition and data formats
are preserved. The old inventory and expiring links still resolve. The UI
uses saved item photos as thumbnails and only shows the review filter when
there are records to review. No reference image is bundled.

Browser checks at 320px, 390px and 1440px covered overflow, portrait and
desktop composition, all filters, search, settings, history and sorting,
the full inventory, item details and both scan entry points.

Froststock preview.7 is the only changed app in the Workspace preview.139
proposal. The managed builder pins the verified preview.138 bundle to
retain installer integration 21 and the current updater and tunnel fixes.

The packaged release passed an isolated fresh setup and an app-only update
from preview.138, including mobile navigation, search, uploading and
cropping scans, and the DeepSeek key settings. Rollback retained the
original photo bytes and inventory. Go, recognition and publisher checks
passed. The 47-app catalog and bundle hashes were verified against the
preview.138 base.
