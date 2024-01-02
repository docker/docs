---
description: See how you can deploy to Kubernetes on Docker Desktop
keywords: deploy, kubernetes, kubectl, orchestration, Docker Desktop
title: Deploy on Kubernetes with Docker Desktop
aliases:
- /docker-for-windows/kubernetes/
- /docker-for-mac/kubernetes/
---

Docker Desktop includes a standalone Kubernetes server and client,
as well as Docker CLI integration that runs on your machine. 

The Kubernetes server runs locally within your Docker instance, is not configurable, and is a single-node cluster. It runs within a Docker container on your local system, and
is only for local testing. 

Turning on Kubernetes allows you to deploy
your workloads in parallel, on Kubernetes, Swarm, and as standalone containers. Turning on or off the Kubernetes server does not affect your other
workloads.

## Install and turn on Kubernetes

1. From the Docker Dashboard, select the **Settings**.
2. Select **Kubernetes** from the left sidebar. 
3. Next to **Enable Kubernetes**, select the checkbox.
4. Select **Apply & Restart** to save the settings and then select **Install** to confirm. This instantiates images required to run the Kubernetes server as containers, and installs the `/usr/local/bin/kubectl` command on your machine.
   > **Important**
   >
   > The `kubectl` binary is not automatically packaged with Docker Desktop for Linux. To install the kubectl command for Linux, see [Kubernetes documentation](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/). It should be installed at `/usr/local/bin/kubectl`.
   { .important}

By default, Kubernetes containers are hidden from commands like `docker ps`, because managing them manually is not supported. Most users do not need this option. To see these internal containers, select **Show system containers (advanced)**. 

When Kubernetes is turned on and running, an additional status bar in the Docker Dashboard footer and Docker menu displays. 

> **Note**
>
> Docker Desktop does not upgrade your Kubernetes cluster automatically after a new update. To upgrade your Kubernetes cluster to the latest version, select **Reset Kubernetes Cluster**.

## Use the kubectl command

Kubernetes integration provides the Kubernetes CLI command
at `/usr/local/bin/kubectl` on Mac and at `C:\Program Files\Docker\Docker\Resources\bin\kubectl.exe` on Windows. This location may not be in your shell's `PATH`
variable, so you may need to type the full path of the command or add it to
the `PATH`.

If you have already installed `kubectl` and it is
pointing to some other environment, such as `minikube` or a GKE cluster, ensure you change the context so that `kubectl` is pointing to `docker-desktop`:

```console
$ kubectl config get-contexts
$ kubectl config use-context docker-desktop
```

> **Tip**
>
> Run the `kubectl` command in a CMD or PowerShell terminal, otherwise `kubectl config get-contexts` may return an empty result. 
>
> If you are using a different terminal and this happens, you can try setting the `kubeconfig` environment variable to the location of the `.kube/config` file. 
{ .tip }

If you installed `kubectl` using Homebrew, or by some other method, and
experience conflicts, remove `/usr/local/bin/kubectl`.

You can test the command by listing the available nodes:

```console
$ kubectl get nodes

NAME                 STATUS    ROLES     AGE       VERSION
docker-desktop       Ready     master    3h        v1.19.7
```

For more information about `kubectl`, see the
[`kubectl` documentation](https://kubernetes.io/docs/reference/kubectl/overview/).

## Turn off and uninstall Kubernetes

To turn off Kubernetes in Docker Desktop:
1. From the Docker Dashboard, select the **Settings** icon.
2. Select **Kubernetes** from the left sidebar. 
3. Next to **Enable Kubernetes**, clear the checkbox
4. Select **Apply & Restart** to save the settings.This stops and removes Kubernetes containers, and also removes the `/usr/local/bin/kubectl` command.