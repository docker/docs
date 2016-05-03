<!--[metadata]>
+++
title = "DTR tool reference"
description = "Learn about the options available on the docker/trusted-registry image."
keywords = ["docker, dtr, install, uninstall, configure"]
[menu.main]
parent="workw_dtr_install"
identifier="dtr_menu_reference"
weight=60
+++
<![end-metadata]-->

# DTR tool reference

This tool has commands to install, configure, and backup Docker
Trusted Registry (DTR). It also allows uninstalling DTR.
By default the tool runs in interactive mode. It prompts you for
the values needed.
For running this tool in non-interactive mode, there are three
ways you can use to pass values:

```bash
$ docker run -it --rm docker/dtr command --option value
$ docker run -e --rm docker/dtr command ENV_VARIABLE=value
$ docker run -e --rm docker/dtr command ENV_VARIABLE
```

Additional help is available for each command with the '--help' option.

## Usage

```bash
$ docker run -it --rm docker/dtr \
    command [command options]
```

## Options

| Option      | Description |
|:------------|:------------|
| `--help, h` | Show help   |

## Commands

| Option        | Description                                                                     |
|:--------------|:--------------------------------------------------------------------------------|
| `install`     | Install Docker Trusted Registry on this Docker Engine                           |
| `join`        | Add a new replica to an existing DTR cluster                                    |
| `reconfigure` | Change DTR configurations                                                       |
| `remove`      | Remove a replica from a DTR cluster                                             |
| `restore`     | Create a new DTR cluster from an existing backup                                |
| `backup`      | Backup a DTR cluster to a tar file and stream it to stdout                      |
| `migrate`     | Migrate configurations, accounts, and repository metadata from DTR 1.4.3 to 2.0 |
| `dumpcerts`   | Dump out the TLS certificates used by this DTR instance                         |
