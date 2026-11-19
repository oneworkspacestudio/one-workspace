import { DEFAULTS, validSettings } from './routine.mjs';

const STORAGE_KEY = 'sunny-moon.routine.v1';
const deviceStore = globalThis.SunnyRoutineStore;

function saveRoutine(value) {
  const saved = JSON.stringify(value);
  if (deviceStore) {
    if (!deviceStore.write(saved)) throw new Error('Device save failed');
  } else {
    localStorage.setItem(STORAGE_KEY, saved);
  }
}

function validRoutine(value) {
  return value && validSettings(value.settings) &&
    ['lastVisit', 'nightVisit', 'expires', 'revision'].every(key => Number.isSafeInteger(value[key]) && value[key] >= 0) &&
    value.nightVisit <= value.lastVisit &&
    (value.lastVisit === 0 ? value.expires === 0 : value.expires > value.lastVisit);
}

function readRoutine() {
  const deviceSaved = deviceStore ? deviceStore.read() : null;
  const saved = deviceSaved ?? localStorage.getItem(STORAGE_KEY);
  const value = saved === null
    ? { settings: { ...DEFAULTS }, lastVisit: 0, nightVisit: 0, expires: 0, revision: 0 }
    : JSON.parse(saved);
  if (!validRoutine(value)) throw new Error('Saved routine is invalid');
  if (deviceStore && deviceSaved === null && saved !== null) saveRoutine(value);
  if (value.expires > 0 && value.expires <= Date.now()) {
    Object.assign(value, { lastVisit: 0, nightVisit: 0, expires: 0, revision: value.revision + 1 });
    saveRoutine(value);
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
    saveRoutine(routine);
  }
  return new Response(JSON.stringify({ routine, csrf: 'on-device' }), {
    headers: { 'Content-Type': 'application/json', 'X-Workspace-Account': 'on-device' },
  });
}
