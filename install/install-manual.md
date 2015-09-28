+++
title = "Manually install Trusted Registry"
description = "Manually install Trusted Registry"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub,  registry"]
[menu.main]
parent="smn_dhe_install"
+++


# Manually install Trusted Registry

This document describes the process of obtaining, installing, and securing
Docker Trusted Registry. Trusted Registry is installed from Docker containers.
If you have not already done so, make sure you have first read the [installation
overview](index.md) for Trusted Registry.  

In general, the steps for installing the Engine for AWS AMI (BDS) are:

1. Acquire a license by purchasing Docker Trusted Registry or signing up for a trial license.
2. Install the commercially supported Docker Engine.
3. Install Docker Trusted Registry.
4. Apply your license to your Docker Trusted Registry instance.
5. Configure your settings.


## Overview
Docker Trusted Registry runs on the following platforms:

* Ubuntu 14.04 LTS
* RHEL 7.0 and 7.1
* CentOS 7.1

Docker Trusted Registry requires the following:

* The latest commercially supported Docker Engine, running
on a supported host.

> **Note:** In order to remain in compliance with your Docker Trusted Registry support agreement,
> you **must** use the current version of commercially supported Docker Engine.
> Running the open source version of Engine is **not** supported.

* The Docker daemon needs to be listening to the unix socket (the default) so
that it can be bind-mounted into the Docker Trusted Registry management containers, allowing
Docker Trusted Registry to manage itself and its updates. For this reason, your Docker Trusted Registry host will also
need internet connectivity so it can access the updates.

* Your host needs to have TCP ports `80` and `443` available for the Docker Trusted Registry
container port mapping.

* You need the Docker Hub user-name and password used when obtaining
the Docker Trusted Registry license (or the user-name of an administrator of the Hub organization
that obtained an Enterprise license).

## Get your license

