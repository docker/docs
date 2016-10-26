<!--[metadata]>
+++
title ="install"
description="Install Docker Trusted Registry"
keywords= ["docker, dtr, cli, install"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_install"
+++
<![end-metadata]-->

# docker/dtr install

Install Docker Trusted Registry

## Usage

```bash
docker run -it --rm docker/dtr \
    install [command options]
```

## Description


This command installs Docker Trusted Registry (DTR) on a node managed by
Docker Universal Control Plane (UCP).

After installing DTR, you can join additional DTR replicas using the 'join'
command.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--http-proxy`|The HTTP proxy used for outgoing requests|
|`--https-proxy`|The HTTPS proxy used for outgoing requests|
|`--no-proxy`|Don't use a proxy for these domains. Format acme.org[, acme.com]|
|`--replica-http-port`|The public HTTP port for the DTR replica. Default is 80|
|`--replica-https-port`|The public HTTPS port for the DTR replica. Default is 443|
|`--log-protocol`|The protocol for sending container logs: tcp, tcp+tls, udp or internal. Default: internal|
|`--log-host`|Endpoint to send logs to, required if --log-protocol is tcp or udp|
|`--log-level`|Log level for container logs. Default: INFO|
|`--log-tls-ca-cert`|PEM-encoded TLS CA cert for DTR logging driver. Ignored if the logging protocol is not tcp+tls|
|`--log-tls-cert`|PEM-encoded TLS cert for DTR logging driver. Ignored if the logging protocol is not tcp+tls|
|`--log-tls-key`|PEM-encoded TLS key for DTR logging driver. Ignored if the address protocol is not tcp+tls|
|`--log-tls-skip-verify`|Disable TLS verification for the logging service. Ignored if the logging address is not tcp+tls|
|`--dtr-external-url`|URL of the host or load balancer clients use to reach DTR. Format https://host[:port]|
|`--dtr-storage-volume`|Full path or volume name to store Docker images in the local filesystem|
|`--nfs-storage-url`|NFS to store Docker images. Requires NFS client libraries. Format nfs://<ip|hostname>/<mountpoint>|
|`--enable-pprof`|Enables pprof profiling of the server|
|`--etcd-heartbeat-interval`|Frequency in milliseconds that the key-value store leader notifies followers|
|`--etcd-election-timeout`|Timeout in milliseconds for key-value store membership|
|`--etcd-snapshot-count`|Number of changes between key-value store snapshots|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-node`|The hostname of the node to install DTR|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--unsafe`|Allow DTR to be installed on a UCP manager node|
|`--extra-envs`|Environment variables or swarm constraints for DTR containers. Format var=val[&var=val]|

