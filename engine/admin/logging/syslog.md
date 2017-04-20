---
description: Describes how to use the syslog logging driver.
keywords: syslog, docker, logging, driver
redirect_from:
- /engine/reference/logging/syslog/
title: Syslog logging driver
---

The `syslog` logging driver routes logs to a `syslog` server.
The `syslog` protocol uses a raw string as log message supporting a very little number of metadata, for a syslog message to be valid the syslog message must be formatted in a specific way, allowing the receiver to extract the following information:

- a priority: is this a debug message, a warning, something purely informational, a critical error, etc.;
- a timestamp indicating when the thing happened;
- a hostname indicating where the thing happened (i.e. on which machine);
- a facility indicating if the message comes from the mail system, the kernel, and such and such;
- a process name and number;

The format has been standardized in [RFC 5424](https://tools.ietf.org/html/rfc5424) and this logging driver implements the [ABNF reference](https://tools.ietf.org/html/rfc5424#section-6) in the following way

```
                TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID
                    +          +             +           |        +
                    |          |             |           |        |
                    |          |             |           |        |
       +------------+          +----+        |           +----+   +---------+
       v                            v        v                v             v
2017-04-01T17:41:05.616647+08:00 a.vm {taskid:aa,version:} 1787791 {taskid:aa,version:}
```

## Usage

You can configure the default logging driver by passing the `--log-driver` and `--log-opt` options to the Docker daemon:

```bash
dockerd \
  --log-driver syslog \
  --log-opt syslog-address=udp://1.2.3.4:1111
```

**Please note that** the syslog-address supports both **udp** and **tcp**.

Or you can configure it in `/etc/docker/daemon.json` with:


```json
{
  "log-driver": "syslog",
  "log-opts":  {
    "syslog": "udp://1.2.3.4:1111"
  }
}
```

You can set the logging driver for a specific container by using the
`--log-driver` option to `docker run`:

```bash
docker run \
      -–log-driver syslog –-log-opt syslog=udp://1.2.3.4:1111 \
      alpine echo hello world
```

## Options

The following logging options are supported for the `syslog` logging driver:

| Option               | Description              | Example value                             |
|----------------------|--------------------------|-------------------------------------------|
| `syslog-address`     | The address of an external `syslog` server. The URI specifier may be `[tcp|udp|tcp+tls]://host:port`, `unix://path`, or `unixgram://path`. If the transport is `tcp`, `udp`, or `tcp+tls`, the default port is `514`.| `--log-opt syslog-address=tcp+tls://192.168.1.3:514`, `--log-opt syslog-address=unix:///tmp/syslog.sock` |
| `syslog-facility`    | The `syslog` facility to use. Can be the number or name for any valid `syslog` facility. See the [syslog documentation](https://tools.ietf.org/html/rfc5424#section-6.2.1). | `--log-opt syslog-facility=daemon` |
| `syslog-tls-ca-cert` | The absolute path to the trust certificates signed by the CA. **Ignored if the address protocol is not `tcp+tls`.** | `--log-opt syslog-tls-ca-cert=/etc/ca-certificates/custom/ca.pem` |
| `syslog-tls-cert`    | The absolute path to the TLS certificate file. **Ignored if the address protocol is not `tcp+tls`**. | `--log-opt syslog-tls-cert=/etc/ca-certificates/custom/cert.pem` |
| `syslog-tls-key`     | The absolute path to the TLS key file. **Ignored if the address protocol is not `tcp+tls`.** | `--log-opt syslog-tls-key=/etc/ca-certificates/custom/key.pem` |
| `syslog-tls-skip`    | If set to `true`, TLS verification is skipped when connecting to the `syslog` daemon. Defaults to `false`. **Ignored if the address protocol is not `tcp+tls`.** | `--log-opt syslog-tls-skip-verify=true` |
| `tag`                | A string that is appended to the `APP-NAME` in the `syslog` message. By default, Docker uses the first 12 characters of the container ID to tag log messages. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format. | `--log-opt tag=mailer` |
| `syslog-format`      | The `syslog` message format to use. If not specified the local UNIX syslog format is used, without a specified hostname. Specify `rfc3164` for the RFC-3164 compatible format, `rfc5424` for RFC-5424 compatible format, or `rfc5424micro` for RFC-5424 compatible format with microsecond timestamp resolution. | `--log-opt syslog-format=rfc5424micro` |
| `labels`             | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`                | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |
