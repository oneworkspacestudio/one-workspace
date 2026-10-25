# My Cinema

My Cinema is a private streaming library for your films and series. Open it from
Workspace, select a viewing profile, then browse, search or resume a title. Adult
profiles use a PIN. Kids have separate watch history, My List and rating limits.

After unlocking an adult profile, the Workspace owner can choose **Add videos**
in the header. Connect a destination folder on the Workspace PC, then upload from
a computer or phone with progress and cancellation. A new folder can be created
explicitly; originals are retained and existing videos are never overwritten.
For videos already in My Drive, connect the same physical folders to index them
in place. Removing a library keeps every file. No media folder is chosen automatically.

The app includes a pinned Jellyfin portable engine and compatible FFmpeg. It runs
on loopback at the Cinema port plus 100. No public engine port, Docker, Node or
Windows service is required. Workspace manages real MyAccounts code/PKCE sign-in,
the launcher, process lifetime and verified backed-up updates. Bridge credentials
use Windows DPAPI and must be re-linked or recovered on the original Windows
account after relocation.

Use Workspace Backups for a complete stopped-process backup and restore, including
Cinema's profile database, encrypted bridge, Jellyfin configuration and metadata.
The standalone generic apphost backup flag is not a complete media-engine backup;
see the app's backup guidance. Large original videos need their own file backup.

See [My Cinema](https://github.com/oneworkspacestudio/my-cinema) for setup, exact
compatibility, security checks and the feature-completion matrix. Optional imports,
offline downloads and native TV packages are not advertised as shipped features.
