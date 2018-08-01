#!/usr/bin/env python3

import argparse
import struct
import sys

# because of pyyaml's serialization, we're going to be using lists instead of
# tuples in this script.
import yaml


parser = argparse.ArgumentParser(description="read data from an oos rom.",
        epilog="""
if action is "getroom", two additional hex integer parameters must be
privided for the group ID and room ID of a specific room to get data
from.

if action is "searchchests", an optional hex integer parameters may be
provided to limit the search to a given group ID and music ID.

if action is "searchobjects", two additional hex integer parameters must
be provided for the interaction mode and ID of the objects to search for.
""".strip())
parser.add_argument("romfile", type=str, help="file path of rom to read")
parser.add_argument("action", type=str, help="type of operation to perform")
parser.add_argument("args", type=str, nargs="*", help="action parameters")
args = parser.parse_args()


def fatal(*args):
    print("%s: error:" % __file__, *args, file=sys.stderr)
    exit(2)


def dict_presenter(dumper, data):
    return dumper.represent_dict(data.items())

def hexint_presenter(dumper, data):
    return dumper.represent_int('0x%02x' % data)

yaml.add_representer(dict, dict_presenter)
yaml.add_representer(int, hexint_presenter)


def full_addr(bank_num, offset):
    if bank_num > 2:
        return 0x4000 * (bank_num - 1) + offset
    return offset


MUSIC_PTR_TABLE = 0x04, 0x483c
OBJECT_PTR_TABLE = 0x11, 0x5b38
CHEST_PTR_TABLE = 0x15, 0x53af

MUSIC = {
    0x03: "overworld",
    0x0a: "horon village",
    0x0d: "essence room",
    0x0e: "house",
    0x0f: "fairy fountain",
    0x12: "hero's cave",
    0x13: "gnarled root dungeon",
}

INTERACTION_MODES = {
    0xf1: "NV interaction",
    0xf2: "DV interaction",
    0xf6: "random entities",
    0xf7: "specific entity",
}

NV_INTERACTIONS = {}

ENTITIES = {
    0x09: ("octorok", {
        0x00: "red 0x00",
        0x01: "red 0x01",
    }),
    0x0a: ("goriya", {
        0x00: "boomerang",
    }),
    0x0e: ("trap", {
        0x00: "spinner",
        0x01: "blade",
    }),
    0x31: ("stalfos", {
        0x00: "blue",
    }),
    0x32: ("keese", {
        0x00: "normal",
    }),
    0x34: ("zol", {
        0x01: "red",
    }),
    0x35: ("floormaster", {}),
    0x38: ("great fairy", {}),
    0x43: ("gel", {}),
    0x53: ("dragonfly", {}),
    0x59: ("fixed drop", {
        0x00: "fairy",
        0x05: "ember seeds",
    }),
    0x70: ("goriya bros", {}),
    0x78: ("aquamentus", {}),
}

DV_INTERACTIONS = {
    0x12: ("dungeon", {
        0x00: "entry text",
        0x01: "small key when room cleared",
        0x02: "chest when room cleared",
        0x04: "stairs when room cleared",
    }),
    0x13: ("push block trigger", {}),
    0x1e: ("doors", {
        0x04: "N opens on trigger",
        0x08: "N opens when room cleared",
        0x09: "E opens when room cleared",
        0x0a: "S opens when room cleared",
        0x0b: "W opens when room cleared",
        0x14: "N opens for torches",
        0x15: "W opens for torches",
    }),
    # 0x20 0x00 used in d1 mid boss room
    # 0x20 0x01 used for button -> small key chest in d0
    # 0x20 0x02 used for button -> small key chest in d1
    # 0x20 0x03 used for boss room in d1
    # 0x21 and 0x22 are outside the d1 entrance
    # 0x25 0x00 and 0x01 are on the cat-stuck-in-tree screen
    # 0x26 0x00 and 0x01 are also on the cat-stuck-in-tree screen
    0x38: ("d1 old man", {}),
    # 0x44 0x09 is in impa's house
    0x6b: ("placed item", {
        0x1f: "gasha seed",
        0x20: "seed satchel",
    }),
    0x78: ("toggle tile", {}),
    0x7e: ("miniboss portal", {}),
    0x7f: ("essence", {}),
    0x9d: ("impa", {}),
    0xc6: ("wooden sword", {}),
    0xdc: ("warp", {
        0x01: "doorway",
        0x02: "chimney",
    }),
    # 0xa5 0x09 used on screen where link falls in the intro
    # 0xdc 0x01 and 0x02 outside hero's cave. entrance ??
    0xe2: ("statue eyes", {}),
}

