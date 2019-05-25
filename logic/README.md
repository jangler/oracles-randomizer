# oracles-randomizer `logic` notes

These files define (most of) the directed graph that's used for randomizer
logic. There are five types of node:

- An `and` node (`[<parent>...]`, the default) is reachable iff all of its
  parents are reachable. One with no parents is reachable.
- An `or` node (`or: [<parent>...]`) is reachable iff any of its parents is
  reachable. One with no parents is unreachable.
- A `nand` node (`not: [<parent>...]`) is a negated `and` node, and a `nor`
  node (`nor: [<parent>...]`) a negated `or` node.
- A `count` node (`count: [<min>, <parent>]`) is reachable iff its parent
  (there must be exactly one) has at least a certain number of parents which
  are also reachable.
- An `either` node (`either: []`) is always reachable, and transcends negated
  nodes. In other words, `either: []` and `not: [either []]` are both
  reachable. It is **not** an exclusive or.

Potential YAML gotchas:

- Names containing commas need to be enclosed in quotes if they appear in a
  comma-separated list.
- Associative arrays (`or`, `count`, etc) need to be enclosed in `{}` if they
  appear outside a list.

Use four spaces for indentation, and don't allow lines longer than 80
characters. Beyond that, there aren't any strict formatting rules in place.
