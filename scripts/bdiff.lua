#!/usr/bin/env lua

-- prints bytes that differ in binary files.
-- arg[1] = first input file
-- arg[2] = second input file
-- arg[3] = limit on number of diffs printed (default 10)

local f1 = io.open(arg[1], 'rb')
local f2 = io.open(arg[2], 'rb')

local ndiffs = 10
if arg[3] then
    ndiffs = tonumber(arg[3])
end

-- get length of shortest file
local len = f1:seek('end')
f1:seek('set')
local len2 = f2:seek('end')
f2:seek('set')
if len2 < len then
    len = len2
    io.stderr:write('file sizes differ\n')
end

for i = 1, len do
    local b1, b2 = string.byte(f1:read(1)), string.byte(f2:read(1))
    if b1 ~= b2 then
        print(string.format('%08x %02x %02x', i - 1, b1, b2))

        ndiffs = ndiffs - 1
        if ndiffs == 0 then
            io.stderr:write('too many diffs; specify limit as third CLI arg\n')
            os.exit(2)
        end
    end
end
