import { connect, allItems, saveItem, updateItem, removeItem, addItems, itemPhotoURL, syncPending, jobPhotos, driveStatus } from './storage.js';
import { isDate, offsetDate, daysFromToday, formatDate, formatToday } from './dates.js';
import { scanLabels } from './recognition.js';

const $ = selector => document.querySelector(selector);
const icons = 'snowflake layout-grid boxes clock-3 history settings camera chevron-right scan-line search shield-check calendar-clock triangle-alert package image x check loader-circle clipboard-list gallery-horizontal-end calendar-days download list-filter check-circle-2 image-plus upload'.split(' ');
const symbol = name => '<i aria-hidden="true" class="icon" data-icon="' + name + '"></i>';
const iconStyle = document.createElement('style');
iconStyle.textContent = icons.map(name => '[data-icon="' + name + '"]{--icon:url("./assets/icons/' + name + '.svg")}').join('\n');
document.head.append(iconStyle);
const esc = value => String(value ?? '').replace(/[&<>"']/g, char => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char]);
const id = () => globalThis.crypto?.randomUUID?.() ?? Date.now().toString(36) + '-' + Math.random().toString(36).slice(2);
const scanner = $('#scanner');
const details = $('#details');
const video = $('#camera');
const toast = $('#toast');
let toastTimer;
let items = [];
let demo = localStorage.getItem('froststock-demo') === 'on';
let filter = 'all';
let activePage;
let inQueue = false;
let syncing = false;
let queueMessage = '';
let currentDetail;
let stream;
let cameraRequest = 0;
let startingCamera = false;
let processingPhoto = false;
let itemPhoto;
let firstPreview;
let sessionCount = 0;
let persistentRequested = false;
let ownerReady = false;
const thumbnailURLs = new Map();
const detailURLs = new Set();
const preferences = {
  dateOrder: localStorage.getItem('froststock-date-order') || 'dmy',
  vision: localStorage.getItem('froststock-vision') === 'on',
};
const channel = typeof BroadcastChannel === 'function' ? new BroadcastChannel('froststock-inventory') : null;

const sample = [
  { id: 'sample-salmon', name: 'Atlantic salmon fillets', category: 'Fish & seafood', expiry: offsetDate(2), sampleImage: 'salmon', dateType: 'Use by' },
  { id: 'sample-berries', name: 'Mixed summer berries', category: 'Fruit', expiry: offsetDate(6), sampleImage: 'berries', dateType: 'Best before' },
  { id: 'sample-peas', name: 'Garden peas', category: 'Vegetables', expiry: offsetDate(67), sampleImage: 'peas', dateType: 'Best before' },
  { id: 'sample-dumplings', name: 'Chicken gyoza', category: 'Prepared meals', expiry: offsetDate(-3), sampleImage: 'dumplings', dateType: 'Best before' },
].map(item => ({ ...item, demo: true, state: 'ready', quantity: 1, created: Date.now() }));

const routes = {
  overview: ['Overview', 'Freezer overview', 'See what’s in your freezer and what needs using first.'],
  inventory: ['Inventory', 'Freezer inventory', 'Search, organise and review the items you’ve scanned.'],
  expiring: ['Expiring soon', 'Expiring soon', 'See what’s due in the next 7 days or past its label date.'],
  history: ['History', 'Freezer history', 'Find items you’ve marked as used.'],
  settings: ['Settings'],
};
const states = {
  expired: ['Past date', 'triangle-alert'], soon: ['Expiring soon', 'calendar-clock'],
  ready: ['In date', 'check-circle-2'], review: ['Needs review', 'clipboard-list'],
  queued: ['Reading labels', 'loader-circle'], pending: ['Saving to Drive', 'loader-circle'],
  used: ['Used', 'history'], none: ['No expiry date', 'calendar-days'],
};
const categories = ['Meat & poultry', 'Fish & seafood', 'Fruit', 'Vegetables', 'Prepared meals', 'Bread & bakery', 'Dairy & desserts', 'Other'];
const dateTypes = ['Expiry date', 'Use by', 'Best before', 'Sell by'];

function notification(message) {
  clearTimeout(toastTimer);
  toast.textContent = message;
  toast.hidden = false;
  toastTimer = setTimeout(() => toast.hidden = true, 3800);
}

function storeDemo(value) {
  demo = value;
  localStorage.setItem('froststock-demo', value ? 'on' : 'off');
  render();
}

function status(item) {
  if (item.used) return 'used';
  if (!item.demo && !item.revision) return 'pending';
  if (item.state === 'queued') return 'queued';
  if (item.review) return 'review';
  if (!item.expiry) return 'none';
  const days = daysFromToday(item.expiry);
  return days < 0 ? 'expired' : days <= 7 ? 'soon' : 'ready';
}

function nameOf(item) { return item.name || (item.state === 'queued' ? 'New scan' : 'Unnamed item'); }

function previewURL(item) {
  if (item.demo) return './assets/' + item.sampleImage + '.webp';
  if (item.revision) return itemPhotoURL(item.id, 'item.jpg');
  if (!item.itemPhoto) return '';
  if (!thumbnailURLs.has(item.id)) thumbnailURLs.set(item.id, URL.createObjectURL(item.itemPhoto));
  return thumbnailURLs.get(item.id);
}

function refreshSoonBadge(count) {
  const badge = $('#nav-soon');
  badge.hidden = count === 0;
  badge.textContent = count;
}

function render() {
  if (!activePage) return;
  const records = demo ? items.concat(sample) : items;
  const current = records.filter(item => !item.used);
  const counts = { all: current.length, soon: 0, expired: 0, review: 0, queued: 0, pending: 0 };
  for (const item of current) {
    const itemStatus = status(item);
    if (Object.hasOwn(counts, itemStatus)) counts[itemStatus]++;
  }
  refreshSoonBadge(counts.soon);
  $('#toggle-demo').textContent = demo ? 'Clear examples' : 'Show examples';
  if (activePage === 'settings') return;
  $('#metrics').innerHTML = [
    ['all', 'Items in freezer', 'boxes'], ['soon', 'Expiring in 7 days', states.soon[1]],
    ['expired', 'Past label date', states.expired[1]], ['review', 'Needs review', states.review[1]],
  ].map(([key, label, icon]) => '<button type="button" class="stat ' + key + '" data-action="metric" data-value="' + key +
    '"><span class="stat-icon">' + symbol(icon) + '</span><span><strong class="stat-number">' + counts[key] +
    '</strong><span class="stat-label">' + label + '</span></span></button>').join('');
  const options = activePage === 'expiring' ? [['all', 'All dates'], ['soon', 'Next 7 days'], ['expired', 'Past date']]
    : [['all', 'All items'], ['soon', 'Expiring soon'], ['expired', 'Past date'], ['review', 'Needs review']];
  $('#filters').hidden = activePage === 'history';
  $('#filters').innerHTML = options.map(([key, label]) => '<button type="button" class="filter" data-action="filter" data-value="' + key +
    '" aria-pressed="' + (key === filter) + '">' + label + '</button>').join('');
  $('#demo-bar').hidden = !demo || activePage === 'history';
  const source = activePage === 'history' ? records.filter(item => item.used)
    : activePage === 'expiring' ? current.filter(item => ['soon', 'expired'].includes(status(item))) : current;
  const query = $('#search').value.trim().toLowerCase();
  const found = source.filter(item => (filter === 'all' || status(item) === filter) && (!query ||
    [nameOf(item), item.category, item.rawItem, item.rawDate].some(text => text?.toLowerCase().includes(query))));
  const sort = $('#sort').value;
  found.sort((a, b) => {
    if (sort === 'name') return nameOf(a).localeCompare(nameOf(b));
    if (sort === 'newest' || activePage === 'history') return (b.used || b.created) - (a.used || a.created);
    if (a.demo && b.demo && activePage === 'overview' && filter === 'all') {
      return sample.indexOf(a) - sample.indexOf(b);
    }
    return (a.expiry || '9999').localeCompare(b.expiry || '9999') || b.created - a.created;
  });
  const shown = activePage === 'overview' && filter === 'all' && !query ? found.slice(0, 5) : found;
  const countText = (found.length === shown.length ? String(found.length) : shown.length + ' of ' + found.length) + (found.length === 1 ? ' item' : ' items');
  $('#inventory-count').textContent = countText + (counts.queued + counts.pending && ['overview', 'inventory'].includes(activePage)
    ? ' · ' + (counts.queued + counts.pending) + ' processing' : '');
  const drive = driveStatus();
  $('#drive-status').textContent = drive.unavailable ? 'My Drive is unavailable. Photos waiting to sync will stay on this device.'
    : drive.incomplete ? 'Some item folders in My Drive have missing or invalid records.' : '';
  $('#drive-status').hidden = !$('#drive-status').textContent;
  const list = $('#inventory-list');
  if (!shown.length) {
    const filtered = Boolean(query) || filter !== 'all';
    const title = filtered ? 'No matching items' : activePage === 'history' ? 'No used items yet'
      : activePage === 'expiring' ? 'Nothing due this week' : 'Your freezer is empty';
    const note = filtered ? 'Try another search or filter.' : activePage === 'history' ? 'Items appear here when you mark them as used.'
      : activePage === 'expiring' ? 'Items with an approaching label date will appear here.' : 'Scan an item and its expiry label to get started.';
    list.innerHTML = '<div class="empty">' + symbol(filtered ? 'search' : 'package') + '<strong>' + title + '</strong><p>' + note + '</p>' +
      (!filtered && ['overview', 'inventory'].includes(activePage) ? '<button type="button" class="button primary" data-action="scan">' + symbol('camera') + ' Scan first item</button>' : '') + '</div>';
    return;
  }
  list.innerHTML = shown.map(item => {
    const state = status(item);
    const [label] = states[state];
    const photo = previewURL(item);
    const date = item.expiry ? formatDate(item.expiry) : state === 'queued' ? 'Reading label…' : 'No date';
    return '<button type="button" class="item" data-action="details" data-value="' + esc(item.id) + '">' +
      '<span class="item-main">' + (photo ? '<img class="item-photo" alt="" src="' + esc(photo) + '">' :
        '<span class="item-photo">' + symbol('package') + '</span>') +
      '<span class="item-name"><strong>' + esc(nameOf(item)) + '</strong><small>' + esc(item.category || 'Uncategorised') +
      ' · ' + (item.quantity || 1) + ' pc' + (item.quantity > 1 ? 's' : '') + '</small></span></span>' +
      '<span class="item-date">' + symbol('calendar-days') + date + '<small>' + esc(item.dateType || 'Expiry date') + '</small></span>' +
      '<span class="badge ' + state + '">' + label + '</span><span class="item-arrow">' + symbol('chevron-right') + '</span></button>';
  }).join('');
}

function navigate(page, nextFilter = 'all') {
  if (location.hash !== '#' + page) {
    history.pushState({}, '', '#' + page);
    $('#search').value = '';
    window.scrollTo(0, 0);
  }
  filter = nextFilter;
  applyRoute();
}

function applyRoute() {
  const route = location.hash.slice(1);
  const next = Object.hasOwn(routes, route) ? route : 'overview';
  if (activePage !== next) $('#search').value = '';
  activePage = next;
  const [crumb, title, note] = routes[activePage];
  $('#crumb').textContent = crumb;
  for (const link of document.querySelectorAll('[data-nav]')) {
    if (link.dataset.nav === activePage) link.setAttribute('aria-current', 'page');
    else link.removeAttribute('aria-current');
  }
  $('#dashboard').hidden = activePage === 'settings';
  $('#settings').hidden = activePage !== 'settings';
  $('#hero').hidden = activePage !== 'overview';
  $('#metrics').hidden = activePage !== 'overview';
  $('#view-inventory').hidden = activePage !== 'overview';
  $('#page-title').textContent = title || '';
  $('#page-note').textContent = note || '';
  $('#today').textContent = formatToday();
  $('#inventory-title').textContent = activePage === 'history' ? 'Used items' : activePage === 'expiring' ? 'Items by label date' : 'Freezer inventory';
  render();
}

async function loadItems() {
  items = await allItems();
  for (const [key, url] of thumbnailURLs) {
    if (items.some(item => item.id === key)) continue;
    URL.revokeObjectURL(url);
    thumbnailURLs.delete(key);
  }
  render();
}

async function changed() {
  await loadItems();
  channel?.postMessage('changed');
}

function nextJob() {
  return items.filter(item => item.revision && item.state === 'queued').sort((a, b) => a.created - b.created)[0];
}

async function processQueue() {
  if (inQueue || !nextJob()) return;
  inQueue = true;
  async function run() {
    await loadItems();
    let job;
    while ((job = nextJob())) {
      queueMessage = 'Reading labels on this device…';
      renderScanner();
      const token = job.job;
      const commit = async (reading, pending) => {
        await updateItem(job.id, stored => stored.job === token && stored.state === 'queued'
          ? { ...stored, ...reading, state: pending ? 'queued' : 'ready' } : null);
        await changed();
      };
      try {
        await scanLabels(await jobPhotos(job), preferences, message => { queueMessage = message; renderScanner(); }, commit);
      } catch (error) {
        if ([401, 403, 409, 503].includes(error.status) || error instanceof TypeError) throw error;
        console.warn('Scan interrupted', error);
        await commit({ state: 'ready', review: true, reviewReason: 'Label reading was interrupted. Retry the scan or enter the details.' }, false);
      }
    }
  }
  try {
    if (navigator.locks?.request) await navigator.locks.request('froststock-recognition', run);
    else await run();
  } catch (error) {
    console.error(error);
    if (error.status === 401 || error.status === 403) storageError(error);
    else notification('The scan queue will retry. Your photos remain saved.');
  } finally {
    queueMessage = '';
    inQueue = false;
    renderScanner();
  }
}

async function uploadQueue() {
  if (syncing || !items.some(item => !item.revision)) return;
  syncing = true;
  try {
    const sync = async () => {
      while (items.some(item => !item.revision)) {
        await syncPending();
        await changed();
        void processQueue();
      }
    };
    if (navigator.locks?.request) await navigator.locks.request('froststock-drive', sync);
    else await sync();
  } catch (error) {
    storageError(error);
    if (ownerReady) {
      await loadItems().catch(storageError);
      void processQueue();
    }
  } finally {
    syncing = false;
  }
}

function displayScanError(message = '') {
  $('#scan-error').hidden = !message;
  $('#scan-error').textContent = message;
}

function renderScanner() {
  if (!scanner.open) return;
  const step = itemPhoto ? 2 : 1;
  $('#step-one').className = 'step ' + (step === 1 ? 'active' : 'done');
  $('#step-one').firstElementChild.innerHTML = step === 1 ? '1' : symbol('check');
  $('#step-two').className = 'step ' + (step === 2 ? 'active' : '');
  $('#viewfinder').dataset.step = step === 2 ? 'date' : 'item';
  $('#camera-heading').textContent = step === 2 ? 'Photograph the expiry date' : 'Photograph the item';
  $('#camera-tip').textContent = step === 2 ? 'Move close and keep the printed date in focus.' : 'Include the front of the pack or its label.';
  $('#camera-hint').textContent = $('#camera-tip').textContent;
  $('#camera-hint').hidden = !stream;
  $('#retake').hidden = step === 1;
  $('#retake').disabled = processingPhoto;
  const thumb = $('#first-thumb');
  thumb.hidden = !firstPreview;
  if (firstPreview) thumb.src = firstPreview;
  else thumb.removeAttribute('src');
  video.hidden = !stream;
  $('#camera-placeholder').hidden = Boolean(stream);
  $('#capture').querySelector('span').textContent = processingPhoto ? 'Saving photo…' : startingCamera ? 'Opening camera…'
    : stream ? (step === 1 ? 'Take item photo' : 'Take expiry photo') : 'Open camera';
  $('#capture').disabled = processingPhoto || startingCamera;
  scanner.querySelector('[data-action="upload"]').disabled = processingPhoto || startingCamera;
  scanner.querySelector('[data-action="close-scan"]').disabled = processingPhoto;
  $('#scan-session').textContent = (sessionCount ? sessionCount + ' captured. ' : '') +
    (queueMessage || (sessionCount ? 'Ready for the next item.' : 'Both photos are saved together.'));
}

function clearFirst() {
  itemPhoto = undefined;
  if (firstPreview) URL.revokeObjectURL(firstPreview);
  firstPreview = undefined;
}

function openScanner() {
  if (scanner.open || !ownerReady) return;
  sessionCount = 0;
  displayScanError();
  scanner.showModal();
  renderScanner();
}

function stopCamera() {
  cameraRequest++;
  startingCamera = false;
  stream?.getTracks().forEach(track => track.stop());
  stream = undefined;
  video.srcObject = null;
}

async function capture() {
  if (processingPhoto || startingCamera) return;
  displayScanError();
  if (stream) {
    if (video.readyState < 2 || !video.videoWidth) {
      displayScanError('The camera is starting. Try again in a moment.');
      return;
    }
    await acceptPhoto(() => photoBlob(video, video.videoWidth, video.videoHeight, Boolean(itemPhoto)));
    return;
  }
  if (!navigator.mediaDevices?.getUserMedia) {
    $('#camera-file').click();
    return;
  }
  const request = ++cameraRequest;
  startingCamera = true;
  renderScanner();
  try {
    const result = await navigator.mediaDevices.getUserMedia({ video: { facingMode: { ideal: 'environment' }, width: { ideal: 1920 }, height: { ideal: 1080 } }, audio: false });
    if (request !== cameraRequest || !scanner.open) {
      result.getTracks().forEach(track => track.stop());
      return;
    }
    stream = result;
    video.srcObject = stream;
    await video.play();
  } catch {
    if (request === cameraRequest) {
      stopCamera();
      displayScanError('Camera unavailable. Check permissions or upload a photo.');
    }
  } finally {
    if (request === cameraRequest) startingCamera = false;
    renderScanner();
  }
}

function photoBlob(image, width, height, closeup) {
  const max = closeup ? 2400 : 1600;
  const scale = Math.min(1, max / Math.max(width, height));
  const canvas = document.createElement('canvas');
  canvas.width = Math.max(1, Math.round(width * scale));
  canvas.height = Math.max(1, Math.round(height * scale));
  const ctx = canvas.getContext('2d', { alpha: false });
  ctx.fillStyle = '#fff';
  ctx.fillRect(0, 0, canvas.width, canvas.height);
  ctx.drawImage(image, 0, 0, canvas.width, canvas.height);
  return new Promise((resolve, reject) => {
    const encode = quality => canvas.toBlob(blob => {
      if (!blob) reject(new Error('Photo encoding failed'));
      else if (blob.size <= 3 * 1024 * 1024) resolve(blob);
      else if (quality > .5) encode(quality - .15);
      else reject(new Error('This photo is too large to save. Try a closer crop.'));
    }, .89);
    encode(.89);
  });
}

async function convertFile(file, closeup) {
  if (file.size > 40 * 1024 * 1024) throw new Error('Choose a photo under 40 MB.');
  const url = URL.createObjectURL(file);
  try {
    const image = new Image();
    image.src = url;
    try { await image.decode(); }
    catch { throw new Error('This photo couldn’t be opened. Try a JPEG, PNG or WebP.'); }
    return await photoBlob(image, image.naturalWidth, image.naturalHeight, closeup);
  } finally { URL.revokeObjectURL(url); }
}

async function acceptPhoto(makeBlob) {
  if (processingPhoto || !scanner.open) return;
  processingPhoto = true;
  displayScanError();
  renderScanner();
  try {
    const blob = await makeBlob();
    if (!itemPhoto) {
      itemPhoto = blob;
      firstPreview = URL.createObjectURL(blob);
    } else {
      const record = {
        id: id(), created: Date.now(), job: id(), state: 'queued', quantity: 1, itemPhoto, datePhoto: blob,
      };
      await saveItem(record);
      clearFirst();
      sessionCount++;
      if (demo) storeDemo(false);
      if (!persistentRequested) {
        persistentRequested = true;
        navigator.storage?.persist?.().catch(() => {});
      }
      items = items.concat(record);
      render();
      void uploadQueue();
    }
  } catch (error) {
    displayScanError(error?.message || 'Couldn’t save this photo. Please try again.');
  } finally {
    processingPhoto = false;
    renderScanner();
  }
}

function options(values, selected) {
  return values.map(value => '<option' + (value === selected ? ' selected' : '') + '>' + esc(value) + '</option>').join('');
}

function showDetail(key) {
  const item = items.find(record => record.id === key) || (demo && sample.find(record => record.id === key));
  if (!item) return;
  for (const url of detailURLs) URL.revokeObjectURL(url);
  detailURLs.clear();
  currentDetail = key;
  const title = nameOf(item);
  const preview = previewURL(item);
  const dateURL = item.revision ? itemPhotoURL(item.id, 'expiry.jpg') : item.datePhoto ? URL.createObjectURL(item.datePhoto) : '';
  if (dateURL.startsWith('blob:')) detailURLs.add(dateURL);
  details.innerHTML = '<div class="dialog-head"><div><p class="kicker">' + (item.demo ? 'SAMPLE ITEM' : 'ITEM DETAILS') +
    '</p><h2>' + esc(title) + '</h2></div><button type="button" class="icon-button" data-action="close-details" aria-label="Close item details">' + symbol('x') +
    '</button></div>' + (preview ? '<div class="detail-photos"><figure><img src="' + esc(preview) + '" alt="Food item"><figcaption>Item photo</figcaption></figure>' +
    (dateURL ? '<figure><img src="' + esc(dateURL) + '" alt="Expiry label"><figcaption>Expiry photo</figcaption></figure>' : '') + '</div>' : '') +
    (item.demo ? '<p class="detail-note">' + esc(item.category) + ' · ' + esc(item.dateType) + ' ' + formatDate(item.expiry) + '</p>'
      : '<form id="item-form"><label>Item name<input name="name" maxlength="100" value="' + esc(item.name) + '" placeholder="Enter the product name" required></label>' +
      '<div class="detail-columns"><label>Category<select name="category"><option value="">Uncategorised</option>' + options(categories, item.category) +
      '</select></label><label>Quantity<input name="quantity" type="number" min="1" max="999" value="' + (item.quantity || 1) + '" required></label></div>' +
      '<div class="detail-columns"><label>Label date<input name="expiry" type="date" min="1900-01-01" max="2199-12-31" value="' + esc(item.expiry) +
      '"></label><label>Label type<select name="dateType">' + options(dateTypes, item.dateType) + '</select></label></div>' +
      (item.state === 'queued' || item.review ? '<p class="detail-note">' + esc(item.state === 'queued'
        ? item.revision ? 'Reading labels. You can enter the details yourself and save them now.' : 'This pair will appear in My Drive after upload.'
        : item.reviewReason || 'Check the photos and confirm the label details. Leave the date blank if the label has no date.') + '</p>' : '') +
      '<details class="label-text"><summary>Extracted label text</summary><strong>Item photo</strong><pre>' + esc(item.rawItem || 'No text found') +
      '</pre><strong>Expiry photo</strong><pre>' + esc(item.rawDate || 'No text found') + '</pre>' +
      (item.visionItem || item.visionDate ? '<strong>Vision model</strong><pre>' + esc([item.visionItem, item.visionDate].filter(Boolean).join('\n')) + '</pre>' : '') +
      '</details><div class="detail-actions"><button class="button primary" type="submit"' + (item.revision ? '' : ' disabled') + '>' + symbol('check') + (item.review ? ' Confirm details' : ' Save details') +
      '</button><button class="button subtle" type="button" data-action="used"' + (item.revision ? '' : ' disabled') + '>' + (item.used ? 'Return to freezer' : 'Mark as used') + '</button>' +
      (item.state !== 'queued' && (item.revision || item.itemPhoto && item.datePhoto) ? '<button class="text-button" type="button" data-action="retry">Read labels again</button>' : '') +
      '<button class="text-button danger" type="button" data-action="delete">Delete</button></div></form>');
  if (!details.open) details.showModal();
}

function fileData(blob) {
  if (!blob) return null;
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(blob);
  });
}

