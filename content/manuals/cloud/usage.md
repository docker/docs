---
title: Docker Cloud usage
linktitle: Usage
weight: 30
description: Learn about Docker Cloud usage and how to monitor your cloud resources.
keywords: cloud, usage, cloud minutes, shared cache, top repositories, cloud builder, Docker Cloud
---

Docker Cloud provides visibility into how your team is using cloud resources to
build and run containers. You can monitor your organization's activity from the
**Cloud overview** page in Docker Cloud.

## Credits and consumption

TBD

## Cloud minutes

This widget shows the total number of cloud minutes used over time. Cloud
minutes represent the time spent running builds and containers in the cloud. You
can use this chart to:

- Track your cloud usage trends over time.
- Spot spikes in usage, which may indicate CI changes or build issues.
- Estimate usage against your subscription limits.

## Shared cache usage

This widget displays data about cache re-use across all builds, helping you
understand how effectively Docker Cloud is using the shared build cache. It
provides insight into:

- The percentage of cache hits vs. misses.
- How much estimated build time is saved by reusing cache layers.
- Opportunities to improve cache efficiency by tuning your Dockerfiles or build
  strategy.

## Top repositories built

This widget highlights the repositories with the highest build activity in
Docker Cloud. This widget helps you understand which projects consume the most
cloud resources and how efficiently they’re being built.

It includes both aggregated metrics and per-repository details to give you a
comprehensive view.

Use this widget to:

- Identify build hotspots: See which repositories are consuming the most build
  time and resources.
- Spot trends: Monitor how build activity evolves across your projects.
- Evaluate efficiency: Check which repositories benefit most from cache re-use.
- Target improvements: Flag repositories with low cache hits or high failure
  rates for optimization.
