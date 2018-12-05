---
title: docker/dtr upgrade
description: Upgrade DTR 2.4.x cluster to this version
keywords: dtr, cli, upgrade
---

Upgrade DTR 2.4.x cluster to this version

## Usage

```bash
docker run -it --rm docker/dtr \
    upgrade [command options]
```

## Description


This command upgrades DTR 2.4.x to the current version of this image.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica.To add, remove or modify DTR, you must connect to an existing  healthy replica's database.. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP.Download the UCP TLS CA certificate from https://<ucp-url>/ca, and  use --ucp-ca "$(cat ca.pem)". |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP.The installation uses TLS but always trusts  the TLS certificate used by UCP, which can lead to man-in-the-middle attacks.  For production deployments, use --ucp-ca "$(cat ca.pem)" instead. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |

