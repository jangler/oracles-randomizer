package rom

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

var globalRomBanks *romBanks // TODO: get rid of this with the other globals

// return e.g. "\x2d\x79" for 0x792d
func addrString(addr uint16) string {
	return string([]byte{byte(addr), byte(addr >> 8)})
}

type romBanks struct {
	endOfBank []uint16
	assembler *assembler
}

// used for unmarshaling asm data from yaml.
type asmData struct {
	filename string
	Defines  map[string]uint16
	Appends  map[byte][]map[string]string
	Replaces map[byte][]AsmReplacement
	FreeCode map[string]string `yaml:"freeCode"`
}

// also loaded from yaml, then converted to yaml.
type metaAsmData struct {
	filename string
	Common   yaml.MapSlice
	Floating yaml.MapSlice
	Seasons  yaml.MapSlice
	Ages     yaml.MapSlice
}

type AsmReplacement struct {
	Addr     uint16
	Old, New string
}

var codeMutables = map[string]*MutableRange{}

// designates a position at which the translated asm will overwrite whatever
// else is there, and associates it with a given label (or a generated label if
// the given one is blank). if the replacement extends beyond the end of the
// bank, the EOB point is moved to the end of the replacement. if the bank
// offset of `addr` is zero, the replacement will start at the existing EOB
// point.
func (r *romBanks) replaceAsm(addr Addr, label, asm string) {
	if data, err := r.assembler.compile(asm); err == nil {
		r.replaceRaw(addr, label, data)
	} else {
		exitWithAsmError(label, err)
	}
}

// as replaceAsm, but interprets the data as a literal byte string.
func (r *romBanks) replaceRaw(addr Addr, label, data string) {
	if addr.offset == 0 {
		addr.offset = r.endOfBank[addr.bank]
	}

	if label == "" {
		label = fmt.Sprintf("replacement at %02x:%04x", addr.bank, addr.offset)
	}

	if end := addr.offset + uint16(len(data)); end > r.endOfBank[addr.bank] {
		if end > 0x8000 {
			panic(fmt.Sprintf("not enough space for %s in bank %02x",
				label, addr.bank))
		}
		r.endOfBank[addr.bank] = end
	}

	codeMutables[label] = &MutableRange{
		Addrs: []Addr{addr},
		New:   []byte(data),
	}
	r.assembler.define(label, addr.offset)
}

// exit with an error code and an assembler error message.
func exitWithAsmError(funcName string, err error) {
	fmt.Fprintf(os.Stderr, "assembler error in %s:\n%v\n", funcName, err)
	os.Exit(1)
}

