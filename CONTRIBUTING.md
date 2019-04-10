# Contributing to oracles-randomizer

If you plan to build the randomizer or contribute to development, here are some
things to know:


## Building

An environment with Git, Bash, and Go is required to build the randomizer.
Python 3 is required for some auxiliary scripts, like the one that generates
the HTML checklists.

`go generate` must be run before `go build` if you are starting with a fresh
repository. `go generate` should also be run after each commit (or when
switching branches) in order to keep the version string up to date.

In the future, `go generate` may be required after each code change to specific
files (and perhaps in packages other than the main package, such as
`oracles-randomizer/logic`).


## Branches

There are three main branches in the repository:

- **master**, which is for tagged release versions and documentation changes.
- **patch**, which is for bugfixes.
- **dev**, which is for new features and may not necessarily be fully
  functional.

Other branches for specific features may branch off **dev** to be merged back
into **dev** later. If you intend to make a pull request, make sure to base
your changes on the appropriate branch.


## Forking

Forking Go repositories is pretty much a mess, because Go has no notion of
relative imports. I have no specific advice on the matter, but Google might.
