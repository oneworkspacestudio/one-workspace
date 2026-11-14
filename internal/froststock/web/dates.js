const months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec'];
const monthName = '(jan(?:uary)?|feb(?:ruary)?|mar(?:ch)?|apr(?:il)?|may|jun(?:e)?|jul(?:y)?|aug(?:ust)?|sep(?:t(?:ember)?)?|oct(?:ober)?|nov(?:ember)?|dec(?:ember)?)';
const expiryWords = /\b(?:best\s*before(?:\s*end)?|use\s*by|sell\s*by|exp(?:iry|iration|ires)?(?:\s*date)?|bbe|bb)\b/i;
const productionWords = /\b(?:packed|produced|manufactur(?:ed|ing)|mfd|mfg|frozen|production)\b/i;
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
  const normalized = text.replace(/[–—]/g, '-').replace(/\r/g, '').trim();
  const lines = normalized.split('\n');
  const matches = new Map();
  const add = (value, context) => {
    if (!value || context === 'production') return;
    if (context === 'expiry' || !matches.has(value)) matches.set(value, context);
  };
  let precedingContext = '';
  for (const line of lines) {
    if (!line.trim()) continue;
    const tokens = [...line.matchAll(/\b(?:best\s*before(?:\s*end)?|use\s*by|sell\s*by|exp(?:iry|iration|ires)?(?:\s*date)?|bbe|bb|packed|produced|manufactur(?:ed|ing)|mfd|mfg|frozen|production)\b/gi)];
    let context = precedingContext;
    let start = 0;
    for (const token of [...tokens, { index: line.length, 0: '' }]) {
      const segment = line.slice(start, token.index);
      for (const [, year, month, day] of segment.matchAll(/(?<![\d/.-])(\d{4})[-/.](\d{1,2})[-/.](\d{1,2})(?![\d/.-])/g)) {
        add(toISO(+year, +month, +day), context);
      }
      for (const [, a, b, year] of segment.matchAll(/(?<![\d/.-])(\d{1,2})\s*[-/. ]\s*(\d{1,2})\s*[-/. ]\s*(\d{2}|\d{4})(?![\d/.-])/g)) {
        const monthFirst = order === 'mdy' && Number(b) <= 31 && Number(a) <= 12 || Number(b) > 12;
        add(toISO(yearNumber(year), +(monthFirst ? a : b), +(monthFirst ? b : a)), context);
      }
      const dayFirst = new RegExp('(?<![a-z0-9])(\\d{1,2})(?:st|nd|rd|th)?[\\s./-]*' + monthName + '[\\s.,/-]*(\\d{4}|\\d{2})(?![a-z0-9])', 'gi');
      const monthFirst = new RegExp('(?<![a-z0-9])' + monthName + '[\\s./-]*(\\d{1,2})(?:st|nd|rd|th)?[\\s.,/-]+(\\d{4}|\\d{2})(?![a-z0-9])', 'gi');
      for (const [, day, month, year] of segment.matchAll(dayFirst)) add(toISO(yearNumber(year), months.indexOf(month.toLowerCase().slice(0, 3)) + 1, +day), context);
      for (const [, month, day, year] of segment.matchAll(monthFirst)) add(toISO(yearNumber(year), months.indexOf(month.toLowerCase().slice(0, 3)) + 1, +day), context);
      start = token.index + token[0].length;
      if (expiryWords.test(token[0])) context = 'expiry';
      else if (productionWords.test(token[0])) context = 'production';
    }
    precedingContext = tokens.length && !/\d/.test(line.slice(tokens.at(-1).index)) ? context : '';
  }
  const explicit = [...matches].filter(([, context]) => context === 'expiry');
  const candidates = explicit.length ? explicit : [...matches];
  const type = /\buse\s*by\b/i.test(text) ? 'Use by' : /\b(best\s*before|bbe|bb)\b/i.test(text) ? 'Best before' : /\bsell\s*by\b/i.test(text) ? 'Sell by' : 'Expiry date';
  return { value: candidates.length === 1 ? candidates[0][0] : '', type, ambiguous: candidates.length > 1 };
}

export function extractName(text) {
  const lines = text.split(/\n/).map(line => line.replace(/\s+/g, ' ').replace(/^[^\p{L}\p{N}]+|[^\p{L}\p{N})]+$/gu, '').trim());
  const candidate = lines.find(line => {
    const letters = line.match(/\p{L}/gu)?.length ?? 0;
    return letters >= 3 && letters / line.length > .5 && !expiryWords.test(line) &&
      !/^(?:net\s|ingredients?\b|nutrition\b|keep frozen\b|keep refrigerated\b|cooking\b|storage\b|frozen food\b|product of\b|warning\b|lot\b)/i.test(line);
  });
  if (!candidate) return '';
  const name = candidate.slice(0, 80);
  return name === name.toUpperCase() ? name.toLowerCase().replace(/(^|\s)\p{L}/gu, c => c.toUpperCase()) : name;
}

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
