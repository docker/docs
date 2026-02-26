---
description: Follow this hands-on tutorial to learn how to use Docker Compose from defining application
  dependencies to experimenting with commands.
keywords: docker compose example, docker compose tutorial, how to use docker compose,
  running docker compose, how to run docker compose, docker compose build image, docker
  compose command example, run docker compose file, how to create a docker compose
  file, run a docker compose file
title: Docker Compose Quickstart
linkTitle: Quickstart
weight: 30
aliases:
- /compose/samples-for-compose/
- /compose/support-and-feedback/samples-for-compose/
---

This tutorial aims to introduce fundamental concepts of Docker Compose by guiding you through the development of a basic Python web application. 

Using the Flask framework, the application features a hit counter in Redis, providing a practical example of how Docker Compose can be applied in web development scenarios. 

The concepts demonstrated here should be understandable even if you're not familiar with Python. 

This is a non-normative example that demonstrates core Compose functionality. 

## Prerequisites

Make sure you have:

- [Installed the latest version of Docker Compose](/manuals/compose/install/_index.md)
- A basic understanding of Docker concepts and how Docker works

## Step 1: Set up the project

1. Create a directory for the project:

   ```console
   $ mkdir compose-demo
   $ cd compose-demo
   ```

2. Create `app.py` in your project directory and paste the following code in:

   ```python
   import os
   import redis
   from flask import Flask

   app = Flask(__name__)
   cache = redis.Redis(
       host=os.getenv("REDIS_HOST", "redis"),
       port=int(os.getenv("REDIS_PORT", "6379")),
   )

   @app.route("/")
   def hello():
       count = cache.incr("hits")
       return f"Hello from Docker! I have been seen {count} time(s).\n"
   ```

   The app reads its Redis connection details from environment variables, with sensible defaults so it works out of the box.

3. Create `requirements.txt` in your project directory and paste the following code in:

   ```text
   flask
   redis
   ```

4. Create a `Dockerfile`:

   ```dockerfile
   # syntax=docker/dockerfile:1
   FROM python:3.12-alpine
   WORKDIR /code
   ENV FLASK_APP=app.py
   ENV FLASK_RUN_HOST=0.0.0.0
   RUN apk add --no-cache gcc musl-dev linux-headers
   COPY requirements.txt .
   RUN pip install -r requirements.txt
   COPY . .
   EXPOSE 5000
   CMD ["flask", "run", "--debug"]
   ```

   {{< accordion title="Understand the Dockerfile" >}}

   This tells Docker to:

   * Build an image starting with the Python 3.10 image.
   * Set the working directory to `/code`.
   * Set environment variables used by the `flask` command.
   * Install gcc and other dependencies
   * Copy `requirements.txt` and install the Python dependencies.
   * Add metadata to the image to describe that the container is listening on port 5000
   * Copy the current directory `.` in the project to the workdir `.` in the image.
   * Set the default command for the container to `flask run --debug`.

   {{< /accordion >}}

   > [!IMPORTANT]
   >
   > Make sure the file is named `Dockerfile` with no extension. Some editors add `.txt`
   > automatically, which causes the build to fail.

   For more information on how to write Dockerfiles, see the [Dockerfile reference](/reference/dockerfile/).

5. Create a `.env` file to hold configuration values:

   ```text
   APP_PORT=8000
   REDIS_HOST=redis
   REDIS_PORT=6379
   ```

   Compose automatically reads `.env` and makes these values available for interpolation
   in your `compose.yaml`. For this example the gains are modest, but in practice,
   keeping configuration out of the Compose file makes it easier to change values across
   environments without editing YAML, avoid committing secrets to
   version control, and reuse values across multiple services.

6. Create a `.dockerignore` file to keep unnecessary files out of your build context:

   ```text
   .env
   *.pyc
   __pycache__
   redis-data
   ```

   Docker sends everything in your project directory to the daemon when it builds an image.
   Without `.dockerignore`, that includes your `.env` file (which may contain secrets) and
   any cached Python bytecode. Excluding them keeps builds fast and avoids accidentally
   baking sensitive values into an image layer.

