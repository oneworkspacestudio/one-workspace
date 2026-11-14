const modelId = 'HuggingFaceTB/SmolVLM-256M-Instruct';
let engine;

async function initialize() {
  const { AutoModelForVision2Seq, AutoProcessor, RawImage, env } =
    await import('./assets/vision/transformers.web.min.js');
  env.allowLocalModels = false;
  env.backends.onnx.wasm.numThreads = 1;
  env.backends.onnx.wasm.wasmPaths = new URL('./assets/vision/', import.meta.url).href;
  const processor = await AutoProcessor.from_pretrained(modelId);
  const model = await AutoModelForVision2Seq.from_pretrained(modelId, {
    device: 'webgpu',
    dtype: { embed_tokens: 'fp16', vision_encoder: 'q4', decoder_model_merged: 'q4' },
  });
  return { processor, model, RawImage };
}

self.onmessage = async ({ data: { blob, purpose } }) => {
  let inputs;
  let generated;
  try {
    engine ??= initialize().catch(error => { engine = undefined; throw error; });
    const { processor, model, RawImage } = await engine;
    const image = await RawImage.fromBlob(blob);
    const instruction = purpose === 'date'
      ? 'Transcribe the exact printed expiry or best before date in the image. Reply with the printed text only. If illegible reply UNKNOWN.'
      : 'Read the food product name printed on the package. Reply with the exact name only. If illegible reply UNKNOWN.';
    const prompt = processor.apply_chat_template([{ role: 'user', content: [{ type: 'image' }, { type: 'text', text: instruction }] }], { add_generation_prompt: true });
    inputs = await processor(prompt, [image], { do_image_splitting: false });
    generated = await model.generate({ ...inputs, max_new_tokens: 32, do_sample: false });
    const answer = processor.tokenizer.decode(generated.tolist()[0].slice(inputs.input_ids.dims.at(-1)), { skip_special_tokens: true });
    self.postMessage({ text: answer.trim() });
  } catch (error) {
    self.postMessage({ error: error.message });
  } finally {
    if (inputs) for (const value of Object.values(inputs)) value?.dispose?.();
    generated?.dispose();
  }
};
