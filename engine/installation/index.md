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

## Docker editions

Docker is available in two editions:

* **Stable** gives you reliable updates every quarter
* **Edge** gives you new features every month

For more information about Docker CE, see
[Docker Community Edition](https://www.docker.com/community-edition/){: target="_blank" class="_" }.

  Docker CE has both **stable** and **edge** channels.

  - **Stable** builds are released once per quarter and are supported for 4
    months.
  - **Edge** builds are released once per month, and are supported for that
    month only. If you subscribe to the Edge channel on Linux distributions, you
    should also subscribe to the Stable channel.

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

We also provide official Docker solutions for running on AWS and Azure. You can read
up on what Docker for AWS and Docker for Azure have to offer you
[here](/docker-for-aws/why/) and [here](/docker-for-azure/why/) respectively.

* [Docker for AWS](/docker-for-aws/)
* [Docker for Azure](/docker-for-azure/)

## Get started

After setting up Docker, try learning the basics over at
[Getting started with Docker](/get-started/).
