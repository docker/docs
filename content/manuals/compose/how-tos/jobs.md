---
title: Run jobs in Compose
linkTitle: Run jobs
weight: 140
description: Declare one-shot and scheduled workloads as jobs in a Compose file, run them on demand or on a cron schedule, and check the result of each run.
keywords: docker compose jobs, compose jobs, cron, scheduled jobs, one-shot, migrations, backups, docker job
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Compose jobs" >}}

Jobs let you declare one-shot work tasks, such as database migrations and backups, in
your Compose file alongside your services. You run a job on demand with
`docker compose run`, or you give it a cron schedule and let Docker Engine run
it for you. Each execution is recorded with its state and exit code, so you can
check the result after the container is gone.

> [!IMPORTANT]
>
> Jobs are experimental. The Compose file syntax and the behavior of jobs
> might change between releases.
>
> Docker Engine ships jobs as an extension that's off by default. You must
> [turn on the jobs extension](/manuals/engine/daemon/jobs.md#turn-on-the-jobs-extension)
> before you declare jobs in a Compose file. While the extension is off,
> Compose and `docker job` fail with an error.

## When to use jobs

A job is a container that runs to completion. You define it with the top-level
[`jobs`](/reference/compose-file/jobs.md) element with the attributes of a
service, plus a `triggers` block that says when the job runs.

Without jobs, one-shot work ends up as a service in a profile that you start
with `docker compose run`. That approach has two limits: Docker Engine
discards the result once the container is removed, and the work only runs when
you run the command.

Jobs address both. Docker Engine keeps a run record for every execution,
including the exit code, and evaluates each schedule itself. A nightly backup
keeps firing after the `docker compose` command exits, and it keeps firing when
your Compose project targets a remote engine.

| Feature                        | Service             | `docker compose run` | Job                       |
| ------------------------------ | ------------------- | -------------------- | ------------------------- |
| Lifetime                       | Runs until stopped  | Runs to completion   | Runs to completion        |
| Result after container removal | None                | None                 | Run record with exit code |
| Triggered by                   | `docker compose up` | You                  | You, or a cron schedule   |

Jobs don't fit every use case:

- Work that must finish before a service starts, such as a migration that
  the application waits for. Use [init containers](init-containers.md)
  instead: a service's `depends_on` accepts only other services, so a service
  can never wait for a job.
- Processes that run until you stop them. Use a service instead.

## Prerequisites

Turn on the jobs extension in Docker Engine. Add the following to the
`daemon.json` file, then restart the daemon:

```json
{
  "enable-extensions": ["org.mobyproject.jobs.v1"]
}
```

To check that the extension is on, run `docker job ls`. With the extension on,
the command prints a table header and no jobs:

```console
$ docker job ls
NAME  TRIGGER  STATE  PAUSED  LAST RUN  EXIT  NEXT FIRE
```

For the other ways to turn on the extension, and for how it interacts with
authorization plugins, see
[Manage jobs in Docker Engine](/manuals/engine/daemon/jobs.md#turn-on-the-jobs-extension).

## Quickstart

This quickstart walks you through running a PostgreSQL database with three jobs:

- `migrate` creates and fills a table. You run it on demand.
- `backup` dumps the database every minute.
- `report` fails on purpose, so you can see what a failed run looks like.

1. Create a project directory with a `backups` directory:

   ```console
   $ mkdir -p jobs-demo/backups
   $ cd jobs-demo
   ```

1. Create a `compose.yaml` file:

   ```yaml
   name: jobs-demo

   x-db-env: &db-env
     DB_URL: postgresql://postgres:demo@db/shop

   services:
     db:
       image: postgres:16-alpine
       environment:
         POSTGRES_PASSWORD: demo
         POSTGRES_DB: shop
       healthcheck:
         test: ["CMD-SHELL", "pg_isready -U postgres -d shop"]
         interval: 2s
         timeout: 3s
         retries: 15

   jobs:
     # Manual: creates a table and fills it the first time it runs.
     migrate:
       image: postgres:16-alpine
       environment: *db-env
       command:
         - sh
         - -c
         - |
           set -e
           psql "$$DB_URL" -q -v ON_ERROR_STOP=1 \
             -c "CREATE TABLE IF NOT EXISTS orders (id serial PRIMARY KEY, amount numeric(10, 2) NOT NULL)" \
             -c "INSERT INTO orders (amount) SELECT round((random() * 200)::numeric, 2) FROM generate_series(1, 100) WHERE NOT EXISTS (SELECT 1 FROM orders)"
           echo "migration applied"
           psql "$$DB_URL" -c "SELECT count(*) AS orders FROM orders"
       triggers:
         manual: true
       depends_on:
         db:
           condition: service_healthy

     # Scheduled: waits for the database, then writes a timestamped dump to
     # ./backups every minute.
     backup:
       image: postgres:16-alpine
       environment: *db-env
       volumes:
         - ./backups:/backups
       command:
         - sh
         - -c
         - |
           until pg_isready -q -d "$$DB_URL"; do sleep 1; done
           pg_dump --no-owner -f /backups/shop-$$(date +%H%M%S).sql "$$DB_URL"
       triggers:
         schedule:
           - cron: "* * * * *"

     # Manual: fails on purpose because the invoices table doesn't exist.
     report:
       image: postgres:16-alpine
       environment: *db-env
       command: sh -c 'psql "$$DB_URL" -v ON_ERROR_STOP=1 -c "SELECT sum(amount) FROM invoices"'
       triggers:
         manual: true
       depends_on:
         db:
           condition: service_healthy
   ```

1. Start the project:

   ```console
   $ docker compose up -d
   ```

   Compose starts the `db` service and registers `backup` with Docker Engine.
   It also prints a warning that lists `migrate` and `report`, because `up`
   doesn't run manual jobs.

1. List the jobs:

   ```console
   $ docker job ls
   NAME              TRIGGER    STATE  PAUSED  LAST RUN  EXIT  NEXT FIRE
   jobs-demo.backup  * * * * *  idle   false   -         -     in 26 seconds
   ```

   Docker Engine counts down to the next run of `backup`. The manual jobs
   appear after their first run.

1. Run the migration:

   ```console
   $ docker compose run migrate
   migration applied
    orders
   --------
       100
   (1 row)
   ```

   Compose prints progress messages while it waits for `db` to be healthy,
   then streams the output of the run.

1. Wait a minute, then list the dumps that `backup` wrote on its own:

   ```console
   $ ls backups
   shop-133000.sql
   ```

   A new dump appears every minute. The file names contain the time in UTC.

1. Run the report and check its exit code:

   ```console
   $ docker compose run report
   ERROR:  relation "invoices" does not exist
   LINE 1: SELECT sum(amount) FROM invoices
                                   ^
   $ echo $?
   1
   ```

   `docker compose run` returns the exit code of the job.

1. Review the jobs, the runs of the failed job, and its logs:

   ```console
   $ docker job ls
   NAME               TRIGGER    STATE  PAUSED  LAST RUN                        EXIT  NEXT FIRE
   jobs-demo.backup   * * * * *  idle   false   succeeded (23 seconds ago)      0     in 36 seconds
   jobs-demo.migrate  manual     idle   false   succeeded (About a minute ago)  0     -
   jobs-demo.report   manual     idle   false   failed (About a minute ago)     1     -
   $ docker job runs jobs-demo.report
   ITERATION  STATE   EXIT  TRIGGER  STARTED             FINISHED            CONTAINER
   1          failed  1     manual   About a minute ago  About a minute ago  6fc96fd5b02d
   $ docker job logs jobs-demo.report
   ERROR:  relation "invoices" does not exist
   LINE 1: SELECT sum(amount) FROM invoices
                                   ^
   ```

1. Remove the jobs, the services, and the failed container:

   ```console
   $ docker compose down --remove-orphans
   ```

   The `backup` job stops firing. The dump files stay in the `backups`
   directory.

## Run a job on demand

`docker compose run JOB` starts a run. Compose:

- Starts the services listed in the job's `depends_on`, waits for their conditions, and runs any jobs it depends on to completion.
- Starts the run and streams the container output to your terminal.
- Exits with the exit code of the run.

| Result of the run                   | Exit code of `docker compose run` |
| ----------------------------------- | --------------------------------- |
| Succeeded                           | `0`                               |
| Failed                              | The exit code of the container    |
| Failed without starting a container | `1`                               |
| Canceled with `docker job cancel`   | `130`                             |

Every job accepts `docker compose run`, including scheduled jobs. To prevent
manual runs of a job, set
[`manual: false`](/reference/compose-file/jobs.md#manual).

Commands that take service names, such as `up JOB`, `start`, `stop`, and
`create`, reject a job name and tell you to use `docker compose run`.

## Run a job on a schedule

`docker compose up` registers every active job that has a `schedule` with
Docker Engine. `docker compose up --no-start` registers scheduled jobs without
starting any service.

The following job takes a backup every night at 03:00 in the Europe/Paris time
zone. `missed_fires: one` covers downtime: if Docker Engine is stopped at
03:00, the job runs once as soon as the Engine starts again, so the night still
gets its backup:

```yaml
jobs:
  backup:
    image: example/backup
    triggers:
      schedule:
        - cron: "0 3 * * *"
          timezone: Europe/Paris
          missed_fires: one
```

A cron expression has five fields: minute, hour, day of the month, month, and
day of the week. Some common expressions:

| Expression     | Runs                               |
| -------------- | ---------------------------------- |
| `*/15 * * * *` | Every 15 minutes                   |
| `0 3 * * *`    | Every day at 03:00                 |
| `0 6 * * 1-5`  | At 06:00 from Monday to Friday     |
| `30 2 1 * *`   | At 02:30 on the first of the month |

Without `timezone`, Docker Engine evaluates the expression in UTC. For the
accepted cron syntax and for the `concurrency` and `missed_fires` settings, see
the [schedule reference](/reference/compose-file/jobs.md#schedule).

Docker Engine doesn't start or wait for the services in `depends_on` when it
fires a scheduled job, so make sure those services are running.

Pause a schedule with `docker job pause` and restart it with
`docker job resume`:

```console
$ docker job pause jobs-demo.backup
$ docker job resume jobs-demo.backup
```

## Change a registered job

Registered jobs are immutable. To change a job, remove it and register it
again.

When you edit a job's definition, or when the image it uses changes,
`docker compose up` and `docker compose run` fail with a `job ... has changed`
error. To register the new definition:

1. Run `docker compose down`, or `docker job rm <PROJECT>.<JOB>` to remove a
   single job and leave the rest of the project registered.
1. Run the `up` or `run` command again.

## Remove jobs and failed containers

`docker compose down` removes the jobs of the project from Docker Engine, along
with the containers and networks it already removes. The schedules stop firing,
and Docker Engine keeps the run history of the removed jobs.

`down` leaves the containers of failed runs. To remove them as well, run
`docker compose down --remove-orphans`.

## Troubleshoot

### Find out why a job failed

Jobs that Compose registers are named `<PROJECT>.<JOB>`, for example
`jobs-demo.backup`.

1. Run `docker job ls` to see the latest run of every job. The `LAST RUN`
   column gives its state, and `EXIT` gives its exit code.
1. Run `docker job runs <PROJECT>.<JOB>` to see the state and exit code of
   every run, not only the latest.
1. Run `docker job logs <PROJECT>.<JOB>` to read the output of the latest run
   that still has a container.

Compose registers jobs with a retention policy: Docker Engine keeps the
container of a failed run and removes the container of a successful run. A
successful run keeps its record and its exit code, but `docker job logs` fails
with `no logs available`.

Compose also limits the run history of each job to the five most recent runs,
rather than the Docker Engine default of 10 000.

A run with state `failed` and no exit code usually means that no container
ran. To see why, run `docker job inspect <PROJECT>.<JOB>` and read the `Error`
field of the latest run.

For more `docker job` commands, see
[Manage jobs with the CLI](/manuals/engine/daemon/jobs.md#manage-jobs-with-the-cli).

### Find out why a schedule doesn't fire

Run `docker job ls` and check the `PAUSED` column. A paused job doesn't fire
until you run `docker job resume`. A scheduled job appears in the list after
`docker compose up` registers it.

## Reference information

- [`jobs`](/reference/compose-file/jobs.md)
- [Manage jobs with the CLI](/manuals/engine/daemon/jobs.md#manage-jobs-with-the-cli).
