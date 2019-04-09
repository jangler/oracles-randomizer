#!/usr/bin/env bash

echo -n "package main

// Code generated - DO NOT EDIT.

const version = \"$(git describe --all --long | sed 's/.\+\///;s/-.\+-g/-/')\"
" > version.go
