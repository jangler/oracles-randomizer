#!/usr/bin/env bash

echo -n "package randomizer

// Code generated - DO NOT EDIT.

const version = \"$(git describe --all --long | sed 's/.\+\///;s/-.\+-g/-/')\"
" > randomizer/version.go
