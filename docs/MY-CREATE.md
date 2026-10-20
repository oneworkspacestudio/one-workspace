# My Create

My Create is a local design editor for text, shapes, images and multiple pages.
It imports user-provided SVG, PowerPoint, PDF, images and supported Photoshop
layers, preserves source files, and reports conversion losses before editing.

The application uses the shared MyAccounts identity, Workspace navigation and
normal update/backup contracts. Default route: `/create`; default port: 9017
(base port + 27). Private data is stored in the dedicated `create` data folder.
Drive copies use the current user's own Drive session and permissions.

Source: [oneworkspacestudio/my-create](https://github.com/oneworkspacestudio/my-create).
Template packs: [my-create-templates](https://github.com/oneworkspacestudio/my-create-templates).

Imported PDF objects are reconstructed drawing operations, not original design
layers. Pixel images contain no native editable text. Optional Windows OCR
creates review-required inferred text and does not erase original pixels.
PPTX preview is reconstructed; it is not an independent PowerPoint renderer.
Canva links/private catalogue exports are not supported. Use files exported by
their authorized owner. No mandatory AI download or paid cloud service.

My Create 0.1.0-preview.1 integrates with Workspace57. The synthetic upgrade
verified retained MyAccounts sign-in, project saving, same-subject My Drive copies
and served icon hashes. See docs/evidence/my-create-native-20261021.json.
The first preview is unsigned; advanced AI reconstruction, a direct Canva
companion and full complex-effect fidelity are not included.
