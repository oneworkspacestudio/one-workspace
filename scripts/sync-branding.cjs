// SVGs are the originals. PNG/ICO files are generated without changing the artwork.
const fs = require('node:fs');
const path = require('node:path');
const sharp = require(process.env.SHARP_MODULE || 'sharp');

const root = path.resolve(__dirname, '..');
const projects = process.argv[2] ? path.resolve(process.argv[2]) : path.dirname(root);
const icons = path.join(root, 'internal/branding/icons');

function copy(source, target) {
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.copyFileSync(source, target);
}

async function main() {
  const generated = path.join(root, 'docs/branding');
  fs.mkdirSync(generated, { recursive: true });
  for (const name of ['one_workspace', 'my_drive', 'my_office', 'my_photos', 'my_git', 'my_music', 'photo_editor', 'my_mail', 'my_accounts', 'my_notes', 'my_calendar', 'my_contacts']) {
    const svg = path.join(icons, name + '.svg');
    const sizes = [16, 24, 32, 48, 64, 128, 256];
    const images = await Promise.all(sizes.map(size => sharp(svg).resize(size, size).png().toBuffer()));
    const header = Buffer.alloc(6 + images.length * 16);
    header.writeUInt16LE(1, 2);
    header.writeUInt16LE(images.length, 4);
    let offset = header.length;
    images.forEach((image, i) => {
      const entry = 6 + i * 16;
      header[entry] = sizes[i] === 256 ? 0 : sizes[i];
      header[entry + 1] = header[entry];
      header.writeUInt16LE(1, entry + 4);
      header.writeUInt16LE(32, entry + 6);
      header.writeUInt32LE(image.length, entry + 8);
      header.writeUInt32LE(offset, entry + 12);
      offset += image.length;
    });
    fs.writeFileSync(path.join(generated, name + '.ico'), Buffer.concat([header, ...images]));
    for (const size of [32, 180, 192, 512]) {
      await sharp(svg).resize(size, size).png().toFile(path.join(generated, `${name}-${size}.png`));
    }
  }
  const at = (repo, relative) => path.join(projects, repo, relative);
  const asset = (name, extension) => path.join(generated, name + extension);
  const svg = name => path.join(icons, name + '.svg');
  for (const name of ['my_drive', 'my_office', 'my_photos', 'photo_editor']) {
    copy(svg(name), at('filebrowser-drive', `frontend/src/assets/branding/${name}.svg`));
  }
  copy(svg('my_drive'), at('filebrowser-drive', 'frontend/public/img/icons/favicon.svg'));
  copy(asset('my_drive', '-512.png'), at('filebrowser-drive', 'frontend/public/img/icons/favicon.png'));
  for (const name of ['logo', 'favicon']) {
    copy(svg('photo_editor'), at('minipaint-studio', `images/${name}.svg`));
    copy(svg('my_git'), at('gitea', `assets/${name}.svg`));
    copy(svg('my_git'), at('gitea', `public/assets/img/${name}.svg`));
    copy(svg('my_git'), at('gitea', `integrations/workspace/branding/${name}.svg`));
  }
  copy(asset('photo_editor', '-192.png'), at('minipaint-studio', 'images/favicon.png'));
  for (const directory of ['public/assets/img', 'integrations/workspace/branding']) {
    copy(asset('my_git', '-32.png'), at('gitea', `${directory}/favicon.png`));
    copy(asset('my_git', '-180.png'), at('gitea', `${directory}/apple-touch-icon.png`));
    copy(asset('my_git', '-512.png'), at('gitea', `${directory}/logo.png`));
  }
  copy(svg('my_music'), at('swingmusic', 'src/swingmusic/assets/my_music.svg'));
  copy(asset('my_music', '.ico'), at('swingmusic', 'src/swingmusic/assets/logo-fill.light.ico'));
  for (const size of [192, 512]) {
    copy(asset('my_music', `-${size}.png`), at('swingmusic', `src/swingmusic/assets/my_music-${size}.png`));
  }
  copy(svg('my_mail'), at('stalwart-mail', 'workspace-mail/assets/mail.svg'));
  copy(svg('my_calendar'), at('my-calendar', 'web/public/calendar.svg'));
  copy(svg('my_contacts'), at('my-calendar', 'web/public/contacts.svg'));
  copy(svg('my_notes'), path.join(process.env.NOTES_ROOT || at('notes', ''), 'web/public/notes.svg'));
  copy(svg('my_accounts'), path.join(process.env.MYACCOUNTS_ROOT || at('myaccounts', ''), 'web/public/myaccounts.svg'));
  console.log('Updated suite branding assets in ' + projects);
}

main().catch(error => { console.error(error); process.exitCode = 1; });
