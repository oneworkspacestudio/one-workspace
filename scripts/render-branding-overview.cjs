const fs = require('node:fs');
const path = require('node:path');
const { createHash } = require('node:crypto');
const sharp = require(process.env.SHARP_MODULE || 'sharp');
const { apps: directory, escape } = require('./app-branding.cjs');

const root = path.resolve(__dirname, '..');
const apps = directory.map(app => [app.icon, app.name]);

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
    return `<text x="${x}" y="${y}" text-anchor="middle" font-family="Arial, sans-serif" font-size="22" fill="#24324b">${escape(name)}</text>`;
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
  const readmePath = path.join(root, 'README.md');
  let readme = fs.readFileSync(readmePath, 'utf8');
  function replaceSection(name, content) {
    const expression = new RegExp(`<!-- ${name}:start -->[\\s\\S]*?<!-- ${name}:end -->`);
    if (!expression.test(readme)) throw new Error(`Missing README section: ${name}`);
    readme = readme.replace(expression, `<!-- ${name}:start -->\n${content}\n<!-- ${name}:end -->`);
  }
  replaceSection('app-overview', `![The complete One Workspace app collection](docs/branding/app-suite-${hash}.png)`);
  const rows = directory.map(app => `    <tr>\n      <td align="center" width="56"><img src="internal/branding/icons/${app.icon}.svg" alt="" width="40" height="40"></td>\n      <td><a href="${escape(app.url)}"><strong>${escape(app.name)}</strong></a></td>\n      <td>${escape(app.description)}</td>\n    </tr>`).join('\n');
  replaceSection('app-directory', `<table>\n  <thead><tr><th colspan="2">App</th><th>What you can do</th></tr></thead>\n  <tbody>\n${rows}\n  </tbody>\n</table>`);
  fs.writeFileSync(readmePath, readme);
  // The small-size preview uses the same directory, including newly added apps.
  const previewPath = path.join(root, 'docs/branding/index.html');
  const preview = fs.readFileSync(previewPath, 'utf8').replace(/const apps=\[[\s\S]*?;/, `const apps=${JSON.stringify(apps)};`);
  fs.writeFileSync(previewPath, preview);
  console.log(output);
  console.log(versioned);
}

main().catch(error => { console.error(error); process.exitCode = 1; });
