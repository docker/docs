---
description: Download rate limit
keywords: Docker, pull requests, download, limit,
title: Download rate limit
---

Docker has enabled download rate limits for pull requests on 
Docker Hub. Limits are determined based on the account type. 
For more information, see [Docker Hub Pricing](https://hub.docker.com/pricing){: target="_blank" rel="noopener" class="_"}.

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

When you issue a pull request and you are over the limit for your account type, Hub will return a `429` response code with the following body when the manifest is requested:
> You have reached your download rate limit: https://docs.docker.com/docker-hub/download-rate-limit/. Upgrade your plan to increase your limit at: https://www.docker.com/pricing

You will see this error message in the Docker CLI or Docker Engine logs.

## How can I check my current rate

Valid, non-rate-limited manifest API reqests to Hub will include the following rate limit headers in the response:

> RateLimit-Limit: \<value\>  
RateLimit-Remaining \<value\>

If you have a proxy or other layer in place that logs your requests, you can inspect the headers of these responses directly.

Otherwise, you can use curl to view these. You will need `curl` and `jq` installed.

> $ TOKEN=$(curl "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)  
$ curl -I -H "Authorization: Bearer $TOKEN" https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest  

Returns for example:

> HTTP/1.1 200 OK  
Content-Length: 2782  
Content-Type: application/vnd.docker.distribution.manifest.v1+prettyjws  
Docker-Content-Digest: sha256:767a3815c34823b355bed31760d5fa3daca0aec2ce15b217c9cd83229e0e2020  
Docker-Distribution-Api-Version: registry/2.0  
Etag: "sha256:767a3815c34823b355bed31760d5fa3daca0aec2ce15b217c9cd83229e0e2020"  
Date: Fri, 23 Oct 2020 20:27:41 GMT  
Strict-Transport-Security: max-age=31536000  

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
- [LayerCI](https://layerci.com/docs/advanced-workflows#logging-in-to-docker){: target="_blank" rel="noopener" class="_"}

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
response. The image download limit returns a longer error message that 
includes a link to this page.
