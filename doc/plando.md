# oracles-randomizer plando support

4.0 adds the option to use a partial or complete spoiler log as input to the
randomizer. If the spoiler log is partial, the other variables are filled in as
they would be in normal seed generation. The configuration must be valid (i.e.
completeable) in randomizer logic. The names in the log may be either external
spoiler log names or internal ones.

Currently the only way to create a plando is via the `-plan` command-line
option.


## Sections

### `-- items --` (or any of the usual item section names)

The distribution of items need not match the vanilla distribution, although due
to technical limitations there may not be more rings than usual, and there can
be only one type of flute. Instances of items can be removed from the pool
altogether by placing them in slots named `null`, although this does *not*
apply to seed trees and rings. If there are fewer items in the pool than there
are slots, the remainder are filled with Gasha Seeds. Spheres and divisions
between progression items / keys / etc do not need to be specified.


### `-- dungeon entrances --`

Dungeon shuffle is automatically enabled if any randomized dungeon entrances
are specified.


### `-- subrosia portals --`

Seasons only. Subrosia portal shuffle is automatically enabled if any
randomized portals are specified.


### `-- default seasons --`

Seasons only. No special notes.


### `-- hints --`

Owl text may consist only of printable ASCII characters in the range ' ' to
'z'.  Not all puntuation characters will actually print correctly in-game.
There is no specific limit on the length of a hint, but words may not be longer
than 16 characters, and the cumulative length may not exceed the available bank
space.
