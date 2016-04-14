<!--[metadata]>
+++
title = "Upgrade"
description = "Learn how to upgrade your Docker Trusted Registry."
keywords = ["docker, dtr, upgrade, install"]
[menu.main]
parent="workw_dtr_install"
identifier="dtr_upgrade"
weight=40
+++
<![end-metadata]-->


# Upgrade the Trusted Registry and the CS Engine

This document describes the steps to upgrade Docker Trusted Registry and the
commercially supported Engine (CS Engine). When you first install, the general
order is to install the CS Engine, then install the Trusted Registry. However,
when you upgrade, you reverse that order. Ensure when upgrading the Trusted
Registry, that you also upgrade to the latest CS Engine.

## Upgrade Docker Trusted Registry

Periodic upgrades to the Trusted Registry trigger a notification to appear in
your Admin dashboard if you have enabled Upgrade checking. This is located in
the General > Settings section of the Trusted Registry Admin dashboard. To
perform this upgrade, you should schedule it during your downtime and allow
about 15 minutes.

To upgrade, perform the following steps:

1. Load the Trusted Registry Dashboard in your browser and navigate to
**Settings > Updates**.

2. Click **Updates** in the Settings navigation bar.

    You can see the currently installed version and a message stating that
    the version is either current or an update is available. If an update
    is available, the message states:

    "System Update Available and an enabled button displays Update
    to version X.X.X.""

3. Click Update to start the update process.

    The process may take longer than what the message indicates.
    To check the status of the install, SSH into the Trusted Registry
    host through a command line:

    ```
    $ sudo docker logs -f $(sudo docker ps -a --no-trunc | grep 'manager execute-upgrade' | head -n1 | awk '{print $1}')
    ```

4. Refresh your screen to see the latest changes.

    The Dashboard displays a message that the upgrade successfully
    completed and that you need to upgrade to the latest CS Engine.

## Upgrade Docker Trusted Registry offline

To upgrade the Trusted Registry offline, perform the following steps:

1. Since you are retrieving a large file, use the `wget` command in your
command line to get the Trusted Registry files. The following
command is an example getting DTR 1.4.3. Ensure to get your correct version.

    ```
    $ wget https://packages.docker.com/dtr/1.4/dtr-1.4.3.tar
    ```

2. After downloading, move the `tar` file to the offline machine you
want to install the Trusted Registry.

3. On that machine, verify that the CS Engine is installed.
If it is not, see the [CS Engine install directions](../cs-engine/install.md).

    ```bash
    $ docker --version
    ```

    > **Note:** To remain compliant with your Docker Trusted Registry support
    > agreement, you **must** use the current version of commercially supported
    > Docker Engine. Running the open source version of Engine is **not**
    > supported.

5. Open a terminal window on that machine and load the `tar` file using the
following command. Again, ensure you get the correct version.

    ```bash
    $ sudo docker load < dtr-1.4.3.tar
```

6. Upgrade the Trusted Registry with the following command:

    ```bash
    $ sudo bash -c "$(docker run docker/trusted-registry upgrade latest)"
    ```

    > **Note**: sudo is needed for `docker/trusted-registry` commands to
    > ensure that the Bash script is run with full access to the Docker host.

### What is updated in the Trusted Registry?

The Trusted Registry pulls new container images from Docker Hub.
Then it deploys those containers. Finally, it stops and removes the
old containers.

If the CS Engine is upgraded first, then the Trusted Registry can still be
upgraded from a command line by running the following command. Ensure to put the
correct version that you want.

```bash
$ sudo bash -c "$(sudo docker run docker/trusted-registry:1.3.3 upgrade 1.4.3)"
```

## See also

* [Install DTR](install-dtr.md)
* [Install DTR offline](install-dtr-offline.md)
