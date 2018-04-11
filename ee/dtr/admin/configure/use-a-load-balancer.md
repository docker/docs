---
title: Use a load balancer
description: Learn how to configure a load balancer to balance user requests across multiple Docker Trusted Registry replicas.
keywords: dtr, load balancer
---

Once you’ve joined multiple DTR replicas nodes for
[high-availability](set-up-high-availability.md), you can configure your own
load balancer to balance user requests across all replicas.

![](../../images/use-a-load-balancer-1.svg)


This allows users to access DTR using a centralized domain name. If a replica
goes down, the load balancer can detect that and stop forwarding requests to
it, so that the failure goes unnoticed by users.

DTR exposes several endpoints you can use to assess if a DTR replica is healthy
or not:

* `/_ping`: Is an unauthenticated endpoint that checks if the DTR replica is
healthy. This is useful for load balancing or other automated health check tasks.
* `/nginx_status`: Returns the number of connections being handled by the
NGINX front-end used by DTR.
* `/api/v0/meta/cluster_status`: Returns extensive information about all DTR
replicas.

## Load balance DTR

DTR does not provide a load balancing service. You can use an on-premises
or cloud-based load balancer to balance requests across multiple DTR replicas.

You can use the unauthenticated `/_ping` endpoint on each DTR replica,
to check if the replica is healthy and if it should remain in the load balancing
pool or not.

Also, make sure you configure your load balancer to:

* Load balance TCP traffic on ports 80 and 443.
* Make sure the load balancer is not buffering requests.
* Make sure the load balancer is forwarding the `Host` HTTP header correctly.
* Make sure there's no timeout for idle connections, or set it to more than 10 minutes.

The `/_ping` endpoint returns a JSON object for the replica being queried of
the form:

```json
{
  "Error": "error message",
  "Healthy": true
}
```

A response of `"Healthy": true` means the replica is suitable for taking
requests. It is also sufficient to check whether the HTTP status code is 200.

An unhealthy replica will return 503 as the status code and populate `"Error"`
with more details on any one of these services:

* Storage container (registry)
* Authorization (garant)
* Metadata persistence (rethinkdb)
* Content trust (notary)

Note that this endpoint is for checking the health of a single replica. To get
the health of every replica in a cluster, querying each replica individually is
the preferred way to do it in real time.


## Configuration examples

Use the following examples to configure your load balancer for DTR.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#nginx" data-group="nginx">NGINX</a></li>
  <li><a data-toggle="tab" data-target="#haproxy" data-group="haproxy">HAProxy</a></li>
  <li><a data-toggle="tab" data-target="#aws">AWS LB</a></li>
</ul>
<div class="tab-content">
  <div id="nginx" class="tab-pane fade in active" markdown="1">
```conf
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

stream {
    upstream dtr_80 {
        server <DTR_REPLICA_1_IP>:80  max_fails=2 fail_timeout=30s;
        server <DTR_REPLICA_2_IP>:80  max_fails=2 fail_timeout=30s;
        server <DTR_REPLICA_N_IP>:80   max_fails=2 fail_timeout=30s;
    }
    upstream dtr_443 {
        server <DTR_REPLICA_1_IP>:443 max_fails=2 fail_timeout=30s;
        server <DTR_REPLICA_2_IP>:443 max_fails=2 fail_timeout=30s;
        server <DTR_REPLICA_N_IP>:443  max_fails=2 fail_timeout=30s;
    }
    server {
        listen 443;
        proxy_pass dtr_443;
    }

    server {
        listen 80;
        proxy_pass dtr_80;
    }
}
```
  </div>
  <div id="haproxy" class="tab-pane fade" markdown="1">
