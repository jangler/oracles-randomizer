#!/usr/bin/env bash

# create zip files for distribution on the three major OSes
#
# strictly speaking nothing here should need quoting but that's just my
# environment

version="$(git tag --contains HEAD)"
appname="$(basename "$PWD")"

unix2dos -n README.md README.txt
unix2dos -n doc/checklist.html checklist.html

mkdir -p "dist/$version"
GOOS=windows GOARCH=386 go build
apack "dist/$version/$appname"_win32_"$version.zip" "$appname.exe" README.txt \
	checklist.html
GOOS=darwin GOARCH=amd64 go build
apack "dist/$version/$appname"_macos64_"$version.zip" "$appname" README.txt \
	checklist.html
GOOS=linux GOARCH=amd64 go build
apack "dist/$version/$appname"_linux64_"$version.zip" "$appname" README.txt \
	checklist.html

rm README.txt checklist.html

echo "============================="
echo "MAKE SURE TO UPDATE VERSION!!"
echo "============================="
