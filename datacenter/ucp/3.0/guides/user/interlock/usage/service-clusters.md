---
title: Service clusters
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
---

In this example we will configure an eight (8) node Swarm cluster that uses service clusters
to route traffic to different proxies.  There are three (3) managers
and five (5) workers.  Two of the workers are configured with node labels to be dedicated
ingress cluster load balancer nodes.  These will receive all application traffic.

This example will not cover the actual deployment of infrastructure.
It assumes you have a vanilla Swarm cluster (`docker init` and `docker swarm join` from the nodes).
See the [Swarm](https://docs.docker.com/engine/swarm/) documentation if you need help
getting a Swarm cluster deployed.

![Interlock Service Clusters](interlock_service_clusters.png)

We will configure the load balancer worker nodes (`lb-00` and `lb-01`) with node labels in order to pin the Interlock Proxy
service.  Once you are logged into one of the Swarm managers run the following to add node labels
to the dedicated ingress workers:

```bash
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-east lb-00
lb-00
$> docker node update --label-add nodetype=loadbalancer --label-add region=us-west lb-01
lb-01
```

You can inspect each node to ensure the labels were successfully added:

```bash
{% raw %}
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-00
map[nodetype:loadbalancer region:us-east]
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-01
map[nodetype:loadbalancer region:us-west]
{% endraw %}
```

Next, we will create a configuration object for Interlock that contains multiple extensions with varying service clusters:

```bash
$> cat << EOF | docker config create service.interlock.conf -
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.us-east]
    Image = "interlockpreview/interlock-extension-nginx:2.0.0-preview"
    Args = ["-D"]
    ServiceName = "interlock-ext-us-east"
    ProxyImage = "nginx:alpine"
    ProxyArgs = []
    ProxyServiceName = "interlock-proxy-us-east"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ServiceCluster = "us-east"
    PublishMode = "host"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 443
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
    Image = "interlockpreview/interlock-extension-nginx:2.0.0-preview"
    Args = ["-D"]
    ServiceName = "interlock-ext-us-west"
    ProxyImage = "nginx:alpine"
    ProxyArgs = []
    ProxyServiceName = "interlock-proxy-us-west"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ServiceCluster = "us-west"
    PublishMode = "host"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 443
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
Note that we are using "host" mode networking in order to use the same ports (`80` and `443`) in the cluster.  We cannot use ingress
networking as it reserves the port across all nodes.  If you want to use ingress networking you will have to use different ports
for each service cluster.

Next we will create a dedicated network for Interlock and the extensions:

```bash
$> docker network create -d overlay interlock
```

Now we can create the Interlock service:

```bash
$> docker service create \
    --name interlock \
    --mount src=/var/run/docker.sock,dst=/var/run/docker.sock,type=bind \
    --network interlock \
    --constraint node.role==manager \
    --config src=service.interlock.conf,target=/config.toml \
    interlockpreview/interlock:2.0.0-preview -D run -c /config.toml
sjpgq7h621exno6svdnsvpv9z
```

## Configure Proxy Services
Once we have the node labels we can re-configure the Interlock Proxy services to be constrained to the
workers for each region.  Again, from a manager run the following to pin the proxy services to the ingress workers:

```bash
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    --constraint-add node.labels.region==us-east \
    interlock-proxy-us-east
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    --constraint-add node.labels.region==us-west \
    interlock-proxy-us-west
```

We are now ready to deploy applications.  First we will create individual networks for each application:

```bash
$> docker network create -d overlay demo-east
$> docker network create -d overlay demo-west
```

Next we will deploy the application in the `us-east` service cluster:

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

Now we deploy the application in the `us-west` service cluster:

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

Only the service cluster that is designated will be configured for the applications.  For example, the `us-east` service cluster
will not be configured to serve traffic for the `us-west` service cluster and vice versa.  We can see this in action when we
send requests to each service cluster.

When we send a request to the `us-east` service cluster it only knows about the `us-east` application (be sure to ssh to the `lb-00` node):

```bash
{% raw %}
$> curl -H "Host: demo-east.local" http://$(docker node inspect -f '{{ .Status.Addr  }}' lb-00)/ping
{"instance":"1b2d71619592","version":"0.1","metadata":"us-east","request_id":"3d57404cf90112eee861f9d7955d044b"}
$> curl -H "Host: demo-west.local" http://$(docker node inspect -f '{{ .Status.Addr  }}' lb-00)/ping
<html>
<head><title>404 Not Found</title></head>
<body bgcolor="white">
<center><h1>404 Not Found</h1></center>
<hr><center>nginx/1.13.6</center>
</body>
</html>
{% endraw %}
```

Application traffic is isolated to each service cluster.  Interlock also ensures that a proxy will only be updated if it has corresponding updates
to its designated service cluster.  So in this example, updates to the `us-east` cluster will not affect the `us-west` cluster.  If there is a problem
the others will not be affected.

