package logic

// subrosia has several large areas which are guaranteed to be traverseable as
// long as you can get there in the first place:
//
// 1. "temple": rosa portal, dance hall, temple, smithy
// 2. "beach": swamp portal, market, beach
// 3. "hide and seek": H&S, mountain portal, spring tower
// 4. "pirate house": village portal, pirates
// 5. "furnace": lake portal, furnace, bomb flower
// 6. "bridge": bridge area (large but not visited in any%)
//
// the other locations are isolated and only traverseable with some combination
// of jumping and boulder removal.

var subrosiaNodes = map[string]*Node{
	"temple": Or("rosa portal",
		And("beach", "ribbon"),
		And("beach", "jump 2"),
		And("hide and seek", "bomb jump 4"),
		And("bridge", "jump 2")),

	"beach": Or("swamp portal",
		And("hide and seek", "jump 2", "bracelet",
			Or("jump 3", "magnet gloves", Hard("bombs"))),
		And("furnace", "bracelet", "jump 2"),
		And("furnace", Or("jump 4", Hard("jump 3"))),
		And("furnace", "jump 2", "magnet gloves"),
		And("temple", "jump 2")),

	"hide and seek": Or("mountain portal",
		And("pirate house", "jump 2"),
		And("bomb jump 4", Or("temple", "bridge"))),

	"pirate house": Or("village portal", And("hide and seek", "jump 2")),

	"furnace": Or("lake portal",
		And("beach", Or("jump 4", Hard("jump 3"))),
		And("beach", "magnet gloves", "jump 2")),

	"bridge": Or(
		And("temple", "jump 2"),
		And("remains portal", "bracelet", "bomb jump 3"),
		And("hide and seek", "bomb jump 4")),

	"dance hall prize": AndSlot("temple"),
	"rod gift":         AndSlot("temple"),
	"star ore spot":    AndSlot("beach", "shovel"),
	"pirate's bell":    And("temple", "rusty bell"),
	"winter tower":     AndSlot("temple", Or("hit far switch", "jump 2")),
	"summer tower":     AndSlot("beach", "ribbon", "bracelet"),
	"spring tower":     AndSlot("hide and seek", "jump 2"),
	"autumn tower":     AndSlot("temple", "jump 2", "bomb flower"),
	"blue ore chest": AndSlot("hide and seek",
		Or("jump 4", "magnet gloves")),
	"red ore chest": AndSlot("furnace", "jump 2",
		Or("jump 4", "magnet gloves")),
	"non-rosa gasha chest": AndSlot("bridge"),
	"rosa gasha chest":     AndSlot("beach", "ribbon", "jump 2"),
	"subrosian market 1":   AndSlot("beach", "star ore"),
	"subrosian market 2":   AndSlot("beach", "ore chunks", "ember seeds"),
	"subrosian market 5":   AndSlot("beach", "ore chunks"),
	"hard ore slot": AndSlot("furnace", "red ore", "blue ore",
		"temple", "bomb flower"),
	"iron shield gift": AndSlot("temple", "hard ore"),

	"enter d8": Or("d8 portal"),
}
