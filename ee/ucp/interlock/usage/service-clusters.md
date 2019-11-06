---
title: Implement service clusters
description: Learn how to route traffic to different proxies using a service cluster.
keywords: ucp, interlock, load balancing, routing
---

Reconfiguring Interlock's proxy can take a long time (1-2 seconds) per overlay network managed by that proxy. In order to scale up to larger number of Docker networks and services routed to by Interlock, you may consider implementing *service clusters*: multiple proxy services managed by Interlock (rather than the default single proxy service), each responsible for routing to a separate set of Docker services and their corresponding networks, thereby minimizing proxy reconfiguration time.

## Prerequisites

In this example, we'll assume you have a UCP cluster set up with at least two worker nodes, `ucp-node-0` and `ucp-node-1`; we'll use these as dedicated proxy servers for two independent Interlock service clusters. 

We'll also assume you've already enabled Interlock, per the [instructions here](https://docs.docker.com/ee/ucp/interlock/deploy/).

## Setting Up Interlock Service Clusters

First, apply some node labels to the UCP workers you've chosen to use as your proxy servers. From a UCP manager:

```bash
docker node update --label-add nodetype=loadbalancer --label-add region=east ucp-node-0
docker node update --label-add nodetype=loadbalancer --label-add region=west ucp-node-1
```

We've labeled `ucp-node-0` to be the proxy for our `east` region, and `ucp-node-1` to be the proxy for our `west` region.

Let's also create a dedicated overlay network for each region's proxy to manage traffic on. We could create many for each, but bear in mind the cumulative performance hit that incurs:

```bash
docker network create --driver overlay eastnet
docker network create --driver overlay westnet
```

Next, modify Interlock's configuration to create two service clusters. Start by writing its current configuration out to a file which you can modify:

```bash
CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0).ConfigName }}' ucp-interlock)
docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > config.toml
```

Open up the file `config.toml` in your text editor of choice, and make the following replacements:

 - Replace `[Extensions.default]` with `[Extensions.east]`
 - Change `ServiceName` to `"ucp-interlock-extension-east"`
 - Change `ProxyServiceName` to `"ucp-interlock-proxy-east"`
 - Add the constraint `"node.labels.region==east"` to the list `ProxyConstraints`
 - Add the key `ServiceCluster="east"` immediately below and inline with `ProxyServiceName`
 - Add the key `Networks=["eastnet"]` immediately below and inline with `ServiceCluster` (*Note this list can contain as many overlay networks as you like; Interlock will _only_ connect to the specified networks, and will connect to them all at startup.*)
 - Change `PublishMode="ingress"` to `PublishMode="host"`
 - Change the section title `[Extensions.default.Labels]` to `[Extensions.east.Labels]`
 - Add the key `"ext_region" = "east"` under the `[Extensions.east.Labels]` section
 - Change the section title `[Extensions.default.ContainerLabels]` to `[Extensions.east.ContainerLabels]`
 - Change the section title `[Extensions.default.ProxyLabels]` to `[Extensions.east.ProxyLabels]`
 - Add the key `"proxy_region" = "east"` under the `[Extensions.east.ProxyLabels]` section
 - Change the section title `[Extensions.default.ProxyContainerLabels]` to `[Extensions.east.ProxyContainerLabels]`
 - Change the section title `[Extensions.default.Config]` to `[Extensions.east.Config]`
 - [Optional] change `ProxyReplicas=2` to `ProxyReplicas=1`, necessary only if there is a single node labeled to be a proxy for each service cluster.

Finally, cut-and-paste the entire `[Extensions.east]` block below, and change every instance of `east` to `west` to create a second service cluster for your `west` region.

Create a new `docker config` object from this configuration file:

```bash
NEW_CONFIG_NAME="com.docker.ucp.interlock.conf-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"
docker config create $NEW_CONFIG_NAME config.toml
```

Update the `ucp-interlock` service to start using this new configuration:

```bash
docker service update \
  --config-rm $CURRENT_CONFIG_NAME \
  --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
  ucp-interlock
```

Finally, do a `docker service ls`. You should see two services providing Interlock proxies, `ucp-interlock-proxy-east` and `-west`. If you only see one Interlock proxy service, delete it with `docker service rm`. After a moment, the two new proxy services should be created, and Interlock will be successfully configured with two service clusters.

## Deploying Services in Separate Service Clusters

Now that you've set up your service clusters, you can deploy services to be routed to by each proxy by using the `service_cluster` label. Create two example services:

```bash
docker service create --name demoeast \
        --network eastnet \
        --label com.docker.lb.hosts=demo.A \
        --label com.docker.lb.port=8000 \
        --label com.docker.lb.service_cluster=east \
        training/whoami:latest

docker service create --name demowest \
        --network westnet \
        --label com.docker.lb.hosts=demo.B \
        --label com.docker.lb.port=8000 \
        --label com.docker.lb.service_cluster=west \
        training/whoami:latest
```

Recall that `ucp-node-0` was your proxy for the `east` service cluster. Attempt to reach your `whoami` service there:

```bash
curl -H "Host: demo.A" http://<ucp-node-0 public IP>
```

You should receive a response indicating the container ID of the `whoami` container declared by the `demoeast` service. Attempt the same `curl` at `ucp-node-1`'s IP, and it will fail: the Interlock proxy running there only routes traffic to services with the `service_cluster=west` label, connected to the `westnet` Docker network you listed in that service cluster's configuration.

Finally, make sure your second service cluster is working analogously to the first:

```bash
curl -H "Host: demo.B" http://<ucp-node-1 public IP>
```

The service routed by `Host: demo.B` is reachable via (and only via) the Interlock proxy mapped to port 80 on `ucp-node-1`. At this point, you have successfully set up and demonstrated that Interlock can manage multiple proxies routing only to services attached to a select subset of Docker networks.

