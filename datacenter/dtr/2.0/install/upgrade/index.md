---
description: Learn how to upgrade your Docker Trusted Registry to a new minor version or patch release.
keywords: docker, dtr, upgrade, install
redirect_from:
- /docker-trusted-registry/install/upgrade/upgrade-minor/
title: Upgrade from 2.0.0
---

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

To upgrade DTR you use the `docker/dtr upgrade` command.

1.  Download a UCP client bundle.

    Having a UCP client bundle allows you to run Docker commands on a UCP
    cluster. Download a UCP client bundle and set up your CLI client to use it.

2.  Find the DTR replica Id.

    When you upgrade your installation, you need to specify the Id of a replica
    that is part of the cluster. If you have a highly-available installation,
    you can provide the Id of any replica.

    You can find the DTR replica Ids on the **Applications** page of Docker
    Universal Control Plane.

3.  Pull the latest docker/dtr image.

    ```bash
    $ docker pull docker/dtr
    ```

4.  Run the upgrade command.

    The upgrade command upgrades all DTR replicas that are part of your cluster:

    ```bash
    # Get the certificates used by UCP
    $ curl -k https://$UCP_URL/ca > ucp-ca.pem

    $ docker run \
      -it \
      --rm \
      docker/dtr upgrade \
        --ucp-url $UCP_URL \
        --existing-replica-id $DTR_REPLICA_ID \
        --ucp-username $USER \
        --ucp-password $PASSWORD \
        --ucp-ca "$(cat ucp-ca.pem)"
    ```

## Where to go next

* [Upgrade to DTR 2.0](upgrade-major.md)
* [Monitor DTR](../../monitor-troubleshoot/index.md)