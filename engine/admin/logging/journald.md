---
description: Describes how to use the fluentd logging driver.
keywords: Journald, docker, logging, driver
redirect_from:
- /engine/reference/logging/journald/
title: Journald logging driver
---

The `journald` logging driver sends container logs to the
[`systemd` journal](http://www.freedesktop.org/software/systemd/man/systemd-journald.service.html).
Log entries can be retrieved using the `journalctl` command, through use of the
`journal` API, or using the `docker logs` command.

In addition to the text of the log message itself, the `journald` log driver
stores the following metadata in the journal with each message:

| Field                       | Description |
------------------------------|-------------|
| `CONTAINER_ID`              | The container ID truncated to 12 characters. |
| `CONTAINER_ID_FULL`         | The full 64-character container ID. |
| `CONTAINER_NAME`            | The container name at the time it was started. If you use `docker rename` to rename a container, the new name is not reflected in the journal entries. |
| `CONTAINER_TAG`             | The container tag ([log tag option documentation](log_tags.md)). |
| `CONTAINER_PARTIAL_MESSAGE` | A field that flags log integrity. Improve logging of long log lines. |

## Usage

Configure the default logging driver by passing the `--log-driver` option to the
Docker daemon:

```bash
$ dockerd --log-driver=journald
```

To configure the logging driver for a specific container, use the `--log-driver`
flag on the `docker run` command.

```bash
$ docker run --log-driver=journald ...
```

## Options

Use the `--log-opt NAME=VALUE` flag to specify additional `journald` logging
driver options.

### `tag`

Specify template to set `CONTAINER_TAG` value in `journald` logs. Refer to
[log tag option documentation](log_tags.md) to customize the log tag format.

### `labels` and `env`

The `labels` and `env` options each take a comma-separated list of keys. If
there is collision between `label` and `env` keys, the value of the `env` takes
precedence. Each option adds additional metadata to the journal with each
message.

## Note regarding container names

The value logged in the `CONTAINER_NAME` field is the name of the container that
was set at startup. If you use `docker rename` to rename a container, the new
name **is not reflected** in the journal entries. Journal entries will continue
to use the original name.

## Retrieving log messages with `journalctl`

Use the `journalctl` command to retrieve log messages. You can apply filter
expressions to limit the retrieved messages to those associated with a specific
container:

```bash
$ sudo journalctl CONTAINER_NAME=webserver
```

You can use additional filters to further limit the messages retrieved. The `-b`
flag only retrieves messages generated since the last system boot:

```bash
$ sudo journalctl -b CONTAINER_NAME=webserver
```

The `-o` flag specifies the format for the retried log messages. Use `-o json`
to return the log messages in JSON format.

```bash
$ sudo journalctl -o json CONTAINER_NAME=webserver
```

## Retrieving log messages with the `journal` API

This example uses the `systemd` Python module to retrieve container
logs:

```python
import systemd.journal

reader = systemd.journal.Reader()
reader.add_match('CONTAINER_NAME=web')

    for msg in reader:
      print '{CONTAINER_ID_FULL}: {MESSAGE}'.format(**msg)
```
