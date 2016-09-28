<!--[metadata]>
+++
title = "migrate"
description = "Learn about the options available on the docker/dtr image."
keywords = ["docker, dtr, install, uninstall, configure"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_migrate"
+++
<![end-metadata]-->

# docker/dtr migrate

Migrate configurations, accounts, and repository metadata from DTR 1.4.3 to 2.0

## Usage

```bash
$ docker run -it --rm docker/dtr \
    migrate [command options]
```

## Description

This command migrates configurations, accounts, and repository metadata from
Docker Trusted Registry 1.4.3 to 2.0.

To upgrade from DTR 1.4.3 to 2.0 you first need to do a fresh installation of
DTR 2.0. Then you run this command on the host where DTR 1.4.3 is running, to
migrate the data to the new installation. Finally, you decommission your
DTR 1.4.3 by uninstalling it.


| Option                 | Description                                                                          |
|:-----------------------|:-------------------------------------------------------------------------------------|
| `--ucp-url `           | Specify the UCP controller URL [$UCP_URL]                                            |
| `--ucp-username`       | Specify the UCP admin username [$UCP_USERNAME]                                       |
| `--ucp-password`       | Specify the UCP admin password [$UCP_PASSWORD]                                       |
| `--debug`              | Enable debug mode, provides additional logging [$DEBUG]                              |
| `--hub-username`       | Specify the Docker Hub username for pulling images [$HUB_USERNAME]                   |
| `--hub-password`       | Specify the Docker Hub password for pulling images [$HUB_PASSWORD]                   |
| `--ucp-insecure-tls`   | Disable TLS verification for UCP [$UCP_INSECURE_TLS]                                 |
| `--ucp-ca`             | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]                               |
| `--run-full-migration` | Run full migration procedure instead of dumping configurations [$RUN_FULL_MIGRATION] |
| `--dtr-load-balancer`  | Specify the domain name and port for the DTR load balancer [$DTR_LOAD_BALANCER]      |
| `--dtr-insecure-tls`   | Disable TLS verification for DTR [$DTR_INSECURE_TLS]                                 |
| `--dtr-ca`             | PEM-encoded TLS CA cert for DTR [$DTR_CA]                                            |
| `--http-proxy`         | Set the HTTP proxy for outgoing requests [$DTR_HTTP_PROXY]                           |
| `--https-proxy`        | Set the HTTPS proxy for outgoing requests [$DTR_HTTPS_PROXY]                         |
| `--no-proxy`           | Set the list of domains to not proxy to [$DTR_NO_PROXY]                              |
