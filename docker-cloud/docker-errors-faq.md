---
description: Known Docker Engine issues in Docker Cloud
keywords: Engine, issues, troubleshoot
redirect_from:
- /docker-cloud/faq/docker-errors-faq/
title: Known issues in Docker Cloud
toc_max: 2
toc_min: 1
---

This is a list of known issues with current versions of Docker Engine along with
our recommended workaround. You might encounter these errors in Docker Cloud.

## Errors and messages

---

## Get i/o timeout

<!-- span tag prevents irritating autolinker from interpreting this as a link -->

*Get https<span></span>://index.docker.io/v1/repositories/\<image\>/images: dial tcp: lookup \<registry host> on \<ip>:53: read udp \<ip>:53: i/o timeout*

### Description

The DNS resolver configured on the host cannot resolve the registry's hostname.

### GitHub link

N/A

### Workaround

Retry the operation, or if the error persists, use another DNS resolver. You can do this by updating your `/etc/resolv.conf` file with these or other DNS servers:

	nameserver 8.8.8.8
	nameserver 8.8.4.4

---

## 500 Server Error: userland proxy

*500 Server Error: Internal Server Error ("Cannot start container \<id>: Error starting userland proxy: listen tcp 0.0.0.0:\<port>: bind: address already in use")*

### Description

Docker Cloud is trying to deploy a container publishing a port which is already
being used by a process on the host (like the SSH server listening in port
`22`).

### GitHub link

N/A

### Workaround

Either choose another port, or SSH into the node and manually stop the process
which is using the port that you are trying to use.

---

## 500 Server Error: bind failed

*500 Server Error: Internal Server Error ("Cannot start container \<id>: Bind for 0.0.0.0:\<port> failed: port is already allocated")*

### Description

Docker Cloud is trying to deploy a container publishing a port which is already
used by another container outside of the scope of Docker Cloud.

### GitHub link

N/A

### Workaround

Either choose another port, or SSH into the node and manually stop the container
which is using the port that you are trying to use.

---

## 500 Server Error: cannot start, executable not found

*500 Server Error: Internal Server Error ("Cannot start container \<id>: [8] System error: exec: "\<path>": executable file not found in $PATH")*

### Description

The `run` command you specified for the container does not exist on the
container.

### GitHub link

N/A

### Workaround

Edit the service to fix the run command.

---

## Timeout when pulling image from the registry

*Timeout when pulling image from the registry*

### Description

Timeouts occur when pulling the image takes more than 10 minutes. This can
sometimes be caused by the Docker daemon waiting for a nonexistent process while
pulling the required image.

### GitHub link


[docker/docker#12823](https://github.com/moby/moby/issues/12823){: target="_blank" class="_" }

### Workaround

Restart the `dockercloud-agent` service (`sudo service dockercloud-agent
restart`) on the node, or restart the node.

---

## Docker Cloud CLI does not currently support Python 3

### Description

The `docker-cloud` command line interface (CLI) does not currently support
Python 3.x.


### GitHub link

[docker/docker-cloud#21](https://github.com/docker/dockercloud-cli/issues/21){: target="_blank" class="_"}

### Workarounds

* Use Python 2.x with the Docker Cloud CLI.

---

## Problems installing and running Docker Cloud with Python 3

### Description

* Some users have encountered problems installing and/or running
Docker Cloud with Anaconda Python 3.5.2 on a Windows host.

* Some users running Python on Windows have encountered problems
running `docker-cloud` inside a container using `docker run`.

### GitHub link

[docker/for-win#368](https://github.com/docker/for-win/issues/368){: target="_blank" class="_" }

[docker/dockercloud-cli#45](https://github.com/docker/dockercloud-cli/issues/45){: target="_blank" class"_" }

### Workarounds

* If you encounter problems with the installation, use Python 2.x.

* Before attempting to run `docker-cloud` in a container with `docker run`,
make sure that you [have Linux containers
enabled](/docker-for-windows/index.md#switch-between-windows-and-linux-containers){: target="_blank" class"_" }.
