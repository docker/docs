---
description: Docker Desktop overview
keywords: Desktop, Docker, GUI, run, docker, local, machine
title: Docker Desktop overview
toc_min: 1
toc_max: 2
redirect_from:
- /desktop/opensource/
- /docker-for-mac/opensource/
- /docker-for-windows/opensource/
---

> **Update to the Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) now requires a paid
> subscription. The grace period for those that will require a paid subscription
> ends on January 31, 2022. [Learn more](https://www.docker.com/blog/the-grace-period-for-the-docker-subscription-service-agreement-ends-soon-heres-what-you-need-to-know/){:
 target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

Docker Desktop is an easy-to-install application for your Mac or Windows environment
that enables you to build and share containerized applications and microservices.
Docker Desktop includes [Docker Engine](../engine/index.md), Docker CLI client,
[Docker Compose](../compose/index.md), [Docker Content Trust](../engine/security/trust/index.md),
[Kubernetes](https://github.com/kubernetes/kubernetes/), and 
[Credential Helper](https://github.com/docker/docker-credential-helpers/).

Docker Desktop works with your choice of development tools and languages and
gives you access to a vast library of certified images and templates in
[Docker Hub](https://hub.docker.com/). This enables development teams to extend
their environment to rapidly auto-build, continuously integrate, and collaborate
using a secure repository.

Some of the key features of Docker Desktop include:

* Ability to containerize and share any application on any cloud platform, in multiple languages and frameworks
* Easy installation and setup of a complete Docker development environment
* Includes the latest version of Kubernetes
* Automatic updates to keep you up to date and secure
* On Windows, the ability to toggle between Linux and Windows Server environments to build applications
* Fast and reliable performance with native Windows Hyper-V virtualization
* Ability to work natively on Linux through WSL 2 on Windows machines
* Volume mounting for code and data, including file change notifications and easy access to running containers on the localhost network
* In-container development and debugging with supported IDEs

## Download and install

Docker Desktop is available for Mac and Windows. For download information, system requirements, and installation instructions, see:

* [Install Docker Desktop on Mac](mac/install.md)
* [Install Docker Desktop on Windows](windows/install.md)

For information about Docker Desktop licensing, see [Docker Desktop License Agreement](../subscription/index.md#docker-desktop-license-agreement).

## Sign in to Docker Desktop

After you’ve successfully installed and started Docker Desktop, we recommend that you authenticate using the **Sign in/Create ID** option from the Docker menu.

Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](../docker-hub/download-rate-limit.md).

Authenticating through Docker Desktop also allows administrators to configure Docker Desktop and set up guardrails using the [Image Access Management](../docker-hub/image-access-management.md) feature to allow team members to:

* Only have access to Trusted Content on Docker Hub
* Pull only from the specified categories of images

To authenticate using Docker Desktop, click the **Sign in/Create ID** option on the Docker menu and follow the on-screen instructions.

## Configure Docker Desktop

To learn about the various UI options and their usage, see:

* [Docker Desktop for Mac user manual](mac/index.md)
* [Docker Desktop for Windows user manual](windows/index.md)

## Release notes

For information about new features, improvements, and bug fixes in Docker Desktop releases, see:

* [Docker Desktop for Mac Release notes](mac/release-notes/index.md)
* [Docker Desktop for Windows Release notes](windows/release-notes/index.md)
