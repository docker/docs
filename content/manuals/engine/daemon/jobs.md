---
title: Manage jobs in Docker Engine
linkTitle: Jobs
description: Turn on the experimental jobs extension in Docker Engine so that Docker Compose can run one-shot and cron-scheduled jobs.
keywords: docker engine, daemon, jobs, extension, extension-enable, cron, scheduled jobs
weight: 50
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Docker Engine jobs" >}}

A job is a one-shot workload that the Docker daemon runs on demand or on a cron
schedule. The daemon keeps a record of every execution, so you can check
whether a backup or a migration succeeded after its container is gone. Because
the daemon evaluates schedules itself, scheduled jobs keep running when no
client is connected.

You declare jobs with the top-level [`jobs`](/reference/compose-file/jobs.md)
element in a Compose file. Docker Compose registers scheduled jobs with the
daemon when you run `docker compose up`, and creates manual jobs when you run
`docker compose run`.

This page covers turning on the jobs extension. To declare jobs in a Compose
file, see [Run jobs in Compose](/manuals/compose/how-tos/jobs.md).

> [!IMPORTANT]
>
> Jobs are experimental. The behavior of jobs
> might change between releases.
>
> Jobs need Docker Engine
> [29.10.0](/manuals/engine/release-notes/29.md#29100) or later (API version
> 1.57). Docker Engine ships jobs as an extension that's off by default. Until
> you [turn it on](#turn-on-the-jobs-extension), `docker compose up`,
> `create`, and `start` refuse a project that declares scheduled jobs.

## Turn on the jobs extension

> [!WARNING]
>
> The daemon serves the jobs API over gRPC on every address it listens on,
> including TCP and TLS. Authorization plugins only evaluate requests to the
> HTTP API, so they don't apply to jobs. If you rely on an authorization
> plugin, take this into account before you turn on the extension.

The jobs extension has the ID `org.mobyproject.jobs.v0`. Turn it on by passing
the ID to `dockerd`:

```console
$ dockerd --extension-enable org.mobyproject.jobs.v0
```

Or add it to the `extensions.enable` list in the
[`daemon.json` file](/manuals/engine/daemon/_index.md#configuration-file):

```json
{
  "extensions": {
    "enable": ["org.mobyproject.jobs.v0"]
  }
}
```

Restart the daemon after you change the list.

To check that the extension is on, run `docker compose up` or
`docker compose create` in a project that declares a scheduled job. With the
extension off, the command fails with an error saying that scheduled jobs need
the jobs extension.

## What happens when the extension is off

With the extension off, Docker Compose behaves as follows:

- `docker compose up`, `create`, and `start` refuse only a project that
  declares scheduled jobs. A project with manual jobs only starts normally.
- `docker compose run <job>` falls back to running the job as a one-off
  service and prints a warning. Docker Engine keeps no run record for it.
- `docker compose down` ignores jobs.

## How jobs work

A job is a container definition with a trigger. The trigger is either manual,
which means the job runs only when you ask for it, or a cron schedule.

Each execution of a job is a run. A run starts as `pending`, moves to
`running`, and ends as `succeeded`, `failed`, `timed_out`, or `cancelled`. The
daemon writes the run record before it creates the container. The record,
including the exit code, survives the removal of the container and daemon
restarts.

Registering a job again with the same name and definition changes nothing. To
change a job, remove it and register it again.

After a restart, the daemon resumes tracking the runs that were in progress and
schedules the next runs. The `missed_fires` setting of each schedule decides
whether the daemon runs a job that was due while it was stopped. For the
scheduling settings, see the
[Compose file reference](/reference/compose-file/jobs.md#schedule).

By default, the daemon keeps up to 10 000 run records for each job and removes
the oldest finished runs first.

Run containers carry the labels `com.docker.job.id` and
`com.docker.job.run-id`. Use them to find job containers with the standard
container commands:

```console
$ docker ps --all --filter label=com.docker.job.id
```
