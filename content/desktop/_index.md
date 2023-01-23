---
description: Docker Desktop overview
keywords: Docker Desktop, Docker, features, GUI, linux, mac, windows
title: Docker Desktop
aliases:
- /desktop/opensource/
- /docker-for-mac/dashboard/
- /docker-for-mac/opensource/
- /docker-for-windows/dashboard/
- /docker-for-windows/opensource/
grid:
  - title: "Install"
    description: |
      Learn how to install Docker Desktop on
      [Mac](/desktop/install/mac-install/),
      [Windows](/desktop/install/windows-install/), or
      [Linux](/desktop/install/linux-install/).
    icon: "download"
  - title: "Explore"
    description: "Navigate Docker Desktop and learn about its key features."
    icon: "explore"
    link: "/desktop/use-desktop"
  - title: "Release notes"
    description: "Find out about new features, improvements, and bug fixes."
    icon: "note_add"
    link: "/desktop/release-notes"
  - title: "Browse common FAQs"
    description: "Explore general FAQs or FAQs for specific platforms."
    icon: "help"
    link: "/desktop/faqs/"
  - title: "Find additional resources"
    description:
      "Find information on networking features, deploying on Kubernetes and more."
    icon: "all_inbox"
    link: "/desktop/additional-resources"
  - title: "Give feedback"
    description:
      "Provide feedback on Docker Desktop or Docker Desktop features."
    icon: "sms"
    link: "/desktop/feedback"
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a paid
> subscription.

Docker Desktop is a one-click-install application for your Mac, Linux, or Windows environment
that enables you to build and share containerized applications and microservices. 

It provides a straightforward GUI (Graphical User Interface) that lets you manage your containers, applications, and images directly from your machine. Docker Desktop can be used either on it's own or as a complementary tool to the CLI. 

Docker Desktop reduces the time spent on complex setups so you can focus on writing code. It takes care of port mappings, file system concerns, and other default settings, and is regularly updated with bug fixes and security updates.

{{< tabs >}}
{{< tab name="What's included in Docker Desktop?" >}}

- [Docker Engine](../engine/_index.md)
- Docker CLI client
- [Docker Buildx](../build/_index.md)
- [Extensions](extensions/_index.md)
- [Docker Compose](../compose/_index.md)
- [Docker Content Trust](../engine/security/trust/_index.md)
- [Kubernetes](https://github.com/kubernetes/kubernetes/)
- [Credential Helper](https://github.com/docker/docker-credential-helpers/)

{{< /tab >}}
{{< tab name="What are the key features of Docker Desktop?" >}}

- Ability to containerize and share any application on any cloud platform, in multiple languages and frameworks.
- Quick installation and setup of a complete Docker development environment.
- Includes the latest version of Kubernetes.
- On Windows, the ability to toggle between Linux and Windows Server environments to build applications.
- Fast and reliable performance with native Windows Hyper-V virtualization.
- Ability to work natively on Linux through WSL 2 on Windows machines.
- Volume mounting for code and data, including file change notifications and easy access to running containers on the localhost network.

{{< /tab >}}
{{< /tabs >}}

Docker Desktop works with your choice of development tools and languages and
gives you access to a vast library of certified images and templates in
[Docker Hub](https://hub.docker.com/). This enables development teams to extend
their environment to rapidly auto-build, continuously integrate, and collaborate
using a secure repository.

{{< grid >}}
