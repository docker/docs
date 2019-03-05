---
title: Default host
description: Learn how to publish a service to be a default host.
keywords: routing, proxy
---

# Publishing a default host service

The following example publishes a service to be a default host. The service responds
whenever there is a request to a host that is not configured.

First, create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next, create the initial service:

```bash
$> docker service create \
    --name demo-default \
    --network demo \
    --detach=false \
    --replicas=1 \
    --label com.docker.lb.default_backend=true \
    --label com.docker.lb.port=8080 \
    ehazlett/interlock-default-app
```

Interlock detects when the service is available and publishes it. After tasks are running
and the proxy service is updated, the application is available via any url that is not
configured:


![Default Backend](../../images/interlock_default_backend.png)
