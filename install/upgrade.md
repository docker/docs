+++
title = "Upgrade Trusted Registry and CS Engine"
description = "Upgrade Trusted Registry and CS Engine"
keywords = ["docker, documentation, about, technology, hub, upgrade, enterprise"]
[menu.main]
parent="smn_dhe_install"
+++

<!--
not certain of this order. and there may be missing steps if the license has to be enabled again. cfh
-->


# Upgrade the Trusted Registry and the CS Engine

This document describes the process and steps necessary to upgrade Docker Trusted Registry and the commercially supported engine (CS engine). The general steps are to:

* Get the latest version of the CS engine and install it.
* Get the latest version of Docker Trusted Registry.
* Turn off the Trusted Registry and restart it again with the latest CS engine.
* Make any changes in your configuration.
* Verify you have completed the upgrade process with no errors.

## Upgrade to the latest version of the CS engine

The CS engine installation script set up the RHEL/Ubuntu package repositories,
so upgrading the CS engine only requires you to run the update commands on your server.

### CentOS 7.1 & RHEL 7.0/7.1 upgrade

The following commands will stop the running Docker Trusted Registry, upgrade CS  engine,
and then start the Trusted Registry again:

```
    $ sudo bash -c "$(sudo docker run docker/trusted-registry stop)"
    $ sudo yum update
    $ sudo systemctl daemon-reload && sudo systemctl restart docker
    $ sudo bash -c "$(sudo docker run docker/trusted-registry start)"
```

### Ubuntu 14.04 LTS upgrade

The following commands will stop the running Docker Trusted Registry, upgrade CS  engine,
and then start the Trusted Registry again:

```
    $ sudo bash -c "$(sudo docker run docker/trusted-registry stop)"
    $ sudo apt-get update && sudo apt-get dist-upgrade docker-engine-cs
    $ sudo bash -c "$(sudo docker run docker/trusted-registry start)"
```


## Upgrade Docker Trusted Registry

1. Load the Docker Trusted Registry Dashboard in your browser and click Settings in the global nav bar.
2. Click Updates in the Settings nav bar. You can see the currently installed version and a message stating that the version is either current or that there is an update available. If an update is available, the message states: System Update Available and an enabled button Update to Version X.XX.
3. Click the Update button to start the update process.

      Docker Trusted Registry pulls new Docker Trusted Registry container images from the Docker Hub. If you have not already connected to Docker Hub, Docker Trusted Registry prompts you to log in.

      The upgrade process requires a small amount of downtime to complete depending on your connection speed.

      Docker Trusted Registry:

      * Connects to the Docker Hub to pull new container images with the new version of Docker Trusted Registry.
      * Deploys those containers.
      * Shuts down the old containers.
      * Resolve any necessary links/urls.

> **Note**: If the CS engine is upgraded first, then
> the Trusted Registry can still be upgraded from the command line by running the following command. Ensure to put the correct version that you want.
>
> `sudo bash -c "$(sudo docker run docker/trusted-registry:1.1.0 upgrade 1.1.1)"`

## See also

* To configure for your environment, see the
[configuration instructions]({{< relref "configuration.md" >}}).
* To use Docker Trusted Registry, see [the User guide]({{< relref "userguide.md" >}}).
* See [installing the CS engine]({{< relref "csengineinstall.md" >}}).
* To make administrative changes, see [the Admin guide]({{< relref "adminguide.md" >}}).
* To see previous changes, see [the release notes]({{< relref "release-notes.md" >}}).
