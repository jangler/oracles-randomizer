# developer notes

## notable code addresses

the names of these, when present, correspond to the ones in drenn's
ages-disasm. the most useful are:

- $10bf = when this executes, (hl) and (hl+1) here are the current chest item
  ID and sub ID.
- $0205 = checks for bit a in the flags starting at hl. $1717 does this
  specifically at $c692 (treasure flags), and $30b3 does this specifically at
  $c6ca (global flags).
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
- $3ab2 = getFreeInteractionSlot, called when a new "interaction" is needed.
  this is used when a room's interactions are loaded (almost anything that's
  not a static tile), but it's also used when creating new objects, like the
  floodgate key. if there's a `ld (hl),60` afterward, that means it's an item
  interaction, and the item ID and subID are usually in registers b and c.
- $3f:440a = when this executes, hl-1 is the start of the object's three-byte
  graphics data starting at $3f:63a3 (see the rom addresses section below for
  details).

others that might be good to know:

- $0e3b = drawObject
	- ID $60 animation = $13:409a
	- ID $59 animation = $14:4130
- $16eb = giveTreasure, which uses param as the second byte, not sub-ID
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
- $074e = copy byte from hl to de, incrementing both

## functions

- 0:041a = getRandomNumber
- 0:2a15 = setLinkIDOverride
- these aren't functions, but they have something to do with tiles for a given
  room when they execute:
	- 0:3944
	- 0:39d6
	- 0:39e0
- 2:4f90 = openMenu
- 2:4fdd = closeMenu
- 4:460c = getTransformedLinkID
- 5:5468 = checkLinkForceState
- 5:5471 = linkSetState
- 6:4865 = checkUseItems
- 6:4911 = checkItem
- 6:4925 = initializeParentItem
- 6:4931 = chooseParentItemSlot
- 6:4994 = parentItemUpdate (the good stuff; what happens when an item is used)
- 7:4f36 = galeSeedTryToWarpLink

## notable ram addresses

these all get checked in a normal frame, just for display purposes:

- $c680-$c6?? = inventory (starting with equipped items)
- $c6a2/$c6a3 = health / max health
- $c6a5-$c6a6 = rupees
- $c6b5-$c6b9 = seed count (ember, scent, pegasus, gale, mystery)
- $c6c5 = active ring
- $c6ca-$c6d9 = some global flags
- $cc4c = active room

other things:

- $c63f-$c640 = bought shop items
- $c643-$c646 = companion state (ricky, dimitri, moosh, then misc.)
	- shop checks for bit 5 of ricky
- $c692-$c6a1 = item flags
- $c6bb = obtained essence flags
- $cc48 = high byte of link object address (in object table starting at $d000)
- $ccea = disable interactions (?)
- $c6c5 = wActiveRing

## notable rom addresses

- $15:57FD + $4X = index of ring given by param X; params below 4 don't
  (normally?) work. this is a generalization of the information described for
  $15:466b.
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

## horon village shop (also syrup?)

- $08:4cce is a table of shop item data. these are two-byte ID, sub ID pairs.
  changing the item in a slot doesn't change the sprite, price, text, or logic
  associated with it, just the item you get after you buy it.
- $11:6646 is the start of the object data for the shop. the objects are
  four-byte blocks: interaction ID, sub ID (acting as offset into the shop data
  table), and (y,x) coords. changing the sub ID still doesn't change the item
  logic (?)
- $08:4c93 is the start of the item price table.
- $08:4c6b is the start of the item price display coords table. this is weird
  and doesn't precisely correspond to (x,y) coords?
- after changing *all* these values, i still can't get item replacement to
  work. after you say yes to buying the item, it just goes back to its place
  without any other text.

## room loading

### wram

- c65c = wGashaMaturity
- cc49 = wActiveGroup (overworld, subrosia, dungeons, etc)
- cc4c = wActiveRoom
- cc63-cc65 = loading room info
- ccc5 = seasons-specific? wRotatingCubePos in ages-disasm

### rom

- 09a8 = flagLocationGroupTable, for rooms
- 01:4662 = ? something to do with room loading ?
- 01:5db5 = ? something to do with room loading ?
- 04:483c = groupMusicPointerTable (offset = wActiveGroup * 2, then
  (hl) + wActiveRoom into actual music table)
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

### code

- 1955 = getThisRoomFlags; a <- flags, hl <- addr of flags
	- 1962 = getRoomFlags; takes a = group, b = room; a <- flags, hl <- addr of
	  flags (h = (table + group), l = room)
- 3003 = initializeRoom
- 01:5ece = ld a,(wActiveGroup); or a; ret nz
- 17fa = add a to wGashaMaturity, +5 when entering room
- 3276 = loadScreenMusic
- getObjectDataAddress:
	- uses group * 2 as index into pointer table, adds room * 2 to the value,
	  loads the address at that value into de… not present in seasons lol, it's
	  inlined in getObjectData

### tiles

- d9 = snow pile
