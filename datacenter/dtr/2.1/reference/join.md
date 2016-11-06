<!--[metadata]>
+++
title ="join"
description="Add a new replica to an existing DTR cluster"
keywords= ["docker, dtr, cli, join"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_join"
+++
<![end-metadata]-->

# docker/dtr join

Add a new replica to an existing DTR cluster



## Description


This command creates a replica of an existing DTR on a node managed by
Docker Universal Control Plane (UCP).

For setting DTR for high-availability, create 3, 5, or 7 replicas of DTR.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-node`|The hostname of the node to install DTR|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--unsafe`|Allow DTR to be installed on a UCP manager node|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--replica-http-port`|The public HTTP port for the DTR replica. Default is 80|
|`--replica-https-port`|The public HTTPS port for the DTR replica. Default is 443|
|`--skip-network-test`|Don't test if overlay networks are working correctly between UCP nodes|
|`--extra-envs`|Environment variables or swarm constraints for DTR containers. Format var=val[&var=val]|

