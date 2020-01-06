---
title: Install Cluster Ingress (Experimental)
description: Learn how to deploy ingress rules using Kubernetes manifests.
keywords: ucp, cluster, ingress, kubernetes
---

>{% include enterprise_label_shortform.md %}

{% include experimental-feature.md %}

## Overview

Cluster Ingress for Kubernetes is currently deployed manually outside of UCP.
Future plans for UCP include managing the full lifecycle of the Ingress
components themselves. This guide describes how to manually deploy Ingress using
Kubernetes deployment manifests.

## Offline Installation

If you are installing Cluster Ingress on a UCP cluster that does not have access
to the Docker Hub, you will need to pre-pull the Ingress container images. If
your cluster has access to the Docker Hub, you can move on to [deploying cluster
ingress](#deploy-cluster-ingress).

Without access to the Docker Hub, you will need to download the container images
on a workstation with access to the internet. Container images are distributed
in a `.tar.gz` and can be downloaded from 
[here](https://s3.amazonaws.com/docker-istio/istio-ingress-1.1.2.tgz).

Once the container images have been downloaded, they would then need to be
copied on to the hosts in your UCP cluster, and then side loaded in Docker.
Images can be side loaded with:

```bash
$ docker load -i ucp.tar.gz
```

There images should now be present on your nodes:

```bash
$ docker images
REPOSITORY              TAG                 IMAGE ID            CREATED             SIZE
docker/node-agent-k8s   1.1.2               4ddd06d05d5d        6 days ago          243MB
docker/proxy_init       1.1.2               ff9628f32621        6 days ago          145MB
docker/proxyv2          1.1.2               bebabbe114a4        6 days ago          360MB
docker/pilot            1.1.2               58b6e18f3545        6 days ago          299MB
```

## Deploy Cluster Ingress

This step deploys the Ingress controller components `istio-pilot` and
`istio-ingressgateway`. Together, these components act as the control-plane and
data-plane for ingress traffic. These components are a simplified deployment of
Istio cluster Ingress functionality. Many other custom Kubernetes resources (CRDs) are
also created that aid in the Ingress functionality. 

> Note
> 
> This does not deploy the service mesh capabilities of Istio as its
> function in UCP is for Ingress.

> Note
> 
> As Cluster Ingress is not built into UCP in this release, a Cluster Admin will
> need to manually download and apply the following Kubernetes Manifest [file](https://s3.amazonaws.com/docker-istio/istio-ingress-1.1.2.yaml). 

1. Download the Kubernetes manifest file.
```bash
$ wget https://s3.amazonaws.com/docker-istio/istio-ingress-1.1.2.yaml
```
2. Source a [UCP Client Bundle](/ee/ucp/user-access/cli/).
3. Deploy the Kubernetes manifest file.
```bash
$ kubectl apply -f istio-ingress-1.1.2.yaml
```
4. Verify that the installation was successful. It may take 1-2 minutes for all pods to become ready.

```bash
$ kubectl get pods -n istio-system -o wide
NAME                                    READY   STATUS    RESTARTS   AGE   IP           NODE         NOMINATED NODE   READINESS GATES
istio-ingressgateway-747bc6b4cb-fkt6k   2/2     Running   0          44s   172.0.1.23   manager-02   <none>           <none>
istio-ingressgateway-747bc6b4cb-gr8f7   2/2     Running   0          61s   172.0.1.25   manager-02   <none>           <none>
istio-pilot-7b74c7568b-ntbjd            1/1     Running   0          61s   172.0.1.22   manager-02   <none>           <none>
istio-pilot-7b74c7568b-p5skc            1/1     Running   0          44s   172.0.1.24   manager-02   <none>           <none>

$ kubectl get services -n istio-system -o wide
NAME                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                                      AGE   SELECTOR
istio-ingressgateway   NodePort    10.96.32.197    <none>        80:33000/TCP,443:33001/TCP,31400:33002/TCP,15030:34420/TCP,15443:34368/TCP,15020:34300/TCP   86s   app=istio-ingressgateway,istio=ingressgateway,release=istio
istio-pilot            ClusterIP   10.96.199.152   <none>        15010/TCP,15011/TCP,8080/TCP,15014/TCP                                                       85s   istio=pilot
```

Now you can test the Ingress deployment. To test that the Envoy proxy is working correctly in the Istio Gateway pods, there is a status port configured on an internal port 15020. From the above output, we can see that port 15020 is exposed as a Kubernetes NodePort. In the output above, the NodePort is 34300, but this could be different in each
environment. 

To check the envoy proxy's status, there is a health endpoint at `/healthz/ready`.

```bash
# Node Port
$ PORT=$(kubectl get service -n istio-system istio-ingressgateway --output jsonpath='{.spec.ports[?(@.name=="status-port")].nodePort}')

# Public IP Address of a Worker or Manager VM in the Cluster
$ IPADDR=51.141.127.241

# Use Curl to check the status port is available
$ curl -vvv http://$IPADDR:$PORT/healthz/ready
*   Trying 51.141.127.241...
* TCP_NODELAY set
* Connected to 51.141.127.241 (51.141.127.241) port 34300 (#0)
> GET /healthz/ready HTTP/1.1
> Host: 51.141.127.241:34300
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Wed, 19 Jun 2019 13:31:53 GMT
< Content-Length: 0
<
* Connection #0 to host 51.141.127.241 left intact
```

If the output is `HTTP/1.1 200 OK`, then Envoy is running correctly, ready to service applications.

## Where to go next

- [Deploy a Sample Application](./ingress/)
