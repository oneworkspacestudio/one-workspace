const assets = new URL('./assets/paddle/', import.meta.url);
let worker;
let ready;
let requestId = 0;

// The bundled PaddleOCR.js worker runs both tiny models with ONNX Runtime WASM.
function send(type, payload, transfer = []) {
  const id = ++requestId;
  const engine = worker;
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => finish(reject, new Error('Local expiry OCR timed out')), type === 'init' ? 45000 : 16000);
    function finish(settle, value) {
      clearTimeout(timer);
      engine.removeEventListener('message', receive);
      engine.removeEventListener('error', failed);
      settle(value);
    }
    function receive({ data }) {
      if (data.kind !== 'worker-transport-response' || data.requestId !== id) return;
      if (data.status === 'success') finish(resolve, data.payload);
      else finish(reject, new Error(data.error?.message || 'Local expiry OCR failed'));
    }
    const failed = event => finish(reject, new Error(event.message || 'Local expiry OCR failed'));
    engine.addEventListener('message', receive);
    engine.addEventListener('error', failed);
    try { engine.postMessage({ kind: 'worker-transport-request', type, requestId: id, payload }, transfer); }
    catch (error) { finish(reject, error); }
  });
}

async function initialize() {
  worker = new Worker(new URL('paddleocr-worker.js', assets), { type: 'module' });
  return send('init', { options: {
    pipelineConfig: {
      modelSelection: { textDetectionModelName: 'PP-OCRv6_tiny_det', textRecognitionModelName: 'PP-OCRv6_tiny_rec' },
      assets: {
        det: { url: new URL('det.tar', assets).href },
        rec: { url: new URL('rec.tar', assets).href },
      },
      textDetectionBatchSize: 1, textRecognitionBatchSize: 1, pipelineBatchSize: 1,
      runtimeDefaults: {}, warnings: [],
    },
    ortOptions: { backend: 'wasm', numThreads: 1, disableWasmProxy: true, wasmPaths: assets.href },
  } });
}

export async function readPaddleLabels(image) {
  if (typeof OffscreenCanvas !== 'function') return [];
  try {
    ready ??= initialize();
    await ready;
    const imageBitmap = await createImageBitmap(image);
    const [result] = await send('predict', {
      sources: [{ kind: 'imageBitmap', imageBitmap }],
      params: { textRecScoreThresh: .55, textDetLimitSideLen: 1280, textDetLimitType: 'max' },
    }, [imageBitmap]);
    return result.items;
  } catch (error) {
    worker?.terminate();
    worker = undefined;
    ready = undefined;
    throw error;
  }
}

// Convert slanted detection boxes to the label's reading coordinates.
export function paddleLayout(lines) {
  const width = line => Math.hypot(line.poly[1][0] - line.poly[0][0], line.poly[1][1] - line.poly[0][1]);
  const visible = lines.filter(line => line.poly?.length === 4);
  const reference = visible.reduce((best, line) => !best || width(line) * line.score > width(best) * best.score ? line : best, null);
  if (!reference) return [];
  const [[x0, y0], [x1, y1]] = reference.poly;
  const angle = Math.atan2(y1 - y0, x1 - x0);
  const cos = Math.cos(angle), sin = Math.sin(angle);
  return visible.map(line => {
    const x = line.poly.map(([px, py]) => px * cos + py * sin);
    const y = line.poly.map(([px, py]) => py * cos - px * sin);
    return { text: line.text, confidence: line.score * 100,
      bbox: { x0: Math.min(...x), x1: Math.max(...x), y0: Math.min(...y), y1: Math.max(...y) } };
  });
}
