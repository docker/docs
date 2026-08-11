---
title: "Docker Build Cloud: Reclaim your time with fast, multi-architecture builds"
linkTitle: Docker Build Cloud
description: |
  Learn how to build and deploy Docker images to the cloud with Docker Build
  Cloud.
summary: |
  Build applications up to 39x faster using cloud-based resources, shared team
  cache, and native multi-architecture support.
keywords: docker build cloud, cloud builds, multi-architecture, shared cache, ci/cd, build performance
aliases:
  - /learning-paths/docker-build-cloud/
  - /guides/docker-build-cloud/ci/
  - /guides/docker-build-cloud/common-questions/
  - /guides/docker-build-cloud/dev/
  - /guides/docker-build-cloud/why/
params:
  tags: [cicd]
  image: images/learning-paths/build-cloud.png
  time: 10 minutes
---


<!-- vale Vale.Spelling = NO -->

98% of developers spend up to an hour every day waiting for builds to finish
([Incredibuild: 2022 Big Dev Build Times](https://www.incredibuild.com/survey-report-2022)).
Heavy, complex builds can become a major roadblock for development teams,
slowing down both local development and CI/CD pipelines.

<!-- vale Vale.Spelling = YES -->

Docker Build Cloud speeds up image build times to improve developer
productivity, reduce frustrations, and help you shorten the release cycle.

## Who’s this for?

- Anyone who wants to tackle common causes of slow image builds: limited local
  resources, slow emulation, and lack of build collaboration across a team.
- Developers working on older machines who want to build faster.
- Development teams working on the same repository who want to cut wait times
  with a shared cache.
- Developers performing multi-architecture builds who don’t want to spend hours
  configuring and rebuilding for emulators.

## What you’ll learn

- Building container images faster locally and in CI
- Accelerating builds for multi-platform images
- Reusing pre-built images to expedite workflows

## Tools integration

Works well with Docker Compose, GitHub Actions, and other CI solutions

<div id="dbc-lp-survey-anchor"></div>

## Why Docker Build Cloud?

Docker Build Cloud is a service that lets you build container images faster,
both locally and in CI. Builds run on cloud infrastructure optimally
dimensioned for your workloads, with no configuration required. The service
uses a remote build cache, ensuring fast builds anywhere and for all team
members.

Docker Build Cloud provides several benefits over local builds:

- Improved build speed
- Shared build cache
- Native multi-platform builds

There’s no need to worry about managing builders or infrastructure — simply
connect to your builders and start building. Each cloud builder provisioned to
an organization is completely isolated to a single Amazon EC2 instance, with a
dedicated EBS volume for build cache and encryption in transit. That means
there are no shared processes or data between cloud builders.

{{< youtube-embed "8AqKhEO2PQA" >}}

<div id="dbc-lp-survey-anchor"></div>

## Demo: set up and use Docker Build Cloud in development

With Docker Build Cloud, you can easily shift the build workload from local machines
to the cloud, helping you achieve faster build times, especially for multi-platform builds.

In this demo, you'll see:

- How to setup the builder locally
- How to use Docker Build Cloud with Docker Compose
- How the image cache speeds up builds for others on your team

{{< youtube-embed "oPGq2AP5OtQ" >}}

<div id="dbc-lp-survey-anchor"></div>

## Demo: Using Docker Build Cloud in CI

Docker Build Cloud can significantly decrease the time it takes for your CI builds
take to run, saving you time and money.

Since the builds run remotely, your CI runner can still use the Docker tooling CLI
without needing elevated permissions, making your builds more secure by default.

In this demo, you will see:

- How to integrate Docker Build Cloud into a variety of CI platforms
- How to use Docker Build Cloud in GitHub Actions to build multi-architecture images
- Speed differences between a workflow using Docker Build Cloud and a workflow running natively
- How to use Docker Build Cloud in a GitLab Pipeline

{{< youtube-embed "wvLdInoVBGg" >}}

<div id="dbc-lp-survey-anchor"></div>

## Common challenges and questions

#### Is Docker Build Cloud a standalone product or a part of Docker Desktop?

Docker Build Cloud is a service that can be used both with Docker Desktop and
standalone. It lets you build your container images faster, both locally and in
CI, with builds running on cloud infrastructure. The service uses a remote
build cache, ensuring fast builds anywhere and for all team members.

When used with Docker Desktop, the [Builds view](/desktop/use-desktop/builds/)
works with Docker Build Cloud out-of-the-box. It shows information about your
builds and those initiated by your team members using the same builder,
enabling collaborative troubleshooting.

To use Docker Build Cloud without Docker Desktop, you must
[download and install](/build-cloud/setup/#use-docker-build-cloud-without-docker-desktop)
a version of Buildx with support for Docker Build Cloud (the `cloud` driver).
If you plan on building with Docker Build Cloud using the `docker compose
build` command, you also need a version of Docker Compose that supports Docker
Build Cloud.

#### How does Docker Build Cloud work with Docker Compose?

Docker Compose works out of the box with Docker Build Cloud. Install the Docker
Build Cloud-compatible client (buildx) and it works with both commands.

#### How many minutes are included in Docker Build Cloud Team plans?

Pricing details for Docker Build Cloud can be found on the [pricing page](https://www.docker.com/pricing?ref=Docs&refAction=DocsGuidesBuildCloudFaq).

#### I’m a Docker personal user. Can I try Docker Build Cloud?

Docker subscribers (Pro, Team, Business) receive a set number of minutes each
month, shared across the account, to use Build Cloud.

If you do not have a Docker subscription, you may sign up for a free Personal
account and start a trial of Docker Build Cloud. Personal accounts are limited to a
single user.

For teams to receive the shared cache benefit, they must either be on a Docker
Team or Docker Business subscription.

#### Does Docker Build Cloud support CI platforms? Does it work with GitHub Actions?

Yes, Docker Build Cloud can be used with various CI platforms including GitHub
Actions, CircleCI, Jenkins, and others. It can speed up your build pipelines,
which means less time spent waiting and context switching.

Docker Build Cloud can be used with GitHub Actions to automate your build,
test, and deployment pipeline. Docker provides a set of official GitHub Actions
that you can use in your workflows.

Using GitHub Actions with Docker Build Cloud is straightforward. With a
one-line change in your GitHub Actions configuration, everything else stays the
same. You don't need to create new pipelines. Learn more in the [CI
documentation](/build-cloud/ci/) for Docker Build Cloud.

<div id="dbc-lp-survey-anchor"></div>
