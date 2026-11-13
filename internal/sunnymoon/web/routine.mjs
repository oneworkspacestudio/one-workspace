export const MINUTE = 60_000;
export const DEFAULTS = Object.freeze({ bedtime: '21:00', wake: '07:00', lead: 60 });
export const LEADS = [15, 30, 45, 60, 90];
const DAY_MINUTES = 24 * 60;

export function minutesOfDay(time) {
  if (typeof time !== 'string' || !/^(?:[01]\d|2[0-3]):[0-5]\d$/.test(time)) return NaN;
  const [hour, minute] = time.split(':').map(Number);
  return hour * 60 + minute;
}

export function validSettings(settings) {
  if (!settings || !LEADS.includes(settings.lead)) return false;
  const daytime = (minutesOfDay(settings.bedtime) - minutesOfDay(settings.wake) + DAY_MINUTES) % DAY_MINUTES;
  return daytime > settings.lead && daytime <= DAY_MINUTES - 60;
}

function atTime(reference, time, offset = 0) {
  const minutes = minutesOfDay(time);
  return new Date(reference.getFullYear(), reference.getMonth(), reference.getDate() + offset,
    Math.floor(minutes / 60), minutes % 60);
}

export function timeInput(date) {
  return [date.getHours(), date.getMinutes()].map(value => String(value).padStart(2, '0')).join(':');
}

export function routineAt(now, settings) {
  const offset = now < atTime(now, settings.wake) ? -1 : 0;
  const start = atTime(now, settings.wake, offset);
  const nextWake = atTime(now, settings.wake, offset + 1);
  const crossesMidnight = minutesOfDay(settings.bedtime) < minutesOfDay(settings.wake);
  const bedtime = atTime(now, settings.bedtime, offset + Number(crossesMidnight));
  const reminder = new Date(Math.max(start.getTime(), bedtime.getTime() - settings.lead * MINUTE));
  const windingDown = new Date(Math.max(start.getTime(), bedtime.getTime() - 105 * MINUTE));
  const sleepy = new Date(Math.max(start.getTime(), bedtime.getTime() - 60 * MINUTE));
  const mood = now <= start ? 'happy' : now >= sleepy ? 'sleepy' : now >= windingDown ? 'winding-down' : 'happy';
  const remaining = (nextWake - Math.max(now.getTime(), bedtime.getTime())) / MINUTE;
  const drowsiness = Math.min(1, Math.max(0, (now - windingDown) / Math.max(MINUTE, sleepy - windingDown)));
  return { start, nextWake, bedtime, reminder, mood, remaining, drowsiness };
}

export function belongsToRoutine(timestamp, routine, now) {
  return Number.isFinite(timestamp) && timestamp >= routine.start.getTime() && timestamp <= now.getTime();
}

export function earlierTime(time, routine, now) {
  if (!Number.isFinite(minutesOfDay(time))) return null;
  const visit = atTime(routine.start, time, Number(time < timeInput(routine.start)));
  return belongsToRoutine(visit.getTime(), routine, now) ? visit.getTime() : null;
}

export function metersAt(now, routine, lastVisit) {
  const sinceVisit = now - (lastVisit ?? routine.start.getTime());
  const fill = (elapsed, interval) => Math.min(100, Math.max(0, Math.floor(100 * elapsed / interval)));
  return {
    bladder: fill(sinceVisit, 240 * MINUTE),
    bowel: fill(sinceVisit, 480 * MINUTE),
    tiredness: fill(now - routine.start, routine.bedtime - routine.start),
  };
}

export function toiletCue(now, routine, meters, bedtimeChecked) {
  if (now <= routine.start || bedtimeChecked) return false;
  if (now >= routine.reminder) return true;
  return Math.max(meters.bladder, meters.bowel) >= 75;
}

export function sleepLabel(minutes) {
  const rounded = Math.floor(minutes);
  if (rounded < 1) return 'Under 1 min until morning';
  const hours = Math.floor(rounded / 60);
  const extra = rounded % 60;
  if (!hours) return `${extra} min of sleep`;
  if (!extra) return `${hours} ${hours === 1 ? 'hour' : 'hours'} of sleep`;
  return `${hours} hr ${extra} min of sleep`;
}