// returns an ordered slice of keys for slot names, so that identical seeds
// produce identical checksums.
func getOrderedSlotKeys() []string {
	keys := make([]string, 0, len(ItemSlots))
	for k := range ItemSlots {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// returns a byte table of (group, room, collect mode) entries for randomized
// items. a mode >7f means to use &7f as an index to a jump table for special
// cases.
func makeCollectModeTable() string {
	b := new(strings.Builder)

	for _, key := range getOrderedSlotKeys() {
		slot := ItemSlots[key]

		// use no pickup animation if item is a key outside a chest
		mode := slot.collectMode
		if mode < 0x80 && slot.Treasure != nil && slot.Treasure.id == 0x30 {
			mode &= 0xf8
		}

		if _, err := b.Write([]byte{slot.group, slot.room, mode}); err != nil {
			panic(err)
		}
		for _, groupRoom := range slot.moreRooms {
			group, room := byte(groupRoom>>8), byte(groupRoom)
			if _, err := b.Write([]byte{group, room, mode}); err != nil {
				panic(err)
			}
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// returns a byte table (group, room, ID, subID) entries for randomized small
// key drops (and other falling items, but those entries won't be used).
func makeRoomTreasureTable(game int) string {
	b := new(strings.Builder)

	for _, key := range getOrderedSlotKeys() {
		slot := ItemSlots[key]

		if slot.collectMode != collectModes["drop"] &&
			(game == GameAges || slot.collectMode != collectModes["d4 pool"]) {
			continue
		}

		// accommodate nil treasures when creating the dummy table before
		// treasures have actually been assigned.
		var err error
		if slot.Treasure == nil {
			_, err = b.Write([]byte{slot.group, slot.room, 0x00, 0x00})
		} else if slot.Treasure.id == 0x30 {
			// make small keys the normal falling variety, with no text box.
			_, err = b.Write([]byte{slot.group, slot.room, 0x30, 0x01})
		} else {
			_, err = b.Write([]byte{slot.group, slot.room,
				slot.Treasure.id, slot.Treasure.subID})
		}
		if err != nil {
			panic(err)
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// used by iterBankItems.
type bankItem struct {
	bank byte
	item map[string]string
}

// sends the items in the banks of the asmData list in sequence.
func iterBankItems(ads []*asmData) chan bankItem {
	c := make(chan bankItem)

	go func() {
		for _, ad := range ads {
			for bank, items := range ad.Appends {
				for _, item := range items {
					c <- bankItem{bank, item}
				}
			}
		}
		close(c)
	}()

	return c
}

// applies the labels and EOB declarations in the given asmData sets.
func (r *romBanks) applyAsmData(game int, ads []*asmData, metas []*metaAsmData) {
	// get preset addrs and defines
	for _, ad := range ads {
		for k, v := range ad.Defines {
			r.assembler.define(k, v)
		}
	}

	// preprocess map slices
	slices := make([]yaml.MapSlice, 0)
	for _, meta := range metas {
		if game == GameSeasons {
			slices = append(slices, meta.Common, meta.Seasons)
		} else {
			slices = append(slices, meta.Common, meta.Ages)
		}
	}

	// include free code
	freeCode := make(map[string]string)
	for _, ad := range ads {
		for k, v := range ad.FreeCode {
			freeCode[k] = v
		}
	}
	for _, meta := range metas {
		for _, item := range meta.Floating {
			k, v := item.Key.(string), item.Value.(string)
			freeCode[k] = v
		}
	}
	for item := range iterBankItems(ads) {
		for name, body := range item.item {
			if strings.HasPrefix(body, "/include") {
				funcName := strings.Split(body, " ")[1]
				item.item[name] = freeCode[funcName]
			}
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

	// make placeholders for EOB labels
	for item := range iterBankItems(ads) {
		for name := range item.item {
			r.assembler.define(name, 0)
		}
	}
	for _, slice := range slices {
		for _, item := range slice {
			k := item.Key.(string)
			if _, label := parseMetalabel(k); label != "" {
				r.assembler.define(label, 0)
			}
		}
	}

	// save original EOB boundaries
	originalEOBs := make([]uint16, 0x40)
	copy(originalEOBs, r.endOfBank)

	// write EOB asm using placeholders for labels, in order to get real addrs
	for item := range iterBankItems(ads) {
		for name, body := range item.item {
			r.replaceAsm(Addr{item.bank, 0}, name, body)
		}
	}
	for _, slice := range slices {
		for _, item := range slice {
			k, v := item.Key.(string), item.Value.(string)
			if addr, label := parseMetalabel(k); addr.offset == 0 {
				r.replaceAsm(Addr{addr.bank, 0}, label, v)
			}
		}
	}

	// reset EOB boundaries
	copy(r.endOfBank, originalEOBs)

	// rewrite EOB asm, using real addresses for labels
	for item := range iterBankItems(ads) {
		for name, body := range item.item {
			r.replaceAsm(Addr{item.bank, 0}, name, body)
		}
	}
	for _, slice := range slices {
		for _, item := range slice {
			k, v := item.Key.(string), item.Value.(string)
			if addr, label := parseMetalabel(k); addr.offset == 0 {
				r.replaceAsm(addr, label, v)
			}
		}
	}

	// make non-EOB asm replacements
	for _, ad := range ads {
		for bank, items := range ad.Replaces {
			for _, item := range items {
				r.replaceAsm(Addr{bank, item.Addr}, "", item.New)
			}
		}
	}
	for _, slice := range slices {
		for _, item := range slice {
			k, v := item.Key.(string), item.Value.(string)
			if addr, label := parseMetalabel(k); addr.offset != 0 {
				r.replaceAsm(addr, label, v)
			}
		}
	}
}

// applies the labels and EOB declarations in the given asm data files.
func (r *romBanks) applyAsmFiles(game int, oldPaths, newPaths []string) {
	ads := make([]*asmData, len(oldPaths))
	for i, path := range oldPaths {
		ads[i] = new(asmData)
		ads[i].filename = path
		if err := yaml.Unmarshal(
			FSMustByte(false, path), ads[i]); err != nil {
			panic(err)
		}
	}

	metas := make([]*metaAsmData, len(newPaths))
	for i, path := range newPaths {
		metas[i] = new(metaAsmData)
		metas[i].filename = path
		if err := yaml.Unmarshal(
			FSMustByte(false, path), metas[i]); err != nil {
			panic(err)
		}
	}

	r.applyAsmData(game, ads, metas)
}

// ShowAsm writes the disassembly of the specified symbol to the given
// io.Writer.
func ShowAsm(symbol string, w io.Writer) error {
	m := codeMutables[symbol]
	s, err := globalRomBanks.assembler.decompile(string(m.New))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stderr, "%02x:%04x: %s\n",
		m.Addrs[0].bank, m.Addrs[0].offset, symbol)
	_, err = fmt.Fprintln(w, s)
	return err
}

// returns the address and label components of a meta-label such as
// "02/openRingList" or "02/56a1/". see asm/README.md for details.
func parseMetalabel(ml string) (addr Addr, label string) {
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
		FSMustByte(false, "/rom/eob.yaml"), eobs); err != nil {
		panic(err)
	}
	return eobs[game]
}

// loads text, processes it, and applies it to matching labels.
func applyText(b []byte, game string) {
	textMap := make(map[string]map[string]string)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/rom/text.yaml"), textMap); err != nil {
		panic(err)
	}
	for label, rawText := range textMap[game] {
		if mut, ok := codeMutables[label]; ok {
			mut.New = processText(rawText)
			mut.Mutate(b)
		} else {
			println("no code label matches text label " + label)
		}
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
