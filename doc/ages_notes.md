# Oracle of Ages randomization notes

These notes apply specifically to Ages, alongside the game-neutral ones in
[README.md](https://github.com/jangler/oracles-randomizer/blob/master/README.md).


## Randomization

- The harp is progressive, starting with the Tune of Echoes, then Currents,
  then Ages.
- The trading sequence is removed. The second sword is in the item pool, and
  the Poe in Yoll Graveyard gives a randomized item.
- The red soldier who takes you to Ambi's palace for bombs waits in Deku Forest
  instead. Talk to him to trade mystery seeds for an item.
- The first **and** second prizes for Target Carts are randomized. Only the
  first item for each other minigame is randomized, with the exception of the
  Lynna Village Shooting Gallery, which has no randomized prize since its
  normal prize is a strange flute.
- Past and present versions of the same mystical seed tree grow the same type
  of seed.
- The secret shop items are not randomized.


## Other notable changes

- The intro sequence is removed. A chest replaces Impa on the screen where she
  normally gives the sword.
- The time portals on the screens adjacent to the Maku Tree are active
  permanently.
- The Tokays on Crescent Island do not steal your items, and the raft does not
  encounter a storm in the Sea of Storms.
- Dormant time portals are added to Nuun Highlands and Symmetry Village past in
  order to prevent softlocks.
- The dormant portal on the west side of Crescent Island present only responds
  to Currents, not Echoes, in order to prevent softlocks. The sign says so.
- Playing the Tune of Currents triggers reentry into a return portal Link is
  standing on. This is useful if you warp into a patch of bushes without a
  bush-breaking item.
- Ambi's palace courtyard is open from the beginning.


## Logic

"Logic" means the set of plays that the randomizer may require you to make in
order to progress in the game. Anything the vanilla game requires you to do is,
of course, in logic. Beyond that and some easy alternate tactics, the rules in
this document apply.


### Normal logic

In logic:

- Jumping over 2 tiles with feather, or 3 tiles with feather + pegasus seeds
  (subtract 1 tile if the jump is over water or lava instead of a pit)
- Using only ember seeds from the seed satchel as a weapon
- Using only ember, scent, or gale seeds from the seed shooter as a weapon
- Using expert's ring as a weapon
- Using thrown objects (bushes, pots) as weapons
- Flipping spiked beetles using the shovel

Out of logic:

- Required damage (damage boosts, falling in pits)
- Farming rupees, assuming you spend them optimally
- Getting initial scent seeds from the plants in D3
- Getting initial bombs from Goron Shooting Gallery or Head Thwomp
- Lighting torches using mystery seeds
- Using only mystery seeds as a weapon
- Using only scent or gale seeds from the seed satchel as a weapon
- Using only bombs as a weapon, except versus armos
- Using fist ring as a weapon
- Doing Patch's restoration ceremony without sword
- Getting a potion for King Zora from Maple
- Guard skip
- Various stupid tricks in D2, D3, D4, and D5
- Text warps
- Magic rings from non-randomized locations
- Linked secrets


### Hard logic

Choosing hard difficulty enables things that are out of normal logic, with the
exception of:

- Farming rupees, except by shovel RNG manips
- Mystery seeds as a weapon
- Text warps
- Magic rings from non-randomized locations
- Linked secrets
