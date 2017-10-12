---
title: docker/dtr join
description: Add a new replica to an existing DTR cluster
keywords: dtr, cli, join
---

Add a new replica to an existing DTR cluster

## Usage

```bash
docker run -it --rm docker/dtr \
    join [command options]
```

## Description


This command creates a replica of an existing DTR on a node managed by
Docker Universal Control Plane (UCP).

For setting DTR for high-availability, create 3, 5, or 7 replicas of DTR.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--enzi-ca` | $ENZI_TLS_CA | Use a PEM-encoded TLS CA certificate for Enzi. |
| `--enzi-insecure-tls` | $ENZI_TLS_INSECURE | Disable TLS verification for Enzi. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--replica-http-port` | $REPLICA_HTTP_PORT | The public HTTP port for the DTR replica. Default is 80. This allows you to customize the HTTP port where users can reach DTR. Once users access the HTTP port, they are redirected to use an HTTPS connection, using the port specified with --replica-https-port. This port can also be used for unencrypted health checks. |
| `--replica-https-port` | $REPLICA_HTTPS_PORT | The public HTTPS port for the DTR replica. Default is 443. This allows you to customize the HTTPS port where users can reach DTR. Each replica can use a different port. |
| `--replica-id` | $DTR_INSTALL_REPLICA_ID | Assign an ID to the DTR replica. Random by default. |
| `--skip-network-test` | $DTR_SKIP_NETWORK_TEST | Don't test if overlay networks are working correctly between UCP nodes. For high-availability, DTR creates an overlay network between UCP nodes and tests that it is working when joining replicas. Don't use this option for production deployments. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP. Download the UCP TLS CA certificate from https://<ucp-url>/ca, and use --ucp-ca "$(cat ca.pem)". |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP. The installation uses TLS but always trusts the TLS certificate used by UCP, which can lead to man-in-the-middle attacks. For production deployments, use --ucp-ca "$(cat ca.pem)" instead. |
| `--ucp-node` | $UCP_NODE | The hostname of the UCP node to deploy DTR. Random by default. You can find the hostnames of the nodes in the cluster in the UCP web UI, or by running `docker node ls` on a UCP manager node. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |
| `--unsafe-join` | $DTR_UNSAFE_JOIN | Join a new replica even if the cluster is unhealthy. Joining replicas to an unhealthy DTR cluster leads to split-brain scenarios and data loss. Don't use this option for production deployments. |

