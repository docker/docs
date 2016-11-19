---
description: Learn about the options available on the docker/dtr image.
keywords: docker, dtr, install, uninstall, configure
title: docker/dtr migrate
---

Migrate configurations, accounts, and repository metadata from DTR 1.4.3 to 2.0

## Usage

```bash
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock docker/dtr \
    migrate [command options]
```

## Description

This command migrates configurations, accounts, and repository metadata from
Docker Trusted Registry 1.4.3 to 2.0.

To upgrade from DTR 1.4.3 to 2.0 you first need to do a fresh installation of
DTR 2.0. Then you run this command on the host where DTR 1.4.3 is running, to
migrate the data to the new installation. Finally, you decommission your
DTR 1.4.3 by uninstalling it.


## Options

| Option                 | Description                                                                                                |
|:-----------------------|:-----------------------------------------------------------------------------------------------------------|
| `--ucp-url`            | Specify the UCP controller URL including domain and port                                                   |
| `--ucp-username`       | Specify the UCP admin username                                                                             |
| `--ucp-password`       | Specify the UCP admin password                                                                             |
| `--debug`              | Enable debug mode, provides additional logging                                                             |
| `--hub-username`       | Specify the Docker Hub username for pulling images                                                         |
| `--hub-password`       | Specify the Docker Hub password for pulling images                                                         |
| `--ucp-insecure-tls`   | Disable TLS verification for UCP                                                                           |
| `--ucp-ca`             | Use a PEM-encoded TLS CA certificate for UCP                                                               |
| `--run-full-migration` | Run full migration procedure instead of dumping configurations                                             |
| `--dtr-external-url`   | Specify the external domain name and port for DTR. If using a load balancer, use its external URL instead. |
| `--dtr-insecure-tls`   | Disable TLS verification for DTR                                                                           |
| `--dtr-ca`             | PEM-encoded TLS CA cert for DTR                                                                            |
| `--http-proxy`         | Set the HTTP proxy for outgoing requests                                                                   |
| `--https-proxy`        | Set the HTTPS proxy for outgoing requests                                                                  |
| `--no-proxy`           | Set the list of domains to not proxy to                                                                    |