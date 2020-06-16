---
description: Overview of docker-compose CLI
keywords: fig, composition, compose, docker, orchestration, cli,  docker-compose
redirect_from:
- /compose/reference/docker-compose/
title: Overview of docker-compose CLI
---

This page provides the usage information for the `docker-compose` Command.

## Command options overview and help

You can also see this information by running `docker-compose --help` from the
command line.

```none
Define and run multi-container applications with Docker.

Usage:
  docker-compose [-f <arg>...] [options] [COMMAND] [ARGS...]
  docker-compose -h|--help

Options:
  -f, --file FILE             Specify an alternate compose file
                              (default: docker-compose.yml)
  -p, --project-name NAME     Specify an alternate project name
                              (default: directory name)
  --verbose                   Show more output
  --log-level LEVEL           Set log level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
  --no-ansi                   Do not print ANSI control characters
  -v, --version               Print version and exit
  -H, --host HOST             Daemon socket to connect to

  --tls                       Use TLS; implied by --tlsverify
  --tlscacert CA_PATH         Trust certs signed only by this CA
  --tlscert CLIENT_CERT_PATH  Path to TLS certificate file
  --tlskey TLS_KEY_PATH       Path to TLS key file
  --tlsverify                 Use TLS and verify the remote
  --skip-hostname-check       Don't check the daemon's hostname against the
                              name specified in the client certificate
  --project-directory PATH    Specify an alternate working directory
                              (default: the path of the Compose file)
  --compatibility             If set, Compose will attempt to convert deploy
                              keys in v3 files to their non-Swarm equivalent

Commands:
  build              Build or rebuild services
  bundle             Generate a Docker bundle from the Compose file
  config             Validate and view the Compose file
  create             Create services
  down               Stop and remove containers, networks, images, and volumes
  events             Receive real time events from containers
  exec               Execute a command in a running container
  help               Get help on a command
  images             List images
  kill               Kill containers
  logs               View output from containers
  pause              Pause services
  port               Print the public port for a port binding
  ps                 List containers
  pull               Pull service images
  push               Push service images
  restart            Restart services
  rm                 Remove stopped containers
  run                Run a one-off command
  scale              Set number of containers for a service
  start              Start services
  stop               Stop services
  top                Display the running processes
  unpause            Unpause services
  up                 Create and start containers
  version            Show the Docker-Compose version information
```

You can use Docker Compose binary, `docker-compose [-f <arg>...] [options]
[COMMAND] [ARGS...]`, to build and manage multiple services in Docker containers.

## Use `-f` to specify name and path of one or more Compose files

Use the `-f` flag to specify the location of a Compose configuration file.

### Specifying multiple Compose files

You can supply multiple `-f` configuration files. When you supply multiple
files, Compose combines them into a single configuration. Compose builds the
configuration in the order you supply the files. Subsequent files override and
add to their predecessors.

For example, consider this command line:

```
$ docker-compose -f docker-compose.yml -f docker-compose.admin.yml run backup_db
```

The `docker-compose.yml` file might specify a `webapp` service.

```
webapp:
  image: examples/web
  ports:
    - "8000:8000"
  volumes:
    - "/data"
```

If the `docker-compose.admin.yml` also specifies this same service, any matching
fields override the previous file. New values, add to the `webapp` service
configuration.

```
webapp:
  build: .
  environment:
    - DEBUG=1
```

Use a `-f` with `-` (dash) as the filename to read the configuration from
`stdin`. When `stdin` is used all paths in the configuration are
relative to the current working directory.

The `-f` flag is optional. If you don't provide this flag on the command line,
Compose traverses the working directory and its parent directories looking for a
`docker-compose.yml` and a `docker-compose.override.yml` file. You must supply
at least the `docker-compose.yml` file. If both files are present on the same
directory level, Compose combines the two files into a single configuration.

The configuration in the `docker-compose.override.yml` file is applied over and
in addition to the values in the `docker-compose.yml` file.

### Specifying a path to a single Compose file

You can use the `-f` flag to specify a path to a Compose file that is not
located in the current directory, either from the command line or by setting up
a [COMPOSE_FILE environment variable](envvars.md#compose_file) in your shell or
in an environment file.

For an example of using the `-f` option at the command line, suppose you are
running the [Compose Rails sample](../rails.md), and
have a `docker-compose.yml` file in a directory called `sandbox/rails`. You can
use a command like [docker-compose pull](pull.md) to get the
postgres image for the `db` service from anywhere by using the `-f` flag as
follows: `docker-compose -f ~/sandbox/rails/docker-compose.yml pull db`

Here's the full example:

```
$ docker-compose -f ~/sandbox/rails/docker-compose.yml pull db
Pulling db (postgres:latest)...
latest: Pulling from library/postgres
ef0380f84d05: Pull complete
50cf91dc1db8: Pull complete
d3add4cd115c: Pull complete
467830d8a616: Pull complete
089b9db7dc57: Pull complete
6fba0a36935c: Pull complete
81ef0e73c953: Pull complete
338a6c4894dc: Pull complete
15853f32f67c: Pull complete
044c83d92898: Pull complete
17301519f133: Pull complete
dcca70822752: Pull complete
cecf11b8ccf3: Pull complete
Digest: sha256:1364924c753d5ff7e2260cd34dc4ba05ebd40ee8193391220be0f9901d4e1651
Status: Downloaded newer image for postgres:latest
```

## Use `-p` to specify a project name

Each configuration has a project name. If you supply a `-p` flag, you can
specify a project name. If you don't specify the flag, Compose uses the current
directory name. See also the [COMPOSE_PROJECT_NAME environment variable](envvars.md#compose_project_name).

## Set up environment variables

You can set [environment variables](envvars.md) for various
`docker-compose` options, including the `-f` and `-p` flags.

For example, the [COMPOSE_FILE environment variable](envvars.md#compose_file)
relates to the `-f` flag, and `COMPOSE_PROJECT_NAME`
[environment variable](envvars.md#compose_project_name) relates to the `-p` flag.

Also, you can set some of these variables in an [environment file](../env-file.md).

## Where to go next

* [CLI environment variables](envvars.md)
* [Declare default environment variables in file](../env-file.md)
