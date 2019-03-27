---
title: Application redirects
description: Learn how to implement redirects using swarm services and the
  layer 7 routing solution for UCP.
keywords: routing, proxy, redirects
---

# Implementing redirects
The following example deploys a simple
service that can be reached at `app.example.org`. Requests to `old.example.org` are redirected to that service.

Create a docker-compose.yml file as shown in the following example:

```yaml
version: "3.2"

services:
  demo:
    image: ehazlett/docker-demo
    deploy:
      replicas: 1
      labels:
        com.docker.lb.hosts: app.example.org,old.example.org
        com.docker.lb.network: demo-network
        com.docker.lb.port: 8080
        com.docker.lb.redirects: http://old.example.org,http://app.example.org
    networks:
      - demo-network

networks:
  demo-network:
    driver: overlay
```

Note that the demo service has labels to signal that traffic for both
`app.example.org` and `old.example.org` should be routed to this service.
There is also a label indicating that all traffic directed to `old.example.org`
should be redirected to `app.example.org`.

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

You can also use the CLI to test if the redirect is working, by running the following command:

```bash
curl --head --header "Host: old.example.org" http://<ucp-ip>:<http-port>
```

You should see something like the following output:

```none
HTTP/1.1 302 Moved Temporarily
Server: nginx/1.13.8
Date: Thu, 29 Mar 2018 23:16:46 GMT
Content-Type: text/html
Content-Length: 161
Connection: keep-alive
Location: http://app.example.org/
```

You can also test that the redirect works from your browser. For that, make sure you add entries for both `app.example.org` and
`old.example.org` to your `/etc/hosts` file and map them to the IP address
of a UCP node.



------------------------SHOULD THE FOLLOWING BE INCLUDED AS WELL? ------------------------------------------

The following example publishes a service and configures a redirect from `old.local` to `new.local`.

First, create an overlay network so that service traffic is isolated and secure:

```bash
$> docker network create -d overlay demo
1se1glh749q1i4pw0kf26mfx5
```

Next, create the service with the redirect:

```bash
$> docker service create \
    --name demo \
    --network demo \
    --detach=false \
    --label com.docker.lb.hosts=old.local,new.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.redirects=http://old.local,http://new.local \
    --env METADATA="demo-new" \
    ehazlett/docker-demo
```

Interlock detects when the service is available and publishes it.  After tasks are running
and the proxy service is updated, the application is available via `http://new.local`
with a redirect configured that sends `http://old.local` to `http://new.local`:

```bash
$> curl -vs -H "Host: old.local" http://127.0.0.1
* Rebuilt URL to: http://127.0.0.1/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 80 (#0)
> GET / HTTP/1.1
> Host: old.local
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 302 Moved Temporarily
< Server: nginx/1.13.6
< Date: Wed, 08 Nov 2017 19:06:27 GMT
< Content-Type: text/html
< Content-Length: 161
< Connection: keep-alive
< Location: http://new.local/
< x-request-id: c4128318413b589cafb6d9ff8b2aef17
< x-proxy-id: 48854cd435a4
< x-server-info: interlock/2.0.0-development (147ff2b1) linux/amd64
<
<html>
<head><title>302 Found</title></head>
<body bgcolor="white">
<center><h1>302 Found</h1></center>
<hr><center>nginx/1.13.6</center>
</body>
</html>
```
