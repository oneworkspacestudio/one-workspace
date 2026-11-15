function canvasFor(width, height) {
  const canvas = document.createElement('canvas');
  canvas.width = Math.round(width);
  canvas.height = Math.round(height);
  return canvas;
}

export async function* labelVariants(blob, kind, cropped = false) {
  const image = await createImageBitmap(blob);
  try {
    const date = kind === 'date';
    const size = date ? 2000 : 1700;
    const scale = Math.min(2.25, size / Math.max(image.width, image.height));
    const source = canvasFor(image.width * scale, image.height * scale);
    const ctx = source.getContext('2d', { willReadFrequently: true });
    ctx.drawImage(image, 0, 0, source.width, source.height);
    enhance(ctx, source.width, source.height);
    yield { canvas: source, mode: '11' };
    if (date) {
      if (!cropped) yield { canvas: cropCentre(source), mode: '6' };
      else yield { canvas: source, mode: '7' };
      for (const angle of [90, 270, 180]) yield { canvas: rotate(source, angle), mode: '11' };
    } else {
      for (const angle of [90, 270, 180]) yield { canvas: rotate(source, angle), mode: '11' };
    }
  } finally {
    image.close();
  }
}

function enhance(ctx, width, height) {
  const frame = ctx.getImageData(0, 0, width, height);
  const { data } = frame;
  const histogram = new Uint32Array(256);
  for (let index = 0; index < data.length; index += 4) {
    const grey = Math.round(.299 * data[index] + .587 * data[index + 1] + .114 * data[index + 2]);
    histogram[grey]++;
    data[index] = grey;
  }
  let low = 0, high = 255, seen = 0;
  const tail = width * height * .02;
  while (seen + histogram[low] < tail) seen += histogram[low++];
  seen = 0;
  while (seen + histogram[high] < tail) seen += histogram[high--];
  const stretch = 255 / Math.max(high - low, 25);
  for (let index = 0; index < data.length; index += 4) {
    const value = (data[index] - low) * stretch;
    data[index] = data[index + 1] = data[index + 2] = value;
  }
  ctx.putImageData(frame, 0, 0);
}

function cropCentre(source) {
  const canvas = canvasFor(source.width * .86, source.height * .58);
  canvas.getContext('2d').drawImage(source, (source.width - canvas.width) / 2, (source.height - canvas.height) / 2,
    canvas.width, canvas.height, 0, 0, canvas.width, canvas.height);
  return canvas;
}

function rotate(source, angle) {
  const sideways = angle !== 180;
  const canvas = canvasFor(sideways ? source.height : source.width, sideways ? source.width : source.height);
  const ctx = canvas.getContext('2d');
  ctx.translate(canvas.width / 2, canvas.height / 2);
  ctx.rotate(angle * Math.PI / 180);
  ctx.drawImage(source, -source.width / 2, -source.height / 2);
  return canvas;
}
