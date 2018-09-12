# developer notes

function names and ram address names, when present, correspond to those in
drenn's ages-disasm.


## functions not documented elsewhere

- 0:008a = interBankCall
- 0:045b = copyMemoryReverse, b is # bytes, de is src, hl is dest
- 0:0462 = copyMemory, b is # bytes, hl is src, de is dest
- 0:0c74 = playSound, a is index
- 0:1432 = get tile at position bc (yyxx), returns a (id) and hl (addr)
- 0:16eb = giveTreasure (a is ID, c is param)
- 0:2215 = objectCopyPositionWithOffset
- 0:2727 = objectCreateExclamationMark
- 0:24d2 = interactionActuallyRunScript
- 0:24fe = interactionSetScript, hl is address in bank b
- 0:250c = runScript, d is object low byte
- 0:2542 = interactionSaveScriptAddress
- 0:393e = loadSmallRoomLayout
	- 0:3958, 0:39ea, 0:39f9 = points for loading room tilemap address
- 0:3ac6 = getFreeInteractionSlot
- 0:3b36 = updateInteraction, d is object low byte
- 3:4cf5 = intro capcomScreen
	- 3:4d68 = state1 (fading in)
- 5:4552 = companionTryToMount
- 5:5471 = linkSetState, a is state, d is object low byte
- 6:4713 = tryToBreakTile body
- 7:497b = itemLoadAttributesAndGraphics
- 7:49ca = itemSetAnimation
- 3f:454e = applyParameter when giving treasure (a is type, c is parameter, de
  is address to write to, b happens to be the treasure index)
- 3f:4445, 3f:444c, 3f:c45a = points for loading sprite data for an object


## ram addresses not documented elsewhere

- c63f-c640 = bought shop items
- c680-c6?? = inventory (starting with equipped items)
- c6a2/c6a3 = health / max health
- c6a5-c6a6, c6a7-c6a8 = rupees, ore chunks
- c6b0 = obtained seasons
- c6b5-c6b9 = seed count (ember, scent, pegasus, gale, mystery)
- c6bb = obtained essence flags
- c6c0-c6c4, c6c5 = rings in box, active ring

- cbb6 = index of room under cursor in map menu
- cc4e = current season
- cc49 = active group
- cc4c = active room
- cc63-cc66 = data about room transition (group, room, ??, link position)
- ccab = allow screen transitions only if zero in treasure H&S
- ccb6 = active tile (under link)? rod of seasons only works when this == 8
- ccea = disable interactions (?)

- cd00 = 0 while screen transitioning, 1 when transition done (useful for wait
  command in scripts)


## flags

treasure flags begin at c692, are indexed by item ID, and determine
whether link is considered to have a given item (regardless of whether it
appears in his equip menu, or whether its other parameters such as quantity or
level are set. treasure flags are checked by 0:1717, checkTreasureObtained.

(some) global flags begin at c6ca. these are pretty general-purpose. they are
checked and set by 0:30c7 and 0:30cd, respectively.

room flags begin at cx00 depending on the group the room is in, starting at
c700 for group 0 (overworld), c800 for group 1 (subrosia) and group 2
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

current room flags are checked by function 0:1956, getThisRoomFlags.

animal companion flags go from c643-c646, with the bytes being for ricky,
dimitri, moosh, then misc. bit 7 determines whether the animal is ready to be
ridden; i think the others are specific to each animal.


## treasures

15:5129 is the treasure data pointer table, indexed by treasure ID, then
incremented by sub ID. each treasure entry is four bytes: collection mode,
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
starting at 11:5b3b and being indexed by room group, then room ID. fx bytes
denote the beginning of an interaction or series of interactions that take the
same parameters, except for f3, fe and ff, which respectively mean "jump to
pointer", "return from pointer", and "end list". unambiguous entries not
beginning with fx jump to the denoted address.

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
commands in bank b. scripts are run concurrently unless a particular script
does something like disable all other objects. bank 15 contains unique
functions called by scripts; general-purpose ones are in bank 0.

for a complete list of script commands, see
ages-disasm/include/script_commands.s, but some are listed here for quick
reference.

- 00 = end script
- 80 = set interaction.state
- 84 = spawn interaction
- 87 = jump table
- 88 = set coordinates, byte = y, byte = x
- 8f = set animation, byte = index
- 98 = show text, word = index
- 9c = set interaction text id, word = index
- a0 = wait until bit of cfc0 is set
- a7 = ? takes two bytes of params ?
- b0 = jump if room flag, byte = flag, word = addr
- b3 = jump if c6xx set, byte = value to bitwise and with word addr
- b5 = jump if global flag, byte = flag, word = addr
- b6 = set global flag, byte = flag
- b8 = disable objects
- bd/be = disable/enable input
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
- fx = ??


## room layouts

room layouts start with the top left tile and proceed left to right. some parts
of rooms are compressed, meaning that groups of the same tile (which need not
be consecutive) are represented by one instance of the tile ID. the
compression is determined periodically by some non-tile bytes. i don't know how
it works beyond that.

large rooms (i.e. dungeons) use a different kind of compression.

the address of a small room's layout can be viewed in function 0:393e,
loadSmallRoomLayout. banks 21, 22, 23, and 24 are used for spring, summer,
autumn, and winter layouts, respectively.

tiles common to different tilesets usually have the same IDs.


## graphics

sprite entries in the table at 3f:63a3 (i think this is the same in jp and
en/us?) are three bytes. the first two determine the index of the sprite, and
the second byte is evens only; odds don't make a difference. the third
determines palette, transformations, and other information. one nybble is
devoted exclusively (?) to which palette the sprite uses, although only 8
exist.


## text

text in the game is encoded using ascii values for letters and some punctuation
from 20 (space) to 7a (lower case z). a value of 01 is a newline, and 00 ends
the text. 02 to 05 are prefixes for dictionaries 00 to 03. for example, 03 1f
would "evaluate" to entry 1f in dictionary 01.

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
