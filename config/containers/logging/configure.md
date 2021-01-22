---
description: Configure logging driver.
keywords: docker, logging, driver
redirect_from:
- /engine/reference/logging/overview/
- /engine/reference/logging/
- /engine/admin/reference/logging/
- /engine/admin/logging/overview/
title: Configure logging drivers
---

Docker includes multiple logging mechanisms to help you
[get information from running containers and services](index.md).
These mechanisms are called logging drivers. Each Docker daemon has a default
logging driver, which each container uses unless you configure it to use a
different logging driver, or "log-driver" for short.

As a default, Docker uses the [`json-file` logging driver](json-file.md), which
caches container logs as JSON internally. In addition to using the logging drivers
included with Docker, you can also implement and use [logging driver plugins](plugins.md).

> **Tip: use the "local" logging driver to prevent disk-exhaustion** 
> 
> By default, no log-rotation is performed. As a result, log-files stored by the
> default [`json-file` logging driver](json-file.md) logging driver can cause 
> a significant amount of disk space to be used for containers that generate much
> output, which can lead to disk space exhaustion.
> 
> Docker keeps the json-file logging driver (without log-rotation) as a default
> to remain backward compatibility with older versions of Docker, and for situations
> where Docker is used as runtime for Kubernetes.
> 
> For other situations, the "local" logging driver is recommended as it performs
> log-rotation by default, and uses a more efficient file format. Refer to the
> [Configure the default logging driver](#configure-the-default-logging-driver)
> section below to learn how to configure the "local" logging driver as a default,
> and the [local file logging driver](local.md) page for more details about the
> "local" logging driver.

## Configure the default logging driver

To configure the Docker daemon to default to a specific logging driver, set the
value of `log-driver` to the name of the logging driver in the `daemon.json`
configuration file. Refer to the "daemon configuration file" section in the
[`dockerd` reference manual](/engine/reference/commandline/dockerd/#daemon-configuration-file)
for details.

The default logging driver is `json-file`. The following example sets the default
logging driver to the [`local` log driver](local.md):

```json
{
  "log-driver": "local"
}
```

If the logging driver has configurable options, you can set them in the
`daemon.json` file as a JSON object with the key `log-opts`. The following
example sets two configurable options on the `json-file` logging driver:

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3",
    "labels": "production_status",
    "env": "os,customer"
  }
}
```

Restart Docker for the changes to take effect for newly created containers.
Existing containers do not use the new logging configuration.

> **Note**
>
> `log-opts` configuration options in the `daemon.json` configuration file must
> be provided as strings. Boolean and numeric values (such as the value for
> `max-file` in the example above) must therefore be enclosed in quotes (`"`).

If you do not specify a logging driver, the default is `json-file`.
To find the current default logging driver for the Docker daemon, run
`docker info` and search for `Logging Driver`. You can use the following
command on Linux, macOS, or PowerShell on Windows:

{% raw %}
```console
$ docker info --format '{{.LoggingDriver}}'

json-file
```
{% endraw %}

