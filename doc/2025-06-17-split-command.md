Implement a custom command that improves upon the existing `split` utility.

The key improvement of this command is to support file-extension-based splitting. For
example, if I split the file `foo-bar.ext` into 2 parts, it should result in the output
files `foo-bar-0.ext` and `foo-bar-0.ext`.

The command should be invoked as a subcommand of the main `mdcli` command, ex. `mdcli
split ...`.

To start, it should only support numeric suffixes (ie. `split -d ...`).

It should support the following split options, similar to `split`:

* `-l`: split into chunks of specified number of lines
* `-c`: split into specified number of chunks
* `-b`: split into chunks of specified size.

Additionally, the `-c` and `-b` options should support an additional `l/N` flag. If the
command is invoked as `mdcli split -c l/5 ...`, the command should split the input into 5
roughly-equal sized chunks, however the split should only occur between line breaks.
Similarly, if the command is invoked as `mdcli split -b l/8Mi`, the command should split
the input into chunks of roughly 8Mi size, with split only occurring on line breaks.

For the `-b` flag, size can be specified as an integer number of bytes, or as a suffixed
quantity. The `K, `Ki`, `M`, `Mi`, `G`, and `Gi` are supported, indicating kilo-, kibi-,
mega-, mebi-, giga-, and gibi- bytes, respectively.

In addition to the flags, the command should take up to 2 args as input, the input file
and the output prefix. If the input file does not have a file extension, the command
should work equivalently to the `split` command. If the input file does have an extension,
the split suffix should be added to the file base-name and the extension should be added
after the suffix.

If the output prefix is not specified (ie. only 1 arg is provided), the output prefix
should be assumed to be the same as the input file without the extension, ie. `mdcli split
-n 10 foo/bar/baz.qux` should split into the `baz-0.qux`, `baz-1.qux`, etc. files in the
`foo/bar` directory.
