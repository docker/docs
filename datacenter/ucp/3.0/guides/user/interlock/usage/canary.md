---
title: Canary application instances
description: Learn how to do canary deployments for your Docker swarm services
keywords: routing, proxy
---

In this example we will publish a service and deploy an updated service as canary instances.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the initial service:

```bash
$> docker service create \
    --name demo-v1 \
    --network demo \
    --detach=false \
    --replicas=4 \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --env METADATA="demo-version-1" \
    ehazlett/docker-demo
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`:

```bash
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to demo.local (127.0.0.1) port 80 (#0)
> GET /ping HTTP/1.1
> Host: demo.local
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: nginx/1.13.6
< Date: Wed, 08 Nov 2017 20:28:26 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 120
< Connection: keep-alive
< Set-Cookie: session=1510172906715624280; Path=/; Expires=Thu, 09 Nov 2017 20:28:26 GMT; Max-Age=86400
< x-request-id: f884cf37e8331612b8e7630ad0ee4e0d
< x-proxy-id: 5ad7c31f9f00
< x-server-info: interlock/2.0.0-development (147ff2b1) linux/amd64
< x-upstream-addr: 10.0.2.4:8080
< x-upstream-response-time: 1510172906.714
<
{"instance":"df20f55fc943","version":"0.1","metadata":"demo-version-1","request_id":"f884cf37e8331612b8e7630ad0ee4e0d"}
```

Notice the `metadata` with `demo-version-1`.

Now we will deploy a "new" version:

```bash
$> docker service create \
    --name demo-v2 \
    --network demo \
    --detach=false \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --env METADATA="demo-version-2" \
    --env VERSION="0.2" \
    ehazlett/docker-demo
```

Since this has a replica of one (1) and the initial version has four (4) replicas 20% of application traffic
will be sent to `demo-version-2`:

```bash
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"23d9a5ec47ef","version":"0.1","metadata":"demo-version-1","request_id":"060c609a3ab4b7d9462233488826791c"}
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"f42f7f0a30f9","version":"0.1","metadata":"demo-version-1","request_id":"c848e978e10d4785ac8584347952b963"}
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"c2a686ae5694","version":"0.1","metadata":"demo-version-1","request_id":"724c21d0fb9d7e265821b3c95ed08b61"}
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"1b0d55ed3d2f","version":"0.2","metadata":"demo-version-2","request_id":"b86ff1476842e801bf20a1b5f96cf94e"}
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
{"instance":"c2a686ae5694","version":"0.1","metadata":"demo-version-1","request_id":"724c21d0fb9d7e265821b3c95ed08b61"}
```

To increase traffic to the new version add more replicas with `docker scale`:

```bash
$> docker service scale demo-v2=4
demo-v2
```

To complete the upgrade, scale the `demo-v1` service to zero (0):

```bash
$> docker service scale demo-v1=0
demo-v1
```

This will route all application traffic to the new version.  If you need to rollback, simply scale the v1 service
back up and v2 down.

