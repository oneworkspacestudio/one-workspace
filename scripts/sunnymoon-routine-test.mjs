import assert from 'node:assert/strict';
import { DEFAULTS, MINUTE, earlierTime, metersAt, routineAt, sleepLabel, validSettings } from '../internal/sunnymoon/web/routine.mjs';

process.env.TZ = 'Europe/London';
const at = (day, hour, minute = 0) => new Date(2026, 10, day, hour, minute);
const routine = now => routineAt(now, DEFAULTS);

assert.equal(routine(at(12, 7)).mood, 'happy');
assert.deepEqual(metersAt(at(12, 7), routine(at(12, 7)), null), { bladder: 0, bowel: 0, tiredness: 0 });
assert.equal(routine(at(12, 19, 14)).mood, 'happy');
assert.equal(routine(at(12, 19, 15)).mood, 'winding-down');
assert.equal(routine(at(12, 20)).mood, 'sleepy');
assert.deepEqual(routine(at(12, 20)).reminder, at(12, 20));
assert.equal(sleepLabel(routine(at(12, 21)).remaining), '10 hours of sleep');
assert.equal(sleepLabel(routine(at(12, 22, 25)).remaining), '8 hr 35 min of sleep');

const afternoon = at(12, 15);
const afterSchool = routine(afternoon);
const schoolVisit = earlierTime('13:00', afterSchool, afternoon);
assert.equal(schoolVisit, at(12, 13).getTime());
assert.equal(earlierTime('15:01', afterSchool, afternoon), null);
assert.deepEqual(metersAt(afternoon, afterSchool, schoolVisit), { bladder: 50, bowel: 25, tiredness: 57 });
assert.deepEqual(metersAt(afternoon, afterSchool, afternoon.getTime()), { bladder: 0, bowel: 0, tiredness: 57 });

assert.equal(routine(at(13, 0)).mood, 'sleepy');
assert.deepEqual(routine(at(13, 0)).start, at(12, 7));
assert.equal(earlierTime('23:05', routine(at(13, 0)), at(13, 0)), at(12, 23, 5).getTime());
assert.equal(routine(at(13, 7)).mood, 'happy');
assert.equal(earlierTime('23:05', routine(at(13, 7)), at(13, 7)), null);

const later = { bedtime: '00:30', wake: '08:30', lead: 45 };
assert.ok(validSettings(later));
assert.equal(validSettings({ ...later, wake: '00:50' }), false);
assert.equal(validSettings({ ...later, wake: 'garbled' }), false);
const afterMidnight = routineAt(at(13, 0), later);
assert.deepEqual(afterMidnight.bedtime, at(13, 0, 30));
assert.deepEqual(afterMidnight.reminder, at(12, 23, 45));
assert.equal(afterMidnight.remaining, 480);

const spring = routine(new Date(2026, 2, 28, 20));
assert.equal((spring.nextWake - spring.bedtime) / MINUTE, 540);
console.log('Passed routine timing, check-ins, rollover, custom settings and DST.');
