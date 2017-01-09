---
datafolder: engine-cli
datafile: docker_container_port
title: docker container port
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

```bash
$ docker ps

CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                                            NAMES
b650456536c7        busybox:latest      top                 54 minutes ago      Up 54 minutes       0.0.0.0:1234->9876/tcp, 0.0.0.0:4321->7890/tcp   test
```

### Find out all the ports mapped

```bash
$ docker container port test

7890/tcp -> 0.0.0.0:4321
9876/tcp -> 0.0.0.0:1234
```

### Find out a specific mapping

```bash
$ docker container port test 7890/tcp

0.0.0.0:4321
```

```bash
$ docker container port test 7890

0.0.0.0:4321
```

### An example showing error for non-existent mapping

```bash
$ docker container port test 7890/udp

2014/06/24 11:53:36 Error: No public port '7890/udp' published for test
```