> **Note**
> 
> Changing the default logging driver or logging driver options in the daemon
> configuration only affects containers that are created after the configuration
> is changed. Existing containers retain the logging driver options that were
> used when they were created. To update the logging driver for a container, the
> container has to be re-created with the desired options.
> Refer to the [configure the logging driver for a container](#configure-the-logging-driver-for-a-container)
> section below to learn how to find the logging-driver configuration of a
> container.

## Configure the logging driver for a container

When you start a container, you can configure it to use a different logging
driver than the Docker daemon's default, using the `--log-driver` flag. If the
logging driver has configurable options, you can set them using one or more
instances of the `--log-opt <NAME>=<VALUE>` flag. Even if the container uses the
default logging driver, it can use different configurable options.

The following example starts an Alpine container with the `none` logging driver.

```console
$ docker run -it --log-driver none alpine ash
```

To find the current logging driver for a running container, if the daemon
is using the `json-file` logging driver, run the following `docker inspect`
command, substituting the container name or ID for `<CONTAINER>`:

{% raw %}
```console
$ docker inspect -f '{{.HostConfig.LogConfig.Type}}' <CONTAINER>

json-file
```
{% endraw %}

## Configure the delivery mode of log messages from container to log driver

Docker provides two modes for delivering messages from the container to the log
driver:

* (default) direct, blocking delivery from container to driver
* non-blocking delivery that stores log messages in an intermediate per-container
  ring buffer for consumption by driver

The `non-blocking` message delivery mode prevents applications from blocking due
to logging back pressure. Applications are likely to fail in unexpected ways when
STDERR or STDOUT streams block.

> **Warning**
>
> When the buffer is full and a new message is enqueued, the oldest message in
> memory is dropped.  Dropping messages is often preferred to blocking the
> log-writing process of an application.
{: .warning}

The `mode` log option controls whether to use the `blocking` (default) or
`non-blocking` message delivery.

The `max-buffer-size` log option controls the size of the ring buffer used for
intermediate message storage when `mode` is set to `non-blocking`. `max-buffer-size`
defaults to 1 megabyte.

The following example starts an Alpine container with log output in non-blocking
mode and a 4 megabyte buffer:

```console
$ docker run -it --log-opt mode=non-blocking --log-opt max-buffer-size=4m alpine ping 127.0.0.1
```

### Use environment variables or labels with logging drivers

Some logging drivers add the value of a container's `--env|-e` or `--label`
flags to the container's logs. This example starts a container using the Docker
daemon's default logging driver (let's assume `json-file`) but sets the
environment variable `os=ubuntu`.

```console
$ docker run -dit --label production_status=testing -e os=ubuntu alpine sh
```

If the logging driver supports it, this adds additional fields to the logging
output. The following output is generated by the `json-file` logging driver:

```json
"attrs":{"production_status":"testing","os":"ubuntu"}
```

## Supported logging drivers

The following logging drivers are supported. See the link to each driver's
documentation for its configurable options, if applicable. If you are using
[logging driver plugins](plugins.md), you may
see more options.

| Driver                        | Description                                                                                                   |
|:------------------------------|:--------------------------------------------------------------------------------------------------------------|
| `none`                        | No logs are available for the container and `docker logs` does not return any output.                         |
| [`local`](local.md)           | Logs are stored in a custom format designed for minimal overhead.                                             |
| [`json-file`](json-file.md)   | The logs are formatted as JSON. The default logging driver for Docker.                                        |
| [`syslog`](syslog.md)         | Writes logging messages to the `syslog` facility. The `syslog` daemon must be running on the host machine.    |
| [`journald`](journald.md)     | Writes log messages to `journald`. The `journald` daemon must be running on the host machine.                 |
| [`gelf`](gelf.md)             | Writes log messages to a Graylog Extended Log Format (GELF) endpoint such as Graylog or Logstash.             |
| [`fluentd`](fluentd.md)       | Writes log messages to `fluentd` (forward input). The `fluentd` daemon must be running on the host machine.   |
| [`awslogs`](awslogs.md)       | Writes log messages to Amazon CloudWatch Logs.                                                                |
| [`splunk`](splunk.md)         | Writes log messages to `splunk` using the HTTP Event Collector.                                               |
| [`etwlogs`](etwlogs.md)       | Writes log messages as Event Tracing for Windows (ETW) events. Only available on Windows platforms.           |
| [`gcplogs`](gcplogs.md)       | Writes log messages to Google Cloud Platform (GCP) Logging.                                                   |
| [`logentries`](logentries.md) | Writes log messages to Rapid7 Logentries.                                                                     |

> **Note**
>
> When using Docker Engine 19.03 or older, the [`docker logs` command](../../../engine/reference/commandline/logs.md)
> is only functional for the `local`, `json-file` and `journald` logging drivers.
> Docker 20.10 and up introduces "dual logging", which uses a local buffer that
> allows  you to use the `docker logs` command for any logging driver. Refer to
> [reading logs when using remote logging drivers](dual-logging.md) for details.

## Limitations of logging drivers

- Reading log information requires decompressing rotated log files, which causes
  a temporary increase in disk usage (until the log entries from the rotated
  files are read) and an increased CPU usage while decompressing. 
- The capacity of the host storage where the Docker data directory resides
  determines the maximum size of the log file information.
