---
description: Deploying to Kubernetes on Docker Desktop
keywords: deploy, kubernetes, kubectl, orchestration
title: Deploy on Kubernetes
redirect_from:
- /docker-for-windows/kubernetes/
- /docker-for-mac/kubernetes/
---

Docker Desktop includes a standalone Kubernetes server and client,
as well as Docker CLI integration that runs on your machine. 

The Kubernetes server runs locally within your Docker instance, is not configurable, and is a single-node cluster. It runs within a Docker container on your local system, and
is only for local testing. 

Enabling Kubernetes allows you to deploy
your workloads in parallel, on Kubernetes, Swarm, and as standalone containers. Enabling or disabling the Kubernetes server does not affect your other
workloads.

## Enable Kubernetes

To enable Kubernetes in Docker Desktop:
1. From the Docker Dashboard, select the **Setting** icon, or **Preferences** icon if you use a macOS.
2. Select **Kubernetes** from the left sidebar. 
3. Next to **Enable Kubernetes**, select the checkbox
4. Select **Apply & Restart** to save the settings and then click **Install** to confirm. This instantiates images required to run the Kubernetes server as containers, and installs the `/usr/local/bin/kubectl` command on your machine.

By default, Kubernetes containers are hidden from commands like `docker ps`, because managing them manually is not supported. Most users do not need this option. To see these internal containers, select **Show system containers (advanced)**. 

When Kubernetes is enabled and running, an additional status bar in the Dashboard footer and Docker menu displays. 

> Note
>
> Docker Desktop does not upgrade your Kubernetes cluster automatically after a new update. To upgrade your Kubernetes cluster to the latest version, select **Reset Kubernetes Cluster**.

## Use the kubectl command

Kubernetes integration provides the Kubernetes CLI command
at `/usr/local/bin/kubectl` on Mac and at `C:\>Program Files\Docker\Docker\Resources\bin\kubectl.exe` on Windows. This location may not be in your shell's `PATH`
variable, so you may need to type the full path of the command or add it to
the `PATH`.

If you have already installed `kubectl` and it is
pointing to some other environment, such as `minikube` or a GKE cluster, ensure you change the context so that `kubectl` is pointing to `docker-desktop`:

```console
$ kubectl config get-contexts
$ kubectl config use-context docker-desktop
```

If you installed `kubectl` using Homebrew, or by some other method, and
experience conflicts, remove `/usr/local/bin/kubectl`.

You can test the command by listing the available nodes:

```console
$ kubectl get nodes

NAME                 STATUS    ROLES     AGE       VERSION
docker-desktop       Ready     master    3h        v1.19.7
```

For more information about `kubectl`, see the
[`kubectl` documentation](https://kubernetes.io/docs/reference/kubectl/overview/){:target="_blank" rel="noopener" class="_"}.

## Disable Kubernetes

To disable Kubernetes in Docker Desktop:
1. From the Docker Dashboard, select the **Setting** icon, or **Preferences** icon if you use a macOS.
2. Select **Kubernetes** from the left sidebar. 
3. Next to **Enable Kubernetes**, clear the checkbox
4. Select **Apply & Restart** to save the settings.This stops and removes Kubernetes containers, and also removes the `/usr/local/bin/kubectl` command.
