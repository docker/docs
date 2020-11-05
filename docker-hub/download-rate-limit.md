---
description: Download rate limit
keywords: Docker, pull requests, download, limit,
title: Download rate limit
---

Docker has enabled download rate limits for pull requests on 
Docker Hub. Limits are determined based on the account type. 
For more information, see [Resource Consumption FAQs](https://www.docker.com/pricing/resource-consumption-updates){: target="_blank" rel="noopener" class="_"} and [Docker Hub Pricing](https://hub.docker.com/pricing){: target="_blank" rel="noopener" class="_"}.

A user's limit will be equal to the highest entitlement of their
personal account or any organization they belong to. To take 
advantage of this, you must log into 
[Docker Hub](https://hub.docker.com/){: target="_blank" rel="noopener" class="_"} 
as an authenticated user. For more information, see
[How do I authenticate pull requests](#how-do-i-authenticate-pull-requests). 
Unauthenticated (anonymous) users will have the limits enforced via IP.

- A pull request is defined as up to two `GET` requests on registry 
manifest URLs (`/v2/*/manifests/*`).
- A normal image pull makes a 
single manifest request.
- A pull request for a multi-arch image makes two 
manifest requests. 
- `HEAD` requests are not counted. 
- Limits are applied based on the user doing the pull, and 
not based on the image being pulled or its owner.

Docker will gradually introduce these rate limits starting November 2nd, 2020.

## How do I know my pull requests are being limited

When you issue a pull request and you are over the limit for your account type, Docker Hub will return a `429` response code with the following body when the manifest is requested:

```
You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limits
```

You will see this error message in the Docker CLI or in the Docker Engine logs.

## How can I check my current rate

Valid manifest API requests to Hub will usually include the following rate limit headers in the response:

```
RateLimit-Limit    
RateLimit-Remaining
```

These headers will be returned on both GET and HEAD requests. Note that using GET emulates a real pull and will count towards the limit; using HEAD will not, so we will use it in this example. To check your limits, you will need `curl`, `grep`, and `jq` installed.

To get a token anonymously (if you are pulling anonymously):

```
$ TOKEN=$(curl "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)
```

To get a token with a user account (if you are authenticating your pulls) - don't forget to insert your username and password in the following command:

```
$ TOKEN=$(curl --user 'username:password' "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)
```

Then to get the headers showing your limits, run the following:

```
$ curl --head -H "Authorization: Bearer $TOKEN" https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest
```

Which should return headers including these:

```
RateLimit-Limit: 100;w=21600
RateLimit-Remaining: 76;w=21600
```

This means my limit is 100 per 21600 seconds (6 hours), and I have 76 pulls remaining.

> Remember that these headers are best-effort and there will be small variations.

### I don't see any RateLimit headers

If you do not see these headers, that means pulling that image would not count towards pull limits. This could be because you are authenticated with a user associated with a Legacy/Pro/Team Docker Hub account, or because the image or your IP is unlimited in partnership with a publisher, provider, or open source organization.

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

If you are using GitHub Actions to build and push Docker images to Docker Hub, see [login action](https://github.com/docker/login-action#dockerhub){: target="_blank" rel="noopener" class="_"}. If you are using another Action, you must add your username and access token in a similar way for authentication.

### Kubernetes

If you are running Kubernetes, follow the instructions in [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/){: target="_blank" rel="noopener" class="_"} for information on authentication.

### Third-party platforms

If you are using any third-party platforms, follow your provider’s instructions on using registry authentication.

- [Artifactory](https://www.jfrog.com/confluence/display/JFROG/Advanced+Settings#AdvancedSettings-RemoteCredentials){: target="_blank" rel="noopener" class="_"}
- [AWS CodeBuild](https://aws.amazon.com/blogs/devops/how-to-use-docker-images-from-a-private-registry-in-aws-codebuild-for-your-build-environment/){: target="_blank" rel="noopener" class="_"}
- [AWS ECS/Fargate](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html){: target="_blank" rel="noopener" class="_"}
- [Azure Pipelines](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml#sep-docreg){: target="_blank" rel="noopener" class="_"}
- [CircleCI](https://circleci.com/docs/2.0/private-images/){: target="_blank" rel="noopener" class="_"}
- [Codefresh](https://codefresh.io/docs/docs/docker-registries/external-docker-registries/docker-hub/){: target="_blank" rel="noopener" class="_"}
- [Drone.io](https://docs.drone.io/pipeline/docker/syntax/images/#pulling-private-images){: target="_blank" rel="noopener" class="_"}
- [GitLab](https://docs.gitlab.com/ee/user/packages/container_registry/#authenticate-with-the-container-registry){: target="_blank" rel="noopener" class="_"}
- [LayerCI](https://layerci.com/docs/advanced-workflows#logging-in-to-docker){: target="_blank" rel="noopener" class="_"}
- [TeamCity](https://www.jetbrains.com/help/teamcity/integrating-teamcity-with-docker.html#Conforming+with+Docker+download+rate+limits){: target="_blank" rel="noopener" class="_"}

## Other limits

Docker Hub also has an overall rate limit to protect the application 
and infrastructure. This limit applies to all requests to Hub 
properties including web pages, APIs, image pulls, etc. The limit is 
applied per-IP, and while the limit changes over time depending on load
and other factors, it is in the order of thousands of requests per 
minute. The overall rate limit applies to all users equally
regardless of account level.

You can differentiate between these limits by looking at the error 
code. The "overall limit" will return a simple `429 Too Many Requests` 
response. The pull limit returns a longer error message that
includes a link to this page.
