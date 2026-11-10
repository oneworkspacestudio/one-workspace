import * as parser from '@codedia/parser';
import * as esbuild from 'esbuild';
import {readFile,writeFile,mkdir,copyFile} from 'node:fs/promises';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const root=path.dirname(fileURLToPath(import.meta.url));process.chdir(root);
// Build dependencies are provided by the locked shared app toolchain. They are
// bundled into the output; node/npm are not needed on an installed computer.
await mkdir('generated',{recursive:true});await mkdir('dist',{recursive:true});
await copyFile('../../branding/icons/game_studio.svg','dist/icon.svg');
const source=await readFile('Studio.aui','utf8');
const tsx=parser.compileReact(parser.parse(source)).replaceAll('@/components/ui','../main');
await writeFile('generated/Studio.tsx',tsx);
await esbuild.build({absWorkingDir:root,entryPoints:['main.ts'],bundle:true,minify:true,format:'esm',outdir:'dist',jsx:'automatic',loader:{'.ts':'tsx','.tsx':'tsx'}});
await writeFile('dist/index.html','<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><link rel="icon" href="./icon.svg"><title>Game Studio · One Workspace</title><link rel="stylesheet" href="./main.css"><script type="module" src="./main.js"></script></head><body><div id="root"></div></body></html>');
