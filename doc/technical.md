# technical notes

function names and ram address names, when present, correspond to those in
drenn's ages-disasm. single addresses are for seasons; double addresses are
ages/seasons, in that order. if an address ends in a slash, it's the same for
both games.


## functions not documented elsewhere

- 0:008a/ = interBankCall
- 0:043e/041a = getRandomNumber
- 0:047f/045b = copyMemoryReverse, b is # bytes, de is src, hl is dest
- 0:0486/0462 = copyMemory, b is # bytes, hl is src, de is dest
- 0:0775 = loadTileset, a is index
- 0:0c98/0c74 = playSound, a is index
- 0:10cc/109a = getChestData
- 0:1435 = get tile at position bc (yyxx), returns a (id) and hl (addr)
- 0:15e9 = interactionInitGraphics
	- 3f:4404/ = interactionLoadGraphics
- 0:171c/16eb = giveTreasure (a is ID, c is param)
- 0:1733/1702 = loseTreasure (a is ID)
- 0:17e5 = refill all seeds
- 0:17b9 = getRandomRingOfGivenTier
- 0:1ddd = lookupCollisionTable, hl = table, scf if a is in table
- 0:21fd, 0:2202, 0:2215 = objectCopyPosition, objectCopyPosition rawAddress,
  objectCopyPositionWithOffset
- 0:2727 = objectCreateExclamationMark
- 0:24d2 = interactionActuallyRunScript
- 0:24fe = interactionSetScript, hl is address in bank b
- 0:250c = runScript, d is object low byte
- 0:2542 = interactionSaveScriptAddress
- 0:2d2a = getThisRoomDungeonFlags
	- bit 4 = has key / boss key (for compass beep)
	- bit 5 = has chest (for map display)
	- bit 6 = ?? if this is set, no compass beep
- 0:393e = loadSmallRoomLayout
	- 0:3958, 0:39ea, 0:39f9 = points for loading room tilemap address
	- 0:3979, 0:3987 = decompressLayoutMode2, decompressLayoutMode2Helper
	- 0:399c, 0:39aa = decompressLayoutMode1, decompressLayoutMode1Helper
	- 0:39cb = decompressLayoutHelper
- 0:3aef/3ac6 = getFreeInteractionSlot
- 0:3b36 = updateInteraction, d is object low byte
- 0:3ea7 = getFreePartSlot
- 1:49e5 = check for compass beep
- 1:5ece = updateSeedTreeRefillData
- 3:4cf5 = intro capcomScreen
	- 3:4d68 = state1 (fading in)
- 5:44aa = specialObjectGetRelativeTileWithDirectionTable
- 5:4552 = companionTryToMount
- 5:493b = companionRetIfNotFinishedWalkingIn
- 5:5471 = linkSetState, a is state, d is object low byte
- 5:5fdb = checkCliffTile, scf if cliff
- 6:4713 = tryToBreakTile body
- 7:497b = itemLoadAttributesAndGraphics
- 7:49ca = itemSetAnimation
- 15:463f = some function for loading treasure data based on object id/subid
- 3f:454e = applyParameter when giving treasure (a is type, c is parameter, de
  is address to write to, b happens to be the treasure index)
- 3f:4445, 3f:444c, 3f:445a = points for loading sprite data for an object


## ram addresses not documented elsewhere

- c63f-c640 = bought shop items
- c688/c680 = inventory (starting with equipped items)
- c6a2-c6a3 = health-maxHealth
- c6a5-c6a6, c6a7-c6a8 = rupees, ore chunks
- c6b0 = obtained seasons
- c6b5-c6b9 = seed count (ember, scent, pegasus, gale, mystery)
- c6bb = obtained essence flags
- c6c0-c6c4, c6cb/c6c5 = rings in box, active ring

- cbb6/ = index of room under cursor in map menu
- cc4e = current season
- cc2d/cc49 = active group
- cc30/cc4c = active room
- cc63-cc66 = data about room transition (group, room, ??, link position)
- cc89 = level of shield that link is using, if using a shield
- ccab = allow screen transitions only if zero in treasure H&S
- ccb6 = active tile (under link)? rod of seasons only works when this == 8
- ccea = disable interactions (?)

- cd00 = 0 while screen transitioning, 1 when transition done (useful for wait
  command in scripts)


## flags