To start, get your copy of Docker Trusted Registry, including a free trial, at the [Docker Subscription page](https://hub-beta.docker.com/enterprise/).

Go to the license [documentation]({{< relref "license.md" >}}) to read about how to download the license. When you're finished, continue to the next step.

### Install the commercially supported Docker Engine

You download and save the CS engine install script from the Licenses page. Follow these [cs engine install ]({{< relref "csengineinstall.md" >}}) directions and then return to this document to continue the install.

#### CentOS 7.1 & RHEL 7.0/7.1 installation

1. Copy the downloaded Bash setup script to your RHEL host.
2. Run the following to install commercially supported Docker Engine and its dependencies.
3. Start the Docker daemon:

```
$ sudo yum update && sudo yum upgrade
$ chmod 755 docker-cs-engine-rpm.sh
$ sudo bash ./docker-cs-engine-rpm.sh
$ sudo yum install docker-engine-cs
$ sudo systemctl enable docker.service
$ sudo systemctl start docker.service
```

In order to simplify using Docker, you can get non-sudo access to the Docker
socket by adding your user to the `docker` group, then logging out and back in
again:

```
$ sudo usermod -a -G docker $USER
$ exit
```

> **Notes**:
>      * You may need to reboot your server to update its RHEL kernel.    
>      * CentOS 7.0 is not supported.


#### Ubuntu 14.04 LTS installation

1. Copy the downloaded Bash setup script to your Ubuntu host.
2. Run the following to install commercially supported Docker Engine and its dependencies:


          $ sudo apt-get update && sudo apt-get upgrade
          $ sudo apt-get install -y linux-image-extra-virtual
          $ sudo reboot
          $ chmod 755 docker-cs-engine-deb.sh
          $ sudo ./docker-cs-engine-deb.sh
          $ sudo apt-get install docker-engine-cs

3. Confirm Docker is running with `sudo service docker start`.

4. (Optional) In order to simplify using Docker, get non-sudo access to the Docker socket by adding your user to the `docker` group, then logging out and back in again:

```
$ sudo usermod -a -G docker $USER
$ exit
```

> **Note**: you may need to reboot your server to update its LTS kernel.


## Install Docker Trusted Registry

Once the commercially supported Docker Engine is installed, you can install Docker Trusted Registry
itself. Docker Trusted Registry is a self-installing application built and distributed using Docker
and the [Docker Hub](https://hub-beta.docker.com/). It is able to restart
and reconfigure itself using the Docker socket that is bind-mounted to its
container.

Install Docker Trusted Registry by running the "docker/trusted-registry" container:

```
	$ sudo bash -c "$(sudo docker run docker/trusted-registry install)"
```

> **Note**: `sudo` is needed for `docker/trusted-registry` commands to
> ensure that the Bash script is run with full access to the Docker host.

The command will execute a shell script that creates the needed
directories and then runs Docker to pull Docker Trusted Registry's images and run its containers.

Depending on your internet connection, this process may take several minutes to
complete.

A successful installation will pull a large number of Docker images and should
display output similar to:

```
$ sudo bash -c "$(sudo docker run docker/trusted-registry install)"
Unable to find image 'docker/trusted-registry:latest' locally
Pulling repository docker/trusted-registry
c46d58daad7d: Pulling image (latest) from docker/trusted-registry
c46d58daad7d: Pulling image (latest) from docker/trusted-registry
c46d58daad7d: Pulling dependent layers
511136ea3c5a: Download complete
fa4fd76b09ce: Pulling metadata
fa4fd76b09ce: Pulling fs layer
ff2996b1faed: Download complete
...
fd7612809d57: Pulling metadata
fd7612809d57: Pulling fs layer
fd7612809d57: Download complete
c46d58daad7d: Pulling metadata
c46d58daad7d: Pulling fs layer
c46d58daad7d: Download complete
c46d58daad7d: Download complete
Status: Downloaded newer image for docker/trusted-registry:latest
Unable to find image 'docker/trusted-registry:1.1.0' locally
Pulling repository docker/trusted-registry
c46d58daad7d: Download complete
511136ea3c5a: Download complete
fa4fd76b09ce: Download complete
1c8294cc5160: Download complete
117ee323aaa9: Download complete
2d24f826cb16: Download complete
33bfc1956932: Download complete
48f0dd6c9414: Download complete
65c30f72ecb2: Download complete
d4b29764d0d3: Download complete
5654f4fe5384: Download complete
9b9faa6ecd11: Download complete
0c275f56ca5c: Download complete
ff2996b1faed: Download complete
fd7612809d57: Download complete
Status: Image is up to date for docker/trusted-registry:1.1.0
INFO  [1.0.0_8ce62a61e058] Attempting to connect to docker engine dockerHost="unix:///var/run/docker.sock"
INFO  [1.0.0_8ce62a61e058] Running install command
<...output truncated...>
Creating container docker_trusted_registry_load_balancer with docker daemon unix:///var/run/docker.sock
Starting container docker_trusted_registry_load_balancer with docker daemon unix:///var/run/docker.sock
Bringing up docker_trusted_registry_log_aggregator.
Creating container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock
Starting container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock
$ docker ps
CONTAINER ID        IMAGE                                          COMMAND                CREATED             STATUS              PORTS                                      NAMES
963ec2a4b047        docker/trusted-registry-nginx:1.1.0            "nginxWatcher"         5 minutes ago       Up 5 minutes        0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   docker_trusted_registry_load_balancer
7eade5529049        docker/trusted-registry-distribution:v2.0.1    "registry /config/st   5 minutes ago       Up 5 minutes        5000/tcp                                   docker_trusted_registry_image_storage_0
b968a8a986f9        docker/trusted-registry-distribution:v2.0.1    "registry /config/st   5 minutes ago       Up 5 minutes        5000/tcp                                   docker_trusted_registry_image_storage_1
390d9d68a33a        docker/trusted-registry-admin-server:1.1.0     "server"               5 minutes ago       Up 5 minutes        80/tcp                                     docker_trusted_registry_admin_server
3f8a53dc5f35        docker/trusted-registry-log-aggregator:1.1.0   "log-aggregator"       5 minutes ago       Up 5 minutes                                                   docker_trusted_registry_log_aggregator
44083421fa16        docker/trusted-registry-garant:1.1.0           "garant /config/gara   5 minutes ago       Up 5 minutes                                                   docker_trusted_registry_auth_server
c4102adf73dc        postgres:9.4.1                                 "/docker-entrypoint.   5 minutes ago       Up 5 minutes        5432/tcp                                   docker_trusted_registry_postgres
```

Once this process completes, you can manage and configure your Docker Trusted Registry instance by pointing your browser to `https://<host-ip>/`.

Your browser warns you that this is an unsafe site, with a self-signed,
untrusted certificate. This is normal and expected; allow this connection
temporarily.

### Setting the Docker Trusted Registry Domain Name

The Docker Trusted Registry Administrator site will also warn that the "Domain Name" is not set.

1. Select "Settings" from the global nav bar at the top of the page, and then set the "Domain Name" to the full host-name of your Docker Trusted Registry server.

2. Click the "Save and Restart Docker Trusted Registry Server" button to generate a new certificate, which will be used
by both the Docker Trusted Registry Administrator web interface and the Docker Trusted Registry server.

3. After the server restarts, you will again need to allow the connection to the untrusted Docker Trusted Registry web admin site.

4. You see a warning notification that this instance of Docker Trusted Registry is unlicensed. You'll correct this in the next section.

### Apply your license

The Docker Trusted Registry services will not start until you apply your license.
To do that, you'll first download your license from the Docker Hub and then
upload it to your Docker Trusted Registry web admin server. Follow these steps:

1. If needed, log back into the [Docker Hub](https://hub.docker.com)
   using the user-name you used when obtaining your license. Under your name, go to Settings to display the Account Settings page. Click the Licenses submenu to display the Licenses page.

2. There is a list of available licenses. Click the download button to
   obtain the license file you want.

3. Go to your Docker Trusted Registry instance in your browser, click Settings in the global nav bar. Click License in the Settings nav bar. Click the Choose File button. It opens a standard file browser. Locate and select the license file you downloaded in the previous step. Approve the selection to close the dialog.

4. Click the Save and restart button. Docker Trusted Registry quits and then restarts with the applied the license.

5. Verify the acceptance of the license by confirming that the "Unlicensed copy"
   warning is no longer present.

### Securing Docker Trusted Registry

Securing Docker Trusted Registry is **required**. You will not be able to push or pull from Docker Trusted Registry until you secure it.

There are several options and methods for securing Docker Trusted Registry. For more information,
see the [configuration documentation]({{< relref "configuration.md#security" >}})

### Using Docker Trusted Registry to push and pull images

Now that you have Docker Trusted Registry configured with a "Domain Name" and have your client
Docker daemons configured with the required security settings, you can test your
setup by following the instructions for
[Using Docker Trusted Registry to Push and pull images]({{< relref "userguide.md" >}}).

### Docker Trusted Registry web interface and registry authentication

By default, there is no authentication set on either the Docker Trusted Registry web admin
interface or the Docker Trusted Registry. You can restrict access using an in-Docker Trusted Registry
configured set of users (and passwords), or you can configure Docker Trusted Registry to use LDAP-
based authentication.

See [Docker Trusted Registry Authentication settings]({{< relref "configuration.md#authentication" >}}) for more
details.

## See also

* To configure for your environment, see the
[configuration instructions]({{< relref "configuration.md" >}}).
* To use Docker Trusted Registry, see [the User guide]({{< relref "userguide.md" >}}).
* To make administrative changes, see [the Admin guide]({{< relref "adminguide.md" >}}).
* To see previous changes, see [the release notes]({{< relref "release-notes.md" >}}).
