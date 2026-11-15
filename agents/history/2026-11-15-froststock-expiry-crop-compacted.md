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
