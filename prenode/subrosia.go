package prenode

// subrosia has several large areas which are guaranteed to be traverseable as
// long as you can get there in the first place:
//
// 1. "temple": rosa portal, dance hall, temple, smithy
// 2. "beach": swamp portal, market, beach
// 3. "hide and seek": H&S, mountain portal, spring tower
// 4. "pirate house": village portal, pirates, desert portal
// 5. "furnace": lake portal, furnace, bomb flower
// 6. "bridge": bridge area (large but not visited in any%)
//
// the other locations are isolated and only traverseable with some combination
// of jumping and boulder removal.

var subrosiaPrenodes = map[string]*Prenode{
	"temple": Or("rosa portal",
		And("beach", "ribbon"),
		And("beach", "jump"),
		And("hide and seek", "pegasus jump L-2"),
		And("bridge", "jump")),

	"beach": Or("swamp portal",
		And("hide and seek", "bracelet", "long jump"),
		And("hide and seek", "jump", "bracelet", "magnet gloves"),
		And("furnace", "bracelet", "jump"),
		And("furnace", "feather L-2"),
		And("furnace", "jump", "magnet gloves"),
		And("temple", "jump")),

	"hide and seek": Or("mountain portal",
		And("pirate house", "jump"),
		And("temple", "pegasus jump L-2"),
		And("bridge", "pegasus jump L-2")),

	"pirate house": Or("village portal", "desert portal",
		And("hide and seek", "jump")),

	"furnace": Or("lake portal",
		And("beach", "feather L-2"),
		And("beach", "magnet gloves")),

	"bridge": Or(
		And("temple", "jump"),
		And("remains portal", "bracelet", "feather L-2"),
		And("hide and seek", "pegasus jump L-2")),

	"boomerang gift":     AndSlot("temple"),
	"rod gift":           AndSlot("temple"),
	"star ore spot":      AndSlot("beach", "shovel"),
	"pirate's bell":      And("temple", "rusty bell"),
	"cross winter tower": Or("hit far switch", "jump"),
	"winter tower":       And("temple", "cross winter tower"),
	"summer tower":       And("beach", "ribbon"),
	"spring tower":       And("hide and seek", "jump"),
	"autumn tower":       And("temple", "jump", "bomb flower"),

	"eruption room": Or("remains portal"),
	"enter d8":      OrStep("d8 portal"),

	// a few places are unaccounted for, but they're irrelevant for now
}
