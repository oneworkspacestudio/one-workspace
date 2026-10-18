# Four-app delivery

Recon: organization admin access confirmed; actual suite main, MyAccounts main,
Notes main, Calendar main, Drive master. No equivalent repositories exist. Reuse
platform/OIDC/brand contracts; standalone local auth remains explicit. See four
ADRs. New defaults use updater policy from the previous repair; preserve that fix.

Order: shared apphost/auth/storage + security tests; Tasks API/UI + recurrence,
permissions/reminders; PDF jobs/viewer/Drive + destructive-content checks; Monitor
registry/checks/incidents + SSRF; Search actual connectors + global header; full
isolated suite install/upgrade/logout/backup acceptance; local release builds,
GitHub digest verification, source/readme screenshots and truthful preview gates.

Source/build/cache stay on I:. Use four named C:/Projects/.cache NTFS app fixture
roots and one named suite root. Track and stop test PIDs; remove only their exact
non-reparse paths. Never use the live owner MyAccounts/Drive for tests. Preserve
live data and newest rollback. No Actions workflows, hosted CI, Docker or WSL.
