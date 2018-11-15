#!/usr/bin/env python3

import re
import sys

# create an HTML checklist (doc/checklist.html) from the logic. this must be
# run from the repository's root folder.

files = {
    'seasons': ('logic/holodrum.go', 'logic/subrosia.go',
        'logic/seasons_dungeons.go'),
    'ages': ('logic/labrynna.go', 'logic/ages_dungeons.go'),
}

version_regexp = re.compile('const version = "(.+?)"')
slot_regexp = re.compile('"(.+?)": +(And|Or)Slot')
name_regexp = re.compile('"(.+?)": +"(.+?)",')

doc_template = """<!DOCTYPE html>
<html>
<style>
h1 { font-size: x-large; }
h2 { font-size: large; }
</style>
<head>
<title>oracles randomizer %s %s checklist</title>
</head>
<body>
<h1>oracles randomizer %s %s checklist</h1>
%s
</body>
</html>
"""

section_template = """<h2>%s</h2>
%s
"""

with open('summary.go') as f:
    for line in f.readlines():
        match = version_regexp.match(line)
        if match:
            version = match.group(1)
            break

names = {}
with open('names.go') as f:
    for line in f.readlines():
        match = name_regexp.search(line)
        if match:
            names[match.group(1)] = match.group(2)

def make_checklist(game, infiles, outfile):
    sections = []

    for filename in infiles:
        slots = []

        with open(filename) as infile:
            for line in infile.readlines():
                match = slot_regexp.search(line)
                if match:
                    name = match.group(1)

                    if name in names:
                        name = names[name]
                    if name[0] == 'd' and name[2] == ' ':
                            name = "D" + name[1:]
                    name = name.replace("map chest", "dungeon map chest")
                    name = name.replace("gasha chest", "gasha seed chest")

                    slots.append(name)

        elements = ['<input type="checkbox"> %s' % name
                for name in slots]

        sections.append(section_template %
                (filename[6:-3].replace(game+'_', ''), '<br>\n'.join(elements)))

    outfile.write(doc_template %
        (version, game, version, game, '\n'.join(sections)))


for game, infiles in files.items():
    with open('checklist/%s.html' % game, 'w') as outfile:
        make_checklist(game, infiles, outfile)
