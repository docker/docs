---
description: Deploying to Kubernetes on Docker for Mac
keywords: mac, edge, kubernetes, kubectl, orchestration
title: Deploy to Kubernetes
---

If you are participating in the [Docker Beta program](https://beta.docker.com/),
you can access the beta for Docker for Mac 17.12 CE, which includes a standalone
Kubernetes server and client, as well as Docker CLI integration.
[You can enable this feature](/docker-for-mac/index.md#kubernetes) so that you
can test deploying your workloads on Kubernetes. The Kubernetes server runs
within a Docker container on your Mac, and is only for local testing. When
Kubernetes support is enabled, you can deploy your workloads to Kubernetes,
swarm mode, and standalone containers in parallel. Enabling or disabling the
Kubernetes server does not affect your other workloads.

The Kubernetes server runs locally within your Docker instance, is not
configurable, and is a single-node cluster. It is provided for development and
testing only.

## Use Docker commands

You can deploy a stack to Kubernetes using the `docker stack deploy` command,
supplying a `docker-compose.yml` file.

```bash
$ docker stack deploy --compose-file /path/to/docker-compose.yml
```

You can see the service deployed using the `kubectl get services` command.

### Specify a namespace

By default, the `default` namespace is used. You can specify a namespace using
the `--namespace` flag.

```bash
$ docker stack deploy --namespace my-app --compose-file /path/to/docker-compose.yml
```

You can then use the `kubectl get services -n my-app` to see only the services
deployed in the `my-app` namespace.

### Override the default orchestrator

While you are testing Kubernetes, you may want to deploy some workloads in swarm
mode. You can use the `DOCKER_ORCHESTRATOR` variable to override the default
orchestrator for a given terminal session or a single Docker command. The
variable can be unset (the default, in which case Kubernetes is the
orchestrator) or set to `swarm` or `kubernetes`. The following command
overrides the orchestrator for a single deployment, by setting the variable at
the start of the command itself.

```bash
DOCKER_ORCHESTRATOR=swarm docker stack deploy --compose-file /path/to/docker-compose.yml
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

Docker has created the following demo app that you can deploy to swarm mode or to to
Kubernetes using the `docker stack deploy` command.

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

TODO create a Kube YAML file to make the same app as the docker-compse.yml
above?