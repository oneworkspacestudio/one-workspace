# Shield private-history update discovery, 30 October 2026

User confirmed workspace-2 still loads old index-CzNDSlaq.js and its Settings
.map crash; My Shield4's fixed build is index-B3l6dnYm.js. Shield4 was published
in Workspace97 at2026-10-28T21:54:32Z; public latest107 catalog also pins Shield4.
The public channel intentionally has an empty my-shield release history.

Found discovery defect in CheckUpdates: different app sources with missing release
timestamps retain the installed contract. Private histories are omitted, so both
timestamps are zero and Shield3 is retained even with verified suite Shield4.
Fixed by allowing a strictly newer canonical same-repository app version from
the verified suite catalog when the public index omits history. Older pins still
retain installed app. No private release metadata is exposed. Regression exercises
public discovery and app-only selection (one payload, retained manager), plus
rejection of a lower suite app version. Targeted regression passed. Whole
internal/distribution package check passed (13.620s).

Source changes: internal/distribution/updates.go, private_suite_test.go.
This correction is NOT yet in a published updater or installed on either PC.
Needs complete test check, commit/push and next unclaimed suite release via mandatory
signed-Android publisher; preserve concurrent latest main/app pins. Do not overwrite
97/107. Need supported remote protected app-only Shield4 update after corrected
updater installed; old updater can falsely say up-to-date. The manifest401 remains
separate: browser link lacks use-credentials, Workspace protects app paths; do not
open app APIs publicly to silence it. Existing primary Shield4 was verified earlier.
Native Chrome automation now observes via @oai/sky but input was stopped for inability
to determine current browser URL; no remote update submitted.
