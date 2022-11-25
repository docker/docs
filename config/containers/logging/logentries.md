---
title: Logentries logging driver
description: Describes how to use the logentries logging driver.
keywords: logentries, docker, logging, driver
redirect_from:
- /engine/admin/logging/logentries/
---

The `logentries` logging driver sends container logs to the
[Logentries](https://logentries.com/) server.

## Usage

Some options are supported by specifying `--log-opt` as many times as needed:

 - `logentries-token`: specify the logentries log set token
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

You need to provide your log set token for logentries driver to work:

```console
$ docker run --log-driver=logentries --log-opt logentries-token=abcd1234-12ab-34cd-5678-0123456789ab
```

### line-only

You could specify whether to send log message wrapped into container data (default) or to send raw log line

```console
$ docker run --log-driver=logentries --log-opt logentries-token=abcd1234-12ab-34cd-5678-0123456789ab --log-opt line-only=true
```
