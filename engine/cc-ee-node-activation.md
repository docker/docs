---
description: Engine
keywords: Engine, CE, EE, node activate
title: CE-EE Node Activate
---

The Docker Enterprise Engine for Docker `18.09` release introduces a new feature called CE-EE Node Activate, and allows a user to perform an in-place seamless activation of the EE feature-set on a given node through the Docker Command Line Interface (CLI) without a re-install.

The Docker Enterprise Engine for Docker version must be version `18.09` or higher.

The activation can be performed either online with connection to Docker Hub, or offline with Linux distribution packages.

## Limitations

* This feature is only supported on x86 Linux nodes
* Windows nodes are not currently supported
* Node level Engine activation between CE and EE is only supported in the same version of Docker Enterprise Engine for Docker
* Prior version of the Docker Enterprise Engine for Docker to not support this feature

## Docker Engine 18.09 CE to EE Node Activation Process

1. Check the current Docker version.
```
root@docker-node:~# docker version
Client:
 Version:           18.09.0-ce
 API version:       1.39
 Go version:        go1.10.3
 Git commit:        ca36ebe
 Built:             Thu Aug 23 18:21:56 2018
 OS/Arch:           linux/amd64
 Experimental:      false
Server:
 Engine:
  Version:          18.09.0-ce
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.10.3
  Git commit:       ca36ebe
  Built:
  OS/Arch:          linux/amd64
  Experimental:     false
```

2. Log into the Docker engine from the command line.

```
root@docker-node:~# docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: beluga
Password:
WARNING! Your password will be stored unencrypted in /home/docker/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store
Login Succeeded
```

3. Activate the EE license. You must use sudo even if your user is part of the docker group.

```
root@docker-node:~# sudo docker engine activate --registry-prefix docker.io/dockereng --version 18.09.0-ee
Looking for existing licenses for beluga...
NUM                 OWNER               PRODUCT ID                     EXPIRES                         PRICING COMPONENTS
0                   beluga          docker-ee-trial                2018-09-13 21:41:12 +0000 UTC   Nodes:10
1                   beluga          docker-ee-trial                2018-08-31 03:17:15 +0000 UTC   Nodes:10
2                   beluga          docker-ee                      2018-10-14 15:30:01 +0000 UTC   Linux (IBM Z) Nodes:20,Linux (x86-64) Nodes:20,Windows (x86-64) Nodes:20
3                   beluga          docker-ee-linux                2019-03-19 20:53:37 +0000 UTC   Nodes:10
4                   docker              docker-ee                      2019-05-11 04:33:27 +0000 UTC   Linux (IBM Z) Nodes:100,Windows (x86-64) Nodes:100,Linux (PowerPC) Nodes:100,Linux (x86-64) Nodes:100
5                   docker              docker-ee                      2021-04-01 18:00:17 +0000 UTC   Linux (x86-64) Nodes:10,Windows (x86-64) Nodes:10
6                   docker              docker-ee-server-oraclelinux   2017-08-12 21:55:40 +0000 UTC   Nodes:10
7                   docker              docker-ee-linux                2019-06-23 17:02:48 +0000 UTC   Nodes:1000
8                   docker              docker-ee-server-rhel          2017-05-27 14:04:12 +0000 UTC   Nodes:10
9                   docker              docker-ee-linux                2019-03-31 07:00:00 +0000 UTC   Nodes:10
```

4. Pick the license of your choice
```
Please pick a license by number: 0
waiting for engine to be responsive... engine is online.
```

5. Check the Docker Enterprise Engine for Docker version. The server engine will now be EE, and the client will stay CE.
```
root@docker-node:~# docker version
Client:
 Version:           18.09.0-ce
 API version:       1.39
 Go version:        go1.10.3
 Git commit:        ca36ebe
 Built:             Thu Aug 23 18:21:56 2018
 OS/Arch:           linux/amd64
 Experimental:      false
Server:
 Engine:
  Version:          18.09.0-ee-1
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.10.3
  Git commit:       b9e7996
  Built:
  OS/Arch:          linux/amd64
  Experimental:     false
```
