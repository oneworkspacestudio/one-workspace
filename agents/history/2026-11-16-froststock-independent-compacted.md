# Independent Froststock distribution

The owner asked to release Froststock as its own Workspace family app. The
application source, builds and release manifests move to
oneworkspacestudio/froststock. Historical Workspace release tags retain the
earlier bundled versions. Preview.140 is the one-time updater migration and
its Workspace bundle excludes Froststock and all its owned package files.

Froststock is registered for optional installation through Updates, with
MyAccounts and My Drive dependencies. Compatibility integration 22 directs
older controllers to the new installer. Installed app IDs, saved ports,
OIDC, My Drive files and private settings retain their format and location.
The new updater accepts only Froststock's organization repository, exact
launch contract, owned paths and verified asset checksums. Install and
update plans leave the Workspace manager outside the Froststock selection.

The app builds against a pinned Workspace AppHost source commit using
public wrappers around the existing private-setting and path checks. The
fully independent release can be promoted without creating a Workspace tag.
Bundled-era catalogs remain readable and old installed Froststock packages
are retained when applying the migration installer.

Go catalog, distribution, suite, SDK and moved app tests passed. Focused
regressions cover independent selection with a suite update available,
an optional install after opt-out, changed launch metadata and preservation
of the installed Froststock package during the manager handoff.
