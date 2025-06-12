---
title: Docker Offload usage and billing
linktitle: Usage & billing
weight: 30
description: Learn about Docker Offload usage and how to monitor your cloud resources.
keywords: cloud, usage, cloud minutes, shared cache, top repositories, cloud builder, Docker Offload
---

{{< summary-bar feature_name="Docker Offload" >}}

## Cloud billing

For Docker Offload, you can view and configure billing on the **Cloud usage**
page in [Docker Home Billing](https://app.docker.com/billing). On this page, you
can:

- View your included budget
- View rates for cloud resources
- Enable or disable on-demand usage
- Add or change payment methods

For more general information about billing, see [Billing](../billing/_index.md).

## Cloud overview

The **Cloud overview** page in Docker Home provides visibility into
how you or your team is using cloud resources to build and run containers.

To view the **Cloud overview** page:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Select the account for which you want to manage Docker Offload.
3. Select **Cloud** > **Cloud overview**.

The following sections describe the available widgets on **Cloud overview**.

### Cloud minutes

This widget shows the total number of cloud minutes used over time. Cloud
minutes represent the time spent running builds and containers in the cloud. You
can use this chart to:

- Track your cloud usage trends over time.
- Spot spikes in usage, which may indicate CI changes or build issues.
- Estimate usage against your subscription limits.

### Shared cache usage

This widget displays data about cache re-use across all builds, helping you
understand how effectively Docker Offload is using the shared build cache. It
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

