import { extractDate, extractName, readTitle } from './dates.js';
import { labelVariants } from './label-image.js';
import { cropLabelPhoto } from './expiry-crop.js';
import { paddleLayout, readPaddleLabels } from './paddle.js';

const assetPath = new URL('./assets/ocr/', import.meta.url).href;
let worker;
let vision;
let visionFailed = false;
let paddleRetryAfter = 0;

async function recognize(image, mode, blocks = false) {
  let timer;
  let engine;
  try {
    const read = async () => {
      worker ??= window.Tesseract.createWorker('eng', 1, {
        workerPath: assetPath + 'worker.min.js', langPath: assetPath, corePath: assetPath,
        workerBlobURL: false,
      });
      engine = await worker;
      const { data } = await engine.recognize(image, { tessedit_pageseg_mode: mode }, { text: true, blocks });
      return { text: data.text.trim(), confidence: data.confidence, blocks: data.blocks };
    };
    return await Promise.race([read(), new Promise((_, reject) => {
      timer = setTimeout(() => reject(new Error('Label reading timed out')), 22000);
    })]);
  } catch (error) {
    const failed = worker;
    worker = undefined;
    if (engine) void engine.terminate().catch(() => {});
    else void failed?.then(instance => instance.terminate()).catch(() => {});
    throw error;
  } finally {
    clearTimeout(timer);
  }
}

const emptyRead = { text: '', confidence: 0 };

async function readPaddle(blob) {
  if (Date.now() < paddleRetryAfter) return [];
  try { return await readPaddleLabels(blob); }
  catch (error) {
    paddleRetryAfter = Date.now() + 60000;
    console.warn('Local OCR model unavailable', error);
    return [];
  }
}

async function readProduct(blob) {
  let reading = emptyRead;
  try { reading = await recognize(blob, '11', true); }
  catch (error) { console.warn('First title scan failed', error); }
  let title = readTitle(reading.text, reading.blocks, reading.confidence);
  if (title.food && title.confidence >= 75 && title.name.includes(' ')) return { reading, title };

  const quality = value => value.name ? value.confidence + (value.food ? 32 : 0) + (value.name.includes(' ') ? 12 : 0) : 0;
  const paddle = paddleLayout(await readPaddle(blob));
  if (paddle.length) {
    const text = paddle.map(line => line.text).join('\n');
    const alternative = readTitle(text, [{ paragraphs: [{ lines: paddle }] }]);
    if (quality(alternative) > quality(title)) {
      title = alternative;
      reading = { text, confidence: title.confidence };
    }
    if (title.food && title.confidence >= 82 && title.name.includes(' ')) return { reading, title };
  }
  try {
    for await (const { canvas, mode } of labelVariants(blob, 'item')) {
      const alternate = await recognize(canvas, mode, true);
      const found = readTitle(alternate.text, alternate.blocks, alternate.confidence);
      if (quality(found) > quality(title)) {
        title = found;
        reading = alternate;
      }
      if (title.food && title.confidence >= 82 && title.name.includes(' ')) break;
    }
  } catch (error) {
    console.warn('Additional title scans stopped', error);
  }
  return { reading, title };
}

function stampDate(text, order) {
  const joined = text.trim().match(/^([0-9oilsbz]{2})([0-9oilsbz]{2})((?:19|20|21)\d{2})$/i);
  const date = extractDate(joined ? joined.slice(1).join('.') : text, order);
  return joined ? { ...date, repaired: true } : date;
}

function dateStamp(lines, order) {
  const rows = lines.filter(line => line.score >= .81 && /^[0-9oilsbz ./-]+$/i.test(line.text.trim()) &&
    stampDate(line.text, order).value);
  if (rows.length !== 2) return;
  const [first, second] = paddleLayout(rows).sort((a, b) => a.bbox.y0 - b.bbox.y0);
  if (!first || !second) return;
  const a = first.bbox, b = second.bbox;
  const height = Math.max(a.y1 - a.y0, b.y1 - b.y0);
  const width = Math.min(a.x1 - a.x0, b.x1 - b.x0);
  if (width <= 0 || b.y0 - a.y0 < height * .6 || b.y0 - a.y0 > height * 2.5 ||
      Math.min(a.x1, b.x1) - Math.max(a.x0, b.x0) < width * .78) return;
  const earlier = stampDate(first.text, order);
  const later = stampDate(second.text, order);
  const days = (Date.parse(later.value) - Date.parse(earlier.value)) / 86400000;
  if (days >= 30 && days <= 1830) return { ...later, suggested: true, candidates: [earlier, later] };
}

