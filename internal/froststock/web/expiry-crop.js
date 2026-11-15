const limit = (value, low, high) => Math.max(low, Math.min(high, value));
const minimum = .015;
const initial = { left: .12, top: .28, right: .88, bottom: .72 };

export function validCrop(crop) {
  return crop && [crop.left, crop.top, crop.right, crop.bottom].every(Number.isFinite) &&
    crop.left >= 0 && crop.top >= 0 && crop.right <= 1 && crop.bottom <= 1 &&
    crop.right - crop.left >= .01 && crop.bottom - crop.top >= .01;
}

function move(crop, dx, dy) {
  dx = limit(dx, -crop.left, 1 - crop.right);
  dy = limit(dy, -crop.top, 1 - crop.bottom);
  return { left: crop.left + dx, top: crop.top + dy, right: crop.right + dx, bottom: crop.bottom + dy };
}

function resize(crop, corner, dx, dy) {
  const result = { ...crop };
  if (corner.includes('l')) result.left = limit(crop.left + dx, 0, crop.right - minimum);
  if (corner.includes('r')) result.right = limit(crop.right + dx, crop.left + minimum, 1);
  if (corner.includes('t')) result.top = limit(crop.top + dy, 0, crop.bottom - minimum);
  if (corner.includes('b')) result.bottom = limit(crop.bottom + dy, crop.top + minimum, 1);
  return result;
}

export async function chooseExpiryCrop(blob, saved) {
  const dialog = document.querySelector('#expiry-crop');
  if (dialog.open || dialog.dataset.loading) return undefined;
  dialog.dataset.loading = 'true';
  const stage = dialog.querySelector('.crop-image');
  const photo = stage.querySelector('img');
  const box = stage.querySelector('.crop-box');
  const url = URL.createObjectURL(blob);
  const events = new AbortController();
  const options = { signal: events.signal };
  let crop = validCrop(saved) ? { ...saved } : { ...initial };
  let drag;
  let result;
  function display() {
    box.style.inset = [crop.top, 1 - crop.right, 1 - crop.bottom, crop.left]
      .map(value => value * 100 + '%').join(' ');
  }
  function position(event) {
    const rect = stage.getBoundingClientRect();
    return { x: limit((event.clientX - rect.left) / rect.width, 0, 1), y: limit((event.clientY - rect.top) / rect.height, 0, 1) };
  }
  try {
    photo.src = url;
    await photo.decode();
    display();
    stage.addEventListener('pointerdown', event => {
      if (drag || event.button !== 0) return;
      event.preventDefault();
      stage.setPointerCapture(event.pointerId);
      drag = { ...position(event), id: event.pointerId, crop: { ...crop },
        kind: event.target.dataset.corner || (box.contains(event.target) ? 'move' : 'new') };
    }, options);
    stage.addEventListener('pointermove', event => {
      if (drag?.id !== event.pointerId) return;
      const { x, y } = position(event);
      const dx = x - drag.x, dy = y - drag.y;
      if (drag.kind === 'move') crop = move(drag.crop, dx, dy);
      else if (drag.kind !== 'new') crop = resize(drag.crop, drag.kind, dx, dy);
      else if (Math.abs(dx) >= minimum && Math.abs(dy) >= minimum) {
        crop = { left: Math.min(x, drag.x), top: Math.min(y, drag.y),
          right: Math.max(x, drag.x), bottom: Math.max(y, drag.y) };
      }
      display();
    }, options);
    stage.addEventListener('lostpointercapture', () => { drag = undefined; }, options);
    const directions = { ArrowLeft: [-1, 0], ArrowRight: [1, 0], ArrowUp: [0, -1], ArrowDown: [0, 1] };
    box.addEventListener('keydown', event => {
      const direction = directions[event.key];
      if (!direction) return;
      event.preventDefault();
      const [dx, dy] = direction.map(value => value * (event.altKey ? .05 : .01));
      crop = event.shiftKey ? resize(crop, 'br', dx, dy) : move(crop, dx, dy);
      display();
    }, options);
    dialog.addEventListener('click', event => {
      const action = event.target.closest('[data-crop]')?.dataset.crop;
      if (!action) return;
      result = action === 'apply' ? crop : action === 'full' ? null : undefined;
      dialog.close();
    }, options);
    const closed = new Promise(resolve => dialog.addEventListener('close', resolve, { once: true, ...options }));
    dialog.showModal();
    await closed;
    return result;
  } finally {
    events.abort();
    delete dialog.dataset.loading;
    photo.removeAttribute('src');
    URL.revokeObjectURL(url);
  }
}

export async function cropLabelPhoto(blob, crop) {
  if (!validCrop(crop)) return blob;
  const image = await createImageBitmap(blob);
  try {
    const left = Math.floor(crop.left * image.width), top = Math.floor(crop.top * image.height);
    const width = Math.ceil(crop.right * image.width) - left;
    const height = Math.ceil(crop.bottom * image.height) - top;
    const scale = Math.min(2, 2000 / Math.max(width, height));
    const canvas = document.createElement('canvas');
    canvas.width = Math.round(width * scale);
    canvas.height = Math.round(height * scale);
    canvas.getContext('2d', { alpha: false }).drawImage(image, left, top, width, height, 0, 0, canvas.width, canvas.height);
    return await new Promise((resolve, reject) => canvas.toBlob(value => {
      if (value) resolve(value);
      else reject(new Error('Could not prepare the selected part of the photo'));
    }, 'image/png'));
  } finally {
    image.close();
  }
}