async function exportBackup() {
  const saved = await allItems();
  if (!saved.length) { notification('There are no saved items to export.'); return; }
  const data = [];
  for (const item of saved) {
    const photos = item.revision ? await Promise.all(['item.jpg', 'expiry.jpg'].map(async name => {
      const response = await fetch(itemPhotoURL(item.id, name));
      if (!response.ok) throw new Error('Couldn’t export both photos for ' + nameOf(item));
      return response.blob();
    })) : [item.itemPhoto, item.datePhoto];
    data.push({ ...item, itemPhoto: await fileData(photos[0]), datePhoto: await fileData(photos[1]) });
  }
  const url = URL.createObjectURL(new Blob([JSON.stringify({ schema: 1, items: data })], { type: 'application/json' }));
  const link = document.createElement('a');
  link.href = url;
  link.download = 'froststock-backup-' + offsetDate(0) + '.json';
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 2000);
  notification('Inventory backup downloaded.');
}

function restorePhoto(data) {
  if (data == null) return null;
  if (typeof data !== 'string' || !/^data:image\/(?:jpeg|png|webp);base64,[a-z0-9+/]+=*$/i.test(data)) {
    throw new Error('The backup contains an invalid photo.');
  }
  const [type, encoded] = data.slice(5).split(';base64,');
  const binary = atob(encoded);
  return new Blob([Uint8Array.from(binary, char => char.charCodeAt(0))], { type });
}

