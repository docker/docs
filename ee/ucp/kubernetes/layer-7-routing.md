---
title: Layer 7 routing
description: Learn how to route traffic to your Kubernetes workloads in Docker Enterprise Edition.
keywords: UCP, Kubernetes, ingress, routing
redirect_from:
  - /ee/ucp/kubernetes/deploy-ingress-controller/
---

When you deploy a Kubernetes application, you may want to make it accessible
to users using hostnames instead of IP addresses.

Kubernetes provides **ingress controllers** for this. This functionality is
specific to Kubernetes. If you're trying to route traffic to Swarm-based
applications, check [layer 7 routing with Swarm](../interlock/index.md).

Use an ingress controller when you want to:

* Give your Kubernetes app an externally-reachable URL.
* Load-balance traffic to your app.

A popular ingress controller within the Kubernetes Community is the [NGINX controller](https://github.com/kubernetes/ingress-nginx), and can be used in Docker Enterprise Edition, but it is not directly supported by Docker, Inc.

Learn about [ingress in Kubernetes](https://v1-11.docs.kubernetes.io/docs/concepts/services-networking/ingress/). 

For an example of a YAML NGINX kube ingress deployment, refer to <https://success.docker.com/article/how-to-configure-a-default-tls-certificate-for-the-kubernetes-nginx-ingress-controller>.
