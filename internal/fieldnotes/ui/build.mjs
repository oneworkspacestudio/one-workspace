import { parse, compileReact } from '@codedia/parser';
import { build } from 'esbuild';
import { mkdir, readFile, writeFile } from 'node:fs/promises';

await mkdir('generated', { recursive: true });
const source = await readFile('Reader.aui', 'utf8');
await writeFile('generated/Reader.tsx', compileReact(parse(source)).replaceAll('@/components/ui', '../components'));
await build({ entryPoints: ['main.ts'], bundle: true, minify: true, format: 'esm', outdir: 'dist', jsx: 'automatic' });
const licenses = await Promise.all(['react', 'react-dom', 'scheduler'].map(async name => name + '\n' + await readFile('node_modules/' + name + '/LICENSE', 'utf8')));
await writeFile('dist/THIRD-PARTY-NOTICES.txt', licenses.join('\n\n'));
await writeFile('dist/index.html', '<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="theme-color" content="#faf8f1"><title>AI Fieldnotes</title><link rel="stylesheet" href="./main.css"><script type="module" src="./main.js"></script></head><body><div id="root"></div></body></html>');
