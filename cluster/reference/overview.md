---
description: Overview of docker cluster CLI
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: Overview of docker cluster CLI
---

This page provides usage information for the `docker cluster` CLI plugin command options.

You can also view this information by running `docker cluster --help` from the
command line.

## Usage
```
docker cluster [Options] [Commands]
```

Options:

- `--dry-run`: Skips resource provisioning.
- `--log-level string`: Specifies the logging level. Valid values include: `trace`,`debug`,`info`,`warn`,`error`, and `fatal`. Defaults to `warn`.

Commands:

- `backup`: Backs up a running cluster.
- `begin`: Creates an example cluster declaration.
- `create`: Creates a new Docker cluster.
- `inspect`: Provides detailed information about a cluster.
- `logs`:TODO: Fetches cluster logs.
- `ls`: Lists all available clusters.
- `restore`: Restores a cluster from a backup.
- `rm`: Removes a cluster.
- `update`: Updates a running cluster's desired state.
- `version`: Displays Version, Commit, and Build type.

Run 'docker cluster [Command] --help' for more information about a command.
```

## Specify name and path of one or more cluster files

Use the `-f` flag to specify the location of a cluster configuration file.

## Set up environment variables

You can set [environment variables](envvars) for various
`docker cluster` options, including the `-f` and `-p` flags.

## Where to go next

* [CLI environment variables](envvars)
