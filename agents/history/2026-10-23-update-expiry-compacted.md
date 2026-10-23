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
