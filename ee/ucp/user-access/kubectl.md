---
title: Install the Kubernetes CLI
description: Learn how to install kubectl, the Kubernetes command-line tool, on Docker Universal Control Plane.
keywords: ucp, cli, administration, kubectl, Kubernetes
---

Docker EE 2.0 and higher deploys Kubernetes as part of a UCP installation.
Deploy, manage, and monitor Kubernetes workloads from the UCP dashboard. Users can
also interact with the Kubernetes deployment through the Kubernetes
command-line tool named kubectl.

To access the UCP cluster with kubectl, install the [UCP client bundle](cli.md).

> Kubernetes on Docker Desktop for Mac and Docker Desktop for Windows
>
> Docker Desktop for Mac and Docker Desktop for Windows provide a standalone Kubernetes server that
> runs on your development machine, with kubectl installed by default. This installation is
> separate from the Kubernetes deployment on a UCP cluster.
> Learn how to [deploy to Kubernetes on Docker Desktop for Mac](/docker-for-mac/kubernetes.md).
{: .important}

## Install the kubectl binary

To use kubectl, install the binary on a workstation which has access to your UCP endpoint.

> Must install compatible version
>
> Kubernetes only guarantees compatibility with kubectl versions that are +/-1 minor versions away from the Kubernetes version.
{: .important}

First, find which version of Kubernetes is running in your cluster. This can be found
within the Universal Control Plane dashboard or at the UCP API endpoint [version](/reference/ucp/3.2/api/).

From the UCP dashboard, click on **About Docker EE** within the **Admin** menu in the top left corner
 of the dashboard. Then navigate to **Kubernetes**.

 ![Find Kubernetes version](../images/kubernetes-version.png){: .with-border}

Once you have the Kubernetes version, install the kubectl client for the relevant
operating system.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#mac">Mac OS</a></li>
  <li><a data-toggle="tab" data-target="#linux">Linux</a></li>
  <li><a data-toggle="tab" data-target="#win">Windows</a></li>
</ul>
<div class="tab-content">
<div id="mac" class="tab-pane fade in active" markdown="1">
```
# Set the Kubernetes version as found in the UCP Dashboard or API
k8sversion=v1.11.5

# Get the kubectl binary.
curl -LO https://storage.googleapis.com/kubernetes-release/release/$k8sversion/bin/darwin/amd64/kubectl

# Make the kubectl binary executable.
chmod +x ./kubectl

# Move the kubectl executable to /usr/local/bin.
sudo mv ./kubectl /usr/local/bin/kubectl
```
<hr>
</div>
<div id="linux" class="tab-pane fade" markdown="1">
```
# Set the Kubernetes version as found in the UCP Dashboard or API
k8sversion=v1.11.5

# Get the kubectl binary.
curl -LO https://storage.googleapis.com/kubernetes-release/release/$k8sversion/bin/linux/amd64/kubectl

# Make the kubectl binary executable.
chmod +x ./kubectl

# Move the kubectl executable to /usr/local/bin.
sudo mv ./kubectl /usr/local/bin/kubectl
```
<hr>
</div>
<div id="win" class="tab-pane fade" markdown="1">
You can download the binary from this [link](https://storage.googleapis.com/kubernetes-release/release/v.1.8.11/bin/windows/amd64/kubectl.exe)

If you have curl installed on your system, you use these commands in Powershell.

```cmd
$env:k8sversion = "v1.11.5"

curl https://storage.googleapis.com/kubernetes-release/release/$env:k8sversion/bin/windows/amd64/kubectl.exe
```
<hr>
</div>
</div>

## Using kubectl with a Docker EE cluster

Docker Enterprise Edition provides users unique certificates and keys to authenticate against
 the Docker and Kubernetes APIs. Instructions on how to download these certificates and how to
 configure kubectl to use them can be found in [CLI-based access.](cli.md#download-client-certificates)

## Install Helm on Docker EE

This section describes how to grant the correct service account/permissions on tiller to install applications into the cluster.

### Prerequisites
Before performing these steps, you must meet the following requirements:
- You must be running a Docker EE 2.1 or higher cluster
- You must have kubectl configured to communicate with the cluster (usually this is done via a client bundle)

To create a service account for tiller and apply the correct grants:

1. Create the tiller service account:

```bash
$ kubectl create serviceaccount --namespace kube-system tiller
```

2. In the UCP UI, navigate to **Access Control > Grants**, and click **Create Role Binding**. Make sure you have selected Kubernetes first.

3. In Subject, select the following:
- Select Subject Type as **Service Account**
- Namespace as **kube-system**
- Service Account as **tiller**

4. Click **Next**.

5. In the Resource Set area, move the slider to enable **Apply Role Binding to all namespace (Cluster Role Binding)**.

6. Click **Next**.

7. In the Role area, select **cluster-admin**.

8. Click **Create**.

9.  At this stage, if you have tiller installed in the cluster already, you need to patch the tiller deployment to use the tiller service account you created:

```bash
$ kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
```

However, if you do not currently have tiller installed, you can install and direct it to use the correct service account at the same time using the following command:

```bash
$ helm init --service-account tiller
```

10. After waiting for the tiller pod to be (re-)created, try to run helm install again:

```bash
$ helm install --name mysql stable/mysql
```

## Where to go next

- [Deploy a workload to a Kubernetes cluster](../kubernetes.md)
- [Deploy to Kubernetes on Docker Desktop for Mac](/docker-for-mac/kubernetes.md)
