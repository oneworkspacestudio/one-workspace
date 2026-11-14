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

## Verified release

Froststock preview.1 shipped in [Workspace preview.128](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.128) from `0be2499`.
The signed publisher verified 960 assets. The complete 47-app catalog has 939
packages. The verified preview.127 bundle served as the integration baseline
and retained existing app payloads byte for byte. The first channel push raced
with another main-branch commit. Rebased only the promotion on the newer main,
pushed `88541cd`, and verified the standard anonymous raw-main URL and updater
after cache propagation. Anonymous downloads of the catalog and all five
Froststock packages matched the published digests.

Go tests and vet passed using an NTFS test temp directory. The exact managed
installer passed fresh setup, restart, optional omission and opt-in upgrade.
Synthetic rollback removed the new app and preserved a Drive file. Two OCR
fixtures produced the expected dates and byte-matched photos. The Drive UI
listed all three ordinary files, and the packaged vision worker read the
expiry fixture on WebGPU. The app picker, owner session, unauthenticated API
denial and mismatched-account/CSRF checks passed. No physical phone tested.

The Windows Workspace installation upgraded from preview.122 through its
protected updater. The transaction completed with verified recovery. Installed
Froststock packages and health matched the release, and the existing owner,
issuer and Drive location were preserved. The remote app route leads to
Workspace sign-in. Synthetic services were stopped and isolated roots removed.

An older recovery copy failed cleanup verification and was retained. The new
recovery copy was verified. The installed My Cinema preview.5 health endpoint
returned HTTP 503 during startup and again after the supported app-only retry.
The update retained the Cinema package. Froststock, My Drive and MyAccounts
were running after the retry. No physical phone or owner-scan fixture was used.
