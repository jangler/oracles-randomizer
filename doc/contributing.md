# Contributing to oracles-randomizer

If you plan to build the randomizer or contribute to development, here are some
things to know:


## Building

An environment with Git, Bash, and Go is required to build the randomizer.
Python 3 is used for some auxiliary scripts, like the one that generates the
HTML checklists, but isn't required. The following instructions are for
building the dev branch.

Clone and set up the repository:

```
git clone https://github.com/jangler/oracles-randomizer.git
git fetch
git checkout dev
git submodule init
git submodule update
```

Install Go dependencies:

```
go get github.com/mjibson/esc
go get github.com/nsf/termbox-go
go get gopkg.in/yaml.v2
```

Generate and build code (do both whenever changes are made):

```
go generate
go build
```


## Branches

There are three main branches in the repository:

- **master**, which is for tagged release versions and documentation changes.
- **patch**, which is for bugfixes.
- **dev**, which is for new features or other major changes and may not be
  fully functional.

Other branches for specific features may branch off **dev** to be merged back
into **dev** later. If you intend to make a pull request, make sure to base
your changes on the appropriate branch.


## Forking

Forking multi-package Go repositories is pretty much a mess, because Go has no
notion of relative imports. The plan is for oracles-randomizer to merge its
"subpackages" into its main package so that this is not a problem, but not all
packages have been merged yet.


## Code style

Always run `go fmt` on each package that has been changed (note that `go fmt`
coerces all line-initial indentation to tabs). Wrap lines longer than 80
characters when possible, assuming 8-space tabs. YAML should also be wrapped at
80 characters when possible and usually indented with 2 spaces.
