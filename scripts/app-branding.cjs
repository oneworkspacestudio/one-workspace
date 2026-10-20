const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const apps = require('../docs/branding/apps.json');
const names = fs.readdirSync(path.join(root, 'internal/branding/icons')).filter(name => name.endsWith('.svg')).map(name => name.slice(0, -4)).sort();
if (new Set(apps.map(app => app.id)).size !== apps.length ||
    new Set(apps.map(app => app.icon)).size !== apps.length ||
    JSON.stringify(apps.map(app => app.icon).sort()) !== JSON.stringify(names)) {
  throw new Error('The app directory must list every SVG icon exactly once. Update docs/branding/apps.json.');
}
const escape = value => value.replace(/[&<>"']/g, char => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[char]));
module.exports = { apps, root, escape };
