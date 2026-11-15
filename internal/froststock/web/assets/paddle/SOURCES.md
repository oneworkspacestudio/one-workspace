# Froststock date OCR

- `paddleocr-worker.js` is based on the worker at `dist/assets/worker-entry-C9UNuyOJ.js` in [PaddleOCR.js 0.4.2](https://www.npmjs.com/package/@paddleocr/paddleocr-js/v/0.4.2). [Package source](https://github.com/PaddlePaddle/PaddleOCR/tree/e5046169b225bcdfbe25d45b4e809ff0f1a69c2c/paddleocr-js).
- `det.tar` and `rec.tar` contain the official `inference.onnx` and `inference.yml` exports. [Detection checkpoint](https://huggingface.co/PaddlePaddle/PP-OCRv6_tiny_det_onnx/tree/2ba1506c0380b8f0b03dd142459aac66d4421f6c) and [recognition checkpoint](https://huggingface.co/PaddlePaddle/PP-OCRv6_tiny_rec_onnx/tree/2612ab37152ae0a677521bae4e1e3d4fb4cf7c30).
- The two `ort-wasm-simd-threaded.jsep` files are unmodified from the `dist` directory of [onnxruntime-web 1.24.3](https://www.npmjs.com/package/onnxruntime-web/v/1.24.3). [Runtime source](https://github.com/microsoft/onnxruntime/tree/v1.24.3).
- The worker bundles [OpenCV.js](https://github.com/TechStark/opencv-js), [js-yaml](https://github.com/nodeca/js-yaml) and [Clipper](https://www.npmjs.com/package/clipper-lib/v/6.4.2). Licence texts and upstream attribution are included here.

Froststock modifies the worker’s four OpenCV binding helpers
(`createNamedFunction`, `embind__requireFunction`, `craftInvokerFunction`, and
`__emval_get_method_caller`). They use closures in place of evaluating JavaScript
strings. The OpenCV WebAssembly and ONNX modules are unchanged.

The worker and runtime are served by Froststock. They only load the two model archives from this same directory.
