#!/usr/bin/env lua

-- prints addresses and verbatim rom data for warps.
-- invalid 
-- arg[1] = rom file (seasons or ages, US)
-- arg[2] = 16-bit room index
-- arg[3] = yx position or n/s/e/w

local rom = io.open(arg[1], 'rb')
local room = tonumber(arg[2], 16)
local yx = tonumber(arg[3], 16)
local edge = ({n = 0, s = 2, e = 1, w = 3})[arg[3]]

if not (yx or edge) then
    print('fatal: invalid parameter: ' .. arg[3])
    os.exit(1)
end

-- check which game the rom is
rom:seek('set', 0x13a)
local game = rom:read(3)
if game ~= 'DIN' and game ~= 'NAY' then
    print('fatal: not an oracles ROM')
    os.exit(1)
end

local bank = (0x04 - 1) * 0x4000

-- read a byte and return it as a number
local function readbyte(file)
    return string.byte(file:read(1))
end

-- follow a 16-bit pointer at the current read cursor
local function followptr(file)
    local l, h = readbyte(file), readbyte(file)
    local hl = (h << 8) + l
    file:seek('set', bank + hl)
end

-- index by high byte
local table = {DIN = 0x7457, NAY = 0x759e}
local hl = table[game] + (room >> 8) * 2
rom:seek('set', bank + hl)
followptr(rom)

-- 04:4670 in seasons, move read cursor to location of positional warp data
local function findposwarp(rom, b)
    local a = readbyte(rom)

    if a & 0x80 ~= 0 then
        -- got one
        rom:seek('cur', 1)
        return true
    elseif a & 0x40 ~= 0 then
        -- compare to room
        if readbyte(rom) == b then
            -- room found, now check for yx
            followptr(rom)
            return findposwarp(rom, yx)
        else
            rom:seek('cur', 2)
            return findposwarp(rom, b)
        end
    elseif a & 0x0f ~= 0 then
        -- next!
        rom:seek('cur', 3)
        return findposwarp(rom, b)
    else
        -- compare to yx
        if readbyte(rom) == b then
            return true
        else
            rom:seek('cur', 2)
            return findposwarp(rom, b)
        end
    end
end

-- 04:46ca in seasons, move read cursor to location of screen edge warp data
local function findedgewarp(rom, b, c)
    local a = string.byte(rom:read(1))

    if a & 0x80 ~= 0 then
        -- no warp found
        return false
    elseif a & 0x40 ~= 0 then
        -- compare to room, possibly follow pointer
        if readbyte(rom) == c then
            followptr(rom)
            return findedgewarp(rom, b, c)
        else
            rom:seek('cur', 2)
            return findedgewarp(rom, b, c)
        end
    else 
        -- compare to room, possibly find warp
        if readbyte(rom) == c and a & b ~= 0 then
            return true
        else
            rom:seek('cur', 2)
            return findedgewarp(rom, b, c)
        end
    end
end

local result = false

if yx then
    found = findposwarp(rom, room & 0xff)
else
    found = findedgewarp(rom, 1 << edge, room & 0xff)
end

if not found then
    print('no warp found')
    os.exit(1)
end

-- read the first set of warp data (group, index)
local addr1 = rom:seek() - bank
local lo, hi = readbyte(rom), readbyte(rom)

-- 04:45d0 in seasons, read (room, position, transition)
-- technically this sometimes accounts for maku tree state in seasons,
-- since the maku tree room is different for each state.
table = {DIN = 0x6d4e, NAY = 0x6f5b}
rom:seek('set', bank + table[game] + (hi >> 4) * 2)
followptr(rom)
rom:seek('cur', lo * 3)
local addr2 = rom:seek() - bank
local room, destyx, trans = readbyte(rom), readbyte(rom), readbyte(rom)

print(string.format('# addr, index, (group << 4) | (something & 0f)'))
print(string.format('%04x,%02x,%02x', addr1, lo, hi))
print(string.format('# addr, room, yx, transition'))
print(string.format('%04x,%02x,%02x,%02x', addr2, room, destyx, trans))
