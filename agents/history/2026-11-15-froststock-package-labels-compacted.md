# Froststock package label recognition

Reproduced failed titles and a wrong expiry year on two angled packages
supplied privately by the owner. Browser PaddleOCR could already read both
food labels and both date lines on the clearer stamp. The previous app used
that model only for dates and could let a conflicting OCR pass replace the
clear two-date result. The running public installation was still preview.128
with Froststock preview.1 when inspected. No owner update was requested.

Froststock preview.5 uses the existing PaddleOCR model for both fields.
Projected text bounding boxes account for slanted packaging. Food terms and
adjacent title lines keep the product cut and description together. A bare
numeric stamp can recover a missing divider when the year is printed. Two
high-confidence unlabelled rows may suggest the later date if the rows are
adjacent, aligned and have a plausible interval. Suggestions remain under
review. Both printed values are saved as options for the owner to select.

Browser emulation on the exact attached photos recovered both complete
names and both later dates from the whole images. The two rows are retained
as candidates and generated dates stay unconfirmed. Cropping to the stamp
and to the later line also worked. Private evidence and photos stay outside
Git and the release. Workspace preview.135 uses the verified preview.134
seed and the standard signed release process.

The isolated DeepSeek-OCR 3B comparison used the same unmodified photos.
Ollama's local CUDA run failed with a PTX compilation error. The CPU run
identified both products and read both dates on the sharper photo. It gave
no usable text for the blurrier stamp, where local PaddleOCR and the date
parser recovered both rows. Warm CPU calls took about 8 to 10 seconds per
image. The test used the official 6.7 GB model and processed the photos
only over loopback. The test runner and its model cache were removed. The
separate DeepSeek Flash feature on main uses a different model.

The signed preview.135 release shipped from `8c31959` and verified 960
public release assets. Channel promotion `de3a3a1` propagated to the normal
public URL. The anonymous preview.134 installer selected Froststock
preview.5 and the expected five packages. Downloaded packages and catalog
matched their SHA-256 digests. The full installer verified 47 apps and 939
package hashes against the preview.134 baseline. Fresh setup and the
app-only update recovered the labels in an isolated mobile-browser test.
Date selection, saved-photo reread, and rollback preserved the Drive
photos. The test services and private fixtures were removed. The live
owner installation still reported preview.128 and Froststock preview.1.