async function importBackup(file) {
  if (file.size > 250 * 1024 * 1024) throw new Error('This backup exceeds 250 MB.');
  let backup;
  try { backup = JSON.parse(await file.text()); }
  catch { throw new Error('Choose a valid Froststock backup.'); }
  if (backup?.schema !== 1 || !Array.isArray(backup.items)) throw new Error('Choose a valid Froststock backup.');
  const existing = new Set((await allItems()).map(item => item.id));
  const restored = [];
  for (const record of backup.items) {
    if (!record || typeof record.id !== 'string' || !/^[\w-]{1,100}$/.test(record.id)) throw new Error('The backup contains an invalid item.');
    if (existing.has(record.id)) continue;
    const itemPhoto = restorePhoto(record.itemPhoto);
    const datePhoto = restorePhoto(record.datePhoto);
    const queued = record.state === 'queued' && itemPhoto && datePhoto;
    const text = (value, length) => typeof value === 'string' ? value.slice(0, length) : '';
    restored.push({
      id: record.id, created: Number.isFinite(record.created) ? record.created : Date.now(),
      used: Number.isFinite(record.used) ? record.used : undefined, state: queued ? 'queued' : 'ready', job: id(),
      itemPhoto, datePhoto, name: text(record.name, 100), category: categories.includes(record.category) ? record.category : '',
      quantity: Math.min(999, Math.max(1, Math.trunc(Number(record.quantity) || 1))),
      expiry: isDate(record.expiry) ? record.expiry : '', dateType: dateTypes.includes(record.dateType) ? record.dateType : 'Expiry date',
      review: Boolean(record.review), reviewReason: text(record.reviewReason, 300),
      rawItem: text(record.rawItem, 20000), rawDate: text(record.rawDate, 20000),
      visionItem: text(record.visionItem, 1000), visionDate: text(record.visionDate, 1000),
    });
    existing.add(record.id);
  }
  const candidates = restored.filter(item => item.itemPhoto && item.datePhoto);
  if (candidates.length !== restored.length) throw new Error('Each restored item needs both photos.');
  for (const item of candidates) {
    if (item.itemPhoto.type !== 'image/jpeg' || item.itemPhoto.size > 3 * 1024 * 1024) item.itemPhoto = await convertFile(item.itemPhoto, false);
    if (item.datePhoto.type !== 'image/jpeg' || item.datePhoto.size > 3 * 1024 * 1024) item.datePhoto = await convertFile(item.datePhoto, true);
  }
  await addItems(candidates);
  if (restored.length) storeDemo(false);
  await changed();
  void uploadQueue();
  void processQueue();
  notification(restored.length ? 'Queued ' + restored.length + (restored.length === 1 ? ' item for My Drive.' : ' items for My Drive.') : 'All items in this backup are already present.');
}

