---
description: Instructions for installing Docker Engine on openSUSE
keywords: requirements, apt, installation, openSUSE, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/opensuse/
- /engine/installation/linux/opensuse/
- /engine/installation/linux/docker-ce/opensuse/
- /install/linux/docker-ce/opensuse/
title: Install Docker Engine on openSUSE
toc_max: 4
---

To get started with Docker Engine on CentOS, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

### OS requirements

To install the docker and docker-compose packages start YaST2, select "Software" and start the module "Software Management". Search for docker and choose to install the Packages "docker" and "python3-docker-compose". (Even though the package is called "python3-docker-compose", it installs "docker-compose" binary). Then click "Accept", and if the installation was successful, "Finish".

To start the docker daemon during boot start YaST2, select "System" and start the module "Services Manager". Select the "docker" service and click "Enable/Disable" and "Start/Stop". To apply your changes click "OK".

To join the docker group that is allowed to use the docker daemon start YaST2, select "Security and Users" and start the module "User and Group Management". Select your user and click "Edit". On the "Details" tab select "docker" in the list of "Additional Groups". Then click "OK" twice.

Now you have to "Log out" of your session and "Log in" again for the changes to take effect. 

### Uninstall old versions

Older versions of Docker were called `docker` or `python3-docker-compose`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo zypper remove docker python3-docker-compose
```

## Installation methods

### with YaST2

To install the docker and docker-compose packages start YaST2, select "Software" and start the module "Software Management". Search for docker and choose to install the Packages "docker" and "python3-docker-compose". (Even though the package is called "python3-docker-compose", it installs "docker-compose" binary). Then click "Accept", and if the installation was successful, "Finish".

To start the docker daemon during boot start YaST2, select "System" and start the module "Services Manager". Select the "docker" service and click "Enable/Disable" and "Start/Stop". To apply your changes click "OK".

To join the docker group that is allowed to use the docker daemon start YaST2, select "Security and Users" and start the module "User and Group Management". Select your user and click "Edit". On the "Details" tab select "docker" in the list of "Additional Groups". Then click "OK" twice.

Now you have to "Log out" of your session and "Log in" again for the changes to take effect. 


### with Command line

To install the docker and docker-compose packages: 

```bash
$ zypper install docker python3-docker-compose

```
To start the docker daemon during boot: 

```bash
$ sudo systemctl enable docker

```
To join the docker group that is allowed to use the docker daemon:

```bash
$ sudo usermod -G docker -a $USER

```
Restart the docker daemon: 
```bash
$ sudo systemctl restart docker

```
Verify docker is running: 

```bash
$ docker version

```

This will pull down and run the, "Hello World" docker container from dockerhub: 
$ docker run --rm hello-world

```
Clean up and remove docker image we pulled down: 

```bash
$ docker images
docker rmi -f IMAGE_ID

```
Where "IMAGE_ID" is the Id value of the "Hello World" container. 
