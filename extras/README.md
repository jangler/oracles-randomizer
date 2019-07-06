# oracles-randomizer `extras/` notes

**This is an experimental feature, and is not "supported" in the same sense
that the rest of the randomizer is. Use at your own risk!**

These files contain additional self-contained assembly code that can be
included using the `-include` command-line option when randomizing a ROM. The
file format is exactly the same as for the files in the `asm/` folder. Included
files are applied as the last step in randomization, so they will overwrite any
of the usual changes that the randomizer has made. Labels without a specified
address (but with a specified bank) are the exception and will be placed at the
end of whatever code the randomizer has already added to the bank. Overflowing
a bank causes a `panic` from Go.

The main intention of this feature is for augmenting plandos (see
`doc/plando.md`) with additional changes that aren't natively supported by the
randomizer.
