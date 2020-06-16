---
description: Describes how to use the syslog logging driver.
keywords: syslog, docker, logging, driver
redirect_from:
- /engine/reference/logging/syslog/
- /engine/admin/logging/syslog/
title: Syslog logging driver
---

The `syslog` logging driver routes logs to a `syslog` server. The `syslog` protocol uses
a raw string as the log message and supports a limited set of metadata. The syslog
message must be formatted in a specific way to be valid. From a valid message, the
receiver can extract the following information:

- **priority**: the logging level, such as `debug`, `warning`, `error`, `info`.
- **timestamp**: when the event occurred.
- **hostname**: where the event happened.
- **facility**: which subsystem logged the message, such as `mail` or `kernel`.
- **process name** and **process ID (PID)**: The name and ID of the process that generated the log.

The format is defined in [RFC 5424](https://tools.ietf.org/html/rfc5424) and Docker's syslog driver implements the
[ABNF reference](https://tools.ietf.org/html/rfc5424#section-6) in the following way:

```none
                TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID
                    +          +             +           |        +
                    |          |             |           |        |
                    |          |             |           |        |
       +------------+          +----+        |           +----+   +---------+
       v                            v        v                v             v
2017-04-01T17:41:05.616647+08:00 a.vm {taskid:aa,version:} 1787791 {taskid:aa,version:}
```

## Usage

To use the `syslog` driver as the default logging driver, set the `log-driver`
and `log-opt` keys to appropriate values in the `daemon.json` file, which is
located in `/etc/docker/` on Linux hosts or
`C:\ProgramData\docker\config\daemon.json` on Windows Server. For more about
configuring Docker using `daemon.json`, see
[daemon.json](../../../engine/reference/commandline/dockerd.md#daemon-configuration-file).

The following example sets the log driver to `syslog` and sets the
`syslog-address` option. The `syslog-address` options supports both UDP and TCP;
this example uses UDP.

```json
{
  "log-driver": "syslog",
  "log-opts": {
    "syslog-address": "udp://1.2.3.4:1111"
  }
}
```

Restart Docker for the changes to take effect.

> **Note**
>
> `log-opts` configuration options in the `daemon.json` configuration file must
> be provided as strings. Numeric and boolean values (such as the value for
> `syslog-tls-skip-verify`) must therefore be enclosed in quotes (`"`).

You can set the logging driver for a specific container by using the
`--log-driver` flag to `docker container create` or `docker run`:

```bash
docker run \
      --log-driver syslog --log-opt syslog-address=udp://1.2.3.4:1111 \
      alpine echo hello world
```

## Options

The following logging options are supported as options for the `syslog` logging
driver. They can be set as defaults in the `daemon.json`, by adding them as
key-value pairs to the `log-opts` JSON array. They can also be set on a given
container by adding a `--log-opt <key>=<value>` flag for each option when
starting the container.

| Option                   | Description                                                                                                                                                                                                                                                                                                      | Example value                                                                                            |
|:-------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------------------------------------------------------------------------------------------------|
| `syslog-address`         | The address of an external `syslog` server. The URI specifier may be `[tcp|udp|tcp+tls]://host:port`, `unix://path`, or `unixgram://path`. If the transport is `tcp`, `udp`, or `tcp+tls`, the default port is `514`.                                                                                            | `--log-opt syslog-address=tcp+tls://192.168.1.3:514`, `--log-opt syslog-address=unix:///tmp/syslog.sock` |
| `syslog-facility`        | The `syslog` facility to use. Can be the number or name for any valid `syslog` facility. See the [syslog documentation](https://tools.ietf.org/html/rfc5424#section-6.2.1).                                                                                                                                      | `--log-opt syslog-facility=daemon`                                                                       |
| `syslog-tls-ca-cert`     | The absolute path to the trust certificates signed by the CA. **Ignored if the address protocol is not `tcp+tls`.**                                                                                                                                                                                              | `--log-opt syslog-tls-ca-cert=/etc/ca-certificates/custom/ca.pem`                                        |
| `syslog-tls-cert`        | The absolute path to the TLS certificate file. **Ignored if the address protocol is not `tcp+tls`**.                                                                                                                                                                                                             | `--log-opt syslog-tls-cert=/etc/ca-certificates/custom/cert.pem`                                         |
| `syslog-tls-key`         | The absolute path to the TLS key file. **Ignored if the address protocol is not `tcp+tls`.**                                                                                                                                                                                                                     | `--log-opt syslog-tls-key=/etc/ca-certificates/custom/key.pem`                                           |
| `syslog-tls-skip-verify` | If set to `true`, TLS verification is skipped when connecting to the `syslog` daemon. Defaults to `false`. **Ignored if the address protocol is not `tcp+tls`.**                                                                                                                                                 | `--log-opt syslog-tls-skip-verify=true`                                                                  |
| `tag`                    | A string that is appended to the `APP-NAME` in the `syslog` message. By default, Docker uses the first 12 characters of the container ID to tag log messages. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format.                                                       | `--log-opt tag=mailer`                                                                                   |
| `syslog-format`          | The `syslog` message format to use. If not specified the local UNIX syslog format is used, without a specified hostname. Specify `rfc3164` for the RFC-3164 compatible format, `rfc5424` for RFC-5424 compatible format, or `rfc5424micro` for RFC-5424 compatible format with microsecond timestamp resolution. | `--log-opt syslog-format=rfc5424micro`                                                                   |
| `labels`                 | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon accepts. Used for advanced [log tag options](log_tags.md).                                                                                                                                                 | `--log-opt labels=production_status,geo`                                                                 |
| `env`                    | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon accepts. Used for advanced [log tag options](log_tags.md).                                                                                                                                  | `--log-opt env=os,customer`                                                                              |
| `env-regex`              | Applies when starting the Docker daemon. Similar to and compatible with `env`. A regular expression to match logging-related environment variables. Used for advanced [log tag options](log_tags.md).                                                                                                            | `--log-opt env-regex=^(os\|customer)`                                                                    |
