---
description: Learn how download rate limits for image pulls on Docker Hub work
keywords: Docker Hub, pulls, download, limit,
title: Docker Hub rate limit
---

## What's the download rate limit on Docker Hub?

Docker Hub limits the number of Docker image downloads, or pulls, based on the account type of the user pulling the image. Pull rate limits are based on individual IP address.

| User type | Rate limit |
| --------- | ---------- |
| Anonymous users | 100 pulls per 6 hours per IP address |
| [Authenticated users](#how-do-i-authenticate-pull-requests)| 200 pulls per 6 hour period |
| Users with a paid [Docker subscription](https://www.docker.com/pricing) | Up to 5000 pulls per day |

If you require a higher number of pulls, you can also purchase an [Enhanced Service Account add-on](service-accounts.md#enhanced-service-account-add-on-pricing).


## Definition of limits

A user's limit is equal to the highest entitlement of their
personal account or any organization they belong to. To take 
advantage of this, you must log in to 
[Docker Hub](https://hub.docker.com/) 
as an authenticated user. For more information, see
[How do I authenticate pull requests](#how-do-i-authenticate-pull-requests). 
Unauthenticated (anonymous) users will have the limits enforced via IP.

- Pulls are accounted to the user doing the pull, not to the owner of the image.
- A pull request is defined as up to two `GET` requests on registry 
manifest URLs (`/v2/*/manifests/*`).
- A normal image pull makes a 
single manifest request.
- A pull request for a multi-arch image makes two 
manifest requests. 
- `HEAD` requests aren't counted.
- Some images are unlimited through our [Docker Sponsored Open Source](https://www.docker.com/blog/expanded-support-for-open-source-software-projects/) and [Docker Verified Publisher](https://www.docker.com/partners/programs) programs.

## How do I know my pull requests are being limited?

When you issue a pull request and you are over the limit, Docker Hub returns a `429` response code with the following body when the manifest is requested:

```
You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limits
```

This error message appears in the Docker CLI or in the Docker Engine logs.

## How can I check my current rate?

Valid API requests to Hub usually include the following rate limit headers in the response:

```
ratelimit-limit    
ratelimit-remaining
docker-ratelimit-source
```

These headers are returned on both GET and HEAD requests.

>**Note**
>
>Using GET emulates a real pull and counts towards the limit. Using HEAD won't. To check your limits, you need `curl`, `grep`, and `jq` installed.

To get a token anonymously, if you are pulling anonymously:

```console
$ TOKEN=$(curl "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)
```

To get a token with a user account, if you are authenticated (insert your username and password in the following command):

```console
$ TOKEN=$(curl --user 'username:password' "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)
```

Then to get the headers showing your limits, run the following:

```console
$ curl --head -H "Authorization: Bearer $TOKEN" https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest
```

Which should return the following headers:

```http
ratelimit-limit: 100;w=21600
ratelimit-remaining: 76;w=21600
docker-ratelimit-source: 192.0.2.1
```

In the example above, the pull limit is 100 pulls per 21600 seconds (6 hours), and there are 76 pulls remaining.

### I don't see any RateLimit headers

If you don't see any RateLimit header, it could be because the image or your IP is unlimited in partnership with a publisher, provider, or an open-source organization. It could also mean that the user you are pulling as is part of a paid Docker plan. Pulling that image won’t count toward pull limits if you don’t see these headers. However, users with a paid Docker subscription pulling more than 5000 times daily require a [Service Account](../docker-hub/service-accounts.md) subscription.

## I'm being limited to a lower rate even though I have a paid Docker subscription

To take advantage of the higher limits included in a paid Docker subscription, you must [authenticate pulls](#how-do-i-authenticate-pull-requests) with your user account.

A Pro, Team, or a Business tier doesn't increase limits on your images for other users. See Docker's [Open Source](https://www.docker.com/blog/expanded-support-for-open-source-software-projects/), [Publisher](https://www.docker.com/partners/programs), or [Large Organization](https://www.docker.com/pricing) offerings.

## How do I authenticate pull requests?

The following section contains information on how to sign in to Docker Hub to authenticate pull requests.

### Docker Desktop

If you are using Docker Desktop, you can sign in to Docker Hub from the Docker Desktop menu.

Select **Sign in / Create Docker ID** from the Docker Desktop menu and follow the on-screen instructions to complete the sign-in process.

### Docker Engine

If you're using a standalone version of Docker Engine, run the `docker login` command from a terminal to authenticate with Docker Hub. For information on how to use the command, see [docker login](../engine/reference/commandline/login.md).

### Docker Swarm

If you're running Docker Swarm, you must use the `-- with-registry-auth` flag to authenticate with Docker Hub. For more information, see [Create a service](../engine/reference/commandline/service_create.md/#create-a-service). If you are using a Docker Compose file to deploy an application stack, see [docker stack deploy](../engine/reference/commandline/stack_deploy.md).

### GitHub Actions

If you're using GitHub Actions to build and push Docker images to Docker Hub, see [login action](https://github.com/docker/login-action#dockerhub). If you are using another Action, you must add your username and access token in a similar way for authentication.

### Kubernetes

If you're running Kubernetes, follow the instructions in [Pull an Image from a Private Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) for information on authentication.

### Third-party platforms

If you're using any third-party platforms, follow your provider’s instructions on using registry authentication.

- [Artifactory](https://www.jfrog.com/confluence/display/JFROG/Advanced+Settings#AdvancedSettings-RemoteCredentials)
- [AWS CodeBuild](https://aws.amazon.com/blogs/devops/how-to-use-docker-images-from-a-private-registry-in-aws-codebuild-for-your-build-environment/)
- [AWS ECS/Fargate](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html)
- [Azure Pipelines](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml#sep-docreg)
- [Chipper CI](https://docs.chipperci.com/builds/docker/#rate-limit-auth)
- [CircleCI](https://circleci.com/docs/2.0/private-images/)
- [Codefresh](https://codefresh.io/docs/docs/docker-registries/external-docker-registries/docker-hub/)
- [Drone.io](https://docs.drone.io/pipeline/docker/syntax/images/#pulling-private-images)
- [GitLab](https://docs.gitlab.com/ee/user/packages/container_registry/#authenticate-with-the-container-registry)
- [LayerCI](https://layerci.com/docs/advanced-workflows#logging-in-to-docker)
- [TeamCity](https://www.jetbrains.com/help/teamcity/integrating-teamcity-with-docker.html#Conforming+with+Docker+download+rate+limits)

## Other limits

Docker Hub also has an overall rate limit to protect the application 
and infrastructure. This limit applies to all requests to Hub 
properties including web pages, APIs, and image pulls. The limit is 
applied per-IP, and while the limit changes over time depending on load
and other factors, it's in the order of thousands of requests per 
minute. The overall rate limit applies to all users equally
regardless of account level.

You can differentiate between these limits by looking at the error 
code. The "overall limit" returns a simple `429 Too Many Requests` 
response. The pull limit returns a longer error message that
includes a link to this page.
