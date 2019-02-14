---
description: Describes how to use the local binary (Protobuf) logging driver.
keywords: local, protobuf, docker, logging, driver
redirect_from:
- /engine/reference/logging/local/
- /engine/admin/logging/local/
title: local binary file Protobuf logging driver
---

This `log-driver` writes to `local` binary files using Protobuf [Protocol Buffers](https://en.wikipedia.org/wiki/Protocol_Buffers)

## Usage

To use the `local` driver as the default logging driver, set the `log-driver`
and `log-opt` keys to appropriate values in the `daemon.json` file, which is
located in `/etc/docker/` on Linux hosts or
`C:\ProgramData\docker\config\daemon.json` on Windows Server. For more information about
configuring Docker using `daemon.json`, see
[daemon.json](/engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `local`.

```json
{
  "log-driver": "local",
  "log-opts": {}
}
```

> **Note**: `log-opt` configuration options in the `daemon.json` configuration
> file must be provided as strings. Boolean and numeric values (such as the value
> for `max-file` in the example above) must therefore be enclosed in quotes (`"`).

Restart Docker for the changes to take effect for newly created containers.

Existing containers will not use the new logging configuration.

You can set the logging driver for a specific container by using the
`--log-driver` flag to `docker container create` or `docker run`:

```bash
$ docker run \
      --log-driver local --log-opt compress="false" \
      alpine echo hello world
```

### Options

The `json-file` logging driver supports the following logging options:

| Option      | Description                                                                                                                                                                                                   | Example  value                           |
|:------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------------|
| `max-size`  | The maximum size of each binary log file before rotation. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`). Defaults to `20m`.                                          | `--log-opt max-size=10m`                 |
| `max-file`  | The maximum number of binary log files. If rotating the logs creates an excess file, the oldest file is removed. **Only effective when `max-size` is also set.** A positive integer. Defaults to `5`.         | `--log-opt max-file=5`                   |
| `compress`  | Whether or not the binary files should be compressed. Defaults to `true`                                                                                                                                      | `--log-opt compress=true`              |
