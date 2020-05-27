#!/usr/bin/env lua

-- prints bytes that differ in binary files.

local function usage()
    io.stderr:write('usage: bdiff [options] file1 file2\n')
    os.exit(1)
end

if #arg < 2 then usage() end

local addr, limit = false, 10
local path1, path2
local parsefunc = nil

for _, v in ipairs(arg) do
    if parsefunc == nil then
        if v == '-a' or v == '--addr' then addr = true
        elseif v == '-l' or v == '--limit' then
            parsefunc = function(s) limit = tonumber(s) end
        elseif path1 == nil then path1 = v
        elseif path2 == nil then path2 = v
        else usage()
        end
    else
        parsefunc(v)
        parsefunc = nil
    end
end

if path1 == nil or path2 == nil then usage() end

local f1 = io.open(path1, 'rb')
local f2 = io.open(path2, 'rb')

-- get length of shortest file
local len = f1:seek('end')
f1:seek('set')
local len2 = f2:seek('end')
f2:seek('set')
if len2 < len then
    len = len2
    io.stderr:write('file sizes differ\n')
end

local function toaddr(i)
    local bank, offset = math.floor(i / 0x4000), i % 0x4000
    if i >= 0x4000 then offset = offset + 0x4000 end
    return bank, offset
end

for i = 1, len do
    local b1, b2 = string.byte(f1:read(1)), string.byte(f2:read(1))
    if b1 ~= b2 then
        if addr then
            local bank, offset = toaddr(i - 1)
            print(string.format('%02x:%04x %02x %02x', bank, offset, b1, b2))
        else
            print(string.format('%08x %02x %02x', i - 1, b1, b2))
        end

        limit = limit - 1
        if limit == 0 then
            io.stderr:write('too many diffs; specify higher --limit\n')
            os.exit(2)
        end
    end
end
