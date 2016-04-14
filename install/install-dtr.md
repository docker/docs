<!--[metadata]>
+++
aliases = [ "/docker-trusted-registry/install/dtr-ami-byol-launch/",
            "/docker-trusted-registry/install/dtr-ami-bds-launch/",
            "/docker-trusted-registry/install/dtr-vhd-azure/"]
title = "Install Docker Trusted Registry"
description = "Learn how to install Docker Trusted Registry for production."
keywords = ["docker, dtr, registry, install"]
[menu.main]
parent="workw_dtr_install"
weight=20
+++
<![end-metadata]-->


# Install Docker Trusted Registry

This document describes the process of obtaining, installing, and securing
Docker Trusted Registry. You can use these instructions if you are installing Trusted Registry on a physical or cloud infrastructure.

If your cloud provider is AWS, you have the option of installing Trusted Registry using an Amazon Machine Image (AMI) instead. For more information, read the [installation overview](index.md) to understand your options.


## Prerequisites

Docker Trusted Registry runs on the following 64-bit platforms:

* Ubuntu 14.04 LTS
* RHEL 7.0 and 7.1
* CentOS 7.1
* SUSE Linux Enterprise 12

Docker Trusted Registry requires the latest commercially supported Docker Engine (CS Engine), running on a supported host.

The Docker daemon listens to the Unix socket (the default) so that it can be
bind-mounted into the Trusted Registry management containers. This allows
Trusted Registry to manage itself and its updates. For this reason, the host you
install on needs internet connectivity so it can access the updates.
Additionally, your host needs to have TCP ports `80` and `443` available for the
Docker Trusted Registry container port mapping.

Installing Trusted Registry requires that you have a login to Docker Hub (or the
user-name of an administrator of the Hub organization that obtained an
Enterprise license. If you already installed CS Engine, you should already have a [Hub account](https://hub.docker.com).

Also, you must have a license for Docker Trusted Registry. This license allows
you to run both Docker Trusted Registry and CS Engine. Before installing,
[purchase a license or sign up for a free, 30 day trial license](https://hub.docker.com/enterprise/).


## Install Docker Trusted Registry

Trusted Registry is a self-installing application built and distributed using
Docker and the [Docker Hub](https://hub.docker.com/). You install Docker Trusted
Registry by running the "docker/trusted-registry" container. Once installed, it
is able to restart and reconfigure itself using the Docker socket that is
bind-mounted to this container.

1. Log in to the machine where you want to install Trusted Registry.

2. Verify that CS Engine is installed.

        $ docker --version

    > **Note:** To remain compliant with your Docker Trusted Registry support agreement, you **must** use the current version of commercially supported Docker Engine. Running the open source version of Engine is **not** supported.

3. Login into the Docker Hub from the command line.

        $ docker login

4. Install the Trusted Registry

	       $ sudo bash -c "$(sudo docker run docker/trusted-registry install)"

    > **Note**: `sudo` is needed for `docker/trusted-registry` commands to
    > ensure that the Bash script is run with full access to the Docker host.

    The command executes a shell script that creates the needed directories,
    pulls the registry's images, and run its containers. Depending on your
    internet connection, this process may take several minutes to complete. A successful outcome completes as follows:

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
          Status: Downloaded newer image for docker/trusted-registry:latest
          Unable to find image 'docker/trusted-registry:1.1.0' locally
          Pulling repository docker/trusted-registry
          c46d58daad7d: Download complete
          511136ea3c5a: Download complete
          ...
          Status: Image is up to date for docker/trusted-registry:1.1.0
          INFO  [1.0.0_8ce62a61e058] Attempting to connect to docker engine dockerHost="unix:///var/run/docker.sock"
          INFO  [1.0.0_8ce62a61e058] Running install command
          <...output truncated...>
          Creating container docker_trusted_registry_load_balancer with docker daemon unix:///var/run/docker.sock
          Starting container docker_trusted_registry_load_balancer with docker daemon unix:///var/run/docker.sock
          Bringing up docker_trusted_registry_log_aggregator.
          Creating container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock
          Starting container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock

5. Use `docker ps` to list all the running containers.

    The listing should show the following were started:

  * `docker_trusted_registry_load_balancer`
  * `docker_trusted_registry_image_storage_0`
  * `docker_trusted_registry_image_storage_1`
  * `docker_trusted_registry_admin_server`
  * `docker_trusted_registry_log_aggregator`
  * `docker_trusted_registry_auth_server`
  * `docker_trusted_registry_postgres`

6. Enter the `https://<host-ip>/` your browser's address bar to run the Trusted Registry interface.

  Your browser warns you that this is an unsafe site, with a self-signed,
  untrusted certificate. This is normal and expected; allow this connection
  temporarily.


## Set the Trusted Registry domain name

The Docker Trusted Registry Administrator site will also warn that the "Domain Name" is not set.

1. Select "Settings" from the global nav bar at the top of the page, and then set the "Domain Name" to the full host-name of your Docker Trusted Registry server.

2. Click the "Save and Restart Docker Trusted Registry Server" button to generate a new certificate, which will be used
by both the Docker Trusted Registry Administrator web interface and the Docker Trusted Registry server.

3. After the server restarts, you will again need to allow the connection to the untrusted Docker Trusted Registry web admin site.

4. You see a warning notification that this instance of Docker Trusted Registry is unlicensed. You'll correct this in the next section.

## Apply your license

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

## Secure the Trusted Registry

Securing Docker Trusted Registry is **required**. You will not be able to push
or pull from Docker Trusted Registry until you secure it.

There are several options and methods for securing Docker Trusted Registry. For
more information, see the [configuration documentation](../configure/configuration.md#security)

## Push and pull images

You have your Trusted Registry configured with a "Domain Name" and your
client Docker daemons configured with the required security settings. But
before you can test your setup by pushing an image, you need to create a
repository first. Follow the instructions for
[Using Docker Trusted Registry to Push and pull images](../repos-and-images/push-and-pull-images.md)
to create a repository and to push and pull images.

## Docker Trusted Registry web interface and registry authentication

By default, there is no authentication set on either the Docker Trusted Registry
web admin interface or the Docker Trusted Registry. You can restrict access
using an in-Docker Trusted Registry configured set of users (and passwords), or
you can configure Docker Trusted Registry to use LDAP- based authentication.

See [Docker Trusted Registry Authentication settings](../configure/configuration.md#authentication) for more details.

## See also

* [Install DTR offline](install-dtr-offline.md)
* [Upgrade DTR](upgrade.md)
