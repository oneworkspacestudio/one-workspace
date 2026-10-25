# My AI

My AI joins One Workspace as an original native application for AI chat and approved code projects. It uses MyAccounts, a dedicated private My Drive source, the shared Apps menu, health checks and normal updater/backups. Runtime has one main EXE with embedded browser assets. Optional Codex and Ollama remain separate processes.

The public app repository contains setup, architecture/licence audit and the detailed support/acceptance matrix: https://github.com/oneworkspacestudio/my-ai . This is an evaluation preview; cloud provider edit/test acceptance needs explicitly configured keys. No Open WebUI code is redistributed. No existing Codex login/config is reused.

Preview4 adds local-first expandable cloud lists, persistent shared server downloads with resumable progress, real time/calculation/public search/page tools and private scoped reads. The bundled public read-only stdio MCP server is configured only in My AI-owned Codex environments. Managed preview4 requires Workspace87. Current86 installations select Workspace and My AI together; other app versions remain unchanged. Restricted Code profiles require explicit supported Windows sandbox setup.

Preview5 makes projects name-only. Workspace88 creates a private Projects area in
My Drive and a private repository in Git, mapped to the same MyAccounts subject.
My AI edits selected project files using the chosen provider's native tools and
automatically commits/pushes completed changes. This basic flow needs no Codex or
host-folder approval. Ordinary members can create and use only their own projects.
No shell/test runner is supplied; cloud-provider turns remain a separate test gate.
Native checks cover real Git/OIDC/Drive/local-LLM editing, retry, safe incoming
changes/conflict refusal, credential exclusions and temporary connection recovery.
Managed preview5 requires88; select Workspace and My AI together on87.
