---
description: Docker Hub overview
keywords: Docker, docker, docker hub, hub, overview
title: Overview
grid:
  - title: "Create a Docker ID"
    description:
      "Sign up and create a new Docker ID."
    icon: "fingerprint"
    link: "/docker-id"
  - title: "Create a repository"
    description: "Create a repository to share your images with your team, customers, or the Docker community."
    icon: "explore"
    link: "/docker-hub/repos"
  - title: "Quickstart"
    description: "Step-by-step instructions on getting started on Docker Hub."
    icon: "checklist"
    link: "/dockder-hub"
  - title: "Manage access tokens"
    description: "Create personal access tokens as an alternative to your password."
    icon: "key"
    link: "/docker-hub/access-tokens"
  - title: "Official images"
    description: "A curated set of Docker repositories hosted on Docker Hub."
    icon: "verified"
    link: "/docker-hub/official_images"
  - title: "Release notes"
    description: "Find out about new features, improvements, and bug fixes."
    icon: "note_add"
    link: "/docker-hub/release-notes"
---

[Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} is a service provided by Docker for
finding and sharing container images with your team. It's the world’s largest repository of container images with an array of content sources including container community developers, open source projects and independent software vendors (ISV) building and distributing their code in containers.

Users get access to free public repositories for storing and sharing images or can choose a [subscription plan](https://www.docker.com/pricing){: target="_blank" rel="noopener" class="_"} for private repositories.

Docker Hub provides the following major features:

* [Repositories](../docker-hub/repos/index.md): Push and pull container images.
* [Teams & Organizations](orgs.md): Manage access to private
repositories of container images.
* [Docker Official Images](official_images.md): Pull and use high-quality
container images provided by Docker.
* [Docker Verified Publisher Images](publish/index.md): Pull and use high-
quality container images provided by external vendors.
* [Docker-Sponsored Open Source Images](dsos-program.md): Pull and use high-
quality container images from non-commercial open source projects.
* [Builds](builds/index.md): Automatically build container images from
GitHub and Bitbucket and push them to Docker Hub.
* [Webhooks](webhooks.md): Trigger actions after a successful push
  to a repository to integrate Docker Hub with other services.

Docker provides a [Docker Hub CLI](https://github.com/docker/hub-tool#readme){: target="_blank" rel="noopener" class="_"} tool (currently experimental) and an API that allows you to interact with Docker Hub. Browse through the [Docker Hub API](/docker-hub/api/latest/){: target="_blank" rel="noopener" class="_"} documentation to explore the supported endpoints.

{{< grid >}}