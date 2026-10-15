# October 15 release status

The supplied eleven app SVGs are applied to app headers, favicons, Windows icons
and repository README logos. Component executables were built locally; GitHub
Actions is disabled for all nine organization repositories.

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
