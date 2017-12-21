---
title: Basic deployment
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

Once Interlock has been deployed you are now ready to launch and publish applications.
Using [Service Labels](https://docs.docker.com/engine/reference/commandline/service_create/#set-metadata-on-a-service--l-label)
the service is configured to publish itself to the load balancer.

Note: the examples below assume a DNS entry (or local hosts entry if you are testing local) exists
for each of the applications.

To publish we will create a Docker Service using two labels:

- `com.docker.lb.hosts`
- `com.docker.lb.port`

The `com.docker.lb.hosts` label instructs Interlock where the service should be available.
The `com.docker.lb.port` label instructs what port the proxy service should use to access
the upstreams.

In this example we will publish a demo service to the host `demo.local`.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will deploy the application:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    ehazlett/docker-demo
6r0wiglf5f3bdpcy6zesh1pzx
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`

```bash
$> curl -s -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"c2f1afe673d4","version":"0.1",request_id":"7bcec438af14f8875ffc3deab9215bc5"}
```

To increase service capacity use the Docker Service [Scale](https://docs.docker.com/engine/swarm/swarm-tutorial/scale-service/) command:

```bash
$> docker service scale demo=4
demo scaled to 4
```

The four service replicas will be configured as upstreams.  The load balancer will balance traffic
across all service replicas.

{% endif %}
