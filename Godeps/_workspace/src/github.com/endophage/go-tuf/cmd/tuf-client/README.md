# go-tuf client CLI

## Install

```
go get github.com/flynn/go-tuf/cmd/tuf-client
```

## Usage

The CLI provides three commands:

* `tuf-client init` - initialize a local file store using root keys (e.g. from
    the output of `tuf root-keys`)
* `tuf-client list` - list available targets and their file sizes
* `tuf-client get` - get a target file and write to STDOUT

All commands require the base URL of the TUF repository as the first non-flag
argument, and accept an optional `--store` flag which is the path to the local
storage.

Run `tuf-client help` from the command line to get more detailed usage
information.

## Examples

```
# init
$ tuf-client init https://example.com/path/to/repo

# init with a custom store path
$ tuf-client init --store /tmp/tuf.db https://example.com/path/to/repo

# list available targets
$ tuf-client list https://example.com/path/to/repo
PATH      SIZE
/foo.txt  1.6KB
/bar.txt  336B
/baz.txt  1.5KB

# get a target
$ tuf-client get https://example.com/path/to/repo /foo.txt
the contents of foo.txt

# the prefixed / is optional
$ tuf-client get https://example.com/path/to/repo foo.txt
the contents of foo.txt
```
