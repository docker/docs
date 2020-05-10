---
title: Deploy a Sample Application with a Canary release (Experimental)
description: Stage a canary release using weight-based load balancing between multiple back-end applications.
keywords: ucp, cluster, ingress, kubernetes
---

>{% include enterprise_label_shortform.md %}

{% include experimental-feature.md %}

## Overview 

This example stages a canary release using weight-based load balancing between
multiple back-end applications.

> Note
> 
> This guide assumes the [Deploy Sample Application](ingress.md)
> tutorial was followed, with the artifacts still running on the cluster. If
> they are not, please go back and follow this guide.

The following schema is used for this tutorial:
- 80% of client traffic is sent to the production v1 service.
- 20% of client traffic is sent to the staging v2 service.
- All test traffic using the header `stage=dev` is sent to the v3 service.

A new Kubernetes manifest file with updated ingress rules can be found [here](./yaml/ingress-weighted.yaml).

1. Source a [UCP Client Bundle](/ee/ucp/user-access/cli/) attached to a cluster with Cluster Ingress installed. 
2. Download the sample Kubernetes manifest file.
```bash
$ wget https://raw.githubusercontent.com/docker/docker.github.io/master/ee/ucp/kubernetes/cluster-ingress/yaml/ingress-weighted.yaml
```
3. Deploy the Kubernetes manifest file.
 
```bash
  $ kubectl apply -f ingress-weighted.yaml
  
  $ kubectl describe vs
   Hosts:
      demo.example.com
    Http:
      Match:
        Headers:
          Stage:
            Exact:  dev
      Route:
        Destination:
          Host:  demo-service
          Port:
            Number:  8080
          Subset:    v3
      Route:
        Destination:
          Host:  demo-service
          Port:
            Number:  8080
          Subset:    v1
        Weight:      80
        Destination:
          Host:  demo-service
          Port:
            Number:  8080
          Subset:    v2
        Weight:      20
```

This virtual service performs the following actions:
- Receives all traffic with host=demo.example.com.
- If an exact match for HTTP header `stage=dev` is found, traffic is routed to v3.
- All other traffic is routed to v1 and v2 in an 80:20 ratio.

Now we can send traffic to the application to view the applied load balancing
algorithms.

```bash
# Public IP Address of a Worker or Manager VM in the Cluster
$ IPADDR=51.141.127.241

# Node Port
$ PORT=$(kubectl get service demo-service --output jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')

$ for i in {1..5}; do curl -H "Host: demo.example.com" http://$IPADDR:$PORT/ping; done
{"instance":"demo-v1-7797b7c7c8-5vts2","version":"v1","metadata":"production","request_id":"d0671d32-48e7-41f7-a358-ddd7b47bba5f"}
{"instance":"demo-v2-6c5b4c6f76-c6zhm","version":"v2","metadata":"staging","request_id":"ba6dcfd6-f62a-4c68-9dd2-b242179959e0"}
{"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"d87601c0-7935-4cfc-842c-37910e6cd573"}
{"instance":"demo-v1-7797b7c7c8-5vts2","version":"v1","metadata":"production","request_id":"4c71ffab-8657-4d99-87b3-7a6933258990"}
{"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production","request_id":"c404471c-cc85-497e-9e5e-7bb666f4f309"}
```

The split between v1 and v2 corresponds to the specified criteria. Within the
v1 service, requests are load-balanced across the three back-end replicas. v3 does
not appear in the requests.

To send traffic to the third service, add the HTTP header `stage=dev`.

```bash
for i in {1..5}; do curl -H "Host: demo.example.com" -H "Stage: dev" http://$IPADDR:$PORT/ping; done
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev","request_id":"52d7afe7-befb-4e17-a49c-ee63b96d0daf"}
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev","request_id":"b2e664d2-5224-44b1-98d9-90b090578423"}
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev","request_id":"5446c78e-8a77-4f7e-bf6a-63184db5350f"}
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev","request_id":"657553c5-bc73-4a13-b320-f78f7e6c7457"}
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev","request_id":"bae52f09-0510-42d9-aec0-ca6bbbaae168"}
```

In this case, 100% of the traffic with the `stage=dev` header is sent to the v3 service.

## Where to go next

- [Deploy the Sample Application with Sticky Sessions](sticky.md)
