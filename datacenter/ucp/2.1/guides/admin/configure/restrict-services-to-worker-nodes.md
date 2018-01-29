---
title: Restrict services to worker nodes
description: Learn how to configure Universal Control Plane to only allow running services in worker nodes.
keywords: docker, ucp, configuration, worker
---

You can configure UCP to only allow users to deploy and run services in
worker nodes. This ensures all cluster management functionality stays
performant, and makes the cluster more secure.

If a user deploys a malicious service that can affect the node where it
is running, it can't affect other nodes in the cluster, or
any cluster management functionality.

To restrict users from deploying to manager nodes, log in with administrator
credentials to the **UCP web UI**, navigate to the **Admin Settings**
page, and choose **Scheduler**.

![](../../images/restrict-services-to-worker-nodes-1.png){: .with-border}

You can then choose if user services should be allowed to run on manager nodes
or not.

## Where to go next

* [Use domain names to access your services](use-domain-names-to-access-services.md)
