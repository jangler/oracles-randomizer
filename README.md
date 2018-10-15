# Oracle of Seasons randomizer

This program reads a Zelda: Oracle of Seasons ROM (US version only), shuffles
the locations of items and mystical seeds, randomizes the default season for
each area, and writes the modified ROM to a new file. It also bypasses essence
checks for overworld events that are necessary for progress, so the dungeons
can be done in any order that the randomized items facilitate. However, you do
need to collect all 8 essences to get the Maku Seed and finish the game.


## Usage

There are three ways to use the randomizer:

1. In Windows, place the randomizer in the same directory as your vanila OoS
   ROM (or vice versa), and run it. The randomizer will automatically choose
   the vanilla ROM and write the randomized ROM and log to a new file.
2. In Windows, drag your vanilla OoS ROM onto the executable. The randomizer
   will write the randomized ROM and log to a new file.
3. Use the command line. Type `oos-randomizer -h` to view the usage summary.
   This is required to specify hard difficulty using the `-hard` option and to
   specify a seed using `-seed`.


## Download

You can download executables for Windows, macOS, and Linux from the
[releases](https://github.com/jangler/oos-randomizer/releases) page. Don't use
the "Download ZIP" link on the main page; that only contains the source code.


## Randomization notes

Items and chests are randomized, with exceptions listed below. The rod of
seasons is split into four items, each of which will give you one season and
the rod itself (if you don't already have it).

There is one flute in the game for a random animal companion, and it's
identified and usable as soon as you get it. Subrosian dancing and Ricky do not
give flutes as they normally would. The Natzu region matches whichever animal
companion the randomized flute calls.

Seed trees and default seasons for each area are also shuffled, and the satchel
and slingshot will start with the type of seeds on the tree in Horon Village.
The duplicate tree (normally a gale tree) has a random seed type instead.

For items that have two levels, the first you obtain will be L-1, and the
second will be L-2, regardless of the order in which you obtain them. The L-2
shield is an exception.

The following items are **not** randomized:

- Renewable shop items (bombs, shield, hearts, etc.)
- Small keys
- Pirate's bell (obtained by polishing rusty bell)
- Gasha seeds and pieces of heart outside of chests
- Subrosian dancing prizes after the first
- Trading sequence items
- Non-essential items given by NPCs
- Subrosian hide and seek items
- Gasha nut contents
- Fixed drops
- Maple drops


## Other notable changes

Other small changes have been made for convenience, to simplify randomization
logic, or to prevent softlocks. The most notable are:

- The intro sequence and pirate cutscene are removed, and the Maku Seed
  cutscenes are abbreviated.
- Mystical seeds grow in all seasons, and can be collected with a slingshot as
  well as a satchel.
- Rosa doesn't appear in the overworld, and her portal is activated by default.
- Fool's ore is randomized (the Strange Brothers trade you nothing for your
  feather). Shovel is not required to retrieve the stolen feather.
- Holding start while closing the map screen outdoors (in the overworld or in
  Subrosia) warps to the seed tree in Horon Village. This also sets your
  save/respawn point to that screen. Tree warping has a one-hour cooldown
  unless the `-freewarp` flag is specified. Tree warp is not supported as a
  "feature" and has no warranty, so consider possible consequences before using
  it.
- In some situations, the game will give you warnings about what you're doing
  or about to do. **If you receive one, what you're doing is out of logic and
  could potentially lead to a softlock**—but in some cases you can also be fine
  as long as you're careful.


## FAQ

**Q: Is [thing] in logic?**

A: See
[logic.md](https://github.com/jangler/oos-randomizer/blob/master/doc/logic.md).

**Q: I'm softlocked. Now what do I do?**

A: If you're softlocked by location, use tree warp. In any case, open an issue
about it or tell me in Discord, and provide the log file.

**Q: Are you going to make a randomizer for Oracle of Ages too?**

A: Ages support will probably be the next feature priority for the randomizer,
although I don't have immediate plans to start working on it. I don't know Ages
nearly as well as I know Seasons, so expect Ages randomization to be initially
poor at avoiding softlock cases. Perfoming D2 skip or text warps will
definitely void any warranty of safety from the randomizer.

**Q: Is there a place to discuss the randomizer?**

A: Yes, the Oracles Discord server (link
[here](https://www.speedrun.com/oos/thread/3qwe1)). The server is mainly focused
on speedrunning, but the #randomizer channel is for anything pertaining to the
randomizer.

**Q: Is there an item tracker I can use?**

A: [EmoTracker](http://emosaru.com/index.php/emotracker/) has an oos-randomizer
plugin developed by Herreteman.
