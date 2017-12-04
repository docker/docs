---
title: Tune Interlock
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
---

It is [recommended](/install/production/) to constrain the proxy service to multiple dedicated worker nodes.
Here are a few other tips for tuning:

## Stopping
You can adjust the stop signal and period by using the `stop-signal` and `stop-grace-period` settings.  For example,
to set the stop signal to `SIGTERM` and grace period to ten (10) seconds use the following:

```bash
$> docker service update --stop-signal=SIGTERM --stop-grace-period=10s interlock-proxy
```

## Update Actions
In the event of an update failure the default Swarm action is to "pause".  This will cause no more updates to happen from
Interlock until operator intervention.  You can change this behavior by setting the `update-failure-action` setting.  For example,
to automatically rollback to the previous configuration upon failure use the following:

```bash
$> docker service update --update-failure-action=rollback interlock-proxy
```

## Update Interval
By default Interlock will configure the proxy service by rolling update.  If you would like to have more time between proxy
updates, such as to let the service settle, you can use the `update-delay` setting.  For example, if you want to have
thirty (30) seconds between updates you can use the following:

```bash
$> docker service update --update-delay=30s interlock=proxy
```
