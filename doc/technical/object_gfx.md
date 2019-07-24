# Object graphics technical summary

The document covers:

- How animations work
- How OAM data works
- How graphics tile data works

This document may contain errors, since I'm not an expert on any of this. I
know just enough about how these things work to manipulate them for the
purposes of the randomizer, and these notes are primarily intended for future
reference by myself and anyone else who wants to work on this sort of thing.

See http://gbdev.gg8.se/wiki/articles/Video_Display for a more detailed
description about how displaying graphics works on the GB/C in general.

Names for functions and struct fields used here are from
https://github.com/Drenn1/oracles-disasm.


## Nomenclature

Some ambiguous terms exist in this context:

- "Graphics" is a general term that I will use to describe "OAM data" and "tile
  data" in combination.
- "Animation frame" means an atomic step of an animation, while a "display
  frame" is the usual kind that happens at ~59.7 Hz on a GB/C.
- "Tile" here usually means 8x8 2bpp display data; the static 16x16 kind will
  be referred to as "room tiles".
- Memory addresses separated by a slash are Ages first, then Seasons.


## Objects

The different types of object (special objects, items, interactions, enemies,
and parts) have their own routines for loading/setting graphics and animations,
but the fundamentals are the same. Objects are differentiated from room tiles
in that they're displayed as sprites instead of background, and they have more
detailed state information, although room tiles can "become" (or rather create)
objects when interacted with, and vice versa. Loaded object state is stored
$1:d000-dfff, with each object getting a $40-byte block. General description
of objects is beyond the scope of this document; just know that the information
here is generally about them.

A "special object" is one of: Link, Maple, Ricky, Dimitri, Moosh, a minecart
that Link is riding, or the raft from Ages. Link is always loaded, and only one
other special object can be loaded at a time.


## Animations

An animation is described as a sequence of three-byte entries of the form
{duration, index, parameter}. The duration is the number of display frames
before the animation advances to the next entry (assuming no other event
happens to change the object's state); the index is used to determine which
graphics to load for that animation frame. If $ff is the first byte in the
entry, the pointer to the animation data will be shifted by -$100 plus the
value of the next byte. This is used to loop animations.

The duration and parameter are loaded into Object.animCounter and
.animParameter, and the address of the next animation entry is loaded into
.animPointer. The index is used as an offset into the tile and OAM pointer
tables for that Object.id.

In Seasons, Link's animation data is in bank 6. The address for a special
object animation is calculated starting at $6:4420.


## OAM 

OAM, or Object Attribute Memory, describes how tile data should be displayed
onscreen. $fe00-$fe9f hold the OAM data currently in use.

An object has "base" OAM values stored in Object.oamFlagsBackup, .oamFlags, and
.oamTileIndexBase, along with its .x, .y, and .z coordinates. These values are
taken together with the values pointed to by Object.oamDataAddress (determined
by using its current animation's index into a pointer table) to determine the
final values that are written to OAM.

The first byte in an OAM data entry in ROM equals the number of 8x16 sprites
that are to be drawn; 4 bytes for each sprite follow. The first is Y offset
relative to the object's nominal display position, the second is X offset, the
third is tile index offset (added to Object.oamTileIndexBase), and the fourth
is xor'ed with .oamFlags, which determine palette, flipping, VRAM bank, and
object/background priority. 

In Seasons, Link's OAM data is stored in bank $12, and address calculation for
special object OAM data begins at $6:4523. Base OAM data for special objects is
set based purely on object ID.


## Tiles

Tiles are stored in VRAM at $0:8000-97ff and $1:8000-97ff. Each tile is 8x8
pixels; this means that each sprite is composed of two tiles. The format is 2
bits per pixel, with two bytes per line: the first byte determines the low
bits, and the second determines the high bits.

Tiles for most(?) interactions are stored compressed, and usually all the tiles
for an object are in VRAM at once. The object then uses the OAM data indicated
by its animation frames in order to switch between those tiles for different
animations. By contrast, Link's tiles are stored uncompressed, and only 4 tiles
are loaded for Link at a time—storing all Link's tiles in VRAM at once would
not be practical (and maybe not even possible), and decompressing tile data
every time Link's animation changes could be costly for the CPU.

In Seasons, Link's tiles are stored in bank $1a, starting at the very
beginning. DMA transfers for special object tiles are queued at $6:4fff.
Tile address calculation for special objects begins at $6:4514.


## DMA transfers

VRAM DMA (Direct Memory Access) is a hardware feature allowing the GB/C to copy
data to VRAM (for OAM, tiles, and more) on behalf of and independent of
whatever the CPU is doing. There are different types of DMA transfer and the
details aren't important for the purposes of this document, but in general data
is not directly copied to VRAM from ROM by code; instead a queue of DMA
transfers is maintained so that the main "thread" of execution does not have to
wait for VRAM to become accessible during VBlank (etc?) before continuing.

The queueDmaTransfer function at $0:058a/0566 is useful for seeing when and
what data is transferred to VRAM (or sometimes WRAM). See oracles-disasm/ages.s
for details.


## Link sprite replacement

A common misconception seems to be that the game "prevents" Link from using
items while transformed, and that if the restriction could be lifted, then Link
would be able to do the things he normally does while using the transformed
graphics. In reality, transformations are different object IDs from regular
Link, and therefore don't have animation or OAM data corresponding to anything
that Link does except walk (Link's ID reverts to normal when doing things like
swimming, or being in shops).

It is possible to "manually" override Link's effective ID for the purposes of
display only and coerce his animation indexes into the 0-7 range used by
transformations, but even with "perfect" coercion (i.e. at the very least Link
is facing the proper direction during all his different actions) things can get
very silly. Animations like swimming, riding minecarts, and riding animals use
specific OAM and tile data; without it Link's sprite would simply be drawn on
top of the water / minecart / animal with an incorrect positional offset (and
incorrect palette in the case of charging an attack with Ricky or Moosh).

In summary: the most practical way to adequately replace Link's graphics would
be to replace the tiles themselves. This is technically unchallenging for
two-sprite objects since Link's graphics data is uncompressed, but the volume
of tiles that would need to be substituted is large, and many of Link's tiles
cannot be adequately matched or generated from the existing tiles of another
character.  Most mobile NPCs have 8 animation frames total, while Link has
about 20 just for riding Moosh (although animation frames and tiles do not
correspond to each other on a 1:4 basis; a complete walking animation usually
uses 16 8x8 tiles for its 8 16x16 frames and achieves the other half by
flipping the tiles).

Other problems:

- Some of Link's animations assume that Link will be symmetrical about the Y
  axis.
- Many Link-sized sprites use a third or fourth sprite for headgear, etc, as in
  the case of Ralph, Din, Stockwell (the shopkeeper), Zoras, Piratians, and
  others.
