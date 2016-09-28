<!--[metadata]>
+++
title = "dumpcerts"
[menu.main]
description = "Docker Trusted Registry dumpcerts command reference."
keywords = ["docker, registry, reference, dumpcerts"]
parent="dtr_menu_reference"
identifier="dtr_reference_dumpcerts"
+++
<![end-metadata]-->

# docker/dtr dumpcerts

Dump out the TLS certificates used by this DTR instance

## Usage

```bash
$ docker run -it --rm docker/dtr \
    dumpcerts [command options] > backup.tar
```

## Description

This command creates a backup of the certificates used by DTR for
communicating across replicas with TLS.


| Option                  | Description                                                        |
|:------------------------|:-------------------------------------------------------------------|
| `--ucp-url `            | Specify the UCP controller URL [$UCP_URL]                          |
| `--ucp-username `       | Specify the UCP admin username [$UCP_USERNAME]                     |
| `--ucp-password`        | Specify the UCP admin password [$UCP_PASSWORD]                     |
| `--debug`               | Enable debug mode, provides additional logging [$DEBUG]            |
| `--hub-username`        | Specify the Docker Hub username for pulling images [$HUB_USERNAME] |
| `--hub-password `       | Specify the Docker Hub password for pulling images [$HUB_PASSWORD] |
| `--ucp-insecure-tls`    | Disable TLS verification for UCP [$UCP_INSECURE_TLS]               |
| `--ucp-ca`              | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]             |
| `--existing-replica-id` | ID of an existing replica in a cluster [$DTR_EXISTING_REPLICA_ID]  |
