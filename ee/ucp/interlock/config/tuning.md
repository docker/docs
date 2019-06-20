---
title: Tune the proxy service
description: Learn how to tune the proxy service for environment optimization
keywords: routing, proxy, interlock
---

## Constrain the proxy service to multiple dedicated worker nodes
Refer to [Proxy service constraints](../deploy/production.md) for information on how to constrain the proxy service to multiple dedicated worker nodes.

## Stop
To adjust the stop signal and period, use the `stop-signal` and `stop-grace-period` settings.  For example,
to set the stop signal to `SIGTERM` and grace period to ten (10) seconds, use the following command:

```bash
$> docker service update --stop-signal=SIGTERM --stop-grace-period=10s interlock-proxy
```

## Update actions
In the event of an update failure, the default Swarm action is to "pause".  This prevents Interlock updates from happening
without operator intervention.  You can change this behavior using the `update-failure-action` setting.  For example,
to automatically rollback to the previous configuration upon failure, use the following command:

```bash
$> docker service update --update-failure-action=rollback interlock-proxy
```

## Update interval
By default, Interlock configures the proxy service using rolling update. For more time between proxy
updates, such as to let a service settle, use the `update-delay` setting.  For example, if you want to have
thirty (30) seconds between updates, use the following command:

```bash
$> docker service update --update-delay=30s interlock=proxy
```
