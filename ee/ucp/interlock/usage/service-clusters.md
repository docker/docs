---
title: Implement service clusters
description: Learn how to route traffic to different proxies using a service cluster.
keywords: ucp, interlock, load balancing, routing
---

<<<<<<< HEAD
Reconfiguring Interlock's proxy can take 1-2 seconds per overlay
network managed by that proxy. To scale up to larger number of Docker 
networks and services routed to by Interlock, you may consider implementing 
*service clusters*. Service clusters are multiple proxy services managed by 
Interlock (rather than the default single proxy service), each responsible for 
routing to a separate set of Docker services and their corresponding networks, 
thereby minimizing proxy reconfiguration time.
=======
>{% include enterprise_label_shortform.md %}

## Configure Proxy Services
With the node labels, you can re-configure the Interlock Proxy services to be constrained to the
workers for each region. For example, from a manager, run the following commands to pin the proxy services to the ingress workers:
>>>>>>> master

## Prerequisites

In this example, we'll assume you have a UCP cluster set up with at least two 
worker nodes, `ucp-node-0` and `ucp-node-1`; we'll use these as dedicated 
proxy servers for two independent Interlock service clusters. 

<<<<<<< HEAD
We'll also assume you've already enabled Interlock, per the 
[instructions here](../deploy/), 
with an HTTP port of 80 and an HTTPS port of 8443.

