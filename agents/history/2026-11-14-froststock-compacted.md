# Froststock in One Workspace

The user corrected the original delivery target to the One Workspace ecosystem
and My Drive. This branch started at public main 456d754 without changing the
older local suite or One Workspace OS. The earlier Sites preview is separate.
Use its backup export for any existing scans, then import in the managed app.

Froststock uses the supervised apphost and MyAccounts owner session. It is an
optional General app, ID froststock, route /froststock, port 9048 (manager +58),
with Accounts and Drive dependencies, installer integration 18, and distinct
app update ownership. The binary embeds the mobile UI, OCR, and optional vision
runtime. Allow outbound model weights only from the vision worker. Tesseract
reads English labels first, with optional SmolVLM 256M on WebGPU.

My Drive is the authoritative inventory. Each pair commits as a directory
containing item.jpg, expiry.jpg, and item.json in Files/Froststock. A browser
outbox protects unsent pairs across reloads. Stable IDs support retry, revisions
prevent stale edits across devices, and only the owner can call app APIs.
The app uses a bounded Drive folder and existing Drive visibility controls.
Use Workspace Backups for ordinary Drive files. App-only updates do not roll
back the contents of Drive.

Release target Workspace preview.128. Validate the complete managed catalog and
actual installer with synthetic Accounts/Drive roots under the task-local C
cache, with programs/release assets on the build drive. Use the mandatory signed
publisher and verify anonymous channel discovery before claiming shipment.
