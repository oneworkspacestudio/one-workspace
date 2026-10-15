<div align="center">
  <img src="internal/branding/icons/one_workspace.svg" alt="" width="64" height="64">
  <h1>One Workspace</h1>
  <p>Your apps and files, together on your own Windows PC.</p>
  <p><a href="https://github.com/oneworkspacestudio/one-workspace/releases"><strong>Download for Windows</strong></a> · <a href="#get-started">Get started</a> · <a href="docs/TECHNICAL.md">Technical guide</a></p>
</div>

![One Workspace brings files, documents, photos and more together](docs/images/one-workspace-showcase.png)

*An illustrated look at the workspace. Check the download notes for the apps included in each package.*

## A place for your everyday work

Browse files, edit documents, enjoy your photos and music, and move between apps with one account. One Workspace brings them together in a browser, while your files stay on the computer and in the folders you choose.

- **Your files stay yours.** Work with ordinary folders on your own storage.
- **One account across your apps.** Sign in and open what you need from the workspace.
- **Room to grow.** Add more apps as your needs change.

## Meet your apps

| App | What you can do |
| --- | --- |
| <img src="internal/branding/icons/my_drive.svg" alt="" width="28" height="28"> **Drive** | Browse, upload, organise and share your files. |
| <img src="internal/branding/icons/my_office.svg" alt="" width="28" height="28"> **Office** | Create and edit documents, spreadsheets, slides and PDFs. |
| <img src="internal/branding/icons/my_photos.svg" alt="" width="28" height="28"> **Photos** | Browse your pictures and explore folders as albums. |
| <img src="internal/branding/icons/my_music.svg" alt="" width="28" height="28"> **Music** | Listen to your music collection. |
| <img src="internal/branding/icons/my_git.svg" alt="" width="28" height="28"> **Git** | Keep code projects organised and collaborate. |
| **MyAccounts** | Manage people, sign-ins and passwords. |

**More options:** [Photo Editor](https://github.com/fahimc/miniPaint) adds drawing, layers and image editing as a separate integration. Mail includes a Gmail-style inbox, automatic MyAccounts mailboxes and guided hosting setup.

Included apps vary by download. The [release notes](https://github.com/oneworkspacestudio/one-workspace/releases) show what is available now.

## Get started

You need a **64-bit Windows PC** and at least **2 GB of free space**, plus room for your files. See the [storage requirements](docs/TECHNICAL.md#requirements-and-storage) before choosing your installation folders.

1. **Download the installer.** Open [Downloads](https://github.com/oneworkspacestudio/one-workspace/releases) and choose `WorkspaceSetup-windows-amd64.exe`.
2. **Make it yours.** Run the setup wizard and enter your name, email address and a password.
3. **Choose one folder.** Workspace organises apps, private settings and your files inside it. Choose local access or connect your domain with Cloudflare.
4. **Install the latest apps.** Setup downloads the newest published Workspace bundle from GitHub, verifies its checksums, and shows progress. Launch Workspace and sign in.

Your computer needs to stay on while you use the apps. Access through your own web address is optional; the [online setup guide](docs/TECHNICAL.md#local-and-public-addresses) explains what is needed.

<details>
<summary>Take a look at the setup wizard</summary>

![One Workspace setup wizard](docs/screenshots/setup-new-logo.png)

</details>

> **Currently in preview.** One Workspace is still being developed. Keep backups of important files and read the download notes before installing or changing an existing workspace.

### Keeping new installs current

Normal installation checks the public GitHub releases each time and downloads the
newest published Workspace application bundle, including preview releases. Each
bundle declares the versions of its apps as a compatible set. An installer you
saved earlier will therefore install the latest published bundle. Unreleased Git
commits are not installed. Setup stops with a retry message if GitHub or checksum
verification fails; it does not silently install the old embedded versions.

Setup also checks [Mail's organization releases](https://github.com/oneworkspacestudio/mail/releases)
independently, and installs the newest published Mail executable with a supported
Workspace compatibility manifest. It verifies GitHub's SHA256 and size before
installation and records the exact Mail version and release source. An invalid,
unavailable or incompatible Mail release stops online installation with guidance.
Download an updated installer if a future Mail integration version is unsupported.

For an intentional offline installation, run
`WorkspaceSetup-windows-amd64.exe --offline`. This uses the versions bundled with
that installer. Existing installations are not automatically upgraded by this
new-install behavior. Mail hosting checks incoming SMTP connectivity separately
and shows guidance for router, Windows firewall or ISP restrictions.

## A few useful things to know

### Where are my files kept?

In the folders you choose on your own computer. They stay as ordinary files that you can access outside Workspace too.

### Can I invite other people?

Yes. Use **Manage accounts** to add people and reset passwords. The starting Drive folder is shared, so review who can access it before inviting others. [Account and sharing details](docs/TECHNICAL.md#accounts-and-administration).

### Do I need an email service?

No. You can create accounts and reset passwords without one. Email verification and password-recovery emails are optional features you can set up later.

### Can I add more apps?

Yes. Workspace can be extended with additional apps. Connecting a new app currently takes some configuration; see [adding apps](docs/TECHNICAL.md#add-another-app).

## Help and further reading

| Looking for… | Start here |
| --- | --- |
| A download or what changed | [Downloads and release notes](https://github.com/oneworkspacestudio/one-workspace/releases) |
| Setup, online access or backups | [Technical guide](docs/TECHNICAL.md) |
| A problem to report | [Report an issue](https://github.com/oneworkspacestudio/one-workspace/issues) |
| Licensing and component information | [Licenses and sources](docs/COMPONENTS.md) |
