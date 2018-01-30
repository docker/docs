---
description: Configure general settings for Docker Trusted Registry
keywords: docker, documentation, about, technology, understanding, enterprise, hub, general, domain name, HTTP, HTTPS ports, Notary, registry
redirect_from:
- /docker-trusted-registry/configure/config-general/
title: Configure general settings
---

This document describes the general settings you need to configure including
using Trusted Content through setting up your Notary server.

## Configure your domain name and port settings

Each setting on this page is explained in the Docker Trusted Registry UI.

* *Domain Name*: **required**. By default it is an empty string. It is the fully qualified domain name assigned to the Docker Trusted Registry host.
* *HTTP Port*: defaults to 80 and is used as the entry point for the image storage service. To see load balancer status, you can query
http://&lt;dtr-host&gt;/load_balancer_status.
* *HTTPS Port*: defaults to 443, used as the secure entry point for the image storage service.
* *HTTP proxy*: defaults to an empty string, proxy server for HTTP requests.
* *HTTPS proxy*: defaults to an empty string, proxy server for HTTPS requests.
* *No proxy*: defaults to an empty string, proxy bypass for HTTP and HTTPS requests.
* *Upgrade checking*: enables or disables automatic checking for the Trusted Registry software updates.

If you need the Trusted Registry to re-generate a self-signed certificate at
some point, you can change the domain name. Whenever the domain name does not
match the current certificate, a new self-signed certificate is generated
for the new domain. This also works with IP addresses.


## Docker Content Trust

The Trusted Registry's includes integration with of Docker Notary to provide
Content Trust functionality, allowing your organization to push and pull
trusted images. After pushing images in the Trusted Registry, you can see
which image tags were signed by viewing the appropriate repositories through
Trusted Registry's web interface.

To configure your Docker client to push signed images to Docker
Trusted Registry refer to the CLI Reference's [Environment Variables
Section](/engine/reference/commandline/cli.md#environment-variables) and
[Notary Section](/engine/reference/commandline/cli.md#notary).

This requires you to set the `DOCKER_CONTENT_TRUST` variable and configure
your system to trust Docker Trusted Registry's TLS certificate if it doesn't
already.

## See also

* [Configure storage settings](config-storage.md)