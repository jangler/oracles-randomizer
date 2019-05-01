#!/usr/bin/env lua

-- determine the amount of free space at the end of each bank in an oracles
-- rom.

NUM_BANKS = 0x40
BANK_SIZE = 0x4000

local f = io.open(arg[1], 'rb')
local rom = f:read('*a')
f:close()

for bank = 0, NUM_BANKS - 1 do
    local run = 0

    for i = 1, BANK_SIZE do
        local byte = string.byte(rom, bank * BANK_SIZE + i)

        -- ages banks are padded with zeroes.
        -- seasons banks are padded with the bank number.
        if byte == 0 or byte == bank then
            run = run + 1
        else
            run = 0
        end
    end

    print(string.format('%02x %x', bank, run))
end
