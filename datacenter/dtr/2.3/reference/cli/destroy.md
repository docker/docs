---
title: docker/dtr destroy
description: Destroy a DTR replica's data
keywords: dtr, cli, destroy
---

Destroy a DTR replica's data

## Usage

```bash
docker run -it --rm docker/dtr \
    destroy [command options]
```

## Description


This command forcefully removes all containers and volumes associated with
a DTR replica without notifying the rest of the cluster. Use this command
on all replicas uninstall DTR.

Use the 'remove' command to gracefully scale down your DTR cluster.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--replica-id` | $DTR_DESTROY_REPLICA_ID | The ID of the replica to destroy. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP. The installation uses TLS but always trusts the TLS certificate used by UCP, which can lead to man-in-the-middle attacks. For production deployments, use --ucp-ca "$(cat ca.pem)" instead. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP. Download the UCP TLS CA certificate from https://<ucp-url>/ca, and  use --ucp-ca "$(cat ca.pem)". |
| `--enzi-insecure-tls` | $ENZI_TLS_INSECURE | Disable TLS verification for Enzi. |
| `--enzi-ca` | $ENZI_TLS_CA | Use a PEM-encoded TLS CA certificate for Enzi. |

