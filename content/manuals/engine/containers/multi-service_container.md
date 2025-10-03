---
description: Learn how to run more than one process in a single container
keywords: docker, supervisor, process management
title: Run multiple processes in a container
weight: 20
aliases:
  - /articles/using_supervisord/
  - /engine/admin/multi-service_container/
  - /engine/admin/using_supervisord/
  - /engine/articles/using_supervisord/
  - /config/containers/multi-service_container/
---

A container's main running process is the [`ENTRYPOINT`](https://docs.docker.com/reference/dockerfile/#entrypoint)
and/or [`CMD`](https://docs.docker.com/reference/dockerfile/#cmd) at the
end of the `Dockerfile`. It's best practice to separate areas of concern by
using one service per container. That service may fork into multiple
processes (for example, Apache web server starts multiple worker processes).
It's ok to have multiple processes, but to get the most benefit out of Docker,
avoid one container being responsible for multiple aspects of your overall
application. You can connect multiple containers using user-defined networks and
shared volumes.

The container's main process is responsible for managing all processes that it
starts. In some cases, the main process isn't well-designed, and doesn't handle
"reaping" (stopping) child processes gracefully when the container exits or signal forwarding.
If your process falls into this category, you can use the
[`--init` option](https://docs.docker.com/reference/cli/docker/container/run/#init) when you
run the container. The `--init` flag inserts a tiny init-process into the
container as the main process, and handles reaping of all processes when the
container exits. Handling such processes this way is superior to using a
full-fledged init process such as `sysvinit` or `systemd` to handle process
lifecycle within your container.

If you need to run more than one service within a container, you can achieve
this in a few different ways. Bear in mind that this always comes with a trade-off.

## Use a wrapper script

Put all of your commands in a wrapper script, complete with testing and
debugging information. Run the wrapper script as your `CMD`. The following is a
naive example. First, the wrapper script:

```bash
#!/bin/bash

# Start the first process
./my_first_process &

# Start the second process
./my_second_process &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?
```

Next, the Dockerfile:

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu:latest
COPY my_first_process my_first_process
COPY my_second_process my_second_process
COPY my_wrapper_script.sh my_wrapper_script.sh
CMD ./my_wrapper_script.sh
```

If you combine this approach with the `--init` flag mentioned above, you will get
the benefits of reaping zombie processes but no signal forwarding. Subprocesses may
not terminate gracefully. If your application requires a gracefull shutdown, be aware
of this pitfall.

## Use Bash job controls

If you have one main process that needs to start first and stay running but you
temporarily need to run some other processes (perhaps to interact with the main
process) then you can use bash's job control. First, the wrapper script:

```bash
#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
./my_main_process &

# Start the helper process
./my_helper_process

# the my_helper_process might need to know how to wait on the
# primary process to start before it does its work and returns


# now we bring the primary process back into the foreground
# and leave it there
fg %1
```

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu:latest
COPY my_main_process my_main_process
COPY my_helper_process my_helper_process
COPY my_wrapper_script.sh my_wrapper_script.sh
CMD ./my_wrapper_script.sh
```

## Use a process manager

This is more involved than the other options, as it requires you to
bundle the process manager binary and its configuration into your image
(or base your image on one that includes it), along with
the different applications it manages. Then you start the process manager as PID 1, which
manages your processes for you.

As with process manager, you do not need the `--init` parameter.

### supervisord

The following Dockerfile example shows this approach for [`supervisord`](https://supervisord.org/).
The example assumes that these files exist at the root of the build context:

- `supervisord.conf`
- `my_first_process`
- `my_second_process`

```dockerfile
# syntax=docker/dockerfile:1
FROM ubuntu:latest
RUN apt-get update && apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY my_first_process my_first_process
COPY my_second_process my_second_process
CMD ["/usr/bin/supervisord"]
```

If you want to make sure both processes output their `stdout` and `stderr` to
the container logs, you can add the following to the `supervisord.conf` file:

```ini
[supervisord]
nodaemon=true
logfile=/dev/null
logfile_maxbytes=0

[program:app]
stdout_logfile=/dev/fd/1
stdout_logfile_maxbytes=0
redirect_stderr=true
```

The obvious downside to this approach is that your container will run indefinetly
as `supervisord` is not designed to terminate when a supervised process terminates.
If you aim for a small image size of your container, this might be an issue to, as it
requires to have a full python runtime available.

### s6 and s6-overlay

[s6-overlay](https://github.com/just-containers/s6-overlay) is a layer that can be installed
on top of your container and uses the [s6](https://skarnet.org/software/s6/overview.html)
process manager.

To make you of it in your container, you have to pull it into your dockerfile first:

```dockerfile
# Use your favorite image
FROM ubuntu
ARG S6_OVERLAY_VERSION=3.2.1.0

RUN apt-get update && apt-get install -y nginx xz-utils
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz

ENTRYPOINT ["/init"]
```

Depending on your target cpu architecture you may need to pull in different releases.

For every service that should run in parallel you then create a shell script
that contains instructions to run or stop the process.

The following example will start an nginx server and a dummy process.

In `services/nginx/` create a file `run`

```sh
#!/usr/bin/with-contenv sh
echo >&2 "Starting: 'nginx'"
exec /usr/sbin/nginx̃
```

A file `finish` in the same directory

```sh
#!/usr/bin/env sh
echo >&2 "Exit(code=${1}): 'nginx'"
```

In `services/hello-world/` create a file `run`

```sh
#!/usr/bin/with-contenv sh
echo >&2 "Starting: 'hello-world'"
exec sleep 3600
```

A file `finish` in the same directory

```sh
#!/usr/bin/env sh
echo >&2 "Exit(code=${1}): 'hello-world'"
```

This example will behave the same way as the `supervisord` and not terminate,
if all supervised process are stopped. Using `exec s6-svscanctl -t /var/run/s6/services`
at the end of a `finish` script can enable this behavior.
