# Froststock

Choose **Froststock** under **General** when installing or updating One Workspace. Open it from the Workspace app picker and sign in with the Workspace owner account.

Take or upload one photo of the food and a second photo of its date label. The scanner returns to the first step as soon as the pair is queued. Keep scanning while your browser saves the pairs to My Drive and extracts the text. A temporarily disconnected browser keeps unsent pairs and retries when the app is open again.

The dashboard shows items due within seven days, items past their printed date, and items that need a closer look. Open an item to verify the photos and edit the label details. Mark used items to move them to History.

Your photos and inventory records live in **My Drive → Froststock**. A second device signed in to the same Workspace sees the same inventory. Each item folder holds `item.jpg`, `expiry.jpg`, and `item.json`. Keep those files together. My Drive access controls who can view the underlying Drive files. Froststock itself permits the Workspace owner.

On the expiry-photo step, select **Crop expiry photo after capture**. Move the box or drag its corners around the date and its “use by” or “best before” label, then tap **Use crop**. The scanner remembers your choice on this device. **Use full photo** queues the whole image. For saved items, use **Crop & read again** by the expiry photo. The original remains in My Drive, so you can adjust a crop later.

Date recognition runs in the browser. Froststock uses bundled English Tesseract and retries difficult dates with the tiny PP-OCRv6 detection and recognition models, adjusted contrast and rotation. The optional SmolVLM 256M vision model uses WebGPU and downloads about 200 MB of weights from Hugging Face. The inference stays on the device. Confirm dates marked **Needs review** against the original photo. Change the numeric date order in **Settings**.

Choose **Read labels again** on a saved scan to reread both labels with the current crop.

For the standalone preview, export a backup in its Settings, then import it in Froststock under One Workspace. It writes each imported pair to My Drive. Managed app updates keep the Drive files. Use Workspace Backups to protect My Drive and the app’s private sign-in data.
