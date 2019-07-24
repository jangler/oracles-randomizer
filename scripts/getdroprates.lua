#!/usr/bin/env lua

-- prints bytes that differ in binary files, with rom bank addresses.
-- arg[1] = rom file (seasons or ages, US)
-- arg[2] = drop table index
-- arg[3] = high byte of room
-- arg[4] = if arg[2] is zero, enemy value at $dxc2

local bank_offset = (0x3f - 1) * 0x4000
local seek_table = {DIN = 0x4a75, NAY = 0x4a46}
local prob_table = {DIN = 0x481d, NAY = 0x47fe}
local drop_table = {DIN = 0x47dd, NAY = 0x47be}
local grp1_table = {DIN = 0x4865, NAY = nil} -- idk when this applies for ages
local name_table = {
    'fairy',         -- 0
    'heart',         -- 1
    '1 rupee',       -- 2
    '5 rupees',      -- 3
    'bombs',         -- 4
    'ember seeds',   -- 5
    'scent seeds',   -- 6
    'pegasus seeds', -- 7
    'gale seeds',    -- 8
    'mystery seeds', -- 9
    nil,             -- a
    nil,             -- b
    '1 ore',         -- c
    '5 ores',        -- d
    nil,             -- e
    'enemy',         -- f
}

local rom = io.open(arg[1], 'rb')
local c = tonumber(arg[2]) -- register c passed to drop function
local group = tonumber(arg[3])
local enemy_var = tonumber(arg[4])

rom:seek('set', 0x13a)
local game = rom:read(3)
if game ~= "DIN" and game ~= "NAY" then
    print("fatal: not an oracles ROM")
    os.exit(1)
end

-- get offset into probability table
print(string.format('index: 0x%x', c))
local a = c | 0x80
if c == 0 then
    a = enemy_var
end
local hl = seek_table[game] + a
rom:seek('set', bank_offset + hl)
a = string.byte(rom:read(1))
c = a
a = ((a & 0xf) << 4) | ((a & 0xf0) >> 4) -- swap a
a = (a >> 1) | ((a & 1) << 7) -- rrca
a = a & 7

-- get probability/64 that a drop will occur
hl = prob_table[game] + a * 2
rom:seek('set', bank_offset + hl)
local l = string.byte(rom:read(1))
local h = string.byte(rom:read(1))
hl = (h << 8) + l
rom:seek('set', bank_offset + hl)
local prob_block = rom:read(8)
local num_set_bits = 0
for i = 1, #prob_block do
    local byte = string.byte(prob_block, i, i)
    for j = 0, 7 do
        if byte & (1 << j) ~= 0 then
            num_set_bits = num_set_bits + 1
        end
    end
end
print(table.pack(string.gsub(
    string.format('droprate: %f', num_set_bits / 0x40), '0+$', ''))[1])

-- get probability of specific drop indexes
a = c & 0x1f
hl = drop_table[game] + a * 2
rom:seek('set', bank_offset + hl)
l = string.byte(rom:read(1))
h = string.byte(rom:read(1))
hl = (h << 8) + l
rom:seek('set', bank_offset + hl)
local drop_block = rom:read(0x20)
local drop_probs = {}
for i = 1, #drop_block do
    local byte = string.byte(drop_block, i, i)
    drop_probs[byte] = drop_probs[byte] or 0
    drop_probs[byte] = drop_probs[byte] + 1
end

-- group 1 (subrosia) converts index
if group == 1 then
    hl = grp1_table[game]
    rom:seek('set', bank_offset + hl)
    local grp1_block = rom:read(0x20)
    local new_probs = {}
    for i = 1, #grp1_block do
        if drop_probs[i-1] ~= nil then
            local byte = string.byte(grp1_block, i, i)
            new_probs[byte] = new_probs[byte] or 0
            new_probs[byte] = new_probs[byte] + drop_probs[i-1]
        end
    end
    drop_probs = new_probs
end

-- convert probabilities to sorted unit fractions
local named_probs = {}
for i = 1, 0x20 do
    if drop_probs[i-1] ~= nil then
        named_probs[i] = {
            name = name_table[i] or i-1,
            prob = drop_probs[i-1] / 0x20
        }
    else
        named_probs[i] = {name = '', prob = 0}
    end
end
table.sort(named_probs, function (a, b) return a.prob > b.prob end)

-- print drop probabilities
print('drops:')
for i, drop in ipairs(named_probs) do
    if drop.prob ~= 0 then
        local prob_string = table.pack(string.gsub(
            string.format('%f', drop.prob), '%.?0+$', ''))[1]
        if prob_string == '1' then
            prob_string = '1.0'
        end
        print(string.format('  %s: %s', drop.name, prob_string))
    end
end
