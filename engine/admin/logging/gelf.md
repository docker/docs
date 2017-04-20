---
description: Describes how to use the Graylog Extended Format logging driver.
keywords: graylog, gelf, logging, driver
redirect_from:
- /engine/reference/logging/gelf/
title: Graylog Extended Format logging driver
---

The `gelf` logging driver is a convenient format that is understood by a number of tools like [Graylog](https://www.graylog.org/), [Logstash](https://www.elastic.co/products/logstash), [Fluentd](http://www.fluentd.org/), and many more.

In GELF, every log message is a dict with the following fields:

- version;
- host (who sent the message in the first place);
- timestamp;
- short and long version of the message;
- any extra field you would like!



## Usage

You can configure the default logging driver by passing the `--log-driver`
and `--log-opt` options to the Docker daemon:

```bash
dockerd
  -–log-driver gelf –-log-opt gelf-address=udp://1.2.3.4:12201 \
```

Or you can configure it in `/etc/docker/daemon.json` with:

```json
{
  "log-driver": "gelf",
  "log-opts":  {
    "gelf-address": "udp://1.2.3.4:12201"
  }
}
```

You can set the logging driver for a specific container by using the
`--log-driver` option to `docker run`:

```bash
docker run \
      -–log-driver gelf –-log-opt gelf-address=udp://1.2.3.4:12201 \
      alpine echo hello world
```

### GELF Options

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


