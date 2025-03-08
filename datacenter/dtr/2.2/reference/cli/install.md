---
title: docker/dtr install
keywords: docker, dtr, cli, install
description: Install Docker Trusted Registry
---

Install Docker Trusted Registry

## Usage

```bash
docker run -it --rm docker/dtr \
    install [command options]
```

## Description


This command installs Docker Trusted Registry (DTR) on a node managed by
Docker Universal Control Plane (UCP).

After installing DTR, you can join additional DTR replicas using the `join`
command.

Example usage:

    $ docker run -it --rm docker/dtr:2.2.0 install \
	    --ucp-node <UCP_NODE_HOSTNAME> \
	    --ucp-insecure-tls

**Note**: We recommend `--ucp-ca "$(cat ca.pem)"` instead of `--ucp-insecure-tls` for a production deployment.

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug`|Enable debug mode for additional logging|
|`--dtr-ca`|Use a PEM-encoded TLS CA certificate for DTR. If not provided, one will be generated at install time.|
|`--dtr-cert`|Use a PEM-encoded TLS certificate for DTR. If not provided, one will be generated at install time.|
|`--dtr-external-url`|URL of the host or load balancer clients use to reach DTR. Format https://host[:port]|
|`--dtr-key`|Use a PEM-encoded TLS private key for DTR. If not provided, one will be generated at install time.|
|`--dtr-storage-volume`|Full path or volume name to store Docker images in the local filesystem|
|`--enable-pprof`|Enable pprof profiling of the server|
|`--extra-envs`|Environment variables or swarm constraints for DTR containers. Format var=val[&var=val]|
|`--http-proxy`|The HTTP proxy used for outgoing requests|
|`--https-proxy`|The HTTPS proxy used for outgoing requests|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--log-host`|Endpoint to send logs to, required if `--log-protocol` is tcp or udp|
|`--log-level`|Log level for container logs. Default: INFO|
|`--log-protocol`|The protocol for sending container logs: tcp, tcp+tls, udp or internal. Default: internal|
|`--nfs-storage-url`|NFS to store Docker images. Requires NFS client libraries. Format nfs://<ip\|hostname>/<mountpoint>|
|`--no-proxy`|Don't use a proxy for these domains. Format acme.org[, acme.com]|
|`--overlay-subnet`|The subnet used by the dtr-ol overlay network. Example: 10.0.0.0/24|
|`--replica-http-port`|The public HTTP port for the DTR replica. Default is 80|
|`--replica-https-port`|The public HTTPS port for the DTR replica. Default is 443|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-node`|The hostname of the target UCP node. Set to empty string or "_random_" to pick one at random.|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--unsafe`|Allow DTR to be installed on any engine version|

