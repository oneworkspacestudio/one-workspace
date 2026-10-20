const fs = require('node:fs');
const path = require('node:path');
const { createHash } = require('node:crypto');
const sharp = require(process.env.SHARP_MODULE || 'sharp');

const root = path.resolve(__dirname, '..');
const apps = [
  ['one_workspace', 'Workspace'], ['my_drive', 'Drive'], ['my_office', 'Docs'],
  ['my_photos', 'Photos'], ['my_music', 'Music'], ['my_git', 'Git'],
  ['my_accounts', 'MyAccounts'], ['photo_editor', 'Photo Editor'], ['my_mail', 'Mail'],
  ['my_calendar', 'Calendar'], ['my_contacts', 'Contacts'], ['my_notes', 'Notes'],
  ['my_tube', 'My Tube'], ['my_chat', 'Chat'], ['my_vault', 'My Vault'],
  ['my_hosting', 'My Hosting'], ['my_tasks', 'My Tasks'], ['my_pdf', 'My PDF'], ['my_monitor', 'My Monitor'], ['my_search', 'My Search'], ['backups', 'Backups'],
  ['my_slicer', 'My Slicer'], ['my_social', 'My Social'], ['my_clips', 'My Clips'],
  ['my_maps', 'My Maps'], ['my_meet', 'My Meet'],
];

const columns = 6;
const cellWidth = 210;
const width = columns * cellWidth;
const rowHeight = 192;
const height = 136 + Math.ceil(apps.length / columns) * rowHeight;

function position(i) {
  const row = Math.floor(i / columns);
  const count = Math.min(columns, apps.length - row * columns);
  return {
    x: ((columns - count) / 2 + i % columns + 0.5) * cellWidth,
    y: 132 + row * rowHeight,
  };
}

async function main() {
  const labels = apps.map(([, name], i) => {
    const { x, y: top } = position(i);
    const y = top + 150;
    return `<text x="${x}" y="${y}" text-anchor="middle" font-family="Arial, sans-serif" font-size="22" fill="#24324b">${name}</text>`;
  }).join('');
  const canvas = Buffer.from(`<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}"><rect width="${width}" height="${height}" rx="24" fill="#f6f8fc"/><text x="${width / 2}" y="62" text-anchor="middle" font-family="Arial, sans-serif" font-size="36" font-weight="700" fill="#162640">One Workspace</text><text x="${width / 2}" y="100" text-anchor="middle" font-family="Arial, sans-serif" font-size="21" fill="#52647d">Your apps, together</text>${labels}</svg>`);
  const layers = await Promise.all(apps.map(async ([icon], i) => ({
    input: await sharp(path.join(root, 'internal/branding/icons', icon + '.svg')).resize(128, 128).png().toBuffer(),
    left: position(i).x - 64,
    top: position(i).y,
  })));
  const output = path.join(root, 'docs/branding/app-suite.png');
  fs.mkdirSync(path.dirname(output), { recursive: true });
  const png = await sharp(canvas).composite(layers).png().toBuffer();
  fs.writeFileSync(output, png);
  const hash = createHash('sha256').update(png).digest('hex').slice(0, 8);
  const versioned = path.join(root, `docs/branding/app-suite-${hash}.png`);
  fs.writeFileSync(versioned, png);
  console.log(output);
  console.log(versioned);
}

main().catch(error => { console.error(error); process.exitCode = 1; });
