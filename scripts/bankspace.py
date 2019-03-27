#!/usr/bin/env python3

# determine the amount of free space at the end of each bank in an oracles rom.

import sys


with open(sys.argv[1], 'rb') as f:
    rom = f.read()

for i in range(0x40):
    run = 0
    for j in range(0x4000):
        if rom[i*0x4000+j] in (0, i):
            run += 1
        else:
            run = 0
    print(hex(i), hex(run))