## Step 2: Define and start your services

Compose simplifies the control of your entire application stack, making it easy to manage services, networks, and volumes in a single, comprehensible YAML configuration file. 

1. Create `compose.yaml` in your project directory and paste the following:

   ```yaml
   services:
   web:
      build: .
      ports:
         - "${APP_PORT}:5000"
      environment:
         - REDIS_HOST=${REDIS_HOST}
         - REDIS_PORT=${REDIS_PORT}

   redis:
      image: redis:alpine
   ```

   This Compose file defines two services: `web` and `redis`. 

   The `web` service uses an image that's built from the `Dockerfile` in the current directory.
   It then binds the container and the host machine to the exposed port, `8000`. This example service uses the default port for the Flask web server, `5000`.

   The `redis` service uses a public [Redis](https://registry.hub.docker.com/_/redis/) image pulled from the Docker Hub registry.

   For more information on the `compose.yaml` file, see [How Compose works](compose-application-model.md).

2. From your project directory, start up your application: 

   ```console
   $ docker compose up
   ```

   With a single command, you create and start all the services from your configuration file. Compose builds your web image, pulls the Redis image, and starts both containers. 

3. Open `http://localhost:8000`. You should see:

   ```text
   Hello from Docker! I have been seen 1 time(s).
   ```

   Refresh the page — the counter increments on each visit.

   This minimal setup works, but it has two problems you'll fix in the next steps:

   - Startup race: `web` starts at the same time as `redis`. If Redis isn't ready yet,
   the Flask app fails to connect and crashes.
   - No persistence: If you run `docker compose down` followed by `docker compose up`, the
   counter resets to zero. `docker compose down` removes the containers, and with them
   any data written to the container's writable layer. `docker compose stop` preserves
   the containers so data survives, but you can't rely on that in production where
   containers are regularly replaced.

4.    Stop the stack before moving on:

   ```console
   $ docker compose down
   ```

## Step 3: Fix the startup race with health checks

To fix the startup race, Compose needs to wait until `redis` is confirmed healthy before
starting `web`.

1. Update `compose.yaml`:

   ```yaml
   services:
   web:
      build: .
      ports:
         - "${APP_PORT}:5000"
      environment:
         - REDIS_HOST=${REDIS_HOST}
         - REDIS_PORT=${REDIS_PORT}
      depends_on:
         redis:
         condition: service_healthy

   redis:
      image: redis:alpine
      healthcheck:
         test: ["CMD", "redis-cli", "ping"]
         interval: 5s
         timeout: 3s
         retries: 5
         start_period: 10s
   ```

   The `healthcheck` block tells Compose how to test whether Redis is ready:

   - `test` is the command Compose runs inside the container to check its health.
     `redis-cli ping` connects to Redis and expects a `PONG` response — if it gets one,
     the container is healthy.
   - `start_period` gives Redis 10 seconds to initialise before health checks begin.
     Any failures during this window don't count toward the retry limit.
   - `interval` runs the check every 5 seconds after the start period has elapsed.
   - `timeout` gives each check 3 seconds to respond before treating it as a failure.
   - `retries` sets how many consecutive failures are allowed before Compose marks the
     container as unhealthy. With `interval: 5s` and `retries: 5`, Compose will wait up
     to 25 seconds before giving up.

2. Start the stack to confirm the ordering is fixed:

   ```console
      $ docker compose up
   ```

   You should see something similar to:

   ```text
      [+] Running 2/2
      ✔ Container compose-test-redis-1  Healthy                       0.0s
   ```

3. Open `http://localhost:8000` to confirm the app is still working, then stop the stack before moving on:

   ```console
      $ docker compose down
   ```

## Step 4: Enable Compose Watch for live updates

Now that startup order is handled, add Compose Watch so that code changes sync into the
running container automatically

1. Update `compose.yaml` to add the `develop.watch` block to the `web` service:
   
   ```yaml
   services:
      web:
         build: .
         ports:
            - "${APP_PORT}:5000"
         environment:
            - REDIS_HOST=${REDIS_HOST}
            - REDIS_PORT=${REDIS_PORT}
         depends_on:
            redis:
            condition: service_healthy
         develop:
            watch:
            - action: sync+restart
               path: .
               target: /code
            - action: rebuild
               path: requirements.txt

   redis:
      image: redis:alpine
      healthcheck:
         test: ["CMD", "redis-cli", "ping"]
         interval: 5s
         timeout: 3s
         retries: 5
         start_period: 10s
   ```

   The `watch` block defines two rules. The `sync+restart` action syncs any changes
   under `.` into `/code` inside the running container, then restarts the Flask process
   so it picks up the new files — this is more reliable than depending on Flask's
   reloader to detect the change itself. The `rebuild` action on `requirements.txt`
   triggers a full image rebuild whenever you add a new dependency, since that can't be
   handled by a simple file sync.

2. Start the stack with Watch enabled:

   ```console
      $ docker compose up --watch
   ```

3. Make a live change. Open `app.py` and update the greeting:

   ```python
      return f"Hello from Compose Watch! I have been seen {count} time(s).\n"
   ```

4. Save the file. Compose Watch detects the change and syncs it immediately:

   ```text
      Syncing service "web" after changes were detected
   ```

5. Refresh `http://localhost:8000`. The updated greeting appears without any restart
   and the counter should still be incrementing.

   > [!NOTE]
   >
   > For this example to work, the `--debug` option is added to the `Dockerfile`. The
   > `--debug` option in Flask enables automatic code reload, making it possible to work
   > on the backend API without the need to restart or rebuild the container.
   > After changing the `.py` file, subsequent API calls will use the new code, but the
   > browser UI will not automatically refresh in this small example. Most frontend
   > development servers include native live reload support that works with Compose.

5. Stop the stack before moving on:

   ```console
      $ docker compose down
   ```

   For more information on how Compose Watch works, see [Use Compose Watch](/manuals/compose/how-tos/file-watch.md).

## Step 4: Persist data with named volumes

Each time you stop and restart the stack the visit counter resets to zero. Redis data
lives inside the container, so it disappears when the container is removed. A named
volume fixes this by storing the data on the host, outside the container lifecycle.

1. Update `compose.yaml`:

   ```yaml
   services:
   web:
      build: .
      ports:
         - "${APP_PORT}:5000"
      environment:
         - REDIS_HOST=${REDIS_HOST}
         - REDIS_PORT=${REDIS_PORT}
      depends_on:
         redis:
         condition: service_healthy
      develop:
         watch:
         - action: sync+restart
            path: .
            target: /code
         - action: rebuild
            path: requirements.txt

   redis:
      image: redis:alpine
      volumes:
         - redis-data:/data
      healthcheck:
         test: ["CMD", "redis-cli", "ping"]
         interval: 5s
         timeout: 3s
         retries: 5
         start_period: 10s

   volumes:
   redis-data:
   ```

   The `redis-data:/data` entry under `redis.volumes` mounts the named volume at `/data`, the path where Redis
   writes its data files. The top-level `volumes` key registers it with Docker so it
   persists between `compose down` and `compose up` cycles.

2. Start the stack with `docker compose up --watch` and refresh `http://localhost:8000` a few times to build up a count.

3. Tear down the stack with `docker compose down` and then bring it back up again with `docker compose up --watch`.

4. Open `http://localhost:8000` — the counter continues from where it left off.

5. Now reset the counter with `docker compose down -v`. 

   The `-v` flag removes named volumes along with the containers. Use this intentionally — it permanently deletes the stored data.

## Step 5: Structure your project with multiple Compose files

As applications grow, a single `compose.yaml` becomes harder to maintain. The `include`
top-level elements lets you split services across multiple files while keeping them part of the
same application.

This is especially useful when different teams own different parts of the stack, or when
you want to reuse infrastructure definitions across projects.

1. Create a new file in your project directory called `infra.yaml` and move the Redis service and volume into it:

   ```yaml
   services:
     redis:
       image: redis:alpine
       volumes:
         - redis-data:/data
       healthcheck:
         test: ["CMD", "redis-cli", "ping"]
         interval: 5s
         timeout: 3s
         retries: 5

   volumes:
     redis-data:
   ```

2. Update `compose.yaml` to include `infra.yaml`:

   ```yaml
   include:
    - path: ./infra.yaml

   services:
     web:
       build: .
       ports:
         - "${APP_PORT}:5000"
       environment:
         - REDIS_HOST=${REDIS_HOST}
         - REDIS_PORT=${REDIS_PORT}
       depends_on:
         redis:
           condition: service_healthy
       develop:
         watch:
            - action: sync+restart
               path: .
               target: /code
            - action: rebuild
               path: requirements.txt
   ```

3. Run the application to confirm everything still works:

   ```console
   $ docker compose up --watch
   ```

   Compose merges both files at startup. The `web` service can still reference `redis`
   by name because all included services share the same default network.

   This is a simplified example, but it demonstrates the basic principle of `include` and how it can make it easier to modularize complex applications into sub-Compose files. For more information on `include` and working with multiple Compose files, see [Working with multiple Compose files](/manuals/compose/how-tos/multiple-compose-files/_index.md). For more advanced patterns — including environment-specific overrides and service inheritance — see [Use multiple Compose files](/compose/how-tos/multiple-compose-files/).

4. Stop the stack before moving on:

   ```console
      $ docker compose down
   ```

## Step 6: Inspect and debug your running stack

With a fully configured stack, you can observe what's happening inside your containers
without stopping anything. This step covers the core commands for inspecting the resolved configuration, streaming logs, and running commands
inside a running container.

Before starting the stack, verify that Compose has resolved your `.env` variables and
merged all files correctly:

```console
$ docker compose config
```

`docker compose config` doesn't require the stack to be running — it works purely from
your files. A few things worth noting in the output:

- `${APP_PORT}`, `${REDIS_HOST}`, and `${REDIS_PORT}` have all been replaced with the
  values from your `.env` file.
- Short-form port notation (`"8000:5000"`) is expanded into its canonical fields
  (`target`, `published`, `protocol`).
- The default network and volume names are made explicit, prefixed with the project name
  `compose-test`.
- prints the fully resolved configuration, merging any files
  brought in via `include` into a single view.

Use `docker compose config` any time you want to confirm what Compose will actually
apply, especially when debugging variable substitution or working with muliple Compose files.

Now start the stack in detached mode so the terminal stays free for the commands that
follow:

```console
$ docker compose up --watch -d
```
### Stream logs from all services

```console
$ docker compose logs -f
```

The `-f` flag follows the log stream in real time, interleaving output from both
containers with color-coded service name prefixes. Refresh `http://localhost:8000` a
few times and watch the Flask request logs appear. To follow logs for a single service,
pass its name:

```console
$ docker compose logs -f web
```

Press `Ctrl+C` to stop following logs. The containers keep running.

### Run commands inside a running container

`docker compose exec` runs a command inside an already-running container without
starting a new one. This is the primary tool for live debugging.

#### Verify environment variables are set correctly

```console
$ docker compose exec web env | grep REDIS
```

```text
REDIS_HOST=redis
REDIS_PORT=6379
```

#### Test that the `web` container can reach Redis using the service name as the hostname

```console
$ docker compose exec web python -c "import redis; r = redis.Redis(host='redis'); print(r.ping())"
```

```text
True
```

This uses the same `redis` library your app uses, so a `True` response confirms that
service discovery, networking, and the Redis connection are all working end to end.

#### Inspect the live value of the hit counter in Redis

```console
$ docker compose exec redis redis-cli GET hits
```

## Where to go next

- [Explore the full list of Compose commands](/reference/cli/docker/compose/)
- [Explore the Compose file reference](/reference/compose-file/_index.md)
- [Check out the Learning Docker Compose video on LinkedIn Learning](https://www.linkedin.com/learning/learning-docker-compose/)
- [Set environment variables in Compose](/compose/how-tos/environment-variables/set-environment-variables/) — go deeper on `.env`, interpolation, and precedence
- [OCI artifact applications](/compose/how-tos/oci-artifact/) — package and distribute your Compose app from a registry


