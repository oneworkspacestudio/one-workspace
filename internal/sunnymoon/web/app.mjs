import { MINUTE, validSettings, routineAt, belongsToRoutine, timeInput,
  earlierTime, metersAt, sleepLabel } from './routine.mjs';

const element = id => document.getElementById(id);
const clock = new Intl.DateTimeFormat('en-GB', { hour: 'numeric', minute: '2-digit', hour12: true });
const weekday = new Intl.DateTimeFormat('en-GB', { weekday: 'long' });
const sunCard = element('sun-card');
const bedtimeVisit = element('bedtime-visit');
const settingsForm = element('settings-form');
const earlierForm = element('earlier-form');
const settingsDialog = element('settings-dialog');
const earlierDialog = element('earlier-dialog');
const feedback = element('feedback');
const sessionPanel = element('session-panel');
const visitButtons = [element('visit'), bedtimeVisit, element('open-earlier'), element('open-settings'),
  settingsForm.querySelector('[type=submit]'), earlierForm.querySelector('[type=submit]')];
const faces = {
  happy: ['Sunny is happy!', 'Let’s take today one step at a time.'],
  'winding-down': ['Sunny is slowing down', 'Time to get a little cosy.'],
  sleepy: ['Sunny is sleepy', 'A fresh, happy start in the morning.'],
};

let state;
let csrf;
let account;
let busy = false;
let displayedMorning;
let celebration;

function setText(id, value) {
  const target = element(id);
  if (target.textContent !== value) target.textContent = value;
}

function displayTime(date) {
  return clock.format(date).replaceAll('.', '');
}

function showSession(message, signIn = false) {
  state = undefined;
  account = undefined;
  csrf = undefined;
  displayedMorning = undefined;
  for (const dialog of [settingsDialog, earlierDialog]) {
    if (dialog.open) dialog.close();
  }
  setText('session-message', message);
  element('retry').hidden = signIn;
  element('sign-in').hidden = !signIn;
  element('routine').hidden = true;
  element('open-settings').hidden = true;
  sessionPanel.hidden = false;
}

async function request(value) {
  const headers = {};
  if (account) headers['X-Workspace-Expected-Account'] = account;
  if (value) {
    headers['Content-Type'] = 'application/json';
    headers['X-CSRF-Token'] = csrf;
  }
  const options = value ? { method: 'PUT', body: JSON.stringify(value) } : {};
  const result = await fetch('./api/routine', { ...options, headers, signal: AbortSignal.timeout(12_000) });
  if (result.status === 401) {
    showSession('Sign in to open your routine', true);
    return false;
  }
  const currentAccount = result.headers.get('X-Workspace-Account');
  if (currentAccount && account && currentAccount !== account) {
    showSession('The signed-in account changed');
    return false;
  }
  if (result.status === 409) return false;
  if (!result.ok) throw new Error('Routine unavailable');
  const data = await result.json();
  ({ csrf, routine: state } = data);
  account = currentAccount;
  return true;
}

function withCurrentVisits(value, routine, now) {
  const lastVisit = belongsToRoutine(value.lastVisit, routine, now) ? value.lastVisit : 0;
  const nightVisit = lastVisit && belongsToRoutine(value.nightVisit, routine, now) &&
    value.nightVisit >= routine.reminder.getTime() ? value.nightVisit : 0;
  return { ...value, lastVisit, nightVisit, expires: lastVisit ? routine.nextWake.getTime() : 0 };
}

function render(now = new Date()) {
  if (!state) return;
  const routine = routineAt(now, state.settings);
  const { start, bedtime, mood, reminder, remaining } = routine;
  const { lastVisit, nightVisit } = withCurrentVisits(state, routine, now);
  if (displayedMorning !== start.getTime()) {
    displayedMorning = start.getTime();
    feedback.textContent = '';
  }
  const afterBedtime = now >= bedtime;
  const reminderDue = now >= reminder;
  const checked = nightVisit !== 0;
  setText('today', weekday.format(start));
  setText('greeting', mood === 'happy' ? 'Hello, sunshine!' : 'Let’s get cosy');
  sunCard.dataset.mood = mood;
  setText('sun-message', faces[mood][0]);
  setText('sun-note', faces[mood][1]);
  setText('bedtime', displayTime(bedtime));
  element('bedtime-status').hidden = !afterBedtime;
  setText('sleep-estimate', sleepLabel(remaining));
  setText('sleep-note', afterBedtime ? 'If you get cosy now' : 'If you go to bed on time');
  setText('reminder-title', checked ? 'Toilet trip done!' : reminderDue ? 'Time for a toilet trip' : 'Toilet before bed');
  setText('reminder-note', checked ? 'All set for bedtime.' : reminderDue ? 'One quick stop before bed.' : `A toilet trip at ${displayTime(reminder)}`);
  bedtimeVisit.hidden = checked || !reminderDue;

  for (const [id, value] of Object.entries(metersAt(now, routine, lastVisit || null))) {
    element(id).value = value;
  }
  return routine;
}

