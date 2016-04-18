+++
title = "install"
keywords= ["install, ucp"]
description = "Install UCP controller"
[menu.main]
identifier = "ucp_install"
parent = "ucp_ref"
+++

# install

Install UCP on this engine.

```bash
docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    install [OPTIONS]
```

## Description

Install the UCP controller on a machine. You can only install on machines where
Docker Engine is already installed. If you intend to install a multi-node
cluster, you must open firewall ports between the Engines for the following
ports:

* 443 (customizable using the `--controller-port` option)
* 12376
* 12379 through 12382
* 2376 (customizable using the `--swarm-port` option)

You can optionally use an externally generated and signed certificate for the
UCP controller by using the `--external-server-cert`. Create a storage volume named
`ucp-controller-server-certs` with ca.pem, cert.pem, and key.pem in the root directory
before running the install.

A license file can optionally be added during install by volume
mounting the file at `/docker_subscription.lic` in the tool.

```bash
-v /path/to/my/docker_subscription.lic:/docker_subscription.lic
```

## Options

| Option                                                     | Description                                                                                    |
|:-----------------------------------------------------------|:-----------------------------------------------------------------------------------------------|
| `--debug`, `-D`                                            | Enable debug.                                                                                  |
| `--jsonlog`                                                | Produce json formatted output for easier parsing.                                              |
| `--interactive`, `-i`                                      | Enable interactive mode.,You will be prompted to enter all required information.               |
| `--fresh-install`                                          | Destroy any existing state and start fresh.                                                    |
| `--san` `[--san option --san option]`                      | Additional Subject Alternative Names for certs (e.g. `--san foo1.bar.com --san foo2.bar.com`). |
| `--host-address`                                           | Specify the visible IP/hostname for this node.                                                 |
| `--swarm-port "2376"`                                      | Select what port to run the local Swarm manager on.                                            |
| `--controller-port "443"`                                  | Select what port to run the local Controller on.                                               |
| `--dns` `[--dns option --dns option]`                      | Set custom DNS servers for the UCP infrastructure containers.                                  |
| `--dns-opt` `[--dns-opt option --dns-opt option]`          | Set DNS options for the UCP infrastructure containers.                                         |
| `--dns-search` `[--dns-search option --dns-search option]` | Set custom DNS search domains for the UCP infrastructure containers.                           |
| `--disable-tracking`                                       | Disable anonymous tracking and analytics.                                                      |
| `--disable-usage`                                          | Disable anonymous usage reporting.                                                             |
| `--external-server-cert`                                        | Set up UCP with an external CA.                                                                |
| `--preserve-certs`                                         | Don't (re)generate certs on the host if existing ones are found.                               |
| `--binpack`                                                | Set Swarm scheduler to binpack mode (default spread).                                          |
| `--random`                                                 | Set Swarm scheduler to random mode (default spread).                                           |
| `--pull "missing"`                                         | Specify image pull behavior (`always`, when `missing`, or `never`).                            
| `--swarm-experimental`                                    | Enable experimental Swarm features. Note: Use only for install, not join).                      |
