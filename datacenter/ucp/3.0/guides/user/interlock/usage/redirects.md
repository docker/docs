---
title: Application redirects
description: Learn how to implement redirects using swarm services and the
  layer 7 routing solution for UCP.
keywords: routing, proxy, redirects
---

Once the [layer 7 routing solution is enabled](../deploy/index.md), you can
start using it in your swarm services. In this example we'll deploy a simple
service that can be reached at `app.example.org`. We'll also redirect
requests to `old.example.org` to that service.

To do that, create a docker-compose.yml file with:

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
There's also a label indicating that all traffic directed to `old.example.org`
should be redirected to `app.example.org`.

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

You can also use the CLI to test if the redirect is working, by running:

```bash
curl --head --header "Host: old.example.org" http://<ucp-ip>:<http-port>
```

You should see something like:

```none
HTTP/1.1 302 Moved Temporarily
Server: nginx/1.13.8
Date: Thu, 29 Mar 2018 23:16:46 GMT
Content-Type: text/html
Content-Length: 161
Connection: keep-alive
Location: http://app.example.org/
```

You can also test that the redirect works from your browser. For that, you
need to make sure you add entries for both `app.example.org` and
`old.example.org` to your `/etc/hosts` file, mapping them to the IP address
of a UCP node.
