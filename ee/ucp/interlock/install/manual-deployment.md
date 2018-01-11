---
title: Deploy Interlock manually
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

## Requirements

- [Docker](https://www.docker.com) version 17.06+ is required to use Interlock
- Docker must be running in [Swarm mode](https://docs.docker.com/engine/swarm/)
- Internet access (see [Offline Installation](offline.md) for installing without internet access)

## Deployment
Interlock uses a configuration file for the core service. The following is an example config
to get started.  In order to utilize the deployment and recovery features in Swarm we will
create a Docker Config object:

```bash
$> cat << EOF | docker config create service.interlock.conf -
ListenAddr = ":8080"
DockerURL = "unix:///var/run/docker.sock"
PollInterval = "3s"

[Extensions]
  [Extensions.default]
    Image = "interlockpreview/interlock-extension-nginx:2.0.0-preview"
    Args = ["-D"]
    ProxyImage = "nginx:alpine"
    ProxyArgs = []
    ProxyConfigPath = "/etc/nginx/nginx.conf"
    ServiceCluster = ""
    PublishMode = "ingress"
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

Next we will create a dedicated network for Interlock and the extensions:

```bash
$> docker network create -d overlay interlock
```

Now we can create the Interlock service.  Note the requirement to constrain to a manager.  The
Interlock core service must have access to a Swarm manager, however the extension and proxy services
are recommended to run on workers.  See the [Production](production.md) section for more information
on setting up for an production environment.

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

There should be three (3) services created.  One for the Interlock service,
one for the extension service and one for the proxy service:

```bash
$> docker service ls
ID                  NAME                MODE                REPLICAS            IMAGE                                                       PORTS
lheajcskcbby        modest_raman        replicated          1/1                 nginx:alpine                                                *:80->80/tcp *:443->443/tcp
oxjvqc6gxf91        keen_clarke         replicated          1/1                 interlockpreview/interlock-extension-nginx:2.0.0-preview
sjpgq7h621ex        interlock           replicated          1/1                 interlockpreview/interlock:2.0.0-preview
```

The Interlock traffic layer is now deployed.  Continue with the [Deploying Applications](/usage/index.md) to publish applications.

{% endif %}
