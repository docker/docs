<!--[metadata]>
+++
title = "install"
description = "Docker Trusted Registry install command reference."
keywords = ["docker, registry, reference, install"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_install"
+++
<![end-metadata]-->

# docker/dtr install

Install Docker Trusted Registry on this Docker Engine

## Usage

```bash
$ docker run -it --rm docker/dtr \
    install [command options]
```

## Description

This command installs DTR on the Docker Engine that runs the command.
After installing DTR, you can add more nodes to a DTR cluster with
the 'join' command.

## Options

| Option                     | Description                                                                                      |
|:---------------------------|:-------------------------------------------------------------------------------------------------|
| `--ucp-url`                | Specify the UCP controller URL [$UCP_URL]                                                        |
| `--ucp-username`           | Specify the UCP admin username [$UCP_USERNAME]                                                   |
| `--ucp-password`           | Specify the UCP admin password [$UCP_PASSWORD]                                                   |
| `--debug`                  | Enable debug mode, provides additional logging [$DEBUG]                                          |
| `--hub-username`           | Specify the Docker Hub username for pulling images [$HUB_USERNAME]                               |
| `--hub-password`           | Specify the Docker Hub password for pulling images [$HUB_PASSWORD]                               |
| `--http-proxy`             | Set the HTTP proxy for outgoing requests [$DTR_HTTP_PROXY]                                       |
| `--https-proxy`            | Set the HTTPS proxy for outgoing requests [$DTR_HTTPS_PROXY]                                     |
| `--no-proxy`               | Set the list of domains to not proxy to [$DTR_NO_PROXY]                                          |
| `--replica-http-port "0"`  | Specify the public HTTP port for the DTR replica [$REPLICA_HTTP_PORT]                            |
| `--replica-https-port "0"` | Specify the public HTTPS port for the DTR replica [$REPLICA_HTTPS_PORT]                          |
| `--log-protocol`           | The protocol for sending container logs: tcp, udp or internal. Default: internal [$LOG_PROTOCOL] |
| `--log-host`               | Endpoint to send logs to, required if --log-protocol is tcp or udp [$LOG_HOST]                   |
| `--log-level`              | Log level for container logs. Default: INFO [$LOG_LEVEL]                                         |
| `--dtr-load-balancer`      | Specify the domain name and port for the DTR load balancer [$DTR_LOAD_BALANCER]                  |
| `--ucp-insecure-tls`       | Disable TLS verification for UCP [$UCP_INSECURE_TLS]                                             |
| `--ucp-ca`                 | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]                                           |
| `--ucp-node`               | Specify the host to install Docker Trusted Registry [$UCP_NODE]                                  |
| `--replica-id`             | Specify the replica Id. Must be unique per replica, leave blank for random [$DTR_REPLICA_ID]     |
