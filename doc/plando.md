# oracles-randomizer plando support

4.0 adds the option to use a partial or complete spoiler log as input to the
randomizer. If the spoiler log is partial (or empty!), the other variables are
left vanilla, with the exception of the ring pool, which is still random for
unspecified rings. The names in the log may be external spoiler log names (like
"wooden/noble sword") or internal ones (like "sword").

Currently the only way to create a plando is via the `-plan` command-line
option. Multiworld plandos are not supported.


## Sections

### `-- items --` (or any of the usual item section names)

The number of each item need not match the vanilla distribution, although due
to technical limitations there may not be more rings than usual, and there can
be only one type of flute. Spheres and divisions between progression items /
keys / etc do not need to be specified.


### `-- dungeon entrances --`

Multiple entrances can technically link to the same dungeon, but each dungeon
can only have one exit, which is randomly chosen from its entrances.


### `-- subrosia portals --`

Seasons only. Multiple Holodrum portals linking to the same Subrosia portal
will have the same issue that dungeon entrances do.


### `-- default seasons --`

Seasons only. No special notes.
