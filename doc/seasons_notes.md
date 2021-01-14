# Oracle of Seasons randomizaton notes

These notes apply specifically to Seasons, alongside the game-neutral ones in
[README.md](https://github.com/jangler/oracles-randomizer/blob/master/README.md).


## Randomization

- The default season for each area is randomized, with the exception of regions
  that have only one season anyway.
- The rod of seasons is broken into four items (one for each season). Obtaining
  a season gives you the rod as well.
- Fool's ore is randomized, since it's actually the most powerful weapon in the
  game. The Strange Brothers trade you nothing for your feather (or cape), and
  shovel is not required to retrieve the stolen item.

The following items are **not** randomized:

- Pirate's bell (obtained by polishing the rusty bell)
- Subrosian dancing prizes after the first
- Subrosian hide and seek items
- Trading sequence items, since the vanilla trading sequence reward is
  technically just a text box and not an item.


## Other notable changes

- The intro sequence is removed, dropping you just north of Horon Village
  at the start of the game instead.
- Mystical seeds grow in all seasons.
- Rosa doesn't appear in the overworld, and her portal is activated by default.
- In some situations, the game will give you a warning about what you're doing
  or about to do. **If you receive one, what you're doing is out of logic and
  could potentially lead to a softlock**—but in some cases you can also be fine
  as long as you're careful.


## Logic

"Logic" means the set of plays that the randomizer may require you to make in
order to progress in the game. Anything the vanilla game requires you to do is,
of course, in logic. Beyond that and some easy alternate tactics, the rules in
this document apply.


### Normal logic

In logic:

- Jumping over 2 tiles with feather, 3 tiles with feather + pegasus seeds, 4
  tiles with cape, or 6 tiles with cape + pegasus seeds (subtract 1 tile if the
  jump is over water or lava instead of a pit)
- Using only ember seeds from the seed satchel as a weapon
- Using only ember, scent, or gale seeds from the slingshot as a weapon
- Using expert's ring as a weapon
- Using thrown objects (bushes, pots) as weapons
- Flipping spiked beetles using the shovel
- Farming ore chunks

Out of logic:

- Anything that gives an explicit in-game warning (these are also potential
  softlocks)
- Required damage (damage boosts, falling in pits)
- Farming rupees
- Getting a new type of mystical seed from anything other than a seed tree
- Getting initial bombs from regenerating plants
- Lighting torches using mystery seeds
- Using only mystery seeds as a weapon
- Using only scent or gale seeds from the seed satchel as a weapon
- Using only bombs as a weapon, except versus enemies immune to sword
- Using fist ring as a weapon
- Carrying bushes or pots between rooms
- Poe skip
- Fighting Frypolar without slingshot
- D8 sidescrollers without cape
- Magic rings from non-randomized locations
- Linked secrets


### Hard logic

Choosing hard difficulty enables things that are out of normal logic, with the
exception of:

- Warnings
- Farming rupees, except by shovel RNG manips
- Lighting more than two torches per room using mystery seeds
- Mystery seeds as a weapon
- Bombs as a weapon, except versus enemies immune to sword
- Double damage boosts
- Trading seeds in Subrosia Market without having a seed item
- Swimming against currents without Swimmer's Ring
- Magic rings from non-randomized locations
- Linked secrets

See
[seasons_hard_guide.md](https://github.com/jangler/oracles-randomizer/blob/master/doc/seasons_hard_guide.md)
for more information on specific tricks in hard logic.
