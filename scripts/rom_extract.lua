#!/usr/bin/env lua

-- writes a section of a rom to a new file.
-- arg[1] = input file (probably a rom)
-- arg[2] = output file
-- arg[3] = hex bank number
-- arg[4] = hex offset within bank + 0x4000
-- arg[5] = hex data size

if #arg ~= 5 then
    print(string.format('usage: %s infile outfile bank offset size', arg[0]))
    os.exit(1)
end

local adjusted_bank = math.max(0, tonumber(arg[3], 16) - 1)
local total_offset = (0x4000 * adjusted_bank) + tonumber(arg[4], 16)
local size = tonumber(arg[5], 16)

local infile = io.open(arg[1], 'rb')
infile:seek('set', total_offset)
local data = infile:read(size)
infile:close()

local outfile = io.open(arg[2], 'wb')
outfile:write(data)
outfile:close()
