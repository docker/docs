---
description: Install Docker Universal Control Plane.
keywords: install, ucp
title: docker/ucp install
---

Install UCP on this Docker Engine.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  install [command options]
```

## Description

The 'install' command will install the UCP controller on the
local engine. If you intend to install a multi-node cluster,
you must open firewall ports between the engines for the
following ports:

Ports: 443, 12376, 12379, 12380, 12381, 12382 and 2376 or the '--swarm-port'

You can optionally use an externally generated and signed certificate
for the UCP controller by specifying '--external-server-cert'.  Create a storage
volume named 'ucp-controller-server-certs' with ca.pem, cert.pem, and key.pem
in the root directory before running the install.

You can inject a license file during install by mounting the file as a volume
at '/docker_subscription.lic' in the tool. For example,
`-v /path/to/my/docker_subscription.lic:/docker_subscription.lic`

## Options

| Option                                                     | Description                                                                                             |
|:-----------------------------------------------------------|:--------------------------------------------------------------------------------------------------------|
| `--debug`, `-D`                                            | Enable debug.                                                                                           |
| `--jsonlog`                                                | Produce json formatted output for easier parsing.                                                       |
| `--interactive`, `-i`                                      | Enable interactive mode.,You will be prompted to enter all required information.                        |
| `--admin-username`                                         | Specify the UCP admin username [$UCP_ADMIN_USER]                                                        |
| `--admin-password`                                         | Specify the UCP admin password [$UCP_ADMIN_PASSWORD]                                                    |
| `--fresh-install`                                          | Destroy any existing state and start fresh.                                                             |
| `--san` `[--san option --san option]`                      | Additional Subject Alternative Names for certs. For example, `--san foo1.bar.com --san foo2.bar.com`.          |
| `--host-address`                                           | Specify the visible IP for this node.                                                                   |
| `--swarm-port "2376"`                                      | Select what port to run the local Swarm manager on (default: 2376)                                      |
| `--controller-port "443"`                                  | Select what port to run the local Controller on (default: 443)                                          |
| `--dns` `[--dns option --dns option]`                      | Set custom DNS servers for the UCP infrastructure containers.                                           |
| `--dns-opt` `[--dns-opt option --dns-opt option]`          | Set DNS options for the UCP infrastructure containers.                                                  |
| `--dns-search` `[--dns-search option --dns-search option]` | Set custom DNS search domains for the UCP infrastructure containers.                                    |
| `--kv-timeout`                                             | Timeout in milliseconds for the KV store (set higher for a multi-datacenter cluster)                    |
| `--kv-snapshot-count`                                      | Number of changes between KV store snapshots (all controllers must use the same value) (default: 10000) |
| `--registry-username`                                      | Specify the username to pull required images with [$REGISTRY_USERNAME]                                  |
| `--registry-password`                                      | Specify the password to pull required images with [$REGISTRY_PASSWORD]                                  |
| `--swarm-experimental`                                     | Enable experimental Swarm features. Note: Use only for install, not join).                              |
| `--disable-tracking`                                       | Disable anonymous tracking and analytics.                                                               |
| `--disable-usage`                                          | Disable anonymous usage reporting.                                                                      |
| `--external-server-cert`                                   | Set up UCP with an external CA.                                                                         |
| `--preserve-certs`                                         | Don't (re)generate certs on the host if existing ones are found.                                        |
| `--binpack`                                                | Set Swarm scheduler to binpack mode (default spread).                                                   |
| `--random`                                                 | Set Swarm scheduler to random mode (default spread).                                                    |
| `--pull "missing"`                                         | Specify image pull behavior (`always`, when `missing`, or `never`) (default: "missing")                 |
| `--skip-engine-discovery`                                  | Do not configure engine for clustering                                                                  |