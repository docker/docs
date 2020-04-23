---
description: Describes how to use the Graylog Extended Format logging driver.
keywords: graylog, gelf, logging, driver
redirect_from:
- /engine/reference/logging/gelf/
- /engine/admin/logging/gelf/
title: Graylog Extended Format logging driver
---

The `gelf` logging driver is a convenient format that is understood by a number of tools such as
[Graylog](https://www.graylog.org/), [Logstash](https://www.elastic.co/products/logstash), and
[Fluentd](http://www.fluentd.org/). Many tools use this format.

In GELF, every log message is a dict with the following fields:

- version
- host (who sent the message in the first place)
- timestamp
- short and long version of the message
- any custom fields you configure yourself

## Usage

To use the `gelf` driver as the default logging driver, set the `log-driver` and
`log-opt` keys to appropriate values in the `daemon.json` file, which is located
in `/etc/docker/` on Linux hosts or `C:\ProgramData\docker\config\daemon.json`
on Windows Server. For more about configuring Docker using `daemon.json`, see
[daemon.json](../../../engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `gelf` and sets the `gelf-address`
option.

```json
{
  "log-driver": "gelf",
  "log-opts": {
    "gelf-address": "udp://1.2.3.4:12201"
  }
}
```

Restart Docker for the changes to take effect.

> **Note**
>
> `log-opts` configuration options in the `daemon.json` configuration file must
> be provided as strings. Boolean and numeric values (such as the value for
> `gelf-tcp-max-reconnect`) must therefore be enclosed in quotes (`"`).

You can set the logging driver for a specific container by setting the
`--log-driver` flag when using `docker container create` or `docker run`:

```bash
$ docker run \
      --log-driver gelf –-log-opt gelf-address=udp://1.2.3.4:12201 \
      alpine echo hello world
```

### GELF options

The `gelf` logging driver supports the following options:

| Option                     | Required  | Description                                                                                                                                                                                                                                                                         | Example value                                       |
| :------------------------- | :-------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------------------------- |
| `gelf-address`             | required  | The address of the GELF server. `tcp` and `udp` are the only supported URI specifier and you must specify the port.                                                                                                                                                                 | `--log-opt gelf-address=udp://192.168.0.42:12201`   |
| `gelf-compression-type`    | optional  | `UDP Only` The type of compression the GELF driver uses to compress each log message. Allowed values are `gzip`, `zlib` and `none`. The default is `gzip`. **Note that enabled compression leads to excessive CPU usage, so it is highly recommended to set this to `none`**.       | `--log-opt gelf-compression-type=gzip`              |
| `gelf-compression-level`   | optional  | `UDP Only` The level of compression when `gzip` or `zlib` is the `gelf-compression-type`. An integer in the range of `-1` to `9` (BestCompression). Default value is 1 (BestSpeed). Higher levels provide more compression at lower speed. Either `-1` or `0` disables compression. | `--log-opt gelf-compression-level=2`                |
| `gelf-tcp-max-reconnect`   | optional  | `TCP Only` The maximum number of reconnection attempts when the connection drop. An positive integer. Default value is 3.                                                                                                                                                           | `--log-opt gelf-tcp-max-reconnect=3`                |
| `gelf-tcp-reconnect-delay` | optional  | `TCP Only` The number of seconds to wait between reconnection attempts. A positive integer. Default value is 1.                                                                                                                                                                     | `--log-opt gelf-tcp-reconnect-delay=1`              |
| `tag`                      | optional  | A string that is appended to the `APP-NAME` in the `gelf` message. By default, Docker uses the first 12 characters of the container ID to tag log messages. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format.                            | `--log-opt tag=mailer`                              |
| `labels`                   | optional  | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon accepts. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md).                                        | `--log-opt labels=production_status,geo`            |
| `env`                      | optional  | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon accepts. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md).                         | `--log-opt env=os,customer`                         |
| `env-regex`                | optional  | Similar to and compatible with `env`. A regular expression to match logging-related environment variables. Used for advanced [log tag options](log_tags.md).                                                                                                                        | `--log-opt env-regex=^(os|customer)`                |

### Examples

This example configures the container to use the GELF server running at
`192.168.0.42` on port `12201`.

```bash
$ docker run -dit \
    --log-driver=gelf \
    --log-opt gelf-address=udp://192.168.0.42:12201 \
    alpine sh
```
