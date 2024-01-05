---
title: ENTRYPOINT
id: entrypoint
short_description: >
  An optional definition for the first part of the command to be run.
---

In a Dockerfile, an `ENTRYPOINT` is an optional definition for the first part
of the command to be run. If you want your Dockerfile to be runnable without
specifying additional arguments to the `docker run` command, you must specify
either `ENTRYPOINT`, `CMD`, or both.

  - If `ENTRYPOINT` is specified, it is set to a single command. Most official
    Docker images have an `ENTRYPOINT` of `/bin/sh` or `/bin/bash`. Even if you
    do not specify `ENTRYPOINT`, you may inherit it from the base image that you
    specify using the `FROM` keyword in your Dockerfile. To override the
    `ENTRYPOINT` at runtime, you can use `--entrypoint`. The following example
    overrides the entrypoint to be `/bin/ls` and sets the `CMD` to `-l /tmp`.

    ```bash
    $ docker run --entrypoint=/bin/ls ubuntu -l /tmp
    ```

  - `CMD` is appended to the `ENTRYPOINT`. The `CMD` can be any arbitrary string
    that is valid in terms of the `ENTRYPOINT`, which allows you to pass
    multiple commands or flags at once. To override the `CMD` at runtime, just
    add it after the container name or ID. In the following example, the `CMD`
    is overridden to be `/bin/ls -l /tmp`.

    ```bash
    $ docker run ubuntu /bin/ls -l /tmp
    ```
In practice, `ENTRYPOINT` is not often overridden. However, specifying the
`ENTRYPOINT` can make your images more flexible and easier to reuse.