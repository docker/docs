+++
title = "join"
keywords= ["join, ucp"]
description = "Joins a node to an existing Docker Universal Control Plane cluster."
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_join"
+++

# docker/ucp join

The join command is no longer used.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  join [command options]
```

## Description

To join a node to UCP, simply run `docker swarm join`

## Options

| Option                                                     | Description                                                                                       |
|:-----------------------------------------------------------|:--------------------------------------------------------------------------------------------------|
| `--debug, -D`                                              | Enable debug.                                                                                     |
| `--jsonlog`                                                | Produce json formatted output for easier parsing.                                                 |
| `--interactive, -i`                                        | Enable interactive mode. You will be prompted to enter all required information.                  |
| `--admin-username`                                         | Specify the UCP admin username [$UCP_ADMIN_USER]                                                  |
| `--admin-password`                                         | Specify the UCP admin password [$UCP_ADMIN_PASSWORD]                                              |
| `--fresh-install`                                          | Destroy any existing state and start fresh.                                                       |
| `--san` `[--san option --san option]`                      | Additional Subject Alternative Names for certs (e.g. `--san foo1.bar.com --san foo2.bar.com`).    |
| `--host-address`                                           | Specify the visible IP/hostname for this node. (override automatic detection) [$UCP_HOST_ADDRESS] |
| `--swarm-port "2376"`                                      | Select what port to run the local Swarm manager on (default: 2376)                                |
| `--controller-port "443"`                                  | Select what port to run the local Controller on (default: 443)                                    |
| `--swarm-grpc-port value`                                  | Select what port to run Swarm GRPC on (default: 2377)                                             |
| `--dns` `[--dns option --dns option]`                      | Set custom DNS servers for the UCP infrastructure containers. [$DNS]                              |
| `--dns-opt` `[--dns-opt option --dns-opt option]`          | Set DNS options for the UCP infrastructure containers. [$DNS_OPT]                                 |
| `--dns-search` `[--dns-search option --dns-search option]` | Set custom DNS search domains for the UCP infrastructure containers. [$DNS_SEARCH]                |
| `--registry-username`                                      | Specify the username to pull required images with [$REGISTRY_USERNAME]                            |
| `--registry-password`                                      | Specify the password to pull required images with [$REGISTRY_PASSWORD]                            |
