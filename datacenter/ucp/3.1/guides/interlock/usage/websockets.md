---
title: Use websockets
description: Learn how to use websockets in your swarm services.
keywords: routing, proxy, websockets
---

First, create an overlay network to isolate and secure service traffic:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next, create the service with websocket endpoints:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --detach=false \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.websocket_endpoints=/ws \
    ehazlett/websocket-chat
```

> **Note**: for this to work, you must have an entry for `demo.local` in your local hosts (i.e. `/etc/hosts`) file.
> This uses the browser for websocket communication, so you must have an entry or use a routable domain.

Interlock detects when the service is available and publishes it. Once tasks are running
and the proxy service is updated, the application should be available via `http://demo.local`. Open
two instances of your browser and text should be displayed on both instances as you type.
