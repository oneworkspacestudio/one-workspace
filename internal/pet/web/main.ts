import {createElement as h, useEffect, useState} from 'react';
import {createRoot} from 'react-dom/client';
import {Pet} from './generated/Pet';
import {Fields, type PetState, type SettingsState, type World} from './ui';
import './style.css';

type Session = {csrf: string; user: {id: string; role: string}};
type Bootstrap = {oidc: boolean; workspace: string};
const times = (minute: number) => String(Math.floor(minute / 60)).padStart(2, '0') + ':' + String(minute % 60).padStart(2, '0');
const verbs: Record<string, string> = {
  explore: 'Exploring the meadow', sleep: 'At the burrow', eat: 'Foraging for a snack',
  play: 'Playing in the meadow', draw: 'Working on a painting', build: 'Building a sculpture',
};
const signs: Record<string, string> = {
  explore: '✧', home: '⌂', treat: '✿', sleep: '☾', play: '◈', art: '◒', build: '▦', eat: '●',
};
async function json(path: string, options?: RequestInit) {
  const response = await fetch(path, {cache: 'no-store', ...options});
  const body = await response.json();
  if (!response.ok) throw Object.assign(new Error(body.error || 'Pet could not finish the request.'), {status: response.status});
  return body;
}

function App() {
  const [config, setConfig] = useState<Bootstrap>();
  const [session, setSession] = useState<Session>();
  const [state, setState] = useState<PetState>();
  const [loading, setLoading] = useState(true);
  const [expired, setExpired] = useState(false);
  const [notice, setNotice] = useState('');
  const [offline, setOffline] = useState(false);
  const [busy, setBusy] = useState(false);
  const [name, setName] = useState('');
  const [journalOpen, setJournalOpen] = useState(false);
  useEffect(() => {
    let active = true;
    const ctl = new AbortController();
    async function open() {
      try {
        const bootstrap = await json('./api/bootstrap', {signal: ctl.signal});
        if (active) setConfig(bootstrap);
        const me = await json('./api/me', {signal: ctl.signal});
        if (active) setSession(me);
      } catch (e: any) {
        if (active && e.status !== 401) setNotice(e.message);
      } finally {
        if (active) setLoading(false);
      }
    }
    void open();
    return () => { active = false; ctl.abort(); };
  }, []);
  const account = session?.user.id;
  useEffect(() => {
    if (!account) return;
    const ctl = new AbortController();
    let timer: number;
    const poll = async () => {
      try {
        const next = await json('./api/world', {
          signal: ctl.signal, headers: {'X-Workspace-Expected-Account': account},
        }) as PetState;
        if (ctl.signal.aborted) return;
        setState(prev => {
          if (prev?.world && (!next.world || next.world.tick < prev.world.tick ||
              next.world.tick === prev.world.tick && next.world.treatAfter < prev.world.treatAfter)) {
            return {...next, world: prev.world};
          }
          return next;
        });
        setOffline(false);
      } catch (e: any) {
        if (ctl.signal.aborted) return;
        if (e.status === 401 || e.status === 409) {
          setSession(undefined);
          setState(undefined);
          setJournalOpen(false);
          setExpired(true);
          setNotice('Your session changed. Continue with MyAccounts to return to the hollow.');
          return;
        }
        setOffline(true);
      }
      timer = window.setTimeout(poll, document.hidden ? 8000 : 1800);
    };
    void poll();
    return () => { ctl.abort(); clearTimeout(timer); };
  }, [account]);

  const write = async (url: string, method: string, data?: unknown) => {
    if (!session) return;
    try {
      return await json(url, {method, body: data === undefined ? undefined : JSON.stringify(data), headers: {
        'Content-Type': 'application/json', 'X-CSRF-Token': session.csrf, 'X-Workspace-Expected-Account': session.user.id,
      }});
    } catch (e: any) {
      if (e.status === 401) {
        setSession(undefined);
        setState(undefined);
        setJournalOpen(false);
        setExpired(true);
      }
      throw e;
    }
  };

  const settings = state && {
    mode: state.mode, modelStatus: state.modelStatus, localModel: state.localModel || '',
    paired: !!state.paired, hasService: !!state.hasService,
  } as SettingsState;
  const save = async (values: {mode: string; model: string; key: string; clearKey: boolean}) => {
    await write('./api/decisions', 'PUT', values);
    const next = await json('./api/world', {headers: {'X-Workspace-Expected-Account': account!}});
    setState(next);
  };
  const world = state?.world;
  const cooldown = Math.max(0, (world?.treatAfter || 0) - (world?.saved || 0));
  const action = async (kind: string) => {
    if (kind === 'journal') {
      setJournalOpen(true);
      return;
    }
    if (busy) return;
    setBusy(true);
    setNotice('');
    try {
      if (kind === 'logout') {
        await write('./auth/logout', 'POST');
        setJournalOpen(false);
        setState(undefined);
        setSession(undefined);
      }
      if (kind === 'adopt' || kind === 'treat') {
        const updated = await write('./api/'+kind, 'POST', kind === 'adopt' ? {name: name.trim() || 'Pip'} : undefined) as World;
        setState(prev => prev && (!prev.world || prev.world.tick <= updated.tick) ? {...prev, world: updated} : prev);
        if (kind === 'treat') setNotice('A clover treat for '+updated.name+'!');
      }
    } catch (e: any) {
      setNotice(e.message);
    } finally {
      setBusy(false);
    }
  };
  const clock = world?.minute ?? 510;
  const dark = clock < 360 || clock >= 1200;
  const creations = [...(world?.creations ?? [])].sort((a,b) => b.id-a.id);
  return h(Fields.Provider, {value: {name, setName, save}},
    h(Pet, {data: {
      workspace: config?.workspace || '/', oidc: config?.oidc || expired, session: !!session,
      loading: loading || !!session && !state && !offline,
      signedOut: !loading && !session, unadopted: state && !world,
      owner: session?.user.role === 'owner', notice: notice || (offline ? 'Connection interrupted. Reconnecting to the hollow…' : ''),
      busy, world, weather: world?.weather === 'Rain' ? 'Rain in the meadow' : dark ? 'A clear night' : 'Clear skies',
      sun: dark ? '☾' : world?.weather === 'Rain' ? '☂' : '☀', clock: times(clock), day: 'DAY '+(world?.day || 1),
      pageStyle: world ? 'playing' : undefined, journalOpen, closeJournal: () => setJournalOpen(false),
      activity: world ? (world.mood === 'Dreaming' ? 'Asleep in the burrow' : verbs[world.action]) : '',
      events: (world?.journal ?? []).slice(-6).reverse().map(event => ({
        ...event, icon: signs[event.kind] || '·', time: times((event.tick*2+510)%1440),
      })),
      creations, galleryCount: String(creations.length), emptyGallery: !creations.length,
      treatDisabled: busy || offline || cooldown > 0, treatLabel: cooldown ? 'Next treat in '+cooldown+'s' : 'Give a clover treat',
      settings,
    }, onAction: action}));
}

createRoot(document.querySelector('#root')!).render(h(App));
