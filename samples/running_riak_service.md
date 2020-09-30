---
description: Build a Docker image with Riak pre-installed
keywords: docker, example, package installation, networking,  riak
title: Dockerize a Riak service
redirect_from:
  - /engine/examples/running_riak_service/
---

The goal of this example is to show you how to build a Docker image with
Riak pre-installed.

## Create a Dockerfile

Create an empty file called `Dockerfile`:

```bash
$ touch Dockerfile
```

Next, define the parent image you want to use to build your image on top
of. We use [Ubuntu](https://hub.docker.com/_/ubuntu/) (tag:
`trusty`), which is available on [Docker Hub](https://hub.docker.com):

```dockerfile
# Riak
#
# VERSION       0.1.1

# Use the Ubuntu parent image provided by dotCloud
FROM ubuntu:trusty
```

After that, we install the curl which is used to download the repository setup
script and we download the setup script and run it.

```dockerfile
# Install Riak repository before we do apt-get update, so that update happens
# in a single step
RUN apt-get install -q -y curl && \
    curl -fsSL https://packagecloud.io/install/repositories/basho/riak/script.deb | sudo bash
```

Then we install and setup a few dependencies:

 - `supervisor` is used manage the Riak processes
 - `riak=2.0.5-1` is the Riak package coded to version 2.0.5

```dockerfile
# Install and setup project dependencies
RUN apt-get update && \
    apt-get install -y supervisor riak=2.0.5-1

RUN mkdir -p /var/log/supervisor

RUN locale-gen en_US en_US.UTF-8

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
```

After that, we modify Riak's configuration:

```dockerfile
# Configure Riak to accept connections from any host
RUN sed -i "s|listener.http.internal = 127.0.0.1:8098|listener.http.internal = 0.0.0.0:8098|" /etc/riak/riak.conf
RUN sed -i "s|listener.protobuf.internal = 127.0.0.1:8087|listener.protobuf.internal = 0.0.0.0:8087|" /etc/riak/riak.conf
```

Then, we expose the Riak Protocol Buffers and HTTP interfaces:

```dockerfile
# Expose Riak Protocol Buffers and HTTP interfaces
EXPOSE 8087 8098
```

Finally, run `supervisord` so that Riak is started:

```dockerfile
CMD ["/usr/bin/supervisord"]
```

## Create a supervisord configuration file

Create an empty file called `supervisord.conf`. Make
sure it's at the same directory level as your `Dockerfile`:

```bash
touch supervisord.conf
```

Populate it with the following program definitions:

```ini
[supervisord]
nodaemon=true

[program:riak]
command=bash -c "/usr/sbin/riak console"
numprocs=1
autostart=true
autorestart=true
user=riak
environment=HOME="/var/lib/riak"
stdout_logfile=/var/log/supervisor/%(program_name)s.log
stderr_logfile=/var/log/supervisor/%(program_name)s.log
```

## Build the Docker image for Riak

Now you can build a Docker image for Riak:

```bash
$ docker build -t "<yourname>/riak" .
```

## Next steps

Riak is a distributed database. Many production deployments consist of
[at least five nodes](
https://riak.com/why-your-riak-cluster-should-have-at-least-five-nodes/).
See the [docker-riak](https://github.com/hectcastro/docker-riak) project
details on how to deploy a Riak cluster using Docker and Pipework.
