# My Git and Hosting delivery repair

Task: accept repository name, browser URL and .git URL in Hosting, then release
and install. Hosting2 / Workspace110 add a capability-authenticated loopback
read-only broker. It maps issuer/subject through read-only Git OAuth SQLite
records and uses Gitea Sudo for per-user authorization on every source operation.
No privileged Git credential leaves the manager. No authority reaches build jobs.
Ambiguous names ask for owner/repo; missing linked Git identities need one visit
to My Git. Existing saved GitHub owner/repo sources remain GitHub. New GitHub
inputs use explicit URLs. External Gitea/Forgejo HTTPS uses a repository-bound
read token and connection-time SSRF/redirect checks.

Hosting2's Source step validates early and detects default branch/framework.
Plain files need no package.json; npm suggestions include build dependencies,
using ci when lockfile exists. Immutable source archive and branch membership
checks cover all providers. Existing signed GitHub webhooks retained, native
Gitea hook signature supported on a narrow public POST route.

Verification: all Workspace Go packages passed; app Go suite, production UI
build and Chrome desktop/mobile creation/error tests passed. Actual native
Gitea fixture resolved all3 forms, read private sources and denied another user's
private repository/disabled users. Actual Hosting/Caddy static publication and
npm/Node build from that private repo passed with synthetic isolated NTFS roots.
The ordinary builder replaces all existing Hosting payloads while preserving
its namespace; app update requires110/native-hosting-git-v1. Git is optional.
No VM, physical phone, reboot or SCM claim. Public release/install still pending.

This task uses I:/Projects/workspace-hosting-git-20261102, cloned from d5c4bce;
origin advanced to b7be993 for concurrent Social11 delivery; incorporate it before
pushing. Preserve owner's completed Social11 update and all unrelated versions.
Verified seed103 exists in prior task's range-cache/seed103-verified.zip; its
cache/<sha>-workspace-bundle.zip is a truncated348MB file and MUST NOT be used.
Owner current install is C:/tmp/workspace; do not use its private state as a test.
