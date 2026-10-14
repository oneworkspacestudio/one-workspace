# Suite logos

The seven supplied SVGs live in `internal/branding/icons/`. The original vector
artwork is unchanged, including the transparent area outside each rounded tile.

| Logo | Application surfaces |
| --- | --- |
| One Workspace | Portal header, setup wizard, favicon, manager and installer EXE |
| My Drive | Drive header, login fallback, favicon, PWA icons and FileBrowser EXE |
| My Office | Workspace Docs tile and the built-in Office editor toolbar |
| My Photos | Portal tile, integrated/standalone Photos header, sidebar, favicon and EXE |
| My Git | Portal tile, Gitea logo/favicon/touch icon and standalone launcher EXE |
| My Music | Portal tile, Music header, favicon, PWA icons and Music EXE |
| Photo Editor | Studio toolbar/menu, Drive sidebar, favicon and Studio EXE |

MyAccounts keeps its existing branding because no replacement was supplied.
App names, routes, authentication and upstream license notices are unchanged.
New add-on apps use the Workspace icon until they receive their own branding.

To refresh the copies in sibling repositories, install Sharp in a local tools
folder (`npm install --prefix dist/branding-tools sharp`), then run:

```powershell
$env:SHARP_MODULE = "$PWD/dist/branding-tools/node_modules/sharp"
node scripts/sync-branding.cjs I:/Projects
./scripts/build-branding-resources.ps1 -Go I:/Projects/tooling/go1.27.2/go/bin/go.exe
```

The sync command generates transparent PNG fallbacks and multi-size Windows ICOs
in `docs/branding/`; the resource command embeds those ICOs in Windows amd64 Go
resources. Photos and Studio build scripts carry those resources into staging.
Music's client adapter applies its assets before the Vite build; PyInstaller uses
the same ICO. The Gitea launcher uses Gitea's supported custom public-assets path
so the bundled upstream executable receives the branding without a core rebuild.

Build the apps locally after syncing. Running installations and already published
release assets continue to use their embedded assets until deliberately updated.
GitHub Actions is not needed or enabled. Open `docs/branding/index.html` for a
light/dark background and small-size preview of all seven icons.
