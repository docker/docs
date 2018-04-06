---
title: Route traffic to a simple swarm service
description: Learn how to do canary deployments for your Docker swarm services
keywords: routing, proxy
---

Once the [layer 7 routing solution is enabled](../deploy/index.md), you can
start using it in your swarm services.

In this example we'll deploy a simple service which:

* Has a JSON endpoint that returns the ID of the task serving the request.
* Has a web UI that shows how many tasks the service is running.
* Can be reached at `http://app.example.org`.

## Deploy the service

Create a `docker-compose.yml` file with:

```yaml
version: "3.2"

services:
  demo:
    image: ehazlett/docker-demo
    deploy:
      replicas: 1
      labels:
        com.docker.lb.hosts: app.example.org
        com.docker.lb.network: demo-network
        com.docker.lb.port: 8080
    networks:
      - demo-network

networks:
  demo-network:
    driver: overlay
```

Note that:

* The `com.docker.lb.hosts` label defines the hostname for the service. When
the layer 7 routing solution gets a request containing `app.example.org` in
the host header, that request is forwarded to the demo service.
* The `com.docker.lb.network` defines which network the `ucp-interlock-proxy`
should attach to in order to be able to communicate with the demo service.
To use layer 7 routing, your services need to be attached to at least one network.
If your service is only attached to a single network, you don't need to add
a label to specify which network to use for routing.
* The `com.docker.lb.port` label specifies which port the `ucp-interlock-proxy`
service should use to communicate with this demo service.
* Your service doesn't need to expose a port in the swarm routing mesh. All
communications are done using the network you've specified.

Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker stack deploy --compose-file docker-compose.yml demo
```

The `ucp-interlock` service detects that your service is using these labels
and automatically reconfigures the `ucp-interlock-proxy` service.

## Test using the CLI

To test that requests are routed to the demo service, run:

```bash
curl --header "Host: app.example.org" \
  http://<ucp-address>:<routing-http-port>/ping
```

Where:

* `<ucp-address>` is the domain name or IP address of a UCP node.
* `<routing-http-port>` is the [port you're using to route HTTP traffic](../deploy/index.md).

If everything is working correctly, you should get a JSON result like:

```json
{"instance":"63b855978452", "version":"0.1", "request_id":"d641430be9496937f2669ce6963b67d6"}
```

## Test using a browser

Since the demo service exposes an HTTP endpoint, you can also use your browser
to validate that everything is working.

Make sure the `/etc/hosts` file in your system has an entry mapping
`app.example.org` to the IP address of a UCP node. Once you do that, you'll be
able to start using the service from your browser.

![browser](../../images/route-simple-app-1.png){: .with-border }

