# PDF files in Drive

One Workspace preview.40 includes My PDF preview.3. New managed PDF uploads,
image conversions and processed PDF copies automatically save as ordinary files
in **Drive / My PDF**, within your authorised Drive scope. Leave My PDF open until
the document says **Saved to Drive**. The source and destination path appear
beside the document. Duplicate filenames become `(1)`, `(2)`, etc.; existing files
are never replaced by automatic saves.

## Existing documents and recovery

- Existing private-library documents and imported copies have **Save to Drive**.
- **Retry Drive save** retries a failed connection or transfer. Pending requests
  persist across browser/server restarts and resume when My PDF is open.
- The same MyAccounts session connects Drive. Drive permissions and quotas apply;
  My PDF does not use an administrator token or bypass Drive authorisation.
- **Save selected copy to Drive** chooses a different folder and filename.
- After a lost upload response, a retry checks the destination bytes before
  making another copy. Expiring leases prevent simultaneous tabs duplicating work.
- The automatic-save checkbox applies to new uploads and processing. Turning it
  off keeps those new documents in the private library.

Your private My PDF library remains the working copy. Removing a document there
does not delete the ordinary saved Drive file. Editable drawing marks remain in
My PDF until you explicitly create an annotated PDF copy. Original and signed
Drive files remain unchanged. Raster exports lose searchable text, forms, links,
accessibility and cryptographic signatures in the new copy. Split ZIPs, page PNGs
and extracted text are separate downloads. Standalone mode remains local.

## Update and compatibility

Use Workspace Updates to install My PDF preview.3 or use the preview.40 complete
installer. Both carry the same verified Windows binary. The private library is
upgraded to schema2 for save receipts; downgrade requires restoring compatible
app data from the verified pre-update backup. Current normal updates preserve
the library, identity, file locations and the latest verified rollback.

## Verification scope

Local native Windows build, all PDF Go tests/vet and eleven frontend contracts
passed. Ten isolated managed browser checks verified MyAccounts sign-in, Drive
connection and automatic upload, duplicate naming, native rotation, browser page
numbers, custom-folder save receipts, transfer failure/retry, lost-response
reconciliation without duplicates, local-document save and reload persistence.
Saved files were downloaded through Drive and compared. A final native smoke
also checked the completion notice and exact downloaded bytes.

Tests used only disposable synthetic accounts and a separate NTFS data folder.
No owner installation served as a feature fixture. These checks do not certify
every feature, hostile-input security, a clean physical host, Windows service
startup, actual reboot, Internet mail or Vault decryption. Existing preview and
component limitations remain in RELEASE-STATUS.md and ECOSYSTEM-REVIEW.md.
