# October 16 release status

The supplied eleven app SVGs are applied to app headers, favicons, Windows icons
and repository README logos. Component executables were built locally; GitHub
Actions is disabled for all nine organization repositories.

My Tube is now included as a twelfth app entry, pinned to its published
0.1.0-preview.1 release. Its separately supplied logo is used in the README,
launcher and switcher. This suite update is built locally; the previously
verified My Tube component binary is reused without modification. Its optional
FFmpeg/FFprobe helpers remain separate and are not shipped in the suite.

## Music uses the previous published version

Defender detected `Trojan:Win32/FakeFolder!pz` while staging the newly built
Music executable. This is not confirmed to be an incorrect detection. The exact
77,280,038-byte file was submitted to Microsoft Security Intelligence with the
owner's approval. Microsoft accepted the submission on October 15, 2026; its
determination is pending.

SHA256: `5f0d121cd9733775b327ac3b39b46b32cf5d9cbdef1b014880a89d6421239837`

The new Music binary is not released. The complete Workspace bundle uses the
previous published Music version (`workspace-20261014`), with its exact source,
SHA256 and size recorded in package provenance. Its internal icon/profile/cache
updates from October 15 are therefore not included. Workspace's Music launcher
and shared app switcher use the new supplied logo. The organization build script
rejects the withheld binary's exact hash.
Do not disable Defender, add exclusions or rename/repackage the executable to
avoid the detection.

After Microsoft's determination, review its findings, update security
intelligence normally and validate the exact intended release binary before
publishing an appropriate Music update and changing the bundled pin.

## Component previews

The other app previews include matching source archives, license notices,
build provenance and SHA256 sums. Their notes distinguish previous sampled
functional checks from the new branding builds. A requested manual Defender
scan failed to start on the build machine; antivirus clearance is not claimed.
The component versions selected for the next bundle are recorded in the lock
file. Existing installers continue to select the latest published complete
Workspace bundle; they do not assemble unreleased source commits.

Mail's Google Safe Browsing review remains pending, and provider browser logout
was blocked during earlier checks. These changes do not assert a new successful
Mail delivery, clean-host, service/reboot or full identity acceptance run.

## Current limitations

- Mail Contacts search and Drive attachment selection are future integrations.
- Calendar invitation delivery and migration of existing private Mail data into
  shared Calendar/Contacts storage are not implemented.
- Notes' direct Drive picker/Calendar adapter, Joplin synchronization and E2EE
  are not implemented. Its own login can drop query parameters.
- Independent users/devices, simultaneous Office merges, complex Office/PDF,
  large interrupted uploads and clean-host/service/reboot acceptance remain open.
- The Windows executables are unsigned previews. This bundle does not establish
  antivirus clearance; the requested manual scan failed to start on the builder.

Back up private configuration, databases, history and ordinary files before
upgrading. An existing workspace is not silently migrated by the new installer.

## Update and recovery window

Preview23 adds local owner maintenance for existing installations: select compatible
Workspace/app releases, verify downloads, snapshot private/program/file storage,
health-check startup and restore on failure. A durable journal and cached controller
support interrupted recovery and legacy launchers. Newer app releases without a
managed-update contract are shown for review. See UPDATES.md and UPDATE-PUBLISHING.md.
Backups remain until explicitly removed. Existing storage, identity, private Mail
and withheld-Music constraints remain. This does not claim clean-host, real reboot,
power-loss or Windows service acceptance.

## Branding cache refresh

Workspace app tiles, the shared switcher, the hub logo and favicon now use
content-addressed SVG filenames. Setup and update windows use those filenames too.
A new SVG receives a new URL, so browser/CDN entries for previous artwork cannot
hide a branding change. Versioned files have immutable caching and ETags; legacy
filenames remain available with revalidation. App component pins are unchanged.

## My Vault Windows preview

My Vault is the thirteenth app in fresh installations. Preview 26 pins the public
organization release 0.2.0 RC2, with branded light screens and optional MyAccounts
OIDC/session reuse. Its separate HTTPS origin and independent master-password
route do not require MyAccounts. Browser sign-in/save/edit/unlock, official SDK
registration, retained-data executable replacement, encrypted export/import,
offline CLI access and full age-encrypted restore/attachment checks passed in
isolated same-host Windows fixtures. Browser Create account submission, clean-host,
desktop/extension, service/reboot and independent security review remain open.
Read the installed
docs/MyVault/acceptance.json and SECURITY-REVIEW.md. This is an unsigned preview
and is not approved for important production credentials. Mobile is deferred.
Existing installations keep their selected app set during update; no new vault
or account is silently enabled. The installer does not update a live suite as
part of building or publishing this release.
