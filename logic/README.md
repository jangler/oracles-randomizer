# oracles-randomizer `logic` notes

These files define (most of) the directed graph that's used for randomizer
logic. There are three types of node:

- An `and` node (`[<parent>...]`, the default) is true iff all of its parents
  are true. One with no parents is true.
- An `or` node (`or: [<parent>...]`) is true iff any of its parents is true.
  One with no parents is false.
- A `count` node (`count: [<min>, <parent>]`) is true iff its parent (there
  must be exactly one) has at least a certain number of parents which are also
  true.

"True" generally means that the node is reachable in a hypothetical game state.

Potential YAML gotchas:

- Names containing commas need to be enclosed in quotes if they appear in a
  comma-separated list.
- Associative arrays (`or`, `count`) need to be enclosed in `{}` if they appear
  outside a list.

Use four spaces for indentation, and don't allow lines longer than 80
characters. Beyond that, there aren't any strict formatting rules in place.
