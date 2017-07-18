---
title: docker/ucp install
description: Install UCP on this node
keywords: docker, ucp, cli, install
---

Install UCP on this node

## Description

This command initializes a new swarm, turns this node into a manager, and installs
Docker Universal Control Plane (UCP).

When installing UCP you can customize:

  * The certificates used by the UCP web server. Create a volume
    named `ucp-controller-server-certs` and copy the ca.pem, cert.pem, and key.pem
    files to the root directory. Then run the install command with the
    `--external-server-cert` flag.

  * The license used by UCP, by bind-mounting the file at
    `/config/docker_subscription.lic` in the tool. E.g. -v /path/to/my/config/docker_subscription.lic:/config/docker_subscription.lic
    or by specifying with `--license "$(cat license.lic)"`

If you're joining more nodes to this swarm, open the following ports in your
firewall:

  * 443 or the `--controller-port`
  * 2376 or the `--swarm-port`
  * 12376, 12379, 12380, 12381, 12382, 12383, 12384, 12385, 12386, 12387
  * 4789 (udp) and 7946 (tcp/udp) for overlay networking


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--interactive, i`|Run in interactive mode and prompt for configuration values|
|`--admin-username`|The UCP administrator username|
|`--admin-password`|The UCP administrator password|
|`--san`|Add subject alternative names to certificates (e.g. --san www1.acme.com --san www2.acme.com)|
|`--host-address`|The network address to advertise to other nodes. Format: IP address or network interface name|
|`--swarm-port`|Port for the Docker Swarm manager. Used for backwards compatibility|
|`--controller-port`|Port for the web UI and API|
|`--swarm-grpc-port`|Port for communication between nodes|
|`--dns`|Set custom DNS servers for the UCP containers|
|`--dns-opt`|Set DNS options for the UCP containers|
|`--dns-search`|Set custom DNS search domains for the UCP containers|
|`--unlock-key`|The unlock key for this swarm-mode cluster, if one exists.|
|`--pull`|Pull UCP images: `always`, when `missing`, or `never`|
|`--registry-username`|Username to use when pulling images|
|`--registry-password`|Password to use when pulling images|
|`--kv-timeout`|Timeout in milliseconds for the key-value store|
|`--kv-snapshot-count`|Number of changes between key-value store snapshots|
|`--swarm-experimental`|Enable Docker Swarm experimental features. Used for backwards compatibility|
|`--disable-tracking`|Disable anonymous tracking and analytics|
|`--disable-usage`|Disable anonymous usage reporting|
|`--external-server-cert`|Customize the certificates used by the UCP web server|
|`--preserve-certs`|Don't generate certificates if they already exist|
|`--binpack`|Set the Docker Swarm scheduler to binpack mode. Used for backwards compatibility|
|`--random`|Set the Docker Swarm scheduler to random mode. Used for backwards compatibility|
|`--external-service-lb`|Set the external service load balancer reported in the UI|
|`--enable-profiling`|Enable performance profiling|
|`--license`|Add a license: e.g. `--license "$(cat license.lic)"`|
