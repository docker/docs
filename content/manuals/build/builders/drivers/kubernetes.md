---
title: Kubernetes driver
description: |
  The Kubernetes driver lets you run BuildKit in a Kubernetes cluster.
  You can connect to, and run your builds in, the cluster using Buildx.
keywords: build, buildx, driver, builder, kubernetes
aliases:
  - /build/buildx/drivers/kubernetes/
  - /build/building/drivers/kubernetes/
  - /build/drivers/kubernetes/
---

The Kubernetes driver lets you connect your local development or CI
environments to builders in a Kubernetes cluster to allow access to more
powerful compute resources, optionally on multiple native architectures.

## Synopsis

Run the following command to create a new builder, named `kube`, that uses the
Kubernetes driver:

```console
$ docker buildx create \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  --driver-opt=[key=value,...]
```

The following table describes the available driver-specific options that you
can pass to `--driver-opt`:

| Parameter                    | Type         | Default                                 | Description                                                                                                                          |
| ---------------------------- | ------------ | --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `image`                      | String       |                                         | Sets the image to use for running BuildKit.                                                                                          |
| `namespace`                  | String       | Namespace in current Kubernetes context | Sets the Kubernetes namespace.                                                                                                       |
| `default-load`               | Boolean      | `false`                                 | Automatically load images to the Docker Engine image store.                                                                          |
| `replicas`                   | Integer      | 1                                       | Sets the number of Pod replicas to create. See [scaling BuildKit][1]                                                                 |
| `requests.cpu`               | CPU units    |                                         | Sets the request CPU value specified in units of Kubernetes CPU. For example `requests.cpu=100m` or `requests.cpu=2`                 |
| `requests.memory`            | Memory size  |                                         | Sets the request memory value specified in bytes or with a valid suffix. For example `requests.memory=500Mi` or `requests.memory=4G` |
| `requests.ephemeral-storage` | Storage size |                                         | Sets the request ephemeral-storage value specified in bytes or with a valid suffix. For example `requests.ephemeral-storage=2Gi`     |
| `limits.cpu`                 | CPU units    |                                         | Sets the limit CPU value specified in units of Kubernetes CPU. For example `requests.cpu=100m` or `requests.cpu=2`                   |
| `limits.memory`              | Memory size  |                                         | Sets the limit memory value specified in bytes or with a valid suffix. For example `requests.memory=500Mi` or `requests.memory=4G`   |
| `limits.ephemeral-storage`   | Storage size |                                         | Sets the limit ephemeral-storage value specified in bytes or with a valid suffix. For example `requests.ephemeral-storage=100M`      |
| `buildkit-root-volume-memory`| Memory size  | Using regular file system               | Mounts `/var/lib/buildkit` on an `emptyDir` memory-backed volume, with `SizeLimit` as the value. For example, `buildkit-root-folder-memory=6G`     |
| `nodeselector`               | CSV string   |                                         | Sets the pod's `nodeSelector` label(s). See [node assignment][2].                                                                    |
| `annotations`                | CSV string   |                                         | Sets additional annotations on the deployments and pods.                                                                             |
| `labels`                     | CSV string   |                                         | Sets additional labels on the deployments and pods.                                                                                  |
| `tolerations`                | CSV string   |                                         | Configures the pod's taint toleration. See [node assignment][2].                                                                     |
| `serviceaccount`             | String       |                                         | Sets the pod's `serviceAccountName`.                                                                                                 |
| `schedulername`              | String       |                                         | Sets the scheduler responsible for scheduling the pod.                                                                               |
| `timeout`                    | Time         | `120s`                                  | Set the timeout limit that determines how long Buildx will wait for pods to be provisioned before a build.                           |
| `rootless`                   | Boolean      | `false`                                 | Run the container as a non-root user. See [rootless mode][3].                                                                        |
| `loadbalance`                | String       | `sticky`                                | Load-balancing strategy (`sticky` or `random`). If set to `sticky`, the pod is chosen using the hash of the context path.            |
| `qemu.install`               | Boolean      | `false`                                 | Install QEMU emulation for multi platforms support. See [QEMU][4].                                                                   |
| `qemu.image`                 | String       | `tonistiigi/binfmt:latest`              | Sets the QEMU emulation image. See [QEMU][4].                                                                                        |

[1]: #scaling-buildkit
[2]: #node-assignment
[3]: #rootless-mode
[4]: #qemu

## Scaling BuildKit

