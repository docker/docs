---
description: Collecting UCP cluster metrics with Prometheus
keywords: prometheus, metrics
title: Collect UCP cluster metrics with Prometheus
redirect_from:
- /engine/admin/prometheus/
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and
alerting toolkit. You can configure Docker as a Prometheus target. This topic
shows you how to configure Docker, set up Prometheus to run as a Docker
container, and monitor your Docker instance using Prometheus.

In UCP 3.0, Prometheus servers were standard containers. In UCP 3.1, Prometheus runs as a Kubernetes deployment. By default, this will be a daemonset that runs on every manager node. One benefit of this change is you can set the daemonset to not schedule on any nodes, which effectively disables Prometheus if you don’t use the UCP web interface.

The data is stored locally on disk for each Prometheus server, so data is not replicated on new managers or if you schedule Prometheus to run on a new node. Metrics are not kept longer than 24 hours.

> **Warning**: Upgrading UCP from 3.0.x to 3.1.x causes loss of metrics data.