## Setting up Interlock service clusters
=======
Add the networks to the Interlock configuration file. Interlock automatically adds networks to the proxy service upon the next proxy update. See *Minimizing the number of overlay networks* in [Interlock architecture](https://docs.docker.com/ee/ucp/interlock/architecture/) for more information.

> Note
>
> Interlock will _only_ connect to the specified networks, and will connect to them all at startup.
>>>>>>> master

First, apply some node labels to the UCP workers you've chosen to use 
as your proxy servers. From a UCP manager:

```bash
docker node update --label-add nodetype=loadbalancer --label-add region=east ucp-node-0
docker node update --label-add nodetype=loadbalancer --label-add region=west ucp-node-1
```

We've labeled `ucp-node-0` to be the proxy for our `east` region, and 
`ucp-node-1` to be the proxy for our `west` region.

Let's also create a dedicated overlay network for each region's proxy to manage 
traffic on. We could create many for each, but bear in mind the cumulative 
performance hit that incurs:

```bash
docker network create --driver overlay eastnet
docker network create --driver overlay westnet
```

<<<<<<< HEAD
Next, modify Interlock's configuration to create two service clusters. Start 
by writing its current configuration out to a file which you can modify:
=======
Application traffic is isolated to each service cluster.  Interlock also ensures that a proxy is updated only if it has corresponding updates to its designated service cluster. In this example, updates to the `us-east` cluster do not affect the `us-west` cluster.  If there is a problem, the others are not affected.

## Usage

The following example configures an eight (8) node Swarm cluster that uses service clusters
to route traffic to different proxies. This example includes:

- Three (3) managers and five (5) workers
- Four workers that are configured with node labels to be dedicated
ingress cluster load balancer nodes. These nodes receive all application traffic.

This example does not cover infrastructure deployment.
It assumes you have a vanilla Swarm cluster (`docker init` and `docker swarm join` from the nodes).
See the [Swarm](https://docs.docker.com/engine/swarm/) documentation if you need help
getting a Swarm cluster deployed.

![Interlock Service Clusters](../../images/interlock_service_clusters.png)

Configure four load balancer worker nodes (`lb-00` through `lb-03`) with node labels in order to pin the Interlock Proxy
service for each Interlock service cluster.  After you log in to one of the Swarm managers, run the following commands to add node labels to the dedicated ingress workers:

```bash
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-east lb-00
lb-00
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-east lb-01
lb-01
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-west lb-02
lb-02
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-west lb-03
lb-03
```

Inspect each node to ensure the labels were successfully added:
>>>>>>> master

```bash
CURRENT_CONFIG_NAME=$(docker service inspect --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0).ConfigName }}' ucp-interlock)
docker config inspect --format '{{ printf "%s" .Spec.Data }}' $CURRENT_CONFIG_NAME > old_config.toml
```

Make a new config file called `config.toml` with the following content, 
which declares two service clusters, `east` and `west`. 

<<<<<<< HEAD
> **Note** you will have to change the UCP version 
> (`3.2.3` in the example below) to match yours, 
> as well as all instances of `*.ucp.InstanceID` 
> (`vl5umu06ryluu66uzjcv5h1bo` below):

```
=======
> Important
>
> The configuration object specified in the following code sample applies to
> UCP versions 3.0.10 and later, and versions 3.1.4 and later. If you are
> working with UCP version 3.0.0 - 3.0.9 or 3.1.0 - 3.1.3, the config object
> should be named `com.docker.ucp.interlock.service-clusters.conf`.

```bash
$> cat << EOF | docker config create com.docker.ucp.interlock.conf-1 -
>>>>>>> master
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
AllowInsecure = false
PollInterval = "3s"

[Extensions]
  [Extensions.east]
    Image = "docker/ucp-interlock-extension:3.2.3"
    ServiceName = "ucp-interlock-extension-east"
    Args = []
    Constraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux"]
    ConfigImage = "docker/ucp-interlock-config:3.2.3"
    ConfigServiceName = "ucp-interlock-config-east"
    ProxyImage = "docker/ucp-interlock-proxy:3.2.3"
    ProxyServiceName = "ucp-interlock-proxy-east"
    ServiceCluster="east"
    Networks=["eastnet"]
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 1
    ProxyStopSignal = "SIGQUIT"
    ProxyStopGracePeriod = "5s"
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux", "node.labels.region==east"]
    PublishMode = "host"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 8443
    TargetSSLPort = 443
    [Extensions.east.Labels]
      "ext_region" = "east"
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.east.ContainerLabels]
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.east.ProxyLabels]
      "proxy_region" = "east"
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.east.ProxyContainerLabels]
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.east.Config]
      Version = ""
      HTTPVersion = "1.1"
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      MaxConnections = 1024
      ConnectTimeout = 5
      SendTimeout = 600
      ReadTimeout = 600
      IPHash = false
      AdminUser = ""
      AdminPass = ""
      SSLOpts = ""
      SSLDefaultDHParam = 1024
      SSLDefaultDHParamPath = ""
      SSLVerify = "required"
      WorkerProcesses = 1
      RLimitNoFile = 65535
      SSLCiphers = "HIGH:!aNULL:!MD5"
      SSLProtocols = "TLSv1.2"
      AccessLogPath = "/dev/stdout"
      ErrorLogPath = "/dev/stdout"
      MainLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" '\n\t\t    '$status $body_bytes_sent \"$http_referer\" '\n\t\t    '\"$http_user_agent\" \"$http_x_forwarded_for\"';"
      TraceLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" $status '\n\t\t    '$body_bytes_sent \"$http_referer\" \"$http_user_agent\" '\n\t\t    '\"$http_x_forwarded_for\" $reqid $msec $request_time '\n\t\t    '$upstream_connect_time $upstream_header_time $upstream_response_time';"
      KeepaliveTimeout = "75s"
      ClientMaxBodySize = "32m"
      ClientBodyBufferSize = "8k"
      ClientHeaderBufferSize = "1k"
      LargeClientHeaderBuffers = "4 8k"
      ClientBodyTimeout = "60s"
      UnderscoresInHeaders = false
      UpstreamZoneSize = 64
      ServerNamesHashBucketSize = 128
      GlobalOptions = []
      HTTPOptions = []
      TCPOptions = []
      HideInfoHeaders = false

  [Extensions.west]
    Image = "docker/ucp-interlock-extension:3.2.3"
    ServiceName = "ucp-interlock-extension-west"
    Args = []
    Constraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux"]
    ConfigImage = "docker/ucp-interlock-config:3.2.3"
    ConfigServiceName = "ucp-interlock-config-west"
    ProxyImage = "docker/ucp-interlock-proxy:3.2.3"
    ProxyServiceName = "ucp-interlock-proxy-west"
    ServiceCluster="west"
    Networks=["westnet"]
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 1
    ProxyStopSignal = "SIGQUIT"
    ProxyStopGracePeriod = "5s"
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux", "node.labels.region==west"]
    PublishMode = "host"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 8443
    TargetSSLPort = 443
    [Extensions.west.Labels]
      "ext_region" = "west"
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.west.ContainerLabels]
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.west.ProxyLabels]
      "proxy_region" = "west"
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.west.ProxyContainerLabels]
      "com.docker.ucp.InstanceID" = "vl5umu06ryluu66uzjcv5h1bo"
    [Extensions.west.Config]
      Version = ""
      HTTPVersion = "1.1"
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      MaxConnections = 1024
      ConnectTimeout = 5
      SendTimeout = 600
      ReadTimeout = 600
      IPHash = false
      AdminUser = ""
      AdminPass = ""
      SSLOpts = ""
      SSLDefaultDHParam = 1024
      SSLDefaultDHParamPath = ""
      SSLVerify = "required"
      WorkerProcesses = 1
      RLimitNoFile = 65535
      SSLCiphers = "HIGH:!aNULL:!MD5"
      SSLProtocols = "TLSv1.2"
      AccessLogPath = "/dev/stdout"
      ErrorLogPath = "/dev/stdout"
      MainLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" '\n\t\t    '$status $body_bytes_sent \"$http_referer\" '\n\t\t    '\"$http_user_agent\" \"$http_x_forwarded_for\"';"
      TraceLogFormat = "'$remote_addr - $remote_user [$time_local] \"$request\" $status '\n\t\t    '$body_bytes_sent \"$http_referer\" \"$http_user_agent\" '\n\t\t    '\"$http_x_forwarded_for\" $reqid $msec $request_time '\n\t\t    '$upstream_connect_time $upstream_header_time $upstream_response_time';"
      KeepaliveTimeout = "75s"
      ClientMaxBodySize = "32m"
      ClientBodyBufferSize = "8k"
      ClientHeaderBufferSize = "1k"
      LargeClientHeaderBuffers = "4 8k"
      ClientBodyTimeout = "60s"
      UnderscoresInHeaders = false
      UpstreamZoneSize = 64
      ServerNamesHashBucketSize = 128
      GlobalOptions = []
      HTTPOptions = []
      TCPOptions = []
      HideInfoHeaders = false
