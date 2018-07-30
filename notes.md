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
- $11:58df = parseGivenObjectData. when this executes, de is the address of the
  start of an object's data. it is called n+1 times if there are n objects in a
  room as you enter. objects include enemies, puzzles, and special behaviors
  like what normally happens in the sword room of the hero's cave.
- $15:466b = hl-1 here is the index of the treasure item's info (collection
  mode, param, text, and sprite, in that order). in other words, (hl) is the
  treasure item's param.
- $3ab2 = getFreeInteractionSlot, called when a new "interaction" is needed.
  this is used when a room's interactions are loaded (almost anything that's
  not a static tile), but it's also used when creating new objects, like the
  floodgate key. if there's a `ld (hl),60` afterward, that means it's an item
  interaction, and the item ID and subID are usually in registers b and c.

others that might be good to know:

- $0e3b = drawObject
- $16eb = giveTreasure, which i believe offsets the treasure param such that it
  needs to be passed as one higher than usual. e.g. if you want $00, pass $01.
- $271a = createTreasure
- $15e9 = interactionInitGraphics
- $3b22 = updateInteraction
- $074e = copy byte from hl to de, incrementing both

## notable ram addresses

these all get checked in a normal frame, just for display purposes:

- $c680-$c6?? = inventory (starting with equipped items)
- $c6a2/$c6a3 = health / max health
- $c6a5-$c6a6 = rupees
- $c6b5-$c6b9 = seed count (ember, ?, ?, ?, ?)
- $c6c5 = active ring
- $c6ca-$c6d9 = some global flags
- $c8a6 = ?
- $cc4c = active room

other things:

- $c63f-$c640 = bought shop items
- $c643-$c646 = companion state (ricky, dimitri, moosh, then misc.)
	- shop checks for bit 5 of ricky
- $c692-$c6a1 = item flags
- $c6bb = obtained essence flags
- $cc77 = ?
- $cc48 = high byte of link object address (in object table starting at $d000)
- $ccea = disable interactions (?)

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
