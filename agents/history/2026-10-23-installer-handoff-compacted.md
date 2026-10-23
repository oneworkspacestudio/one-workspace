# Verified installer handoff from Updates

User requests fixing the newer-installer error, automatic installer download and
reinstall/restart from the update page, then source push, release and local bootstrap
so they can run reinstall. Preserve I:/Workspace owner installation and recovery;
the much older I:/.workspace installation is unrelated. No agents or Actions.

Read latest repository histories and fast-forwarded clean I:/one-workspace to
c72168f. Work branch codex/installer-handoff. Public69 integration3 catalog is
rejected by installed66/67 integration2 controller before update choices appear.
Current owner runtime is running; config autostart=true and HKCU Run targets its
protected manager. No StartupApproved disable entry was found. Actual reboot
failure not yet diagnosed, and no real reboot acceptance is claimed.

New stable installer metadata survives future catalog integration/schema guards.
Only same-release organization assets with GitHub size/SHA256 metadata are used.
Update page explicitly backs up/reinstalls via verified full setup --repair DATA
--offline, using its exact embedded snapshot and existing journalled repair.
Separate app versions/accounts/selection/paths/files retained. Child progress is
relayed; result checks every expected process and connector, with app retry link.
No account reset path from maintenance. Future controller can report completion
using the stable settings header without interpreting unknown catalog semantics.

Targeted Go tests passed. Full tests/vet, native exact-installer damaged-program
repair and identity/file checks, publication and owner bootstrap still pending.
Artifacts I:/Projects/tooling/installer-handoff-20261023; only named C NTFS fixture
C:/Projects/.cache/installer-handoff-20261023 is disposable after process/reparse
checks. Go1.27.2 ZIP verified against official digest. Public69 bundle digest
verified. Keep all25 app contracts and payload bytes from69 in the new release.
