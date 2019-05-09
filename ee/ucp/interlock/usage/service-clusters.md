---
title: Implement service clusters
description: Learn how to route traffic to different proxies using a service cluster.
keywords: ucp, interlock, load balancing, routing
---

## Configure Proxy Services
With the node labels, you can re-configure the Interlock Proxy services to be constrained to the
workers for each region. FOr example, from a manager, run the following commands to pin the proxy services to the ingress workers:

```bash
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    --constraint-add node.labels.region==us-east \
    ucp-interlock-proxy-us-east
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    --constraint-add node.labels.region==us-west \
    ucp-interlock-proxy-us-west
```

You are now ready to deploy applications. First, create individual networks for each application:

```bash
$> docker network create -d overlay demo-east
$> docker network create -d overlay demo-west
```

Next, deploy the application in the `us-east` service cluster:

```bash
$> docker service create \
    --name demo-east \
    --network demo-east \
    --detach=true \
    --label com.docker.lb.hosts=demo-east.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.service_cluster=us-east \
    --env METADATA="us-east" \
    ehazlett/docker-demo
```

Now deploy the application in the `us-west` service cluster:

```bash
$> docker service create \
    --name demo-west \
    --network demo-west \
    --detach=true \
    --label com.docker.lb.hosts=demo-west.local \
    --label com.docker.lb.port=8080 \
    --label com.docker.lb.service_cluster=us-west \
    --env METADATA="us-west" \
    ehazlett/docker-demo
```

Only the designated service cluster is configured for the applications. For example, the `us-east` service cluster
is not configured to serve traffic for the `us-west` service cluster and vice versa. You can observe this when you
send requests to each service cluster.

When you send a request to the `us-east` service cluster, it only knows about the `us-east` application. This example uses IP address lookup from the swarm API, so you must `ssh` to a manager node or configure your shell with a UCP client bundle before testing:

```bash
{% raw %}
$> curl -H "Host: demo-east.local" http://$(docker node inspect -f '{{ .Status.Addr  }}' lb-00):8080/ping
{"instance":"1b2d71619592","version":"0.1","metadata":"us-east","request_id":"3d57404cf90112eee861f9d7955d044b"}
$> curl -H "Host: demo-west.local" http://$(docker node inspect -f '{{ .Status.Addr  }}' lb-00):8080/ping
<html>
<head><title>404 Not Found</title></head>
<body bgcolor="white">
<center><h1>404 Not Found</h1></center>
<hr><center>nginx/1.13.6</center>
</body>
</html>
{% endraw %}
```

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

```bash
{% raw %}
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-00
map[nodetype:loadbalancer region:us-east]
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-02
map[nodetype:loadbalancer region:us-west]
{% endraw %}
```

Next, create an Interlock configuration object that contains multiple extensions with varying service clusters.

< Important: The configuration object specified in the following code sample applies to UCP versions 3.0.10 and later, and versions 3.1.4 and later.

If you are working with UCP version 3.0.0 - 3.0.9 or 3.1.0 - 3.1.3, specify `com.docker.ucp.interlock.service-clusters.conf`.

```bash
$> cat << EOF | docker config create com.docker.ucp.interlock.conf-1 -
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.us-east]
    Image = "{{ page.ucp_org }}/ucp-interlock-extension:{{ page.ucp_version }}"
    Args = []
    ServiceName = "ucp-interlock-extension-us-east"
    ProxyImage = "{{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }}"
    ProxyArgs = []
    ProxyServiceName = "ucp-interlock-proxy-us-east"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 2
    ServiceCluster = "us-east"
    PublishMode = "host"
    PublishedPort = 8080
    TargetPort = 80
    PublishedSSLPort = 8443
    TargetSSLPort = 443
    [Extensions.us-east.Config]
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      WorkerProcesses = 1
      RlimitNoFile = 65535
      MaxConnections = 2048
    [Extensions.us-east.Labels]
      ext_region = "us-east"
    [Extensions.us-east.ProxyLabels]
      proxy_region = "us-east"

  [Extensions.us-west]
    Image = "{{ page.ucp_org }}/ucp-interlock-extension:{{ page.ucp_version }}"
    Args = []
    ServiceName = "ucp-interlock-extension-us-west"
    ProxyImage = "{{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }}"
    ProxyArgs = []
    ProxyServiceName = "ucp-interlock-proxy-us-west"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 2
    ServiceCluster = "us-west"
    PublishMode = "host"
    PublishedPort = 8080
    TargetPort = 80
    PublishedSSLPort = 8443
    TargetSSLPort = 443
    [Extensions.us-west.Config]
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      WorkerProcesses = 1
      RlimitNoFile = 65535
      MaxConnections = 2048
    [Extensions.us-west.Labels]
      ext_region = "us-west"
    [Extensions.us-west.ProxyLabels]
      proxy_region = "us-west"
EOF
oqkvv1asncf6p2axhx41vylgt
```
Note that "host" mode networking is used in order to use the same ports (`8080` and `8443`) in the cluster. You cannot use ingress
networking as it reserves the port across all nodes. If you want to use ingress networking, you must use different ports
for each service cluster.

Next, create a dedicated network for Interlock and the extensions:

```bash
$> docker network create -d overlay ucp-interlock
```

Now [enable the Interlock service](../deploy/index.md#enable-layer-7-routing).
