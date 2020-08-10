---
description: Download rate limit
keywords: Docker, pull requests, download, limit,
title: Download rate limit
---

Docker has enabled download rate limits for downloads and pull requests on Docker Hub.

A Docker image contains multiple layers. Each layer in a pull request represents a download object. For example, when you download the latest Python image from Docker Hub, you’ll be downloading eight layers and indexes. The download rate limit introduced by Docker caps the number of objects that users can download within a specified timeframe. Any downloads beyond this limit will result in the `429 Too Many Requests` error message.

Docker will gradually impose download rate limits with an eventual limit of  300 downloads per six hours for anonymous users.

Logged in users will not be affected at this time. Therefore, we recommend that you log into [Docker Hub](https://hub.docker.com/){: target="_blank" class="_"} as an authenticated user. For more information, see the following section [How do I authenticate pull requests](#how-do-i-authenticate-pull-requests).

## How do I authenticate pull requests

The following section contains information on how to log into on Docker Hub to authenticate pull requests.

### Docker Desktop

If you are using Docker Desktop, you can log into Docker Hub from the Docker Desktop menu.

Click **Sign in / Create Docker ID** from the Docker Desktop menu and follow the on-screen instructions to complete the sign-in process.

### Docker Engine

If you are using a standalone version of Docker Engine, run the `docker login` command from a terminal to authenticate with Docker Hub. For information on how to use the command, see [docker login](../engine/reference/commandline/login.md).

### Docker Swarm

If you are running Docker Swarm, you must use the `-- with-registry-auth` flag to authenticate with Docker Hub. For more information, see [docker service create](../engine/reference/commandline/service_create.md/#create-a-service). If you are using a Docker Compose file to deploy an application stack, see [docker stack deploy](../engine/reference/commandline/stack_deploy.md).

### GitHub Actions

If you are using GitHub Actions to build and push Docker images to Docker Hub, see [username](https://github.com/docker/build-push-action#username){: target="_blank" class="_"}. If you are using another Action, you must add your username and access token in a similar way for authentication.

### Kubernetes

If you are running Kubernetes, follow the instructions in [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/){: target="_blank" class="_"} for information on authentication.

### Third-party platforms

If you are using any third-party platforms, follow your provider’s instructions on using registry authentication.

- [CircleCI](https://circleci.com/docs/2.0/private-images/){: target="_blank" class="_"}
- [Drone.io](https://docs.drone.io/pipeline/docker/syntax/images/#pulling-private-images){: target="_blank" class="_"}
- [Codefresh](https://codefresh.io/docs/docs/docker-registries/external-docker-registries/docker-hub/){: target="_blank" class="_"}
- [AWS ECS/Fargate](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html){: target="_blank" class="_"}
- [AWS CodeBuild](https://aws.amazon.com/blogs/devops/how-to-use-docker-images-from-a-private-registry-in-aws-codebuild-for-your-build-environment/){: target="_blank" class="_"}
