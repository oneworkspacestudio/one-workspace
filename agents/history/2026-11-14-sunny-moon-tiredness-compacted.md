# Sunny & Moon evening tiredness

Sunny & Moon preview.4 uses the existing drowsiness value to boost the
tiredness meter over the last three hours before bedtime. In a routine with
7 am wake-up and 9 pm bedtime, the bar reads 94% at 7:15 pm, 98% at 8 pm
and 100% at bedtime. The meter resets at wake-up. Toilet cues keep their
icon and animation while the eyes continue to follow Sunny’s bedtime mood.
The settings, check-in data and managed update contract are compatible.

The routine regression tests and the 11-size browser matrix passed. The
390×844 browser checked the rendered 94% bar and drooped eyelids with the
toilet cue active. The build targets Workspace preview.131 with the verified
public preview.130 seed. Use the mandatory signed Android publisher. Build
and test artifacts are under the tooling folder `sunny-moon-tiredness-20261114`.
