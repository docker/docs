---
title: Docker Context
description: Learn about Docker Context
keywords: engine, context, cli, kubernetes
---


## Introduction

This guide shows how _contexts_ make it easy for a **single Docker CLI** to manage multiple Swarm clusters, multiple Kubernetes clusters, and multiple individual Docker nodes.

A single Docker CLI can have multiple contexts. Each context contains all of the endpoint and security information required to manage a different cluster or node. The `docker context` command makes it easy to configure these contexts and switch between them.

As an example, a single Docker client on your company laptop might be configured with two contexts; **dev-k8s** and **prod-swarm**. **dev-k8s** contains the endpoint data and security credentials to configure and manage a Kubernetes cluster in a development environment. **prod-swarm** contains everything required to manage a Swarm cluster in a production environment. Once these contexts are configured, you can use the top-level `docker context use <context-name>` to easily switch between them.

For information on using Docker Context to deploy your apps to the cloud, see [Deploying Docker containers on Azure](aci-integration.md) and [Deploying Docker containers on ECS](ecs-integration.md).

## Prerequisites

To follow the examples in this guide, you'll need:

- A Docker client that supports the top-level `context` command

Run `docker context` to verify that your Docker client supports contexts.

You will also need one of the following:

- Docker Swarm cluster
- Single-engine Docker node
- Kubernetes cluster

## The anatomy of a context

A context is a combination of several properties. These include:

- Name
- Endpoint configuration
- TLS info
- Orchestrator

The easiest way to see what a context looks like is to view the **default** context.

```
$ docker context ls
NAME          DESCRIPTION     DOCKER ENDPOINT                KUBERNETES ENDPOINT      ORCHESTRATOR
default *     Current...      unix:///var/run/docker.sock                             swarm
```

This shows a single context called "default". It's configured to talk to a Swarm cluster through the local `/var/run/docker.sock` Unix socket. It has no Kubernetes endpoint configured.

The asterisk in the `NAME` column indicates that this is the active context. This means all `docker` commands will be executed against the "default" context unless overridden with environment variables such as `DOCKER_HOST` and `DOCKER_CONTEXT`, or on the command-line with the `--context` and `--host` flags.

Dig a bit deeper with `docker context inspect`. In this example, we're inspecting the context called `default`.

```
$ docker context inspect default
[
    {
        "Name": "default",
        "Metadata": {
            "StackOrchestrator": "swarm"
        },
        "Endpoints": {
            "docker": {
                "Host": "unix:///var/run/docker.sock",
                "SkipTLSVerify": false
            }
        },
        "TLSMaterial": {},
        "Storage": {
            "MetadataPath": "\u003cIN MEMORY\u003e",
            "TLSPath": "\u003cIN MEMORY\u003e"
        }
    }
]
```

This context is using "swarm" as the orchestrator (`metadata.stackOrchestrator`). It is configured to talk to an endpoint exposed on a local Unix socket at `/var/run/docker.sock` (`Endpoints.docker.Host`), and requires TLS verification (`Endpoints.docker.SkipTLSVerify`).

### Create a new context

You can create new contexts with the `docker context create` command.

The following example creates a new context called "docker-test" and specifies the following:

- Default orchestrator = Swarm
- Issue commands to the local Unix socket `/var/run/docker.sock`

```
$ docker context create docker-test \
  --default-stack-orchestrator=swarm \
  --docker host=unix:///var/run/docker.sock

Successfully created context "docker-test"
```

The new context is stored in a `meta.json` file below `~/.docker/contexts/`. Each new context you create gets its own `meta.json` stored in a dedicated sub-directory of `~/.docker/contexts/`.

> **Note:** The default context behaves differently than manually created contexts. It does not have a `meta.json` configuration file, and it dynamically updates based on the current configuration. For example, if you switch your current Kubernetes config using `kubectl config use-context`, the default Docker context will dynamically update itself to the new Kubernetes endpoint.

You can view the new context with `docker context ls` and `docker context inspect <context-name>`.

The following can be used to create a config with Kubernetes as the default orchestrator using the existing kubeconfig stored in `/home/ubuntu/.kube/config`. For this to work, you will need a valid kubeconfig file in `/home/ubuntu/.kube/config`. If your kubeconfig has more than one context, the current context (`kubectl config current-context`) will be used.

```
$ docker context create k8s-test \
  --default-stack-orchestrator=kubernetes \
  --kubernetes config-file=/home/ubuntu/.kube/config \
  --docker host=unix:///var/run/docker.sock

Successfully created context "k8s-test"
```

You can view all contexts on the system with `docker context ls`.

```
$ docker context ls
NAME           DESCRIPTION   DOCKER ENDPOINT               KUBERNETES ENDPOINT               ORCHESTRATOR
default *      Current       unix:///var/run/docker.sock   https://35.226.99.100 (default)   swarm
k8s-test                     unix:///var/run/docker.sock   https://35.226.99.100 (default)   kubernetes
docker-test                  unix:///var/run/docker.sock                                     swarm
```