TREASURES = {
    0x03: ("bombs", {
        0x00: "10 count",
    }),
    0x28: ("rupees", {
        0x03: "20 count",
        0x04: "30 count",
    }),
    0x2d: ("ring", {
        0x04: "discovery ring",
    }),
    0x30: ("small key", {
        0x03: "in chest",
    }),
    0x31: ("boss key", {
        0x03: "in chest",
    }),
    0x32: ("compass", {
        0x02: "in chest",
    }),
    0x33: ("dungeon map", {
        0x02: "in chest",
    }),
    0x34: ("gasha seed", {}),
}


def lookup_entry(table, entry_id, param):
    if entry_id in table:
        entry = table[entry_id]

        if param in entry[1]:
            return entry[0], entry[1][param]

        return entry[0], param

    return entry_id, param


def read_byte(buf, bank, addr, increment=0):
    if increment > 0:
        return buf[full_addr(bank, addr)], addr + increment

    return buf[full_addr(bank, addr)]


def read_ptr(buf, bank, addr):
    return struct.unpack_from('<H', buf, offset=full_addr(bank, addr))[0]


def read_music(buf, group, room, name=True):
    bank, addr = MUSIC_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2) + room

    value = read_byte(buf, bank, addr)
    if name:
        if value in MUSIC:
            return MUSIC[value]

    return value


def read_objects(buf, group, room, name=True):
    # read initial pointer
    bank, addr = OBJECT_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2) + room * 2
    addr = read_ptr(buf, bank, addr)

    # read objects (recursively if more pointers are involved)
    objects = []
    while read_byte(buf, bank, addr) != 0xff:
        new_objects, addr = read_interaction(buf, bank, addr, name)
        objects += new_objects

    return objects


def loop_read_interaction(buf, bank, addr, name=True):
    objects = []

    while read_byte(buf, bank, addr) not in (0xfe, 0xff):
        new_objects, addr = read_interaction(buf, bank, addr, name)
        objects += new_objects

    return objects, addr


def read_interaction(buf, bank, addr, name=True):
    objects = []

    # read interaction type
    mode, addr = read_byte(buf, bank, addr, 1)

    if mode == 0xf0:
        print("skipped interaction type", hex(mode), "@", hex(addr - 1),
                file=sys.stderr)
        # TODO
        while read_byte(buf, bank, addr) < 0xf0:
            addr += 1
    elif mode == 0xf1:
        # "no-value interaction"
        nv_interactions = []
        while read_byte(buf, bank, addr) < 0xf0:
            if name:
                kind = list(lookup_entry(NV_INTERACTIONS,
                        read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)))
            else:
                kind = [read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)]
            addr += 2

            objects.append({
                "mode": "NV interaction" if name else mode,
                "variety": kind,
            })
    elif mode == 0xf2:
        # "double-value interaction"
        dv_interactions = []
        while read_byte(buf, bank, addr) < 0xf0:
            if name:
                kind = list(lookup_entry(DV_INTERACTIONS,
                        read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)))
            else:
                kind = [read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)]
            addr += 2

            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append({
                "mode": "DV interaction" if name else mode,
                "variety": kind,
                "coords": [x, y],
            })
    elif mode in (0xf3, 0xf4, 0xf5):
        # pointer to other interaction
        ptr = read_ptr(buf, bank, addr)
        addr += 2
        new_objects, _ = loop_read_interaction(buf, bank, ptr, name)
        objects += new_objects
    elif mode == 0xf6:
        # randomly placed entities
        count = read_byte(buf, bank, addr) >> 5
        param = read_byte(buf, bank, addr) & 0x0f
        addr += 1

        if name:
            kind = list(lookup_entry(ENTITIES,
                    read_byte(buf, bank, addr), read_byte(buf, bank, addr+1)))
        else:
            kind = [read_byte(buf, bank, addr), read_byte(buf, bank, addr+1)]
        addr += 2

        objects.append({
            "mode": "random entities" if name else mode,
            "count": count,
            "param": param,
            "variety": kind,
        })
    elif mode == 0xf7:
        # specifically placed entities
        param, addr = read_byte(buf, bank, addr, 1)

        entities = []
        while read_byte(buf, bank, addr) < 0xf0:
            if name:
                kind = list(lookup_entry(ENTITIES,
                        read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)))
            else:
                kind = [read_byte(buf, bank, addr),
                        read_byte(buf, bank, addr+1)]
            addr += 2

            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append({
                "mode": "specific entity" if name else mode,
                "param": param,
                "variety": kind,
                "coords": [x, y]
            }),
    elif mode in (0xf8, 0xf9, 0xfa):
        # TODO
        print("skipped interaction type", hex(mode), "@", hex(addr - 1),
                file=sys.stderr)
        while read_byte(buf, bank, addr) < 0xf0:
            addr += 1
    elif mode == 0xfe:
        # end data at pointer
        addr += 1
    elif mode == 0xff:
        # no more interactions to read
        pass
    else:
        print("unknown interaction type", hex(mode), "@", hex(addr - 1),
                file=sys.stderr)

    return objects, addr


