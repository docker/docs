---
title: Context/path based routing
description: Learn how to do route traffic to your Docker swarm services based
  on a url path
keywords: routing, proxy
---

In this example we will publish a service using context or path based routing.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the initial service:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --detach=false \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.context_root=/app \
    --label com.docker.lb.context_root_rewrite=true \
    --env METADATA="demo-context-root" \
    ehazlett/docker-demo
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`:

```bash
$> curl -vs -H "Host: demo.local" http://127.0.0.1/app/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 80 (#0)
> GET /app/ HTTP/1.1
> Host: demo.local
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: nginx/1.13.6
< Date: Fri, 17 Nov 2017 14:25:17 GMT
< Content-Type: text/html; charset=utf-8
< Transfer-Encoding: chunked
< Connection: keep-alive
< x-request-id: 077d18b67831519defca158e6f009f82
< x-proxy-id: 77c0c37d2c46
< x-server-info: interlock/2.0.0-dev (732c77e7) linux/amd64
< x-upstream-addr: 10.0.1.3:8080
< x-upstream-response-time: 1510928717.306
...
```

