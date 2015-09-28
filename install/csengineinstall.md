+++
title = "Manually Install the CS Docker Engine"
description = "Install instructions for the commercially supported Docker Engine"
keywords = ["docker, documentation, about, technology, enterprise, hub, commercially supported Docker Engine, CS engine, registry"]
[menu.main]
parent="smn_dhe_install"
+++


# Manually Install the CS Docker Engine

You must first install the latest version of the commercially supported Docker Engine (CS engine) before you install the latest version of Docker Trusted Registry. Do this with an RPM or DEB package which you access using a script downloaded from your
[Docker Hub Licenses page](https://hub-beta.docker.com/account/licenses/).

If you have not already done so, make sure you have first read the [installation overview](index.md) for Trusted Registry.

## Download the installation script

1. To download the commercially supported Docker Engine Bash installation script,
log in to the [Docker Hub](https://hub-beta.docker.com) with the user-name used to
obtain your license.
2. Once you're logged in, go to the
[Licenses](https://hub-beta.docker.com/account/licenses/) page in your Hub account's Settings section. This is accessed through the gear icon located in the upper right of your page.
3. Click the button at the top right of the Licenses page that corresponds to your intended host operating system. The Bash setup script
downloads.
4. Next, select your operating system and follow the applicable steps.

### Install CentOS 7.1 & RHEL 7.0/7.1

1. Copy the downloaded Bash setup script to your RHEL host.
2. Run the following to install the engine and its dependencies.
3. Start the Docker daemon.

  The following example illustrates these steps.

  ```
  $ sudo yum update && sudo yum upgrade
  $ chmod 755 docker-cs-engine-rpm.sh
  $ sudo ./docker-cs-engine-rpm.sh
  $ sudo yum install docker-engine-cs
  $ sudo systemctl enable docker.service
  $ sudo systemctl start docker.service
  ```

  To make using Docker easier, you can get non-sudo access to the Docker
  socket by adding your user to the `docker` group.

(*Optional*) Log out and back in again as seen in the next example.


```
$ sudo usermod -a -G docker $USER
$ exit
```
> **Notes**:
>   You may need to reboot your server to update its RHEL kernel.
>   CentOS 7.0 is not supported.

### Install Ubuntu 14.04 LTS

1. Copy the downloaded Bash setup script to your Ubuntu host.
2. Run the following to install the engine and its dependencies.

      The following example illustrates these steps.

      ```
      $ sudo apt-get update && sudo apt-get upgrade
      $ sudo apt-get install -y linux-image-extra-virtual
      $ sudo reboot
      $ chmod 755 docker-cs-engine-deb.sh
      $ sudo ./docker-cs-engine-deb.sh
      $ sudo apt-get install docker-engine-cs
      ```

3. Confirm Docker is running with `sudo service docker start`.

        To make using Docker easier, you can get non-sudo access to the Docker socket by adding your user to the `docker` group.

4. (*Optional*) Log out and back in again as seen in the next example.

  ```
  $ sudo usermod -a -G docker $USER
  $ exit
  ```

> **Note**: You may need to reboot your server to update its LTS kernel.

## Next step
Your ready to install [Docker Trusted Registry]({{< relref "install.md" >}}).

## See also

* If you were manually installing Docker Trusted Registry, then continue with [those steps]({{< relref "install.md" >}}).
* To configure for your environment, see the
[Configuration instructions]({{< relref "configuration.md" >}}).
* To use Docker Trusted Registry, see the [User guide]({{< relref "userguide.md" >}}).
* To make administrative changes, see the [Admin guide]({{< relref "adminguide.md" >}}).
* To see previous changes, see the [release notes]({{< relref "release-notes.md" >}}).

## See also

* To configure for your environment, see
[Configuration instructions]({{< relref "configuration.md" >}}).
* To use Docker Trusted Registry, see [the User guide]({{< relref "userguide.md" >}}).
* To make administrative changes, see [the Admin guide]({{< relref "adminguide.md" >}}).
* To see previous changes, see [the release notes]({{< relref "release-notes.md" >}}).
