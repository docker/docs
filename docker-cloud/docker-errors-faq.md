---
description: Known Docker Engine issues in Docker Cloud
keywords: Engine, issues, troubleshoot
redirect_from:
- /docker-cloud/faq/docker-errors-faq/
title: Known issues in Docker Cloud
---

This is a list of known issues with current versions of Docker Engine along with our recommended workaround. You might encounter these errors in Docker Cloud.

---

## Error message: Get i/o timeout

<!-- span tag prevents irritating autolinker from interpreting this as a link -->

*Get https<span></span>://index.docker.io/v1/repositories/\<image\>/images: dial tcp: lookup \<registry host> on \<ip>:53: read udp \<ip>:53: i/o timeout*

#### Description

The DNS resolver configured on the host cannot resolve the registry's hostname.

#### GitHub link

N/A

#### Workaround

Retry the operation, or if the error persists, use another DNS resolver. You can do this by updating your `/etc/resolv.conf` file with these or other DNS servers:

	nameserver 8.8.8.8
	nameserver 8.8.4.4

---

## Error message: 500 Server Error: userland proxy

*500 Server Error: Internal Server Error ("Cannot start container \<id>: Error starting userland proxy: listen tcp 0.0.0.0:\<port>: bind: address already in use")*

#### Description

Docker Cloud is trying to deploy a container publishing a port which is already being used by a process on the host (like the SSH server listening in port `22`).

#### GitHub link

N/A

#### Workaround

Either choose another port, or SSH into the node and manually stop the process which is using the port that you are trying to use.

---

## Error message: 500 Server Error: bind failed

*500 Server Error: Internal Server Error ("Cannot start container \<id>: Bind for 0.0.0.0:\<port> failed: port is already allocated")*

#### Description

Docker Cloud is trying to deploy a container publishing a port which is already used by another container outside of the scope of Docker Cloud.

#### GitHub link

N/A

#### Workaround

Either choose another port, or SSH into the node and manually stop the container which is using the port that you are trying to use.

---

## Error message: 500 Server Error: cannot start, executable not found

*500 Server Error: Internal Server Error ("Cannot start container \<id>: [8] System error: exec: "\<path>": executable file not found in $PATH")*

#### Description

The `run` command you specified for the container does not exist on the container.

#### GitHub link

N/A

#### Workaround

Edit the service to fix the run command.

---

## Error message: Timeout when pulling image from the registry

*Timeout when pulling image from the registry*

#### Description

Timeouts occur when pulling the image takes more than 10 minutes. This can sometimes be caused by the Docker daemon waiting for a nonexistent process while pulling the required image.

#### GitHub link

<a href="https://github.com/docker/docker/issues/12823" target="_blank">docker/docker#12823</a>

#### Workaround

Restart the `dockercloud-agent` service (`sudo service dockercloud-agent restart`) on the node, or restart the node.