The current context is indicated with an asterisk ("\*").

## Use a different context

You can use `docker context use` to quickly switch between contexts.

The following command will switch the `docker` CLI to use the "k8s-test" context.

```
$ docker context use k8s-test

k8s-test
Current context is now "k8s-test"
```

Verify the operation by listing all contexts and ensuring the asterisk ("\*") is against the "k8s-test" context.

```
$ docker context ls
NAME            DESCRIPTION                               DOCKER ENDPOINT               KUBERNETES ENDPOINT               ORCHESTRATOR
default         Current DOCKER_HOST based configuration   unix:///var/run/docker.sock   https://35.226.99.100 (default)   swarm
docker-test                                               unix:///var/run/docker.sock                                     swarm
k8s-test *                                                unix:///var/run/docker.sock   https://35.226.99.100 (default)   kubernetes
```

`docker` commands will now target endpoints defined in the "k8s-test" context.

You can also set the current context using the `DOCKER_CONTEXT` environment variable. This overrides the context set with `docker context use`.

Use the appropriate command below to set the context to `docker-test` using an environment variable.

Windows PowerShell:

```
> $Env:DOCKER_CONTEXT=docker-test
```

Linux:

```
$ export DOCKER_CONTEXT=docker-test
```

Run a `docker context ls` to verify that the "docker-test" context is now the active context.

You can also use the global `--context` flag to override the context specified by the `DOCKER_CONTEXT` environment variable. For example, the following will send the command to a context called "production".

```
$ docker --context production container ls
```

## Exporting and importing Docker contexts

The `docker context` command makes it easy to export and import contexts on different machines with the Docker client installed.

You can use the `docker context export` command to export an existing context to a file. This file can later be imported on another machine that has the `docker` client installed.

By default, contexts will be exported as a _native Docker contexts_. You can export and import these using the `docker context` command. If the context you are exporting includes a Kubernetes endpoint, the Kubernetes part of the context will be included in the `export` and `import` operations.

There is also an option to export just the Kubernetes part of a context. This will produce a native kubeconfig file that can be manually merged with an existing `~/.kube/config` file on another host that has `kubectl` installed. You cannot export just the Kubernetes portion of a context and then import it with `docker context import`. The only way to import the exported Kubernetes config is to manually merge it into an existing kubeconfig file.

Let's look at exporting and importing a native Docker context.

### Exporting and importing a native Docker context

The following example exports an existing context called "docker-test". It will be written to a file called `docker-test.dockercontext`.

```
$ docker context export docker-test
Written file "docker-test.dockercontext"
```

Check the contents of the export file.

```
$ cat docker-test.dockercontext
meta.json0000644000000000000000000000022300000000000011023 0ustar0000000000000000{"Name":"docker-test","Metadata":{"StackOrchestrator":"swarm"},"Endpoints":{"docker":{"Host":"unix:///var/run/docker.sock","SkipTLSVerify":false}}}tls0000700000000000000000000000000000000000000007716 5ustar0000000000000000
```

This file can be imported on another host using `docker context import`. The target host must have the Docker client installed.

```
$ docker context import docker-test docker-test.dockercontext
docker-test
Successfully imported context "docker-test"
```

You can verify that the context was imported with `docker context ls`.

The format of the import command is `docker context import <context-name> <context-file>`.

Now, let's look at exporting just the Kubernetes parts of a context.

### Exporting a Kubernetes context

You can export a Kubernetes context only if the context you are exporting has a Kubernetes endpoint configured. You cannot import a Kubernetes context using `docker context import`.

These steps will use the `--kubeconfig` flag to export **only** the Kubernetes elements of the existing `k8s-test` context to a file called "k8s-test.kubeconfig". The `cat` command will then show that it's exported as a valid kubeconfig file.

```
$ docker context export k8s-test --kubeconfig
Written file "k8s-test.kubeconfig"
```

Verify that the exported file contains a valid kubectl config.

```
$ cat k8s-test.kubeconfig
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data:
    <Snip>
    server: https://35.226.99.100
  name: cluster
contexts:
- context:
    cluster: cluster
    namespace: default
    user: authInfo
  name: context
current-context: context
kind: Config
preferences: {}
users:
- name: authInfo
  user:
    auth-provider:
      config:
        cmd-args: config config-helper --format=json
        cmd-path: /snap/google-cloud-sdk/77/bin/gcloud
        expiry-key: '{.credential.token_expiry}'
        token-key: '{.credential.access_token}'
      name: gcp
```

You can merge this with an existing `~/.kube/config` file on another machine.

## Updating a context

You can use `docker context update` to update fields in an existing context.

The following example updates the "Description" field in the existing `k8s-test` context.

```
$ docker context update k8s-test --description "Test Kubernetes cluster"
k8s-test
Successfully updated context "k8s-test"
```