async function readExpiry(blob, order, cropped) {
  const reads = [];
  function record(reading, minimum = 76) {
    const date = extractDate(reading.text, order);
    const score = (date.explicit ? 45 : 0) + (date.repaired ? 0 : 23) + reading.confidence;
    reads.push({ reading, date, score });
    return date.value && date.explicit && !date.repaired && reading.confidence >= minimum &&
      !reads.some(hit => hit.date.value && hit.date.value !== date.value);
  }
  async function read(image, mode) {
    return record(await recognize(image, mode));
  }
  let certain = false;
  try { certain = await read(blob, cropped ? '6' : '11'); }
  catch (error) { console.warn('First expiry scan failed', error); }
  if (!certain) {
    const lines = await readPaddle(blob);
    if (lines.length) {
      const text = lines.map(line => line.text).join('\n');
      const confidence = 100 * lines.reduce((sum, line) => sum + line.score, 0) / lines.length;
      const reading = { text, confidence };
      certain = record(reading, 89);
      const parsed = reads.at(-1).date;
      const hasExplicit = reads.some(hit => hit.date.candidates.some(candidate => candidate.explicit));
      if (parsed.candidates.length && !hasExplicit) {
        const suggested = dateStamp(lines, order);
        if (suggested) return { reading, date: suggested };
      }
      if (parsed.ambiguous && !hasExplicit) {
        const clearDates = lines.filter(line => line.score >= .85).flatMap(line => extractDate(line.text, order).candidates);
        const values = new Set(clearDates.map(candidate => candidate.value));
        if (values.size >= 2) return { reading, date: parsed };
      }
    }
  }
  if (!certain) {
    try {
      for await (const { canvas, mode } of labelVariants(blob, 'date', cropped)) {
        if (await read(canvas, mode)) break;
      }
    } catch (error) {
      console.warn('Additional date scans stopped', error);
    }
  }
  const valid = reads.filter(hit => hit.date.value).sort((a, b) => b.score - a.score);
  const best = valid[0];
  if (best) {
    const conflicting = valid.filter(hit => hit.score >= best.score - 18 && hit.date.value !== best.date.value);
    if (!conflicting.length) return best;
    const choices = [...new Map([best, ...conflicting].map(hit => [hit.date.value, hit.date])).values()];
    return { reading: best.reading, date: { ...best.date, value: '', ambiguous: true, candidates: choices } };
  }
  return reads.find(hit => hit.date.ambiguous) || reads.sort((a, b) => b.reading.confidence - a.reading.confidence)[0];
}

function readWithVision(blob, purpose) {
  return new Promise((resolve, reject) => {
    vision ??= new Worker(new URL('./vision.js', import.meta.url), { type: 'module' });
    const timeout = setTimeout(() => fail(new Error('Vision model timed out')), 120000);
    function fail(error) {
      clearTimeout(timeout);
      vision.terminate();
      vision = undefined;
      reject(error);
    }
    vision.onmessage = ({ data }) => {
      if (data.error) return fail(new Error(data.error));
      clearTimeout(timeout);
      resolve(data.text);
    };
    vision.onerror = event => fail(new Error(event.message));
    vision.postMessage({ blob, purpose });
  });
}

function deepSeekDate(reading, order) {
  const text = reading.rawDate?.trim().slice(0, 4900) || '';
  const evidence = extractDate(text, order);
  const stamp = reading.expiryText?.trim().slice(0, 80);
  if (!stamp || !text.toLowerCase().includes(stamp.toLowerCase())) return null;
  const suggested = extractDate(stamp, order);
  if (!suggested.value || evidence.value !== suggested.value) return null;
  return { text, date: evidence };
}

