---
description: Joins a node to an existing Docker Universal Control Plane cluster.
keywords: join, ucp
title: docker/ucp join
---

Join this engine to an existing UCP.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  join [command options]
```

## Description

When running the 'join' command, you must run this tool
on the engine you wish to join to an existing UCP.  The UCP controller
must be running on a different engine.  Both engines must have network
visibility to each other and have the required ports open on their firewall
settings.  If your system has multiple IP addresses, you may need to
specify the '--host-address' option to ensure the correct address is used.

Ports: 443, 12376, 12379, 12380, 12381, 12382 and 2376 or the '--swarm-port' from 'install'

To enable high-availability, you must join at least one node with the
'replica' flag.  When joining a replica, you may optionally provide
a backup tar file containing the Root CA material from the initial
controller saved via the 'backup' command.  Only the Root CA material
will be retrieved from that backup. Add the following volume mount:
`-v /path/to/my/backup.tar:/backup.tar`

## Options

| Option                                                     | Description                                                                                       |
|:-----------------------------------------------------------|:--------------------------------------------------------------------------------------------------|
| `--debug, -D`                                              | Enable debug.                                                                                     |
| `--jsonlog`                                                | Produce json formatted output for easier parsing.                                                 |
| `--interactive, -i`                                        | Enable interactive mode. You will be prompted to enter all required information.                  |
| `--admin-username`                                         | Specify the UCP admin username [$UCP_ADMIN_USER]                                                  |
| `--admin-password`                                         | Specify the UCP admin password [$UCP_ADMIN_PASSWORD]                                              |
| `--fresh-install`                                          | Destroy any existing state and start fresh.                                                       |
| `--san` `[--san option --san option]`                      | Additional Subject Alternative Names for certs. For example, `--san foo1.bar.com --san foo2.bar.com`.    |
| `--host-address`                                           | Specify the visible IP/hostname for this node. (override automatic detection) [$UCP_HOST_ADDRESS] |
| `--swarm-port "2376"`                                      | Select what port to run the local Swarm manager on (default: 2376)                                |
| `--controller-port "443"`                                  | Select what port to run the local Controller on (default: 443)                                    |
| `--dns` `[--dns option --dns option]`                      | Set custom DNS servers for the UCP infrastructure containers.                                     |
| `--dns-opt` `[--dns-opt option --dns-opt option]`          | Set DNS options for the UCP infrastructure containers.                                            |
| `--dns-search` `[--dns-search option --dns-search option]` | Set custom DNS search domains for the UCP infrastructure containers.                              |
| `--registry-username`                                      | Specify the username to pull required images with [$REGISTRY_USERNAME]                            |
| `--registry-password`                                      | Specify the password to pull required images with [$REGISTRY_PASSWORD]                            |
| `--url`                                                    | The connection URL for the remote UCP controller [$UCP_URL]                                       |
| `--fingerprint`                                            | The fingerprint of the UCP controller you trust [$UCP_FINGERPRINT]                                |
| `--replica`                                                | Configure this node to be a UCP controller replica.                                               |
| `--external-server-cert`                                   | (Replica only) Use externally signed certificates for the controller.                             |
| `--pull "missing"`                                         | Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing")           |
| `--passphrase`                                             | Decrypt the Root CA tar file with the provided passphrase [$UCP_PASSPHRASE]                       |
| `--skip-engine-discovery`                                  | Do not configure engine for clustering                                                            |