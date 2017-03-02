---
description: Configure logging driver.
keywords: docker, logging, driver, Fluentd
redirect_from:
title: Configure logging drivers
---

Docker includes multiple logging mechanisms to help you get information from
running containers and services. These mechanisms are called logging
drivers.

Each Docker daemon has a default logging driver, which each container uses
unless you configure it to use a different logging driver.

## Configure the default logging driver for the Docker daemon

To configure the Docker daemon to default to a specific logging driver, use the
`--log-driver=<VALUE>` flag. If the logging driver has configurable options,
you can set them using one or more instances of the `--log-opt <NAME>=<VALUE>`
flag.

If you do not specify a logging driver, the default is `json-file`. Thus,
the default output for commands such as `docker inspect <CONTAINER>` is JSON.

To find the current default logging driver for the Docker daemon, run
`docker info` and search for `Logging Driver`. You can use the following
command on Linux or macOS:

```bash
$ docker info |grep 'Logging Driver'

Logging Driver: json-file
```

## Configure the logging driver for a container

When you start a container, you can configure it to use a different logging
driver than the Docker daemon's default. If the logging driver has configurable
options, you can set them using one or more instances of the
`--log-opt <NAME>=<VALUE>` flag. Even if the container uses the default logging
driver, it can use different configurable options.

To find the current logging driver for a running container, if the daemon
is using the `json-file` logging driver, run the following `docker inspect`
command, substituting the container name or ID for `<CONTAINER>`:

```bash
{% raw %}
$ docker inspect -f '{{.HostConfig.LogConfig.Type}}' <CONTAINER>

json-file
{% endraw %}
```

## Supported logging drivers

The following logging drivers are supported. See each driver's section below
for its configurable options, if applicable.

| Driver      | Description                                                    |
|-------------|----------------------------------------------------------------|
| `none`      | No logs will be available for the container and `docker logs` will not return any output. |
| `json-file` | The logs are formatted as JSON. The default logging driver for Docker. |
| `syslog`    | Writes logging messages to the `syslog` facility. The `syslog` daemon must be running on the host machine. |
| `journald`  | Writes log messages to `journald`. The `journald` daemon must be running on the host machine. |
| `gelf`      | Writes log messages to a Graylog Extended Log Format (GELF) endpoint such as Graylog or Logstash. |
| `fluentd`   | Writes log messages to `fluentd` (forward input). The `fluentd` daemon must be running on the host machine. |                                          |
| `awslogs`   | Writes log messages to Amazon CloudWatch Logs.                 |
| `splunk`    | Writes log messages to `splunk` using the HTTP Event Collector.|
| `etwlogs`   | Writes log messages as Event Tracing for Windows (ETW) events. Only available on Windows platforms. |
| `gcplogs`   | Writes log messages to Google Cloud Platform (GCP) Logging.    |
| `nats`      | NATS logging driver for Docker. Publishes log entries to a NATS server.|

## Limitations of logging drivers

- The `docker logs` command is not available for drivers other than `json-file`
  and `journald`.

## Examples

### Configure the logging driver using labels or environment variables

If your container uses labels or environment variables specified in the
Dockerfile or at runtime, some logging drivers can use these labels or
environment variables to control logging behavior. If a collision occurs
between a label and an environment variable, the environment variable takes
precedence.

Specify logging attributes and options when you start the Docker daemon. For
example, to manually start the daemon with the `json-file` driver, and set a
label and two environment variables, use the following command:

```bash
$ dockerd \
         --log-driver=json-file \
         --log-opt labels=production_status \
         --log-opt env=os,customer
```

Next, run a container and specify values for the `labels` or `env`. For
example, you might use this:

```bash
$ docker run -dit --label production_status=testing -e os=ubuntu alpine sh
```

If the logging driver supports it, this adds additional fields to the logging
output. The following output is for `json-file`:

```json
"attrs":{"production_status":"testing","os":"ubuntu"}
```

## `none`

The `none` driver disables logging for the Docker daemon (if set on the daemon
at start-up) or for an individual container at runtime. It has no options.

