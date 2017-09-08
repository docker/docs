---
title: docker/dtr dumpcerts
description: Print the TLS certificates used by DTR
keywords: dtr, cli, certificate, cert, tls
---

 Print the TLS certificates used by DTR

## Usage

```bash
 docker run -i --rm docker/dtr \
    dumpcerts [command options] > backup.tar
```

## Description

This command creates a backup of the certificates used by DTR for
communicating across replicas with TLS.

## Options

|         Option          | Environment Variable |                   Description                   |
| :---------------------- | :------------------- | :---------------------------------------------- |
| `--debug`               | $DEBUG               | Enable debug mode for additional logs.          |
| `--existing-replica-id` | $DTR_REPLICA_ID      | The ID of an existing DTR replica.              |
| `--help-extended`       | $DTR_EXTENDED_HELP   | Display extended help text for a given command. |
| `--ucp-ca`              | $UCP_CA              | Use a PEM-encoded TLS CA certificate for UCP.   |
| `--ucp-insecure-tls`    | $UCP_INSECURE_TLS    | Disable TLS verification for UCP.               |
| `--ucp-password`        | $UCP_PASSWORD        | The UCP administrator password.                 |
| `--ucp-url`             | $UCP_URL             | The UCP URL including domain and port.          |
| `--ucp-username`        | $UCP_USERNAME        | The UCP administrator username.                 |





