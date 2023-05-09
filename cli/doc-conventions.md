---
title: CLI documentation conventions
description: Patterns and formats used in CLI documentation examples
keywords: cli, docs, conventions, examples, formats, patterns
---

Command line examples use the Unix shell format as a default. In some cases,
Windows examples are shown alongside the Unix version of the command.

## Command prompt

A dollar sign (`$`) at the beginning of a line represents the command prompt.

```console
$ docker build .
```

Examples showing commands that run as root inside a container use a pound sign
(`#`) for the prompt. The pound sign is preceded by the current working
directory.

```console
$ docker run -it alpine sh
/# echo "hello world"
hello world
/# exit
$ █
```

## Multi-line commands

A backslash `\` at the end of a line represents a multi-line command.

```console
$ docker network create \
  --driver=bridge \
  --subnet=172.28.0.0/16 \
  --ip-range=172.28.5.0/24 \
  --gateway=172.28.5.254 \
  br0
```

[Here documents](https://en.wikipedia.org/wiki/Here_document){: target="blank" rel="noopener"}
are also used to show multi-line commands.

```console
$ docker build - <<EOF
FROM alpine
RUN echo "hello world" > /foo.txt
EOF
```

In examples showing PowerShell commands, a single backtick (`\``) at the end of
a line represents a newline escape.

## Command synopsis

Examples use the following syntax to represent the interface of a command:

- Angle brackets (`<>`) surround placeholder values
- Square brackets (`[]`) surround optional arguments
- Vertical bars (`|`) separate choices
- Ellipses (`...`) can be repeated
  
Ellipses are also used for truncated command outputs.
