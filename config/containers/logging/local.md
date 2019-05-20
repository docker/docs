---
<<<<<<< HEAD
description: Describes how to use the local logging driver.
keywords: local, docker, logging, driver
redirect_from:
- /engine/reference/logging/local/
- /engine/admin/logging/local/
title: Local File logging driver
---

The `local` logging driver captures output from container's stdout/stderr and
writes them to an internal storage that is optimized for performance and disk
use.

By default the `local` driver preserves 100MB of log messages per container and
uses automatic compression to reduce the size on disk.

> *Note*: the `local` logging driver currently uses file-based storage. The
> file-format and storage mechanism are designed to be exclusively accessed by
> the Docker daemon, and should not be used by external tools as the
> implementation may change in future releases.
=======
description: Describes how to use the local binary (Protobuf) logging driver.
keywords: local, protobuf, docker, logging, driver
redirect_from:
- /engine/reference/logging/local/
- /engine/admin/logging/local/
title: local binary file Protobuf logging driver
---

This `log-driver` writes to `local` binary files using Protobuf [Protocol Buffers](https://en.wikipedia.org/wiki/Protocol_Buffers)
>>>>>>> Sync forked amberjack branch with docs-private (#1068)

## Usage

To use the `local` driver as the default logging driver, set the `log-driver`
and `log-opt` keys to appropriate values in the `daemon.json` file, which is
located in `/etc/docker/` on Linux hosts or
<<<<<<< HEAD
`C:\ProgramData\docker\config\daemon.json` on Windows Server. For more about
configuring Docker using `daemon.json`, see
[daemon.json](/engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `local` and sets the `max-size`
option.
=======
`C:\ProgramData\docker\config\daemon.json` on Windows Server. For more information about
configuring Docker using `daemon.json`, see
[daemon.json](/engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `local`.
>>>>>>> Sync forked amberjack branch with docs-private (#1068)

```json
{
  "log-driver": "local",
<<<<<<< HEAD
  "log-opts": {
    "max-size": "10m"
  }
}
```

Restart Docker for the changes to take effect for newly created containers. Existing containers do not use the new logging configuration.
=======
  "log-opts": {}
}
```

> **Note**: `log-opt` configuration options in the `daemon.json` configuration
> file must be provided as strings. Boolean and numeric values (such as the value
> for `max-file` in the example above) must therefore be enclosed in quotes (`"`).

Restart Docker for the changes to take effect for newly created containers.

Existing containers will not use the new logging configuration.
>>>>>>> Sync forked amberjack branch with docs-private (#1068)

You can set the logging driver for a specific container by using the
`--log-driver` flag to `docker container create` or `docker run`:

```bash
$ docker run \
<<<<<<< HEAD
      --log-driver local --log-opt max-size=10m \
=======
      --log-driver local --log-opt compress="false" \
>>>>>>> Sync forked amberjack branch with docs-private (#1068)
      alpine echo hello world
```

### Options

The `local` logging driver supports the following logging options:

| Option      | Description                                                                                                                                                                                                   | Example  value                           |
|:------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------------|
<<<<<<< HEAD
| `max-size`  | The maximum size of the log before it is rolled. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`). Defaults to 20m.                                          | `--log-opt max-size=10m`                 |
| `max-file`  | The maximum number of log files that can be present. If rolling the logs creates excess files, the oldest file is removed. **Only effective when `max-size` is also set.** A positive integer. Defaults to 5. | `--log-opt max-file=3`                   |
| `compress`  | Toggle compression of rotated log files. Enabled by default. | `--log-opt compress=false` |

### Examples

This example starts an `alpine` container which can have a maximum of 3 log
files no larger than 10 megabytes each.

```bash
$ docker run -it --log-opt max-size=10m --log-opt max-file=3 alpine ash
```
=======
| `max-size`  | The maximum size of each binary log file before rotation. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`). Defaults to `20m`.                                          | `--log-opt max-size=10m`                 |
| `max-file`  | The maximum number of binary log files. If rotating the logs creates an excess file, the oldest file is removed. **Only effective when `max-size` is also set.** A positive integer. Defaults to `5`.         | `--log-opt max-file=5`                   |
| `compress`  | Whether or not the binary files should be compressed. Defaults to `true`                                                                                                                                      | `--log-opt compress=true`              |
>>>>>>> Sync forked amberjack branch with docs-private (#1068)
