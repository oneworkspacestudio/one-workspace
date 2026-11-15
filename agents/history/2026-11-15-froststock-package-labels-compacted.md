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
