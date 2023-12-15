---
description: Learn about Docker Hub features and the administrative tasks available
  in Docker Hub
keywords: Docker, docker, docker hub, hub, overview
title: Overview of Docker Hub
grid:
- title: Create a Docker ID
  description: Register and create a new Docker ID.
  icon: fingerprint
  link: /docker-id
- title: Create a repository
  description: Create a repository to share your images with your team, customers,
    or the Docker community.
  icon: inbox
  link: /docker-hub/repos
- title: Quickstart
  description: Step-by-step instructions on getting started on Docker Hub.
  icon: explore
  link: /docker-hub/quickstart
- title: Release notes
  description: Find out about new features, improvements, and bug fixes.
  icon: note_add
  link: /docker-hub/release-notes
---

Docker Hub is a service provided by Docker for finding and sharing container images.

It's the world’s largest repository of container images with an array of content sources including container community developers, open source projects, and independent software vendors (ISV) building and distributing their code in containers.

Docker Hub is also where you can go to [change your Docker account settings and carry out administrative tasks](/admin/).

{{< tabs >}}
{{< tab name="What key features are included in Docker Hub?" >}}
* [Repositories](../docker-hub/repos/index.md): Push and pull container images.
* [Builds](builds/index.md): Automatically build container images from
GitHub and Bitbucket and push them to Docker Hub.
* [Webhooks](webhooks.md): Trigger actions after a successful push
  to a repository to integrate Docker Hub with other services.
* [Docker Hub CLI](https://github.com/docker/hub-tool#readme) tool (currently experimental) and an API that allows you to interact with Docker Hub.
  * Browse through the [Docker Hub API](/docker-hub/api/latest/) documentation to explore the supported endpoints.
{{< /tab >}}
{{< tab name="What administrative tasks can I perform in Docker Hub?" >}}
* [Create and manage teams and organizations](orgs.md)
* [Create a company](creating-companies.md)
* [Enforce sign in](configure-sign-in.md)
* Set up [SSO](../security/for-admins/single-sign-on/index.md) and [SCIM](../security/for-admins/scim.md)
* Use [Group mapping](group-mapping.md)
* [Carry out domain audits](domain-audit.md)
* [Use Image Access Management](image-access-management.md) to control developers' access to certain types of images
* [Turn on Registry Access Management](../security/for-admins/registry-access-management.md)
{{< /tab >}}
{{< /tabs >}}

{{< grid >}}
