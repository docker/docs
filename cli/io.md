---
title: Input and output streams
description: How Docker CLI uses the standard streams for input and output
keywords: cli, io, stdin, stdout, stderr, streams, redirect
---

This page describes how the Docker CLI makes use of I/O streams
(`stdin`, `stdout`, `stderr`).

## Standard error

The CLI prints status messages, logs, and errors print to `stderr`.
If you want to redirect progress, logs, or error messages to a file, you must
redirect the `stderr` stream, not `stdout`.

The following example shows how to redirect `stderr` to a file (`log.txt`)
using Bash:

```console
$ docker run postgres 2> log.txt
```

For most commands that print progress messages, the Docker CLI buffers output.
The log output gets rewritten while the command executes.

## Standard output

Output from a command that executes successfully prints to `stdout`.
For example:

```console
$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED         STATUS        PORTS      NAMES
a324480dfe9c   redis     "docker-entrypoint.s…"   2 seconds ago   Up 1 second   6379/tcp   eloquent_satoshi
```

If a command results in an error, the error message prints to `stderr`, and
`stdout` is empty.

```console
$ docker ps --filter=oops
invalid argument "oops" for "-f, --filter" flag: bad format of filter (expected name=value)
See 'docker ps --help'.
```

## Standard input

Some commands can read from `stdin`, and some have flag options that let you
pass data using `stdin`. Most commands don't make use of `stdin` at all.

The `docker build` command is one of the commands that can read from `stdin`.
Running `docker build -` creates an image by reading a Dockerfile from `stdin`.
The single dash (`-`) instructs the command to read from `stdin` and use an
empty build context.

```console
$ docker build - <<EOF
FROM alpine
RUN echo "hello world" > /file
EOF
```

Another command that uses `stdin`, by default, is `docker load`. This command
lets you load a tarball into the Docker Engine's image store.

```console
$ docker load < busybox.tar.gz
```

For commands that have flags that expect files as input, you can use the single
dash (`-`) as the flag argument to make that flag read it's file from `stdin`.
In the following example, the `docker buildx bake` command reads it's bake file
from `stdin`:

```console
$ docker buildx bake --file - <<EOT
target "default" {
  args = {
    VERSION = "v1"
  }
}
EOT
```
