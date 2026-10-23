# Updates from an idle maintenance window

User reports HTTP409 applying selected updates and requests installer fix, source
push, release and local installation. Owner I:/Workspace is preserved; unrelated
I:/.workspace untouched. Read AGENTS, RTK and prior installer-handoff history.
No agents or GitHub Actions. Public70 immutable; next release71.

Reproduced on installed70: release check 65+minutes old, HTTP409 "Release check
expired; check again". UI request error was overwritten by the next1.5s poll.

New server path refreshes expired checks before install or reinstall. Continues
only if selected resulting catalog and exact download metadata remain identical;
installer repair compares validated exact installer metadata. Changed files or
installation require review. Network errors retain the old check and stop safely.
HTTP errors stay visible/focused until the next action. Progress scrolls into view.
Freshness also rejects missing/future dates. Existing capabilities, origin/host,
operation locking, checksum validation and backed-up transactions remain active.

Artifacts I:/Projects/tooling/update-expiry-20261023. Only disposable NTFS root:
C:/Projects/.cache/update-expiry-20261023. Track fixture processes and stop before
bounded cleanup. All tests/vet initially pass; final native/public/local acceptance
will be recorded before completion. No physical reboot claim.


## Delivery complete

Public v0.3.0-preview.71 source8eac4d9e1f6b6bf615c093be3912eef7f5021448
pushed to main and immutable tag. All25 public release asset sizes/SHA256 digests
verified; public manager downloaded and verified. Full installer also copied and
verified in Downloads. All25 app contracts and134 unchanged payloads retained
from70;136 payloads verified. No Actions.

All Go packages, vet and final updater JS syntax pass. Browser synthetic HTTP409
error remains visible/focused after background status polls. Native Windows test
uses a real65-minute-old plan, public release refresh and actual verified Drive4
download/install through POST/api/apply, returning202. Verified backup, owner
and credentials, user document and unselected apps preserved; all4 selected
process services restarted. Original owner login verified before update. Native
test196.69seconds. First run hit our UI mock on the fixture Tasks port; stopped
mock, started healthy fixture and reran successfully. This was a harness conflict.
No full71 setup repair fixture or physical reboot acceptance is claimed.

After publication, installed71 locally through the normal backed-up setup update.
Then the new maintenance HTTP flow installed Drive4, Mail17, Slicer4 and Clips6.
Both completed with verified recovery copies. Local and public health report
all24 services running. Owner/credentials, locations, connection settings,
selection and unselected app contracts preserved. Startup registration verified.
Updated browser shows completed operation; screenshot in docs/screenshots.

The owner started the four-app operation in the newly opened browser while the
delivery script was preparing the same request. The duplicate API request was
correctly rejected as busy; observed the existing operation through completion
instead of issuing another request. This was separate from the original expired
check409. A local Python locale decoding error was also corrected before mutation.
The installed helper SHA256 matches the released71 manager; page includes the
new error handling and shows Update complete with an empty error field.

All test processes stopped; removed only named C cache after process/reparse
checks (1,671,704,240bytes). No fixture processes remain. Live owner runtime and
normal installed maintenance helper intentionally remain running. Evidence:
docs/evidence/update-expiry-20261023.json. Later evidence commit must not move
the immutable source tag or rebuild the release.
