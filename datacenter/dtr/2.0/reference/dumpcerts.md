---
description: Docker Trusted Registry dumpcerts command reference.
keywords: docker, registry, reference, dumpcerts
title: docker/dtr dumpcerts
---

Dump out the TLS certificates used by this DTR instance

## Usage

```bash
docker run -i --rm docker/dtr \
    dumpcerts [command options] > backup.tar
```

## Description

This command creates a backup of the certificates used by DTR for
communicating across replicas with TLS.


## Options

| Option                  | Description                                              |
|:------------------------|:---------------------------------------------------------|
| `--ucp-url`             | Specify the UCP controller URL including domain and port |
| `--ucp-username`        | Specify the UCP admin username                           |
| `--ucp-password`        | Specify the UCP admin password                           |
| `--debug`               | Enable debug mode, provides additional logging           |
| `--hub-username`        | Specify the Docker Hub username for pulling images       |
| `--hub-password`        | Specify the Docker Hub password for pulling images       |
| `--ucp-insecure-tls`    | Disable TLS verification for UCP                         |
| `--ucp-ca`              | Use a PEM-encoded TLS CA certificate for UCP             |
| `--existing-replica-id` | ID of an existing replica in a cluster                   |