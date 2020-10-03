# Building and contributing to oracles-randomizer

If you plan to build the randomizer from source or contribute to development,
here are some things to know:


## Building

An environment with Git and Go is required to build the randomizer. The
following instructions are for building the dev branch.

First, clone and set up the repository:

```
go get github.com/jangler/oracles-randomizer
cd $GOPATH/src/github.com/jangler/oracles-randomizer
git fetch
git checkout dev
git submodule init
git submodule update
```

You'll get warnings about being unable to build the code, which is expected and
since some of the code isn't tracked by the repository and hasn't been
generated locally yet. Then install Go dependencies:

```
go get github.com/mjibson/esc
go get github.com/gdamore/tcell
go get github.com/yuin/gopher-lua
go get gopkg.in/yaml.v2
```

Last, generate and build the code (do both whenever changes are made):

```
go generate
go build
```

You'll probably have to add `$GOPATH/bin` to your "path" environment variable
in order for `esc` to work for code generation. Alternately, you can copy
`$GOPATH/bin/esc` to somewhere that's already in your "path".


## Branches

There are three main branches in the repository:

- **master**, which is for tagged release versions and documentation changes.
- **patch**, which is for bugfixes.
- **dev**, which is for new features or other major changes and may not be
  fully functional.

Other branches for specific features may branch off **dev** to be merged back
into **dev** later. If you intend to make a pull request, make sure to base
your changes on the appropriate branch. Also ask in advance unless you're
making a simple bugfix. *Also* also run `go test ./randomizer/` to make sure
tests pass before making commits.


## Organization

Go code is in `randomizer/`, but some types of changes don't even need to touch
the Go code. Logic is in `logic/`, GBC assembly code is in `asm/`, owl hint
names are in `hints/`, and various ROM addresses and values are in `romdata/`.
All the non-Go directories use YAML, although sometimes the contents of the
YAML amount to something more like a domain-specific language.


## Code style

Always run `go fmt ./randomizer/` before commits (note that `go fmt` coerces
all line-initial indentation to tabs). Wrap lines (in Go and YAML) at 80
characters when possible, assuming 8-space tabs for Go. YAML is necessarily
indented using spaces; the standard varies between 2 and 4 spaces, depending on
the conventions of the directory.
