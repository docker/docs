---
title: Use a load balancer
description: Learn how to configure a load balancer to balance user requests across multiple Docker Trusted Registry replicas.
keywords: docker, dtr, load balancer
---

Once you’ve joined multiple DTR replicas nodes for high-availability, you can
configure your own load balancer to balance user requests across all replicas.

![](../../images/use-a-load-balancer-1.svg)


This allows users to access DTR using a centralized domain name. If a replica
goes down, the load balancer can detect that and stop forwarding requests to
it, so that the failure goes unnoticed by users.

## Load-balancing DTR

DTR does not provide a load balancing service. You can use an on-premises
or cloud-based load balancer to balance requests across multiple DTR replicas.

Make sure you configure your load balancer to:

* Load-balance TCP traffic on ports 80 and 443 (default)
* Not terminate HTTPS connections. Load balancer needs to be TCP-passthrough.
* Use the `/health` endpoint on each DTR replica, to check if
the replica is healthy or not

> Note: If you choose non-default TCP ports for HTTP and HTTPS, ensure that your load balancer forwarding logic and healthchecks match the chosen ports.


## Load-Balancing Configuration Examples

Below are sample **nginx**, **haproxy**, and **AWS Elastic Loadbalancer** configuration files that you can use to set up your DTR external load-balancer.

1. NGINX

	Here is a sample NGINX config. You need to replace the DTR_REPLICA_IPs with your DTR replica IPs. Also if you use non-standard HTTP/HTTPS ports, you need to use those instead. 
	
	```
	user  nginx;
	worker_processes  1;
	
	error_log  /var/log/nginx/error.log warn;
	pid        /var/run/nginx.pid;
	
	events {
	    worker_connections  1024;
	}
	
	stream {
	    upstream dtr_443 {
	        server <DTR_REPLICA_1_IP>:80  max_fails=2 fail_timeout=30s;
	        server <DTR_REPLICA_2_IP>:80  max_fails=2 fail_timeout=30s;
	        server <DTR_REPLICA_N_IP>:80   max_fails=2 fail_timeout=30s;
	    }
	    upstream dtr_80 {
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
	
	You can easily launch NGINX as a container on your designated LB machine using the following command. This assumes that you have `nginx.conf` in current directory. 
	
	```
	docker run -d -p 80:80 -p 443:443 --restart=unless-stopped --name dtr-lb -v ${PWD}/nginx.conf:/etc/nginx/nginx.conf:ro nginx:stable-alpine
	```


2. HAPROXY

	Here is a sample HAProxy config. You need to change the DTR_REPLICA_IPs and the DTR_FQDN. Also if you use non-standard HTTP/HTTPS ports, you need to use those instead.
	
	
	```
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
	        option httpchk GET /health HTTP/1.1\r\nHost:\ <DTR_FQDN>
	        server node01 <DTR_REPLICA_1_IP>:80 check weight 100
	        server node02 <DTR_REPLICA_2_IP>:80 check weight 100
	        server node03 <DTR_REPLICA_N_IP>:80 check weight 100
	backend dtr_upstream_servers_443
	        mode tcp
	        option httpchk GET /health HTTP/1.1\r\nHost:\ <DTR_FQDN>
	        server node01 <DTR_REPLICA_1_IP>:443 weight 100 check ssl verify none
	        server node02 <DTR_REPLICA_2_IP>:443 weight 100 check ssl verify none
	        server node03 <DTR_REPLICA_N_IP>:443 weight 100 check ssl verify none
	```
	
	You can easily launch HAProxy as a container on your designated LB machine using the following command. This assumes that you have `haproxy.cfg` in current directory.
	
	
	```
	docker run -d -p 443:443 -p 80:80 -p 8181:8181 --restart=unless-stopped --name dtr-lb -v ${PWD}/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro haproxy:1.7-alpine haproxy -d -f /usr/local/etc/haproxy/haproxy.cfg
	```

3. AWS ELB

	Here is a sample configuration for DTR's ELB. You can use aws cli or Console when configuring the ELB.
	
	```
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
	                "Target": "HTTPS:443/health",
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
	        },
	        
	```

## Where to go next

* [Backups and disaster recovery](../backups-and-disaster-recovery.md)
* [DTR architecture](../../architecture.md)