### Examples

This example starts an `alpine` container with the `none` log driver.

```bash
$ docker run -it --log-driver none alpine ash
```

## `json-file`

`json-file` is the default logging driver, and returns logging output in JSON
format.

### Options

The `json-file` logging driver supports the following logging options:

| Option     | Description              | Example  value                            |
|------------|--------------------------|-------------------------------------------|
| `max-size` | The maximum size of the log before it is rolled. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`). | `--log-opt max-size=10m` |
| `max-file` | The maximum number of log files that can be present. If rolling the logs creates excess files, the oldest file is removed. **Only effective when `max-size` is also set.** A positive integer. | `--log-opt max-file=3` |
| `labels`   | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`      | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |

> **Note**: If `max-size` and `max-file` are set, `docker logs` only returns the
> log lines from the newest log file.

### Examples

This example starts an `alpine` container which can have a maximum of 3 log
files no larger than 10 megabytes each.

```bash
$ docker run --it --log-opt max-size=10m --log-opt max-file=3 alpine ash
```

## `syslog`

### Options

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

### Examples

This example sends the container's logging output to a `syslog` remote server at
`192.168.0.42` on port `123`, using the the `daemon` facility:

```bash
$ docker run \
         --log-driver=syslog \
         --log-opt syslog-address=tcp://192.168.0.42:123 \
         --log-opt syslog-facility=daemon \
         alpine ash
```

This example connects to `syslog` using TCP+TLS transport and specifies the
trust certificate, certificate, and key to use.

```bash
$ docker run \
         --log-driver=syslog \
         --log-opt syslog-address=tcp+tls://192.168.0.42:123 \
         --log-opt syslog-tls-ca-cert=syslog-tls-ca-cert=/etc/ca-certificates/custom/ca.pem \
         --log-opt syslog-tls-cert=syslog-tls-ca-cert=/etc/ca-certificates/custom/cert.pem \
         --log-opt syslog-tls-key=syslog-tls-ca-cert=/etc/ca-certificates/custom/key.pem \
         alpine ash
```

## `journald`

The `journald` logging driver stores the container id in the journal's
`CONTAINER_ID` field. For detailed information on working with this logging
driver, see [the journald logging driver](journald.md) reference documentation.

### Options

| Option     | Description              | Example value                             |
|------------|--------------------------|-------------------------------------------|
| `tag`      | A template for setting the `CONTAINER_TAG` value in `journald` logs. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format. | `--log-opt tag=mailer` |
| `labels`   | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`      | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |


### Examples

```bash
$ docker run \
         --log-driver=journald \
         alpine ash
```

## `gelf`

### Options

The `gelf` logging driver supports the following options:

| Option               | Description              | Example value                             |
|----------------------|--------------------------|-------------------------------------------|
| `gelf-address`       | The address of the GELF server. `udp` is the only supported URI specifier and you must specify the port.| `--log-opt gelf-address=udp://192.168.0.42:12201` |
| `gelf-compression-type` | The type of compression the GELF driver uses to compress each log message. Allowed values are `gzip`, `zlib` and `none`. The default is `gzip`. | `--log-opt gelf-compression-type=gzip` |
| `gelf-compression-level` | The level of compression when `gzip` or `zlib` is the `gelf-compression-type`. An integer in the range of `-1` to `9` (BestCompression). Default value is 1 (BestSpeed). Higher levels provide more compression at lower speed.| `--log-opt gelf-compression-level=2` |
| `tag`                | A string that is appended to the `APP-NAME` in the `gelf` message. By default, Docker uses the first 12 characters of the container ID to tag log messages. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format. | `--log-opt tag=mailer` |
| `labels`             | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`                | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |


### Examples

This example connects the container to the GELF server running at
`192.168.0.42` on port `12201`.

```bash
$ docker run -dit \
             --log-driver=gelf \
             --log-opt gelf-address=udp://192.168.0.42:12201 \
             alpine sh
