# My Remote

My Remote is included in One Workspace104 and later. It appears under
Administration in setup and in the shared app picker. Earlier installations
receive it as a new default app when Workspace is updated, unless explicitly
excluded. Use the normal protected updater; a separate desktop shortcut is not
required to install the app.

Sign in with your MyAccounts owner identity. Desktop access initially stays off.
Choose **Set up remote access**, enter a computer name and your own PIN, and keep
that Windows user signed in. Owner connections require both sign-in and PIN;
support visitors need a one-use invitation and your approval. Stop sharing to
revoke a session. Signing out of Workspace immediately closes owner desktop
connections. Approved support grants have separate stop/expiry controls. Do not publish the private VNC or bridge listener.

The app is available through Workspace's `/remote/` route. Network access uses
your existing authenticated HTTPS Workspace gateway. A public page alone does
not grant desktop access. Each claimed remote deployment still needs its own
network/desktop acceptance; synthetic phone layouts do not prove a physical
phone or two-PC session.

Workspace owns the complete runtime and its dedicated Private/remote data folder.
Updates use the same protected backup, verification and recovery controls as the
other apps. Missing or modified runtime files cause startup to fail visibly.

An unconfigured standalone preview2 can migrate its known launchers with backups
in its old data folder. Existing programs and data remain. An enrolled standalone
host requires explicit migration review; unrelated TightVNC settings are never
overwritten. If a Workspace rollback removes the new app, the retained legacy
launcher backups remain available for explicit recovery.

This preview does not include access before Windows sign-in, secure desktop,
audio, file transfer or individual monitor selection. See the
[app repository](https://github.com/oneworkspacestudio/my-remote) for capabilities,
licences and component releases. Maintainers must follow [APP-LIFECYCLE.md](APP-LIFECYCLE.md).
