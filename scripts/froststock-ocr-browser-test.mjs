import assert from 'node:assert/strict';
import { chromium } from 'playwright';

// Serve internal/froststock/web on loopback, with no real inventory or account.
const origin = process.env.FROSTSTOCK_TEST_URL || 'http://127.0.0.1:37133/';
const browser = await chromium.launch({ executablePath: process.env.FROSTSTOCK_CHROMIUM });
try {
  const page = await browser.newPage({ viewport: { width: 390, height: 844 } });
  await page.route(origin, route => route.fulfill({
    contentType: 'text/html',
    body: '<!doctype html><script src="./assets/ocr/tesseract.min.js"></script>',
  }));
  await page.goto(origin);
  const result = await page.evaluate(async () => {
    const { scanLabels } = await import('./recognition.js');
    const canvas = (width, height) => Object.assign(document.createElement('canvas'), { width, height });
    const photo = (c, quality = .66) => new Promise(resolve => c.toBlob(resolve, 'image/jpeg', quality));
    const draw = (c, bg, fg, lines) => {
      const ctx = c.getContext('2d');
      ctx.fillStyle = bg;
      ctx.fillRect(0, 0, c.width, c.height);
      for (const [text, size, x, y, color] of lines) {
        ctx.font = 'bold ' + size + 'px Arial';
        ctx.fillStyle = color || fg;
        ctx.fillText(text, x, y);
      }
    };
    const product = canvas(980, 650);
    draw(product, '#415a6d', '#fff', [
      ['MARKET', 66, 98, 130], ['ATLANTIC SALMON', 78, 90, 278], ['FILLETS', 78, 90, 365],
      ['COOK FROM FROZEN', 32, 96, 507], ['SERVING SUGGESTION', 24, 96, 574],
    ]);
    const stamp = canvas(720, 380);
    draw(stamp, '#acafaa', '#747870', [
      ['BATCH 3821', 24, 55, 74], ['BEST BEFORE', 44, 170, 173],
      ['22 NOV 2027', 54, 172, 237], ['07:32', 22, 188, 297],
    ]);
    const sideways = canvas(stamp.height, stamp.width);
    const ctx = sideways.getContext('2d');
    ctx.translate(sideways.width, 0);
    ctx.rotate(Math.PI / 2);
    ctx.drawImage(stamp, 0, 0);
    const mixed = canvas(870, 290);
    draw(mixed, '#fff', '#202020', [
      ['PACKED 24.11.2026', 42, 100, 95], ['USE BY 12.01.2027', 59, 100, 220],
    ]);
    const incomplete = canvas(720, 330);
    draw(incomplete, '#fff', '#222', [
      ['BEST BEFORE', 54, 100, 130], ['22 NOV', 56, 100, 225],
    ]);
    const itemPhoto = await photo(product);
    async function run(c) {
      let scanned;
      await scanLabels({ itemPhoto, datePhoto: await photo(c) }, { dateOrder: 'dmy', vision: false }, () => {}, value => { scanned = value; });
      return scanned;
    }
    return [await run(sideways), await run(mixed), await run(incomplete)];
  });
  assert.equal(result[0].name, 'Atlantic Salmon Fillets', JSON.stringify(result[0]));
  assert.equal(result[0].expiry, '2027-11-22', JSON.stringify(result[0]));
  assert.equal(result[0].dateType, 'Best before');
  assert.equal(result[1].expiry, '2027-01-12', JSON.stringify(result[1]));
  assert.equal(result[1].dateType, 'Use by');
  assert.equal(result[2].expiry, '', JSON.stringify(result[2]));
  assert.equal(result[2].review, true);
  console.log('Local OCR passed wrapped package title, sideways low-contrast date, competing packed/use-by dates and missing year.');
} finally {
  await browser.close();
}
