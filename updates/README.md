# Public update channel

Workspace installers and maintenance windows without a GitHub token read
`releases-v1.json` here over HTTPS. One request covers the Workspace suite and
all managed app repositories. Public release assets use their public download
URLs, with the same repository, tag, size and SHA256 verification as API assets.
The 60-request anonymous GitHub API allowance is not used by this path.

This is the publisher's promoted update channel, not a live GitHub API mirror.
Its `publishedAt` date is shown in the update window. Releases become available
to these clients after promotion. A missing, invalid or unavailable channel
stops the check with an error; it does not claim everything is up to date.
Explicitly authenticated/private or custom-repository installations retain the
API path and display rate-limit retry guidance if GitHub refuses a request.

## Required release promotion

After **every** public Workspace or managed app release:

1. Verify all public release assets, sizes and GitHub SHA256 digests.
2. Run `python scripts/publish-update-index.py` with publisher access. It uses
   GH_TOKEN/GITHUB_TOKEN or existing Git Credential Manager credentials in memory.
   No token or account information is written to the output.
3. Review the JSON changes, including the new tag and every expected asset.
4. Commit and push the generated file to main. Do not move the source release tag.
5. Verify the raw public URL contains the promoted tag, then perform an
   unauthenticated installer/update check and verify its selected assets.

Do not announce delivery until this promotion is complete. Run publishing and
builds locally; GitHub Actions must remain disabled. Draft releases are excluded.
Do not replace a published release's files; publish and promote a new version.

`withheld-releases.json` excludes explicitly withheld builds during every promotion. Keep those immutable historical artifacts out of the default channel; never remove a withholding record merely because a newer build is published.
