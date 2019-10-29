# Seasons

## Important todos
- tarm entrances either shouldn't lead to dungeons, or should require gale seeds to be in logic
- mt cucco top cucco - add another one that travels the other way
- non-random desert quicksand can be a potential soft lock
- reconsidering resetting season on exit to overworld

## Soft lock warnings
Most of these checks can be easily done if seasons reset on overworld exit
- eastern suburbs check cliff (up the springbloom flower)
- upper windmill currently set as 1-way to avoid soft locks, but can change back
- jumping down from d5, and from d5 area and portal to lake
- jumping down from mt cucco - down the springbloom flower
- jumping down from west coast graveyard
- jumping down from d6 area (springbloom flower and snowpath)
- jumping down from tarm 1st deku area
- Additionally, re-checking checks for existing warnings

## Smoother experience 
- currently Moblin keep outer entrances are just not considered in logic - have moblin keep destroyed flag put 2 stairs in those outer entrances
- fix sprites in some subrosia houses, notably smithy and minigame  
- (if adding minigame into rando, need to fix exit warping to vanilla location)
- correct rupee/ore chunk display
- ore chunk farming logic

## Future ideas
- hard logic idea - 2d spring tower section with bombs/seeds
- hide and seek without feather
- check d2 deku scrub - for logical bombs
- add back in the d2 alternate entrances (how will it work with essences)

## Relating to shared items
- getSharedItemIds - better failure option
- seasons_slots.yaml - accurate room values for standing items
- winter woods hp - animal usage?

# Ages
- Soft locks

# Common
- tree warp in logic
- soft reset in logic (for one-ways, eg bomb caves, ember bushes and keys)
- Press up in doors specifically to go back through them?
- put gale seeds into logic in more places so that you can jump down usual soft lock cliffs 
- test state-related rooms, eg Black Tower