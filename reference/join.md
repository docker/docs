+++
title = "join"
keywords= ["join, ucp"]
description = "Join this Engine to an existing UCP"
[menu.main]
identifier = "ucp_join"
parent = "ucp_ref"
+++

# join

Join this Engine to an existing UCP

## Usage

```
docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    join [command options]
```

## Description

When running the `join` command, you must run this tool
on the engine you wish to join to an existing UCP.  The UCP controller
must be running on a different engine.  Both engines must have network
visibility to eachother and have the required ports open on their firewall
settings.  If your system has multiple IP addresses, you may need to
specify the `--host-address` option to ensure the correct address is used.

Ports: 443 (customizable using `install --controller-port`),
12376, 12379, 12380, 12381, 12382 and 2376
(customizable using `install --swarm-port`).

To enable high-availability, you must join at least one node with the
`--replica` flag.

## Options

| Option                                                     | Description                                                                                         |
|:-----------------------------------------------------------|:----------------------------------------------------------------------------------------------------|
| `--debug`, `-D`                                            | Enable debug.                                                                                       |
| `--jsonlog`                                                | Produce json formatted output for easier parsing.                                                   |
| `--interactive`, `-i`                                      | Enable interactive mode. You will be prompted to enter all required information.                    |
| `--fresh-install`                                          | Destroy any existing state and start fresh.                                                         |
| `--san` [`--san option` `--san option`]                    | Additional Subject Alternative Names for certs (e.g. `--san foo1.bar.com --san foo2.bar.com`).      |
| `--host-address`                                           | Specify the visible IP/hostname for this node. (override automatic detection) [`$UCP_HOST_ADDRESS`] |
| `--swarm-port "2376"`                                      | Select what port to run the local Swarm manager on.                                                 |
| `--controller-port "443"`                                  | Select what port to run the local Controller on.                                                    |
| `--dns` [`--dns option --dns option`]                      | Set custom DNS servers for the UCP infrastructure containers.                                       |
| `--dns-opt `[`--dns-opt option --dns-opt option`]          | Set DNS options for the UCP infrastructure containers.                                              |
| `--dns-search` [`--dns-search option --dns-search option`] | Set custom DNS search domains for the UCP infrastructure containers.                                |
| `--disable-tracking`                                       | Disable anonymous tracking and analytics.                                                           |
| `--disable-usage`                                          | Disable anonymous usage reporting.                                                                  |
| `--url`                                                    | The connection URL for the remote UCP controller [`$UCP_URL`]                                       |
| `--fingerprint `                                           | The fingerprint of the UCP controller you trust [`$UCP_FINGERPRINT`]                                |
| `--replica`                                                | Configure this node as a full Orca controller replica.                                              |
| `--external-server-cert`                                        | (Replica only) Use externally signed certificates for the controller.                               |
| `--pull "missing"`                                         | Specify image pull behavior (`always`, when `missing`, or `never`).                                 |
