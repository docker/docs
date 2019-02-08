---
description: Engine
keywords: Engine, CE, EE, node, activation, enterprise, patching
title: CE-EE Node Activate
---

The Docker Engine `18.09` release introduces a new feature called **CE-EE Node Activate**, which allows a user to perform an in-place seamless activation of the Enterprise engine feature set on an existing Community Edition (CE) node through the Docker command line.

CE-EE Node Activate will apply a license, and switch the Docker engine to the Enterprise engine binary.

## Requirements:
* The Docker Community Edition (CE) version must be 18.09 or higher.
* All of the Docker packages must be installed: docker-cli, docker-server, and containerd.

## Limitations

* This feature is only supported on x86 Linux nodes.
* Windows nodes are not currently supported.
* Node-level Engine activation between CE and EE is only supported in the same version of Docker Enterprise Engine for Docker.
* Prior versions of Docker CE do not support this feature.

## Notes on patching after running CE to EE Node Activation

Docker recommends replacing the apt or yum repository from CE with the EE repository that appears 
on your hub/store account after starting the trial or paid license. This allows apt/yum 
upgrade operations to work as expected and keep them current as long as your license is still
valid and has not expired.

> **Note**: You can use the `docker engine update` command. However, if you continue to use
> the CE packages, the OS package will no longer replace the active daemon binary during apt/yum 
> updates, so you are responsible for performing the `docker engine update` operation periodically 
> to keep your engine up to date.

## Docker Engine 18.09 CE to EE Node Activation Process

The activation can be performed either online with connection to Docker Hub, or offline.

1. Check the current Docker version. Both the Docker client and server (`containerd`) need to be installed.  Your output may vary slightly from what is displayed on this page.

```
$ docker version
Client:
 Version:           18.09.0
 API version:       1.39
 Go version:        go1.10.4
 Git commit:        4d60db4
 Built:             Wed Nov  7 00:48:22 2018
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          18.09.0
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.10.4
  Git commit:       4d60db4
  Built:            Wed Nov  7 00:19:08 2018
  OS/Arch:          linux/amd64
  Experimental:     false
```

2. Log into the Docker engine from the command line.

> **Note**: When running the command `docker login`, the shell stores the credentials in the current user's home
> directory. RHEL and Ubuntu-based Linux distributions have different behavior for sudo. RHEL sets $HOME to point
> to `/root` while Ubuntu leaves `$HOME` pointing to the user's home directory who ran `sudo` and this can cause
> permission and access problems when switching between `sudo` and non-sudo'd commands.


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

For CentOS, use `sudo`:
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

4. Activate the EE license. You must use `sudo` even if your user is part of the docker group.

```
$ sudo docker engine activate --license ee-engine-license-1000-nodes.lic
License: Quantity: 400 Nodes	Expiration date: 2019-05-12	License is currently active
18.09.0: resolved
e3cba72cdf2e: done [==================================================>]     434B/434B
3642e2b52398: done [==================================================>]  1.161kB/1.161kB
eb6fe2abc4ad: done [==================================================>]  4.544MB/4.544MB
f2f08b0292f2: done [==================================================>]  25.65MB/25.65MB
a539281ee17b: done [==================================================>]  1.122MB/1.122MB
515c4dc2b0fe: done [==================================================>]  333.9kB/333.9kB
2cf04a6ee63e: done [==================================================>]   4.84kB/4.84kB
Successfully activated engine.
Restart docker with 'systemctl restart docker' to complete the activation.
```

5. Check the Docker Engine version. The engine server will become EE, and the engine client will stay CE.

```
$ docker version
Client:
 Version:           18.09.0
 API version:       1.39
 Go version:        go1.10.4
 Git commit:        4d60db4
 Built:             Wed Nov  7 00:48:22 2018
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Enterprise
 Engine:
  Version:          18.09.0
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.10.4
  Git commit:       33a45cd
  Built:            Wed Nov  7 00:19:46 2018
  OS/Arch:          linux/amd64
  Experimental:     false
```

**NOTE**: Your output may vary slightly from what is displayed on this page.

6. If you are running a Swarm cluster with CE, please repeat these steps on each node.

## Offline CE-EE node activation

For offline CE-EE node activation, you'll need to get the Docker Enterprise Engine onto the system. The recommended model is to download the EE `.deb` or `.rpm` packages manually and copy them to the target systems. Afterward, download the license manually, and copy that license to the target systems. Use the `--license <path/to/license.file>` command line option to the activate command.
