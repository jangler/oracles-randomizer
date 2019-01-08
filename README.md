# Zelda Oracles Randomizer

This program reads a Zelda: Oracle of Seasons or Oracle of Ages ROM (US
versions only), shuffles the locations of (most) items and mystical seeds, and
writes the modified ROM to a new file. In Seasons, the default seasons for each
area are also randomized. Most arbitrary overworld checks for essences and
other game flags are removed, so the dungeons and other checks can be done in
any order that the randomized items facilitate. However, you do need to collect
all 8 essences to get the Maku Seed and finish the game.


## Usage

There are three ways to use the randomizer:

1. Place the randomizer in the same directory as your vanila ROM(s) (or vice
   versa), and run it. The randomizer will automatically find your vanilla
   ROM(s) and prompt for further options.
2. In Windows, drag your vanilla ROM onto the executable. Same deal as above,
   except that the ROM and randomizer don't have to be in the same folder.
3. Use the command line. Type `./oracles-randomizer -h` to view the usage
   summary.


## Download

You can download executables for Windows, macOS, and Linux from the
[releases](https://github.com/jangler/oracles-randomizer/releases) page. Don't
use the "Download ZIP" link on the main page; that only contains the source
code. The download also contains a rudimentary location checklist and item
tracker. If you're looking for a more detailed item and map tracker,
[EmoTracker](https://emotracker.net/) has a pack developed by Herreteman.


## Randomization notes

General details common to both games:

- Items and chests are randomized, with these exceptions:
    - Renewable shop and business scrub items (bombs, shield, hearts, etc.)
	- Small keys
	- Gasha seeds and pieces of heart outside of chests
	- NPCs that give non-progression items in the vanilla game
	- Gasha nut contents
	- Fixed drops (from bushes, pots, etc.)
	- Maple drops
- Mystical seed trees are randomized, with no more than two trees of each type.
  Items that use seeds for ammunition start with the type of seed that's on the
  Horon Village or Lynna City tree.
- For items that have two levels, the first you obtain will be L-1, and the
  second will be L-2, regardless of the order in which you obtain them. The L-2
  shield is an exception.
- There is one flute in the game for a random animal companion, and it's
  identified and usable as soon as you get it. Only the 150-rupee item in the
  shop is randomized; the other two usual means of getting a strange flute
  don't give anything special. The animal companion regions (Natzu in Seasons
  and Nuun in Ages) match whatever flute is in the seed.
- Rings are instantly appraised when you get them, and the ring list can be
  accessed from the inventory ring box icon. For convenience, the L-3 ring box
  is given at the start. The punch rings can be used with only one equip slot
  empty.
- If tree warp is enabled, holding start while closing the map screen outdoors
  warps to the seed tree in Horon Village or Lynna City. Tree warp comes with
  no warranty and is not supported as a "feature", so think carefully before
  using it.

For game-specific notes on randomization and logic, see
[seasons_notes.md](https://github.com/jangler/oracles-randomizer/blob/master/doc/seasons_notes.md)
and
[ages_notes.md](https://github.com/jangler/oracles-randomizer/blob/master/doc/ages_notes.md)


## FAQ

**Q: Is there a place to discuss the randomizer?**

A: Yes, the Oracles Discord server (link
[here](https://www.speedrun.com/oos/thread/3qwe1)). The server is mainly focused
on speedrunning, but the #randomizer channel is for anything pertaining to the
randomizer.

**Q: I found a problem. What do I do?**

A: Open an issue about it on GitHub or bring it up in the #randomizer channel
in the Oracles discord. Provide your seed's log file either way.

**Q: Will you make a cross-game randomizer that combines Ages and Seasons into
one ROM?**

A: no

**Q: Can I at least do a linked game?**

A: Not yet. Probably in a future version.


## Thanks to:

- Drenn for [ages-disasm](https://github.com/Drenn1/ages-disasm) and additional
  code.
- Herreteman, dragonc0, Phoenomenom714, and jaysee87 for help with logic,
  playtesting, design, and "customer support".
- Everyone who helped playtest prerelease versions of the randomizer.
