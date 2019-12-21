#!/usr/bin/env bash

# create zip files for distribution on the three major OSes
#
# strictly speaking nothing here should need quoting but that's just my
# environment

go generate
python3 scripts/checklist.py
unix2dos -n README.md README.txt

version="$(grep -o '".\+"' randomizer/version.go | tr -d '"')"
appname="$(basename "$PWD")"

mkdir -p "dist/$version"

function buildfor() {
	echo "building for $1/$2"
	GOOS=$1 GOARCH=$2 go build
	zip -r "dist/$version/$appname"_$3_"$version.zip" "$appname$4" \
		README.txt tracker/
}

buildfor windows 386 win32 .exe
buildfor darwin amd64 macos64
buildfor linux amd64 linux64

rm README.txt

echo "archives written to dist/$version/"
