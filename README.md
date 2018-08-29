# Oracle of Seasons randomizer

This program reads a Zelda: Oracle of Seasons ROM (US or JP version) shuffles
the locations of items and mystical seeds, randomizes the default season for
each area, and writes the modified ROM to a new file. It also bypasses essence
checks for overworld events that are necessary for progress, so the dungeons
can be done in any order that the randomized items facilitate. However, you do
have to collect all 8 essences to get the Maku Seed and finish the game.

The randomizer is relatively new and under active development, so consider it
"beta" for now. See the [issue
tracker](https://github.com/jangler/oos-randomizer/issues) for known problems.


## Usage

The randomizer uses a command-line interface, and I currently have no plans to
implement a graphical one. It's a simple program (from the user's perspective),
and command lines are not very hard.

The normal usage is `./oos-randomizer oos_original.gbc oos_randomized.gbc` (or
whatever filenames you want), but there are additional flags you can pass
before the filename arguments, as displayed in the usage (`./oos-randomizer
-h`) message:

    Usage of ./oos-randomizer:
      -freewarp
            allow unimited tree warp (no cooldown)
      -keyonly
            only randomize key item locations
      -seed string
            specific random seed to use (32-bit hex number)
      -update
            update already randomized ROM to this version
      -verbose
            print more detailed output to terminal


## Download

You can download executables for Windows, macOS, and Linux from the
[releases](https://github.com/jangler/oos-randomizer/releases) page.


## Randomization notes

Essential items and chests are randomied, with exceptions listed below. The rod
of seasons is split into four items, each of which will give you one season and
the rod itself (if you don't already have it).

Seed trees and default seasons for each area are also shuffled, and the satchel
and slingshot will start with the type of seeds on the tree in Horon Village.

The following items are **not** randomized:

- Shop items
- Small keys and boss keys
- Pirate's bell, hard ore, and iron shield
- Found items (gasha seeds and pieces of heart outside of chests)
- Subrosian dancing prizes after the first
- Trading sequence items
- Non-essential items given by NPCs
- Subrosian hide and seek items
- Gasha nut contents
- Fixed drops
- Maple drops

If the `-keyonly` flag is specified, only key items (the items required to
complete a normal game) and their locations are shuffled. Speedrunners should
note that the first Subrosian dancing prize could still be important.


## Other notable changes

Other small changes have been made for convenience, to simplify randomization
logic, or to prevent softlocks. The most notable are:

- The intro sequence and pirate cutscene are almost entirely removed.
- Mystical seeds grow in all seasons.
- Seeds can be collected if the player has either a slingshot or the satchel.
- The cliff between Eastern Suburbs and Sunken City has stairs instead of a
  spring flower.
- Rosa doesn't appear in the overworld, and her portal is activated by default.
- The diving spot at the south end of Sunken City is removed.
- **Holding start while closing the map screen warps to the seed tree in Horon
  Village.** Tree warping has a one-hour cooldown unless the `-freewarp` flag
  is specified.

## FAQ

**Q: When I run the randomizer, a command prompt window opens and closes
without doing anything. What do I do?**

A: If you're still lost after reading the "usage" section of the readme, either
Google how to use the command prompt or ask a friend. Replace command prompt
with Unix shell if you're on macOS.

**Q: Do I have to do HSS skip or Poe skip?**

A: No, but you can if you want to, and the randomizer accounts for those and other
sequence breaks.

**Q: I'm softlocked. Now what do I do?**

A: If you're softlocked by location, use tree warp. Otherwise, open an issue
about it or tell me in Discord, and provide the log file. Depending on the
problem, you may be able to `-update` your ROM using the next patch version to
un-softlock.

**Q: Are you going to make a randomizer for Oracle of Ages too?**

A: Maybe, but not until the Seasons randomizer is reasonably feature-complete
(as i see it). Ages also has some big sequence breaks and the Crescent Island
trading sequence, which would both be tricky to account for in the logic unless
they're just removed.