export async function scanLabels(item, preferences, report, commit, readCloud) {
  let dateResult = emptyRead, itemResult = emptyRead;
  let date = extractDate('', preferences.dateOrder);
  let title = { name: '', confidence: 0 };
  let readFailed = false;
  let datePhoto = item.datePhoto;
  try {
    datePhoto = await cropLabelPhoto(item.datePhoto, item.dateCrop);
  } catch (error) {
    readFailed = true;
    console.warn('Expiry photo crop failed', error);
  }
  const localReading = 'Reading labels on this device…';
  let cloud, cloudDate, cloudName, cloudError;
  if (readCloud && !readFailed) {
    try {
      cloud = await readCloud(item.dateOnly ? null : item.itemPhoto, datePhoto, report);
      if (cloud) {
        cloudDate = deepSeekDate(cloud, preferences.dateOrder);
        if (!item.dateOnly) cloudName = extractName(cloud.name || '');
      }
    } catch (error) {
      const message = error?.message;
      if (message?.startsWith('DeepSeek ')) cloudError = message;
    }
  }
  if (!cloudDate || (!item.dateOnly && !cloudName)) {
    report(cloud ? 'Checking unreadable fields on this device…'
      : cloudError ? cloudError + ' ' + localReading : localReading);
  }
  if (cloudDate) {
    date = cloudDate.date;
    dateResult = { text: cloudDate.text };
  } else {
    try {
      const result = await readExpiry(datePhoto, preferences.dateOrder, datePhoto !== item.datePhoto);
      if (!result) readFailed = true;
      else ({ reading: dateResult, date } = result);
    } catch (error) { readFailed = true; console.warn('Expiry OCR failed', error); }
  }
  if (!item.dateOnly && !cloudName) {
    try { ({ reading: itemResult, title } = await readProduct(item.itemPhoto)); }
    catch (error) { readFailed = true; console.warn('Item OCR failed', error); }
  }
  const name = cloudName || (title.confidence >= 38 ? title.name : '');
  const cloudItemText = item.dateOnly ? '' : cloud?.rawItem?.trim().slice(0, 4900) || '';
  const otherDate = extractDate(itemResult.text || cloudItemText, preferences.dateOrder);
  const fromItem = !item.dateCrop && !date.value && !date.ambiguous && otherDate.value && otherDate.explicit;
  const expiry = fromItem ? otherDate : date;
  const inferredName = Boolean(cloudName && cloud?.nameInferred);
  const localNameReview = !name || title.confidence < 65 || (!title.food && !name.includes(' '));
  const nameReview = item.dateOnly ? !item.name : cloudName ? inferredName : localNameReview;
  const review = Boolean(cloudDate || !expiry.value || fromItem || expiry.suggested || expiry.repaired ||
    dateResult.confidence < 75 || nameReview);
  const result = {
    name: name || item.name || '',
    expiry: expiry.value || (date.ambiguous && item.review ? '' : item.expiry || ''),
    dateType: expiry.value ? expiry.type : item.dateType || 'Expiry date',
    rawItem: (cloudName ? cloudItemText : itemResult.text) || cloudItemText || item.rawItem || '',
    rawDate: dateResult.text || cloud?.rawDate?.trim().slice(0, 4900) || (item.dateCrop ? '' : item.rawDate || ''),
    visionDate: '',
    dateCandidates: expiry.candidates.length > 1 ? expiry.candidates.slice(0, 6).map(candidate => candidate.value) : [],
    review,
    reviewReason: expiry.suggested ? 'Two dates were read from the stamp. The later date is suggested. Confirm it against the pack.'
      : date.ambiguous ? 'Several dates were found. Choose the expiry date printed on the pack.'
      : fromItem ? 'The date was read from the item photo. Check it against the expiry photo.'
      : expiry.repaired ? 'Some digits were unclear. Confirm the suggested date against the photo.'
      : cloudDate && inferredName ? 'DeepSeek suggested the food name and expiry date. Confirm both against the photos.'
      : cloudDate ? 'DeepSeek suggested this expiry date. Confirm it against the photo.'
      : inferredName ? 'DeepSeek suggested a food name from the image. Confirm it against the photo.'
      : readFailed ? 'Label reading was interrupted. Read the labels again or enter the details.'
      : !expiry.value ? 'The expiry date could not be read. Crop around the date or enter it.'
      : nameReview && !name ? 'The food title could not be read. Check the package photo or enter the name.' : '',
    state: 'ready',
  };
  const needsTitle = !item.dateOnly && !name;
  const tryVision = preferences.vision && navigator.gpu && !visionFailed &&
    ((!expiry.value && !date.ambiguous) || needsTitle);
  await commit(result, tryVision);
  if (!tryVision) return;
  report('Reading difficult labels with the local vision model…');
  try {
    if (!expiry.value && !date.ambiguous) {
      const text = await readWithVision(datePhoto, 'date');
      const suggestion = extractDate(text, preferences.dateOrder);
      result.visionDate = text;
      if (suggestion.value && !item.expiry) {
        result.expiry = suggestion.value;
        result.dateType = suggestion.type;
        result.reviewReason = 'The vision model suggested this date. Confirm it against the photo.';
      }
    }
    if (needsTitle) {
      const text = await readWithVision(item.itemPhoto, 'item');
      result.visionItem = text;
      if (!/\b(?:unknown|illegible|cannot read)\b/i.test(text)) result.name = extractName(text) || result.name;
    }
  } catch (error) {
    visionFailed = true;
    result.reviewReason ||= 'The local vision model was unavailable. Check the photos and enter the label details.';
    console.warn('Local vision model unavailable', error);
  }
  await commit(result, false);
}
