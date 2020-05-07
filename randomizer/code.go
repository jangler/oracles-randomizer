package randomizer

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// loaded from yaml, then converted to asm.
type asmData struct {
	filename string
	Common   yaml.MapSlice
	Floating yaml.MapSlice
	Seasons  yaml.MapSlice
	Ages     yaml.MapSlice
}

// designates a position at which the translated asm will overwrite whatever
// else is there, and associates it with a given label (or a generated label if
// the given one is blank). if the replacement extends beyond the end of the
// bank, the EOB point is moved to the end of the replacement. if the bank
// offset of `addr` is zero, the replacement will start at the existing EOB
// point.
//
// returns the final label of the replacement.
func (rom *romState) replaceAsm(addr address, label, asm string) string {
	if data, err := rom.assembler.compile(asm); err == nil {
		return rom.replaceRaw(addr, label, data)
	} else {
		panic(fmt.Sprintf("assembler error in %s:\n%v\n", label, err))
	}
}

// as replaceAsm, but interprets the data as a literal byte string.
func (rom *romState) replaceRaw(addr address, label, data string) string {
	if addr.offset == 0 {
		addr.offset = rom.bankEnds[addr.bank]
	}

	if label == "" {
		label = fmt.Sprintf("replacement at %02x:%04x", addr.bank, addr.offset)
	} else if strings.HasPrefix(label, "dma_") && addr.offset%0x10 != 0 {
		addr.offset += 0x10 - (addr.offset % 0x10) // align to $xxx0
	}

	end := addr.offset + uint16(len(data))
	if end > rom.bankEnds[addr.bank] {
		if end > 0x8000 || (end > 0x4000 && addr.bank == 0) {
			panic(fmt.Sprintf("not enough space for %s in bank %02x",
				label, addr.bank))
		}
		rom.bankEnds[addr.bank] = end
	}

	rom.codeMutables[label] = &mutableRange{
		addr: addr,
		new:  []byte(data),
	}
	rom.assembler.define(label, addr.offset)

	return label
}

