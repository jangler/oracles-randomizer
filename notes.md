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

## notable ram addresses

these all get checked in a normal frame, just for display purposes:

- $c680-$c6?? = inventory (starting with equipped items)
- $c6a2/$c6a3 = health / max health
- $c6a6 = rupees
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

## notable rom addresses

- $15:57FD + $4X = index of ring given by param X; params below 4 don't
  (normally?) work. this is a generalization of the information described for
  $15:466b.
