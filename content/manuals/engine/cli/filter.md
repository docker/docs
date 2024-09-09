---
title: Filter commands
weight: 30
description: |
  Use the filtering function in the CLI to selectively include resources
  that match the pattern you define.
keywords: cli, filter, commands, output, include, exclude
aliases:
  - /config/filter/
---

You can use the `--filter` flag to scope your commands. When filtering, the
commands only include entries that match the pattern you specify.

## Using filters

The `--filter` flag expects a key-value pair separated by an operator.

```console
$ docker COMMAND --filter "KEY=VALUE"
```

The key represents the field that you want to filter on.
The value is the pattern that the specified field must match.
The operator can be either equals (`=`) or not equals (`!=`).

For example, the command `docker images --filter reference=alpine` filters the
output of the `docker images` command to only print `alpine` images.

```console
$ docker images
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
ubuntu       20.04     33a5cc25d22c   36 minutes ago   101MB
ubuntu       18.04     152dc042452c   36 minutes ago   88.1MB
alpine       3.16      a8cbb8c69ee7   40 minutes ago   8.67MB
alpine       latest    7144f7bab3d4   40 minutes ago   11.7MB
busybox      uclibc    3e516f71d880   48 minutes ago   2.4MB
busybox      glibc     7338d0c72c65   48 minutes ago   6.09MB
$ docker images --filter reference=alpine
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
alpine       3.16      a8cbb8c69ee7   40 minutes ago   8.67MB
alpine       latest    7144f7bab3d4   40 minutes ago   11.7MB
```

The available fields (`reference` in this case) depend on the command you run.
Some filters expect an exact match. Others handle partial matches. Some filters
let you use regular expressions.

Refer to the [CLI reference description](#reference) for each command to learn
about the supported filtering capabilities for each command.

## Combining filters

You can combine multiple filters by passing multiple `--filter` flags. The
following example shows how to print all images that match `alpine:latest` or
`busybox` - a logical `OR`.

```console
$ docker images
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
ubuntu       20.04     33a5cc25d22c   2 hours ago   101MB
ubuntu       18.04     152dc042452c   2 hours ago   88.1MB
alpine       3.16      a8cbb8c69ee7   2 hours ago   8.67MB
alpine       latest    7144f7bab3d4   2 hours ago   11.7MB
busybox      uclibc    3e516f71d880   2 hours ago   2.4MB
busybox      glibc     7338d0c72c65   2 hours ago   6.09MB
$ docker images --filter reference=alpine:latest --filter=reference=busybox
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
alpine       latest    7144f7bab3d4   2 hours ago   11.7MB
busybox      uclibc    3e516f71d880   2 hours ago   2.4MB
busybox      glibc     7338d0c72c65   2 hours ago   6.09MB
```

### Multiple negated filters

Some commands support negated filters on [labels](/manuals/engine/manage-resources/labels.md).
Negated filters only consider results that don't match the specified patterns.
The following command prunes all containers that aren't labeled `foo`.

```console
$ docker container prune --filter "label!=foo"
```

There's a catch in combining multiple negated label filters. Multiple negated
filters create a single negative constraint - a logical `AND`. The following 
command prunes all containers except those labeled both `foo` and `bar`.
Containers labeled either `foo` or `bar`, but not both, will be pruned.

```console
$ docker container prune --filter "label!=foo" --filter "label!=bar"
```

## Reference

For more information about filtering commands, refer to the CLI reference
description for commands that support the `--filter` flag:

- [`docker config ls`](/reference/cli/docker/config/ls.md)
- [`docker container prune`](/reference/cli/docker/container/prune.md)
- [`docker image prune`](/reference/cli/docker/image/prune.md)
- [`docker image ls`](/reference/cli/docker/image/ls.md)
- [`docker network ls`](/reference/cli/docker/network/ls.md)
- [`docker network prune`](/reference/cli/docker/network/prune.md)
- [`docker node ls`](/reference/cli/docker/node/ls.md)
- [`docker node ps`](/reference/cli/docker/node/ps.md)
- [`docker plugin ls`](/reference/cli/docker/plugin/ls.md)
- [`docker container ls`](/reference/cli/docker/container/ls.md)
- [`docker search`](/reference/cli/docker/search.md)
- [`docker secret ls`](/reference/cli/docker/secret/ls.md)
- [`docker service ls`](/reference/cli/docker/service/ls.md)
- [`docker service ps`](/reference/cli/docker/service/ps.md)
- [`docker stack ps`](/reference/cli/docker/stack/ps.md)
- [`docker system prune`](/reference/cli/docker/system/prune.md)
- [`docker volume ls`](/reference/cli/docker/volume/ls.md)
- [`docker volume prune`](/reference/cli/docker/volume/prune.md)
