---
title: CLI flags
description: Using CLI flags with Docker commands.
keywords: cli, flags, options, option stacking, values, types
---

The Docker CLI uses flags to configure the behavior of commands. This page
describes how to set flags, and how flag values work.

## Option stacking

Single character command line options can be combined.
Rather than typing `docker run -i -t --name test busybox sh`,
you can write `docker run -it --name test busybox sh`.

## Repeating flags

You can specify options like `-a=[]` multiple times in a single command line,
for example in these commands:

```console
$ docker run -a stdin -a stdout -i -t ubuntu /bin/bash

$ docker run -a stdin -a stdout -a stderr ubuntu /bin/ls
```

Sometimes, multiple options can call for a more complex value string as for
`-v`:

```console
$ docker run -v /host:/container example/mysql
```

> **Note**
>
> Do not use the `-t` and `-a stderr` options together due to
> limitations in the `pty` implementation. All `stderr` in `pty` mode
> simply goes to `stdout`.

## Flag values

Flags can accept map, list, string, integer, or Boolean values.

### Boolean

Boolean options take the form `-d=false`. The value you see in the help text is
the default value which is set if you do **not** specify that flag. If you
specify a Boolean flag without a value, this will set the flag to `true`,
irrespective of the default value.

For example, running `docker run -d` will set the value to `true`, so your
container **will** run in "detached" mode, in the background.

Options which default to `true` (e.g., `docker build --rm=true`) can only be
set to the non-default value by explicitly setting them to `false`:

```console
$ docker build --rm=false .
```

### Strings and integers

Options like `--name=""` expect a string, and they can only be specified once.
Options like `-c=0` expect an integer, and they can only be specified once.

### Maps

CLI flags that take a map value expect a CSV-formatted string for a value.

For example, the `docker run --mount` flag expects a map value with the
following keys:

- `type`
- `source` (or `src`)
- `target` (or `dst`)

You define the map on the CLI as follows:

```console
$ docker run --mount type=bind,source=/home/me/code,target=/app <image>
```

## Flag inheritance

The base `docker` command sets a number of flags that child commands, such as
`run` and `build`, inherit. Refer to
[`docker` base command](../engine/reference/commandline/index.md).
