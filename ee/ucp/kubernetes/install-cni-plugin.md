---
title: Install a CNI plugin
description: Learn how to install a Container Networking Interface plugin on Docker Universal Control Plane.
keywords: ucp, cli, administration, kubectl, Kubernetes, cni, Container Networking Interface, flannel, weave, ipip, calico
---

For Docker Universal Control Plane, [Project Calico](https://docs.projectcalico.org/v3.0/introduction/) 
provides the secure networking functionality for the container communication with Kubernetes.

Docker EE supports Calico and installs the 
built-in [Calico](https://github.com/projectcalico/cni-plugin) plugin, but you can override that and 
install a Docker certified plugin.

> **Note**: The `--cni-installer-url` option is deprecated as of UCP 3.1. It is replaced by the `--unmanaged-cni` option.

# Install UCP with a custom CNI plugin

Modify the [UCP install command-line](../admin/install/index.md#step-4-install-ucp)
to add the `--cni-installer-url` [option](/reference/ucp/3.0/cli/install.md),
providing a URL for the location of the CNI plugin's YAML file:

```bash
docker container run --rm -it --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <node-ip-address> \
   --unmanaged-cni <true|false> \
  --interactive
```

> **Note**: Setting `--unmanaged-cni` to `true` value installs UCP without a managed CNI plugin. UCP and the
> Kubernetes components will be running but pod-to-pod networking will not function until a CNI plugin is manually
> installed. This will impact some functionality of UCP until a CNI plugin is running.

You must provide a correct YAML installation file for the CNI plugin, but most
of the default files work on Docker EE with no modification.

## YAML files for CNI plugins

Use the following commands to get the YAML files for popular CNI plugins. 

- [Flannel](https://github.com/coreos/flannel)
  ```bash
  # Get the URL for the Flannel CNI plugin.
  CNI_URL="https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml"
  ```
- [Weave](https://www.weave.works/)
  ```bash
  # Get the URL for the Weave CNI plugin. 
  CNI_URL="https://cloud.weave.works/k8s/net?k8s-version=Q2xpZW50IFZlcnNpb246IHZlcnNpb24uSW5mb3tNYWpvcjoiMSIsIE1pbm9yOiI5IiwgR2l0VmVyc2lvbjoidjEuOS4zIiwgR2l0Q29tbWl0OiJkMjgzNTQxNjU0NGYyOThjOTE5ZTJlYWQzYmUzZDA4NjRiNTIzMjNiIiwgR2l0VHJlZVN0YXRlOiJjbGVhbiIsIEJ1aWxkRGF0ZToiMjAxOC0wMi0wN1QxMjoyMjoyMVoiLCBHb1ZlcnNpb246ImdvMS45LjIiLCBDb21waWxlcjoiZ2MiLCBQbGF0Zm9ybToibGludXgvYW1kNjQifQpTZXJ2ZXIgVmVyc2lvbjogdmVyc2lvbi5JbmZve01ham9yOiIxIiwgTWlub3I6IjgrIiwgR2l0VmVyc2lvbjoidjEuOC4yLWRvY2tlci4xNDMrYWYwODAwNzk1OWUyY2UiLCBHaXRDb21taXQ6ImFmMDgwMDc5NTllMmNlYWUxMTZiMDk4ZWNhYTYyNGI0YjI0MjBkODgiLCBHaXRUcmVlU3RhdGU6ImNsZWFuIiwgQnVpbGREYXRlOiIyMDE4LTAyLTAxVDIzOjI2OjE3WiIsIEdvVmVyc2lvbjoiZ28xLjguMyIsIENvbXBpbGVyOiJnYyIsIFBsYXRmb3JtOiJsaW51eC9hbWQ2NCJ9Cg=="
  ```
  If you have kubectl available, for example by using
  [Docker Desktop for Mac](/docker-for-mac/kubernetes.md), you can use the following
  command to get the URL for the [Weave](https://www.weave.works/) CNI plugin:
  ```bash
  # Get the URL for the Weave CNI plugin. 
  CNI_URL="https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
  ```
- [Romana](http://docs.romana.io/)
  ```bash
  # Get the URL for the Romana CNI plugin.
  CNI_URL="https://raw.githubusercontent.com/romana/romana/master/docs/kubernetes/romana-kubeadm.yml"
  ```

## Disable IP in IP overlay tunneling

The Calico CNI plugin supports both overlay (IPIP) and underlay forwarding
technologies. By default, Docker UCP uses IPIP overlay tunneling.

If you're used to managing applications at the network level through the 
underlay visibility, or you want to reuse existing networking tools in the
underlay, you may want to disable the IPIP functionality. Run the following
commands on the Kubernetes master node to disable IPIP overlay tunneling.

```bash
# Exec into the Calico Kubernetes controller container.
docker exec -it $(docker ps --filter name=k8s_calico-kube-controllers_calico-kube-controllers -q) sh

# Download calicoctl
wget https://github.com/projectcalico/calicoctl/releases/download/v3.1.1/calicoctl && chmod +x calicoctl

# Get the IP pool configuration. 
./calicoctl get ippool -o yaml > ippool.yaml

# Edit the file: Disable IPIP in ippool.yaml by setting "ipipMode: Never".

# Apply the edited file to the Calico plugin.
./calicoctl apply -f ippool.yaml

```

These steps disable overlay tunneling, and Calico uses the underlay networking,
in environments where it's supported.

## Where to go next

- [Install UCP for production](../admin/install.md)
- [Deploy a workload to a Kubernetes cluster](../kubernetes.md)
