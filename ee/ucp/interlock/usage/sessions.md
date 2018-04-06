---
title: Persistent (sticky) sessions
description: Learn how to configure your swarm services with persistent sessions
  using UCP.
keywords: routing, proxy
---

In this example we will publish a service and configure the proxy for persistent (sticky) sessions.

# Cookies
In the following example we will show how to configure sticky sessions using cookies.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the service with the cookie to use for sticky sessions:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --detach=false \
    --replicas=5 \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.sticky_session_cookie=session \
    --label com.docker.lb.port=8080 \
    --env METADATA="demo-sticky" \
    ehazlett/docker-demo
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`
and configured to use sticky sessions:

```bash
$> curl -vs -c cookie.txt -b cookie.txt -H "Host: demo.local" http://127.0.0.1/ping
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 80 (#0)
> GET /ping HTTP/1.1
> Host: demo.local
> User-Agent: curl/7.54.0
> Accept: */*
> Cookie: session=1510171444496686286
>
< HTTP/1.1 200 OK
< Server: nginx/1.13.6
< Date: Wed, 08 Nov 2017 20:04:36 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 117
< Connection: keep-alive
* Replaced cookie session="1510171444496686286" for domain demo.local, path /, expire 0
< Set-Cookie: session=1510171444496686286
< x-request-id: 3014728b429320f786728401a83246b8
< x-proxy-id: eae36bf0a3dc
< x-server-info: interlock/2.0.0-development (147ff2b1) linux/amd64
< x-upstream-addr: 10.0.2.5:8080
< x-upstream-response-time: 1510171476.948
<
{"instance":"9c67a943ffce","version":"0.1","metadata":"demo-sticky","request_id":"3014728b429320f786728401a83246b8"}
```

Notice the `Set-Cookie` from the application.  This is stored by the `curl` command and sent with subsequent requests
which are pinned to the same instance.  If you make a few requests you will notice the same `x-upstream-addr`.

# IP Hashing
In this example we show how to configure sticky sessions using client IP hashing.  This is not as flexible or consistent
as cookies but enables workarounds for some applications that cannot use the other method.

First we will create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next we will create the service with the cookie to use for sticky sessions using IP hashing:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --detach=false \
    --replicas=5 \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.ip_hash=true \
    --env METADATA="demo-sticky" \
    ehazlett/docker-demo
```

Interlock will detect once the service is available and publish it.  Once the tasks are running
and the proxy service has been updated the application should be available via `http://demo.local`
and configured to use sticky sessions:

```bash
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 80 (#0)
> GET /ping HTTP/1.1
> Host: demo.local
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: nginx/1.13.6
< Date: Wed, 08 Nov 2017 20:04:36 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 117
< Connection: keep-alive
< x-request-id: 3014728b429320f786728401a83246b8
< x-proxy-id: eae36bf0a3dc
< x-server-info: interlock/2.0.0-development (147ff2b1) linux/amd64
< x-upstream-addr: 10.0.2.5:8080
< x-upstream-response-time: 1510171476.948
<
{"instance":"9c67a943ffce","version":"0.1","metadata":"demo-sticky","request_id":"3014728b429320f786728401a83246b8"}
```

You can use `docker service scale demo=10` to add some more replicas.  Once scaled, you will notice that requests are pinned
to a specific backend.

Note: due to the way the IP hashing works for extensions, you will notice a new upstream address when scaling replicas.  This is
expected as internally the proxy uses the new set of replicas to decide on a backend on which to pin.  Once the upstreams are
determined a new "sticky" backend will be chosen and that will be the dedicated upstream.

