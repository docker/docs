---
description: Learn how to use the local logging driver with Docker Engine
keywords: local, docker, logging, driver, file
title: Local file logging driver
aliases:
  - /engine/reference/logging/local/
  - /engine/admin/logging/local/
  - /config/containers/logging/local/
---

The `local` logging driver captures output from container's stdout/stderr and
writes them to an internal storage that's optimized for performance and disk
use.

By default, the `local` driver preserves 100MB of log messages per container and
uses automatic compression to reduce the size on disk. The 100MB default value is based on a 20M default size
for each file and a default count of 5 for the number of such files (to account for log rotation).

> [!WARNING]
>
> The `local` logging driver uses file-based storage. These files are designed
> to be exclusively accessed by the Docker daemon. Interacting with these files
> with external tools may interfere with Docker's logging system and result in
> unexpected behavior, and should be avoided.

## Usage

To use the `local` driver as the default logging driver, set the `log-driver`
and `log-opt` keys to appropriate values in the `daemon.json` file. For more
about configuring Docker using `daemon.json`, see
[daemon.json](/reference/cli/dockerd.md#daemon-configuration-file).

{{% include "daemon-cfg-desktop.md" %}}

The following example sets the log driver to `local` and sets the `max-size`
option.

```json
{
  "log-driver": "local",
  "log-opts": {
    "max-size": "10m"
  }
}
```

Restart Docker for the changes to take effect for newly created containers.
Existing containers don't use the new logging configuration automatically.

You can set the logging driver for a specific container by using the
`--log-driver` flag to `docker container create` or `docker run`:

```console
$ docker run \
      --log-driver local --log-opt max-size=10m \
      alpine echo hello world
```

Note that `local` is a bash reserved keyword, so you may need to quote it in scripts.

### Options

The `local` logging driver supports the following logging options:

| Option         | Description                                                                                                                                                                    | Example value                                      |
| :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------- |
| `max-size`     | The maximum size of the log before it's rolled. A positive integer plus a unit modifier (`k`, `m`, or `g`; case-insensitive). Defaults to 20m.                                | `--log-opt max-size=10m`                           |
| `max-file`     | The maximum number of log files that can be present. If rolling the logs creates excess files, the oldest file is removed. A positive integer. Defaults to 5.                | `--log-opt max-file=3`                             |
| `compress`     | Toggle compression of rotated log files. Enabled by default.                                                                                                                   | `--log-opt compress=false`                         |
| `labels`       | A comma-separated list of logging-related labels this daemon accepts. Used for advanced [log tag options](log_tags.md).                                                        | `--log-opt labels=production_status,geo`           |
| `labels-regex` | Similar to and compatible with `labels`. A regular expression to match logging-related labels. Used for advanced [log tag options](log_tags.md).                              | `--log-opt labels-regex=^(production_status\|geo)` |
| `env`          | A comma-separated list of logging-related environment variables this daemon accepts. Used for advanced [log tag options](log_tags.md).                                        | `--log-opt env=os,customer`                        |
| `env-regex`    | Similar to and compatible with `env`. A regular expression to match logging-related environment variables. Used for advanced [log tag options](log_tags.md).                 | `--log-opt env-regex=^(os\|customer)`              |
| `tag`          | Specify a template to set the `tag` value in log messages. Refer to the [log tag option documentation](log_tags.md) to customize the log tag format.                          | `--log-opt tag={{.ImageName}}`                     |

### Examples

This example starts an `alpine` container which can have a maximum of 3 log
files no larger than 10 megabytes each.

```console
$ docker run -it --log-driver local --log-opt max-size=10m --log-opt max-file=3 alpine ash
```
