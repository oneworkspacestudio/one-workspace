import { extractDate, extractName, readTitle } from './dates.js';
import { labelVariants } from './label-image.js';
import { cropLabelPhoto } from './expiry-crop.js';
import { readPaddleLabels } from './paddle.js';

const assetPath = new URL('./assets/ocr/', import.meta.url).href;
let worker;
let vision;
let visionFailed = false;
let paddleFailed = false;

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

async function readProduct(blob) {
  let reading = await recognize(blob, '11', true);
  let title = readTitle(reading.text, reading.blocks, reading.confidence);
  if (title.confidence >= 75 && title.name.includes(' ')) return { reading, title };
  try {
    for await (const { canvas, mode } of labelVariants(blob, 'item')) {
      const alternate = await recognize(canvas, mode, true);
      const found = readTitle(alternate.text, alternate.blocks, alternate.confidence);
      const quality = value => (value.name ? value.confidence : 0) + (value.name.includes(' ') ? 12 : 0);
      if (quality(found) > quality(title)) {
        title = found;
        reading = alternate;
      }
      if (title.confidence >= 82 && title.name.includes(' ')) break;
    }
  } catch (error) {
    console.warn('Additional title scans stopped', error);
  }
  return { reading, title };
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
  if (!certain && !paddleFailed) {
    try {
      const lines = await readPaddleLabels(blob);
      if (lines.length) {
        const text = lines.map(line => line.text).join('\n');
        const confidence = 100 * lines.reduce((sum, line) => sum + line.score, 0) / lines.length;
        certain = record({ text, confidence }, 89);
      }
    } catch (error) {
      paddleFailed = true;
      console.warn('Local expiry model unavailable', error);
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
    const conflict = valid.some(hit => hit.date.value !== best.date.value && hit.score >= best.score - 18);
    if (!conflict) return best;
    return { reading: best.reading, date: { ...best.date, value: '', ambiguous: true } };
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

export async function scanLabels(item, preferences, report, commit) {
  let dateResult = emptyRead, itemResult = emptyRead;
  let date = extractDate('', preferences.dateOrder);
  let title = { name: '', confidence: 0 };
  let readFailed = false;
  let datePhoto = item.datePhoto;
  try {
    datePhoto = await cropLabelPhoto(item.datePhoto, item.dateCrop);
    const result = await readExpiry(datePhoto, preferences.dateOrder, datePhoto !== item.datePhoto);
    if (!result) readFailed = true;
    else ({ reading: dateResult, date } = result);
  } catch (error) { readFailed = true; console.warn('Expiry OCR failed', error); }
  if (!item.dateOnly) {
    try { ({ reading: itemResult, title } = await readProduct(item.itemPhoto)); }
    catch (error) { readFailed = true; console.warn('Item OCR failed', error); }
  }
  const name = title.confidence >= 38 ? title.name : '';
  const otherDate = extractDate(itemResult.text, preferences.dateOrder);
  const fromItem = !item.dateCrop && !date.value && !date.ambiguous && otherDate.value && otherDate.explicit;
  const expiry = fromItem ? otherDate : date;
  const nameReview = item.dateOnly ? !item.name : !name || title.confidence < 65;
  const review = Boolean(!expiry.value || fromItem || expiry.repaired || dateResult.confidence < 75 || nameReview);
  const result = {
    name: name || item.name || '',
    expiry: expiry.value || item.expiry || '', dateType: expiry.value ? expiry.type : item.dateType || 'Expiry date',
    rawItem: itemResult.text || item.rawItem || '', rawDate: dateResult.text || (item.dateCrop ? '' : item.rawDate || ''),
    visionDate: '',
    review,
    reviewReason: date.ambiguous ? 'Several dates were found. Crop around the expiry date or enter it.'
      : fromItem ? 'The date was read from the item photo. Check it against the expiry photo.'
      : expiry.repaired ? 'Some digits were unclear. Confirm the suggested date against the photo.'
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
