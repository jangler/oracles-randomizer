#!/usr/bin/env python3

import re
import sys

# create an HTML checklist (scripts/checklist.html) from the prenodes. this
# must be run from the repository's root folder.


files = 'prenode/holodrum.go', 'prenode/subrosia.go', 'prenode/dungeons.go'

version_regexp = re.compile('const version = "(.+?)"')
slot_regexp = re.compile('"(.+?)": +(And|Or)Slot')

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

with open('scripts/checklist.html', 'w') as outfile:
    sections = []

    for filename in files:
        slots = []

        with open(filename) as infile:
            for line in infile.readlines():
                match = slot_regexp.search(line)
                if match:
                    slots.append(match.group(1))

        elements = ['<input type="checkbox"> %s</input>' % name
                for name in slots]

        sections.append(section_template %
                (filename[8:-3], '<br>\n'.join(elements)))

    outfile.write(doc_template % (version, version, '\n'.join(sections)))
