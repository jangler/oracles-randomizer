# developer notes

## notable code addresses

the names of these, when present, correspond to the ones in drenn's
ages-disasm. the most useful are:

- $0b:4409 = when this executes, (hl) and (hl+1) are the given item ID and sub ID.
  this is for keys falling from ceilings, npcs giving items, most other items
  you don't receive from chests.
- $11:58b5 = parseObjectData. falls through to parseGivenObjectData.
- $11:58df = parseGivenObjectData. when this executes, de is the address of the
  start of an object's data. objects include enemies, puzzles, and special
  behaviors like what normally happens in the sword room of the hero's cave.
  **this is not always called directly when parsing interations.**
- $15:466b = hl-1 here is the index of the treasure item's info (collection
  mode, param, text, and sprite, in that order). in other words, (hl) is the
  treasure item's param.
- $3f:440a = when this executes, hl-1 is the start of the object's three-byte
  graphics data starting at $3f:63a3 (see the rom addresses section below for
  details).

others that might be good to know:

- $0e3b = drawObject
	- ID $60 animation = $13:409a
	- ID $59 animation = $14:4130
- $271a = createTreasure
- graphics:
	- $15e9 = interactionInitGraphics, which calls the following:
	- $3f:4404 = interactionLoadGraphics
		- takes d = object struct high byte, returns a = animation index
		- returns a = animation index
		- looks up the interaction's graphics data in 3-byte table $3f:63a3
		  based on ID, loads the results into d+$1c to d+$1f (?)
	- $25ca = interactionSetAnimation
		- takes a = anim index, d = object struct high byte
		- looks up the interataction's animation data in table $3f:4bb5 based
		  on ID, loads the results into object struct
	- $1e41 = objectSetVisible
	- $1e57 = objectSetVisible80 (redundant)
	- $24fd = maybe not present in ages? copies some object data to other parts
	  of itself, like… direction?
- $239a = interactionIncState
- $3b22 = updateInteraction

## functions / code

these are jp:

- 0:2a15 = setLinkIDOverride
- 2:4f90 = openMenu
- 2:4fdd = closeMenu
- 4:460c = getTransformedLinkID
- 5:5468 = checkLinkForceState
- 6:4865 = checkUseItems
- 6:4911 = checkItem
- 6:4925 = initializeParentItem
- 6:4931 = chooseParentItemSlot
- 6:4994 = parentItemUpdate (the good stuff; what happens when an item is used)

these are en/us:

- 0:045b = copyMemoryReverse, b is # bytes, de is src, hl is dest
- 0:0462 = copyMemory, b is # bytes, hl is src, de is dest
- 0:0c74 = playSound, a is index
- 0:1432 = get tile at position bc (yyxx), returns a (id) and hl (addr)
- 0:16eb = giveTreasure (a is ID, c is param)
- 0:184b = showText, bc is index
- 0:1956 = getThisRoomFlags
- 0:24fe = interactionSetScript, hl is address in bank b
- 0:30c7, 0:30cd = checkGlobalFlag / setGlobalFlag, a is bit starting at c6ca
- 0:3958, 0:39ea, 0:39f9 = points for loading room tilemap address
- 0:3ac6 = getFreeInteractionSlot
- 3:4cf5 = intro capcomScreen
	- 3:4d68 = state1 (fading in)
- 5:4552 = companionTryToMount
- 5:5471 = linkSetState, a is state, d is object low byte
- 7:497b = itemLoadAttributesAndGraphics
- 7:49ca = itemSetAnimation
- 3f:454e = applyParameter when giving treasure (a is type, c is parameter, de
  is address to write to, b happens to be the treasure index)
- 3f:4445, 3f:444c, 3f:c45a = points for loading sprite data for an object

## script commands (most are documented in more detail in ages-disasm)

