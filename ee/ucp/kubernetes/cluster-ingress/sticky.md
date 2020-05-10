---
title: Deploy a Sample Application with Sticky Sessions (Experimental)
description: Learn how to use cookies with Ingress host and path routing.
keywords: ucp, cluster, ingress, kubernetes
---

>{% include enterprise_label_shortform.md %}

{% include experimental-feature.md %}

## Overview

With persistent sessions, the Ingress controller can use a predetermined header
or dynamically generate a HTTP header cookie for a client session to use, so
that a clients requests are sent to the same back end.

> Note
> 
> This guide assumes the [Deploy Sample Application](ingress.md)
> tutorial was followed, with the artifacts still running on the cluster. If
> they are not, please go back and follow this guide.

This is specified within the Istio Object `DestinationRule` via a
`TrafficPolicy` for a given host. In the following example configuration,
consistentHash is chosen as the load balancing method and a cookie named
“session” is used to determine the consistent hash. If incoming requests do not
have the “session” cookie set, the Ingress proxy sets it for use in future
requests.

> Note 
> 
> A Kubernetes manifest file with an updated DestinationRule can be found [here](./yaml/ingress-sticky.yaml).

1. Source a [UCP Client Bundle](/ee/ucp/user-access/cli/) attached to a cluster with Cluster Ingress installed. 
2. Download the sample Kubernetes manifest file.
```bash
$ wget https://raw.githubusercontent.com/docker/docker.github.io/master/ee/ucp/kubernetes/cluster-ingress/yaml/ingress-sticky.yaml
```
3. Deploy the Kubernetes manifest file with the new DestinationRule. This file includes the consistentHash loadBalancer policy. 
```bash
$ kubectl apply -f ingress-sticky.yaml
```
4. Curl the service to view how requests are load balanced without using cookies. In this example, requests are bounced between different v1 services.

```bash
  # Public IP Address of a Worker or Manager VM in the Cluster
  $ IPADDR=51.141.127.241

  # Node Port
  $ PORT=$(kubectl get service demo-service --output jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')

  $ for i in {1..5}; do curl -H "Host: demo.example.com" http://$IPADDR:$PORT/ping; done
  {"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production","request_id":"b40a0294-2629-413b-b876-76b59d72189b"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"721fe4ba-a785-484a-bba0-627ee6e47188"}
  {"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production","request_id":"77ed801b-81aa-4c02-8cc9-7e3bd3244807"}
  {"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production","request_id":"36d8aaed-fcdf-4489-a85e-76ea96949d6c"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"4693b6ad-286b-4470-9eea-c8656f6801ae"}
```

Now curl again and inspect the headers returned from the proxy.

```bash
  $ curl -i -H "Host: demo.example.com" http://$IPADDR:$PORT/ping
  HTTP/1.1 200 OK
  set-cookie: session=1555389679134464956; Path=/; Expires=Wed, 17 Apr 2019 04:41:19 GMT; Max-Age=86400
  date: Tue, 16 Apr 2019 04:41:18 GMT
  content-length: 131
  content-type: text/plain; charset=utf-8
  x-envoy-upstream-service-time: 0
  set-cookie: session="d7227d32eeb0524b"; Max-Age=60; HttpOnly
  server: istio-envoy

  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"011d5fdf-2285-4ce7-8644-c2df6481c584"}
```

The Ingress proxy sets a 60 second TTL cookie named `session` on this HTTP request. A browser or other client application can use that value in future
requests.

Now curl the service again using the flags that save cookies persistently across sessions. The header information shows the session is being set,
persisted across requests, and that for a given session header, the responses are coming from the same back end.

```bash
  $ for i in {1..5}; do curl -c cookie.txt -b cookie.txt -H "Host: demo.example.com" http://$IPADDR:$PORT/ping; done
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"72b35296-d6bd-462a-9e62-0bd0249923d7"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"c8872f6c-f77c-4411-aed2-d7aa6d1d92e9"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"0e7b8725-c550-4923-acea-db94df1eb0e4"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"9996fe77-8260-4225-89df-0eaf7581e961"}
  {"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"d35c380e-31d6-44ce-a5d0-f9f6179715ab"}
```

When the HTTP uses the cookie that is set by the Ingress proxy, all requests are sent to the same back end, `demo-v1-7797b7c7c8-kw6gp`.

## Where to go next

- [Cluster Ingress Overview](ingress.md)
