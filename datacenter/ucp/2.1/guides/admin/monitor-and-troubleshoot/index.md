---
description: Monitor your Docker Universal Control Plane installation, and learn how
  to troubleshoot it.
keywords: Docker, UCP, troubleshoot
title: Monitor your cluster
---

This article gives you an overview of how to monitor your Docker UCP cluster.

## Check the cluster status from the UI

There are two ways of monitoring the cluster status from the UCP web UI. The
first one by visiting the **Nodes** screen on the UCP web app.

![UCP dashboard](../../images/monitor-ucp-1.png){: .with-border}

In the nodes screen you can see if all the nodes in the cluster are healthy, or
if there is any problem. You can learn more about what a specific node state
means and what you can do about it by reviewing the [full list of node
states](./troubleshoot-node-messages.md).

The second way of monitoring cluster status is by viewing any warning banners
that may appear at the top of the screen. These banners will warn you if the
cluster is experiencing node outages, clock skew or other issues, and will often
include a URL you can follow to learn how to resolve the specific issue reported
by the banner. Admin users may be presented with more banners than normal users.

## Check the cluster status from the CLI

You can also monitor the status of a UCP cluster using the Docker CLI client, by
using a [client certificate
bundle](../../user/access-ucp/cli-based-access.md). 

After that, you can view the status of all UCP nodes with the following command:

```none
$ docker node ls
```

If any node is reported as `Down`, you may cross-reference whether that is an
expected state by looking up the node's status message on the [full list of node
states](./troubleshoot-node-messages.md). As a rule of thumb, if the status
message starts with `[Pending]`, then the current state is transient and the
node is expected to correct itself back into a healthy state eventually.

## Automated status checking

You can use the `https://<manager-url>/_ping` endpoint to check the health of a
single UCP manager node. When you access this endpoint, the UCP manager
validates that all its internal components are working, and returns one of the
following HTTP error codes:

* 200, if all components are healthy
* 500, if one or more components are not healthy

If an admin's client certificate is used as a TLS client certificate for the
`_ping` endpoint, then a detailed error message will be returned if any
components are unhealthy. 

If you're accessing the `_ping` endpoint through a load balancer, you'll have no
way of knowing which UCP manager node is not healthy, so make sure to connect
directly to the URL of a manager node.

## Where to go next

* [Troubleshoot with logs](troubleshoot-with-logs.md)
* [Troubleshoot node states](./troubleshoot-node-messages.md)