```conf
global
    log /dev/log    local0
    log /dev/log    local1 notice

defaults
        mode    tcp
        option  dontlognull
        timeout connect 5000
        timeout client 50000
        timeout server 50000
### frontends
# Optional HAProxy Stats Page accessible at http://<host-ip>:8181/haproxy?stats
frontend dtr_stats
        mode http
        bind 0.0.0.0:8181
        default_backend dtr_stats
frontend dtr_80
        mode tcp
        bind 0.0.0.0:80
        default_backend dtr_upstream_servers_80
frontend dtr_443
        mode tcp
        bind 0.0.0.0:443
        default_backend dtr_upstream_servers_443
### backends
backend dtr_stats
        mode http
        option httplog
        stats enable
        stats admin if TRUE
        stats refresh 5m
backend dtr_upstream_servers_80
        mode tcp
        option httpchk GET /_ping HTTP/1.1\r\nHost:\ <DTR_FQDN>
        server node01 <DTR_REPLICA_1_IP>:80 check weight 100
        server node02 <DTR_REPLICA_2_IP>:80 check weight 100
        server node03 <DTR_REPLICA_N_IP>:80 check weight 100
backend dtr_upstream_servers_443
        mode tcp
        option httpchk GET /_ping HTTP/1.1\r\nHost:\ <DTR_FQDN>
        server node01 <DTR_REPLICA_1_IP>:443 weight 100 check check-ssl verify none
        server node02 <DTR_REPLICA_2_IP>:443 weight 100 check check-ssl verify none
        server node03 <DTR_REPLICA_N_IP>:443 weight 100 check check-ssl verify none
```
  </div>
  <div id="aws" class="tab-pane fade" markdown="1">
```json
{
    "Subnets": [
        "subnet-XXXXXXXX",
        "subnet-YYYYYYYY",
        "subnet-ZZZZZZZZ"
    ],
    "CanonicalHostedZoneNameID": "XXXXXXXXXXX",
    "CanonicalHostedZoneName": "XXXXXXXXX.us-west-XXX.elb.amazonaws.com",
    "ListenerDescriptions": [
        {
            "Listener": {
                "InstancePort": 443,
                "LoadBalancerPort": 443,
                "Protocol": "TCP",
                "InstanceProtocol": "TCP"
            },
            "PolicyNames": []
        }
    ],
    "HealthCheck": {
        "HealthyThreshold": 2,
        "Interval": 10,
        "Target": "HTTPS:443/_ping",
        "Timeout": 2,
        "UnhealthyThreshold": 4
    },
    "VPCId": "vpc-XXXXXX",
    "BackendServerDescriptions": [],
    "Instances": [
        {
            "InstanceId": "i-XXXXXXXXX"
        },
        {
            "InstanceId": "i-XXXXXXXXX"
        },
        {
            "InstanceId": "i-XXXXXXXXX"
        }
    ],
    "DNSName": "XXXXXXXXXXXX.us-west-2.elb.amazonaws.com",
    "SecurityGroups": [
        "sg-XXXXXXXXX"
    ],
    "Policies": {
        "LBCookieStickinessPolicies": [],
        "AppCookieStickinessPolicies": [],
        "OtherPolicies": []
    },
    "LoadBalancerName": "ELB-DTR",
    "CreatedTime": "2017-02-13T21:40:15.400Z",
    "AvailabilityZones": [
        "us-west-2c",
        "us-west-2a",
        "us-west-2b"
    ],
    "Scheme": "internet-facing",
    "SourceSecurityGroup": {
        "OwnerAlias": "XXXXXXXXXXXX",
        "GroupName":  "XXXXXXXXXXXX"
    }
}
```
  </div>
</div>


You can deploy your load balancer using:

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#nginx-2" data-group="nginx">NGINX</a></li>
  <li><a data-toggle="tab" data-target="#haproxy-2" data-group="haproxy">HAProxy</a></li>
</ul>
<div class="tab-content">
  <div id="nginx-2" class="tab-pane fade in active" markdown="1">
```conf
# Create the nginx.conf file, then
# deploy the load balancer

docker run --detach \
  --name dtr-lb \
  --restart=unless-stopped \
  --publish 80:80 \
  --publish 443:443 \
  --volume ${PWD}/nginx.conf:/etc/nginx/nginx.conf:ro \
  nginx:stable-alpine
```
  </div>
  <div id="haproxy-2" class="tab-pane fade" markdown="1">
```conf
# Create the haproxy.cfg file, then
# deploy the load balancer

docker run --detach \
  --name dtr-lb \
  --publish 443:443 \
  --publish 80:80 \
  --publish 8181:8181 \
  --restart=unless-stopped \
  --volume ${PWD}/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro \
  haproxy:1.7-alpine haproxy -d -f /usr/local/etc/haproxy/haproxy.cfg
```
  </div>
</div>

## Where to go next

- [Backups and disaster recovery](../disaster-recovery/index.md)
- [Monitor and troubleshoot](../monitor-and-troubleshoot/index.md)