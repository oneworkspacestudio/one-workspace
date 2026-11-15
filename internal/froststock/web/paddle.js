const assets = new URL('./assets/paddle/', import.meta.url);
let worker;
let ready;
let requestId = 0;

// The bundled PaddleOCR.js worker runs both tiny models with ONNX Runtime WASM.
function send(type, payload, transfer = []) {
  const id = ++requestId;
  const engine = worker;
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => finish(reject, new Error('Local expiry OCR timed out')), 16000);
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
