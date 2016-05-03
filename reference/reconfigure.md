<!--[metadata]>
+++
title = "reconfigure"
description = "Docker Trusted Registry reconfigure command reference."
keywords = ["docker, registry, reference, reconfigure"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_reconfigure"
+++
<![end-metadata]-->

# docker/dtr reconfigure

Change DTR configurations

# Usage

```bash
$ docker run -it --rm docker/dtr \
  reconfigure [command options]
```

## Description

This command performs changes to DTR configurations.

WARNING: DTR is restarted for the new configurations to take
effect. To have no down time, configure your DTR for high-availability.


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
| `--existing-replica-id`    | ID of an existing replica in a cluster [$DTR_EXISTING_REPLICA_ID]                                |
