# Component provenance

For the complete organization build, see [Organization apps](ORG-APPS.md). Its
catalog.json and licenses/components/organization-components.json record the
actual app versions, organization release URLs, sizes and checksums. All app
catalog sources belong to oneworkspacestudio. The seed notes below describe the
original development package and do not identify the current packaged binaries.

## Original development seed (historical)

The release's `catalog.json` records exact SHA256 digests and byte sizes for
every extracted file. The installer verifies them before first launch, and the
supervisor verifies each process executable again and holds a deny-write/delete
handle during its lifetime. The manager and its trusted state live on protected
NTFS storage; large component executables may remain on the program volume.

| Component | Version/source | Provenance |
| --- | --- | --- |
| Workspace manager and setup | This repository, 0.1.0-preview.1 | Go 1.27.2, Windows amd64, CGO disabled |
| MyAccounts | [preview-20261014.2](https://github.com/oneworkspacestudio/myaccounts/releases/tag/preview-20261014.2) | Published executable, SHA256 `8626e7b6ec628d46d57a35bd6adee8f468256c890e0a4e0e79b06fbcbd49ae33` |
| FileBrowser Drive + Docs | Local `photos-preview` development build | SHA256 `304c32928cdd422f79083038c4189b66de6e7970b7790ecae29c6c43dbdcda23` |
| Photos companion | Local development build from the custom FileBrowser tree | SHA256 `cee5b11051327c3a84b9972e4d28b242a42ef62cacb8c79f5a67197da49f52bb` |
| cloudflared | Official Cloudflare 2026.10.0 Windows amd64 | Exact digest in the release catalog; optional at runtime |

The latest Photos source was uncommitted in the existing FileBrowser worktree.
This repository does not claim that it has a corresponding published upstream
tag, or publish changes back into that worktree. The component seed release
includes its source snapshot (with editor source) and links to separate engine
source archives. That snapshot
captures the available development tree; byte-for-byte reproducibility of the
pre-existing development executable has not been established.

The editor derives from `ranuts/document` at
`1330b3c525ca4d7211f99592f5e94df749e5202a`, including downstream embedding and
collaboration patches. Its upstream AGPL license/notice and the FileBrowser
third-party notices are extracted into the program folder. Engine sources and
the pre-existing converter source snapshot are linked from the component seed.
The converter upstream did not pin its compiler/source for exact reproduction;
the source snapshot is not asserted to reproduce its WebAssembly byte for byte.

Workspace orchestration source is Apache-2.0. Separately licensed components
retain their own terms; see `LICENSE`, `NOTICE`, extracted `licenses/`, and the
component source archives. Do not treat the orchestration license as relicensing
the bundled editor or its combined FileBrowser distribution.
