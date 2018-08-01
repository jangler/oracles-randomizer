#!/usr/bin/env python3

import struct
import sys

# because of pyyaml's serialization, we're going to be using lists instead of
# tuples in this script.
import yaml


def hexint_presenter(dumper, data):
    return dumper.represent_int('0x%02x' % data)

yaml.add_representer(int, hexint_presenter)


def full_addr(bank_num, offset):
    if bank_num > 2:
        return 0x4000 * (bank_num - 1) + offset
    return offset


MUSIC_PTR_TABLE = 0x04, 0x483c
OBJECT_PTR_TABLE = 0x11, 0x5b38
CHEST_PTR_TABLE = 0x15, 0x53af

ENTITIES = {
    0x0a: ("goriya", {
        0x00: "boomerang",
    }),
    0x0e: ("blade trap", {
        0x01: "blue",
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
    0x59: ("fixed drop", {
        0x00: "fairy",
        0x05: "ember seeds",
    }),
}

DV_INTERACTIONS = {
    0x12: ("dungeon", {
        0x00: "entry text",
        0x01: "small key when room cleared",
        0x02: "chest when room cleared",
    }),
    0x13: ("push block trigger", {}),
    0x1e: ("doors", {
        0x0a: "S opens when room cleared",
        0x0b: "W opens when room cleared",
        0x14: "N opens for torches",
        0x15: "W opens for torches",
    }),
    0x38: ("d1 old man", {}),
    0x78: ("toggle tile", {}),
    0x7e: ("minecart?", {}),
    0xe2: ("statue eyes", {}),
}

TREASURES = {
    0x32: ("compass", {
        0x02: "in chest",
    }),
    0x33: ("dungeon map", {
        0x02: "in chest",
    }),
    0x34: ("gasha seed", {
        0x01: "in chest",
    }),
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


def read_music(buf, group, room):
    bank, addr = MUSIC_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2) + room
    return read_byte(buf, bank, addr)


def read_objects(buf, group, room):
    # read initial pointer
    bank, addr = OBJECT_PTR_TABLE
    addr = read_ptr(buf, bank, addr + group * 2) + room * 2
    addr = read_ptr(buf, bank, addr)

    # read objects (recursively if more pointers are involved)
    objects = []
    while read_byte(buf, bank, addr) != 0xff:
        new_objects, addr = read_interaction(buf, bank, addr)
        objects += new_objects

    return objects


def loop_read_interaction(buf, bank, addr):
    objects = []

    while read_byte(buf, bank, addr) not in (0xfe, 0xff):
        new_objects, addr = read_interaction(buf, bank, addr)
        objects += new_objects

    return objects, addr


def read_interaction(buf, bank, addr):
    objects = []

    # read interaction type
    mode, addr = read_byte(buf, bank, addr, 1)

    if mode == 0xf0:
        # TODO
        while read_byte(buf, bank, addr) < 0xf0:
            addr += 1
    elif mode == 0xf2:
        # "double-value interaction"
        dv_interactions = []
        while read_byte(buf, bank, addr) < 0xf0:
            kind = list(lookup_entry(DV_INTERACTIONS,
                    read_byte(buf, bank, addr), read_byte(buf, bank, addr+1)))
            addr += 2
            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append(["DV interaction", kind, [x, y]])
    elif mode in (0xf3, 0xf4, 0xf5):
        # pointer to other interaction
        ptr = read_ptr(buf, bank, addr)
        addr += 2
        new_objects, _ = loop_read_interaction(buf, bank, ptr)
        objects += new_objects
    elif mode == 0xf6:
        # randomly placed entities
        count = read_byte(buf, bank, addr) >> 5
        param = read_byte(buf, bank, addr) & 0x0f
        addr += 1

        kind = list(lookup_entry(ENTITIES,
                read_byte(buf, bank, addr), read_byte(buf, bank, addr+1)))
        addr += 2

        objects.append(["random entities", count, param, kind])
    elif mode == 0xf7:
        # specifically placed entities
        param, addr = read_byte(buf, bank, addr, 1)

        entities = []
        while read_byte(buf, bank, addr) < 0xf0:
            kind = list(lookup_entry(ENTITIES,
                    read_byte(buf, bank, addr), read_byte(buf, bank, addr+1)))
            addr += 2
            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append(["specific entity", param, kind, [x, y]])
    elif mode in (0xf8, 0xf9, 0xfa):
        # TODO
        while read_byte(buf, bank, addr) < 0xf0:
            addr += 1
    elif mode == 0xfe:
        # end data at pointer
        addr += 1
    elif mode == 0xff:
        # no more interactions to read
        pass

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


if len(sys.argv) != 4:
    print("usage: %s <romfile> <group> <room>" % __file__)
    exit(2)

with open(sys.argv[1], 'rb') as f:
    rom = f.read()

group = int(sys.argv[2], 16)
room = int(sys.argv[3], 16)

room_data = {
    "group": group,
    "room": room,
    "music": read_music(rom, group, room),
    "objects": read_objects(rom, group, room),
    "chest": read_chest(rom, group, room),
}

yaml.dump(room_data, sys.stdout)
