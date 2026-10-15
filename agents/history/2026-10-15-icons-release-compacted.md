# Application icon refresh and local Windows release

User authorized publishing all app fixes and new releases to oneworkspacestudio,
then requested the supplied eleven SVG app icons and README logos first.
This source is prepared for `v0.3.0-preview.20` in `oneworkspacestudio/one-workspace`.
Builds run locally on Windows; organization Actions are disabled, including the
previously enabled Notes/MyAccounts repositories. No Actions job is used.

The supplied SVG paths/gradients remain unchanged; PNG/ICO variants and Windows
resources are generated from them. Web branding, favicons, suite launchers and
repository README logos are updated. The original pack lacked a hub icon; the
owner subsequently supplied one-workspace.svg. It now replaces the hub artwork
in the app, setup, Windows resources and README. Build scripts preserve original
license notices and pinned third-party engines. Do not copy live/private data
into any source or release artifact.

Prior app-flow-fixes handoff records implemented fixes and sampled acceptance.
The icon refresh is a build/publication task; it is not a new full browser or
clean-host acceptance run. Mail Safe Browsing review and provider logout browser
block remain unresolved; do not claim delivery or provider sign-out accepted.
Notes/Mail/Calendar integration feature gaps remain explicit in release notes.

Release preparation/evidence outside Git: projects tooling workspace-release-20261015.
The suite handoff records the final publication receipt. Verify GitHub release
state/assets and source tag before claiming publication. Never rewrite published
tags/assets or use GitHub Actions for this release.

## Distribution decision and current Music limitation

Defender quarantined the new Music.exe during staging (Trojan:Win32/FakeFolder!pz).
The exact SHA256 is recorded in docs/RELEASE-STATUS.md and the lock file. With
explicit owner approval, Microsoft accepted the binary for analysis; final
determination is pending, not a confirmed incorrect detection. Private submission
receipt remains outside Git. No Defender exclusion, disabling or quarantine
restoration was used. Requested manual Defender scans failed to start.

Owner then requested the complete bundle with issues treated as current limits.
Build it with the seven updated component releases and the previous published
Music workspace-20261014 binary, whose SHA256 is pinned. Do not include the newly
flagged executable. Music's new internal artwork/profile/cache changes therefore
remain unreleased; suite launcher/switcher icons are updated. build-org rejects
the withheld hash and downloads only pinned, published organization releases.
Drive/Photos now use their own new organization release, replacing seed reuse.

No new full E2E/clean-host acceptance was requested for this branding release.
Current integration and browser security limits are in docs/RELEASE-STATUS.md;
include it in the bundle and release notes. The previous complete release is
untouched. Preserve the active installation; no new test installation is needed.

## Publication completed

Seven component previews are published in oneworkspacestudio: Mail preview16,
Photo Editor studio-v0.1.1, Git workspace-20261015, Calendar/Contacts preview2,
MyAccounts preview-20261015.1, Drive/Docs/Photos drive-v0.5.0-preview.2, and Notes
preview3. Their GitHub asset SHA256 digests and byte sizes matched the staged
source/license/build-info/checksum files. Music source was pushed, with no new
Music binary release. All nine repository Actions settings are disabled.

During packaging, an independently authorized Notes editor fix advanced Notes
main to4c6e422 and published preview4. Its new complete Workspace preview21 retains
all supplied branding and the previous Music pin. Stopped this task's preview20
publisher before publication; its tag/assets remain in a clearly superseded,
unpublished draft. Do not publish/reuse20 or rewrite its existing tag/assets.

Final complete release:
https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.21
Published2026-10-15T21:28:47Z. All11 release asset sizes and GitHub SHA256 digests
were independently checked against its checksum file. Local package checks confirm
all11 org app sections, Notes preview4 matches the published Notes executable,
the earlier Music hash is used, and the flagged Music build is excluded. The
source archive's Workspace SVG matches the separately supplied original artwork,
Windows icon resources match, and README references the new logo. Preview21 is
the newest published complete bundle; online installation therefore selects it.

Detailed receipts remain outside Git in workspace-release-20261015, including
preview21-local-verification.json and suite-publication-receipt.json. New build
checks do not replace the earlier sampled functional acceptance or clear current
limits. Microsoft's Music determination is still pending. No active installation
was replaced by publication. Named Studio/Music C: build caches were removed after
source/output preservation, canonical-path/reparse and process-reference checks.

## README logo visibility correction

GitHub rendered the supplied SVGs, but inline40/42px title icons were easy to
miss beside older screenshots. README now leads with a112px labelled logo.
Workspace also has a current-logo overview and explicit organization repository
links for every app. Documentation-only follow-up; published binaries/tags and
current preview limitations are unchanged. No Actions or runtime tests needed.

The documentation correction is pushed to the default branches of all nine
organization app repositories. Live GitHub browser checks confirm every header
logo loads (natural width256, displayed112px); Workspace's overview loads at its
original1440px width and its app table includes explicit repository links.
Public-page proof and per-repository push receipts remain outside Git in the
release tooling folder. These are README image checks, not app runtime tests.

Owner requested retaining the original illustrated showcase too. Its image file
was never deleted; restore it as the README hero and keep the new logo overview
under Meet your apps. Future branding edits must retain both images.
