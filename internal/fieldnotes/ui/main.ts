import {createElement as h, useEffect, useRef, useState} from 'react';
import {createRoot} from 'react-dom/client';
import {Reader} from './generated/Reader';
import {Controls} from './components';
import {encodePack, readPack, type Article, type Device, type DisplayArticle, type Library} from './model';
import starter from './starter.json';
import './style.css';

declare global {
  interface Window {
    FieldnotesDevice?: Device;
    fieldnotesImport?: (value: string | null) => void;
    fieldnotesBack?: () => boolean;
  }
}

const device = window.FieldnotesDevice;
const offline = !!device;
const empty: Library = {articles: [], saved: [], jobs: [], deepseek: false, browser: false, browserReady: false};
const topics = [
  {title: 'Loop engineering', query: 'loop engineering agent feedback loops UI software engineering'},
  {title: 'AI maturity', query: 'AI engineering maturity capabilities software engineering DORA'},
  {title: 'Design systems', query: 'AI design system component governance frontend engineering case studies'},
  {title: 'Standards & accessibility', query: 'WCAG AI generated user interfaces accessibility W3C frontend testing'},
  {title: 'Evaluations', query: 'agent evaluation frontend visual regression accessibility UI engineering'},
  {title: 'What comes next', query: 'AI-native engineering operating model harness engineering UI teams'},
];
const releaseLink = 'https://github.com/oneworkspacestudio/one-workspace/blob/main/docs/AI_FIELDNOTES.md';
const readHash = () => location.hash.startsWith('#read-') ? location.hash.slice(6) : '';
const dateFormat = new Intl.DateTimeFormat(undefined, {month: 'short', day: 'numeric', year: 'numeric', timeZone: 'UTC'});
const starterIDs = new Set(starter.map(a => a.id));

