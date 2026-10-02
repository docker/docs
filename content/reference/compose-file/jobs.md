---
linkTitle: Jobs
title: Define jobs in Docker Compose
description: Explore all the attributes the jobs top-level element can have, including triggers and schedules.
keywords: compose, compose specification, jobs, triggers, schedule, cron, compose file reference
weight: 160
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Compose jobs" >}}

The top-level `jobs` element declares containers that run to completion, such
as database migrations and backups. A job runs when you start it with
`docker compose run`, on a cron schedule, or both. Docker Engine records every
execution of a job with its state and exit code.

`jobs` is a map whose keys are job names and whose values are job definitions.
Job names can contain letters, digits, `.`, `_`, and `-`. A job name can't
match the name of a service, because both are targets of `docker compose run`.

> [!IMPORTANT]
>
> Jobs are experimental and need the
> [jobs extension](/manuals/engine/daemon/jobs.md#turn-on-the-jobs-extension)
> turned on in Docker Engine.

For a walkthrough, see [Run jobs in Compose](/manuals/compose/how-tos/jobs.md).

## Example

The following example declares a manual job and a scheduled job. Docker Engine
runs `sync` every 5 minutes, waits for a run in progress to finish before it
starts the next one, and skips the runs it missed while it was stopped:

```yaml
jobs:
  migrate:
    image: example/migrations
    triggers:
      manual: true

  sync:
    image: example/sync
    triggers:
      schedule:
        - cron: "*/5 * * * *"
          concurrency: queue
          missed_fires: skip
```

## Attributes

A job definition accepts the following attributes.

### Service attributes

A job accepts every attribute of a [service](services.md) except the ones that
manage the lifecycle of a long-running service: `attach`, `container_name`,
`deploy`, `develop`, `extends`, `external_links`, `links`, `post_start`,
`pre_start`, `pre_stop`, `provider`, `restart`, and `scale`.

A job can set [`depends_on`](services.md#depends_on) to services and to other
jobs. `docker compose run` starts the services and runs the jobs that a job
depends on first. A job that another job depends on runs as a regular
container and has no run record. Docker Engine doesn't evaluate `depends_on`
when it starts a scheduled job. A service can't set `depends_on` to a job.

### `profiles`

`profiles` defines a list of named profiles for the job. A job with
`profiles` is only active when one of its profiles is active, unless you target
it with `docker compose run`. In that case, Compose activates its profiles. See
[Profiles](profiles.md).

### `triggers`

`triggers` defines when the job runs. It's required, and it must set at least
one of `manual` and `schedule`.

### `manual`

`triggers.manual` sets whether you can start the job with
`docker compose run`. It's a boolean:

- Not set: you can start the job manually. This is the default.
- `true`: same as not set. Compose rejects `manual: true` together with
  `schedule`.
- `false`: you can't start the job manually. `docker compose run` fails,
  including when another job depends on this one. A job with `manual: false`
  and no `schedule` never runs.

```yaml
jobs:
  rotate-logs:
    image: example/rotate
    triggers:
      manual: false
      schedule:
        - cron: "0 0 1 * *"
```

### `schedule`

`triggers.schedule` is a list of schedules on which Docker Engine runs the job.
Compose supports one entry for each job. An entry is either a cron expression
written as a string, or a mapping with the attributes that follow.

The following two jobs have the same schedule:

```yaml
jobs:
  short-syntax:
    image: example/job
    triggers:
      schedule:
        - "*/10 * * * *"

  long-syntax:
    image: example/job
    triggers:
      schedule:
        - cron: "*/10 * * * *"
```

A manual run of a scheduled job is an extra run. It doesn't change the time of
the next scheduled run.

Docker Engine parses the cron expression when `docker compose up` or
`docker compose run` registers the job, so an invalid expression fails then and
not when Compose loads the file. A second entry fails at the same moment.

#### `cron`

`cron` is required. It's a crontab expression with five fields separated by
spaces: minute, hour, day of the month, month, and day of the week.

| Field        | Accepted values               |
| ------------ | ----------------------------- |
| Minute       | 0-59                          |
| Hour         | 0-23                          |
| Day of month | 1-31                          |
| Month        | 1-12                          |
| Day of week  | 0-7, where 0 and 7 are Sunday |

Each field accepts numbers, lists such as `1,15`, ranges such as `1-5`, and
steps such as `*/10` or `0-30/10`. Names such as `mon` or `jan` and shortcuts
such as `@daily` aren't supported. When both the day of the month and the day
of the week are restricted, the job runs on days that match either field.

#### `timezone`

`timezone` is the IANA time zone name that Docker Engine uses to evaluate
`cron`, for example `Europe/Paris`. When it's not set, Docker Engine uses UTC.

#### `concurrency`

`concurrency` sets what happens when the schedule fires while the previous run
of the job is still in progress. The value is one of:

- `forbid` skips the new run. This is the default.
- `queue` waits for the current run to finish and then starts a single new
  run.

#### `missed_fires`

`missed_fires` sets what happens to runs that were due while Docker Engine was
stopped. The value is one of:

- `one` starts a single catch-up run when Docker Engine starts again. This is
  the default.
- `skip` drops the missed runs.