async function update(change, error) {
  if (busy || !state) return false;
  busy = true;
  for (const button of visitButtons) button.disabled = true;
  try {
    for (let attempt = 0; attempt < 2; attempt++) {
      const value = change();
      if (!value) return false;
      if (await request(value)) {
        render();
        return true;
      }
      if (!state || !await request()) return false;
    }
    throw new Error('Routine changed');
  } catch {
    error.textContent = 'Couldn’t save this check-in. Try again.';
    return false;
  } finally {
    busy = false;
    for (const button of visitButtons) button.disabled = false;
  }
}

async function refresh() {
  if (busy) return;
  busy = true;
  try {
    if (await request()) {
      render();
      sessionPanel.hidden = true;
      element('routine').hidden = false;
      element('open-settings').hidden = false;
    }
  } catch {
    showSession('Couldn’t open your routine. Try again.');
  } finally {
    busy = false;
  }
}

async function recordVisit({ bedtimeCheck = false, earlier = false } = {}) {
  const error = earlier ? element('earlier-error') : feedback;
  let message;
  const saved = await update(() => {
    const now = new Date();
    const routine = routineAt(now, state.settings);
    const timestamp = earlier ? earlierTime(earlierForm.elements.time.value, routine, now) : now.getTime();
    if (timestamp === null) {
      error.textContent = `Choose a time between ${displayTime(routine.start)} ${weekday.format(routine.start)} and now.`;
      return null;
    }
    if (bedtimeCheck && timestamp < routine.reminder.getTime()) {
      render(now);
      return null;
    }
    const current = withCurrentVisits(state, routine, now);
    if (earlier && timestamp < current.lastVisit) {
      earlierDialog.close();
      feedback.textContent = 'Your newer check-in is already saved.';
      return null;
    }
    message = earlier ? `Visit saved for ${displayTime(new Date(timestamp))}.` : 'Nice listening to your body!';
    return { ...current, lastVisit: timestamp, nightVisit: bedtimeCheck ? timestamp : current.nightVisit,
      expires: routine.nextWake.getTime() };
  }, error);
  if (!saved) return;
  if (earlier) earlierDialog.close();
  if (bedtimeCheck) element('reminder-title').focus({ preventScroll: true });
  feedback.textContent = message;
  clearTimeout(celebration);
  sunCard.classList.remove('is-celebrating');
  requestAnimationFrame(() => {
    sunCard.classList.add('is-celebrating');
    celebration = setTimeout(() => sunCard.classList.remove('is-celebrating'), 560);
  });
}

element('visit').addEventListener('click', () => recordVisit());
bedtimeVisit.addEventListener('click', () => recordVisit({ bedtimeCheck: true }));
element('retry').addEventListener('click', refresh);

element('open-settings').addEventListener('click', () => {
  for (const [name, value] of Object.entries(state.settings)) {
    settingsForm.elements.namedItem(name).value = value;
  }
  element('settings-error').textContent = '';
  settingsDialog.showModal();
});

settingsForm.addEventListener('submit', async event => {
  event.preventDefault();
  const form = new FormData(settingsForm);
  const settings = { bedtime: form.get('bedtime'), wake: form.get('wake'), lead: Number(form.get('lead')) };
  const error = element('settings-error');
  if (!validSettings(settings)) {
    error.textContent = 'Choose times with at least an hour for sleep and a toilet reminder after wake-up.';
    return;
  }
  if (await update(() => {
    const now = new Date();
    return withCurrentVisits({ ...state, settings }, routineAt(now, settings), now);
  }, error)) settingsDialog.close();
});

element('open-earlier').addEventListener('click', () => {
  const now = new Date();
  const routine = render(now);
  const suggested = new Date(Math.max(routine.start.getTime(), now.getTime() - 30 * MINUTE));
  earlierForm.elements.time.value = timeInput(suggested);
  setText('visit-date', `Visit time since ${displayTime(routine.start)} ${weekday.format(routine.start)}`);
  element('earlier-error').textContent = '';
  earlierDialog.showModal();
});
earlierForm.addEventListener('submit', event => {
  event.preventDefault();
  recordVisit({ earlier: true });
});

for (const button of document.querySelectorAll('[data-close]')) {
  button.addEventListener('click', () => element(button.dataset.close).close());
}

window.addEventListener('pageshow', event => { if (event.persisted) refresh(); });
document.addEventListener('visibilitychange', () => { if (!document.hidden) refresh(); });
refresh();
setInterval(() => render(), 5_000);
