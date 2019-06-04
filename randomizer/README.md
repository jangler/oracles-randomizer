# oracles-randomizer coding tips

- Names in Go are conventially lowerMixedCaps by default, but UpperMixedCaps is
  semantically different: it "exports" the name so that external packages can
  access it. This is necessary for the `Main()` function, which is called from
  `main.go` in the repository's root directory, and structs members which are
  being unmarshaled from YAML by the `yaml` package. If any names beyond these
  are UpperMixedCaps, it's probably vestigial from when the randomizer was
  split into multiple Go packages (and should be corrected).
- It's OK to `panic()` in unrecoverable situations where some part of the
  randomizer itself is incorrect, like if embedded YAML can't be loaded. For
  extrinsic situations like invalid user input (including input ROM), functions
  should return `error`s instead, which the `Main()` function should ultimately
  report and then exit gracefully.
- Go lacks a ternary operator, so the randomizer defines a `ternary()` function
  to help make some code more concise. But this is *not* the same as a ternary
  operator, since arguments to functions are evaluated before passing them. So
  `ternary(a < b, a, b)` will work as intended, but
  `ternary(len(a) < 2, a[1], a[2])` will panic if the condition isn't met.

If you're unfamiliar with Go as a language,
<https://golang.org/doc/effective_go.html> is probably the best reference for
most aspects of it. Other documents are available at <https://golang.org/doc/>.
