---
description: Learn how to read container logs locally when using a third party logging solution.
keywords: docker, logging, driver
title: Use docker logs to read container logs for remote logging drivers
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
to [Configure the default logging driver](/config/containers/logging/configure)
for additional information. 


## Prerequisites 
 
No configuration changes are needed to use dual logging. Docker Engine 20.10 and
up automatically enable dual logging if the configured logging driver does not
support reading logs.

The following examples show the result of running a `docker logs` command with
and without dual logging availability:

### Without dual logging capability:

When a container or `dockerd` was configured with a remote logging driver such
as `splunk`, an error was displayed when attempting to read container logs
locally:

- Step 1: Configure Docker daemon

    ```console
    $ cat /etc/docker/daemon.json
    {
      "log-driver": "splunk",
      "log-opts": {
        ...
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

### With dual logging capability:

To configure a container or docker with a remote logging driver such as splunk:

- Step 1: Configure Docker daemon

    ```console
    $ cat /etc/docker/daemon.json
    {
      "log-driver": "splunk",
      "log-opts": {
        ...
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
> For a local driver, such as `json-file` and `journald`, there is no difference in
> functionality before or after the dual logging capability became available.
> The log is locally visible in both scenarios.


## Limitations

- You cannot specify more than one log driver. 
- If a container using a logging driver or plugin that sends logs remotely
  suddenly has a "network" issue, no ‘write’ to the local cache occurs. 
- If a write to `logdriver` fails for any reason (file system full, write
  permissions removed), the cache write fails and is logged in the daemon log.
  The log entry to the cache is not retried.
- Some logs might be lost from the cache in the default configuration because a
  ring buffer is used to prevent blocking the stdio of the container in case of
  slow file writes. An admin must repair these while the daemon is shut down.
