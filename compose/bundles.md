---
advisory: experimental
description: Description of Docker and Compose's experimental support for application bundles
keywords: documentation, docs, docker, compose, bundles, stacks
title: Docker stacks and distributed application bundles (experimental)
---

> **Note**: This is a modified copy of the [Docker Stacks and Distributed Application
> Bundles](https://github.com/moby/moby/blob/v1.12.0-rc4/experimental/docker-stacks-and-bundles.md)
> document in the [docker/docker-ce repo](https://github.com/docker/docker-ce). It's been updated to accurately reflect newer releases.

## Overview

A Dockerfile can be built into an image, and containers can be created from
that image. Similarly, a `docker-compose.yml` can be built into a **distributed
application bundle**, and **stacks** can be created from that bundle. In that
sense, the bundle is a multi-services distributable image format.

Docker Stacks and Distributed Application Bundles started as experimental
features introduced in Docker 1.12 and Docker Compose 1.8, alongside the concept
of swarm mode, and nodes and services in the Engine API. Neither Docker Engine
nor the Docker Registry support distribution of bundles, and the concept of a
`bundle` is not the emphasis for new releases going forward.

However, [swarm mode](/engine/swarm/index.md), multi-service applications, and
stack files now are fully supported. A stack file is a particular type of
[version 3 Compose file](/compose/compose-file/index.md).

If you are just getting started with Docker and want to learn the best way to
deploy multi-service applications, a good place to start is the [Get Started
walkthrough](/get-started/). This shows you how to define
a service configuration in a Compose file, deploy the app, and use
the relevant tools and commands.

## Produce a bundle

The easiest way to produce a bundle is to generate it using `docker-compose`
from an existing `docker-compose.yml`. Of course, that's just *one* possible way
to proceed, in the same way that `docker build` isn't the only way to produce a
Docker image.

From `docker-compose`:

```bash
$ docker-compose bundle
WARNING: Unsupported key 'network_mode' in services.nsqd - ignoring
WARNING: Unsupported key 'links' in services.nsqd - ignoring
WARNING: Unsupported key 'volumes' in services.nsqd - ignoring
[...]
Wrote bundle to vossibility-stack.dab
```

## Create a stack from a bundle

> **Note**: Because support for stacks and bundles is in the experimental stage,
> you need to install an experimental build of Docker Engine to use it.
>
> If you're on Mac or Windows, download the “Beta channel” version of
> [Docker for Mac](/docker-for-mac/) or
> [Docker for Windows](/docker-for-windows/) to install
> it. If you're on Linux, follow the instructions in the
> [experimental build README](https://github.com/docker/cli/blob/master/experimental/README.md).

A stack is created using the `docker deploy` command:

```bash
# docker deploy --help

Usage:  docker deploy [OPTIONS] STACK

Create and update a stack

Options:
      --file   string        Path to a Distributed Application Bundle file (Default: STACK.dab)
      --help                 Print usage
      --with-registry-auth   Send registry authentication details to Swarm agents
```

Let's deploy the stack created before:

```bash
# docker deploy vossibility-stack
Loading bundle from vossibility-stack.dab
Creating service vossibility-stack_elasticsearch
Creating service vossibility-stack_kibana
Creating service vossibility-stack_logstash
Creating service vossibility-stack_lookupd
Creating service vossibility-stack_nsqd
Creating service vossibility-stack_vossibility-collector
```

We can verify that services were correctly created:

```bash
# docker service ls
ID            NAME                                     REPLICAS  IMAGE
COMMAND
29bv0vnlm903  vossibility-stack_lookupd                1 nsqio/nsq@sha256:eeba05599f31eba418e96e71e0984c3dc96963ceb66924dd37a47bf7ce18a662 /nsqlookupd
4awt47624qwh  vossibility-stack_nsqd                   1 nsqio/nsq@sha256:eeba05599f31eba418e96e71e0984c3dc96963ceb66924dd37a47bf7ce18a662 /nsqd --data-path=/data --lookupd-tcp-address=lookupd:4160
4tjx9biia6fs  vossibility-stack_elasticsearch          1 elasticsearch@sha256:12ac7c6af55d001f71800b83ba91a04f716e58d82e748fa6e5a7359eed2301aa
7563uuzr9eys  vossibility-stack_kibana                 1 kibana@sha256:6995a2d25709a62694a937b8a529ff36da92ebee74bafd7bf00e6caf6db2eb03
9gc5m4met4he  vossibility-stack_logstash               1 logstash@sha256:2dc8bddd1bb4a5a34e8ebaf73749f6413c101b2edef6617f2f7713926d2141fe logstash -f /etc/logstash/conf.d/logstash.conf
axqh55ipl40h  vossibility-stack_vossibility-collector  1 icecrime/vossibility-collector@sha256:f03f2977203ba6253988c18d04061c5ec7aab46bca9dfd89a9a1fa4500989fba --config /config/config.toml --debug
```

## Manage stacks

Stacks are managed using the `docker stack` command:

```bash
# docker stack --help

Usage:  docker stack COMMAND

Manage Docker stacks

Options:
      --help   Print usage

Commands:
  config      Print the stack configuration
  deploy      Create and update a stack
  rm          Remove the stack
  services    List the services in the stack
  tasks       List the tasks in the stack

Run 'docker stack COMMAND --help' for more information on a command.
```

## Bundle file format

Distributed application bundles are described in a JSON format. When bundles
are persisted as files, the file extension is `.dab`.

A bundle has two top-level fields: `version` and `services`. The version used
by Docker 1.12 tools is `0.1`.

`services` in the bundle are the services that comprise the app. They
correspond to the new `Service` object introduced in the 1.12 Docker Engine API.

A service has the following fields:

<dl>
    <dt>
        Image (required) <code>string</code>
    </dt>
    <dd>
        The image that the service runs. Docker images should be referenced
        with full content hash to fully specify the deployment artifact for the
        service. Example:
        <code>postgres@sha256:e0a230a9f5b4e1b8b03bb3e8cf7322b0e42b7838c5c87f4545edb48f5eb8f077</code>
    </dd>
    <dt>
        Command <code>[]string</code>
    </dt>
    <dd>
        Command to run in service containers.
    </dd>
    <dt>
        Args <code>[]string</code>
    </dt>
    <dd>
        Arguments passed to the service containers.
    </dd>
    <dt>
        Env <code>[]string</code>
    </dt>
    <dd>
        Environment variables.
    </dd>
    <dt>
        Labels <code>map[string]string</code>
    </dt>
    <dd>
        Labels used for setting meta data on services.
    </dd>
    <dt>
        Ports <code>[]Port</code>
    </dt>
    <dd>
        Service ports (composed of <code>Port</code> (<code>int</code>) and
        <code>Protocol</code> (<code>string</code>). A service description can
        only specify the container port to be exposed. These ports can be
        mapped on runtime hosts at the operator's discretion.
    </dd>

    <dt>
        WorkingDir <code>string</code>
    </dt>
    <dd>
        Working directory inside the service containers.
    </dd>

    <dt>
        User <code>string</code>
    </dt>
    <dd>
        Username or UID (format: <code>&lt;name|uid&gt;[:&lt;group|gid&gt;]</code>).
    </dd>

    <dt>
        Networks <code>[]string</code>
    </dt>
    <dd>
        Networks that the service containers should be connected to. An entity
        deploying a bundle should create networks as needed.
    </dd>
</dl>

> **Note**: Some configuration options are not yet supported in the DAB format,
> including volume mounts.

## Related topics

* [Get started walkthrough](/get-started/)

* [docker stack deploy](/engine/reference/commandline/stack_deploy/) command

* [deploy](/compose/compose-file/index.md#deploy) option in [Compose files](/compose/compose-file/index.md)
