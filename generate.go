package main

// the only point of this package is so that the top-level project directory
// doesn't get cluttered with a ton of .go files. for convenience (mainly when
// forking), `go generate` automatically generates main.go (a file which the
// git repo is configured to ignore) importing the appropriate local path.

//go:generate go run generate/generate.go
//go:generate esc -o randomizer/embed.go -pkg randomizer asm/ hints/ logic/ romdata/ lgbtasm/lgbtasm.lua
