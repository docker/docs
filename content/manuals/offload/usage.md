---
title: Docker Offload usage and billing
linktitle: Usage & billing
weight: 30
description: Learn about Docker Offload usage and how to monitor your cloud resources.
keywords: cloud, usage, cloud minutes, shared cache, top repositories, cloud builder, Docker Offload
---

{{< summary-bar feature_name="Docker Offload" >}}

> [!NOTE]
>
> All free trial usage granted for the Docker Offload Beta expire after 90
> days from the time they are granted. To continue using Docker Offload after
> your usage expires, you can enable on-demand usage at [Docker Home
> Billing](https://app.docker.com/billing).

## Docker Offload billing

For Docker Offload, you can view and configure billing on the **Docker Offload**
page in [Docker Home Billing](https://app.docker.com/billing). On this page, you
can:

- View your included usage
- View rates for cloud resources
- Enable or disable on-demand usage
- Add or change payment methods

For more general information about billing, see [Billing](../billing/_index.md).

## Docker Offload overview

The Docker Offload overview page in Docker Home provides visibility into
how you or your team is using cloud resources to build and run containers.

To view the **Overview** page:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Select the account for which you want to manage Docker Offload.
3. Select **Offload** > **Overview**.

The following sections describe the available widgets on **Overview**.

### Offload minutes

This widget shows the total number of offload minutes used over time. Offload
minutes represent the time spent running builds and containers in the Offload
environment. You can use this chart to:

- Track your Offload usage trends over time.
- Spot spikes in usage, which may indicate CI changes or build issues.
- Estimate usage against your subscription limits.

### Build cache usage

This widget displays data about cache re-use across all builds, helping you
understand how effectively Docker Offload is using the build cache. It
provides insight into:

- The percentage of cache hits vs. misses.
- How much estimated build time is saved by reusing cache layers.
- Opportunities to improve cache efficiency by tuning your Dockerfiles or build
  strategy.

### Top repositories built

This widget highlights the repositories with the highest build activity for
Docker Offload. This widget helps you understand which projects consume the most
cloud resources and how efficiently they're being built.

It includes both aggregated metrics and per-repository details to give you a
comprehensive view.

Use this widget to:

- Identify build hotspots: See which repositories are consuming the most build
  time and resources.
- Spot trends: Monitor how build activity evolves across your projects.
- Evaluate efficiency: Check which repositories benefit most from cache re-use.
- Target improvements: Flag repositories with low cache hits or high failure
  rates for optimization.

### Top 10 images

This widget shows the top 10 images used in Docker Offload in run sessions. It
provides insight into which images are most frequently used, helping you
understand your team's container usage patterns.
