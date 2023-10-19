---
description: Frequently asked questions for Docker Compose
keywords: documentation, docs,  docker, compose, faq, docker compose vs docker-compose
title: Compose FAQs
---

{{< include "compose-eol.md" >}}

### How do I get help?

Docker Compose is under active development. If you need help, would like to
contribute, or simply want to talk about the project with like-minded
individuals, we have a number of open channels for communication.

* To report bugs or file feature requests, use the [issue tracker on Github](https://github.com/docker/compose/issues).

* To talk about the project with people in real time, join the
  `#docker-compose` channel on the [Docker Community Slack](https://dockr.ly/slack).

* To contribute code submit a [pull request on Github](https://github.com/docker/compose/pulls).

### Where can I find example Compose files?

There are [many examples of Compose files on GitHub](https://github.com/docker/awesome-compose).

### What is the difference between `docker compose` and `docker-compose`

Version one of the Docker Compose command-line binary was first released in 2014. It was written in Python, and is invoked with `docker-compose`. Typically, Compose V1 projects include a top-level version element in the compose.yml file, with values ranging from 2.0 to 3.8, which refer to the specific file formats.

Version two of the Docker Compose command-line binary was announced in 2020, is written in Go, and is invoked with `docker compose`. Compose V2 ignores the version top-level element in the compose.yml file.

For further information, see [History and development of Compose](history.md).

### What's the difference between `up`, `run`, and `start`?

Typically, you want `docker compose up`. Use `up` to start or restart all the
services defined in a `compose.yml`. In the default "attached"
mode, you see all the logs from all the containers. In "detached" mode (`-d`),
Compose exits after starting the containers, but the containers continue to run
in the background.

The `docker compose run` command is for running "one-off" or "adhoc" tasks. It
requires the service name you want to run and only starts containers for services
that the running service depends on. Use `run` to run tests or perform
an administrative task such as removing or adding data to a data volume
container. The `run` command acts like `docker run -ti` in that it opens an
interactive terminal to the container and returns an exit status matching the
exit status of the process in the container.

The `docker compose start` command is useful only to restart containers
that were previously created but were stopped. It never creates new
containers.

### Why do my services take 10 seconds to recreate or stop?

The `docker compose stop` command attempts to stop a container by sending a `SIGTERM`. It then waits
for a [default timeout of 10 seconds](../engine/reference/commandline/compose_stop.md). After the timeout,
a `SIGKILL` is sent to the container to forcefully kill it.  If you
are waiting for this timeout, it means that your containers aren't shutting down
when they receive the `SIGTERM` signal.

There has already been a lot written about this problem of
[processes handling signals](https://medium.com/@gchudnov/trapping-signals-in-docker-containers-7a57fdda7d86)
in containers.

To fix this problem, try the following:

* Make sure you're using the exec form of `CMD` and `ENTRYPOINT`
in your Dockerfile.

  For example use `["program", "arg1", "arg2"]` not `"program arg1 arg2"`.
  Using the string form causes Docker to run your process using `bash` which
  doesn't handle signals properly. Compose always uses the JSON form, so don't
  worry if you override the command or entrypoint in your Compose file.

* If you are able, modify the application that you're running to
add an explicit signal handler for `SIGTERM`.

* Set the `stop_signal` to a signal which the application knows how to handle:

```yaml
services:
  web:
    build: .
    stop_signal: SIGINT
```

* If you can't modify the application, wrap the application in a lightweight init
system (like [s6](https://skarnet.org/software/s6/)) or a signal proxy (like
[dumb-init](https://github.com/Yelp/dumb-init) or
[tini](https://github.com/krallin/tini)).  Either of these wrappers takes care of
handling `SIGTERM` properly.

### Can I control service startup order?

Yes, see [Controlling startup order](startup-order.md).

### How do I run multiple copies of a Compose file on the same host?

Compose uses the project name to create unique identifiers for all of a
project's  containers and other resources. To run multiple copies of a project,
set a custom project name using the [`-p` command line option](reference/index.md)
or the [`COMPOSE_PROJECT_NAME` environment variable](environment-variables/envvars.md#compose_project_name).

### Can I use JSON instead of YAML for my Compose file?

Yes. [YAML is a superset of JSON](https://stackoverflow.com/a/1729545/444646) so
any JSON file should be valid YAML. To use a JSON file with Compose,
specify the filename to use, for example:

```console
$ docker compose -f docker-compose.json up
```

### Should I include my code with `COPY`/`ADD` or a volume?

You can add your code to the image using `COPY` or `ADD` directive in a
`Dockerfile`.  This is useful if you need to relocate your code along with the
Docker image, for example when you're sending code to another environment
(production, CI, etc).

You should use a `volume` if you want to make changes to your code and see them
reflected immediately, for example when you're developing code and your server
supports hot code reloading or live-reload.

There may be cases where you want to use both. You can have the image
include the code using a `COPY`, and use a `volume` in your Compose file to
include the code from the host during development. The volume overrides
the directory contents of the image.