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

To declare jobs in a Compose file, see
[Run jobs in Compose](/manuals/compose/how-tos/jobs.md).

> [!IMPORTANT]
>
> Jobs are experimental. The daemon ships the feature as a built-in extension
> that is off by default, and its behavior might change between releases.

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

There is no `docker job` command to create a job. Declare jobs in a Compose file
and let Docker Compose register them with the daemon.

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
| `docker job prune`       | Remove idle manual jobs                                         |

`JOB` is the name or the ID of the job. Jobs that Docker Compose registers are
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

List the runs of one job:

```console
$ docker job runs jobs-demo.report
ITERATION  STATE   EXIT  TRIGGER  STARTED        FINISHED       CONTAINER
1          failed  1     manual   2 minutes ago  2 minutes ago  3f9a1c0d2b7e
```

### Run a job and read its logs

`docker job run` returns right away without streaming output. With
`--reschedule`, the run replaces the next scheduled run. To follow a run to its
end, use `docker job wait`, which exits with the exit code of the run. For jobs
that Docker Compose registers, `docker compose run` also starts the services
that the job depends on and streams the output:

```console
$ docker job run jobs-demo.report
$ docker job wait jobs-demo.report
```

To read the output, use `docker job logs`. With `--follow`, it keeps reading
and switches to the next run when the current one ends:

```console
$ docker job logs jobs-demo.report
$ docker job logs --follow jobs-demo.report
```

`docker job logs` prints the logs of the most recent run whose container still
exists. For when Docker Compose keeps containers, see
[Find out why a job failed](/manuals/compose/how-tos/jobs.md#find-out-why-a-job-failed).

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
