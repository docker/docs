---
description: Collecting UCP cluster metrics with Prometheus
keywords: prometheus, metrics, ucp
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

## Deploy Prometheus on worker nodes

To deploy Prometheus on worker nodes in a cluster:

1. Begin by sourcing an admin bundle.

2. Verify that ucp-metrics pods are running on all managers.

```
$ kubectl -n kube-system get pods -l k8s-app=ucp-metrics -o wide
NAME                READY     STATUS    RESTARTS   AGE       IP              NODE
ucp-metrics-hvkr7   3/3       Running   0          4h        192.168.80.66   3a724a-0
```

3. Add a Kubernetes node label to one or more workers.  Here we add a label with key "ucp-metrics" and value "".

```
$ kubectl label node 3a724a-1 ucp-metrics=
node "noah-3a724a-1" labeled
```
