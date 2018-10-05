---
description: Engine
keywords: Engine, CE, EE, node, activation, enterprise
title: CE-EE Node Activate
---

The Docker Engine release `18.09` release introduces a new feature called CE-EE Node Activate, which allows a user to perform an in-place seamless activation of the Enterprise engine feature-set on an existing Community Edition (CE) node through the Docker command line.

CE-EE Node Activate will apply a license, and if you aren't already running the Enterprise engine, and switch the Docker engine to the enterprise engine binary.

The Docker Community Edition version must be 18.09 or higher.

The activation can be performed either online with connection to Docker Hub, or offline.

## Limitations

* This feature is only supported on x86 Linux nodes
* Windows nodes are not currently supported
* Node level Engine activation between CE and EE is only supported in the same version of Docker Enterprise Engine for Docker
* Prior versions of Docker CE do not support this feature

## Docker Engine 18.09 CE to EE Node Activation Process

1. Check the current Docker version. Both the Docker client and server (`containerd`) need to be installed.  Your output may vary slightly from what is displayed on this page.

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

**NOTE:** When the running the command `docker login`, the shell stores the credentials in the current users's home directory. For CentOS and Red Hat, the `sudo` command overwrites overwrites the $HOME environment variable.


For Ubuntu or Debian:
```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: beluga
Password:
WARNING! Your password will be stored unencrypted in /home/docker/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store
Login Succeeded
```

For CentOS or Red Hat Linux, use `sudo`:
```
$ sudo docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: beluga
Password:
WARNING! Your password will be stored unencrypted in /home/docker/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store
Login Succeeded
```

3. [Download your Docker Enterprise license](https://success.docker.com/article/where-is-my-docker-enterprise-edition-license) and distribute it to your Docker engines.

4. Activate the EE license. You must use sudo even if your user is part of the docker group.

```
root@docker-node:~# sudo docker engine actviate --license dockersub.lic
Looking for existing licenses for beluga...
NUM                 OWNER               PRODUCT ID                     EXPIRES                         PRICING COMPONENTS
0                   beluga          docker-ee-trial                2018-09-13 21:41:12 +0000 UTC   Nodes:10
1                   beluga          docker-ee-trial                2018-08-31 03:17:15 +0000 UTC   Nodes:10
```

**NOTE:** If the Docker EE engines are in a swarm cluster, you only need to activate the license on the 
manager node. This action stores the license as a swarm configuration, which is compatible with UCP.

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
