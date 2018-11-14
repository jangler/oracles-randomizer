package logic

var NodeValues = map[string]int{
	"village shop 1":  -20,
	"village shop 2":  -30,
	"village shop 3":  -150,
	"member's shop 1": -300,
	"member's shop 2": -300,
	"member's shop 3": -200,

	"shop, 30 rupees":  -30,
	"shop, 150 rupees": -150,
	"wild tokay game":  -10,

	"blaino gift": -10,

	"king zora":           -300,
	"goron dance past":    -20,
	"goron dance present": -10,
	"target carts 1":      -10,
	"target carts 2":      -10,
	"shooting gallery":    -20,

	"goron mountain old man":      300,
	"western coast old man":       300,
	"holodrum plain east old man": 200,
	"horon village old man":       100,
	"north horon old man":         100,

	"tarm ruins old man":          -200,
	"woods of winter old man":     -50,
	"holodrum plain west old man": -100,

	// rng is involved; each rupee is either worth 1, 5, or 10
	"d2 rupee room": 200,
	"d6 rupee room": 150,
}

var RupeeValues = map[string]int{
	"rupees, 1":   1,
	"rupees, 5":   5,
	"rupees, 10":  10,
	"rupees, 20":  20,
	"rupees, 30":  30,
	"rupees, 50":  50,
	"rupees, 100": 100,
	"rupees, 200": 200,
}
