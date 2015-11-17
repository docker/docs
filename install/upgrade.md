+++
title = "Upgrade Trusted Registry and the CS Engine"
description = "Upgrade Trusted Registry and the CS Engine"
keywords = ["docker, documentation, about, technology, hub, upgrade, enterprise"]
[menu.main]
parent="smn_dhe_install"
+++


# Upgrade the Trusted Registry and the CS Engine

This document describes the process and steps necessary to upgrade Docker
Trusted Registry and the commercially supported engine (CS engine). When you
first install, the general order is to install the CS engine, then install the
Trusted Registry. However, when you upgrade, you reverse that order.

  **Note**: Ensure when upgrading the Trusted Registry, that you also upgrade to the latest CS Engine.

The following are overall steps, which are explained in detail in this document:

* Upgrade to latest version of Docker Trusted Registry.
* Get the latest version of the CS engine.
* Turn off the Trusted Registry and restart it again with the latest CS engine.
* Make any changes in your configuration.
* Verify you have completed the upgrade process with no errors.

## Upgrade Docker Trusted Registry

Periodic upgrades to the Trusted Registry triggers a notification to appear in your Admin dashboard if you have enabled Upgrade checking in the
General > Settings section of the user interface (UI).

To upgrade, perform the following steps:

1. Load the Trusted Registry Dashboard in your browser and navigate to > Settings > Updates.

2. Click Updates in the Settings nav bar. You can see the currently installed version and a message stating that the version is either current or an update is
available. If an update is available, the message states: System Update
Available and an enabled button displays Update to version X.X.X.

3. Click Update to start the update process. The process may take longer than what the message indicates. To check the status of the install, SSH into the Trusted Registry host through a command line:

      `$ sudo docker logs -f $(sudo docker ps -a --no-trunc | grep 'manager execute-upgrade' | head -n1 | awk '{print $1}')`

4. Refresh your screen to see the latest changes.

      The Dashboard displays a message that the upgrade successfully completed and that you need to upgrade to the latest CS Engine.

### What is updated?

* The Trusted Registry pulls new container images from Docker Hub.
* It then deploys those containers.
* It stops and removes the old containers.

The upgrade process requires a small amount of downtime to complete.

> **Note**: If the CS engine is upgraded first, then
> the Trusted Registry can still be upgraded from the command line by running the following command. Ensure to put the correct version that you want.
>
> `$ sudo bash -c "$(sudo docker run docker/trusted-registry:1.3.3 upgrade 1.4.0)"`


## Upgrade to the latest version of the CS engine

1. Stop the Trusted Registry prior to upgrading the CS engine.

    `$ sudo bash -c "$(sudo docker run docker/trusted-registry:1.4.0 stop)"`

**WARNING**: If you stop the CS Engine, while the Trusted Registry is running, the Trusted Registry may not perform as expected and you must restart it.

### CentOS 7.1 & RHEL 7.0/7.1 (YUM-based systems)

Perform the following commands in your terminal to remove your current CS
engine, and install the new version. When complete, restart the Trusted Registry
again.

1. Remove the current engine:

    `$ sudo yum remove docker-engine-cs`

2. Add Docker's public key for CS packages:

    `$ sudo rpm --import "https://pgp.mit.edu/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"`

3. Install yum-utils if necessary:

    `$ sudo yum install -y yum-utils`

4. Add the repository for the new version and disable the old one:

    ```
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.9/yum/repo/main/centos/7
    $ sudo yum-config-manager --disable 'Docker_cs*'
    ```
5. Install the new package:

    `$ sudo yum install docker-engine`

6. Enable the Docker daemon as a service and then start it.

    ```
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

7. Restart the Trusted Registry:  

    ```
    $ sudo bash -c "$(sudo docker run docker/trusted-registry start)"
    ```

### Ubuntu 14.04 LTS (APT-based systems)

Perform the following commands in your terminal to remove your current CS
engine, and install the new version. When complete, restart the Trusted Registry
again.

1. Remove the current engine:

    `$ sudo apt-get remove docker-engine-cs`

2. Add Docker's public key for CS packages:

    `$ curl -s 'https://pgp.mit.edu/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import`

3. Install the HTTPS helper for apt (your system may already have it):

    `$ sudo apt-get update && sudo apt-get install apt-transport-https`

4. Install additional virtual drivers not in the base image.

        $ sudo apt-get install -y linux-image-extra-virtual

      You may need to reboot your server after updating the LTS kernel.

5. Add the repository for the new version:

    `$ echo "deb https://packages.docker.com/1.9/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list`

        **Note**: modify the "ubuntu-trusty" string for your flavor of ubuntu or debian.
        * debian-jessie (Debian 8)
        * debian-stretch (future release)
        * debian-wheezy (Debian 7)
        * ubuntu-precise (Ubuntu 12.04)
        * ubuntu-trusty (Ubuntu 14.04)
        * ubuntu-utopic (Ubuntu 14.10)
        * ubuntu-vivid (Ubuntu 15.04)
        * ubuntu-wily (Ubuntu 15.10)

6. Install the new package:

    `$ sudo apt-get update && sudo apt-get install docker-engine`

7. Restart the Trusted Registry:  

     `$ sudo bash -c "$(sudo docker run docker/trusted-registry start)"`

## LDAP Configure options

With this release, there are several changes to the LDAP configuration options
that affect authentication and global roles.

* Performance for LDAP user authentication has been significantly increased,
reducing the number of required LDAP requests to only a single BIND request to
authenticate a user.

* The "Read-Write Search Filter" and "Read-Only Search Filter" fields have been
deprecated. You can now create organization accounts and teams in the Trusted
Registry to allow for more fine grained access control. Team member lists can be
synced with a group in LDAP.

* An "Admin Password" is now required. Use this password to login as the user
admin in case the Trusted Registry is unable to authenticate you using your LDAP
server. This account can be used to login to the Trusted Registry and correct
identity and authentication settings.

* Users on your LDAP server are now synced to the Trusted Registry's local
database using your configured "User Search Filter". Objects in LDAP that match
this filter and have a valid "User Login Attribute" are created as a local user
with the "User Login Attribute" as their username. Only these users are able to
login to Docker Trusted Registry.

* The "Admin LDAP DN" must now be specified to identify the group object on your
LDAP server. This should be synced to the system administrators list. The "Admin
Group Member Attribute" should be set to the name of the attribute on this group
object which corresponds to the Distinguished Name of the group member objects.
This setting deprecates the old "Admin Search Filter" field.

## See also

* To configure for your environment, see the
[configuration instructions](../configuration.md).
* To use Docker Trusted Registry, see [the User guide](../userguide.md).
* See [installing the CS engine](install-csengine.md).
* To make administrative changes, see [the Admin guide](../adminguide.md).
* To see previous changes, see [the release notes](../release-notes.md).
