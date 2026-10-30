# Publish a release locally

Every new release in `oneworkspacestudio/one-workspace` includes a signed Android
APK, including desktop Workspace releases and Android-only releases. Users install
one Android launcher and open the web versions of all their apps from it.

GitHub Actions remains disabled. Publishing runs on the maintainer's computer.
Do not create a source-only release through GitHub's web interface or raw `gh`
commands. The supported entry point is `scripts/publish-release.ps1`.

Before adding or maintaining an app, complete [APP-LIFECYCLE.md](APP-LIFECYCLE.md),
including the managed registry and fresh-install/upgrade acceptance.

## Before publishing

1. Finish the intended changes. For Android changes, increase `versionCode` and
   `versionName` in `android/app/build.gradle`. Update the README download links
   and release notes. Keep runtime/phone acceptance claims accurate.
2. Commit and push the complete source. The publisher requires a clean worktree
   and a source commit accessible in the public repository.
3. Use JDK17+, Android SDK35 and Gradle8.11.1. The wrapper is included. Set
   `JAVA_HOME` and pass the SDK location. `-Cache` and `-Gradle` can reuse an
   existing Gradle cache/distribution. Keep large artifacts on the selected
   build drive; never stage private data in a public asset directory.
4. Supply the existing permanent signing key through
   `WORKSPACE_ANDROID_KEYSTORE`, `WORKSPACE_ANDROID_STORE_PASSWORD`,
   `WORKSPACE_ANDROID_KEY_ALIAS`, `WORKSPACE_ANDROID_KEY_PASSWORD` in the local
   process environment. Never commit key files/passwords or use debug keys.
   `android/release-policy.json` pins the public signing certificate fingerprint.
5. Authenticate `gh` locally. Use an external release notes file and a dedicated
   directory containing only intended public assets. Desktop tags also require
   the built installer, manager, catalog and application bundle plus their
   existing source, licence and acceptance materials.

## Publish

For an Android-only release:

```powershell
./scripts/publish-release.ps1 -Tag android-v0.1.0-preview.4 `
  -Title 'One Workspace Android preview 4' `
  -AssetsDirectory 'I:/Projects/tooling/android4-release' `
  -NotesFile 'I:/Projects/tooling/android4-notes.md' `
  -Sdk "$env:LOCALAPPDATA/Android/Sdk"
```

For a desktop Workspace release, pass its `v...` tag, matching catalog version
and prepared desktop asset directory to the same command. Do not omit Android.
The APK stays a separate release download; it does not require adding another
desktop application bundle download to an Android-only update.

The command builds the signed APK, verifies its certificate/package/version and
embedded commit/tree provenance against the committed source, and packages its
source/notices. It creates or resumes a matching draft, uploads missing assets,
verifies all sizes/SHA256 digests, then publishes. Differing draft assets stop
publication for review; published releases are never overwritten. It verifies an
anonymous APK download, regenerates the public index, commits only that index and
pushes to main. It then checks the exact public index URL used by installers.

If publication succeeds but channel push or CDN propagation fails, the error
identifies the remaining step. Inspect the release, run
`scripts/publish-update-index.py` if needed, commit/push its generated index and
verify the exact public raw-main URL before announcing delivery. Do not rerun a
published tag or replace its files. A matching unfinished draft can be resumed.

Complete the anonymous installer check from [updates/README.md](../updates/README.md)
with temporary downloads on the build drive. Record acceptance, cleanup and
delivery in the compact history. Do not claim a phone or live-install test that
was not performed.

This local guard and repository policy apply to One Workspace releases. GitHub
itself permits maintainers to bypass a local script, so use this entry point for
every release. Server app repositories retain their own release workflows.

## Windows publication reliability

The publisher decodes GitHub/Git output as UTF-8 and batches asset arguments below
Windows UTF-16 command-line limits. `python scripts/test-publisher.py` covers
large catalogs, Unicode paths and empty/overlong input without network writes.

For very large assets, use current official GitHub CLI tooling and NTFS staging
if another filesystem stalls reads. A failed draft upload may be resumed with
the same source and exact bytes. Verify retained assets, retry only missing files,
then rerun the mandatory publisher to verify the full set and promote. Never
clobber published assets or announce an unverified channel.
