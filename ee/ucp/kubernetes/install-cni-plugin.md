---
title: Install an unmanaged CNI plugin
description: Learn how to install a Container Networking Interface (CNI) plugin on Docker Universal Control Plane.
keywords: ucp, kubernetes, cni, container networking interface, flannel, weave, calico
---

>{% include enterprise_label_shortform.md %}

For Docker Universal Control Plane (UCP), [Calico](https://docs.projectcalico.org/v3.7/introduction/) 
provides the secure networking functionality for container-to-container communication within
Kubernetes. UCP handles the lifecycle of Calico and packages it with UCP
installation and upgrade. Additionally, the Calico deployment included with
UCP is fully supported with Docker providing guidance on the 
[CNI components](https://github.com/projectcalico/cni-plugin).

At install time, UCP can be configured to install an alternative CNI plugin 
to support alternative use cases. The alternative CNI plugin may be certified by 
Docker and its partners, and published on Docker Hub. UCP components are still 
fully supported by Docker and respective partners. Docker will provide
pointers to basic configuration, however for additional guidance on managing third-party 
CNI components, the platform operator will need to refer to the partner documentation 
or contact that third party.

UCP does manage the version or configuration of alternative CNI plugins. UCP
upgrade will not upgrade or reconfigure alternative CNI plugins. To switch
between managed and unmanaged CNI plugins or vice versa, you must uninstall and
then reinstall UCP.

## Install an unmanaged CNI plugin on Docker UCP

Once a platform operator has complied with [UCP system
requirements](/ee/ucp/admin/install/system-requirements/) and
taken into consideration any requirements for the custom CNI plugin, you can 
[run the UCP install command](/reference/ucp/3.1/cli/install/) with the `--unmanaged-cni` flag
to bring up the platform.

This command will install UCP, and bring up components
like the user interface and the RBAC engine. UCP components that
require Kubernetes Networking, such as Metrics, will not start and will stay in
a `Container Creating` state in Kubernetes, until a CNI is installed. 

### Install UCP without a CNI plugin

Once connected to a manager node with the Docker Enterprise Engine installed, 
you are ready to install UCP with the `--unmanaged-cni` flag.

```bash
docker container run --rm -it --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <node-ip-address> \
  --unmanaged-cni \
  --interactive
```

Once the installation is complete, you will be able to access UCP in the browser. 
Note that the manager node will be unhealthy as the kubelet will 
report `NetworkPluginNotReady`. Additionally, the metrics in the UCP dashboard
will also be unavailable, as this runs in a Kubernetes pod.

### Configure CLI access to UCP

Next, a platform operator should log into UCP, download a UCP client bundle, and
configure the Kubernetes CLI tool, `kubectl`. See [CLI Based
Access](/ee/ucp/user-access/cli/#download-client-certificates) for more details.
   
With `kubectl`, you can see that the UCP components running on
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

### Install an unmanaged CNI plugin

You can use`kubectl` to install a custom CNI plugin on UCP. 
Alternative CNI plugins are Weave, Flannel, Canal, Romana and many more. 
Platform operators have complete flexibility on what to install, but Docker 
will not support the CNI plugin.

The steps for installing a CNI plugin typically include: 
- Downloading the relevant upstream CNI binaries from
https://github.com/containernetworking/cni/releases/tag/
- Placing them in `/opt/cni/bin`
- Downloading the relevant CNI plugin's Kubernetes Manifest YAML, and 
- Running `$ kubectl apply -f <your-custom-cni-plugin>.yaml`
   
Follow the CNI plugin documentation for specific installation 
instructions.

> Note
> 
> While troubleshooting a custom CNI plugin, you may wish to access logs
> within the kubelet. Connect to a UCP manager node and run
> `$ docker logs ucp-kubelet`.

### Verify the UCP installation

Upon successful installation of the CNI plugin, the related UCP components should have
a `Running` status as pods start to become available.

```
$ kubectl get pods -n kube-system -o wide
NAME                           READY     STATUS    RESTARTS   AGE       IP            NODE         NOMINATED NODE
compose-565f7cf9ff-gq2gv       1/1       Running   0          21m       10.32.0.2     manager-01   <none>
compose-api-574d64f46f-r4c5g   1/1       Running   0          21m       10.32.0.3     manager-01   <none>
kube-dns-6d96c4d9c6-8jzv7      3/3       Running   0          22m       10.32.0.5     manager-01   <none>
ucp-metrics-nwt2z              3/3       Running   0          22m       10.32.0.4     manager-01   <none>
weave-net-wgvcd                2/2       Running   0          8m        172.31.6.95   manager-01   <none>
```

> Note
>
> The above example deployment uses Weave. If you are using an alternative 
> CNI plugin, look for the relevant name and review its status.

## Where to go next

- [Make your Cluster Highly Available](../admin/install/index.md#step-6-join-manager-nodes)
- [Deploy a Sample Application with Ingress](cluster-ingress/ingress.md)
