#!/usr/bin/env python3

import re
import sys

# create an HTML checklist (scripts/checklist.html) from the logic. this must
# be run from the repository's root folder.


files = 'logic/holodrum.go', 'logic/subrosia.go', 'logic/dungeons.go'

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
<title>oos-randomizer %s checklist</title>
</head>
<body>
<h1>oos-randomizer %s checklist</h1>
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

with open('scripts/checklist.html', 'w') as outfile:
    sections = []

    for filename in files:
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
                (filename[6:-3], '<br>\n'.join(elements)))

    outfile.write(doc_template % (version, version, '\n'.join(sections)))
