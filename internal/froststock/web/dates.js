const months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec'];
const monthName = '(jan(?:uary)?|feb(?:ruary)?|mar(?:ch)?|apr(?:il)?|may|jun(?:e)?|jul(?:y)?|aug(?:ust)?|sep(?:t(?:ember)?)?|oct(?:ober)?|nov(?:ember)?|dec(?:ember)?)';
const expiryToken = 'best[ .-]*before(?:[ .-]*end)?|use[ .-]*by|sell[ .-]*by|exp(?:iry|iration|ires)?(?:[ .-]*date)?|b[ .-]*b[ .-]*e|bb|ub';
const expiryWords = new RegExp('\\b(?:' + expiryToken + ')\\b', 'i');
const productionWords = /\b(?:packed|produced|manufactur(?:ed|ing)|mfd|mfg|frozen|production|batch|lot)\b/i;
const labelTokens = new RegExp('(?<![a-z0-9])(?:' + expiryToken + '|' + productionWords.source + ')(?=[^a-z]|$)', 'gi');
const ocrDigits = value => value.toUpperCase().replace(/[OILSBZ]/g, letter => ({ O: '0', I: '1', L: '1', S: '5', B: '8', Z: '2' })[letter]);
const dayFirst = new RegExp('(?<![a-z0-9])([0-9oilsbz]{1,2})(?:st|nd|rd|th)?[\\s./-]*' + monthName + '[\\s.,/-]*([0-9oilsbz]{4}|[0-9oilsbz]{2})(?![a-z0-9])', 'gi');
const monthFirst = new RegExp('(?<![a-z0-9])' + monthName + '[\\s./-]*([0-9oilsbz]{1,2})(?:st|nd|rd|th)?[\\s.,/-]+([0-9oilsbz]{4}|[0-9oilsbz]{2})(?![a-z0-9])', 'gi');
const dateFormat = new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' });
const fullDate = new Intl.DateTimeFormat('en-GB', { weekday: 'long', day: 'numeric', month: 'long' });

function toISO(year, month, day) {
  if (!(year >= 1900 && year <= 2199)) return '';
  const date = new Date(Date.UTC(year, month - 1, day));
  return date.getUTCMonth() + 1 === month && date.getUTCDate() === day ? date.toISOString().slice(0, 10) : '';
}

function yearNumber(text) {
  return Number(text) + (text.length === 2 ? 2000 : 0);
}

export function extractDate(text, order = 'dmy') {
  const normalized = text.normalize('NFKC').replace(/[–—−]/g, '-').replace(/[·•]/g, '.').replace(/\r/g, '').trim();
  const lines = normalized.split('\n');
  const matches = new Map();
  const add = (value, context, fields) => {
    if (!value || context === 'production') return;
    const raw = fields.join('');
    const repaired = /[a-z]/i.test(raw);
    if (repaired && (!/\d/.test(raw) || raw.match(/[a-z]/gi).length > 2)) return;
    const candidate = { value, type: context || 'Expiry date', explicit: !!context, repaired };
    const previous = matches.get(value);
    if (!previous || (candidate.explicit && !previous.explicit) ||
        (!candidate.repaired && previous.repaired && candidate.explicit === previous.explicit)) {
      matches.set(value, candidate);
    }
  };
  function numeric(a, b, year, context) {
    const first = +ocrDigits(a), second = +ocrDigits(b);
    const monthFirst = second > 12 || (first <= 12 && order === 'mdy');
    add(toISO(yearNumber(ocrDigits(year)), monthFirst ? first : second, monthFirst ? second : first), context, [a, b, year]);
  }
  let precedingContext = '';
  for (const line of lines) {
    if (!line.trim()) continue;
    const tokens = [...line.matchAll(labelTokens)];
    let context = precedingContext;
    let start = 0;
    for (const token of [...tokens, { index: line.length, 0: '' }]) {
      const segment = line.slice(start, token.index);
      for (const [, a, b, c] of segment.matchAll(/(?<![a-z0-9/.-])([0-9oilsbz]{1,4})[ \t]*(?:[/.,-][ \t]*|[ \t]+)([0-9oilsbz]{1,2})[ \t]*(?:[/.,-][ \t]*|[ \t]+)([0-9oilsbz]{1,4})(?![a-z0-9/.-])/gi)) {
        if (a.length === 4 && c.length <= 2) add(toISO(+ocrDigits(a), +ocrDigits(b), +ocrDigits(c)), context, [a, b, c]);
        else if (a.length <= 2 && [2, 4].includes(c.length)) numeric(a, b, c, context);
      }
      if (context && context !== 'production') {
        const compact = segment.trim().match(/^[:.]?[ \t]*(\d{2})(\d{2})(\d{2}|\d{4})(?:[ \t]|$)/);
        if (compact) numeric(...compact.slice(1), context);
      }
      for (const [, day, month, year] of segment.matchAll(dayFirst)) add(toISO(yearNumber(ocrDigits(year)), months.indexOf(month.toLowerCase().slice(0, 3)) + 1, +ocrDigits(day)), context, [day, year]);
      for (const [, month, day, year] of segment.matchAll(monthFirst)) add(toISO(yearNumber(ocrDigits(year)), months.indexOf(month.toLowerCase().slice(0, 3)) + 1, +ocrDigits(day)), context, [day, year]);
      start = token.index + token[0].length;
      if (productionWords.test(token[0])) context = 'production';
      else if (token[0]) {
        const label = token[0].replace(/[ .-]/g, '').toLowerCase();
        context = label.startsWith('best') || label.startsWith('bb') ? 'Best before' : label === 'useby' || label === 'ub'
          ? 'Use by' : label === 'sellby' ? 'Sell by' : 'Expiry date';
      }
    }
    const lastLabel = tokens.at(-1);
    precedingContext = lastLabel && /^[\s:.-]*(?:on[\s:.-]*)?$/i.test(line.slice(lastLabel.index + lastLabel[0].length)) ? context : '';
  }
  const explicit = [...matches.values()].filter(candidate => candidate.explicit);
  const candidates = explicit.length ? explicit : [...matches.values()];
  return candidates.length === 1 ? { ...candidates[0], ambiguous: false } :
    { value: '', type: 'Expiry date', explicit: false, repaired: false, ambiguous: candidates.length > 1 };
}