One of the main advantages of the Kubernetes driver is that you can scale the
number of builder replicas up and down to handle increased build load. Scaling
is configurable using the following driver options:

- `replicas=N`

  This scales the number of BuildKit pods to the desired size. By default, it
  only creates a single pod. Increasing the number of replicas lets you take
  advantage of multiple nodes in your cluster.

- `requests.cpu`, `requests.memory`, `requests.ephemeral-storage`, `limits.cpu`, `limits.memory`, `limits.ephemeral-storage`

  These options allow requesting and limiting the resources available to each
  BuildKit pod according to the official Kubernetes documentation
  [here](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/).

For example, to create 4 replica BuildKit pods:

```console
$ docker buildx create \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  --driver-opt=namespace=buildkit,replicas=4
```

Listing the pods, you get this:

```console
$ kubectl -n buildkit get deployments
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
kube0   4/4     4            4           8s

$ kubectl -n buildkit get pods
NAME                     READY   STATUS    RESTARTS   AGE
kube0-6977cdcb75-48ld2   1/1     Running   0          8s
kube0-6977cdcb75-rkc6b   1/1     Running   0          8s
kube0-6977cdcb75-vb4ks   1/1     Running   0          8s
kube0-6977cdcb75-z4fzs   1/1     Running   0          8s
```

Additionally, you can use the `loadbalance=(sticky|random)` option to control
the load-balancing behavior when there are multiple replicas. `random` selects
random nodes from the node pool, providing an even workload distribution across
replicas. `sticky` (the default) attempts to connect the same build performed
multiple times to the same node each time, ensuring better use of local cache.

