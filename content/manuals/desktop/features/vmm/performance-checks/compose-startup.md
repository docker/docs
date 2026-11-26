---
title: Time a multi-service Compose startup
linkTitle: Compose startup
description: Measure how long it takes to restart a stopped 11-service Supabase Compose stack until it serves a known row over HTTP.
keywords: docker desktop, performance, benchmark, docker compose, supabase, postgres, startup
weight: 70
---

This check measures the wait when you return to a project and start its whole
local environment again. The containers are gone, but the database on disk
remains. It uses Supabase's self-hosted Compose setup, pinned to one commit,
with 11 services including PostgreSQL and a gateway.

The slow preparation happens before timing: downloading images, initializing the
database, writing one known row, and stopping the stack while keeping its data.
The timed part brings the stack back up and ends when an HTTP request through
the gateway returns that row. A healthy stack isn't enough. The timer runs until
the data comes back.

The result is wall-clock time on your machine. Repeat the full preparation before
every run. A second start of the same stack measures warm caches.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

## Additional prerequisites

- Docker Compose v2.24.4 or later.
- `curl` is installed.
- About 12 GB of free disk space for images.
- Port `18000` on localhost is free. The check uses its own container names and
  Compose project name, so it doesn't affect an existing Supabase stack. Stop
  anything that uses those names or that port.

## Run the check

1. Get the pinned Compose setup:

   ```console
   $ mkdir -p ~/perfcheck/compose-startup && cd ~/perfcheck/compose-startup

   $ rm -rf fixture && mkdir fixture
   $ git init -q fixture/repo
   $ git -C fixture/repo remote add origin https://github.com/supabase/supabase.git
   $ git -C fixture/repo fetch -q --depth=1 origin 86854671e95be31e24fa0785cded0579fe692fcb
   $ git -C fixture/repo checkout -q --detach FETCH_HEAD
   ```

2. Add an override file. It renames the containers, moves the gateway to port
   18000, stops one service from publishing ports, and puts the database on a
   named volume. The `!override` marker needs Compose 2.24.4 or later.

   ```console
   $ cat > fixture/repo/docker/docker-compose.perfcheck.yml <<'EOF'
   services:
     studio:
       container_name: perfcheck-supabase-studio
     kong:
       container_name: perfcheck-supabase-kong
       ports: !override
         - "127.0.0.1:18000:8000/tcp"
     auth:
       container_name: perfcheck-supabase-auth
     rest:
       container_name: perfcheck-supabase-rest
     realtime:
       container_name: perfcheck-supabase-realtime
     storage:
       container_name: perfcheck-supabase-storage
     imgproxy:
       container_name: perfcheck-supabase-imgproxy
     meta:
       container_name: perfcheck-supabase-meta
     functions:
       container_name: perfcheck-supabase-functions
     db:
       container_name: perfcheck-supabase-db
       volumes:
         - perfcheck-db-data:/var/lib/postgresql/data
     supavisor:
       container_name: perfcheck-supabase-supavisor
       ports: !override []
   volumes:
     perfcheck-db-data:
   EOF
   ```

3. Write the environment file. The keys are Supabase's published demo API keys.
   They expire far in the future, so the result doesn't depend on when you run
   the check.

   ```console
   $ cd fixture/repo/docker

   $ ANON=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlLWRlbW8iLCJpYXQiOjE2NDE3NjkyMDAsImV4cCI6NDEwMjQ0NDgwMH0.ls8OqwV1EwDuvR9mnGLHB6UF60bOS8KI7Eb9oPaGl9M
   $ SERVICE=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmljZV9yb2xlIiwiaXNzIjoic3VwYWJhc2UtZGVtbyIsImlhdCI6MTY0MTc2OTIwMCwiZXhwIjo0MTAyNDQ0ODAwfQ.MoxRoj5EjPmaLLdlpy5y7JhpP5IwLMDhHodKE58usYk

   $ sed -e "s|^ANON_KEY=.*|ANON_KEY=$ANON|" \
         -e "s|^SERVICE_ROLE_KEY=.*|SERVICE_ROLE_KEY=$SERVICE|" \
         .env.example > .env
   ```

4. Define a helper for the Compose commands:

   ```bash
   compose() {
     docker compose --project-name perfcheck-supabase --env-file .env \
       --file docker-compose.yml --file docker-compose.perfcheck.yml "$@"
   }
   ```

5. Remove anything left from a previous attempt, download the images, and start
   the stack once so the database initializes. This step isn't timed. It's the
   slowest step, and downloads about 12 GB the first time.

   ```console
   $ compose down --volumes --remove-orphans
   $ compose pull
   $ compose up -d --wait --no-build
   ```

6. Write one known row into the database:

   ```console
   $ compose exec -T db psql -v ON_ERROR_STOP=1 -U postgres -d postgres \
     -c "CREATE TABLE public.perfcheck (value text PRIMARY KEY); INSERT INTO public.perfcheck (value) VALUES ('ready'); GRANT SELECT ON public.perfcheck TO service_role; NOTIFY pgrst, 'reload schema';"
   ```

7. Confirm that you can read the row over HTTP. This step isn't timed. The first
   attempts can return a 404 while the API layer detects the new table.

   ```console
   $ attempt=0
   $ until curl -fsS --connect-timeout 1 --max-time 2 \
       -H "apikey: $SERVICE" -H "Authorization: Bearer $SERVICE" \
       'http://127.0.0.1:18000/rest/v1/perfcheck?select=value' > response.json &&
       [ "$(tr -d '\r\n' < response.json)" = '[{"value":"ready"}]' ]; do
     attempt=$((attempt + 1)); [ "$attempt" -lt 30 ] ||
       { echo "PREPARATION FAILED - clean up and start this run again"; break; }
     sleep 1
   done
   ```

8. Stop the stack and keep its data:

   ```console
   $ compose down --remove-orphans
   ```

9. Run the timed command. It starts every service and waits until the row comes
   back.

   ```console
   $ time (
     compose up -d --wait --no-build &&
     attempt=0 &&
     until curl -fsS --connect-timeout 1 --max-time 2 \
         -H "apikey: $SERVICE" -H "Authorization: Bearer $SERVICE" \
         'http://127.0.0.1:18000/rest/v1/perfcheck?select=value' > response.json &&
         [ "$(tr -d '\r\n' < response.json)" = '[{"value":"ready"}]' ]; do
       attempt=$((attempt + 1)); [ "$attempt" -lt 30 ] || exit 1; sleep 1
     done
   )
   ```

   Record the elapsed time.

10. Check that the run is valid. The timed block must finish without an error,
   and `response.json` must contain exactly:

   ```text
   [{"value":"ready"}]
   ```

11. Clean up. Do this even if the run failed. A stack left running affects every
    later measurement.

    ```console
    $ compose down --volumes --remove-orphans
    $ cd ../../..
    $ rm -rf fixture
    ```

    On Linux, the database writes root-owned files to
    `fixture/repo/docker/volumes`. If `rm -rf` can't remove them, change their
    permissions first. Run this from `fixture/repo/docker`, after the containers
    are gone:

    ```console
    $ docker run --rm --user 0:0 --entrypoint sh \
      --volume "$PWD/volumes:/cleanup" supabase/postgres:17.6.1.136 \
      -c 'chmod -R a+rwX /cleanup'
    ```
