const fs = require('node:fs');
const path = require('node:path');
const sharp = require(process.env.SHARP_MODULE || 'sharp');

const root = path.resolve(__dirname, '..');
const apps = [
  ['one_workspace', 'Workspace'], ['my_drive', 'Drive'], ['my_office', 'Docs'],
  ['my_photos', 'Photos'], ['my_music', 'Music'], ['my_git', 'Git'],
  ['my_accounts', 'MyAccounts'], ['photo_editor', 'Photo Editor'], ['my_mail', 'Mail'],
  ['my_calendar', 'Calendar'], ['my_contacts', 'Contacts'], ['my_notes', 'Notes'],
];

async function main() {
  const labels = apps.map(([, name], i) => {
    const x = 120 + (i % 6) * 240;
    const y = 282 + Math.floor(i / 6) * 192;
    return `<text x="${x}" y="${y}" text-anchor="middle" font-family="Arial, sans-serif" font-size="22" fill="#24324b">${name}</text>`;
  }).join('');
  const canvas = Buffer.from(`<svg xmlns="http://www.w3.org/2000/svg" width="1440" height="520" viewBox="0 0 1440 520"><rect width="1440" height="520" rx="24" fill="#f6f8fc"/><text x="720" y="62" text-anchor="middle" font-family="Arial, sans-serif" font-size="36" font-weight="700" fill="#162640">One Workspace</text><text x="720" y="100" text-anchor="middle" font-family="Arial, sans-serif" font-size="21" fill="#52647d">Your apps, together</text>${labels}</svg>`);
  const layers = await Promise.all(apps.map(async ([icon], i) => ({
    input: await sharp(path.join(root, 'internal/branding/icons', icon + '.svg')).resize(128, 128).png().toBuffer(),
    left: 56 + (i % 6) * 240,
    top: 132 + Math.floor(i / 6) * 192,
  })));
  const output = path.join(root, 'docs/branding/app-suite.png');
  fs.mkdirSync(path.dirname(output), { recursive: true });
  await sharp(canvas).composite(layers).png().toFile(output);
  console.log(output);
}

main().catch(error => { console.error(error); process.exitCode = 1; });
