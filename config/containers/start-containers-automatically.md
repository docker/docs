---
description: How to start containers automatically
keywords: containers, restart, policies, automation, administration
redirect_from:
- /engine/articles/host_integration/
- /engine/admin/host_integration/
- /engine/admin/start-containers-automatically/
title: Start containers automatically
---

Docker provides [restart policies](../../engine/reference/run.md#restart-policies---restart)
to control whether your containers start automatically when they exit, or when
Docker restarts. Restart policies ensure that linked containers are started in
the correct order. Docker recommends that you use restart policies, and avoid
using process managers to start containers.

Restart policies are different from the `--live-restore` flag of the `dockerd`
command. Using `--live-restore` allows you to keep your containers running
during a Docker upgrade, though networking and user input are interrupted.

## Use a restart policy

To configure the restart policy for a container, use the `--restart` flag
when using the `docker run` command. The value of the `--restart` flag can be
any of the following:

| Flag             | Description                                                                                     |
|:-----------------|:------------------------------------------------------------------------------------------------|
| `no`             | Do not automatically restart the container. (the default)                                       |
| `on-failure`     | Restart the container if it exits due to an error, which manifests as a non-zero exit code.     |
| `always`         | Always restart the container if it stops. If it is manually stopped, it is restarted only when Docker daemon restarts or the container itself is manually restarted. (See the second bullet listed in [restart policy details](#restart-policy-details)) |
| `unless-stopped` | Similar to `always`, except that when the container is stopped (manually or otherwise), it is not restarted even after Docker daemon restarts. |

The following example starts a Redis container and configures it to always
restart unless it is explicitly stopped or Docker is restarted.

```bash
$ docker run -d --restart unless-stopped redis
```

This command changes the restart policy for an already running container named `redis`.

```bash
$ docker update --restart unless-stopped redis
```

And this command will ensure all currently running containers will be restarted unless stopped.

```bash
$ docker update --restart unless-stopped $(docker ps -q)
```

### Restart policy details

Keep the following in mind when using restart policies:

- A restart policy only takes effect after a container starts successfully. In
  this case, starting successfully means that the container is up for at least
  10 seconds and Docker has started monitoring it. This prevents a container
  which does not start at all from going into a restart loop.

- If you manually stop a container, its restart policy is ignored until the
  Docker daemon restarts or the container is manually restarted. This is another
  attempt to prevent a restart loop.

- Restart policies only apply to _containers_. Restart policies for swarm
  services are configured differently. See the
  [flags related to service restart](../../engine/reference/commandline/service_create.md).


## Use a process manager

If restart policies don't suit your needs, such as when processes outside
Docker depend on Docker containers, you can use a process manager such as
[upstart](http://upstart.ubuntu.com/),
[systemd](http://freedesktop.org/wiki/Software/systemd/), or
[supervisor](http://supervisord.org/) instead.

> **Warning**
>
> Do not try to combine Docker restart policies with host-level process managers,
> because this creates conflicts.
{:.warning}

To use a process manager, configure it to start your container or service using
the same `docker start` or `docker service` command you would normally use to
start the container manually. Consult the documentation for the specific
process manager for more details.

### Using a process manager inside containers

Process managers can also run within the container to check whether a process is
running and starts/restart it if not.

> **Warning**
>
> These are not Docker-aware and just monitor operating system processes within
> the container. Docker does not recommend this approach, because it is
> platform-dependent and even differs within different versions of a given Linux
> distribution.
{:.warning}
