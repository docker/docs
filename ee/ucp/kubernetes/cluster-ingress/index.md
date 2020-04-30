---
title: Kubernetes Cluster Ingress (Experimental)
description: Learn about Ingress host and path routing for Kubernetes applications.
keywords: ucp, cluster, ingress, kubernetes
redirect_from:
  - /ee/ucp/kubernetes/layer-7-routing/
---

>{% include enterprise_label_shortform.md %}

{% include experimental-feature.md %}

## Cluster Ingress capabilities

Cluster Ingress provides L7 services to traffic entering a Docker Enterprise cluster for a variety of different use-cases that help provide application resilience, security, and observability. Ingress provides dynamic control of L7 routing in a highly available architecture that is also high performing.

UCP's Ingress for Kubernetes is based on the [Istio](https://istio.io/) control-plane and is a simplified deployment focused on just providing ingress services with minimal complexity. This includes features such as:

- L7 host and path routing
- Complex path matching and redirection rules
- Weight-based load balancing
- TLS termination
- Persistent L7 sessions
- Hot config reloads
- Redundant and highly available design

For a detailed look at Istio Ingress architecture, refer to the [Istio Ingress docs.](https://istio.io/docs/tasks/traffic-management/ingress/)

To get started with UCP Ingress, the following help topics are provided:

- [Install Cluster Ingress on a UCP Cluster](install.md)
- [Deploy a Sample Application with Ingress Rules](ingress.md)
- [Deploy a Sample Application with a Canary Release](canary.md)
- [Deploy a Sample Application with Sticky Sessions](sticky.md)

## Where to go next

- [Install Cluster Ingress on to a UCP Cluster](install.md)