# Sunny & Moon evening tiredness

Sunny & Moon preview.4 uses the existing drowsiness value to boost the
tiredness meter over the last three hours before bedtime. In a routine with
7 am wake-up and 9 pm bedtime, the bar reads 94% at 7:15 pm, 98% at 8 pm
and 100% at bedtime. The meter resets at wake-up. Toilet cues keep their
icon and animation while the eyes continue to follow Sunny’s bedtime mood.
The settings, check-in data and managed update contract are compatible.

The routine regression tests and the 11-size browser matrix passed. The
390×844 browser checked the rendered 94% bar and drooped eyelids with the
toilet cue active. Sunny & Moon, catalog, distribution and suite Go tests
passed. The suite tests required NTFS for both Go's temp and test data paths.

[Workspace preview.131](https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.131)
shipped from `be9c66b` and was promoted by `d7c2846`. The build verified
all 939 packages and preserved the unrelated app payloads from public
preview.130. The exact installer passed fresh setup, restart, app-only
upgrade and rollback using isolated NTFS data. Account, settings, files and
manager were preserved through the app-only upgrade. Synthetic fixtures
were stopped and removed. The signed publisher verified 959 release assets
and anonymously downloaded the APK. The anonymous installer selected and
verified preview.131 and reported Sunny & Moon preview.4. No owner install
or physical phone test was performed.

Build and test artifacts are under the tooling folder
`sunny-moon-tiredness-20261114`.
