---
description: Deploying to Kubernetes on Docker for Mac
keywords: mac, edge, kubernetes, kubectl, orchestration
title: Deploy to Kubernetes
---

**Kubernetes is only available in Docker for Mac 17.12 CE Edge to participants in the [Docker Beta program](https://beta.docker.com/). To access beta builds, you must be signed in with your Docker ID within Docker for Mac: select ![whale menu](/docker-for-mac/images/whale-x.png){: .inline} -> Sign in / Create Docker ID from the menu bar.**

If you are part of the Docker Beta program, Docker for Mac 17.12 CE Edge
includes a standalone Kubernetes server and client, as well as Docker CLI
integration. The Kubernetes server runs locally within your Docker instance, is
not configurable, and is a single-node cluster.

The Kubernetes server runs within a Docker container on your Mac, and is only
for local testing. When Kubernetes support is enabled, you can deploy your
workloads, in parallel, on Kubernetes, Swarm, and as standalone containers.
Enabling or disabling the Kubernetes server does not affect your other
workloads.

See [Docker for Mac > Getting started](/docker-for-mac/index.md#kubernetes) to
enable Kubernetes and begin testing the deployment of your workloads on
Kubernetes.

> If you independently installed the Kubernetes CLI, `kubectl`, make sure that
> it is pointing to `docker-for-desktop` and not some other context such as
> `minikube` or a GKE cluster. Run: `kubectl config use-context docker-for-desktop`.
> If you experience conflicts with an existing `kubectl` installation, remove `/usr/local/bin/kubectl`.

## Use Docker commands

You can deploy a stack on Kubernetes with `docker stack deploy`, the
`docker-compose.yml` file, and the name of the stack.

```bash
$ docker stack deploy --compose-file /path/to/docker-compose.yml mystack
$ docker stack services mystack
```

You can see the service deployed with the `kubectl get services` command.

### Specify a namespace

By default, the `default` namespace is used. You can specify a namespace with
the `--namespace` flag.

```bash
$ docker stack deploy --namespace my-app --compose-file /path/to/docker-compose.yml mystack
```

Run `kubectl get services -n my-app` to see only the services deployed in the
`my-app` namespace.

### Override the default orchestrator

While testing Kubernetes, you may want to deploy some workloads in swarm mode.
Use the `DOCKER_ORCHESTRATOR` variable to override the default orchestrator for
a given terminal session or a single Docker command. This variable can be unset
(the default, in which case Kubernetes is the orchestrator) or set to `swarm` or
`kubernetes`. The following command overrides the orchestrator for a single
deployment, by setting the variable at the start of the command itself.

```bash
DOCKER_ORCHESTRATOR=swarm docker stack deploy --compose-file /path/to/docker-compose.yml mystack
```

> **Note**: Deploying the same app in Kubernetes and swarm mode may lead to
> conflicts with ports and service names.

## Use the kubectl command

The Docker for Mac Kubernetes integration provides the Kubernetes CLI command
at `/usr/local/bin/kubectl`. This location may not be in your shell's `PATH`
variable, so you may need to type the full path of the command or add it to
the `PATH`. For more information about `kubectl`, see the
[official `kubectl` documentation](https://kubernetes.io/docs/reference/kubectl/overview/).
You can test the command by listing the available nodes:

```bash
$ kubectl get nodes

NAME                 STATUS    ROLES     AGE       VERSION
docker-for-desktop   Ready     master    3h        v1.8.2
```

## Example app

Docker has created the following demo app that you can deploy to swarm mode or
to Kubernetes using the `docker stack deploy` command.

```yaml
version: '3.3'

services:
  web:
    build: web
    image: dockerdemos/lab-web
    volumes:
     - "./web/static:/static"
    ports:
     - "80:80"

  words:
    build: words
    image: dockerdemos/lab-words
    deploy:
      replicas: 5
      endpoint_mode: dnsrr
      resources:
        limits:
          memory: 16M
        reservations:
          memory: 16M

  db:
    build: db
    image: dockerdemos/lab-db
```

If you already have a Kubernetes YAML file, you can deploy it using the
`kubectl` command.
