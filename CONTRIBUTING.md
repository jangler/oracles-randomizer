# Contributing to oracles-randomizer

If you plan to build the randomizer or contribute to development, here are some
things to know:


## Building

An environment with Git, Bash, and Go is required to build the randomizer.
Python 3 is required for some auxiliary scripts, like the one that generates
the HTML checklists.

First, install Go dependencies:

```
go get -u github.com/mjibson/esc
go get -u github.com/nsf/termbox-go
```

`go generate` must be run before `go build` if you are starting with a fresh
repository or if data files (YAML etc) have been changed. `go generate` should
also be run after each commit (or when switching branches) in order to keep the
version string up to date.


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
relative imports. The plan is for oracles-randomizer to merge its "subpackages"
into its main package so that this is not a problem, but not all packages have
been merged yet.


## Code style

Always run `go fmt` on each package that has been changed (note that `go fmt`
coerces all line-initial indentation to tabs). Wrap lines longer than 80
characters when possible, assuming 8-space tabs. YAML should also be wrapped at
80 characters when possible and indented with 4 spaces.
