+++
title = "Install the Trusted Registry offline"
description = "Install the Trusted Registry offline"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, offline, Trusted Registry, registry"]
[menu.main]
parent="workw_dtr_install"
+++


# Install the Trusted Registry offline

This document describes the process of obtaining, installing, and securing
Docker Trusted Registry offline. Since your system is not connected to the internet, there will be no notifications regarding upgrading either the CS Engine or the Trusted Registry. You will also not be able to link from the Trusted Registry UI to our documentation except for the API documentation. Docker recommends that you contact customer support to obtain the latest information.

For more information about installing, read the
[installation overview](index.md) to understand your options.

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

1. Since you are retrieving a large file, use the `wget` command in your command line to get the Trusted Registry files. The following command is an example getting DTR 1.4.3. Ensure to get your correct version.

      `wget https://packages.docker.com/dtr/1.4/dtr-1.4.3.tar`

2. After downloading, move the `tar` file to the offline machine you want to install the Trusted Registry.

3. On that machine, verify that the CS Engine is installed. If it is not, see the [CS Engine install directions](install-csengine.md).

    `$ docker --version`

    > **Note:** To remain compliant with your Docker Trusted Registry support agreement, you **must** use the current version of commercially supported Docker Engine. Running the open source version of Engine is **not** supported.

5. Open a terminal window on that machine and load the `tar` file using the following command. Again, ensure you get the correct version.

      `$ sudo docker load < dtr-1.4.3.tar`

6. Install the Trusted Registry with the following command:

      `$ sudo bash -c "$(sudo docker run docker/trusted-registry install)"`


    > **Note**: `sudo` is needed for `docker/trusted-registry` commands to
    > ensure that the Bash script is run with full access to the Docker host.

    The command runs the registry's containers from the images you loaded in the previous step. You will know that you successfully installed by the following in part:

    Image is up to date for docker/trusted-registry:1.4.3


    ```
    Checking for required image: docker/trusted-registry-distribution:v2.2.1
    Checking for required image: postgres:9.4.1
    ...
    INFO  [1.4.3-003501_g657863b] Attempting to connect to docker engine dockerHost="unix:///var/run/docker.sock"
    INFO  [1.4.3-003501_g657863b] Running install command
    INFO  [1.4.3-003501_g657863b] Running pull command
    INFO  [1.4.3-003501_g657863b] Using links? false
    INFO  [1.4.3-003501_g657863b] DTR Network created
    Bringing up docker_trusted_registry_postgres.
    Creating container docker_trusted_registry_postgres with docker daemon unix:///var/run/docker.sock
    Starting container docker_trusted_registry_postgres with docker daemon unix:///var/run/docker.sock
    ...
    Bringing up docker_trusted_registry_log_aggregator.
    Creating container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock
    Starting container docker_trusted_registry_log_aggregator with docker daemon unix:///var/run/docker.sock
    Bringing up docker_trusted_registry_auth_server.
    Creating container docker_trusted_registry_auth_server with docker daemon unix:///var/run/docker.sock
    Starting container docker_trusted_registry_auth_server with docker daemon unix:///var/run/docker.sock
    Bringing up docker_trusted_registry_postgres.
    Creating container docker_trusted_registry_postgres with docker daemon unix:///var/run/docker.sock
    Container already exists for daemon at unix:///var/run/docker.sock: docker_trusted_registry_postgres
    Starting container docker_trusted_registry_postgres with docker daemon unix:///var/run/docker.sock
    Container docker_trusted_registry_postgres is already running for daemon at unix:///var/run/docker.sock
    ```

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

4. Click Save and restart. Docker Trusted Registry quits and then restarts with the applied the license.

5. Verify the acceptance of the license by confirming that the "Unlicensed copy"
   warning is no longer present.

## Secure the Trusted Registry

Securing Docker Trusted Registry is **required**. You will not be able to push
or pull from Docker Trusted Registry until you secure it.

There are several options and methods for securing Docker Trusted Registry. For
more information, see the [configuration documentation](../configuration.md#security)

## Push and pull images

You have your Trusted Registry configured with a "Domain Name" and your
client Docker daemons configured with the required security settings. But
before you can test your setup by pushing an image, you need to create a repository first. Follow the instructions for [Using Docker
Trusted Registry to Push and pull images](../userguide.md) to create a repository and to push and pull images.

## Docker Trusted Registry web interface and registry authentication

By default, there is no authentication set on either the Docker Trusted Registry
web admin interface or the Docker Trusted Registry. You can restrict access
using an in-Docker Trusted Registry configured set of users (and passwords), or
you can configure Docker Trusted Registry to use LDAP based authentication.

See [Docker Trusted Registry Authentication settings](../configuration.md#authentication) for more details.

## See also

* To configure for your environment, see the
[configuration instructions](../configuration.md).
* To use Docker Trusted Registry, see [the User guide](../userguide.md).
* To make administrative changes, see [the Admin guide](../adminguide.md).
* To see previous changes, see [the release notes](../release-notes.md).
