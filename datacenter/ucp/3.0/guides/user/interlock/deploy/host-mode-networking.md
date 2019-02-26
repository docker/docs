---
title: Host mode networking
description: Learn how to configure the UCP layer 7 routing solution with
  host mode networking.
keywords: routing, proxy
redirect_from:
  - /ee/ucp/interlock/usage/host-mode-networking/
---

By default the layer 7 routing components communicate with one another using
overlay networks. You can customize the components to use host mode networking
instead.

You can choose to:

* Configure the `ucp-interlock` and `ucp-interlock-extension` services to
communicate using host mode networking.
* Configure the `ucp-interlock-proxy` and your swarm service to communicate
using host mode networking.
* Use host mode networking for all of the components.

In this example we'll start with a production-grade deployment of the layer
7 routing solution and update it so that use host mode networking instead of
overlay networking.

When using host mode networking you won't be able to use DNS service discovery,
since that functionality requires overlay networking.
For two services to communicate, each service needs to know the IP address of
the node where the other service is running.

## Production-grade deployment

If you haven't already, configure the
[layer 7 routing solution for production](production.md).

Once you've done that, the `ucp-interlock-proxy` service replicas should be
running on their own dedicated nodes.

## Update the ucp-interlock config

[Update the ucp-interlock service configuration](configure.md) so that it uses
host mode networking.

Update the `PublishMode` key to:

```toml
PublishMode = "host"
```

When updating the `ucp-interlock` service to use the new Docker configuration,
make sure to update it so that it starts publishes its port on the host:

```bash
docker service update \
  --config-rm $CURRENT_CONFIG_NAME \
  --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
  --publish-add mode=host,target=8080 \
  ucp-interlock
```

The `ucp-interlock` and `ucp-interlock-extension` services are now communicating
using host mode networking.

## Deploy your swarm services

Now you can deploy your swarm services. In this example we'll deploy a demo
service that also uses host mode networking.
Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service:

```bash
docker service create \
  --name demo \
  --detach=false \
  --label com.docker.lb.hosts=app.example.org \
  --label com.docker.lb.port=8080 \
  --publish mode=host,target=8080 \
  --env METADATA="demo" \
  ehazlett/docker-demo
```

Docker allocates a high random port on the host where the service can be reached.
To test that everything is working you can run:

```bash
curl --header "Host: app.example.org" \
  http://<proxy-address>:<routing-http-port>/ping
```

Where:

* `<proxy-address>` is the domain name or IP address of a node where the proxy
service is running.
* `<routing-http-port>` is the [port you're using to route HTTP traffic](index.md).

If everything is working correctly, you should get a JSON result like:

```json
{"instance":"63b855978452", "version":"0.1", "request_id":"d641430be9496937f2669ce6963b67d6"}
```