For more information on scalability, see the options for
[`docker buildx create`](/reference/cli/docker/buildx/create.md#driver-opt).

## Node assignment

The Kubernetes driver allows you to control the scheduling of BuildKit pods
using the `nodeSelector` and `tolerations` driver options.
You can also set the `schedulername` option if you want to use a custom scheduler altogether.

You can use the `annotations` and `labels` driver options to apply additional
metadata to the deployments and pods that's hosting your builders.

The value of the `nodeSelector` parameter is a comma-separated string of
key-value pairs, where the key is the node label and the value is the label
text. For example: `"nodeselector=kubernetes.io/arch=arm64"`

The `tolerations` parameter is a semicolon-separated list of taints. It accepts
the same values as the Kubernetes manifest. Each `tolerations` entry specifies
a taint key and the value, operator, or effect. For example:
`"tolerations=key=foo,value=bar;key=foo2,operator=exists;key=foo3,effect=NoSchedule"`

These options accept CSV-delimited strings as values. Due to quoting rules for
shell commands, you must wrap the values in single quotes. You can even wrap all
of `--driver-opt` in single quotes, for example:

```console
$ docker buildx create \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  '--driver-opt="nodeselector=label1=value1,label2=value2","tolerations=key=key1,value=value1"'
```

## Multi-platform builds

The Kubernetes driver has support for creating
[multi-platform images](/manuals/build/building/multi-platform.md),
either using QEMU or by leveraging the native architecture of nodes.

### QEMU

Like the `docker-container` driver, the Kubernetes driver also supports using
[QEMU](https://www.qemu.org/) (user
mode) to build images for non-native platforms. Include the `--platform` flag
and specify which platforms you want to output to.

For example, to build a Linux image for `amd64` and `arm64`:

```console
$ docker buildx build \
  --builder=kube \
  --platform=linux/amd64,linux/arm64 \
  -t <user>/<image> \
  --push .
```

> [!WARNING]
>
> QEMU performs full-CPU emulation of non-native platforms, which is much
> slower than native builds. Compute-heavy tasks like compilation and
> compression/decompression will likely take a large performance hit.

Using a custom BuildKit image or invoking non-native binaries in builds may
require that you explicitly turn on QEMU using the `qemu.install` option when
creating the builder:

```console
$ docker buildx create \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  --driver-opt=namespace=buildkit,qemu.install=true
```

### Native

If you have access to cluster nodes of different architectures, the Kubernetes
driver can take advantage of these for native builds. To do this, use the
`--append` flag of `docker buildx create`.

First, create your builder with explicit support for a single architecture, for
example `amd64`:

```console
$ docker buildx create \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  --platform=linux/amd64 \
  --node=builder-amd64 \
  --driver-opt=namespace=buildkit,nodeselector="kubernetes.io/arch=amd64"
```

This creates a Buildx builder named `kube`, containing a single builder node
named `builder-amd64`. Assigning a node name using `--node` is optional. Buildx
generates a random node name if you don't provide one.

Note that the Buildx concept of a node isn't the same as the Kubernetes concept
of a node. A Buildx node in this case could connect multiple Kubernetes nodes of
the same architecture together.

With the `kube` builder created, you can now introduce another architecture into
the mix using `--append`. For example, to add `arm64`:

```console
$ docker buildx create \
  --append \
  --bootstrap \
  --name=kube \
  --driver=kubernetes \
  --platform=linux/arm64 \
  --node=builder-arm64 \
  --driver-opt=namespace=buildkit,nodeselector="kubernetes.io/arch=arm64"
```

Listing your builders shows both nodes for the `kube` builder:

```console
$ docker buildx ls
NAME/NODE       DRIVER/ENDPOINT                                         STATUS   PLATFORMS
kube            kubernetes
  builder-amd64 kubernetes:///kube?deployment=builder-amd64&kubeconfig= running  linux/amd64*, linux/amd64/v2, linux/amd64/v3, linux/386
  builder-arm64 kubernetes:///kube?deployment=builder-arm64&kubeconfig= running  linux/arm64*
```

You can now build multi-arch `amd64` and `arm64` images, by specifying those
platforms together in your build command:

```console
$ docker buildx build --builder=kube --platform=linux/amd64,linux/arm64 -t <user>/<image> --push .
```

You can repeat the `buildx create --append` command for as many architectures
that you want to support.

## Rootless mode

The Kubernetes driver supports rootless mode. For more information on how
rootless mode works, and its requirements, see
[here](https://github.com/moby/buildkit/blob/master/docs/rootless.md).

To turn it on in your cluster, you can use the `rootless=true` driver option:

```console
$ docker buildx create \
  --name=kube \
  --driver=kubernetes \
  --driver-opt=namespace=buildkit,rootless=true
```

This will create your pods without `securityContext.privileged`.

Requires Kubernetes version 1.19 or later. Using Ubuntu as the host kernel is
recommended.

## Example: Creating a Buildx builder in Kubernetes

This guide shows you how to:

- Create a namespace for your Buildx resources
- Create a Kubernetes builder.
- List the available builders
- Build an image using your Kubernetes builders

Prerequisites:

- You have an existing Kubernetes cluster. If you don't already have one, you
  can follow along by installing
  [minikube](https://minikube.sigs.k8s.io/docs/).
- The cluster you want to connect to is accessible via the `kubectl` command,
  with the `KUBECONFIG` environment variable
  [set appropriately](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/#set-the-kubeconfig-environment-variable) if necessary.

1. Create a `buildkit` namespace.

   Creating a separate namespace helps keep your Buildx resources separate from
   other resources in the cluster.

   ```console
   $ kubectl create namespace buildkit
   namespace/buildkit created
   ```

2. Create a new builder with the Kubernetes driver:

   ```console
   $ docker buildx create \
     --bootstrap \
     --name=kube \
     --driver=kubernetes \
     --driver-opt=namespace=buildkit
   ```

   > [!NOTE]
   >
   > Remember to specify the namespace in driver options.

3. List available builders using `docker buildx ls`

   ```console
   $ docker buildx ls
   NAME/NODE                DRIVER/ENDPOINT STATUS  PLATFORMS
   kube                     kubernetes
     kube0-6977cdcb75-k9h9m                 running linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
   default *                docker
     default                default         running linux/amd64, linux/386
   ```

4. Inspect the running pods created by the build driver with `kubectl`.

   ```console
   $ kubectl -n buildkit get deployments
   NAME    READY   UP-TO-DATE   AVAILABLE   AGE
   kube0   1/1     1            1           32s

   $ kubectl -n buildkit get pods
   NAME                     READY   STATUS    RESTARTS   AGE
   kube0-6977cdcb75-k9h9m   1/1     Running   0          32s
   ```

   The build driver creates the necessary resources on your cluster in the
   specified namespace (in this case, `buildkit`), while keeping your driver
   configuration locally.

5. Use your new builder by including the `--builder` flag when running buildx
   commands. For example: :

   ```console
   # Replace <registry> with your Docker username
   # and <image> with the name of the image you want to build
   docker buildx build \
     --builder=kube \
     -t <registry>/<image> \
     --push .
   ```

That's it: you've now built an image from a Kubernetes pod, using Buildx.

## Further reading

For more information on the Kubernetes driver, see the
[buildx reference](/reference/cli/docker/buildx/create.md#driver).
