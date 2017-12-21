---
title: Websockets
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

In this example we will publish a service and configure support for websockets.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the service with websocket endpoints:

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

Note: for this to work you must have an entry for `demo.local` in your local hosts (i.e. `/etc/hosts`) file.
This uses the browser for websocket communication so you will need to have an entry or use a routable domain.

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`.  Open
two instances of your browser and you should see text on both instances as you type.

{% endif %}
