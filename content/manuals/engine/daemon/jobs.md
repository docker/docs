---
title: Manage jobs in Docker Engine
linkTitle: Jobs
description: Turn on the experimental jobs extension in Docker Engine, then use the docker job commands to list, run, and inspect jobs.
keywords: docker engine, daemon, jobs, extension, enable-extensions, docker job, cron, scheduled jobs
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

Creating a job happens in Docker Compose; `docker job` then manages it. You declare jobs with the top-level [`jobs`](/reference/compose-file/jobs.md) element in a Compose file, and `docker compose up` registers them with the daemon. No `docker job` command creates a job.

This page covers turning on the jobs extension and using the `docker job`
commands to inspect and control jobs. To declare jobs in a Compose file, see
[Run jobs in Compose](/manuals/compose/how-tos/jobs.md).

> [!IMPORTANT]
>
> Jobs are experimental. The behavior of jobs
> might change between releases.
>
> Docker Engine ships jobs as an extension that's off by default. Until you
> [turn it on](#turn-on-the-jobs-extension), both Compose and `docker job`
> fail with an error saying that the daemon doesn't support jobs.

## Turn on the jobs extension

> [!WARNING]
>
> The daemon serves the jobs API over gRPC on every address it listens on,
> including TCP and TLS. Authorization plugins only evaluate requests to the
> HTTP API, so they don't apply to jobs. If you rely on an authorization
> plugin, take this into account before you turn on the extension.

The jobs extension has the ID `org.mobyproject.jobs.v1`. Turn it on by passing
the ID to `dockerd`:

```console
$ dockerd --enable-extension org.mobyproject.jobs.v1
```

Or add it to the `enable-extensions` list in the
[`daemon.json` file](/manuals/engine/daemon/_index.md#configuration-file):

```json
{
  "enable-extensions": ["org.mobyproject.jobs.v1"]
}
```

Restart the daemon after you change the list.

To check that the extension is on, list the jobs:

```console
$ docker job ls
```

With the extension on, the command prints a table, which is empty until a job
exists. With the extension off, the command fails with an error saying that
the daemon doesn't support jobs.

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

## Manage jobs with the CLI

The `docker job` commands inspect and control jobs that already exist. To
create one, declare it in a Compose file and run `docker compose up`. See
[Run jobs in Compose](/manuals/compose/how-tos/jobs.md).

| Command                  | Description                                                     |
| ------------------------ | --------------------------------------------------------------- |
| `docker job ls`          | List jobs with their trigger, state, and last run               |
| `docker job runs JOB`    | List the runs of a job, paged with `--limit` and `--before`     |
| `docker job inspect JOB` | Print a job as JSON, or a single run with `--run`               |
| `docker job run JOB`     | Start a run now and print its ID                                |
| `docker job logs JOB`    | Print the logs of the most recent run                           |
| `docker job wait JOB`    | Wait for a run to finish and exit with its exit code            |
| `docker job cancel JOB`  | Cancel the run in progress                                      |
| `docker job pause JOB`   | Stop a job from firing on its schedule                          |
| `docker job resume JOB`  | Let a paused job fire on its schedule again                     |
| `docker job rm JOB`      | Remove a job. The `--runs` option sets what happens to its runs |
| `docker job prune`       | Remove idle manual jobs. Leaves scheduled job.                  |

`JOB` is the name or the ID of the job.  Jobs that Docker Compose registers are
named `<PROJECT>.<JOB>`, for example `jobs-demo.backup`.

### List jobs and runs

List all jobs:

```console
$ docker job ls
NAME              TRIGGER    STATE  PAUSED  LAST RUN                    EXIT  NEXT FIRE
jobs-demo.backup  * * * * *  idle   false   succeeded (12 seconds ago)  0     in 47 seconds
jobs-demo.report  manual     idle   false   failed (2 minutes ago)      1     -
```

The `NEXT FIRE` column shows `-` for manual jobs, for paused jobs, and while a
run is in progress.

The daemon writes the run record before it creates the container, so the
record and its exit code survive both the removal of the container and a
daemon restart.

By default, the daemon keeps up to 10 000 run records for each job and removes
the oldest finished runs first.

List the runs of one job:

```console
$ docker job runs jobs-demo.report
ITERATION  STATE   EXIT  TRIGGER  STARTED        FINISHED       CONTAINER
1          failed  1     manual   2 minutes ago  2 minutes ago  3f9a1c0d2b7e
```

### Start a run on demand

`docker job run` starts a run and returns immediately, printing the run ID. It doesn't stream output and doesn't wait for the run to finish. Pair it with `docker job wait`, which blocks until the run ends and exits with the run's exit code:

```console
$ docker job run jobs-demo.report
$ docker job wait jobs-demo.report
```

Reach for `docker job run` when the job name is all you have: from a machine without the project files, against a remote Engine, or from a script.

If you're working in the Compose project, use `docker compose run` instead. It starts the services the job depends on, streams the output, and exits with the run's exit code, all in one command.

By default a manual run of a scheduled job is an extra run and leaves the schedule untouched. Pass `--reschedule` to make it count as the next scheduled run instead.

### Read the logs of a run

`docker job logs` prints the output of the most recent run whose container still exists:

```console
$ docker job logs jobs-demo.report
```

Add `--follow` to keep reading, and to switch to the next run when the current one ends:

```console
$ docker job logs --follow jobs-demo.report
```

Logs live in the run's container, so they're gone once the container is removed, even though the run record survives. For which containers Docker Compose keeps and which it removes, see [Find out why a job failed](/manuals/compose/how-tos/jobs.md#find-out-why-a-job-failed).

### Find the containers of a job

Run containers carry the labels `com.docker.job.id` and
`com.docker.job.run-id`. Use them to find the containers of a job with the
standard container commands:

```console
$ docker ps --all --filter label=com.docker.job.id=<JOB_ID>
```

The ID of a job is the `ID` field in the output of `docker job inspect`.

### Remove jobs

`docker job rm` removes a job. The `--runs` option sets what happens to the run
history:

- `keep` keeps all runs. This is the default.
- `remove` removes all runs.
- `remove-finished` removes only the runs that have finished.

`docker job prune` removes idle manual jobs, including the ones that Docker
Compose registers. Use `--label` to limit the removal to jobs that have a given
label, such as `--label com.docker.compose.project=<PROJECT>`.
