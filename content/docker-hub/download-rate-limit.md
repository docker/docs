---
description: Learn about usage and rate limits for Docker Hub.
keywords: Docker Hub, pulls, download, limit, usage
title: Docker Hub usage and rate limits
---

Understanding your Docker Hub usage and rate limits helps you manage your and
your organization's usage effectively.

## View Docker Hub usage

You can download a CSV file of your or your organization's Docker Hub usage. To
download the file:

1. Sign in to [Docker Hub](https://hub.docker.com).

   If you want to download usage for all members of an organization, you must
   sign in to an account that is an owner for that organization. Otherwise,
   you can only view your own usage. 

2. In Docker Hub, select **Usage**.
3. In the drop-down, select whether to download your personal data or
   data for an organization.
4. In **From** and **To**, select a date range for the data.
5. Select **Send report to email** to have Docker email you a link to the data
   file. Note that email processing time may vary.

The file contains the following comma separated values.

| CSV column         | Definition                                                                                                                                                                                                                 | Usage guidance                                                                                                                                                                |
|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `datehour`           | The date and hour of the pull that resulted in the data transfer.                                                                                                                                                          | This helps in identifying peak usage times and patterns.                                                                                                                      |
| `user_name`          | The Docker ID of the user that pulled the image                                                                                                                                                                            | This lets organization owners track data consumption per user and manage resources effectively.                                                                               |
| `repository`         | The name of the repository of the image that was pulled.                                                                                                                                                                   | This lets you identify which repositories are most frequently accessed and consume most of the data transfer.                                                                 |
| `pat_name`           | Name of the personal access token that was used for authentication with Docker CLI.                                                                                                                                        | Personal access tokens are usually used to authenticate automated tools (Docker Desktop, CI/CD tools, etc.). This is useful for identifying which automated system issued the pull. |
| `ips`                | The IP address that was used to pull the image. This field is aggregated, so more than one IP address may appear, representing all the IPs used to pull an image within the same date and hour.                            | This helps you understand the origin of the data transfer, which is useful for diagnosing and identifying patterns in automated or manual pulls.                              |
| `repository_privacy` | The privacy state of the image repository that was pulled. This can either be `public` or `private`.                                                                                                                       | This distinguishes between public and private repositories to identify which data transfer threshold the pull impacts.                                                        |
| `tag`                | The tag for the image. The tag is only available if the pull request included a tag.                                                                                                                                       | This helps in identifying the image. Tags are often used to identify specific versions or variants of an image.                                                               |
| `digest`             | The unique image digest for the image.                                                                                                                                                                                     | This helps in identifying the image.                                                                                                                                          |
| `version_checks`     | The number of version checks accumulated for the date and hour of each image repository. Depending on the client, a pull request can do a version check to verify the existence of an image or tag without downloading it. | This helps identify the frequency of version checks, which you can use to analyze usage trends and potential unexpected behaviors.                                          |
| `pulls`              | The number of pulls (HEAD & GET requests) accumulated for the date and hour of each image repository.                                                                                                                      | This helps identify the frequency of repository pulls, which you can use to analyze usage trends and potential unexpected behaviors.                                          |

## Rate limit

A user's rate limit is equal to the highest entitlement of their personal
account or any organization they belong to. To take advantage of this, you must
sign in to [Docker Hub](https://hub.docker.com/) as an authenticated user. For
more information, see [How do I authenticate pull
requests](#how-do-i-authenticate-pull-requests). Unauthenticated (anonymous)
users will have the limits enforced via IP.

- Pulls are accounted to the user doing the pull, not to the owner of the image.
- A pull request is defined as up to two `GET` requests on registry manifest
URLs (`/v2/*/manifests/*`).
- A normal image pull makes a single manifest request.
- A pull request for a multi-arch image makes two manifest requests.
- `HEAD` requests aren't counted.
- Some images are unlimited through the [Docker Sponsored Open
  Source](https://www.docker.com/blog/expanded-support-for-open-source-software-projects/)
  and [Docker Verified Publisher](https://www.docker.com/partners/programs)
  programs.

### What's the download rate limit on Docker Hub?

Docker Hub limits the number of Docker image downloads, or pulls, based on the
account type of the user pulling the image. Pull rate limits are based on
individual IP address.

| User type                                                               | Rate limit                           |
|-------------------------------------------------------------------------|--------------------------------------|
| Anonymous users                                                         | 100 pulls per 6 hours per IP address |
| [Authenticated users](#how-do-i-authenticate-pull-requests)             | 200 pulls per 6 hour period          |
| Users with a paid [Docker subscription](https://www.docker.com/pricing) | Up to 5000 pulls per day             |

If you require a higher number of pulls, you can also buy an [Enhanced Service Account add-on](service-accounts.md#enhanced-service-account-add-on-pricing).

### How do I know my pull requests are being limited?

When you issue a pull request and you are over the limit, Docker Hub returns a
`429` response code with the following body when the manifest is requested:

```text
You have reached your pull rate limit. You may increase the limit by authenticating and upgrading: https://www.docker.com/increase-rate-limits
```

This error message appears in the Docker CLI or in the Docker Engine logs.

### How can I check my current rate?

Valid API requests to Hub usually include the following rate limit headers in
the response:

```text
ratelimit-limit
ratelimit-remaining
docker-ratelimit-source
```

These headers are returned on both GET and HEAD requests.

>**Note**
>
> Using GET emulates a real pull and counts towards the limit. Using HEAD won't.
> To check your limits, you need `curl`, `grep`, and `jq` installed.

To get a token anonymously, if you are pulling anonymously:

```console
$ TOKEN=$(curl "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull" | jq -r .token)
```

To get a token with a user account, if you are authenticated (insert your
username and password in the following command):

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

In the previous example, the pull limit is 100 pulls per 21600 seconds (6
hours), and there are 76 pulls remaining.

#### I don't see any RateLimit headers

If you don't see any RateLimit header, it could be because the image or your IP
is unlimited in partnership with a publisher, provider, or an open-source
organization. It could also mean that the user you are pulling as is part of a
paid Docker plan. Pulling that image won’t count toward pull limits if you don’t
see these headers. However, users with a paid Docker subscription pulling more
than 5000 times daily require a [Service
Account](../docker-hub/service-accounts.md) subscription.

### I'm being limited to a lower rate even though I have a paid Docker subscription

To take advantage of the higher limits included in a paid Docker subscription,
you must [authenticate pulls](#how-do-i-authenticate-pull-requests) with your
user account.

A Pro, Team, or a Business tier doesn't increase limits on your images for other
users. See Docker's [Open
Source](https://www.docker.com/blog/expanded-support-for-open-source-software-projects/),
[Publisher](https://www.docker.com/partners/programs), or [Large
Organization](https://www.docker.com/pricing) offerings.

### Other limits

Docker Hub also has an overall rate limit to protect the application and
infrastructure. This limit applies to all requests to Hub properties including
web pages, APIs, and image pulls. The limit is applied per-IP, and while the
limit changes over time depending on load and other factors, it's in the order
of thousands of requests per minute. The overall rate limit applies to all users
equally regardless of account level.

You can differentiate between these limits by looking at the error code. The
"overall limit" returns a simple `429 Too Many Requests` response. The pull
limit returns a longer error message that includes a link to this page.

## How do I authenticate pull requests?

The following section contains information on how to sign in to Docker Hub to
authenticate pull requests.

### Docker Desktop

If you are using Docker Desktop, you can sign in to Docker Hub from the Docker
Desktop menu.

Select **Sign in / Create Docker ID** from the Docker Desktop menu and follow
the on-screen instructions to complete the sign-in process.

### Docker Engine

If you're using a standalone version of Docker Engine, run the `docker login`
command from a terminal to authenticate with Docker Hub. For information on how
to use the command, see [docker login](../reference/cli/docker/login.md).

### Docker Swarm

If you're running Docker Swarm, you must use the `-- with-registry-auth` flag to
authenticate with Docker Hub. For more information, see [Create a
service](../reference/cli/docker/service/create.md/#create-a-service). If you
are using a Docker Compose file to deploy an application stack, see [docker
stack deploy](../reference/cli/docker/stack/deploy.md).

### GitHub Actions

If you're using GitHub Actions to build and push Docker images to Docker Hub,
see [login action](https://github.com/docker/login-action#dockerhub). If you are
using another Action, you must add your username and access token in a similar
way for authentication.

### Kubernetes

If you're running Kubernetes, follow the instructions in [Pull an Image from a
Private
Registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)
for information on authentication.

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