- 00 = end script
- 80 = set state of interaction.state
- 84 = spawn interaction
- 87 = jump table
- 88 = set coordinates, byte = y, byte = x
- 8f = set animation, byte = index
- 98 = show text, word = index
- 9c = set interaction text id, word = index
- a0 = wait for bit of cfc0 is set
- a7 = ? takes two bytes of params ?
- b0 = jump if room flag, byte = flag, word = addr
- b5 = jump if global flag, byte = flag, word = addr
- b6 = set global flag, byte = flag
- b8 = disable objects
- bd = disable input
- be = enable input
- c0 = call another script
- c3 = jump to script addr based on text option
- c4 = jump to script addr
- cd-cf = stop if bit 5/6/7 of room flags is set
- d3 = wait until flag is set, byte = flag, word = addr
- d4 = wait until object byte equals value
- d7 = set counter, byte = value
- de = spawn item on link, word = id, subid (?)
- e0 = call function in bank 15, word = addr
- e1 = call function in bank 15, word = addr, byte = value of a and e
- e3 = play sound, byte = index
- ec-ef = move npc up/down/left/right, byte = frames
- f6 = set object counter? like d7 but for a specific object?
- f7 = ? no params
- f8 = ? no params

## notable ram addresses

- c63f-c640 = bought shop items
- c643-c646 = companion state (ricky, dimitri, moosh, then misc.)
- c680-c6?? = inventory (starting with equipped items)
- c692-c6a1 = item flags
- c6a2/c6a3 = health / max health
- c6a5-c6a6 = rupees
- c6b5-c6b9 = seed count (ember, scent, pegasus, gale, mystery)
- c6bb = obtained essence flags
- c6c5 = active ring
- c6ca-c6d9 = some global flags

- c7xx = overworld room flags
- c8xx = subrosia room flags (& some group 2?)
- c9xx, caxx = etc

- cbb6 = index of room under cursor in map menu
- cc4e = current season
- cc48 = high byte of link object address (in object table starting at d000)
- cc49 = active group
- cc4c = active room
- cc63-cc66 = data about room transition (group, room, ???, position)
- ccab = allow screen transitions only if zero in treasure H&S
- ccb6 = active tile? rod of seasons only works when this == 8
- ccea = disable interactions (?)

- dx58-dx59 = script address

## notable rom addresses (leftover JP stuff)

- $3f:63a3 is a table of three byte graphics data sets. the middle is a flag
  and the other two combine to form the address of the actual data ($97 $80
  $68) becomes de=$6897… except this needs another multiple-of-three offset
  before finishing, based on the location of the desired graphics in wram (?).
  $3f:68a3 is an example used for bombs in the shop; it reads $5c $10 $40,
  where $5c and $10 combine to become $d0 (sprite ID) and $40 becomes $01
  (nybbles swapped) for palette/transform flags. i don't remember where the
  first nybble of $5c goes or what it's for. it looks like these are in the
  same order as the shop item data table.
- $3f:4bb5 = interactionAnimationTable
- $14:4d85 = interactionOamDataTable
- $6:49a7 = item-specific code jump table
- $6:5508 = itemUsageParameterTable
	- offset by 2 * item ID

### rom (leftover JP stuff)

- 04:6d4e = pointer table for room transitions?
	- plus group * 2?
	- then that value plus cc64 * 3? but it's zero when entering hero's cave
	- then (cc64) <- (hl), (cc66) <- (hl+1), (cc65) <- (hl+2) | 80, (cc6a) <-
	  0a
	- then (ff8b) <- (cc49), (cc49) <- (cc63), so (cc63) was the group, and
	  (cc4c) <- (cc64), so (cc64) was the room (so (hl) was the room).
- 04:7457 = pointer table for room transitions?
	- when entering hero's cave (thru main entrance), group * 2 is added to
	  this number, then (hl+1) is loaded into (cc64), and the high and low
	  nybbles of (hl+2) are loaded into (cc63) and (cc65) respectively.
		- for hero's cave entrance, hl+1 is 00 and hl+2 is 44. hero's cave's
		  room and group are both 04.
	- when leaving hero's cave, group * 2 added, … and then some things happen.
- 07:4219 = ? something to do with room loading ?
- 11:5b38 = objectDataGroupTable (pointers to object data by group)
- 15:53af = chest treasure pointer table, offset by 2 * group, then increment
  hl in steps of 4 bytes until (hl+1) == the room number, but abort if (hl) ==
  ff. the next two bytes are the treasure ID and sub ID.
