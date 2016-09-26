+++
title = "install"
description = "Install Docker Universal Control Plane."
keywords= ["install, ucp"]
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_install"
+++

# docker/ucp install

Install UCP on this engine.

## Usage

```
docker run --rm -it \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           install [command options]
```

## Description

The `install` command will install the UCP controller on the
local engine. If you intend to install a multi-node cluster,
you must open firewall ports between the engines for the
following ports:

- 443 or the `--controller-port`
- 2376 or the `--swarm-port`
- 12376, 4789, 12379, 12380, 12381, 12382, 12383, 12384, 12385, 12386
- 4789(udp) and 7946(tcp/udp) for overlay networking

You can optionally use an externally generated and signed certificate
for the UCP controller by specifying `--external-server-cert`. Create a storage
volume named `ucp-controller-server-certs` with ca.pem, cert.pem, and key.pem
in the root directory before running the install.

A license file can optionally be injected during install by volume
mounting the file at '/docker_subscription.lic' in the tool.  E.g.,
`-v /path/to/my/docker_subscription.lic:/docker_subscription.lic`

## Options

```nohighlight
--debug, -D                Enable debug mode
--jsonlog                  Produce json formatted output for easier parsing
--interactive, -i          Enable interactive mode.  You will be prompted to enter all required information
--admin-username value     Specify the UCP admin username [$UCP_ADMIN_USER]
--admin-password value     Specify the UCP admin password [$UCP_ADMIN_PASSWORD]
--san value                Additional Subject Alternative Names for certs (e.g. --san foo1.bar.com --san foo2.bar.com)
--host-address value       Specify the primary IP address or network interface name for this node to advertise to other members in the cluster (required if multiple interfaces present) [$UCP_HOST_ADDRESS]
--swarm-port value         Select what port to run the local Swarm manager on (default: 2376)
--controller-port value    Select what port to run the local UCP Controller on (default: 443)
--swarm-grpc-port value    Select what port to run Swarm GRPC on (default: 2377)
--dns value                Set custom DNS servers for the UCP infrastructure containers [$DNS]
--dns-opt value            Set DNS options for the UCP infrastructure containers [$DNS_OPT]
--dns-search value         Set custom DNS search domains for the UCP infrastructure containers [$DNS_SEARCH]
--pull value               Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing")
--registry-username value  Specify the username to pull required images with [$REGISTRY_USERNAME]
--registry-password value  Specify the password to pull required images with [$REGISTRY_PASSWORD]
--kv-timeout value         Timeout in milliseconds for the KV store (set higher for a multi-datacenter cluster, all controllers must use the same value) (default: 5000) [$KV_TIMEOUT]
--kv-snapshot-count value  Number of changes between KV store snapshots (all controllers must use the same value) (default: 20000) [$KV_SNAPSHOT_COUNT]
--swarm-experimental       Enable experimental swarm features
--disable-tracking         Disable anonymous tracking and analytics
--disable-usage            Disable anonymous usage reporting
--external-server-cert     Set up UCP Controller with an externally signed server certificate
--preserve-certs           Don't (re)generate certs on the host if existing ones are found
--binpack                  Set Swarm scheduler to binpack mode (default spread)
--random                   Set Swarm scheduler to random mode (default spread)
```
