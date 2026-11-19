import { DEFAULTS, validSettings } from './routine.mjs';

const STORAGE_KEY = 'sunny-moon.routine.v1';

function validRoutine(value) {
  return value && validSettings(value.settings) &&
    ['lastVisit', 'nightVisit', 'expires', 'revision'].every(key => Number.isSafeInteger(value[key]) && value[key] >= 0) &&
    value.nightVisit <= value.lastVisit &&
    (value.lastVisit === 0 ? value.expires === 0 : value.expires > value.lastVisit);
}

function readRoutine() {
  const saved = localStorage.getItem(STORAGE_KEY);
  const value = saved === null
    ? { settings: { ...DEFAULTS }, lastVisit: 0, nightVisit: 0, expires: 0, revision: 0 }
    : JSON.parse(saved);
  if (!validRoutine(value)) throw new Error('Saved routine is invalid');
  if (value.expires > 0 && value.expires <= Date.now()) {
    Object.assign(value, { lastVisit: 0, nightVisit: 0, expires: 0, revision: value.revision + 1 });
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
  }
  return value;
}

export async function deviceRoutineRequest(url, options = {}) {
  if (url !== './api/routine') throw new Error('Unknown device request');
  let routine = readRoutine();
  if (options.method === 'PUT') {
    const proposed = JSON.parse(options.body);
    if (!validRoutine(proposed)) return new Response(null, { status: 400 });
    if (proposed.revision !== routine.revision) return new Response(null, { status: 409 });
    routine = { ...proposed, revision: routine.revision + 1 };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(routine));
  }
  return new Response(JSON.stringify({ routine, csrf: 'on-device' }), {
    headers: { 'Content-Type': 'application/json', 'X-Workspace-Account': 'on-device' },
  });
}
