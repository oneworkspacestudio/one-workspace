# Froststock

Choose **Froststock** under **General** when installing or updating One Workspace. Open it from the Workspace app picker and sign in with the Workspace owner account.

Take or upload one photo of the food and a second photo of its date label. The scanner returns to the first step as soon as the pair is queued. Keep scanning while your browser saves the pairs to My Drive and extracts the text. A temporarily disconnected browser keeps unsent pairs and retries when the app is open again.

The dashboard shows items due within seven days, items past their printed date, and items that need a closer look. Open an item to verify the photos and edit the label details. Mark used items to move them to History.

Your photos and inventory records live in **My Drive → Froststock**. A second device signed in to the same Workspace sees the same inventory. Each item folder holds `item.jpg`, `expiry.jpg`, and `item.json`. Keep those files together. My Drive access controls who can view the underlying Drive files. Froststock itself permits the Workspace owner.

OCR runs in the browser using bundled English Tesseract. The optional SmolVLM 256M model uses WebGPU and downloads about 200 MB of model weights from Hugging Face. Inference runs on the device. Check suggested dates against the original label. Change the numeric date order in **Settings**.

For the standalone preview, export a backup in its Settings, then import it in Froststock under One Workspace. It writes each imported pair to My Drive. Managed app updates keep the Drive files. Use Workspace Backups to protect My Drive and the app’s private sign-in data.
