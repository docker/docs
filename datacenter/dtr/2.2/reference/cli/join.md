---
title: docker/dtr join
keywords: docker, dtr, cli, join
description: Add a new replica to an existing DTR cluster
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

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug`|Enable debug mode for additional logging|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--extra-envs`|Environment variables or swarm constraints for DTR containers. Format var=val[&var=val]|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--replica-http-port`|The public HTTP port for the DTR replica. Default is 80|
|`--replica-https-port`|The public HTTPS port for the DTR replica. Default is 443|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--skip-network-test`|Don't test if overlay networks are working correctly between UCP nodes|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-node`|The hostname of the target UCP node. Set to empty string or "_random_" to pick one at random.|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--unsafe`|Allow DTR to be installed on any engine version|
|`--unsafe-join`|Perform the join despite the cluster containing unhealthy replicas.|

