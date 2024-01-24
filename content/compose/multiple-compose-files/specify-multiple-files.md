---
title: Specify multiple compose files
description: Understand how to specify multiple Compose files in the command line
keywords: -f flag, compose file, multiple compose files, command line
---

Use the `-f` flag to specify the location of multiple Compose configuration files.

When you supply multiple files, Compose combines them into a single configuration. Compose builds the
configuration in the order you supply the files. Subsequent files override and
add to their predecessors.

For example, consider this command line:

```console
$ docker compose -f compose.yml -f compose.admin.yml run backup_db
```

The `compose.yml` file might specify a `webapp` service.

```yaml
webapp:
  image: examples/web
  ports:
    - "8000:8000"
  volumes:
    - "/data"
```

If the `compose.admin.yml` also specifies this same service, any matching
fields override the previous file. New values, add to the `webapp` service
configuration.

```yaml
webapp:
  build: .
  environment:
    - DEBUG=1
```

When you use multiple Compose files, all paths in the files are relative to the
first configuration file specified with `-f`. You can use the
`--project-directory` option to override this base path.

Use a `-f` with `-` (dash) as the filename to read the configuration from
`stdin`. When `stdin` is used all paths in the configuration are
relative to the current working directory.

The `-f` flag is optional. If you don't provide this flag on the command line,
Compose traverses the working directory and its parent directories looking for a
`compose.yml` and a `compose.override.yml` file. You must supply
at least the `compose.yml` file. If both files are present on the same
directory level, Compose combines the two files into a single configuration.

The configuration in the `compose.override.yml` file is applied over and
in addition to the values in the `compose.yml` file.

### Specifying a path to a single Compose file

You can use the `-f` flag to specify a path to a Compose file that is not
located in the current directory, either from the command line or by setting up
a [COMPOSE_FILE environment variable](../environment-variables/envvars.md#compose_file) in your shell or
in an environment file.

For an example of using the `-f` option at the command line, suppose you are
running the [Compose Rails sample](https://github.com/docker/awesome-compose/tree/master/official-documentation-samples/rails/README.md), and
have a `compose.yml` file in a directory called `sandbox/rails`. You can
use a command like [docker compose pull](../../engine/reference/commandline/compose_pull.md) to get the
postgres image for the `db` service from anywhere by using the `-f` flag as
follows: `docker compose -f ~/sandbox/rails/compose.yml pull db`

Here's the full example:

```console
$ docker compose -f ~/sandbox/rails/compose.yml pull db
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