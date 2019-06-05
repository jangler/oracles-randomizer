package randomizer

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// implements the -plan flag: read a spoiler log and gen a seed from it,
// instead of vice versa. unspecified variables are left vanilla.

type plan struct {
	source   string
	items    map[string]string
	dungeons map[string]string
	portals  map[string]string
	seasons  map[string]string
	hints    map[string]string
}

func newPlan() *plan {
	return &plan{
		items:    make(map[string]string),
		dungeons: make(map[string]string),
		portals:  make(map[string]string),
		seasons:  make(map[string]string),
		hints:    make(map[string]string),
	}
}

var conditionRegexp = regexp.MustCompile(`(.+?) +<- (.+)`)

// loads conditions from a file in spoiler log format.
func parseSummary(path string, game int) (*plan, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	p := newPlan()
	p.source = string(b)
	section := p.items
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.Replace(line, "\r", "", 1)
		if strings.HasPrefix(line, "--") {
			switch line {
			case "-- items --", "-- progression items --",
				"-- small keys and boss keys --", "-- other items --":
				section = p.items
			case "-- dungeon entrances --":
				section = p.dungeons
			case "-- subrosia portals --":
				section = p.portals
			case "-- default seasons --":
				section = p.seasons
			case "-- hints --":
				section = p.hints
			default:
				return nil, fmt.Errorf("unknown section: %q", line)
			}
		} else {
			submatches := conditionRegexp.FindStringSubmatch(line)
			if submatches != nil {
				if submatches[1] == "null" {
					var nullKey string
					for i := 0; true; i++ {
						nullKey = fmt.Sprintf("null %d", i)
						if section[nullKey] == "" {
							break
						}
					}
					section[nullKey] = ungetNiceName(submatches[2], game)
				} else {
					section[ungetNiceName(submatches[1], game)] =
						ungetNiceName(submatches[2], game)
				}
			}
		}
	}

	return p, nil
}

// like findRoute, but uses a specified configuration instead of a random one.
func makePlannedRoute(rom *romState, p *plan) (*routeInfo, error) {
	ri := &routeInfo{
		companion: sora(rom.game, moosh, dimitri).(int), // shop is default
		entrances: make(map[string]string),
		graph:     newRouteGraph(rom),
		src:       rand.New(rand.NewSource(0)),
		usedItems: list.New(),
		usedSlots: list.New(),
	}

	// must init rings before item placement
	ringValues := make([]string, 0)
	for _, item := range p.items {
		if strings.Contains(item, " ring") {
			ringValues = append(ringValues, item)
		}
	}
	ringMap, err := rom.randomizeRingPool(ri.src, ringValues)
	if err != nil {
		return nil, err
	}

	// item slots
	fluteSet := false // error if different flutes are given
	for slot, item := range p.items {
		// use original ring names
		if ringName, ok := reverseLookup(ringMap, item); ok {
			item = ringName.(string)
		}

		// add given item/slot combo to list and graph
		if _, ok := rom.treasures[item]; !ok {
			return nil, fmt.Errorf("no such item: %s", item)
		}
		if _, ok := ri.graph[slot]; !ok {
			return nil, fmt.Errorf("no such check: %s", slot)
		}
		ri.graph[item] = newNode(item, orNode)
		if !itemFitsInSlot(ri.graph[item], ri.graph[slot]) {
			return nil, fmt.Errorf("%s doesn't fit in %s", item, slot)
		}
		ri.graph[item].addParent(ri.graph[slot])
		ri.usedItems.PushBack(ri.graph[item])
		ri.usedSlots.PushBack(ri.graph[slot])

		// set flute if applicable
		if strings.HasSuffix(item, "flute") {
			prevCompanion := ri.companion
			switch item {
			case "ricky's flute":
				ri.companion = ricky
			case "dimitri's flute":
				ri.companion = dimitri
			case "moosh's flute":
				ri.companion = moosh
			}
			if fluteSet && ri.companion != prevCompanion {
				return nil, fmt.Errorf("can't have multiple types of flute")
			}
			fluteSet = true
		}
	}

	// seasons
	if rom.game == gameSeasons {
		ri.seasons = make(map[string]byte, len(p.seasons))
		for area, season := range p.seasons {
			id := getStringIndex(seasonsById, season)
			if id == -1 {
				return nil, fmt.Errorf("invalid default season: %s", season)
			}
			if getStringIndex(seasonAreas, area) == -1 {
				return nil, fmt.Errorf("invalid season area: %s", area)
			}
			ri.seasons[area] = byte(id)
		}
	} else if len(p.seasons) != 0 {
		return nil, fmt.Errorf("ages doesn't have default seasons")
	}

	// dungeon entrances
	for entrance, dungeon := range p.dungeons {
		entrance = strings.Replace(entrance, " entrance", "", 1)
		for _, s := range []string{entrance, dungeon} {
			if s == "d0" || getStringIndex(dungeonNames[rom.game], s) == -1 {
				return nil, fmt.Errorf("no such dungeon: %s", s)
			}
		}
		ri.entrances[entrance] = dungeon
	}

	// portals
	if rom.game == gameSeasons {
		ri.portals = make(map[string]string, len(p.portals))
		for portal, connect := range p.portals {
			if _, ok := subrosianPortalNames[portal]; !ok {
				return nil, fmt.Errorf("invalid holodrum portal: %s", portal)
			}
			if _, ok := reverseLookup(subrosianPortalNames, connect); !ok {
				return nil, fmt.Errorf("invalid subrosia portal: %s", connect)
			}
			ri.portals[portal] = connect
		}
	} else if len(p.portals) != 0 {
		return nil, fmt.Errorf("ages doesn't have subrosia portals")
	}

	return ri, nil
}

// overwrites regular owl hints with planned ones.
func planOwlHints(p *plan, h *hinter, owlHints map[string]string) error {
	// sanity check first
	for owl, hint := range p.hints {
		hint = strings.Trim(hint, `"`)
		if _, ok := owlHints[owl]; !ok {
			return fmt.Errorf("no such owl: %s", owl)
		}
		if !isValidGameText(hint) {
			return fmt.Errorf("invalid hint text: %q", hint)
		}
	}

	// use hint if planned hint found, placeholder if not
	for owl := range owlHints {
		if hint, ok := p.hints[owl]; ok {
			owlHints[owl] = h.format(strings.Trim(hint, `"`))
		} else {
			owlHints[owl] = "..."
		}
	}

	return nil
}

// returns the index of s in a, or -1 if not found.
func getStringIndex(a []string, s string) int {
	for i, v := range a {
		if v == s {
			return i
		}
	}
	return -1
}
