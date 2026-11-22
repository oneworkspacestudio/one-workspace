export type Source = { title: string; url: string; date?: string; text?: string };
export type Section = { heading: string; body: string; refs: number[] };
export type Article = {
  id: string; title: string; summary: string; topic: string; date: string;
  sections: Section[]; sources: Source[];
};
export type DisplayArticle = Article & { label: string; byline: string; tone: string; saved: boolean };
export type Job = { id: string; kind: string; query: string; status: string; sources: Source[] | null; error?: string };
export type Library = { jobs: Job[]; articles: Article[]; saved: string[]; deepseek: boolean; browser: boolean; browserReady: boolean };
export type Device = { read(): string | null; write(data: string): boolean; pick(): void };

export const packLimit = 2_000_000;
const invalid = () => Error('This file is not a valid AI Fieldnotes reading pack.');
const bounded = (value: unknown, max: number): value is string => typeof value === 'string' && value.length > 0 && value.length <= max;
const articleID = /^(?:[a-f0-9]{32}|loop-engineering|five-stages|inclusive-interface)$/;

export function safeLink(value: unknown): value is string {
  if (!bounded(value, 2048)) return false;
  try {
    const link = new URL(value);
    return link.protocol === 'https:' && !!link.hostname && !link.username && !link.password;
  } catch {
    return false;
  }
}

export function readPack(raw: string): Article[] {
  if (!raw || new TextEncoder().encode(raw).length > packLimit) throw invalid();
  let value: { format?: string; articles?: Article[] };
  try { value = JSON.parse(raw); } catch { throw invalid(); }
  if (value?.format !== 'ai-fieldnotes/1' || !Array.isArray(value.articles) || value.articles.length > 100) throw invalid();
  const ids = new Set();
  return value.articles.map(article => {
    if (!article || !bounded(article.id, 40) || !articleID.test(article.id) || ids.has(article.id) ||
        !bounded(article.title, 240) || !bounded(article.summary, 1200) || !bounded(article.topic, 300) ||
        !bounded(article.date, 10) || !/^\d{4}-\d{2}-\d{2}$/.test(article.date) ||
        !Array.isArray(article.sources) || article.sources.length < 1 || article.sources.length > 8 ||
        !Array.isArray(article.sections) || article.sections.length < 1 || article.sections.length > 8) throw invalid();
    ids.add(article.id);
    const sources = article.sources.map(entry => {
      if (!entry || !bounded(entry.title, 240) || !safeLink(entry.url)) throw invalid();
      return {title: entry.title, url: entry.url, date: bounded(entry.date, 40) ? entry.date : undefined};
    });
    const sections = article.sections.map(part => {
      if (!part || !bounded(part.heading, 220) || !bounded(part.body, 7000) || !Array.isArray(part.refs) ||
          part.refs.length > 8 || part.refs.some(i => !Number.isInteger(i) || i < 1 || i > sources.length)) throw invalid();
      return {heading: part.heading, body: part.body, refs: [...new Set(part.refs)]};
    });
    return {id: article.id, title: article.title, summary: article.summary, topic: article.topic, date: article.date, sections, sources};
  });
}

export function encodePack(articles: Article[]): string {
  const raw = JSON.stringify({format: 'ai-fieldnotes/1', articles});
  if (new TextEncoder().encode(raw).length > packLimit || articles.length > 100) {
    throw Error('This pack is too large. Choose up to 100 briefings and keep the pack under 2 MB.');
  }
  return raw;
}
