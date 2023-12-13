---
title: Logentries logging driver (deprecated)
description: Learn how to use the logentries logging driver with Docker Engine
keywords: logentries, docker, logging, driver
aliases:
  - /engine/admin/logging/logentries/
---

> **Deprecated**
>
> The logentries service is no longer in operation since November 15, 2022,
> and the logentries driver [has been deprecated](../../../engine/deprecated.md#logentries-logging-driver).
> 
> This driver will be removed in Docker Engine v25.0, and you must migrate to
> a supported logging driver before upgrading to Docker Engine v25.0. Read the
> [Configure logging drivers](configure.md) page for an overview of supported
> logging drivers.
{ .warning }


The `logentries` logging driver sends container logs to the
[Logentries](https://logentries.com/) server.

## Usage

Some options are supported by specifying `--log-opt` as many times as needed:

- `logentries-token`: specify the Logentries log set token
- `line-only`: send raw payload only

Configure the default logging driver by passing the
`--log-driver` option to the Docker daemon:

```console
$ dockerd --log-driver=logentries
```

To set the logging driver for a specific container, pass the
`--log-driver` option to `docker run`:

```console
$ docker run --log-driver=logentries ...
```

Before using this logging driver, you need to create a new Log Set in the
Logentries web interface and pass the token of that log set to Docker:

```console
$ docker run --log-driver=logentries --log-opt logentries-token=abcd1234-12ab-34cd-5678-0123456789ab
```

## Options

Users can use the `--log-opt NAME=VALUE` flag to specify additional Logentries logging driver options.

### logentries-token

You need to provide your log set token for the Logentries driver to work:

```console
$ docker run --log-driver=logentries --log-opt logentries-token=abcd1234-12ab-34cd-5678-0123456789ab
```

### line-only

You could specify whether to send log message wrapped into container data (default) or to send raw log line

```console
$ docker run --log-driver=logentries --log-opt logentries-token=abcd1234-12ab-34cd-5678-0123456789ab --log-opt line-only=true
```