const actions = {
  scan: openScanner, 'close-scan': () => { if (!processingPhoto) scanner.close(); }, capture,
  upload: () => $('#photo-file').click(), retake: () => { clearFirst(); displayScanError(); renderScanner(); },
  filter: value => { filter = value; render(); },
  metric: value => navigate(value === 'soon' || value === 'expired' ? 'expiring' : 'inventory', value),
  inventory: () => navigate('inventory'), details: showDetail, 'close-details': () => details.close(),
  'dismiss-demo': () => storeDemo(false), 'toggle-demo': () => storeDemo(!demo),
  export: exportBackup, import: () => $('#backup-file').click(),
  used: async () => {
    const result = await updateItem(currentDetail, item => ({ ...item, used: item.used ? undefined : Date.now() }));
    details.close();
    await changed();
    if (result) notification(result.used ? 'Item moved to history.' : 'Item returned to the freezer.');
  },
  retry: async () => {
    await updateItem(currentDetail, item => ({ ...item, state: 'queued', job: id(), reviewReason: '' }));
    details.close();
    await changed();
    void processQueue();
  },
  delete: async () => {
    if (!window.confirm('Delete this item and its photos?')) return;
    await removeItem(currentDetail);
    details.close();
    await changed();
    notification('Item deleted.');
  },
};

