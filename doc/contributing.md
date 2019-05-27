# Contributing to oracles-randomizer

If you plan to build the randomizer or contribute to development, here are some
things to know:


## Building

An environment with Git and Go is required to build the randomizer. Python 3 is
used for some auxiliary scripts, like the one that generates the HTML
checklists, but isn't required. The following instructions are for building the
dev branch.

Clone and set up the repository:

```
go get github.com/jangler/oracles-randomizer
cd $GOPATH/src/github.com/jangler/oracles-randomizer
git fetch
git checkout dev
git submodule init
git submodule update
```

Install Go dependencies (this should only be required if you `git clone` the
repository instead of using `go get` as described above):

```
go get github.com/mjibson/esc
go get github.com/nsf/termbox-go
go get github.com/yuin/gopher-lua
go get gopkg.in/yaml.v2
```

Generate and build code (do both whenever changes are made):

```
go generate
go build
```

Test by running `go test ./randomizer`.


## Branches

There are three main branches in the repository:

- **master**, which is for tagged release versions and documentation changes.
- **patch**, which is for bugfixes.
- **dev**, which is for new features or other major changes and may not be
  fully functional.

Other branches for specific features may branch off **dev** to be merged back
into **dev** later. If you intend to make a pull request, make sure to base
your changes on the appropriate branch.


## Code style

Always run `go fmt` in the `randomizer` package before commits (note that `go
fmt` coerces all line-initial indentation to tabs). Wrap lines longer than 80
characters when possible, assuming 8-space tabs. YAML should also be wrapped at
80 characters when possible and indented with 2 or 4 spaces, depending on the
conventions of the directory it's in.
