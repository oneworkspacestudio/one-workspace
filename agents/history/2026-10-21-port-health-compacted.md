# Installer ports and MyAccounts health

User requested correction of installation port conflicts and false unhealthy
MyAccounts in Monitor. Existing owner Workspace57 was stopped; its manager was
started using the existing config, and all24 services became healthy. Numeric
127.0.0.1 Host gets403 from Accounts; localhost gets200. No owner data was a fixture.

Fresh Install now tries complete free port groups before identities or app config
are created; all selected app ports, auxiliary TCP and Meet UDP are checked. New
local origin/issuer follow the chosen base; public HTTPS origins are preserved.
Existing runtime/update port assignments are retained. Errors preserve the real
bind cause instead of labelling every failure as already in use.

Monitor3 uses canonical localhost authority only for managed exact registry URLs,
retaining pinned loopback dialling and arbitrary-target SSRF/redirect restrictions.
Suite/setupui tests, Monitor all Go tests/vet and unchanged AUI frontend tests pass.
Native occupied-port fresh setup and deployment acceptance remain to be recorded.

Build artifacts I:/Projects/tooling/workspace-ports-20261021; isolated NTFS tests
C:/Projects/.cache/workspace-ports-20261021. Use Go1.27.2 with GOTMPDIR on this C
test root, shared GOCACHE/GOMODCACHE on I. Git needs per-process global config NUL
and this exact safe.directory to avoid unavailable drive entries. Source app's
pre-existing untracked internal/app/linux.go is unrelated and preserved.
Tags/assets remain immutable. Actions remain disabled. Preserve owner identities,
data, installed Chat9 and the newest verified recovery copy.
