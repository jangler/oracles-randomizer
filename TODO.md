# TODO

## high priority

- have a routine that finds all absolutely required items for a given goal
- verify that you can't reach a L-1 item after the corresponding L-2 item
	- just disable the L-1 item after the L-2 equivalent is placed

## mid priority

- have pointgen optimize nodes
- seeds
	- try changing satchel to start with a different type of seed
	- try putting seeds in chests and see if it affects the drops you can get
	- try randomizing trees
- verify that you can't reach a slingshot before the seed satchel
	- only run this when placing the slingshot
	- this is unlikely to happen because of the routing possibilities, but it's
	  almost certainly possible
	- also just don't enable the slingshots until the satchel is placed
- rod
- verify that dungeons are completeable in *any* possible key order

## low priority

- maybe only return true/false from getmark/peekmark?
- subrosian dance hall -> dimitri's flute
- scramble dungeons? before the randomization step
