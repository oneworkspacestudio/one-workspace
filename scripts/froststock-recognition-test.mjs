import assert from 'node:assert/strict';
import { test } from 'node:test';
import { extractDate, readTitle } from '../internal/froststock/web/dates.js';

test('finds the date tied to the expiry label and ignores packing dates', () => {
  assert.equal(extractDate('MFG 04/06/26 EXP 20/08/27').value, '2027-08-20');
  assert.equal(extractDate('MFG 04/06/26 EXP\n20/08/27').value, '2027-08-20');
  assert.equal(extractDate('PACKED ON 22 NOV 2027').value, '');
  assert.equal(extractDate('BEST BEFORE: 22 NOV 2027\nPACKED 18 NOV 2026').type, 'Best before');
  assert.equal(extractDate('USE BY 13/01/2027').value, '2027-01-13');
});

test('reads compact and faint-stamp transcriptions without inventing missing fields', () => {
  assert.equal(extractDate('1806.2029').value, '2029-06-18');
  assert.equal(extractDate('1806.2029').repaired, true);
  assert.equal(extractDate('18062029').value, '');
  assert.deepEqual(extractDate('19.06.2027\n18.06.2029').candidates.map(date => date.value), ['2027-06-19','2029-06-18']);
  assert.equal(extractDate('BBE 22·11·2027').value, '2027-11-22');
  assert.equal(extractDate('EXP241227').value, '2027-12-24');
  assert.equal(extractDate('L241227').value, '');
  assert.equal(extractDate('BEST BEFORE\nO8 DEC 2O27').value, '2027-12-08');
  assert.equal(extractDate('BEST BEFORE\nO8 DEC 2O27').repaired, true);
  assert.equal(extractDate('USE BY 13/O1/27').value, '2027-01-13');
  assert.equal(extractDate('EXP 2028-02-29').value, '2028-02-29');
  assert.equal(extractDate('EXP 29/02/27').value, '');
  assert.equal(extractDate('BEST BEFORE 11/2027').value, '');
  assert.equal(extractDate('USE BY 12 DEC').value, '');
  assert.equal(extractDate('BBE 05/06/27', 'mdy').value, '2027-05-06');
  assert.equal(extractDate('BBE 05/06/27', 'dmy').value, '2027-06-05');
  assert.equal(extractDate('USE BY 12/01/27 15/01/27').ambiguous, true);
});

test('identifies prominent product titles ahead of a brand and directions', () => {
  const line = (text, y0, y1, x1, x0 = 100) => ({ text, bbox: { x0, y0, x1, y1 }, confidence: 88 });
  const lines = [
    line('MARKET', 4, 37, 200), line('ATLANTIC SALMON', 80, 125, 480),
    line('FILLETS', 132, 172, 335), line('COOK FROM FROZEN', 240, 289, 445),
  ];
  assert.equal(readTitle('', [{ paragraphs: [{ lines }] }]).name, 'Atlantic Salmon Fillets');
  assert.equal(readTitle('MARKET\nGREEN PEAS 800 g\nKEEP FROZEN').name, 'Green Peas');
  assert.equal(readTitle('NUTRITION\nBEST BEFORE 12 JAN 27').name, '');
});

test('keeps the food on a slanted package with a large brand and a short cut on the next line', () => {
  const line = (text, x0, x1, y0, y1) => ({ text, confidence: 95, bbox: { x0, x1, y0, y1 } });
  const lines = [
    line('Brand', 40, 390, 410, 555), line('FROZEN TURKEY', 68, 365, 560, 625),
    line('THIGH', 147, 275, 600, 657), line('PIECES', 160, 269, 668, 702),
  ];
  assert.equal(readTitle('', [{ paragraphs: [{ lines }] }]).name, 'Frozen Turkey Thigh');
  assert.equal(readTitle('FROZEN\nBrand\nTURKEY WING').name, 'Turkey Wing');
  assert.equal(readTitle('Icken Part)').name, '');
});
