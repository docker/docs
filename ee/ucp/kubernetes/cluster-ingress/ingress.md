---
title: Deploy a Sample Application with Ingress (Experimental)
description: Learn how to deploy Ingress rules for Kubernetes applications.
keywords: ucp, cluster, ingress, kubernetes
---

>{% include enterprise_label_shortform.md %}

{% include experimental-feature.md %}

## Overview 

Cluster Ingress is capable of routing based on many HTTP attributes, but most
commonly the HTTP host and path. The following example shows the basics of
deploying Ingress rules for a Kubernetes application. An example application is
deployed from this [deployment manifest](./yaml/demo-app.yaml) and L7 Ingress
rules are applied.

## Deploy a Sample Application

In this example, three different versions of the docker-demo application are
deployed. The docker-demo application is able to display the container hostname,
environment variables or labels in its HTTP responses, therefore is good sample
application for an Ingress controller.

The three versions of the sample application are:

- v1: a production version with three replicas running.
- v2: a staging version with a single replica running.
- v3: a development version also with a single replica. 

> Note
> 
> An example Kubernetes manifest file containing all three deployments can be found [here](./yaml/demo-app.yaml).

1. Source a [UCP Client Bundle](/ee/ucp/user-access/cli/) attached to a cluster with Cluster Ingress installed. 
2. Download the sample Kubernetes manifest file.
```bash
$ wget https://raw.githubusercontent.com/docker/docker.github.io/master/ee/ucp/kubernetes/cluster-ingress/yaml/demo-app.yaml
```
3. Deploy the sample Kubernetes manifest file. 
```bash
$ kubectl apply -f demo-app.yaml
```
4. Verify that the sample applications are running.

```bash
  $ kubectl get pods -n default
  NAME                                  READY   STATUS    RESTARTS   AGE
  demo-v1-7797b7c7c8-5vts2              1/1     Running   0          3h
  demo-v1-7797b7c7c8-gfwzj              1/1     Running   0          3h
  demo-v1-7797b7c7c8-kw6gp              1/1     Running   0          3h
  demo-v2-6c5b4c6f76-c6zhm              1/1     Running   0          3h
  demo-v3-d88dddb74-9k7qg               1/1     Running   0          3h

  $ kubectl get services -o wide
  NAME           TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE   SELECTOR
  demo-service   NodePort    10.96.97.215   <none>        8080:33383/TCP   3h    app=demo
  kubernetes     ClusterIP   10.96.0.1      <none>        443/TCP          1d    <none>
```

This first part of the tutorial deployed the pods and a Kubernetes service. 
Using Kubernetes NodePorts, these pods can be accessed outside of the Cluster
Ingress. This illustrates the standard L4 load balancing that a Kubernetes
service applies. 

```bash
# Public IP Address of a Worker or Manager VM in the Cluster
$ IPADDR=51.141.127.241

# Node Port
$ PORT=$(kubectl get service demo-service --output jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')

# Send traffic directly to the NodePort to bypass L7 Ingress
$ for i in {1..5}; do curl http://$IPADDR:$PORT/ping; done
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev"}
{"instance":"demo-v3-d88dddb74-9k7qg","version":"v3","metadata":"dev"}
{"instance":"demo-v2-6c5b4c6f76-c6zhm","version":"v2","metadata":"staging"}
{"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production"}
{"instance":"demo-v1-7797b7c7c8-gfwzj","version":"v1","metadata":"production"}
```

The L4 load balancing is applied to the number of replicas that exist for each
service. Different scenarios require more complex logic to load balancing.
Make sure to detach the number of back-end instances from the load balancing
algorithms used by the Ingress.

## Apply Ingress Rules to the Sample Application

To leverage the Cluster Ingress for the sample application, there are three custom resources types that need to be deployed:
- Gateway
- Virtual Service 
- Destinationrule

> Note
> 
> For the sample application, an example manifest file with all three objects defined can be found [here](./yaml/ingress-simple.yaml).

1. Source a [UCP Client Bundle](/ee/ucp/user-access/cli/) attached to a cluster with Cluster Ingress installed. 
2. Download the sample Kubernetes manifest file.
```bash
$ wget https://raw.githubusercontent.com/docker/docker.github.io/master/ee/ucp/kubernetes/cluster-ingress/yaml/ingress-simple.yaml
```
3. Deploy the sample Kubernetes manifest file.

```bash
  $ kubectl apply -f ingress-simple.yaml

  $ kubectl describe virtualservice demo-vs
  ...
  Spec:
    Gateways:
      cluster-gateway
    Hosts:
      demo.example.com
    Http:
      Match:  <nil>
      Route:
        Destination:
          Host:  demo-service
          Port:
            Number:  8080
          Subset:    v1
```

This configuration matches all traffic with `demo.example.com` and sends it to
the back end version=v1 deployment, regardless of the quantity of replicas in
the back end. 

Curl the service again using the port of the Ingress gateway. Because DNS is
not set up, use the `--header` flag from curl to manually set the host header.

```bash
# Find the Cluster Ingress Node Port
$ PORT=$(kubectl get service -n istio-system istio-ingressgateway --output jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')

# Public IP Address of a Worker or Manager VM in the Cluster
$ IPADDR=51.141.127.241

$ for i in {1..5}; do curl --header "Host: demo.example.com" http://$IPADDR:$PORT/ping; done
{"instance":"demo-v1-7797b7c7c8-5vts2","version":"v1","metadata":"production","request_id":"2558fdd1-0cbd-4ba9-b104-0d4d0b1cef85"}
{"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"59f865f5-15fb-4f49-900e-40ab0c44c9e4"}
{"instance":"demo-v1-7797b7c7c8-5vts2","version":"v1","metadata":"production","request_id":"fe233ca3-838b-4670-b6a0-3a02cdb91624"}
{"instance":"demo-v1-7797b7c7c8-5vts2","version":"v1","metadata":"production","request_id":"842b8d03-8f8a-4b4b-b7f4-543f080c3097"}
{"instance":"demo-v1-7797b7c7c8-kw6gp","version":"v1","metadata":"production","request_id":"197cbb1d-5381-4e40-bc6f-cccec22eccbc"}
```

To have Server Name Indication (SNI) work with TLS services, use curl's `--resolve` flag.

```bash
$ curl --resolve demo.example.com:$IPADDR:$PORT http://demo.example.com/ping
```

In this instance, the three back-end v1 replicas are load balanced and no
requests are sent to the other versions.

## Where to go next

- [Deploy a Sample Application with a Canary release](./canary/)
- [Deploy a Sample Application with Sticky Sessions](./sticky/)