```


## `fluentd`

### Options

The `fluentd` logging driver supports the following options:

| Option                 | Description              | Example value                             |
|------------------------|--------------------------|-------------------------------------------|
| `fluentd-address`      | The address of the Fluentd server, in the format `host:port` with no protocol specifier.| `--log-opt fluentd-address=192.168.0.42:24224` |
| `fluentd-buffer-limit` | The maximum size of the fluentd log buffer, with a size prefix `KB`, `MB`, or `GB`. Defaults to `8MB`. | `--log-opt fluentd-buffer-limit=8MB` |
| `fluentd-retry-wait`   | The initial delay before retrying after a connection failure. After the initial delay, it increases exponentially. Defaults to `1000ms`.| `--log-opt fluentd-retry-wait=1000ms` |
| `fluentd-max-retries`  | The maximum number of connection attempts before the container stops due to failure to connect. Defaults to `1073741824`, which is effectively infinite. | `--log-opt fluentd-max-retries=200` |
| `fluentd-async-connect`| If set to `false`, Docker blocks on initial connection and the container stops if it cannot connect. Defaults to `false`. | `--log-opt fluentd-async-connect=false` |
| `tag`                  | A string that is appended to the `APP-NAME` in the `fluentd` message. By default, Docker uses the first 12 characters of the container ID to tag log messages. Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format. | `--log-opt tag=mailer` |

### Examples

This example logs container output to the Fluentd server running on `localhost`
at port `24224` and prepends the tag `docker.<CONTAINER_NAME>` to the beginning
of each message.

```bash
{% raw %}
$ docker run -dit \
             --log-driver=fluentd \
             --log-opt fluentd-address=localhost:24224 \
             --log-opt tag="docker.{{.Name}}" \
             alpine sh
{% endraw %}
```

For detailed information on working with the `fluentd` logging driver, see
[the fluentd logging driver](fluentd.md).


## `awslogs`

The Amazon Cloudwatch Logs driver is called `awslogs`.

### Options

The `awslogs` supports the following options:

| Option                 | Description              | Example value                             |
|------------------------|--------------------------|-------------------------------------------|
| `awslogs-region`       | Sets the region where the logs are sent. If not set, the container's region is used. | `--log-opt awslogs-region=us-east-1` |
| `awslogs-group`        | The log group to use for the logs. | `--log-opt awslogs-group=myLogGroup` |
| `awslogs-stream`       | The log stream to use. If not specified, the container ID is used as the log stream. | `--log-opt awslogs-stream=myLogStream` |

### Examples

This exampe sends the logs to region `us-east-1` and uses the log group
`myLogGroup`.

```bash
$ docker run \
         --log-driver=awslogs \
         --log-opt awslogs-region=us-east-1 \
         --log-opt awslogs-group=myLogGroup \
         alpine sh
```

For detailed information on working with this logging driver, see
[the `awslogs` logging driver](awslogs.md) reference documentation.

## `splunk`

The `splunk` logging driver sends container logs to the
[HTTP Event Collector](http://dev.splunk.com/view/event-collector/SP-CAAAE6M)
in Splunk Enterprise and Splunk Cloud.

### Options

The `splunk` logging driver **requires** the following options:

| Option                 | Description              | Example value                             |
|------------------------|--------------------------|-------------------------------------------|
| `splunk-token`  | The Splunk HTTP Event Collector token. | `--log-opt splunk-token=<splunk_http_event_collector_token>` |
| `splunk-url`    | Path to your Splunk Enterprise or Splunk Cloud instance (including port and scheme used by HTTP Event Collector).| `--log-opt splunk-url=https://your_splunk_instance:8088` |

The `splunk` logging driver **allows** the following options:

