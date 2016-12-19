---
description: Upgrade a v2.0.0 or later cluster to this version of DTR
keywords: docker, dtr, cli, upgrade
title: docker/dtr upgrade
---

Upgrade a v2.0.0 or later cluster to this version of DTR

## Usage

```bash
docker run -it --rm docker/dtr \
    upgrade [command options]
```

## Description


This command upgrades an existing DTR 2.0.0 or later cluster to the current version of
this bootstrapper.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|Specify the UCP controller URL including domain and port|
|`--ucp-username`|Specify the UCP admin username|
|`--ucp-password`|Specify the UCP admin password|
|`--debug`|Enable debug mode, provides additional logging|
|`--hub-username`|Specify the Docker Hub username for pulling images|
|`--hub-password`|Specify the Docker Hub password for pulling images|
|`--http-proxy`|Set the HTTP proxy for outgoing requests|
|`--https-proxy`|Set the HTTPS proxy for outgoing requests|
|`--no-proxy`|Set the list of domains to not proxy to|
|`--replica-http-port`|Specify the public HTTP port for the DTR replica; 0 means unchanged/default|
|`--replica-https-port`|Specify the public HTTPS port for the DTR replica; 0 means unchanged/default|
|`--log-protocol`|The protocol for sending container logs: tcp, tcp+tls, udp or internal. Default: internal|
|`--log-host`|Endpoint to send logs to, required if --log-protocol is tcp or udp|
|`--log-level`|Log level for container logs. Default: INFO|
|`--log-tls-ca-cert`|PEM-encoded TLS CA cert for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.|
|`--log-tls-cert`|PEM-encoded TLS cert for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.|
|`--log-tls-key`|PEM-encoded TLS key for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.|
|`--log-tls-skip-verify`|Configures DTR logging driver's TLS verification. This verification is enabled by default, but it can be overrided by setting this option to true. This option is ignored if the address protocol is not tcp+tls.|
|`--dtr-external-url`|Specify the external domain name and port for DTR. If using a load balancer, use its external URL instead.|
|`--enable-pprof`|Enables pprof profiling of the server|
|`--etcd-heartbeat-interval`|Set etcd's frequency (ms) that its leader will notify followers that it is still the leader.|
|`--etcd-election-timeout`|Set etcd's timeout (ms) for how long a follower node will go without hearing a heartbeat before attempting to become leader itself.|
|`--etcd-snapshot-count`|Set etcd's number of changes before creating a snapshot.|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--existing-replica-id`|ID of an existing replica in a cluster|
