---
description: Lists the installation methods
keywords: docker, installation, install, docker ce, docker ee, docker editions, stable, edge
redirect_from:
- /installation/
- /engine/installation/linux/
- /engine/installation/linux/frugalware/
- /engine/installation/frugalware/
title: Install Docker
---

Docker is available in two editions: **Community Edition (CE)** and **Enterprise
Edition (EE)**.

Docker Community Edition (CE) is ideal for developers and small
teams looking to get started with Docker and experimenting with container-based
apps. Docker CE has two update channels, **stable** and **edge**:

* **Stable** gives you reliable updates every quarter
* **Edge** gives you new features every month

For more information about Docker CE, see
[Docker Community Edition](https://www.docker.com/community-edition/){: target="_blank" class="_" }.

Docker Enterprise Edition (EE) is designed for enterprise
development and IT teams who build, ship, and run business critical
applications in production at scale. For more information about Docker EE,
including purchasing options, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

{% include docker_ce_ee.md %}

## Supported platforms

Docker CE and EE are available on multiple platforms, on cloud and on-premises.
Use the following matrix to choose the best installation path for you.

{% include docker_platform_matrix.md %}

See also [Docker Cloud](#on-docker-cloud) for setup instructions for
Digital Ocean, Packet, SoftLink, or Bring Your Own Cloud.

## Time-based release schedule

Starting with Docker 17.03, Docker uses a time-based release schedule, outlined
below.

{% include docker_schedule_matrix.md %}

### Prior releases

Instructions for installing prior releases of Docker can be found in the
[Docker archives](/docsarchive/).

## Docker Cloud

You can use Docker Cloud to automatically provision and manage your cloud instances.

* [Amazon Web Services setup guide](/docker-cloud/infrastructure/link-aws.md)
* [DigitalOcean setup guide](/docker-cloud/infrastructure/link-do.md)
* [Microsoft Azure setup guide](/docker-cloud/infrastructure/link-azure.md)
* [Packet setup guide](/docker-cloud/infrastructure/link-packet.md)
* [SoftLayer setup guide](/docker-cloud/infrastructure/link-softlayer.md)
* [Use the Docker Cloud Agent to Bring your Own Host](/docker-cloud/infrastructure/byoh.md)

## Get started

After setting up Docker, try learning the basics over at
[Getting started with Docker](/get-started/).
