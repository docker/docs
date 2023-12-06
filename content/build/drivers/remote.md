---
title: Remote driver
description: |
  The remote driver lets you connect to a remote BuildKit instance
  that you set up and configure manually.
keywords: build, buildx, driver, builder, remote
aliases:
  - /build/buildx/drivers/remote/
  - /build/building/drivers/remote/
---

The Buildx remote driver allows for more complex custom build workloads,
allowing you to connect to externally managed BuildKit instances. This is useful
for scenarios that require manual management of the BuildKit daemon, or where a
BuildKit daemon is exposed from another source.

## Synopsis

```console
$ docker buildx create \
  --name remote \
  --driver remote \
  tcp://localhost:1234
```

The following table describes the available driver-specific options that you can
pass to `--driver-opt`:

| Parameter    | Type   | Default            | Description                                                            |
| ------------ | ------ | ------------------ | ---------------------------------------------------------------------- |
| `key`        | String |                    | Sets the TLS client key.                                               |
| `cert`       | String |                    | Absolute path to the TLS client certificate to present to `buildkitd`. |
| `cacert`     | String |                    | Absolute path to the TLS certificate authority used for validation.    |
| `servername` | String | Endpoint hostname. | TLS server name used in requests.                                      |

## Example: Remote BuildKit over Unix sockets

This guide shows you how to create a setup with a BuildKit daemon listening on a
Unix socket, and have Buildx connect through it.

1. Ensure that [BuildKit](https://github.com/moby/buildkit) is installed.

   For example, you can launch an instance of buildkitd with:

   ```console
   $ sudo ./buildkitd --group $(id -gn) --addr unix://$HOME/buildkitd.sock
   ```

   Alternatively, [see here](https://github.com/moby/buildkit/blob/master/docs/rootless.md)
   for running buildkitd in rootless mode or [here](https://github.com/moby/buildkit/tree/master/examples/systemd)
   for examples of running it as a systemd service.

2. Check that you have a Unix socket that you can connect to.

   ```console
   $ ls -lh /home/user/buildkitd.sock
   srw-rw---- 1 root user 0 May  5 11:04 /home/user/buildkitd.sock
   ```

3. Connect Buildx to it using the remote driver:

   ```console
   $ docker buildx create \
     --name remote-unix \
     --driver remote \
     unix://$HOME/buildkitd.sock
   ```

4. List available builders with `docker buildx ls`. You should then see
   `remote-unix` among them:

   ```console
   $ docker buildx ls
   NAME/NODE           DRIVER/ENDPOINT                        STATUS  PLATFORMS
   remote-unix         remote
     remote-unix0      unix:///home/.../buildkitd.sock        running linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
   default *           docker
     default           default                                running linux/amd64, linux/386
   ```

You can switch to this new builder as the default using
`docker buildx use remote-unix`, or specify it per build using `--builder`:

```console
$ docker buildx build --builder=remote-unix -t test --load .
```

Remember that you need to use the `--load` flag if you want to load the build
result into the Docker daemon.

## Example: Remote BuildKit in Docker container

This guide will show you how to create setup similar to the `docker-container`
driver, by manually booting a BuildKit Docker container and connecting to it
using the Buildx remote driver. This procedure will manually create a container
and access it via it's exposed port. (You'd probably be better of just using the
`docker-container` driver that connects to BuildKit through the Docker daemon,
but this is for illustration purposes.)

1.  Generate certificates for BuildKit.

    You can use this [bake definition](https://github.com/moby/buildkit/blob/master/examples/create-certs)
    as a starting point:

    ```console
    SAN="localhost 127.0.0.1" docker buildx bake "https://github.com/moby/buildkit.git#master:examples/create-certs"
    ```

    Note that while it's possible to expose BuildKit over TCP without using
    TLS, it's not recommended. Doing so allows arbitrary access to BuildKit
    without credentials.

2.  With certificates generated in `.certs/`, startup the container:

    ```console
    $ docker run -d --rm \
      --name=remote-buildkitd \
      --privileged \
      -p 1234:1234 \
      -v $PWD/.certs:/etc/buildkit/certs \
      moby/buildkit:latest \
      --addr tcp://0.0.0.0:1234 \
      --tlscacert /etc/buildkit/certs/daemon/ca.pem \
      --tlscert /etc/buildkit/certs/daemon/cert.pem \
      --tlskey /etc/buildkit/certs/daemon/key.pem
    ```

    This command starts a BuildKit container and exposes the daemon's port 1234
    to localhost.

3.  Connect to this running container using Buildx:

    ```console
    $ docker buildx create \
      --name remote-container \
      --driver remote \
      --driver-opt cacert=${PWD}/.certs/client/ca.pem,cert=${PWD}/.certs/client/cert.pem,key=${PWD}/.certs/client/key.pem,servername=<TLS_SERVER_NAME> \
      tcp://localhost:1234
    ```

    Alternatively, use the `docker-container://` URL scheme to connect to the
    BuildKit container without specifying a port:

    ```console
    $ docker buildx create \
      --name remote-container \
      --driver remote \
      docker-container://remote-container
    ```

## Example: Remote BuildKit in Kubernetes

This guide will show you how to create a setup similar to the `kubernetes`
driver by manually creating a BuildKit `Deployment`. While the `kubernetes`
driver will do this under-the-hood, it might sometimes be desirable to scale
BuildKit manually. Additionally, when executing builds from inside Kubernetes
pods, the Buildx builder will need to be recreated from within each pod or
copied between them.

1. Create a Kubernetes deployment of `buildkitd`, as per the instructions
   [here](https://github.com/moby/buildkit/tree/master/examples/kubernetes).

   Following the guide, create certificates for the BuildKit daemon and client
   using [create-certs.sh](https://github.com/moby/buildkit/blob/master/examples/kubernetes/create-certs.sh),
   and create a deployment of BuildKit pods with a service that connects to
   them.

2. Assuming that the service is called `buildkitd`, create a remote builder in
   Buildx, ensuring that the listed certificate files are present:

   ```console
   $ docker buildx create \
     --name remote-kubernetes \
     --driver remote \
     --driver-opt cacert=${PWD}/.certs/client/ca.pem,cert=${PWD}/.certs/client/cert.pem,key=${PWD}/.certs/client/key.pem \
     tcp://buildkitd.default.svc:1234
   ```

Note that this only works internally, within the cluster, since the BuildKit
setup guide only creates a `ClusterIP` service. To access a builder remotely,
you can set up and use an ingress, which is outside the scope of this guide.

### Debug a remote builder in Kubernetes

If you're having trouble accessing a remote builder deployed in Kubernetes, you
can use the `kube-pod://` URL scheme to connect directly to a BuildKit pod
through the Kubernetes API. Note that this method only connects to a single pod
in the deployment.

```console
$ kubectl get pods --selector=app=buildkitd -o json | jq -r '.items[].metadata.name'
buildkitd-XXXXXXXXXX-xxxxx
$ docker buildx create \
  --name remote-container \
  --driver remote \
  kube-pod://buildkitd-XXXXXXXXXX-xxxxx
```

Alternatively, use the port forwarding mechanism of `kubectl`:

```console
$ kubectl port-forward svc/buildkitd 1234:1234
```

Then you can point the remote driver at `tcp://localhost:1234`.
