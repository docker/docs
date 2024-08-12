---
title: Docker contexts
description: Learn about managing multiple daemons from a single client with contexts
keywords: engine, context, cli, daemons, remote
aliases:
  - /engine/context/working-with-contexts/
---

## Introduction

This guide shows how you can use contexts to manage Docker daemons from a single client.

Each context contains all information required to manage resources on the daemon.
The `docker context` command makes it easy to configure these contexts and switch between them.

As an example, a single Docker client might be configured with two contexts:

- A default context running locally
- A remote, shared context

Once these contexts are configured,
you can use the `docker context use <context-name>` command
to switch between them.

## Prerequisites

To follow the examples in this guide, you'll need:

- A Docker client that supports the top-level `context` command

Run `docker context` to verify that your Docker client supports contexts.

## The anatomy of a context

A context is a combination of several properties. These include:

- Name and description
- Endpoint configuration
- TLS info

To list available contexts, use the `docker context ls` command.

```console
$ docker context ls
NAME        DESCRIPTION                               DOCKER ENDPOINT               ERROR
default *                                             unix:///var/run/docker.sock
```

This shows a single context called "default".
It's configured to talk to a daemon through the local `/var/run/docker.sock` Unix socket.

The asterisk in the `NAME` column indicates that this is the active context.
This means all `docker` commands run against this context,
unless overridden with environment variables such as `DOCKER_HOST` and `DOCKER_CONTEXT`,
or on the command-line with the `--context` and `--host` flags.

Dig a bit deeper with `docker context inspect`.
The following example shows how to inspect the context called `default`.

```console
$ docker context inspect default
[
    {
        "Name": "default",
        "Metadata": {},
        "Endpoints": {
            "docker": {
                "Host": "unix:///var/run/docker.sock",
                "SkipTLSVerify": false
            }
        },
        "TLSMaterial": {},
        "Storage": {
            "MetadataPath": "\u003cIN MEMORY\u003e",
            "TLSPath": "\u003cIN MEMORY\u003e"
        }
    }
]
```

### Create a new context

You can create new contexts with the `docker context create` command.

The following example creates a new context called `docker-test` and specifies
the host endpoint of the context to TCP socket `tcp://docker:2375`.

```console
$ docker context create docker-test --docker host=tcp://docker:2375
docker-test
Successfully created context "docker-test"
```

The new context is stored in a `meta.json` file below `~/.docker/contexts/`.
Each new context you create gets its own `meta.json` stored in a dedicated sub-directory of `~/.docker/contexts/`.

You can view the new context with `docker context ls` and `docker context inspect <context-name>`.

```console
$ docker context ls
NAME          DESCRIPTION                             DOCKER ENDPOINT               ERROR
default *                                             unix:///var/run/docker.sock
docker-test                                           tcp://docker:2375
```

The current context is indicated with an asterisk ("\*").

## Use a different context

You can use `docker context use` to switch between contexts.

The following command will switch the `docker` CLI to use the `docker-test` context.

```console
$ docker context use docker-test
docker-test
Current context is now "docker-test"
```

Verify the operation by listing all contexts and ensuring the asterisk ("\*") is against the `docker-test` context.

```console
$ docker context ls
NAME            DESCRIPTION                           DOCKER ENDPOINT               ERROR
default                                               unix:///var/run/docker.sock
docker-test *                                         tcp://docker:2375
```

`docker` commands will now target endpoints defined in the `docker-test` context.

You can also set the current context using the `DOCKER_CONTEXT` environment variable.
The environment variable overrides the context set with `docker context use`.

Use the appropriate command below to set the context to `docker-test` using an environment variable.

{{< tabs >}}
{{< tab name="PowerShell" >}}

```ps
> $env:DOCKER_CONTEXT='docker-test'
```

{{< /tab >}}
{{< tab name="Bash" >}}

```console
$ export DOCKER_CONTEXT=docker-test
```

{{< /tab >}}
{{< /tabs >}}

Run `docker context ls` to verify that the `docker-test` context is now the
active context.

You can also use the global `--context` flag to override the context.
The following command uses a context called `production`.

```console
$ docker --context production container ls
```

## Exporting and importing Docker contexts

You can use the `docker context export` and `docker context import` commands
to export and import contexts on different hosts.

The `docker context export` command exports an existing context to a file.
The file can be imported on any host that has the `docker` client installed.

### Exporting and importing a context

The following example exports an existing context called `docker-test`.
It will be written to a file called `docker-test.dockercontext`.

```console
$ docker context export docker-test
Written file "docker-test.dockercontext"
```

Check the contents of the export file.

```console
$ cat docker-test.dockercontext
```

Import this file on another host using `docker context import`
to create context with the same configuration.

```console
$ docker context import docker-test docker-test.dockercontext
docker-test
Successfully imported context "docker-test"
```

You can verify that the context was imported with `docker context ls`.

The format of the import command is `docker context import <context-name> <context-file>`.

## Updating a context

You can use `docker context update` to update fields in an existing context.

The following example updates the description field in the existing `docker-test` context.

```console
$ docker context update docker-test --description "Test context"
docker-test
Successfully updated context "docker-test"
```
