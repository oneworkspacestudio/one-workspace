import {parse, compileReact} from '@codedia/parser';
import {build} from 'esbuild';
import {readFile, writeFile, mkdir, copyFile} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';

process.chdir(fileURLToPath(new URL('.', import.meta.url)));
await mkdir('generated', {recursive: true});
await mkdir('dist', {recursive: true});
await copyFile('../../branding/icons/pet.svg', 'dist/icon.svg');
const source = await readFile('Pet.aui', 'utf8');
const tsx = compileReact(parse(source)).replaceAll('@/components/ui', '../ui');
await writeFile('generated/Pet.tsx', tsx);
await build({
  entryPoints: ['main.ts'], bundle: true, minify: true, format: 'esm', outdir: 'dist',
  assetNames: '[name]', loader: {'.woff2': 'file'}, jsx: 'automatic',
});
await writeFile('dist/index.html', '<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover,interactive-widget=resizes-content"><link rel="icon" href="./icon.svg"><title>Pet · One Workspace</title><link rel="stylesheet" href="./main.css"><script type="module" src="./main.js"></script></head><body><div id="root"></div></body></html>');
const licenses = [
  ['Three.js', 'three/LICENSE'],
  ['React', 'react/LICENSE'],
  ['React DOM', 'react-dom/LICENSE'],
  ['Scheduler', 'scheduler/LICENSE'],
  ['Bricolage Grotesque', '@fontsource-variable/bricolage-grotesque/LICENSE'],
  ['DM Sans', '@fontsource-variable/dm-sans/LICENSE'],
];
const notices = await Promise.all(licenses.map(async ([name, file]) => name+'\n\n'+(await readFile('node_modules/'+file,'utf8')).trim()));
await writeFile('THIRD-PARTY.txt', notices.join('\n\n'+'='.repeat(60)+'\n\n')+'\n');
