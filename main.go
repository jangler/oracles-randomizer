package main

//go:generate bash scripts/generate.sh
//go:generate esc -o randomizer/embed.go -pkg randomizer hints/ logic/
//go:generate esc -o rom/embed.go -pkg rom asm/ hints/ romdata/ lgbtasm/lgbtasm.lua

// the only point of this file is so that the top-level project directory
// doesn't get cluttered with a ton of .go files.

import "github.com/jangler/oracles-randomizer/randomizer"

func main() {
	randomizer.Main()
}
