# One Workspace brief

Provide one Windows x64 setup EXE that installs the existing MyAccounts identity
server, custom FileBrowser Drive, embedded Docs/Sheets/Slides editors and optional
standalone Photos app. Add an app portal and a manifest contract for future apps.

Setup asks for the owner's name, account email and password, separate program,
private-data and file locations, optional Windows sign-in startup, and local or
public connection settings. Email sending is optional and must not require SMTP
for ordinary installations. Account administrators can manage users/passwords in
MyAccounts. Public mode retains Cloudflare tunnelling, with credentials supplied
at install time. Do not hardcode personal domains or credentials into packages.

Keep existing repositories and live previews intact. Store development/build
payloads off the system drive. Preserve original component licensing and honest
provenance, including local development components. Publish a separate extensible
Workspace repository with installation and integration documentation.
