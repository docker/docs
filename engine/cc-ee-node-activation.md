---
description: Engine
keywords: Engine, CE, EE, node, activation, enterprise
title: Enterprise Node Activation
---

The Docker Engine release `18.09` release introduces a new feature called CE-EE Node Activate, which allows a user to perform an in-place seamless activation of the Enterprise engine feature-set on an existing Community Edition (CE) node through the Docker command line.

Enterprise Node Activation will apply a license, and if you aren't already running the Enterprise engine, and switch the Docker engine to the enterprise engine binary.

The Docker Community Edition version must be 18.09 or higher.

The activation can be performed either online with connection to Docker Hub, or offline.

## Limitations

* This feature is only supported on x86 Linux nodes
* Windows nodes are not currently supported
* Node level Engine activation between CE and EE is only supported in the same version of Docker Enterprise Engine for Docker
* Prior versions of Docker CE do not support this feature

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
**NOTE**: Your output may vary slightly from what is displayed on this page.

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
root@docker-node:~# sudo docker engine activate
Looking for existing licenses for beluga...
NUM                 OWNER               PRODUCT ID                     EXPIRES                         PRICING COMPONENTS
0                   beluga          docker-ee-trial                2018-09-13 21:41:12 +0000 UTC   Nodes:10
1                   beluga          docker-ee-trial                2018-08-31 03:17:15 +0000 UTC   Nodes:10
```

4. Pick the license of your choice
```
Please pick a license by number: 0
waiting for engine to be responsive... engine is online.
```

5. Check the Docker Engine version. The server engine will now be EE, and the client will stay CE.
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

**NOTE**: Your output may vary slightly from what is displayed on this page.

## Off-line CE-EE node activation

For offline CE-EE node activation, you'll need to get the enterprise engine onto the system. The recommended model is to download the EE deb or rpm packages manually and copy them to the target systems. Afterward, download the license manually, and copy that license to the target systems. Use the `--license <path/to/license.file>` command line option to the activate command.
