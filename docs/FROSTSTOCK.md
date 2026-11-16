# Froststock

Choose **Froststock** under **General** when installing or updating One Workspace. Open it from the Workspace app picker and sign in with the Workspace owner account.

Take or upload one photo of the food and a second photo of its date label. The scanner returns to the first step as soon as the pair is queued. Keep scanning while your browser saves the pairs to My Drive and extracts the text. A temporarily disconnected browser keeps unsent pairs and retries when the app is open again.

**My freezer** shows photos, quantities and label dates. Use **Expiring soon** for items due within seven days or **Past date** for older items. **Needs review** appears when there are scans to confirm. Search by name, category or label text. Open any card to check the photos and edit the details.

Use the gear to open settings and change the sort order. Items marked as used move to **Settings → Your data → View history**, where you can return them to the freezer.

Your photos and inventory records live in **My Drive → Froststock**. A second device signed in to the same Workspace sees the same inventory. Each item folder holds `item.jpg`, `expiry.jpg`, and `item.json`. Keep those files together. My Drive access controls who can view the underlying Drive files. Froststock itself permits the Workspace owner.

On the expiry-photo step, select **Crop expiry photo after capture**. Move the box or drag its corners around the date and its “use by” or “best before” label, then tap **Use crop**. The scanner remembers your choice on this device. **Use full photo** queues the whole image. For saved items, use **Crop & read again** by the expiry photo. The original remains in My Drive, so you can adjust a crop later.

To use DeepSeek for new scans and **Read labels again**, open **Settings → DeepSeek**, paste an [API key](https://platform.deepseek.com/api_keys), and choose **Save key**. DeepSeek Flash then reads the photos first and suggests the food name and expiry date. It receives the item photo and the expiry crop if you set one. A date-only retry sends only the expiry photo or crop. Froststock stores the key in private Workspace data. Choose **Remove key** to return to device-only reading.

If DeepSeek is unavailable or leaves a field unreadable, Froststock uses its browser scanner. It uses bundled English Tesseract and PP-OCRv6 for difficult food titles and dates. The optional SmolVLM 256M vision model uses WebGPU and downloads about 200 MB of weights from Hugging Face. Its inference stays on the device. Confirm dates marked **Needs review** against the original photo. Set the numeric date order in **Settings**.

Choose **Read labels again** on a saved scan to reread both labels with the current crop. When an unlabelled stamp has two clearly stacked dates, Froststock suggests the later date and shows both dates for review. Tap one of the dates and confirm it against the pack.

For the standalone preview, export a backup in its Settings, then import it in Froststock under One Workspace. It writes each imported pair to My Drive. Managed app updates keep the Drive files. Use Workspace Backups to protect My Drive and the app’s private sign-in data.
