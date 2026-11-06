# My Podcasts integration

App source is a separate `my-podcasts` repository. This integration worktree adds
an optional Creative & Development app, ID `podcasts`, `/podcasts`,
`MyPodcasts.exe`, public service port 9044 and manager-relative offset +54.
MyAccounts callback is `/podcasts/auth/callback`. There are no sidecars or
browser/FFmpeg/Node runtime dependencies. All mutable files and encrypted
provider settings belong under the app's dedicated private NTFS root.

Managed inventory, complete builder, strict package namespace, shared picker,
content-addressed branding, logout cookies, update discovery and independent
recovery storage are registered. Installer integration 14/minimum manager 122
prevents old managers silently dropping the app contract. Complete-catalog and
publisher checks reject omission. The builder pins the verified public preview.2 archive.

The normal builder supports `-MyPodcastsCandidate <archive>` for isolated local
acceptance. Its manifest remains candidate=true and publication rejects it.
Public releases must use the verified component pin and mandatory signed-Android
suite publisher. Keep Actions disabled and verify the public channel and each
authorized installed PC. Never replace public assets.

Choose OpenAI, DeepSeek or the existing AI Service/Ollama as writer. DeepSeek or
OpenAI can supply web research; local automatic research uses selected RSS/direct
X. AI Service requires an owner-issued local app credential for `podcasts`, its
loopback address and an active model selected in the service. No model/grant is
automatically activated. OpenAI speech is separate; reviewed scripts can wait
for voices. Direct X timelines require an optional developer token. Private
APIs remain protected by MyAccounts, native owner authorization and CSRF.
Background schedules survive tab closure and sign-out by design; Cancel stops
authorized production. Owner deployments and paid-provider quality are separate
acceptance stages, not inferred from local tests.
