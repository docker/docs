---
title: Layer 7 routing
description: Learn how to route traffic to your Kubernetes workloads in
  Docker Enterprise Edition.
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

Kubernetes provides an NGINX ingress controller that is based on the Kubernetes community [NGINX controller](https://github.com/kubernetes/ingress-nginx) and that can be used in Docker EE. This example serves as a recommendation for how to deploy an ingress controller on Docker Enterprise, but it is not directly supported by Docker, Inc.
Learn about [ingress in Kubernetes](https://v1-11.docs.kubernetes.io/docs/concepts/services-networking/ingress/). 

## Deploy NGINX ingress controller

The cluster is ready for the ingress controller deployment, which has three
main components:

- a simple HTTP server, named `default-http-backend`,
- an ingress controller, named `nginx-ingress-controller`, and
- a service that exposes the app, named `ingress-nginx`.

Navigate to the **Create Kubernetes Object** page, and in the **Object YAML**
editor, paste your NGINX ingress controller YAML.

Because Kubernetes role based access control (RBAC) is supported, you can download the [Kubenetes YAML file](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/mandatory.yaml) to create a dedicated namespace and default service account.

For an additional example of a YAML NGINX kube ingress deployment, refer to https://success.docker.com/article/how-to-configure-a-default-tls-certificate-for-the-kubernetes-nginx-ingress-controller.

## Check your deployment

The `default-http-backend` provides a simple service that serves a 404 page
at `/` and serves 200 on the `/healthz` endpoint.

1.  Navigate to the **Controllers** page and confirm that the
    **default-http-backend** and **nginx-ingress-controller** objects are
    scheduled.

    > Scheduling latency
    >
    > It may take several seconds for the HTTP backend and the ingress controller's
    > `Deployment` and `ReplicaSet` objects to be scheduled.
    {: .important}

    ![](../images/deploy-ingress-controller-2.png){: .with-border}

2.  When the workload is running, navigate to the **Load Balancers** page
    and click the **ingress-nginx** service.

    ![](../images/deploy-ingress-controller-3.png){: .with-border}

3.  In the details pane, click the first URL in the **Ports** section.

    A new page opens, displaying `default backend - 404`.

## Check your deployment from the CLI

From the command line, confirm that the deployment is running by using
`curl` with the URL that's shown on the details pane of the **ingress-nginx**
service.

```bash
curl -I http://<ucp-ip>:<ingress port>/
```

This command returns the following result.

```
HTTP/1.1 404 Not Found
Server: nginx/1.13.8
```

Test the server's health ping service by appending `/healthz` to the URL.

```bash
curl -I http://<ucp-ip>:<ingress port>/healthz
```

This command returns the following result.

```
HTTP/1.1 200 OK
Server: nginx/1.13.8
```
