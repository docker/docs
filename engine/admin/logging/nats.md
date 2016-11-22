---
description: Describes how to use NATS for publishing log entries
keywords: NATS, nats.io, messaging, docker, logging, driver
title: NATS logging driver
---

Docker logging driver for sending container the logs as events published to NATS in JSON format.

## Usage

You can configure the default logging driver by passing the `--log-driver`
option to the Docker daemon:

```bash
$ dockerd --log-driver=nats
```

You can set the logging driver for a specific container by using the
`--log-driver` option to `docker run`:

```bash
$ docker run --log-driver=nats ...
```

This log driver does not implement a reader so it is incompatible with `docker logs`.

## nats options

You can use the `--log-opt NAME=VALUE` flag to customize the logging driver
for NATS:

| Option                      | Required | Description                                                                                                                                 |
|-----------------------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------|
| `labels`                    | optional | Comma-separated list of keys of labels, which should be included in message, if these labels are specified for container.                   |
| `env`                       | optional | Comma-separated list of keys of environment variables, which should be included in message, if these variables are specified for container. |
| `tag`                       | optional | Specify tag for message.  Refer to the [log tag option documentation](log_tags.md) for customizing the log tag format.                      |
| `nats-servers`              | optional | NATS cluster nodes separated by commas. e.g. `nats://127.0.0.1:4222,nats://127.0.0.1:4223`. Defaults to `localhost:4222`                    |
| `nats-max-reconnect`        | optional | Maximum attempts that the driver will try to connect before giving up. Defaults to infinite (`-1`)                                          |
| `nats-subject`              | optional | Specific subject to which logs will be published. Defaults to using `tag` if not specified                                                  |
| `nats-user`                 | optional | Specify user in case of authentication required                                                                                             |
| `nats-pass`                 | optional | Specify password in case of authentication required                                                                                         |
| `nats-token`                 | optional | Specify token in case of authentication required                                                                                           |
| `nats-tls-ca-cert`          | optional | Specified the absolute path to the trust certificates signed by the CA                                                                      |
| `nats-tls-cert`             | optional | Specifies the absolute path to the TLS certificate file                                                                                     |
| `nats-tls-key`              | optional | Specifies the absolute path to the TLS key file                                                                                             |
| `nats-tls-skip-verify`      | optional | Specifies whether to skip verification by setting it to `true`                                                                              |

Below is an example usage of the driver for sending logs to a node in a
NATS cluster to the `docker.logs` subject:

```bash
$ docker run --log-driver=nats \
             --log-opt nats-subject=docker.logs \
             --log-opt nats-servers=nats://nats-node-1:4222,nats://nats-node-2:4222,nats://nats-node-3:4222 \
             your/application
```

By default, the tag is used as the subject for NATS, so it has to be a valid
subject in case subject it is left unspecified:

```bash
{% raw %}
$ docker run --log-driver nats \
             --log-opt tag="docker.{{.ID}}.{{.ImageName}}"
             your/application
{% endraw %}
```

Secure connection to NATS using TLS can be customized by setting `tls://` scheme
in the URI and absolute paths to the certs and key files:

```bash
docker run --log-driver nats \
           --log-opt nats-tls-key=/srv/configs/certs/client-key.pem \
           --log-opt nats-tls-cert=/srv/configs/certs/client-cert.pem \
           --log-opt nats-tls-ca-cert=/srv/configs/certs/ca.pem \
           --log-opt nats-servers="tls://127.0.0.1:4223,tls://127.0.0.1:4222" \
           your/application
```

Skip verify is enabled by default, in order to deactivate we can specify `nats-tls-skip-verify`:

```bash
  docker run --log-driver nats \
             --log-opt nats-tls-skip-verify \
             --log-opt nats-servers="tls://127.0.0.1:4223,tls://127.0.0.1:4222" \
             your/application
```
