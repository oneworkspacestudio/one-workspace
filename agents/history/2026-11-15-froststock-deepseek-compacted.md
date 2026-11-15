# Froststock DeepSeek label recognition

Froststock can send the item photo and the expiry photo or selected crop to
DeepSeek Flash. Saving an API key in Settings makes DeepSeek the first reader
for captures and saved-photo retries. Date-only retries send only the expiry
image. Removing the key returns to local-only scanning.

The managed service uses an owner-only, CSRF-protected API and stores the key
in its private host data with Workspace's platform secret storage. The key is
never included in settings responses, inventory data or browser storage.
DeepSeek uploads go directly from the service to its fixed Chat Completions
endpoint. No photos are needed to save settings.

The browser checks the suggested expiry against the date transcription and
existing date parser. DeepSeek date suggestions and visually inferred names
require review. Missing fields and provider failures use the existing local
scanner. Capture, background queuing and My Drive records keep their format.

Mobile browser checks with a mock upstream covered activation and removal,
two-image scans, date-only crops, packing dates, missing years, and local
fallback. A live DeepSeek scan requires an owner-provided key.
