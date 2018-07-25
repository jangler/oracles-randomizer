# TODO

## high priority

- only clear marks when necessary
- maybe only return true/false from getmark/peekmark?

## mid priority

- seeds
	- try changing satchel to start with a different type of seed
	- try putting seeds in chests and see if it affects the drops you can get
	- try randomizing trees
- rod
- verify that dungeons are completeable in *any* possible key order
- verify that you can't reach a L-1 item after the corresponding L-2 item
	- only run this when placing a L-1 item
	- maybe just disable the L-1 item after the L-2 equivalent is placed

## low priority

- subrosian dance hall -> dimitri's flute
- scramble dungeons? before the randomization step
- verify that you can't reach a slingshot before the seed satchel
	- only run this when placing the slingshot
	- this is unlikely to happen because of the routing possibilities, but it's
	  almost certainly possible
	- also just don't enable the slingshots until the satchel is placed

## flags

change flag structure:

no flags = randomize, input rom arg(0) and output rom arg(1), all implemented
dungeons.

-dryrun = don't actually write anything, and don't require any additional args.

-pointgen filename = what -op pointgen does now.

-checkgraph = what -op checkGraph does now.

-verify filename = what -op verifyData does now.

-start node = treat the given node as the start location, not horon village.
probably just have this default to horon village, actually, and not have horon
village be a special case in the graph.

-goals nodelist = nodes that must be reached before the randomizer is
satisfied. one string, comma-separated.

-forbid nodelist = nodes that must be impossible to reach in the generated rom.
naturally this can cause the command to fail, maybe after running for a very
long time. one string, comma-separated.

this seems like a good place to start…
