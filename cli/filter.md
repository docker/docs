---
title: Filter command outputs
description: |
  Use the filtering function in the CLI to selectively include output
  that matches the pattern you define.
keywords: cli, filter, output, include, exclude, regex
---

You can use the `--filter`, short form `-f`, flag to prune the output of commands.
When you apply a filter, command output only include the results that match the
pattern you specify.

Commands that support the `--filter` flag are:

- [`docker config ls`](../engine/reference/commandline/config_ls.md)
- [`docker container ls`](../engine/reference/commandline/container_ls.md)
- [`docker container prune`](../engine/reference/commandline/container_prune.md)
- [`docker image ls`](../engine/reference/commandline/image_ls.md)
- [`docker image prune`](../engine/reference/commandline/image_prune.md)
- [`docker images`](../engine/reference/commandline/images.md)
- [`docker network ls`](../engine/reference/commandline/network_ls.md)
- [`docker network prune`](../engine/reference/commandline/network_prune.md)
- [`docker node ls`](../engine/reference/commandline/node_ls.md)
- [`docker node ps`](../engine/reference/commandline/node_ps.md)
- [`docker plugin ls`](../engine/reference/commandline/plugin_ls.md)
- [`docker ps`](../engine/reference/commandline/ps.md)
- [`docker search`](../engine/reference/commandline/search.md)
- [`docker secret ls`](../engine/reference/commandline/secret_ls.md)
- [`docker service ls`](../engine/reference/commandline/service_ls.md)
- [`docker service ps`](../engine/reference/commandline/service_ps.md)
- [`docker stack ps`](../engine/reference/commandline/stack_ps.md)
- [`docker system prune`](../engine/reference/commandline/system_prune.md)
- [`docker volume ls`](../engine/reference/commandline/volume_ls.md)
- [`docker volume prune`](../engine/reference/commandline/volume_prune.md)

The field that you can use with `--filter` depends on the command you invoke.
Refer to the CLI reference for commands to learn which fields they support.

## How to use filters

The `--filter` flag expects a key-value string value:

```text
--filter KEY=VALUE
```

The `KEY` in the value represents the field that you want to filter on.
The `VALUE` is the pattern that the specified field must match.

For example, the `docker images` command supports filtering on the image name
(`reference`) as a field.

```console
$ docker images
REPOSITORY      TAG       IMAGE ID       CREATED            SIZE
alpine          latest    124c7d270790   3 days ago         11.4MB
mysql           latest    a43f6e7e7f3a   2 days ago         804MB
redis           latest    ea30bef6a142   5 days ago         164MB
ubuntu          latest    dfd64a3b4296   46 minutes ago     106MB
$ docker images --filter reference=alpine
REPOSITORY      TAG       IMAGE ID       CREATED            SIZE
alpine          latest    124c7d270790   3 days ago         11.4MB
```

The `reference` filter for the `docker images` command expects an exact match.

```console
$ docker images --filter reference=alp
REPOSITORY      TAG       IMAGE ID       CREATED            SIZE
```

Most filters expect and exact match, and others handle partial matches, and
some even let you use regular expressions.

Refer to the CLI reference description for each command to learn about the
supported filtering capabilities for each command.