treasure flags begin at c69a/c692, are indexed by item ID, and determine
whether link is considered to have a given item (regardless of whether it
appears in his equip menu, or whether its other parameters such as quantity or
level are set. treasure flags are checked by 0:1748/1717,
checkTreasureObtained.

(some) global flags begin at c6d0/c6ca. these are pretty general-purpose. they
are checked and set by 0:31f3/30c7 and 0:31f9/30cd, respectively.

room flags begin at cx00 depending on the group the room is in, starting at
c700 for group 0 (overworld), c800 for group 1 (past/subrosia) and group 2
(buildings), c900 (?) for group 4 (caves and dungeons), and ca00 (?) for group
5 (other caves and dungeons). bit 4 tracks whether the room has been explored,
and bit 5 is commonly used to track whether the treasure in the room has been
obtained (e.g. if a chest has been opened). some treasure rooms such as shops
and npcs check other flags (usually treasure flags) instead of room flags, but
there's a jp-only bug where the master diver and the chest in his room *both*
set bit 5, meaning that if you get the master diver's item, the chest will be
opened the next time you visit the room (but not vice versa, i think). bit 7
appears to immediately delete some interactions if set? bits 6 and 3 are also
used, but i don't know for what in particular.

current room flags are checked by function 0:197d/1956, getThisRoomFlags.

animal companion flags go from c646-c649/c643-c646, with the bytes being for
ricky, dimitri, moosh, then misc. bit 7 determines whether the animal is ready
to be ridden; i think the others are specific to each animal.

### ages flags

- 0a = finished intro (can open menu?)
- 0c = maku tree vanished
- 0e = fairies put the forest back in order
- 11 = maku tree tells you where the seventh essence is
- 12 = past maku tree opened gate
- 13 = finished talking to maku tree after getting maku seed
- 15 = gave rafton rope; spawns ricky in forest of time
- 16 = defeated great moblin… then unset right afterward?
- 18 = set when entering nayru/veran fight
- 19 = set when escaping nayru/veran fight
- 1a = defeated great moblin? or unset when defeating great moblin
- 1b = talked to tingle
- 1c = set after getting bomb capacity upgrade from fairy
- 1d = set after d3 essence, needed for flute, unset when entering nuun
- 20 = set when talking to cursed queen fairy; needed to receive fairy powder
- 21 = landed in the world at start of game
- 22 = reset by bridge-building foreman; needed for animal companion event?
- 23 = checked by bridge-building foreman, needed to retrieve workers
- 24 = checked in fairies' woods?
- 25 = bridge built
- 26 = rafton completed raft
- 27 = cured king zora
- 29 = checked in fairies' woods?
- 2a = entered symmetry city brother's house
- 2b = fairies put the forest back in order
- 2e = listened to symmetry city wife's problem
- 2f = got crown key from goron elder
- 30 = cleansed the zora seas
- 31 = got permission from king zora to enter jabu-jabu
- 32 = saw ralph cutscene outside ambi's palace
- 34 = got eyeball from captain; checked every frame on every screen
- 35 = finished twinrova cutscene after getting maku seed
- 36 = traded mystery seeds for feather
- 37 = traded scent seeds for bracelet
- 3d = set at start
- 3e = maku tree told you to go to yoll graveyard
- 3f = past maku tree opened gate (2)
- 40 = ralph goes back in time
- 41 = surprised guy runs away when entering past
- 42 = checked in fairies' woods?
- 43 = talked to cheval; triggers ralph cutscene outside
- 44 = hit maple
- 46 = got satchel upgrade from tingle


## treasures

16:5332/15:5129 is the treasure data pointer table, indexed by treasure ID,
then incremented by sub ID. each treasure entry is four bytes: collection mode,
parameter used in giveTreasure, text index, and sprite index.

collection mode determines whether an item appears from a chest when opened, is
bought in a shop, is dropped from the ceiling, is simply found lying around,
etc. several collection modes seem to be interchangeable, but the chest
collection modes are important—if a the collection mode of an item in a chest
isn't a chest collection mode, the game will get stuck in the chest-opening
cutscene. dungeon map and compass chests have their own collection mode, which
primarily determines animation, but also requires the item to have a valid text
index; otherwise the game will stick in the chest-opening cutscene as well.

the giveTreasure parameter means something different for different treasure
IDs. for some treasures, it determines the level to set the object to (though
this will not decrease the level), and for some it determines the quantity
given. for many treasures, the parameter has no effect and is always zero.

treasure data can be retrieved based on ID and sub ID using the rom data
script.


## interactions

each room has a list of interactions associated with it, with pointer table
starting at 15:432b/11:5b3b and being indexed by room group, then room ID (in
ages, the pointers are to bank 12, not 15). fx bytes denote the beginning of an
interaction or series of interactions that take the same parameters, except for
f3, fe and ff, which respectively mean "jump to pointer", "return from
pointer", and "end list". unambiguous entries not beginning with fx jump to the
denoted address.

room interaction lists can be retrieved based on group and room ID using the
rom data script.

interactions generally (or always?) create objects (see objects section).


## objects

the dxxx block of ram is devoted to object info. entries for link and animal
companions start at d000 and d100 respectively, and entries for other entities
start at dx40. many functions take dx as the high byte of an object's address
in ram, and supply the low byte to read/write particular variables. these
generally assume that the object's data starts at dx40 and not dx00.

displayed text also uses this block, and seems to overwrite some other object
data while being displayed?


## scripts

each object can have a script associated with it if its dx58-dx59 word is
nonzero. this word determines the current address in a series of script
commands in bank c/b. scripts are run concurrently unless a particular script
does something like disable all other objects. bank 16?/15 contains unique
functions called by scripts; general-purpose ones are in bank 0.

for a complete list of script commands, see
ages-disasm/include/script_commands.s, but some are listed here with addresses.

- 00 = end script
- 80 = 4186, set interaction.state
- 81 = 4197
- 82 = 419d
- 83 = 258f
- 84 = 41a2, spawn interaction
- 85 = 41ca
- 86 = 4a1c
- 87 = 41f8, jump table
- 88 = 4203, set coordinates, byte = y, byte = x
- 89 = 4213
- 8a = 4280
- 8b = 421a
- 8c = 422c
- 8d = 4236
- 8e = 4240
- 8f = 425e, set animation, byte = index
- 90 = 42fd
- 91 = 4319
- 92 = 43fc, or memory, word = addr, byte = value
- 93 = 4252
- 94 = 4247
- 95 = 4221
- 96 = 4290
- 97 = 42a0
- 98 = 42c1, show text, word = index
- 99 = 42ed
- 9a = 42e2
- 9b = 43c7
- 9c = 43de, set interaction text id, word = index
- 9d = 43d1
- 9e = 44e4
- 9f = 42cc
- a0-a7 = 432b, wait until bit of cfc0 is set
- a8-af = 433d, toggle bit of cfc0
- b0 = 43a7, jump if room flag, byte = flag, word = addr
- b1 = 43bb
- b2 = 4103, custom script, set ccaa = 01 and jump to command a0
- b3 = 435d, jump if c6xx set, byte = value to bitwise and with word addr
- b4 = 44ca
- b5 = 4562, jump if global flag, byte = flag, word = addr
- b6 = 4573, set global flag, byte = flag
- b7 = 4103, nop
- b8 = 4168, disable objects
- b9 = 4170
- ba = 4173
- bb = 414c
- bc = 4162
- bd/be = 4147/415e, disable/enable input
- bf = 4103, nop
- c0 = call another script
- c3 = jump to script addr based on text option
- c4 = jump to script addr
- cd-cf = stop if bit 5/6/7 of room flags is set
- d3 = wait until flag is set, byte = flag, word = addr
- d4 = wait until object byte equals value
- d7 = set counter, byte = value
- de = spawn item on link, word = id, subid
- e0 = call function in bank 15, word = addr
- e1 = call function in bank 15, word = addr, byte = value of a and e
- e3 = play sound, byte = index
- ec-ef = move npc up/down/left/right, byte = frames
- fx = ??


## room layouts

room layouts start with the top left tile and proceed left to right. some parts
of rooms are compressed, meaning that groups of the same tile (which need not
be consecutive) are represented by one instance of the tile ID. the
compression is determined periodically by some non-tile bytes.

compression "mode 1" uses one byte as a bitmask of locations in the next 8
tiles where the tile at the next byte appears. following those two bytes, there
is one byte for each unset bit in the bitmask to determine the non-compressed
tiles. the cycle repeats.

large rooms (i.e. dungeons) use a different kind of compression.

the address of a small room's layout can be viewed in function 0:393e,
loadSmallRoomLayout. banks 21, 22, 23, and 24 are used for spring, summer,
autumn, and winter layouts, respectively.

tiles common to different tilesets usually have the same IDs.


## graphics

sprite entries in the table at 3f:6427/6425 are three bytes. the first two
determine the index of the sprite, and the second byte is evens only; odds
don't make a difference. the third determines palette, transformations, and
other information. one nybble is devoted exclusively (?) to which palette the
sprite uses, although only 8 exist.


## text

text in the game is encoded using ascii values for letters and some punctuation
from 20 (space) to 7a (lower case z). a value of 01 is a newline, and 00 ends
the text. 02 to 05 are prefixes for dictionaries 00 to 03. for example, 03 1f
would "evaluate" to entry 1f in dictionary 01. 09 is the prefix for a color
change, with 00 being the default white.

0c seems to be the start token for a text entry, and the following character
determines the position of the text box on the screen. 00 positions it
automatically at the top or bottom, depending on link's y position.

07 is the prefix for a jump command, the next byte being the low byte of the
current text block.  for example, index 2601 ending with 07 03 would jump to
index 2603.

see ages-disasm/text/ for text IDs and dictionary entries.

3f:5c00 is the text pointer table; text IDs are two bytes, so the address at
5c00+2h is read first, then the address at (5c00+2h)+2l is read. the game has
two further tables for "text offsets" at 3f:4fe2 and 3f:4ffa, the first for
text IDs less than 2c00 and the second for those greater or equal.

0:184b is showText, and 0:1936 reads a byte from the text table. watching this
function is probably the easiest way to determine where particular text data is
in the rom. watching 3f:4fa4 lets you determine the location of text *address*
data (stored in hl).
