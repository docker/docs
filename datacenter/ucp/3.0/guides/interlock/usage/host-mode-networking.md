---
title: Host mode networking
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

In some scenarios operators cannot use the overlay networks.  Interlock supports
host mode networking in a variety of ways (proxy only, Interlock only, application only, hybrid).

In this example we will configure an eight (8) node Swarm cluster that uses host mode
networking to route traffic without using overlay networks. There are three (3) managers
and five (5) workers.  Two of the workers are configured with node labels to be dedicated
ingress cluster load balancer nodes.  These will receive all application traffic.

This example will not cover the actual deployment of infrastructure.
It assumes you have a vanilla Swarm cluster (`docker init` and `docker swarm join` from the nodes).
See the [Swarm](https://docs.docker.com/engine/swarm/) documentation if you need help
getting a Swarm cluster deployed.

Note: when using host mode networking you will not be able to use the DNS service discovery as that
requires overlay networking.  You can use other tooling such as [Registrator](https://github.com/gliderlabs/registrator)
that will give you that functionality if needed.

We will configure the load balancer worker nodes (`lb-00` and `lb-01`) with node labels in order to pin the Interlock Proxy
service.  Once you are logged into one of the Swarm managers run the following to add node labels
to the dedicated load balancer worker nodes:

```bash
$> docker node update --label-add nodetype=loadbalancer lb-00
lb-00
$> docker node update --label-add nodetype=loadbalancer lb-01
lb-01
```

You can inspect each node to ensure the labels were successfully added:

```bash
{% raw %}
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-00
map[nodetype:loadbalancer]
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-01
map[nodetype:loadbalancer]
{% endraw %}
```

Next, we will create a configuration object for Interlock that specifies host mode networking:

```bash
$> cat << EOF | docker config create service.interlock.conf -
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.default]
    Image = "interlockpreview/interlock-extension-nginx:2.0.0-preview"
    Args = []
    ServiceName = "interlock-ext"
    ProxyImage = "nginx:alpine"
    ProxyArgs = []
    ProxyServiceName = "interlock-proxy"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    PublishMode = "host"
    PublishedPort = 80
    TargetPort = 80
    PublishedSSLPort = 443
    TargetSSLPort = 443
    [Extensions.default.Config]
      User = "nginx"
      PidPath = "/var/run/proxy.pid"
      WorkerProcesses = 1
      RlimitNoFile = 65535
      MaxConnections = 2048
EOF
oqkvv1asncf6p2axhx41vylgt
```

Note the `PublishMode = "host"` setting.  This instructs Interlock to configure the proxy service for host mode networking.

Now we can create the Interlock service also using host mode networking:

```bash
$> docker service create \
    --name interlock \
    --mount src=/var/run/docker.sock,dst=/var/run/docker.sock,type=bind \
    --constraint node.role==manager \
    --publish mode=host,target=8080 \
    --config src=service.interlock.conf,target=/config.toml \
    interlockpreview/interlock:2.0.0-preview -D run -c /config.toml
sjpgq7h621exno6svdnsvpv9z
```

## Configure Proxy Services
Once we have the node labels we can re-configure the Interlock Proxy services to be constrained to the
workers.  Again, from a manager run the following to pin the proxy services to the load balancer worker nodes:

```bash
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    interlock-proxy
```

Now we can deploy the application:

```bash
$> docker service create \
    --name demo \
    --detach=false \
    --label com.docker.lb.hosts=demo.local \
    --label com.docker.lb.port=8080 \
    --publish mode=host,target=8080 \
    --env METADATA="demo" \
    ehazlett/docker-demo
```

This will run the service using host mode networking.  Each task for the service will have a high port (i.e. 32768) and use
the node IP address to connect.  You can see this when inspecting the headers from the request:

```bash
$> curl -vs -H "Host: demo.local" http://127.0.0.1/ping
curl -vs -H "Host: demo.local" http://127.0.0.1/ping
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
< Date: Fri, 10 Nov 2017 15:38:40 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 110
< Connection: keep-alive
< Set-Cookie: session=1510328320174129112; Path=/; Expires=Sat, 11 Nov 2017 15:38:40 GMT; Max-Age=86400
< x-request-id: e4180a8fc6ee15f8d46f11df67c24a7d
< x-proxy-id: d07b29c99f18
< x-server-info: interlock/2.0.0-preview (17476782) linux/amd64
< x-upstream-addr: 172.20.0.4:32768
< x-upstream-response-time: 1510328320.172
<
{"instance":"897d3c7b9e9c","version":"0.1","metadata":"demo","request_id":"e4180a8fc6ee15f8d46f11df67c24a7d"}
```

{% endif %}
