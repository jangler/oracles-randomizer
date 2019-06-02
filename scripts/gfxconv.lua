#!/usr/bin/env lua

-- converts an imagemagicked raw .gray file to asm-formatted 2bpp tile data.
-- arg[1] = input file
-- arg[2] = output file

local infile = io.open(arg[1], 'rb')
local outfile = io.open(arg[2], 'w')

local len = infile:seek('end')
infile:seek('set')

for tile = 1, len, 64 do -- 64 bytes = one 8x8 tile
    local tilestr = 'db '

    for line = 1, 8 do -- 8 bytes = one 8-pixel row
        local grayline = infile:read(8)
        local high, low = 0, 0 -- 2 bytes = one 8-pixel row

        for j = 1, #grayline do
            local c = string.byte(grayline, j)
            low = low | (1 << (8 - j))
            if c == 0 then
                high = high | (1 << (8 - j))
            end
        end

        if line ~= 1 then
            tilestr = tilestr .. ','
        end
        tilestr = tilestr .. string.format('%02x,%02x', high, low)
    end

    outfile:write(tilestr .. '\n')
end
