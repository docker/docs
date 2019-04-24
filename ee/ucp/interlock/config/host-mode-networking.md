---
title: Configure host mode networking
description: Learn how to configure the UCP layer 7 routing solution with
  host mode networking.
keywords: routing, proxy, interlock, load balancing
redirect_from:
  - /ee/ucp/interlock/usage/host-mode-networking/
  - /ee/ucp/interlock/deploy/host-mode-networking/
---

By default, layer 7 routing components communicate with one another using
overlay networks, but Interlock supports
host mode networking in a variety of ways, including proxy only, Interlock only, application only, and hybrid. 

When using host mode networking, you cannot use DNS service discovery,
since that functionality requires overlay networking.
For services to communicate, each service needs to know the IP address of
the node where the other service is running.

To use host mode networking instead of overlay networking:

1. Perform the configuration needed for a production-grade deployment
2. Update the ucp-interlock configuration
3. Deploy your Swarm services

## Configuration for a production-grade deployment

If you have not done so, configure the
[layer 7 routing solution for production](../deploy/production.md).

The `ucp-interlock-proxy` service replicas should then be
running on their own dedicated nodes.

## Update the ucp-interlock config

[Update the ucp-interlock service configuration](./index.md) so that it uses
host mode networking.

Update the `PublishMode` key to:

```toml
PublishMode = "host"
```

When updating the `ucp-interlock` service to use the new Docker configuration,
make sure to update it so that it starts publishing its port on the host:

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

Now you can deploy your swarm services. 
Set up your CLI client with a [UCP client bundle](../../user-access/cli.md),
and deploy the service. The following example deploys a demo
service that also uses host mode networking:

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

In this example, Docker allocates a high random port on the host where the service can be reached.

To test that everything is working, run the following command:

```bash
curl --header "Host: app.example.org" \
  http://<proxy-address>:<routing-http-port>/ping
```

Where:

* `<proxy-address>` is the domain name or IP address of a node where the proxy
service is running.
* `<routing-http-port>` is the [port you're using to route HTTP traffic](index.md).

If everything is working correctly, you should get a JSON result like:

{% raw %}
```json
{"instance":"63b855978452", "version":"0.1", "request_id":"d641430be9496937f2669ce6963b67d6"}
```
{% endraw %}

The following example describes how to configure an eight (8) node Swarm cluster that uses host mode
networking to route traffic without using overlay networks. There are three (3) managers
and five (5) workers.  Two of the workers are configured with node labels to be dedicated
ingress cluster load balancer nodes.  These will receive all application traffic.

This example does not cover the actual deployment of infrastructure.
It assumes you have a vanilla Swarm cluster (`docker init` and `docker swarm join` from the nodes).
See the [Swarm](https://docs.docker.com/engine/swarm/) documentation if you need help
getting a Swarm cluster deployed.

Note: When using host mode networking, you cannot use the DNS service discovery because that
requires overlay networking.  You can use other tooling such as [Registrator](https://github.com/gliderlabs/registrator)
that will give you that functionality if needed.

Configure the load balancer worker nodes (`lb-00` and `lb-01`) with node labels in order to pin the Interlock Proxy
service.  Once you are logged into one of the Swarm managers run the following to add node labels
to the dedicated load balancer worker nodes:

```bash
$> docker node update --label-add nodetype=loadbalancer lb-00
lb-00
$> docker node update --label-add nodetype=loadbalancer lb-01
lb-01
```

Inspect each node to ensure the labels were successfully added:

{% raw %}
```bash
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-00
map[nodetype:loadbalancer]
$> docker node inspect -f '{{ .Spec.Labels  }}' lb-01
map[nodetype:loadbalancer]
```
{% endraw %}

Next, create a configuration object for Interlock that specifies host mode networking:

```bash
$> cat << EOF | docker config create service.interlock.conf -
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.default]
    Image = "{{ page.ucp_org }}/ucp-interlock-extension:{{ page.ucp_version }}"
    Args = []
    ServiceName = "interlock-ext"
    ProxyImage = "{{ page.ucp_org }}/ucp-interlock-proxy:{{ page.ucp_version }}"
    ProxyArgs = []
    ProxyServiceName = "interlock-proxy"
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ProxyReplicas = 1
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

Note the `PublishMode = "host"` setting. This instructs Interlock to configure the proxy service for host mode networking.

Now create the Interlock service also using host mode networking:

```bash
$> docker service create \
    --name interlock \
    --mount src=/var/run/docker.sock,dst=/var/run/docker.sock,type=bind \
    --constraint node.role==manager \
    --publish mode=host,target=8080 \
    --config src=service.interlock.conf,target=/config.toml \
    { page.ucp_org }}/ucp-interlock:{{ page.ucp_version }} -D run -c /config.toml
sjpgq7h621exno6svdnsvpv9z
```

## Configure proxy services
With the node labels, you can re-configure the Interlock Proxy services to be constrained to the
workers. From a manager run the following to pin the proxy services to the load balancer worker nodes:

```bash
$> docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    interlock-proxy
```

Now you can deploy the application:

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

This runs the service using host mode networking. Each task for the service has a high port (for example, 32768) and uses
the node IP address to connect. You can see this when inspecting the headers from the request:

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