// returns a byte table of (group, room, collect mode, player) entries for
// randomized items. a mode >7f means to use &7f as an index to a jump table
// for special cases.
func makeCollectPropertiesTable(itemSlots map[string]*itemSlot) string {
	b := new(strings.Builder)
	for _, key := range orderedKeys(itemSlots) {
		slot := itemSlots[key]

		// use no pickup animation for falling small keys
		mode := slot.collectMode
		if mode == 0x29 && slot.treasure != nil && slot.treasure.id == 0x30 {
			mode &= 0xf8
		}

		if _, err := b.Write([]byte{slot.group, slot.room, mode, slot.player}); err != nil {
			panic(err)
		}
		for _, groupRoom := range slot.moreRooms {
			group, room := byte(groupRoom>>8), byte(groupRoom)
			if _, err := b.Write([]byte{group, room, mode, slot.player}); err != nil {
				panic(err)
			}
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// returns a byte table (group, room, id, subid) entries for randomized small
// key drops (and other falling items, but those entries won't be used).
func makeRoomTreasureTable(game int, itemSlots map[string]*itemSlot) string {
	b := new(strings.Builder)

	for _, key := range orderedKeys(itemSlots) {
		slot := itemSlots[key]

		if key != "maku path basement" &&
			slot.collectMode != collectModes["drop"] &&
			(game == gameAges || slot.collectMode != collectModes["d4 pool"]) {
			continue
		}

		// accommodate nil treasures when creating the dummy table before
		// treasures have actually been assigned.
		var err error
		if slot.treasure == nil {
			_, err = b.Write([]byte{slot.group, slot.room, 0x00, 0x00})
		} else if slot.treasure.id == 0x30 {
			// make small keys the normal falling variety, with no text box.
			_, err = b.Write([]byte{slot.group, slot.room, 0x30, 0x01})
		} else {
			_, err = b.Write([]byte{slot.group, slot.room,
				slot.treasure.id, slot.treasure.subid})
		}
		if err != nil {
			panic(err)
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// that's correct
type eobThing struct {
	addr         address
	label, thing string
}

// applies the labels and EOB declarations in the asm data sets.
// returns a slice of added labels.
func (rom *romState) applyAsmData(asmFiles []*asmData) []string {
	// preprocess map slices (keys = labels, values = asm blocks)
	slices := make([]yaml.MapSlice, 0)
	for _, asmFile := range asmFiles {
		if rom.game == gameSeasons {
			slices = append(slices, asmFile.Common, asmFile.Seasons)
		} else {
			slices = append(slices, asmFile.Common, asmFile.Ages)
		}
	}

	// include free code
	freeCode := make(map[string]string)
	for _, asmFile := range asmFiles {
		for _, item := range asmFile.Floating {
			k, v := item.Key.(string), item.Value.(string)
			freeCode[k] = v
		}
	}
	for _, slice := range slices {
		for name, item := range slice {
			v := item.Value.(string)
			if strings.HasPrefix(v, "/include") {
				funcName := strings.Split(v, " ")[1]
				slice[name].Value = freeCode[funcName]
			}
		}
	}

	// save original EOB boundaries
	originalBankEnds := make([]uint16, 0x40)
	copy(originalBankEnds, rom.bankEnds)

	// make placeholders for labels and accumulate EOB items
	allEobThings := make([]eobThing, 0, 3000) // 3000 is probably fine
	for _, slice := range slices {
		for _, item := range slice {
			k, v := item.Key.(string), item.Value.(string)
			addr, label := parseMetalabel(k)
			if label != "" {
				rom.assembler.define(label, 0)
			}
			if addr.offset == 0 {
				allEobThings = append(allEobThings,
					eobThing{address{addr.bank, 0}, label, v})
			}
		}
	}

	// defines (which have no labels, by convention) must be processed first
	sort.Slice(allEobThings, func(i, j int) bool {
		return allEobThings[i].label == ""
	})
	// owl text must go last
	for i, thing := range allEobThings {
		if thing.label == "owlText" {
			allEobThings = append(allEobThings[:i], allEobThings[i+1:]...)
			allEobThings = append(allEobThings, thing)
			break
		}
	}

	// write EOB asm using placeholders for labels, in order to get real addrs
	for _, thing := range allEobThings {
		rom.replaceAsm(thing.addr, thing.label, thing.thing)
	}

	// also get labels for labeled replacements
	for _, slice := range slices {
		for _, item := range slice {
			addr, label := parseMetalabel(item.Key.(string))
			if addr.offset != 0 && label != "" {
				rom.assembler.define(label, addr.offset)
			}
		}
	}

	// reset EOB boundaries
	copy(rom.bankEnds, originalBankEnds)

	labels := make([]string, 0, 3000) // 3000 probably still fine

	// rewrite EOB asm, using real addresses for labels
	for _, thing := range allEobThings {
		labels = append(labels,
			rom.replaceAsm(thing.addr, thing.label, thing.thing))
	}

	// make non-EOB asm replacements
	for _, slice := range slices {
		for _, item := range slice {
			k, v := item.Key.(string), item.Value.(string)
			if addr, label := parseMetalabel(k); addr.offset != 0 {
				labels = append(labels, rom.replaceAsm(addr, label, v))
			}
		}
	}

	return labels
}

// applies the labels and EOB declarations in the given asm data files.
func (rom *romState) applyAsmFiles(infos []os.FileInfo) {
	asmFiles := make([]*asmData, len(infos))

	// standard files are embedded
	for i, info := range infos {
		asmFiles[i] = new(asmData)
		asmFiles[i].filename = info.Name()

		// readme etc
		if !strings.HasSuffix(info.Name(), ".yaml") {
			continue
		}

		path := "/asm/" + info.Name()
		if err := yaml.Unmarshal(
			FSMustByte(false, path), asmFiles[i]); err != nil {
			panic(err)
		}
	}

	rom.applyAsmData(asmFiles)
}

// showAsm writes the disassembly of the specified symbol to the given
// io.Writer.
func (rom *romState) showAsm(symbol string, w io.Writer) error {
	mut := rom.codeMutables[symbol]
	if mut == nil {
		return fmt.Errorf("no such label: %s", symbol)
	}
	s, err := rom.assembler.decompile(string(mut.new))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%02x:%04x: %s\n",
		mut.addr.bank, mut.addr.offset, symbol)
	_, err = fmt.Fprintln(w, s)
	return err
}

// returns the address and label components of a meta-label such as
// "02/openRingList" or "02/56a1/". see asm/README.md for details.
func parseMetalabel(ml string) (addr address, label string) {
	switch tokens := strings.Split(ml, "/"); len(tokens) {
	case 1:
		fmt.Sscanf(ml, "%s", &label)
	case 2:
		fmt.Sscanf(ml, "%x/%s", &addr.bank, &label)
	case 3:
		fmt.Sscanf(ml, "%x/%x/%s", &addr.bank, &addr.offset, &label)
	default:
		panic("invalid metalabel: " + ml)
	}
	return
}

// returns a $40-entry slice of addresses of the ends of rom banks for the
// given game.
func loadBankEnds(game string) []uint16 {
	eobs := make(map[string][]uint16)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/eob.yaml"), eobs); err != nil {
		panic(err)
	}
	return eobs[game]
}

// loads text, processes it, and attaches it to matching labels.
func (rom *romState) attachText() {
	// load initial text
	textMap := make(map[string]map[string]string)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/text.yaml"), textMap); err != nil {
		panic(err)
	}
	for label, rawText := range textMap[gameNames[rom.game]] {
		if mut, ok := rom.codeMutables[label]; ok {
			mut.new = processText(rawText)
		} else {
			panic(fmt.Sprintf("no code label matches text label %q", label))
		}
	}

	// insert randomized item names into shop text
	shopNames := loadShopNames(gameNames[rom.game])
	shopMap := map[string]string{
		"shopFluteText": "shop, 150 rupees",
	}
	if rom.game == gameSeasons {
		shopMap["membersShopSatchelText"] = "member's shop 1"
		shopMap["membersShopGashaText"] = "member's shop 2"
		shopMap["membersShopMapText"] = "member's shop 3"
		shopMap["marketRibbonText"] = "subrosia market, 1st item"
		shopMap["marketPeachStoneText"] = "subrosia market, 2nd item"
		shopMap["marketCardText"] = "subrosia market, 5th item"
	}
	for codeName, slotName := range shopMap {
		code := rom.codeMutables[codeName]
		itemName := shopNames[rom.itemSlots[slotName].treasure.displayName]
		code.new = append(code.new[:2],
			append([]byte(itemName), code.new[2:]...)...)
	}
}

var hashCommentRegexp = regexp.MustCompile(" #.+?\n")

// processes a raw text string as a go string literal, converting escape
// sequences to their actual values. "comments" and literal newlines are
// stripped.
func processText(s string) []byte {
	var err error
	s = hashCommentRegexp.ReplaceAllString(s, "")
	s, err = strconv.Unquote("\"" + s + "\"")
	if err != nil {
		panic(err)
	}
	return []byte(s)
}

var articleRegexp = regexp.MustCompile("^(an?|the) ")

// return a map of internal item names to text that should be displayed for the
// item in shops.
func loadShopNames(game string) map[string]string {
	m := make(map[string]string)

	// load names used for owl hints
	itemFiles := []string{
		"/hints/common_items.yaml",
		fmt.Sprintf("/hints/%s_items.yaml", game),
	}
	for _, filename := range itemFiles {
		if err := yaml.Unmarshal(
			FSMustByte(false, filename), m); err != nil {
			panic(err)
		}
	}

	// remove articles
	for k, v := range m {
		m[k] = articleRegexp.ReplaceAllString(v, "")
	}

	return m
}

// set up all the pre-randomization asm changes, and track the state so that
// the randomization changes can be applied later.
func (rom *romState) initBanks() {
	rom.codeMutables = make(map[string]*mutableRange)
	rom.bankEnds = loadBankEnds(gameNames[rom.game])
	asm, err := newAssembler()
	if err != nil {
		panic(err)
	}
	rom.assembler = asm

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	roomTreasureBank := byte(sora(rom.game, 0x3f, 0x38).(int))
	numOwlIds := sora(rom.game, 0x1e, 0x14).(int)
	rom.replaceRaw(address{0x06, 0}, "collectPropertiesTable",
		makeCollectPropertiesTable(rom.itemSlots))
	rom.replaceRaw(address{roomTreasureBank, 0}, "roomTreasures",
		makeRoomTreasureTable(rom.game, rom.itemSlots))
	rom.replaceRaw(address{0x3f, 0}, "owlTextOffsets",
		string(make([]byte, numOwlIds*2)))

	// load all asm files in the asm/ directory.
	dir, err := FS(false).Open("/asm/")
	if err != nil {
		panic(err)
	}
	fi, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	rom.applyAsmFiles(fi)
}

// apply user-included asm files.
func (rom *romState) addIncludes() error {
	asmFiles := make([]*asmData, len(rom.includes))

	// read from filesystem
	for i, path := range rom.includes {
		asmFiles[i] = new(asmData)
		asmFiles[i].filename = path

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(asmFiles[i]); err != nil {
			return err
		}
	}

	// apply immediately
	labels := rom.applyAsmData(asmFiles)
	sort.Strings(labels)
	for _, label := range labels {
		rom.codeMutables[label].mutate(rom.data)
	}

	return nil
}
