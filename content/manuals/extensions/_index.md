---
title: Docker Extensions
weight: 60
description: Add debugging, testing, security, and networking tools to Docker Desktop.
keywords: Docker Extensions, Docker Desktop, Linux, Mac, Windows
params:
  sidebar:
    group: Application development
aliases:
 - /desktop/extensions/
grid:
  - title: "Extensions Marketplace"
    description: "Find and install extensions in Docker Desktop."
    link: /extensions/marketplace/
    icon: search
  - title: "Extensions SDK"
    description: "Build and distribute a Docker Desktop extension."
    link: /extensions/extensions-sdk/
    icon: code
  - title: "Private marketplace"
    description: "Control which extensions your organization can install."
    link: /extensions/private-marketplace/
    icon: settings
---

Docker Extensions let you use third-party tools within Docker Desktop to extend its functionality.

You can seamlessly connect your favorite development tools to your application development and deployment workflows. Augment Docker Desktop with debugging, testing, security, and networking functionalities, and create custom add-ons using the Extensions [SDK](extensions-sdk/_index.md).

Anyone can use Docker Extensions and there is no limit to the number of extensions you can install.

## Explore the documentation

{{< grid >}}

## What extensions are available?

There is a mix of partner and community-built extensions and Docker-built extensions.
You can explore the list of available extensions in [Docker Hub](https://hub.docker.com/search?q=&type=extension) or in the Extensions Marketplace within Docker Desktop.

## Security and trust

Docker Extensions run with elevated privileges on your host machine. They have direct access to the Docker Engine, can read and write files on your filesystem, and can install and run native binaries. 

Docker reviews extensions submitted to the Marketplace, but does not guarantee the security of any extension. Extensions installed outside the Marketplace have not been reviewed at all. Only install extensions from publishers you trust. 

If you're an organization admin, see [Configure a private marketplace](private-marketplace.md) to control which extensions your team can install.