| Option                 | Description                     | Example value                             |
|------------------------|---------------------------------|-------------------------------------------|
| `splunk-source`             | Event source.              | `--log-opt splunk-token=176FCEBF-4CF5-4EDF-91BC-703796522D20` |
| `splunk-sourcetype`         | Event source type.         | `--log-opt splunk-sourcetype=iis`   |
| `splunk-index`              | Event index.               | `--log-opt splunk-index=os` |
| `splunk-capath`             | Path to root certificate.  | `--log-opt splunk-capath=/path/to/cert/cacert.pem` |
| `splunk-caname`             | Name to use for validating server certificate. Defaults to the hostname of the `splunk-url`. | `--log-opt splunk-caname=SplunkServerDefaultCert` |
| `splunk-insecureskipverify` | Ignore server certificate validation. | `--log-opt splunk-insecureskipverify=false` |
| `tag`                       | Specify tag for message, which interpret some markup. Default value is {% raw %}`{{.ID}}`{% endraw %} (12 characters of the container ID). Refer to the [log tag option documentation](log_tags.md) for information about customizing the log tag format. | {% raw %}`--log-opt tag="{{.Name}}/{{.FullID}}"`{% endraw %} |
| `labels`                    | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`                       | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |


### Examples

This examples sets several options for the `splunk` logging driver.

```bash
{% raw %}
$ docker run \
       --log-driver=splunk \
       --log-opt splunk-token=176FCEBF-4CF5-4EDF-91BC-703796522D20 \
       --log-opt splunk-url=https://splunkhost:8088 \
       --log-opt splunk-capath=/path/to/cert/cacert.pem \
       --log-opt splunk-caname=SplunkServerDefaultCert \
       --log-opt tag="{{.Name}}/{{.FullID}}" \
       --log-opt labels=location \
       --log-opt env=TEST \
       --env "TEST=false" \
       --label location=west \
   alpine sh
{% endraw %}
```

For detailed information about working with the `splunk` logging driver, see the
[Splunk logging driver](splunk.md) reference documentation.

## `etwlogs`

### Options

The `etwlogs` logging driver forwards each log message as an ETW event. An ETW
listener can then be created to listen for these events. This driver does not
accept any options.

### Examples

```bash
$ docker run \
         --logging-driver=etwlogs \
         alpine sh
```

The ETW logging driver is only available on Windows. For detailed information
about working with this logging driver, see [the ETW logging driver](etwlogs.md)
reference documentation.

## `gcplogs`

### Options

The Google Cloud Platform logging driver (`gcplogs`) supports the following
options:

| Option                 | Description                     | Example value                             |
|------------------------|---------------------------------|-------------------------------------------|
| `gcp-project`          | The GCP project to log to. Defaults to discovering this value from the GCE metadata service. | `--log-opt gcp-project=myProject` |
| `log-cmd`              | | `--log-opt log-cmd=true` | Whether to log the command that started the container. Defaults to `false`. | `--log-opt log-cmd=false` |
| `labels`               | Applies when starting the Docker daemon. A comma-separated list of logging-related labels this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md).| `--log-opt labels=production_status,geo` |
| `env`                  | Applies when starting the Docker daemon. A comma-separated list of logging-related environment variables this daemon will accept. Adds additional key on the `extra` fields, prefixed by an underscore (`_`). Used for advanced [log tag options](log_tags.md). | `--log-opt env=os,customer` |

### Examples

This example logs the start command and sets a label and an environment
variable, which will be incorporated into the logs if the Docker daemon was
started with the appropriate `--log-opt` options.

```bash
$ docker run --log-driver=gcplogs \
    --log-opt gcp-log-cmd=true \
    --env "TEST=false" \
    --label location=west \
    your/application \
```

For detailed information about working with the Google Cloud logging driver, see
the [Google Cloud Logging driver](gcplogs.md) reference documentation.

## NATS logging options

The NATS logging driver supports the following options:

```none
--log-opt labels=<label1>,<label2>
--log-opt env=<envvar1>,<envvar2>
--log-opt tag=<tag>
--log-opt nats-servers="<comma separated list of nats servers uris>"
--log-opt nats-max-reconnect="<max attempts to connect to a server>"
--log-opt nats-subject="<subject where logs are sent>"
--log-opt nats-tls-ca-cert="<absolute path to cert>"
--log-opt nats-tls-cert="<absolute path to cert>"
--log-opt nats-tls-key="<absolute path to cert>"
--log-opt nats-tls-skip-verify="<value>"
```

For detailed information, see [the NATS logging driver](nats.md) reference
documentation.
