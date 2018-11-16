---
title: docker/dtr emergency-repair
description: Recover DTR from loss of quorum
keywords: dtr, cli, emergency-repair
---

Recover DTR from loss of quorum

## Usage

```bash
docker run -it --rm docker/dtr \
    emergency-repair [command options]
```

## Description


This command repairs a DTR cluster that has lost quorum by reverting your
cluster to a single DTR replica.

There are three steps you can take to recover an unhealthy DTR cluster:

1. If the majority of replicas are healthy, remove the unhealthy nodes from
   the cluster, and join new ones for high availability.
2. If the majority of replicas are unhealthy, use this command to revert your
   cluster to a single DTR replica.
3. If you can't repair your cluster to a single replica, you'll have to
   restore from an existing backup, using the 'restore' command.

When you run this command, a DTR replica of your choice is repaired and
turned into the only replica in the whole DTR cluster.
The containers for all the other DTR replicas are stopped and removed. When
using the 'force' option, the volumes for these replicas are also deleted.

After repairing the cluster, you should use the 'join' command to add more
DTR replicas for high availability.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica.To add, remove or modify DTR, you must connect to an existing  healthy replica's database.. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--overlay-subnet` | $DTR_OVERLAY_SUBNET | The subnet used by the dtr-ol overlay network. Example: 10.0.0.0/24.For high-availalibity, DTR creates an overlay network between UCP nodes. This flag  allows you to choose the subnet for that network. Make sure the subnet you choose is not  used on any machine where DTR replicas are deployed. |
| `--prune` | $PRUNE | Delete the data volumes of all unhealthy replicas.With this option, the volume of the DTR replica you're restoring    is preserved but the volumes for all other replicas are deleted. This has the same    result as completely uninstalling DTR from those replicas.. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP.Download the UCP TLS CA certificate from https://<ucp-url>/ca, and  use --ucp-ca "$(cat ca.pem)". |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP.The installation uses TLS but always trusts  the TLS certificate used by UCP, which can lead to man-in-the-middle attacks.  For production deployments, use --ucp-ca "$(cat ca.pem)" instead. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |
| `--y, yes` | $YES | Answer yes to any prompts. |

