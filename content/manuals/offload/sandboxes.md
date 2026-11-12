---
title: Manage Docker Offload sandboxes
linktitle: Manage Offload sandboxes
weight: 30
description: Create, list, stop, resume, and delete the cloud sandbox that Docker Offload uses to run your containers.
keywords: offload, cloud sandbox, context, idle timeout, keep alive, resume
---

{{< summary-bar feature_name="Docker Offload" >}}

A cloud sandbox is an isolated environment on Docker-managed hardware, running
its own Docker daemon. Docker Offload runs your builds and containers there
instead of on your machine, using the Docker commands you already know. For why
you'd run workloads this way, see [Core concepts](core-concepts.md).

A cloud sandbox persists between sessions. Stopping it
pauses it rather than discarding it, so its containers, images, and volumes are
still there when you come back. This page covers that lifecycle, including when
a sandbox pauses, how to control it, and what's available while it's paused. It
also covers working with local and cloud containers at the same time.

You use one Docker Offload cloud sandbox per account at a time. You can connect to it
from more than one machine.

### Choose the size of a new cloud sandbox

```console
$ docker offload create --cpus 8 --memory 16g --ttl 4h
```

| Flag | Description |
| --- | --- |
| `--cpus <number>` | The number of CPUs: 1, 2, 4, 8, or 16 |
| `--memory <size>` | Memory, as a size such as `16g` |
| `--ttl <duration>` | How long the cloud sandbox keeps running after your machine stops renewing it, as a duration such as `4h` |

These flags apply only when Docker Offload creates a sandbox. If a cloud sandbox is
already running, they're ignored, and `docker offload create` tells you so.

To change the size of a cloud sandbox you already have, create it again with
`--recreate`.

> [!WARNING]
>
> `--recreate` deletes the existing cloud sandbox and everything on it, including
> containers, images, and volumes.

## Check the status of your cloud sandbox

```console
$ docker offload status

  Account: <your-account>

  Sandbox
    ID: sbx_001m49455jm42x4v6epbg8kd5fy
    State: Running
    Created: 6 days ago

  This machine
    State: Started (Active)
    Session started: 21 hours ago
    Idle detection: on (5m)
    Kept awake by: checking activity
```

`Account` is the Docker account the sandbox belongs to. The rest is grouped to
separate the sandbox itself from your connection to it.

Sandbox:

| Field | Description |
| --- | --- |
| ID | The sandbox ID |
| State | Whether the sandbox is running or paused |
| Created | When the sandbox was created |

This machine:

| Field | Description |
| --- | --- |
| State | Whether this machine is connected to the sandbox, and whether it's active |
| Session started | How long this machine has been connected |
| Idle detection | Whether idle detection is on, and the idle timeout in effect. The default is 5m |
| Time until idle | When this machine pauses the sandbox if no activity happens before then |
| Kept awake by | What's keeping the sandbox awake, such as a running container or a Docker command from this machine |

`Time until idle` and `Kept awake by` are alternatives. The countdown appears
while the sandbox is heading for a pause, and `Kept awake by` appears in its
place when something is holding it open.

