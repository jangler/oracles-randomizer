# oracles-randomizer `asm/` notes

Asm files are written in YAML and have the following structure:

```
common:
  (bank/)?address/label: code
  ...
floating:
  label: code
  ...
seasons:
  (bank/)?address/label: code
  ...
ages:
  (bank/)?address/label: code
  ...
```

- `common` code is used in both games.
- `floating` code is defined but not given a bank or address.
- `seasons` and `ages` code only apply to the respective games.

Any of the sections can be omitted.

- A key of the form `02/openRingList` means that the label `openRingList` will
  be attached to its translated value, which will be appended to bank `02`.
- A key of the form `02/56a1/` will overwrite the data at `02:56a1` with its
  translated value. Its label is empty, so it is "anonymous" and cannot be
  referenced by other code. Non-empty labels are allowed but rarely used for
  this type of item.
- A key of the form `removeGashaNutRingText` means that the entire string is a
  label attached to its untranslated value, which is not assigned a location in
  the ROM. Another item can `/include removeGashaNutRingText` in order to use
  that value at a specific address. This type of key is only allowed in the
  `floating` section, and it is the only kind of key that can appear in that
  section.

YAML does not really care how much indentation happens as long as it happens at
all. In most cases I indent by two spaces, but I indent blocks of code by four
for readability.

The code itself is translated by [lgbtasm](https://github.com/jangler/lgbtasm).
