#!/usr/bin/env python3

import re
import sys

import yaml # pyyaml

# create an HTML checklist (doc/checklist.html) from the logic. this must be
# run from the repository's root folder.

files = {
    'seasons': 'romdata/seasons_slots.yaml',
    'ages': 'romdata/ages_slots.yaml',
}

version_regexp = re.compile('const version = "(.+?)"')
slot_regexp = re.compile('"(.+?)": +(And|Or)Slot')
name_regexp = re.compile('"(.+?)": +"(.+?)",')

doc_template = """<!DOCTYPE html>
<html>
<style>
body { font-family: sans-serif; }
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

with open('randomizer/version.go') as f:
    for line in f.readlines():
        match = version_regexp.match(line)
        if match:
            version = match.group(1)
            break

names = {}
with open('randomizer/names.go') as f:
    for line in f.readlines():
        match = name_regexp.search(line)
        if match:
            names[match.group(1)] = match.group(2)

def make_checklist(game, infile, outfile):
    with open(infile) as f:
        slots = yaml.load(f)
    elements = ['<input type="checkbox"> %s' % name for name in slots]
    outfile.write(doc_template %
        (version, game, version, game, '<br>\n'.join(elements)))


for game, infiles in files.items():
    with open('checklist/%s.html' % game, 'w') as outfile:
        make_checklist(game, infiles, outfile)