function titleLine(text) {
  return text.normalize('NFKC').replace(/\s+/g, ' ').replace(/^[^\p{L}\p{N}]+|[^\p{L}\p{N})]+$/gu, '')
    .replace(/\s+\d+(?:[.,]\d+)?\s*(?:g|kg|ml|oz|lb)\b.*$/i, '').trim();
}

function validTitle(line) {
  const letters = line.match(/\p{L}/gu)?.length ?? 0;
  return letters >= 3 && letters / line.length > .5 && line.length <= 100 && !expiryWords.test(line) &&
    !/^(?:keep\b|net\b|ingredients?\b|nutrition\b|cooking\b|storage\b|frozen food\b|product of\b|warning\b|lot\b|batch\b|suitable for\b|serving\b|cook (?:from|before)\b|made with\b|recycl|no (?:artificial|added)\b|serves?\s+\d|per (?:portion|serving|\d)|family pack$)/i.test(line);
}

function formatTitle(text) {
  const name = text.slice(0, 100);
  return name === name.toUpperCase() ? name.toLowerCase().replace(/(^|[\s-])\p{L}/gu, c => c.toUpperCase()) : name;
}

export function readTitle(text, blocks, confidence = 0) {
  const lines = (blocks ?? []).flatMap(block => block.paragraphs ?? []).flatMap(paragraph => paragraph.lines ?? []);
  const candidates = [];
  for (const line of lines) {
    const title = titleLine(line.text);
    if (!line.bbox || !validTitle(title)) continue;
    const height = Math.max(1, line.bbox.y1 - line.bbox.y0);
    const words = title.split(' ').length;
    const score = height * (words === 1 ? .56 : 1 + Math.min(words, 5) * .08) * (.38 + line.confidence / 130);
    candidates.push({ title, height, score, bbox: line.bbox, confidence: line.confidence });
  }
  candidates.sort((a, b) => b.score - a.score);
  const best = candidates[0];
  if (!best) {
    const plain = text.split(/\n/).map(titleLine).filter(validTitle);
    return { name: plain.length ? formatTitle(plain.find(line => line.includes(' ')) || plain[0]) : '', confidence, score: 0 };
  }
  // Join two close, similarly sized and horizontally overlapping title lines.
  const adjacent = candidates.slice(1).find(line => {
    const [a, b] = [best, line].sort((first, second) => first.bbox.y0 - second.bbox.y0);
    const gap = b.bbox.y0 - a.bbox.y1;
    const overlap = Math.min(a.bbox.x1, b.bbox.x1) - Math.max(a.bbox.x0, b.bbox.x0);
    return line.height >= best.height * .65 && line.height <= best.height * 1.35 && line.confidence >= 50 &&
      gap >= -2 && gap <= best.height * .65 && overlap >= Math.min(a.bbox.x1 - a.bbox.x0, b.bbox.x1 - b.bbox.x0) * .55;
  });
  const title = adjacent ? [best, adjacent].sort((a, b) => a.bbox.y0 - b.bbox.y0).map(line => line.title).join(' ') : best.title;
  return { name: formatTitle(title), confidence: best.confidence, score: best.score };
}

export const extractName = text => readTitle(text).name;

export function localISO(date = new Date()) {
  return toISO(date.getFullYear(), date.getMonth() + 1, date.getDate());
}

export function offsetDate(days, date = new Date()) {
  return localISO(new Date(date.getFullYear(), date.getMonth(), date.getDate() + days));
}

export function daysFromToday(value, now = new Date()) {
  return Math.round((Date.parse(value) - Date.parse(localISO(now))) / 86400000);
}
export const formatDate = iso => iso ? dateFormat.format(new Date(iso)) : 'No date';
export const formatToday = () => fullDate.format(new Date());
export const isDate = value => typeof value === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(value) &&
  value === toISO(...value.split('-').map(Number));
