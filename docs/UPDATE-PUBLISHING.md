# Publishing independent app updates

Publish **workspace-update.json** alongside a new Windows app release. Include it
in the draft before publication; do not rewrite previously published tags/assets.
The updater gets release timestamps, URLs, archive hashes and byte sizes from
GitHub. The contract declares the app-owned files and managed integration that
remain compatible with existing storage. Repository and tag must match exactly.

The common format supports every current organization app, including coupled
Drive/Docs/Photos and Calendar/Contacts/helper packages. It only permits a group to
replace its own known package paths and notices. Workspace itself continues to
use the complete bundle catalog, its manager asset and publisher-verified bundle.

Use the current compatible Workspace catalog as the base. Build the app and put
the exact binaries and licenses into the release archive. Make a local package map
with every package owned by that app group. For example, Notes uses:

```json
{
  "notes": {"path":"dist/MyNotes.exe","asset":"notes-windows-amd64.zip","member":"MyNotes.exe"},
  "notes-LICENSE": {"path":"dist/LICENSE","asset":"notes-windows-amd64.zip","member":"LICENSE"},
  "notes-NOTICE": {"path":"dist/NOTICE","asset":"notes-windows-amd64.zip","member":"NOTICE"},
  "notes-THIRD-PARTY-LICENSES.txt": {"path":"dist/THIRD-PARTY-LICENSES.txt","asset":"notes-windows-amd64.zip","member":"THIRD-PARTY-LICENSES.txt"}
}
```

Generate the contract in this repository:

```powershell
go run ./cmd/update-manifest --catalog catalog.json --group notes --version 0.1.0-preview.5 --repo oneworkspacestudio/notes --tag v0.1.0-preview.5 --files notes-update-files.json --out workspace-update.json
```

The tool computes SHA256 and size from each local file and outputs no local paths
or credentials. The package map's asset name must exactly match the uploaded
asset. Its member must exactly match the path within the ZIP. Omit member for a
direct uploaded file. Supply all owned executables, helpers and notices together;
missing files or cross-app replacements are rejected. The updater verifies the
archive digest, then each individual member’s bytes, size and uniqueness.

Generated schema:

```json
{
  "schema": 1,
  "integration": 2,
  "repository": "oneworkspacestudio/notes",
  "tag": "v0.1.0-preview.5",
  "apps": ["catalog App objects retaining the managed launch contract"],
  "packages": {"package-id": {"file":"installed/path","sha256":"64 hex characters","size":123,"asset":"uploaded-file.zip","member":"path/in/archive"}}
}
```

The apps array contains real catalog App objects; the example text is explanatory.
The generator preserves the base launch/storage contract, changes app version and
source URLs, and validates the resulting catalog. A publisher must test the new
release against existing databases and these arguments before declaring integration2.
Existing OIDC callbacks, data paths and private Mail backend must remain compatible.
Breaking migrations need a separately reviewed migration flow, not this marker.

The organization release remains an unsigned preview until the separate platform
acceptance gates are met. This metadata is a compatibility declaration and uses
GitHub release integrity; it is not an executable code-signing certificate.
