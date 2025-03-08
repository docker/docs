---
title: docker/dtr reconfigure
keywords: docker, dtr, cli, reconfigure
description: Change DTR configurations
---

Change DTR configurations

## Usage

```bash
docker run -it --rm docker/dtr \
    reconfigure [command options]
```

## Description


This command changes DTR configuration settings.

DTR is restarted for the new configurations to take effect. To have no down
time, configure your DTR for high-availability.


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
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--http-proxy`|The HTTP proxy used for outgoing requests|
|`--https-proxy`|The HTTPS proxy used for outgoing requests|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--log-host`|Endpoint to send logs to, required if `--log-protocol` is tcp or udp|
|`--log-level`|Log level for container logs. Default: INFO|
|`--log-protocol`|The protocol for sending container logs: tcp, tcp+tls, udp or internal. Default: internal|
|`--nfs-storage-url`|NFS to store Docker images. Requires NFS client libraries. Format nfs://<ip\|hostname>/<mountpoint>|
|`--no-proxy`|Don't use a proxy for these domains. Format acme.org[, acme.com]|
|`--replica-http-port`|The public HTTP port for the DTR replica. Default is 80|
|`--replica-https-port`|The public HTTPS port for the DTR replica. Default is 443|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|

