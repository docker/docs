---
description: Learn how to read container logs locally when using a third party logging solution.
keywords: docker, logging, driver, dual logging, dual-logging, cache, ring-buffer, configuration
title: Use docker logs with remote logging drivers
---

## Overview 

Prior to Docker Engine 20.10, the [`docker logs` command](../../../engine/reference/commandline/logs.md)
could only be used with logging drivers that supported  for containers using the
`local`, `json-file`, or `journald` log drivers. However, many third party logging
drivers had no support for locally reading logs using `docker logs`

This created multiple problems when attempting to gather log data in an
automated and standard way. Log information could only be accessed and viewed
through the third-party solution in the format specified by that
third-party tool. 

Starting with Docker Engine 20.10, you can use `docker logs` to read container
logs regardless of the configured logging driver or plugin. This capability,
referred to as "dual logging", allows you to use `docker logs` to read container
logs locally in a consistent format, regardless of the log driver used, because
the engine is configured to log information to the “local” logging driver. Refer
to [Configure the default logging driver](configure.md) for additional information. 

Dual logging uses the [`local`](local.md) logging driver to act as cache for
reading the latest logs of your containers. By default, the cache has log-file
rotation enabled, and is limited to a maximum of 5 files of 20MB each (before
compression) per container.

Refer to the [configuration options](#configuration-options) section to customize
these defaults, or to the [disable dual-logging](#disable-the-dual-logging-cache)
section to disable this feature.

## Prerequisites 
 
No configuration changes are needed to use dual logging. Docker Engine 20.10 and
up automatically enable dual logging if the configured logging driver does not
support reading logs.

The following examples show the result of running a `docker logs` command with
and without dual logging availability:

### Without dual logging capability

When a container is configured with a remote logging driver such as `splunk`, and
dual logging is disabled, an error is displayed when attempting to read container
logs locally:

- Step 1: Configure Docker daemon

    ```console
    $ cat /etc/docker/daemon.json
    {
      "log-driver": "splunk",
      "log-opts": {
        "cache-disabled": "true",
        ... (options for "splunk" logging driver)
      }
    }
    ```

- Step 2: Start the container

    ```console
    $ docker run -d busybox --name testlog top 
    ```

- Step 3: Read the container logs

    ```console
    $ docker logs 7d6ac83a89a0
    Error response from daemon: configured logging driver does not support reading
    ```

### With dual logging capability

With the dual logging cache enabled, the `docker logs` command can be used to
read logs, even if the logging driver does not support reading logs. The following
examples shows a daemon configuration that uses the `splunk` remote logging driver
as a default, with dual logging caching enabled:

- Step 1: Configure Docker daemon

    ```console
    $ cat /etc/docker/daemon.json
    {
      "log-driver": "splunk",
      "log-opts": {
        ... (options for "splunk" logging driver)
      }
    }
    ```

- Step 2: Start the container

    ```console
    $ docker run -d busybox --name testlog top 
    ```

- Step 3: Read the container logs

    ```console
    $ docker logs 7d6ac83a89a0
    2019-02-04T19:48:15.423Z [INFO]  core: marked as sealed                                          	 
    2019-02-04T19:48:15.423Z [INFO]  core: pre-seal teardown starting                                                                                                 	 
    2019-02-04T19:48:15.423Z [INFO]  core: stopping cluster listeners                                                                                             	 
    2019-02-04T19:48:15.423Z [INFO]  core: shutting down forwarding rpc listeners                                                                                 	 
    2019-02-04T19:48:15.423Z [INFO]  core: forwarding rpc listeners stopped
    2019-02-04T19:48:15.599Z [INFO]  core: rpc listeners successfully shut down
    2019-02-04T19:48:15.599Z [INFO]  core: cluster listeners successfully shut down	
    ```

> **Note**
>
> For logging drivers that support reading logs, such as the `local`, `json-file`
> and `journald` drivers, there is no difference in functionality before or after
> the dual logging capability became available. For these drivers, Logs can be
> read using `docker logs` in both scenarios.


### Configuration options

The "dual logging" cache accepts the same configuration options as the
[`local` logging driver](local.md), but with a `cache-` prefix. These options
can be specified per container, and defaults for new containers can be set using
the [daemon configuration file](/engine/reference/commandline/dockerd/#daemon-configuration-file).

By default, the cache has log-file rotation enabled, and is limited to a maximum
of 5 files of 20MB each (before compression) per container. Use the configuration
options described below to customize these defaults.


| Option           | Default   | Description                                                                                                                                       |
|:-----------------|:----------|:--------------------------------------------------------------------------------------------------------------------------------------------------|
| `cache-disabled` | `"false"` | Disable local caching. Boolean value passed as a string (`true`, `1`, `0`, or `false`).                                                           |
| `cache-max-size` | `"20m"`   | The maximum size of the cache before it is rotated. A positive integer plus a modifier representing the unit of measure (`k`, `m`, or `g`).       |
| `cache-max-file` | `"5"`     | The maximum number of cache files that can be present. If rotating the logs creates excess files, the oldest file is removed. A positive integer. |
| `cache-compress` | `"true"`  | Enable or disable compression of rotated log files. Boolean value passed as a string (`true`, `1`, `0`, or `false`).                              |

## Disable the dual logging cache

Use the `cache-disabled` option to disable the dual logging cache. Disabling the
cache can be useful to save storage space in situations where logs are only read
through a remote logging system, and if there is no need to read logs through
`docker logs` for debugging purposes.

Caching can be disabled for individual containers or by default for new containers,
when using the [daemon configuration file](/engine/reference/commandline/dockerd/#daemon-configuration-file).

The following example uses the daemon configuration file to use the ["splunk'](splunk.md)
logging driver as a default, with caching disabled:

```console
$ cat /etc/docker/daemon.json
{
  "log-driver": "splunk",
  "log-opts": {
    "cache-disabled": "true",
    ... (options for "splunk" logging driver)
  }
}
```

> **Note**
>
> For logging drivers that support reading logs, such as the `local`, `json-file`
> and `journald` drivers, dual logging is not used, and disabling the option has
> no effect.

## Limitations

- If a container using a logging driver or plugin that sends logs remotely
  suddenly has a "network" issue, no ‘write’ to the local cache occurs. 
- If a write to `logdriver` fails for any reason (file system full, write
  permissions removed), the cache write fails and is logged in the daemon log.
  The log entry to the cache is not retried.
- Some logs might be lost from the cache in the default configuration because a
  ring buffer is used to prevent blocking the stdio of the container in case of
  slow file writes. An admin must repair these while the daemon is shut down.
