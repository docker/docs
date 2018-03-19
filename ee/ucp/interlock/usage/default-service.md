---
title: Set a default service
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
- version: ucp-2.2
---

{% if include.version=="ucp-3.0" %}

The default proxy service used by UCP to provide layer 7 routing is NGINX,
so when users try to access a route that hasn't been configured, they will
see the default NGINX 404 page.

![Default NGINX page]()

In this example we will publish a service to be a default host.  This service will respond
whenever there is a request to a host that is not configured.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the initial service:

```bash
$> docker service create \
    --name demo-default \
    --network demo \
    --detach=false \
    --replicas=1 \
    --label com.docker.lb.defaul_backend=true \
    --label com.docker.lb.port=8080 \
    ehazlett/interlock-default-app
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via any url that is not
configured:


![Default Backend](interlock_default_backend.png)

{% endif %}
