# Oracle of Seasons randomizer

This program reads a Zelda: Oracle of Seasons ROM (probably JP version only; I
haven't tested on EN/EU), shuffles the locations of key items, and writes the
modified ROM to a new file. It also bypasses essence checks for overworld
events that allow progress, so the dungeons can be done in any order that the
randomized items facilitate. And they *will* facilitate a path to the goal,
although you may have to be inventive to find it.

The randomizer has not been publicly tested yet, so consider it "beta" for now.
I've done a full run without problems, but it doesn't mean that there aren't
any bugs. I particularly recommend saving before opening any chest, because if
the randomizer screwed up then you'll be stuck and have to hard reset.


## Usage

The randomizer is a command-line interface, and I currently have no plans to
implement a graphical one. It's a simple program (from the user's perspective),
and command lines are not very hard.

The normal usage is `./oos-randomizer oos_original.gbc oos_randomized.gbc` (or
whatever filenames you want), but there are additional flags you can pass
before the filename arguments to fine-tune the randomization, as displayed in
the usage (`./oos-randomizer -h`) message:

    Usage of ./oos-randomizer:
      -devcmd string
        	if given, run developer command
      -dryrun
        	don't write an output file for any operation
      -forbid string
        	comma-separated list of nodes that must not be reachable
      -goal string
        	comma-separated list of nodes that must be reachable (default "done")
      -maxlen int
        	if >= 0, maximum number of slotted items in the route (default -1)
      -start string
        	comma-separated list of nodes treated as given (default "horon village")

Note that some combinations of these flags can result in impossible conditions,
like `-goal 'd1 essence' -forbid 'ember seeds'`. See further below for an
abbreviated list of possible `-goal` and `-forbid` nodes.

Regardless of the value of `-maxlen`, the randomizer will place items in all
available slots. The flag just limits the number of slotted items that are
*necessary* in order to reach the goal(s).


## Randomization notes

Most inventory items (equippable and non-equippable) that are necessary to
complete a casual playthrough are shuffled, with some exceptions:

- Purchasable items (bombs, shield, and strange flute) are not shuffled.
- The rod of seasons is always in the temple.
- The L-2 sword is always in the lost woods.
- The ribbon and pirate's bell are not shuffled.
- Maybe other stuff I'm forgetting

**Items are only placed in locations where you would normally obtain another
key item.**

Sometimes useful rings (primarily the fist/expert's rings) are placed in slots
instead of normal items.

Speedrunners should note that the Subrosian dancing prize could be important!


## Potentially useful goal/forbid nodes

- 'done', meaning defeating Onox
- 'dX essence', where X is a number from 1-8
- 'sword L-1'
- really just look at the files in the prenodes/ folder if you want more

Remember that you can specify multiple nodes as goals/forbids by separating
them with commas. Also remember to quote the strings lol