```

If instead you prefer to modify the config file Interlock creates by default, 
the crucial parts to adjust for a service cluster are:

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
 - Copy the entire `[Extensions.east]` block a second time, changing `east` to `west` for your `west` service cluster.

Create a new `docker config` object from this configuration file:

```bash
NEW_CONFIG_NAME="com.docker.ucp.interlock.conf-$(( $(cut -d '-' -f 2 <<< "$CURRENT_CONFIG_NAME") + 1 ))"
docker config create $NEW_CONFIG_NAME config.toml
```
<<<<<<< HEAD
=======
> Note
>
> "Host" mode networking is used in order to use the same ports (`8080` and `8443`) in the cluster. You cannot use ingress
> networking as it reserves the port across all nodes. If you want to use ingress networking, you must use different ports
> for each service cluster.
>>>>>>> master

Update the `ucp-interlock` service to start using this new configuration:

```bash
docker service update \
  --config-rm $CURRENT_CONFIG_NAME \
  --config-add source=$NEW_CONFIG_NAME,target=/config.toml \
  ucp-interlock
```

Finally, do a `docker service ls`. You should see two services providing 
Interlock proxies, `ucp-interlock-proxy-east` and `-west`. If you only see 
one Interlock proxy service, delete it with `docker service rm`. 
After a moment, the two new proxy services should be created, and Interlock 
will be successfully configured with two service clusters.

## Deploying services in separate service clusters

Now that you've set up your service clusters, you can deploy services to be 
routed to by each proxy by using the `service_cluster` label. Create two example 
services:

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

Recall that `ucp-node-0` was your proxy for the `east` service cluster. 
Attempt to reach your `whoami` service there:

```bash
curl -H "Host: demo.A" http://<ucp-node-0 public IP>
```

You should receive a response indicating the container ID of the `whoami` 
container declared by the `demoeast` service. Attempt the same `curl` at 
`ucp-node-1`'s IP, and it will fail: the Interlock proxy running there only 
routes traffic to services with the `service_cluster=west` label, connected 
to the `westnet` Docker network you listed in that service cluster's 
configuration.

Finally, make sure your second service cluster is working analogously to the 
first:

```bash
curl -H "Host: demo.B" http://<ucp-node-1 public IP>
```

The service routed by `Host: demo.B` is reachable via (and only via) the 
Interlock proxy mapped to port 80 on `ucp-node-1`. At this point, you have 
successfully set up and demonstrated that Interlock can manage multiple 
proxies routing only to services attached to a select subset of Docker 
networks.

