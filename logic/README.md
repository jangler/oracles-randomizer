# oracles-randomizer `logic` notes

These files define (most of) the directed graph that's used for randomizer
logic. There are four types of node:

- An `and` node (`[<parent>...]`, the default) is reachable iff all of its
  parents are reachable. One with no parents is reachable.
- An `or` node (`or: [<parent>...]`) is reachable iff any of its parents is
  reachable. One with no parents is unreachable.
- A `count` node (`count: [<min>, <parent>]`) is reachable iff its parent
  (there should be exactly one) has at least a certain number of parents which
  are also reachable.
- The singular `rupees` node relays the net rupee value of its parents to its
  children, which are `count` nodes.

Potential YAML gotchas:

- Names containing commas need to be enclosed in quotes if they appear in a
  comma-separated list.
- Associative arrays (`or`, `count`, etc) need to be enclosed in `{}` if they
  appear outside a list.

Use four spaces for indentation, and don't allow lines longer than 80
characters. Beyond that, there aren't any strict formatting rules in place.
