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

ENTITIES = {
    0x31: ("stalfos", {
        0x00: "blue",
    }),
    0x59: ("fixed drop", {
        0x00: "fairy",
        0x05: "ember seeds",
    }),
}

DV_INTERACTIONS = {
    0x12: ("dungeon", {
        0x01: "small key falls when room cleared",
    }),
    0x1e: ("doors", {
        0x14: "north shutter",
        0x0b: "open when room cleared",
    }),
}


def lookup_entity(entityID, param):
    if entityID in ENTITIES:
        entity = ENTITIES[entityID]

        if param in entity[1]:
            return entity[0], entity[1][param]

        return entity[0], hex(param)

    return hex(entityID), hex(param)


def lookup_DV(dvID, param):
    if dvID in DV_INTERACTIONS:
        dv = DV_INTERACTIONS[dvID]

        if param in dv[1]:
            return dv[0], dv[1][param]

        return dv[0], hex(param)

    return hex(dvID), hex(param)


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
            kind = lookup_DV(read_byte(buf, bank, addr),
                    read_byte(buf, bank, addr+1))
            addr += 2
            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append(("DV interaction", kind, hex(x), hex(y)))
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

        kind = lookup_entity(read_byte(buf, bank, addr),
                read_byte(buf, bank, addr+1))
        addr += 2

        objects.append(("random entities", hex(count), hex(param), kind))
    elif mode == 0xf7:
        # specifically placed entities
        param, addr = read_byte(buf, bank, addr, 1)

        entities = []
        while read_byte(buf, bank, addr) < 0xf0:
            kind = lookup_entity(read_byte(buf, bank, addr),
                    read_byte(buf, bank, addr+1))
            addr += 2
            x, addr = read_byte(buf, bank, addr, 1)
            y, addr = read_byte(buf, bank, addr, 1)

            objects.append(("specific entity",
                    hex(param), kind, hex(x), hex(y)))
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