document.addEventListener('click', async event => {
  const nav = event.target.closest('[data-nav], .brand');
  if (nav && !event.ctrlKey && !event.metaKey && !event.shiftKey && !event.altKey) {
    event.preventDefault();
    navigate(nav.dataset.nav || 'overview');
    return;
  }
  const button = event.target.closest('[data-action]');
  if (!button || !Object.hasOwn(actions, button.dataset.action)) return;
  try { await actions[button.dataset.action](button.dataset.value); }
  catch (error) { notification(error?.message || 'This action could not be completed.'); }
});

document.addEventListener('submit', async event => {
  if (event.target.id !== 'item-form') return;
  event.preventDefault();
  const fields = new FormData(event.target);
  const name = fields.get('name').trim();
  if (!name) { event.target.elements.name.focus(); return; }
  try {
    await updateItem(currentDetail, item => ({
      ...item, name, category: fields.get('category'), quantity: Number(fields.get('quantity')),
      expiry: fields.get('expiry'), dateType: fields.get('dateType'), state: 'ready',
      review: false, reviewReason: '', job: id(),
    }));
    details.close();
    await changed();
    notification('Item details saved.');
  } catch { notification('The item could not be saved.'); }
});

for (const input of [$('#photo-file'), $('#camera-file')]) {
  input.addEventListener('change', () => {
    const file = input.files[0];
    input.value = '';
    if (file) void acceptPhoto(() => convertFile(file, Boolean(itemPhoto)));
  });
}
$('#backup-file').addEventListener('change', async event => {
  const file = event.target.files[0];
  event.target.value = '';
  try { if (file) await importBackup(file); }
  catch (error) { notification(error?.message || 'The backup could not be restored.'); }
});
$('#search').addEventListener('input', render);
$('#sort').addEventListener('change', render);
$('#date-order').value = preferences.dateOrder;
$('#date-order').addEventListener('change', event => {
  preferences.dateOrder = event.target.value;
  localStorage.setItem('froststock-date-order', preferences.dateOrder);
});
$('#vision-enabled').checked = preferences.vision;
$('#vision-enabled').addEventListener('change', async event => {
  if (event.target.checked && !await navigator.gpu?.requestAdapter().catch(() => null)) {
    event.target.checked = false;
    notification('This browser does not support the WebGPU vision model.');
    return;
  }
  preferences.vision = event.target.checked;
  localStorage.setItem('froststock-vision', preferences.vision ? 'on' : 'off');
});
scanner.addEventListener('cancel', event => { if (processingPhoto) event.preventDefault(); });
scanner.addEventListener('close', stopCamera);
details.addEventListener('close', () => {
  currentDetail = undefined;
  for (const url of detailURLs) URL.revokeObjectURL(url);
  detailURLs.clear();
});
window.addEventListener('hashchange', () => { filter = 'all'; applyRoute(); });
async function refresh() {
  if (!ownerReady) return;
  await loadItems();
  void uploadQueue();
  void processQueue();
}

function storageError(error) {
  if (error.status === 401 || error.status === 403) {
    ownerReady = false;
    if (scanner.open && !processingPhoto) scanner.close();
    if (details.open) details.close();
    $('#main').hidden = true;
    $('#session-screen').hidden = false;
    return;
  }
  notification(error.message || 'My Drive is unavailable. Queued photos will be retried.');
}

window.addEventListener('pageshow', () => refresh().catch(storageError));
document.addEventListener('visibilitychange', () => { if (!document.hidden) refresh().catch(storageError); });
channel?.addEventListener('message', () => refresh().catch(storageError));
setInterval(() => { if (!document.hidden) refresh().catch(storageError); }, 20000);
applyRoute();
connect().then(() => {
  ownerReady = true;
  $('#main').hidden = false;
  $('#session-screen').hidden = true;
  return refresh();
}).catch(storageError);
