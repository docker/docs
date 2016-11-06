<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/install/upgrade/upgrade-minor/"]
title = "Upgrade from 2.0.0"
description = "Learn how to upgrade your Docker Trusted Registry to a new minor version or patch release."
keywords = ["docker, dtr, upgrade, install"]
[menu.main]
parent="menu_dtr_upgrade"
identifier="dtr_upgrade_minor"
weight=10
+++
<![end-metadata]-->

# Upgrade from 2.0.0

<!-- beta: remove after this goes to GA -->

> You can't upgrade from one beta version to another. To test the latest
> beta version, start by uninstalling the previous version.

The first step in upgrading to a new minor version or patch release of DTR 2.0,
is ensuring you're running DTR 2.0. If that's not the case, start by upgrading
your installation to version 2.0.0, and then upgrade to the latest version
available.

There is no downtime when upgrading an highly-available DTR cluster. If your
DTR deployment has a single replica, schedule the upgrade to take place outside
business peak hours to ensure the impact on your business is close to none.

## Step 1. Upgrade DTR to 2.0

Make sure you're running DTR 2.0. If that's not the case, [upgrade your
installation to the 2.0 version](upgrade-major.md).

## Step 2. Upgrade DTR

To upgrade DTR you use the `upgrade` command.

1. Download a UCP client bundle.

    Having a UCP client bundle allows you to run Docker commands on a UCP
    cluster. Download a UCP client bundle and set up your CLI client to use it.

2. Pull the latest `docker/dtr` image.

    ```bash
    $ docker pull docker/dtr:<version>
    ```

    If the node you're upgrading doesn't have access to the internet, you can
    use a machine with internet connection to
    [pull all the DTR images](../install-dtr-offline.md).

4. Run the upgrade command.

    The upgrade command upgrades all DTR replicas that are part of your cluster:

    ```bash
    $ docker run -it --rm \
      docker/dtr:<version> upgrade \
      --ucp-insecure-tls
    ```

    By default the upgrade command runs in interactive mode and prompts you for
    any necessary information. You can also check the
    [reference documentation](../reference/index.md) for other existing flags.

## Where to go next

* [Upgrade to DTR 2.0](upgrade-major.md)
* [Monitor DTR](../../monitor-troubleshoot/index.md)
