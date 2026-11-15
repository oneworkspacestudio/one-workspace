# Froststock expiry crops

Froststock preview.4 adds the optional **Crop expiry photo after capture**
step. Existing items can use **Crop & read again**, which queues a date-only
scan and preserves the saved name and original photos. The bounded crop
coordinates live in the existing My Drive record and survive export and
restore. Records from earlier versions are still readable.

Difficult dates fall back to bundled PP-OCRv6 tiny models with ONNX Runtime
WASM in a background worker. The existing Tesseract path remains the fast
first pass, and numeric date interpretation still uses the configured order.
The model returns text for the existing date parser. Close expiry candidates
are flagged for review. New bundled sources and licences are in
`internal/froststock/web/assets/paddle`.

In synthetic browser checks, the Paddle models read the faint stamp, kept
packing and expiry text distinct, returned a date without inventing a year,
and returned no text for a blank label. SmolVLM 256M matched both complete
dates on the same fixtures. The 500M model chose the packing date in one
case. Both produced a stray year on a blank image that the date parser
discarded. Retained 256M as the optional WebGPU fallback and kept its date
suggestions subject to review. The 390-pixel browser crop test passed
touch, cancel, saved-photo retry, and original-photo preservation.

Workspace preview.134 builds from the published preview.133 seed and uses
the standard signed companion release process.

The first managed package exposed dynamic JavaScript in the OpenCV
bindings. Replaced it with closures in four helpers. A browser fixture
with the same CSP as the managed app passed before the installer rebuild.

The shipping installer verified all 47 apps and 939 package hashes on the
preview.133 baseline. Fresh managed setup passed touch crop at a mobile
viewport and ran PaddleOCR through the authenticated Froststock proxy.
App-only update, saved-photo date-only recognition, and rollback passed
on an isolated preview.133 fixture. Title and original Drive photo bytes
were retained. Go tests and the recognition and publisher checks passed.
The isolated services were stopped and their private fixture roots removed.
No physical phone or owner inventory was used.

Workspace preview.134 shipped from `02658a6` with 960 verified release assets
and the signed Android companion. Promotion `f2b9e3a` became visible on the
normal public index after cache propagation. The anonymous preview.133
installer offered Froststock preview.4 and selected all five expected
package hashes. Anonymous catalog and package downloads matched.
The owner installation was outside this update.
