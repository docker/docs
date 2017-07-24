---
description: Using tmpfs mounts
title: Use tmpfs mounts
keywords: storage, persistence, data persistence, tmpfs
---

[Volumes](volumes.md) and [bind mounts](bind-mounts.md) are mounted into the
container's filesystem by default, and their contents are stored on on the host
machine.

There may be cases where you do not want to store a container's data on the host
machine, but you also don't want to write the data into the container's writable
layer, for performance or security reasons, or if the data relates to
non-persistent application state.

To give the container access to the data without writing it anywhere
permanently, you can use a `tmpfs` mount, which is only stored in the host
machine's memory (or swap, if memory is low). When the container stops, the
`tmpfs` mount is removed. If a container is committed, the `tmpfs` mount is not
saved.

## Limitations of tmpfs containers

- `tmpfs` mounts cannot be shared among containers.
- `tmpfs` mounts only work on Linux containers, and not on Windows containers.

## Use a tmpfs mount in a container

To use a `tmpfs` mount in a container, use the `--tmpfs` flag, or use the
`--mount` flag with `type=tmpfs` and `destination` options. There is no
`source` for `tmpfs` mounts. The following example creates a `tmpfs` mount at
`/app` in a Nginx container. The first example uses the `--tmpfs` flag and the
second uses the `--mount` flag.

```bash
$ docker run -d \
  -it \
  --name tmptest \
  --tmpfs /app \
  nginx:latest
```

```bash
$ docker run -d \
  -it \
  --name tmptest \
  --mount type=tmpfs,destination=/app \
  nginx:latest
```

Verify that the mount is a `tmpfs` mount by running `docker container inspect
tmptest` and looking for the `Mounts` section:

```json
"Tmpfs": {
    "/app": ""
},
```

Remove the container:

```bash
$ docker container stop tmptest

$ Docker container rm tmptest
```

### tmpfs options

`tmpfs` mounts allow for two configuration options, neither of which is
required. If you need to specify these options, you must use the `--mount` flag,
as the `--tmpfs` flag does not support them.

| Option       | Description                                                                                           |
|:-------------|:------------------------------------------------------------------------------------------------------|
| `tmpfs-size` | Size of the tmpfs mount in bytes. Unlimited by default.                                               |
| `tmpfs-mode` | File mode of the tmpfs in octal. For instance, `700` or `0770`. Defaults to `1777` or world-writable. |

The following example sets the `tmpfs-mode` to `1770`, so that it is not
world-readable within the container.

```bash
docker run -d \
  -it \
  --name tmptest \
  --mount type=tmpfs,destination=/app,tmpfs-mode=1770 \
  nginx:latest
```

### Syntax differences for services

The `docker service create` command does not support thd `--tmpfs` flag. When
mounting a `tmpfs` mount into a service's containers, you must use the `--mount`
flag.

## Next steps

- Learn about [volumes](volumes.md)
- Learn about [bind mounts](bind-mounts.md)
- Learn about [storage drivers](/engine/userguide/storagedriver.md)