Your sandbox also appears in the [Docker Agentic Platform
Console](https://agentic-platform.docker.com/) under your account, with a name
that starts with `offload-`. For how it's reported there, see [Docker Offload
usage](usage.md#find-your-sandbox-in-the-console).

## Control when your cloud sandbox pauses

Docker Offload pauses an idle sandbox so that you aren't billed for time you
aren't using. A paused sandbox keeps its containers, images, and volumes.

Set the idle timeout when you create or connect to a sandbox:

```console
$ docker offload create --idle-timeout 30m
```

| Flag | Behavior |
| --- | --- |
| `--idle-timeout <duration>` | Pauses the sandbox after it has been quiet for `<duration>`. Defaults to 5m when you don't pass the flag. The minimum is 1m and there's no maximum. A value of 0 behaves like `--keep-alive`. Values between 0 and 1m are rejected |
| `--keep-alive` | Turns off idle detection, so the sandbox runs until you stop it. Can't be combined with `--idle-timeout` |

Either flag applies to a running sandbox straight away, without restarting it.
`docker offload create` returns once the new timeout is in effect.

Each machine remembers the idle timeout you last used for a sandbox. If you
leave the flag off, Docker Offload reuses that value rather than resetting to
the default. A new sandbox starts at 5m.

> [!WARNING]
>
> With `--keep-alive`, the sandbox keeps running, and billing continues, until
> you run `docker offload stop`. It still pauses if your machine sleeps past
> the sandbox's time to live, and after 24 hours. For more information, see
> [Docker Offload usage](usage.md).

Idle detection runs on the machine that's connected to the sandbox, so it only
applies while that machine is awake and online. If your machine sleeps or goes
offline, the sandbox's time to live takes over instead.

If `idle_timeout` is set in a config override file, it takes priority over both
flags, and `docker offload create` tells you when that happens.

### What counts as activity

A sandbox is idle when both of these are true:

- Nothing is working inside the sandbox. Docker Offload samples the sandbox
  regularly and treats container or Docker daemon activity as work.
- No command, build, or pull from your machine is in progress. An attached
  session or a forwarded port stops counting after a minute with no data
  passing through it.

Some workloads read as idle even though you're not finished with them. A job
that only waits on a remote call is the most common. Raise `--idle-timeout` or
use `--keep-alive` for those.

The reverse also happens. A container running on the host network keeps the
sandbox awake, and so does a large number of running containers, because
listing them is itself work. A sandbox in that state doesn't pause on its own,
so stop it when you're finished.

## Set how long a sandbox lives

The `--ttl` flag sets the sandbox's time to live. The maximum time is 24 hours. It's independent of idle
detection:

```console
$ docker offload create --ttl 8h
```

While a machine is connected, Docker Offload keeps renewing the time to live,
so it counts down only once no machine is connected. In practice, the idle
timeout handles the case where you're at your machine but not running anything,
and the time to live handles the case where your machine has gone away.

## Stop a sandbox

```console
$ docker offload stop
```

Your Docker context switches back to `desktop-linux` immediately, and the
sandbox pauses in the background. Pausing can take several minutes to finish.

A paused sandbox keeps its containers, images, and volumes. Billing stops when
the sandbox starts pausing. For more information, see [Docker Offload
usage](usage.md).

Stopping pauses the sandbox for every machine on the account, not only the one
you run the command on.

## While your sandbox is paused

A paused sandbox keeps its containers, images, and volumes, but it isn't
serving anything until it resumes. While it's paused:

- `docker ps` and `docker images` list nothing. Your containers and images are
  still there and appear again after the sandbox resumes.
- Forwarded ports and file sync don't work.
- Attached sessions and exec sessions have ended. Reattach with `docker attach`
  or start a new `docker exec` after the sandbox resumes.

A sandbox stays paused when your machine wakes, reconnects, signs in, or
restarts Docker Desktop. It resumes when you run a command that needs it.

## Resume a sandbox

You don't resume a sandbox explicitly. Running a command that needs the sandbox
or one of its containers resumes it, including `docker run`, `docker start`,
`docker exec`, `docker attach`, `docker logs`, and `docker inspect`. `docker
offload create` resumes it too.

Commands that don't need the sandbox leave it paused, including `docker ps`,
`docker images`, and `docker offload status`.

Resuming restores the sandbox's previous state. Containers that were running
resume, and data written to volumes and container filesystems is preserved.

If the sandbox is still pausing when you run a command that needs it, Docker
Offload waits for the pause to finish and then resumes the sandbox. If the
pause takes longer than the wait, Docker Offload tells you that the sandbox is
still pausing and asks you to try again shortly.

## Delete a sandbox

```console
$ docker offload rm
```

Docker Offload asks for confirmation first. Pass `--force` to skip the prompt.

Deleting a sandbox removes everything on it, including containers, images, and
volumes, whether it's running or paused. Every machine on the account uses the
same sandbox, so `docker offload rm` deletes it for all of them.

The next `docker offload create` makes a new sandbox, which starts with the
default idle timeout of 5m.

## Connect from more than one machine

You can connect to the same sandbox from more than one machine. Running `docker
offload create` on a second machine connects to your existing sandbox rather
than creating another one, so both machines report the same account and the
same sandbox ID.

Each connection is tracked separately, which means:

- `Session started` counts from when that machine connected, so it restarts on
  a newly connected machine while an earlier machine keeps its own count.
- Idle detection is managed independently by each machine, and the machines
  don't share that state. `Time until idle` can differ between them, and each
  machine can use a different idle timeout.

Because each machine applies its own idle timeout, another machine can pause
the sandbox on its timeout while you're still using it. If you work from two
machines, set the idle timeout you want on both.

> [!IMPORTANT]
>
> Stopping and deleting aren't per-machine. Every machine on the account uses
> the same sandbox, so `docker offload stop` pauses it for all of them, and
> `docker offload rm` deletes it for all of them.

## Work with local and cloud containers

Creating a sandbox doesn't stop your local containers. It changes which Docker
daemon your commands reach.

While your context is `docker-cloud`:

- `docker ps -a` lists containers in the sandbox. Your local containers don't
  appear.
- Images you built locally aren't present in the sandbox, and images you build
  in the sandbox aren't present locally.

To run a command against your local daemon without stopping Docker Offload,
pass the context explicitly:

```console
$ docker --context desktop-linux ps -a
$ docker --context desktop-linux compose up
```

To send a single command to the cloud while your context stays local, use the
`--cloud` flag:

```console
$ docker --cloud build .
```

`--cloud` is a global flag on the `docker` command, so it goes before the
subcommand. Passing it after the subcommand fails:

```console
$ docker run --cloud --rm hello-world
unknown flag: --cloud
```

This applies to every subcommand, not only `run`.

`--cloud` applies to that one command. Unlike `docker offload create`, it
doesn't change the context your other commands use, so the rest of your
terminal keeps working locally. If you don't have a sandbox yet, `--cloud`
creates one and waits for it to be ready, the same as `create` does.

If the sandbox can't be provisioned, the command fails. `--cloud` never falls
back to running locally, so a build or run you meant for the cloud doesn't
silently use your own machine.

### Ports and files

Docker Offload forwards ports from containers in the sandbox to your machine,
and syncs bind-mounted files, so published ports behave as they do locally.
Forwarded ports and file sync stop working while the sandbox is paused, and
start again when it resumes.

This differs from using cloud sandboxes directly through the `sbx` CLI, where
there's no path back to your machine. See [Compare local and cloud
sandboxes](/manuals/ai/sandboxes/cloud/local-vs-cloud.md).

<!-- TODO: Confirm how a port published in the sandbox is addressed from the
host: whether it's localhost:<port> as with a local container, and whether
anything differs for ports already in use locally. -->

<!-- TODO: Confirm which bind mounts sync, and in which direction. -->

Ports published by your local containers aren't reachable through the `docker-cloud`
context. Switch context or use `--context desktop-linux` to reach them.

## Use Docker Compose in a cloud sandbox

Compose commands behave the same way in a sandbox as they do locally, and
Compose state survives a pause:

- `docker compose stop` keeps the containers and their data. `docker compose
  up` resumes them.
- `docker compose down` removes the containers. `docker compose up` then
  creates containers without the previous data.

Because the sandbox persists while paused, a stack you bring up in one session
is still there when you resume, including data written to its volumes.
