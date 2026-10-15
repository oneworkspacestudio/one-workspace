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
