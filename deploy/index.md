---
title: Docker Enterprise Edition
description: |
  Learn about Docker Enterprise Edition, the enterprise-grade cluster
  management solution from Docker.
keywords: ucp, overview, orchestration, clustering, Docker EE
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
cli_tabs:
- version: docker-cli-linux
- version: docker-cli-win
- version: kubectl
next_steps:
- path: /engine/install
  title: Install Docker
- path: /get-started/
  title: Get Started with Docker
---

{% if include.ui %}

{% if include.version=="ucp-3.0" %}
Docker Enterprise Edition (EE) is a secure, scalable, and supported container
platform for building and orchestrating applications across multi-tenant Linux
and Windows environments.

Docker Enterprise Edition (EE) is a Containers-as-a-Service platform for IT that manages and secures diverse applications across disparate infrastructure, both on-premises and in the cloud.


Docker Enterprise Edition (EE) enables deploying your workloads for high
availability (HA) onto the orchestrator of your choice. Docker EE automates
many of the tasks that orchestration requires, like provisioning pods, 
containers, and cluster resources. Self-healing components 

Role-based access control applies to all of the orchestrators, and communcation
within the cluster is secured with TLS. [Docker Content Trust](/engine/security/trust/content_trust/) 
is enforced for images on all of the orchestrators. 


Interoperability between swarm and kubernetes workloads (networking/storage)
Fully supports all EE features including LDAP/AD, RBAC, Scanning, Signing Enforcement, Security Policies, etc.
 

{% elsif include.version=="ucp-2.2" %}
Docker Universal Control Plane (UCP) is the enterprise-grade cluster management
solution from Docker. You install it on-premises or in your virtual private
cloud, and it helps you manage your Docker swarm and applications through a 
single interface.

![](../datacenter/images/ucp.png){: .with-border}

## Centralized swarm management

With Docker, you can join up to thousands of physical or virtual machines
together to create a container cluster, or swarm, allowing you to deploy your
applications at scale. Docker Universal Control Plane extends the
functionality provided by Docker to make it easier to manage your swarm
from a centralized place.

You can manage and monitor your container cluster using a graphical UI.

![](../datacenter/images/try-ddc-2.png){: .with-border}

{% endif %}
{% endif %}

{% if include.cli %}
Since UCP exposes the standard Docker API, you can continue using the tools
you already know, including the Docker CLI client, to deploy and manage your
applications.

For example, you can use the `docker info` command to check the
status of a Docker swarm managed by UCP:

{% if include.version=="docker-cli-linux" %}
```bash
$ docker info

Containers: 38
Running: 23
Paused: 0
Stopped: 15
Images: 17
Server Version: 17.06
...
Swarm: active
NodeID: ocpv7el0uz8g9q7dmw8ay4yps
Is Manager: true
ClusterID: tylpv1kxjtgoik2jnrg8pvkg6
Managers: 1
…
```
{% elsif include.version=="docker-cli-win" %}
```powershell
PS> docker info
```
{% elsif include.version=="kubectl" %}
UCP also exposes the standard Kubernetes API, so you can use the usual
Kubernetes commands, like `kubectl cluster-info`:

```bash
$ kubectl cluster-info
```
{% endif %}

## Deploy, manage, and monitor

With Docker UCP, you can manage from a centralized place all of the computing
resources you have available, like nodes, volumes, and networks.

You can also deploy and monitor your applications and services.

## Built-in security and access control

Docker UCP has its own built-in authentication mechanism and integrates with
LDAP services. It also has role-based access control (RBAC), so that you can
control who can access and make changes to your swarm and applications.
[Learn about role-based access control](access-control/index.md).

![](../datacenter/ucp/2.2/guides/images/overview-3.png){: .with-border}

Docker UCP integrates with Docker Trusted Registry so that you can keep the
Docker images you use for your applications behind your firewall, where they
are safe and can't be tampered with.

You can also enforce security policies and only allow running applications
that use Docker images you know and trust.

{% endif %}