function initialDeviceArticles(): Article[] {
  if (!device) return [];
  try { return readPack(device.read() || encodePack([])); } catch { return []; }
}
function App() {
  const [folder, setFolder] = useState('all');
  const [panel, setPanel] = useState('');
  const [notice, setNotice] = useState('');
  const [fields, setFields] = useState<Record<string, string>>({provider: 'deepseek'});
  const [reading, setReading] = useState(readHash);
  const [owner, setOwner] = useState(false);
  const [workspace, setWorkspace] = useState('/');
  const [library, setLibrary] = useState<Library>(empty);
  const [deviceArticles, setDeviceArticles] = useState<Article[]>(initialDeviceArticles);
  const [busy, setBusy] = useState(false);
  const [selected, setSelected] = useState<number[]>([]);
  const session = useRef<any>(null);
  const item = useRef<string | null>(null);
  const running = library.jobs.find(job => job.status === 'queued' || job.status === 'running');
  const last = library.jobs[0];
  const search = library.jobs.find(job => job.kind === 'search');
  const sourceList = search?.status === 'done' ? search.sources || [] : [];

  function field(name: string, value: string) {
    setFields(previous => ({...previous, [name]: value}));
  }
  function clearPrivate() {
    session.current = null;
    setOwner(false);
    setLibrary(empty);
    setPanel('');
    setFields({provider: 'deepseek'});
    setSelected([]);
  }
  async function api(path: string, method = 'GET', body?: object) {
    const current = session.current;
    const response = await fetch('./api/' + path, {method, cache: 'no-store',
      headers: {'Content-Type': 'application/json', 'X-CSRF-Token': current?.csrf || '',
        'X-Workspace-Expected-Account': current?.user.id || ''},
      body: body === undefined ? undefined : JSON.stringify(body)});
    if (response.status === 401 || (current && response.headers.get('X-Workspace-Account') &&
        response.headers.get('X-Workspace-Account') !== current.user.id)) clearPrivate();
    const result = await response.json();
    if (!response.ok) throw Error(result.error || 'The request could not finish.');
    return result;
  }
  async function refresh() {
    if (!session.current) return;
    const current = session.current;
    try {
      const result = await api('fieldnotes');
      if (session.current === current) setLibrary(result);
    } catch (error) {
      if (session.current === current) setNotice((error as Error).message);
    }
  }
  useEffect(() => {
    const listen = () => { setReading(readHash()); window.scrollTo(0, 0); };
    window.addEventListener('hashchange', listen);
    if (offline) return () => window.removeEventListener('hashchange', listen);
    let active = true;
    (async () => {
      try {
        const boot = await api('bootstrap');
        if (!active) return;
        if (boot.workspace) setWorkspace(boot.workspace);
        const next = await api('me');
        if (!active || next.user.role !== 'owner') return;
        session.current = next;
        setOwner(true);
        await refresh();
      } catch { /* The public reader stays available while signed out. */ }
    })();
    const timer = setInterval(refresh, 3500);
    return () => { active = false; clearInterval(timer); window.removeEventListener('hashchange', listen); };
  }, []);
  useEffect(() => { setSelected([]); }, [search?.id]);
  useEffect(() => {
    if (last?.status === 'done' && last.kind === 'brief' && library.articles[0]?.id === last.id &&
        sessionStorage.getItem('fieldnotes-opened') !== last.id) {
      sessionStorage.setItem('fieldnotes-opened', last.id);
      open(last.id);
    }
  }, [last?.id, last?.status, library.articles[0]?.id]);
  useEffect(() => {
    if (!device) return;
    window.fieldnotesImport = raw => {
      if (raw === null) { setNotice('Could not read that file. Choose a Fieldnotes pack under 2 MB.'); return; }
      try {
        const incoming = readPack(raw);
        const merged = new Map(deviceArticles.map(a => [a.id, a]));
        for (const article of incoming) merged.set(article.id, article);
        const next = [...merged.values()];
        persist(next);
        setFolder('saved');
        open('');
        setNotice(incoming.length ? 'Reading pack imported.' : 'This reading pack contains no briefings.');
      } catch (error) { setNotice((error as Error).message); }
    };
    window.fieldnotesBack = () => {
      if (!reading) return false;
      open('');
      return true;
    };
    return () => { delete window.fieldnotesImport; delete window.fieldnotesBack; };
  }, [deviceArticles, reading]);

  function persist(articles: Article[]) {
    if (!device?.write(encodePack(articles))) throw Error('Could not save this reading list on the device.');
    setDeviceArticles(articles);
  }
  function open(id: string) {
    history.replaceState(null, '', id ? '#read-' + id : location.pathname + location.search);
    setReading(id);
    window.scrollTo(0, 0);
    requestAnimationFrame(() => {
      const title = document.querySelector<HTMLElement>('.reading h1, .feed h1, .tabs button');
      title?.setAttribute('tabindex', '-1');
      title?.focus({preventScroll: true});
    });
  }
  const saved = new Set(offline ? deviceArticles.map(a => a.id) : library.saved);
  const all = [...(offline ? deviceArticles : library.articles), ...starter]
    .filter((a, i, list) => list.findIndex(b => b.id === a.id) === i);
  const display: DisplayArticle[] = all.map(article => {
    const starterArticle = starterIDs.has(article.id);
    const words = article.sections.reduce((sum, part) => sum + part.body.trim().split(/\s+/).length, 0);
    const date = new Date(article.date + 'T00:00:00Z');
    const byline = (starterArticle ? 'Fieldnotes' : 'DeepSeek briefing') + ' · ' +
      (Number.isNaN(date.getTime()) ? article.date : dateFormat.format(date)) + ' · ' + Math.max(1, Math.ceil(words / 180)) + ' min read';
    return {...article, byline, label: starterArticle ? article.topic : 'Research briefing',
      tone: article.id === 'five-stages' ? 'step' : article.id === 'inclusive-interface' ? 'grid-cover' : 'loop',
      sources: article.sources.map((source, index) => ({...source, index: index + 1})), saved: saved.has(article.id)};
  });
  const reader = display.find(article => article.id === reading);
  const featured = folder === 'all' && display[0];
  const stories = folder === 'saved' ? display.filter(article => article.saved) : display.slice(1);

  async function action(name: string, explicit?: string) {
    const value = explicit || item.current || '';
    setNotice('');
    if (name === 'read') { open(value); return; }
    if (name === 'back' || name === 'home' || name === 'all' || name === 'saved') {
      if (name === 'home' || name === 'all') setFolder('all');
      if (name === 'saved') setFolder('saved');
      open('');
      return;
    }
    if (name === 'import') { device?.pick(); return; }
    if (name === 'close') { setPanel(''); return; }
    if (name === 'topic' && offline) {
      const match = display.find(article => article.title.toLowerCase().includes(value.split(' ').slice(0, 2).join(' ')));
      if (match) open(match.id);
      else setNotice('Import a reading pack for more topics.');
      return;
    }
    if (name === 'topic' || name === 'research' || name === 'settings') {
      if (!owner) { setNotice('Sign in with the Workspace owner account to research topics.'); return; }
      open('');
      setPanel(name === 'settings' ? 'settings' : 'research');
      if (name === 'topic') field('query', value);
      return;
    }
    if (busy) return;
    setBusy(true);
    try {
      if (name === 'bookmark') {
        const id = value || reading;
        const article = all.find(a => a.id === id);
        if (!article) return;
        if (offline) persist(saved.has(id) ? deviceArticles.filter(entry => entry.id !== id) : [...deviceArticles, article]);
        else await api('saved', 'PUT', {id, saved: !saved.has(id)});
      } else if (name === 'export') {
        const payload = encodePack(all.filter(a => saved.has(a.id)));
        const url = URL.createObjectURL(new Blob([payload], {type: 'application/json'}));
        const link = document.createElement('a');
        link.download = 'ai-fieldnotes-reading-pack.json';
        link.href = url;
        link.click();
        setTimeout(() => URL.revokeObjectURL(url), 1000);
      } else if (name === 'search' && !running) {
        await api('research', 'POST', {query: fields.query || '', provider: fields.provider});
      } else if (name === 'generate' && !running) {
        await api('research', 'POST', {search: search?.id, indices: selected});
      } else if (name.startsWith('save') || name.startsWith('remove')) {
        const deepseek = name.endsWith('DeepSeek');
        const key = deepseek ? 'deepseek' : 'browser';
        if (name.startsWith('save') && !fields[key]?.trim()) throw Error('Enter the provider credential first.');
        await api('providers', 'PUT', {provider: deepseek ? 'deepseek' : 'workspace',
          key: name.startsWith('save') ? fields[key].trim() : ''});
        field(key, '');
        setNotice(name.startsWith('save') ? 'Credential saved.' : 'Credential removed.');
      }
      await refresh();
    } catch (error) {
      setNotice((error as Error).message);
    } finally {
      setBusy(false);
    }
  }

  const data = {offline, notice, owner, workspace, login: './auth/login', androidLink: releaseLink,
    topics, reader, featured, stories, hasStories: stories.length > 0, emptyList: folder === 'saved' && !stories.length,
    listHeading: folder === 'saved' ? 'Saved briefings' : 'More to read',
    allTab: folder === 'all' ? 'active-tab' : 'tab', savedTab: folder === 'saved' ? 'active-tab' : 'tab',
    bookmarkLabel: reader?.saved ? 'Remove from reading list' : 'Save to reading list',
    bookmarkDisabled: (!offline && !owner) || busy, allowExport: owner && folder === 'saved' && !!saved.size,
    researchOpen: panel === 'research', settingsOpen: panel === 'settings', busy, browser: library.browser,
    deepseekStatus: library.deepseek ? 'Key saved' : 'Enter API key', browserStatus: library.browserReady ? 'Credential saved' : 'Enter credential',
    deepseekMissing: !library.deepseek || busy, browserMissing: !library.browserReady || busy,
    providers: [{id: 'deepseek', label: 'DeepSeek Search', disabled: false},
      {id: 'workspace', label: 'Workspace Web Search', disabled: !library.browserReady}],
    jobStatus: running ? (running.kind === 'search' ? 'Searching the web…' : 'Writing from selected sources…') :
      last?.error || (last?.kind === 'search' && last.status === 'done' ? 'Search complete.' : ''),
    searchDisabled: busy || !!running || !fields.query?.trim() || (fields.provider === 'deepseek' ? !library.deepseek : !library.browserReady),
    generateDisabled: busy || !!running || !selected.length || !library.deepseek,
    sources: sourceList.map((source, index) => ({...source, index, selected: selected.includes(index)})), hasSources: sourceList.length > 0};
  return h(Controls.Provider, {value: {fields, field, item, action, offline, busy, bookmarkDisabled: data.bookmarkDisabled,
    selectedCount: selected.length, pickSource: (i: number) =>
      setSelected(previous => previous.includes(i) ? previous.filter(j => j !== i) : previous.length < 5 ? [...previous, i] : previous)}},
    h(Reader, {data, onAction: action}));
}

createRoot(document.getElementById('root')!).render(h(App));