def read_chest(buf, group, room):
    # read initial pointer
    bank, addr = CHEST_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2)

    # loop through group chests until marker 0xff is reached.
    # that byte must be used for something else too, but i don't know what.
    while True:
        info, chest_room, treasure_id, treasure_subid = \
                buf[full_addr(bank, addr):full_addr(bank, addr)+4]
        if info == 0xff:
            break

        if chest_room == room:
            return list(lookup_entry(TREASURES, treasure_id, treasure_subid))

        addr += 4

    return None


def get_chests(buf, group):
    bank, addr = CHEST_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2)

    # loop through group chests until marker 0xff is reached
    chests = []
    while True:
        info, room, treasure_id, treasure_subid = \
            buf[full_addr(bank, addr):full_addr(bank, addr+4)]
        if info == 0xff:
            break

        chests.append({
            "group": group,
            "room": room,
            "music": read_music(rom, group, room, name=False),
            "treasure": list(lookup_entry(TREASURES,
                    treasure_id, treasure_subid))
        })

        addr += 4

    return chests


def name_objects(objects):
    for obj in objects:
        obj["mode"] = INTERACTION_MODES[obj["mode"]]
        if obj["mode"] == "NV interaction":
            obj["variety"] = list(lookup_entry(NV_INTERACTIONS,
                    *obj["variety"]))
        elif obj["mode"] in ("random entities", "specific entity"):
            obj["variety"] = list(lookup_entry(ENTITIES, *obj["variety"]))
        elif obj["mode"] in "DV interaction":
            obj["variety"] = list(lookup_entry(DV_INTERACTIONS,
                    *obj["variety"]))


with open(args.romfile, "rb") as f:
    rom = f.read()

if args.action == "getroom":
    if len(args.args) != 2:
        fatal("getroom expects 2 args, got", len(args.args))

    group = int(args.args[0], 16)
    room = int(args.args[1], 16)

    room_data = {
        "group": group,
        "room": room,
        "music": read_music(rom, group, room),
        "objects": read_objects(rom, group, room),
        "chest": read_chest(rom, group, room),
    }

    yaml.dump(room_data, sys.stdout)
elif args.action == "searchchests":
    if len(args.args) == 0: # all groups
        chests = []
        for group in range(8):
            chests += get_chests(rom, group)
    elif len(args.args) == 1: # specific group
        group = int(args.args[0], 16)
        chests = get_chests(rom, group)
    elif len(args.args) == 2: # specific group and music
        group = int(args.args[0], 16)
        music = int(args.args[1], 16)

        chests = get_chests(rom, group)

        # filter by music and print music name if possible
        chests = [chest for chest in chests if chest["music"] == music]
        for chest in chests:
            if chest["music"] in MUSIC:
                chest["music"] = MUSIC[chest["music"]]
    elif len(args.args) > 2:
        fatal("searchchests expects 0-2 args, got", len(args.args))

    yaml.dump(chests, sys.stdout)
elif args.action == "searchobjects":
    if len(args.args) != 2:
        fatal("searchobjects expects 2 args, got", len(args.args))

    mode = int(args.args[0], 16)
    obj_id = int(args.args[1], 16)

    # read all interactions in all rooms in all groups, and collate the
    # accumulated objects that match the given ID.
    objects = []
    for group in range(8):
        bank, addr = OBJECT_PTR_TABLE
        addr = read_ptr(rom, bank, addr + group * 2)

        # loop through rooms until the high byte is fxxx, which means that the
        # interaction pointers have ended and the interaction data has started
        room = 0
        while True:
            room_addr = read_ptr(rom, bank, addr + room * 2)
            if room > 0xff or room_addr > 0xf000:
                break

            # read objects (recursively if more pointers are involved)
            room_objects = []
            while read_byte(rom, bank, room_addr) != 0xff:
                new_objects, room_addr = read_interaction(
                        rom, bank, room_addr, name=False)
                room_objects += new_objects

            for obj in room_objects:
                if obj["mode"] == mode and obj["variety"][0] == obj_id:
                    full_obj = {
                        "group": group,
                        "room": room,
                        "music": read_music(rom, group, room),
                    }
                    full_obj.update(obj)
                    objects.append(full_obj)

            room += 1

    name_objects(objects)

    yaml.dump(objects, sys.stdout)
else:
    fatal("unknown action:", args.action)
