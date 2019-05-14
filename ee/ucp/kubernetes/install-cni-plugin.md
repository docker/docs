---
title: Install an Unmanaged CNI plugin
description: Learn how to install a Container Networking Interface plugin on Docker Universal Control Plane.
keywords: ucp, Kubernetes, cni, Container Networking Interface, flannel, weave, calico
---

For Docker Universal Control Plane, [Project
Calico](https://docs.projectcalico.org/v3.7/introduction/) provides the secure
networking functionality for the container to container communication within
Kubernetes. The Universal Control Plane handles the lifecycle of Project Calico,
installing and upgrading the project as part of UCP. Additionally, the Project
Calico deployment installed is fully supported, with Docker Support providing
guidance on the [CNI components](https://github.com/projectcalico/cni-plugin).

At install time, the Universal Control Plane can be configured to not install
the Calico CNI Plugin, allowing Platform Operators to install alternative CNI
plugins post installation to support alternative use cases. In this model all
Universal Control Plane components are still fully supported by Docker Support,
however for guidance and support on the 3rd party CNI components, the platform
operator will need to contact that 3rd party.

## Installing a Unmanaged CNI Plugin on Docker UCP

Once a platform operator has followed the [System
Requirements](/ee/ucp/admin/install/system-requirements/) documentation, and
taken into consideration any requirements for the custom CNI plugin. The Docker
UCP installation command can be ran with a `--unmanaged-cni` flag to bring up
the platform.

This command will install the Universal Control Plane, and bring up components
like the UCP UI and the RBAC engine. Universal Control Planes components that
require Kubernetes Networking, such as Metrics, will not start and will stay in
a `Container Creating` state in Kubernetes, until a CNI is installed. 

### Install UCP without a CNI Plugin

Once connected to a manager node, with the Docker Enterprise Engine installed, you are ready to bring up UCP with the `--unmanaged-cni` flag.

```bash
docker container run --rm -it --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <node-ip-address> \
  --unmanaged-cni true \
  --interactive
```

Once the installation is complete, you will be able to access UCP in web
browser. Note the Manager node will be unhealthy as the Kubelet will be
reporting `NetworkPluginNotReady`. Additionally the metrics in the UCP dashboard
will also be unavailable, as this runs in a Kubernetes pods.

### Configure CLI Access to UCP

Next a platform operator should log into UCP, download a UCP Client Bundle, and
configure the Kubernetes CLI tool `kubectl`. For more information see [CLI Based
Access](ee/ucp/user-access/cli/#download-client-certificates)
   
Using `kubectl` you should be able to see that the UCP components that run on
Kubernetes are still pending, waiting for a CNI driver before becoming
available. 

```bash
$ kubectl get nodes
NAME         STATUS     ROLES     AGE       VERSION
manager-01   NotReady   master    10m       v1.11.9-docker-1
  
$ kubectl get pods -n kube-system -o wide
NAME                           READY     STATUS              RESTARTS   AGE       IP        NODE         NOMINATED NODE
compose-565f7cf9ff-gq2gv       0/1       Pending             0          10m       <none>    <none>       <none>
compose-api-574d64f46f-r4c5g   0/1       Pending             0          10m       <none>    <none>       <none>
kube-dns-6d96c4d9c6-8jzv7      0/3       Pending             0          10m       <none>    <none>       <none>
ucp-metrics-nwt2z              0/3       ContainerCreating   0          10m       <none>    manager-01   <none>
```

### Install an Unmanaged CNI Plugin

Using the `kubectl` tool you are now able to install a custom CNI plugin on to
the Universal Control Plane. Alternative CNI plugins could be Weave, Flannel,
Canal, Romana and many more. As Docker Support will not support the CNI plugin,
a platform operator has complete flexibility on what they install.

Common steps when installing a CNI plugin would be to download the relevant
upstream CNI binaries from
https://github.com/containernetworking/cni/releases/tag/, placing them in
`/opt/cni/bin`, additionally downloading the relevant CNI Plugins Kubernetes
Manifest yaml and running `$ kubectl apply -f <your-custom-cni-plugin>.yaml`.
   
Please follow the relevant CNI Plugins documentation to understand the specific
installation instructions.

> While troubleshooting a Custom CNI Plugin, you may wish to access logs
> within the Kubelet. This can be done while connected to the UCP manager with
> `$ docker logs ucp-kubelet`.

### Verify the UCP installation

Once the CNI plugin has been completely installed, you should now see the
Universal Control Plane components that run as pods start to become available.

```
$ kubectl get pods -n kube-system -o wide
NAME                           READY     STATUS    RESTARTS   AGE       IP            NODE         NOMINATED NODE
compose-565f7cf9ff-gq2gv       1/1       Running   0          21m       10.32.0.2     manager-01   <none>
compose-api-574d64f46f-r4c5g   1/1       Running   0          21m       10.32.0.3     manager-01   <none>
kube-dns-6d96c4d9c6-8jzv7      3/3       Running   0          22m       10.32.0.5     manager-01   <none>
ucp-metrics-nwt2z              3/3       Running   0          22m       10.32.0.4     manager-01   <none>
weave-net-wgvcd                2/2       Running   0          8m        172.31.6.95   manager-01   <none>
```

> Note: In this example you can see we have deployed Weave, you may have
> alternative CNI pods in this namespace. 

## Where to go next

- [Make your Cluster Highly Available](https://docs.docker.com/ee/ucp/admin/install/#step-6-join-manager-nodes)
- [Install an Ingress Controller on Kubernetes](ee/ucp/kubernetes/layer-7-routing/)