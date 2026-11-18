---
title: Run jobs in Compose
linkTitle: Run jobs
weight: 140
description: Declare one-shot and scheduled workloads as jobs in a Compose file, run them on demand or on a cron schedule, and check the result of each run.
keywords: docker compose jobs, compose jobs, cron, scheduled jobs, one-shot, migrations, backups
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
it for you. Docker Engine records each execution with its state and exit code.

> [!IMPORTANT]
>
> Jobs are experimental. The Compose file syntax and the behavior of jobs
> might change between releases.
>
> Docker Engine ships jobs as an extension that's off by default. You must
> [turn on the jobs extension](/manuals/engine/daemon/jobs.md#turn-on-the-jobs-extension)
> before you declare scheduled jobs in a Compose file. While the extension is
> off, `docker compose up`, `create`, and `start` refuse a project that
> declares scheduled jobs. `docker compose run` on a job falls back to a
> one-off service and prints a warning, so Docker Engine keeps no run record.

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

- Docker Compose 5.8.0 or later
- Docker Engine 29.10.0 or later
- [Jobs extensions in Docker Engine](/manuals/engine/daemon/jobs.md#turn-on-the-jobs-extension)

To turn on the jobs extension in Docker Engine, add the following to the
`daemon.json` file, then restart the daemon:

```json
{
  "extensions": {
    "enable": ["org.mobyproject.jobs.v0"]
  }
}
```

To check that the extension is on, run `docker compose up` in a project that
declares a scheduled job. With the extension off, the command fails with an
error that tells you to start `dockerd` with `--extension-enable`.

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

1. Find the container of the failed run:

   ```console
   $ docker compose ps -a
   ```

   Compose keeps the container of a failed run, so `report` shows up in the
   list with an `Exited (1)` status. Read its output with `docker logs` and the
   name of the container from the list:

   ```console
   $ docker logs <CONTAINER>
   ```

   The `migrate` and `backup` runs succeeded, so their containers are already
   gone.

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
| Timed out                           | `124`                             |
| Canceled                            | `130`                             |

Every job accepts `docker compose run`, including scheduled jobs. To prevent
manual runs of a job, set
[`manual: false`](/reference/compose-file/jobs.md#manual).

A manual job doesn't stay registered with Docker Engine. Compose creates it
when you run `docker compose run` and removes it from the engine after the run.
Docker Engine keeps the run record, and Compose keeps the container of a failed
run.

A scheduled job must already be registered, so run `docker compose up` or
`docker compose create` first. `docker compose run` runs a scheduled job
exactly as `up` registered it, so it refuses options that change the job, such
as `-e`, `-v`, `--user`, `--entrypoint`, `--workdir`, `--label`, and a command
override. To change the job, edit the Compose file and run `up` again.

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

## Change a registered job

Registered jobs are immutable. To change a job, remove it and register it
again.

When you edit the definition of a scheduled job, or when the image it uses
changes, `docker compose up` fails with a `job ... has changed` error.
`docker compose run` doesn't compare the definition: it runs a scheduled job as
registered, so it keeps running the old definition until you register the new
one. To register the new definition:

1. Run `docker compose down`.
1. Run `docker compose up` again.

A manual job isn't registered between runs, so an edit applies to the next
`docker compose run`.

## Remove jobs and failed containers

`docker compose down` removes the scheduled jobs of the project from Docker
Engine, along with the containers and networks it already removes. The
schedules stop firing, and Docker Engine keeps the run history of the removed
jobs.

`down` leaves the containers of failed runs. To remove them as well, run
`docker compose down --remove-orphans`.

If you remove a scheduled job from the Compose file, Docker Engine keeps the
registered job and it keeps firing. When a project declares scheduled jobs,
`docker compose up` warns about registered jobs that the Compose file no longer
declares. Run `docker compose up --remove-orphans` to remove them.

## Troubleshoot

### Find out why a job failed

Docker Engine keeps the container of a failed run and removes the container of
a successful run of a scheduled job, so you can only inspect failed runs. Their
containers carry the Compose labels of the project, so `docker compose ps -a`
lists them:

1. Run `docker compose ps -a` and find the container of the job. A failed run
   has an exit status other than `0`.
1. Run `docker compose logs <CONTAINER>` to read its output.

A successful run of a scheduled job leaves no container. Check its effects
instead, such as the dump files that the quickstart writes to `./backups`.

Compose limits the run history of each job to the five most recent runs,
rather than the Docker Engine default of 10 000.

### Find out why a schedule doesn't fire

Check the following:

- The jobs extension is turned on and the daemon was restarted afterward.
- You ran `docker compose up` or `docker compose create` so that Compose
  registered the job.
- The cron expression is in the time zone you expect. Without `timezone`,
  Docker Engine uses UTC.
- With the default `concurrency: forbid`, a schedule that fires while the
  previous run is still in progress skips the new run.

## Reference information

- [`jobs`](/reference/compose-file/jobs.md)
- [Manage jobs in Docker Engine](/manuals/engine/daemon/jobs.md)