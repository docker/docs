---
description: Get Started with Docker Hub
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, docs, documentation, accounts, organizations, repositories, groups, teams
redirect_from:
- /docker-hub/overview/
title: Get Started with Docker Hub
---

[Docker Hub](https://hub.docker.com) is a service provided by Docker for finding and sharing container images with your team.

![Docker Hub Landing Page](/docker-hub/images/getting-started.png)

Docker Hub provides the following major features:

* [Repositories](/docker-hub/repos.md): Push and pull container images. Private
  repositories allow you to share container images with your team. Public
  repositories allow you to share them with anyone.
* [Teams & Organizations](/docker-hub/orgs.md): Manage access to private repositories.
* [Official Repositories](/docker-hub/official_repos.md): Pull & use high-quality container images provided by Docker.
* [Publisher Images](/docker-hub/publish/customer_faq.md): Pull & use high-quality container
  images provided by external vendors. Certified images also include support and guarantee
  compatibility with Docker Enterprise.
* [Builds](/docker-hub/builds.md): Automatically build container images and push them to your repositories when you make changes to code in GitHub or BitBucket.
* [Webhooks](/docker-hub/webhooks.md): Trigger actions after a successful push
  to a repository to integrate Docker Hub with other services.


## Sign up for Docker Hub

Start by creating an [account](/docker-hub/accounts.md) at [https://hub.docker.com](https://hub.docker.com).

> **Note**: You can search for and pull Docker images from Hub without logging
> in, however to push images or share them with your team, you must log in.

## Find & Pull an Official Repositories

Docker Hub contains a number of [Official
Repositories](http://hub.docker.com/explore/). These are curated repositories from vendors and contributors to Docker. They contain Docker images from vendors like Canonical, Oracle, and Elastic that you can use as the basis to build your applications and services.

With Official Repositories you know you're using an optimized and
up-to-date image that was built by experts to power your applications.

> **Note**: If you would like to contribute an Official Repository for your
> organization or product, see the documentation on
> [Official Repositories on Docker Hub](/docker-hub/official_repos.md) for more
> information.

## Create your first repository and push an image to it

To create a repo:
1. Log into [Docker Hub](https://hub.docker.com)
2. Click on Create Repositories on the home page:
3. Name it <your namespace>/<redis>

Next, we'll push an images

1. Download [Docker Desktop](https://docker.com/get-started)
2. `docker pull redis` to pull the Official **redis** image from Docker Hub
3. `docker tag redis <namespace>/redis`
4. `docker push <namespace>/redis`

(insert pic here)

## Upgrading your Plan

Your Docker ID includes one private Docker Hub repository for free. If you need
more private repositories, you can upgrade from your free account to a paid
plan.

To upgrade, log in to Docker Hub and click [Upgrade Plan](https://hub.docker.com/account/billing-plans/), in the dropdown menu.

(insert pic here)

### Next Steps

You've successfully create a repo and pushed a Docker image to it. Next:
- Create an [Organization](/docker-hub/orgs.md) to use Docker Hub with your team.
- Automatically build container images from code via [Builds](/docker-hub/builds/index.md).
- [Explore](https://hub.docker.com/explore) Official & Publisher Images

### Docker Commands to Interact with Docker Hub
Docker itself provides access to Docker Hub services via the
[`docker search`](/engine/reference/commandline/search.md),
[`pull`](/engine/reference/commandline/pull.md),
[`login`](/engine/reference/commandline/login.md), and
[`push`](/engine/reference/commandline/push.md) commands.
