---
description: Using tmpfs mounts
title: Use tmpfs mounts
keywords: storage, persistence, data persistence, tmpfs
---

[Volumes](volumes.md) and [bind mounts](bind-mounts.md) are mounted into the
container's filesystem by default, and their contents are stored on the host
machine.

There may be cases where you do not want to store a container's data on the host
machine, but you also don't want to write the data into the container's writable
layer, for performance or security reasons, or if the data relates to
non-persistent application state. An example might be a temporary one-time
password that the container's application creates and uses as-needed.

To give the container access to the data without writing it anywhere
permanently, you can use a `tmpfs` mount, which is only stored in the host
machine's memory (or swap, if memory is low). When the container stops, the
`tmpfs` mount is removed. If a container is committed, the `tmpfs` mount is not
saved.

![tmpfs on the Docker host](images/types-of-mounts-tmpfs.png)

## Choosing the --tmpfs or --mount flag

Originally, the `--tmpfs` flag was used for standalone containers and
the `--mount` flag was used for swarm services. However, starting with Docker
17.06, you can also use `--mount` with standalone containers. In general,
`--mount` is more explicit and verbose. The biggest difference is that the
`--tmpfs` flag does not support any configurable options.

> **Tip**: New users should use the `--mount` syntax. Experienced users may
> be more familiar with the `--tmpfs` syntax, but are encouraged to
> use `--mount`, because research has shown it to be easier to use.

- **`--tmpfs`**: Mounts a `tmpfs` mount without allowing you to specify any
  configurable options, and can only be used with standalone  containers.

- **`--mount`**: Consists of multiple key-value pairs, separated by commas and each
  consisting of a `<key>=<value>` tuple. The `--mount` syntax is more verbose
  than `-v` or `--volume`, but the order of the keys is not significant, and
  the value of the flag is easier to understand.
  - The `type` of the mount, which can be [`bind`](bind-mounts-md), `volume`, or
    [`tmpfs`](tmpfs.md). This topic discusses `tmpfs`, so the type will always
    be `tmpfs`.
  - The `destination` takes as its value the path where the `tmpfs` mount
    will be mounted in the container. May be specified as `destination`, `dst`,
    or `target`.
  - The `tmpfs-type` and `tmpfs-mode` options. See
    [tmpfs options](#tmpfs-options).

The examples below show both the `--mount` and `--tmpfs` syntax where possible,
and `--mount` is presented first.

### Differences between `--tmpfs` and `--mount` behavior

- The `--tmpfs` flag does not allow you to specify any configurable options.
- The `--tmpfs` flag cannot be used with swarm services. You must use `--mount`.

## Limitations of tmpfs containers

- `tmpfs` mounts cannot be shared among containers.
- `tmpfs` mounts only work on Linux containers, and not on Windows containers.

## Use a tmpfs mount in a container

To use a `tmpfs` mount in a container, use the `--tmpfs` flag, or use the
`--mount` flag with `type=tmpfs` and `destination` options. There is no
`source` for `tmpfs` mounts. The following example creates a `tmpfs` mount at
`/app` in a Nginx container. The first example uses the `--mount` flag and the
second uses the `--tmpfs` flag.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-group="mount" data-target="#mount-run"><code>--mount</code></a></li>
  <li><a data-toggle="tab" data-group="volume" data-target="#tmpfs-run"><code>--tmpfs</code></a></li>
</ul>
<div class="tab-content">
<div id="mount-run" class="tab-pane fade in active" markdown="1">

```bash
$ docker run -d \
  -it \
  --name tmptest \
  --mount type=tmpfs,destination=/app \
  nginx:latest
```

</div><!--mount-->
<div id="tmpfs-run" class="tab-pane fade" markdown="1">

```bash
$ docker run -d \
  -it \
  --name tmptest \
  --tmpfs /app \
  nginx:latest
```

</div><!--volume-->
</div><!--tab-content-->

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

### Specify tmpfs options

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

## Next steps

- Learn about [volumes](volumes.md)
- Learn about [bind mounts](bind-mounts.md)
- Learn about [storage drivers](/engine/userguide/storagedriver.md)
