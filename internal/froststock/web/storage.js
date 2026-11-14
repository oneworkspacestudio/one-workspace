let account, csrf, database;
let lastRemote = [];
let remoteUnavailable = false;
let incomplete = 0;
let mutation = 0;
const base = new URL('./api/items/', import.meta.url);

async function request(url, init = {}) {
  const headers = new Headers(init.headers);
  if (account) headers.set('X-Workspace-Expected-Account', account);
  if (csrf) headers.set('X-CSRF-Token', csrf);
  const response = await fetch(url, { ...init, headers });
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw Object.assign(new Error(body.error || body.message || 'My Drive is unavailable. Try again.'), { status: response.status });
  }
  return response.status === 204 ? null : response.json();
}

export async function connect() {
  const session = await request(new URL('./api/me', import.meta.url));
  if (session.user.role !== 'owner') {
    throw Object.assign(new Error('Open Froststock with the Workspace owner account.'), { status: 403 });
  }
  account = session.user.id;
  csrf = session.csrf;
  database = new Promise((resolve, reject) => {
    const open = indexedDB.open('workspace-froststock-' + account, 1);
    open.onupgradeneeded = () => open.result.createObjectStore('pending', { keyPath: 'id' });
    open.onerror = () => reject(open.error);
    open.onsuccess = () => resolve(open.result);
  });
}

async function transaction(mode, operation) {
  const db = await database;
  return new Promise((resolve, reject) => {
    const tx = db.transaction('pending', mode);
    const result = operation(tx.objectStore('pending'));
    tx.oncomplete = () => resolve(result?.result);
    tx.onabort = () => reject(tx.error);
  });
}

const pending = () => transaction('readonly', store => store.getAll());
const unstage = id => transaction('readwrite', store => store.delete(id));
export const saveItem = item => transaction('readwrite', store => store.add(item));
export const addItems = items => transaction('readwrite', store => {
  for (const item of items) store.add(item);
});
export const itemPhotoURL = (id, name) => new URL(encodeURIComponent(id) + '/' + name, base).href;
export const driveStatus = () => ({ unavailable: remoteUnavailable, incomplete });

export async function allItems() {
  const version = mutation;
  try {
    const response = await request(new URL('./api/items', import.meta.url));
    if (mutation === version) {
      lastRemote = response.items;
      incomplete = response.incomplete;
    }
    remoteUnavailable = false;
  } catch (error) {
    if (error.status === 401 || error.status === 403 || error.status === 409) throw error;
    remoteUnavailable = true;
  }
  const remoteIDs = new Set(lastRemote.map(item => item.id));
  const staged = (await pending()).filter(item => !remoteIDs.has(item.id));
  return lastRemote.concat(staged);
}

function serialized(item) {
  const { itemPhoto, datePhoto, ...value } = item;
  return value;
}

// A successful upload can lose its response. Drive returns the existing entry
// for a repeated ID, so the browser can safely retry the pair.
export async function syncPending() {
  let failure;
  for (const item of await pending()) {
    try {
      const form = new FormData();
      form.set('item', item.itemPhoto, 'item.jpg');
      form.set('expiry', item.datePhoto, 'expiry.jpg');
      form.set('record', JSON.stringify(serialized(item)));
      const synced = await request(new URL(encodeURIComponent(item.id), base), { method: 'POST', body: form });
      mutation++;
      lastRemote = lastRemote.filter(entry => entry.id !== synced.id).concat(synced);
      await unstage(item.id);
    } catch (error) {
      if (![400, 409].includes(error.status)) throw error;
      failure ??= error;
    }
  }
  if (failure) throw failure;
}

export async function jobPhotos(item) {
  if (item.itemPhoto && item.datePhoto) return item;
  const blobs = await Promise.all(['item.jpg', 'expiry.jpg'].map(async name => {
    const response = await fetch(itemPhotoURL(item.id, name), { headers: { 'X-Workspace-Expected-Account': account } });
    if (!response.ok) throw Object.assign(new Error('Both photos must be present in My Drive.'), { status: response.status });
    return response.blob();
  }));
  return { ...item, itemPhoto: blobs[0], datePhoto: blobs[1] };
}

export async function updateItem(id, update) {
  const before = lastRemote.find(item => item.id === id);
  if (!before) throw new Error('This item is waiting to save to My Drive.');
  const after = update(before);
  if (!after) return;
  try {
    const result = await request(new URL(encodeURIComponent(id), base), {
      method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(serialized(after)),
    });
    mutation++;
    lastRemote = lastRemote.map(item => item.id === id ? result : item);
    return result;
  } catch (error) {
    if (error.status === 409) await allItems();
    throw error;
  }
}

export async function removeItem(id) {
  const before = lastRemote.find(item => item.id === id);
  if (!before) {
    const discard = async () => {
      if (lastRemote.some(item => item.id === id)) throw new Error('This item has just synced. Open it again to delete it.');
      await unstage(id);
    };
    if (navigator.locks?.request) return navigator.locks.request('froststock-drive', discard);
    return discard();
  }
  await request(new URL(encodeURIComponent(id), base), {
    method: 'DELETE', headers: { 'If-Match': '"' + before.revision + '"' },
  });
  mutation++;
  lastRemote = lastRemote.filter(item => item.id !== id);
}
