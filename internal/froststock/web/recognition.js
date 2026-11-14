import { extractDate, extractName } from './dates.js';

const assetPath = new URL('./assets/ocr/', import.meta.url).href;
let worker;
let vision;
let visionFailed = false;

async function recognize(blob, expiry = false) {
  worker ??= window.Tesseract.createWorker('eng', 1, {
    workerPath: assetPath + 'worker.min.js', langPath: assetPath, corePath: assetPath,
    workerBlobURL: false, errorHandler: () => {},
  }).catch(error => { worker = undefined; throw error; });
  const engine = await worker;
  const result = await engine.recognize(blob, { tessedit_pageseg_mode: expiry ? '11' : '3' });
  return { text: result.data.text.trim(), confidence: result.data.confidence };
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
  let dateResult = { text: '', confidence: 0 };
  let itemResult = { text: '', confidence: 0 };
  try {
    dateResult = await recognize(item.datePhoto, true);
    if (!extractDate(dateResult.text, preferences.dateOrder).value) {
      const block = await (await worker).recognize(item.datePhoto, { tessedit_pageseg_mode: '6' });
      if (extractDate(block.data.text, preferences.dateOrder).value || block.data.confidence > dateResult.confidence) {
        dateResult = { text: block.data.text.trim(), confidence: block.data.confidence };
      }
    }
  } catch (error) { console.warn('Expiry OCR failed', error); }
  try { itemResult = await recognize(item.itemPhoto); }
  catch (error) { console.warn('Item OCR failed', error); }
  const date = extractDate(dateResult.text, preferences.dateOrder);
  const name = extractName(itemResult.text);
  const result = {
    name: itemResult.confidence >= 40 ? name : '',
    expiry: date.value, dateType: date.type, rawItem: itemResult.text, rawDate: dateResult.text,
    review: !date.value || dateResult.confidence < 75 || !name || itemResult.confidence < 60,
    reviewReason: date.ambiguous ? 'Several dates were found. Check the expiry photo.' : '',
    state: 'ready',
  };
  const tryVision = preferences.vision && navigator.gpu && !visionFailed &&
    ((!date.value && !date.ambiguous) || !result.name);
  await commit(result, tryVision);
  if (!tryVision) return;
  report('Reading difficult labels with the local vision model…');
  try {
    if (!date.value && !date.ambiguous) {
      const text = await readWithVision(item.datePhoto, 'date');
      const suggestion = extractDate(text, preferences.dateOrder);
      result.visionDate = text;
      if (suggestion.value) {
        result.expiry = suggestion.value;
        result.dateType = suggestion.type;
        result.reviewReason = 'The vision model suggested this date. Confirm it against the photo.';
      }
    }
    if (!result.name) {
      const text = await readWithVision(item.itemPhoto, 'item');
      result.visionItem = text;
      if (!/\b(?:unknown|illegible|cannot read)\b/i.test(text)) result.name = extractName(text);
    }
  } catch (error) {
    visionFailed = true;
    result.reviewReason ||= 'The local vision model was unavailable. Check the photos and enter the label details.';
    console.warn('Local vision model unavailable', error);
  }
  await commit(result, false);
}
