---
description: Learn how to upgrade Docker Universal Control Plane with minimal impact
  to your users.
keywords: Docker, UCP, upgrade, update
title: Upgrade UCP
---

This page guides you on upgrading Docker Universal Control Plane (UCP).

Before upgrading to a new version of UCP, check the
[release notes](../release_notes.md) for the version you are upgrading to,
including information about new features, breaking changes, and
other relevant information for upgrading to a particular version.

Before starting an upgrade, make sure your cluster is healthy. If a problem
occurs, that will make it easier to find and troubleshoot any problems.
Also, don't change any cluster configuration during an upgrade. That can lead
to misconfigurations that are difficult to troubleshoot.

## The UCP upgrade command

To upgrade a UCP installation, you run the `docker/ucp upgrade` command on
each node of the cluster.

When you run the upgrade command, it:

1. Pulls the images of the new UCP version from Docker Hub.
2. Checks if it is possible to upgrade directly to the new version.

    Depending on the version you have installed, it might not be possible
    to upgrade directly to the latest version. In that case, you need
    to upgrade to intermediate versions before upgrading to the latest version.

    Check the [release notes](../release_notes.md) to see if its possible to
    upgrade directly or not.

3. Stops and removes the old UCP containers.

    This doesn't affect other running containers. Also, existing cluster
    configurations are not affected, since they are persisted in volumes.

4. Deploys the new UCP containers to the node.


## The upgrade procedure

The upgrade procedure depends on whether your cluster is set up for
high-availability or not.
A cluster that is not set for high-availability, has only one controller node,
while a cluster that supports high-availability has multiple controller nodes.

To check the number of controller nodes in your cluster, navigate to the **UCP
web UI**, and check the **Nodes** page.

![Cluster replicas](../images/multiple-replicas.png)

In this example we have 3 controller nodes set up, which means that this
cluster is set up for high-availability.

After finding the number of controller nodes in your cluster, jump to the
upgrade instructions that apply to you:

* [My cluster is not set for high-availability](upgrade.md#my-cluster-is-not-set-for-high-availability),
if your cluster has only one controller node.
* [My cluster is set for high-availability](upgrade.md#my-cluster-is-set-for-high-availability),
if your cluster has multiple controller nodes.


### My cluster is not set for high-availability

If your cluster is not set up for high-availability (does not have
replica nodes):

1. Make sure your cluster is healthy before starting the upgrade.

    Login into the **UCP UI** and navigate to the **Nodes** page. Make sure
    all nodes are listed and healthy.

2. Log into the controller node using ssh.

3.  Pull the docker/ucp image for the version you want to upgrade to.

    ```bash
    $ docker pull docker/ucp:$UCP_VERSION
    ```

4.  Upgrade the controller node.

    ```none
    $ docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:$UCP_VERSION upgrade -i
    ```

5. Delete your browser cache.

    Your browser can continue serving cached pages from the older version of
    UCP. To avoid errors on the UCP UI, ensure you delete the browser cache.

6. Make sure the controller node is healthy.

    In your browser, navigate to the **UCP web app**, to make
    sure the controller is running. Navigate to the **Nodes** page,
    to make sure the node is healthy.

7. Upgrade all other nodes in the cluster using steps 2-5.

### My cluster is set for high-availability

If your cluster is set up for high-availability (has several controller nodes):

1. Schedule the upgrade to take place outside business hours.

    During an upgrade, all user containers will continue running normally. But
    upgrading outside business hours ensures the impact on your business
    is close to none.

2. Make sure your cluster is healthy before starting the upgrade.

    Login into the **UCP UI** and navigate to the **Nodes** page. Make sure all
    nodes are listed and healthy.

3. Ensure no administrator user makes configuration changes during the upgrade.

    The UCP cluster uses an internal key-value store to save configuration
    settings, like the method used for authenticating users. For
    high-availability that key-value store is replicated across the
    controller and replica nodes.

    During an upgrade, new values are stored on the key-value store. If at the
    same time an administrator makes configuration changes, some nodes might
    use that configuration, while others might not.

    To avoid misconfigurations, ensure no administrator changes UCP
    configurations during the upgrade.

4. Block user access to the controller node.

    This can be done by setting the load balancer to drain existing connections
    to the controller node.

    While connections are being drained, all new user requests are sent to the
    replica nodes, and existing connections to the controller are given some
    time to complete.

    This way users can continue using UCP, while the controller node is
    being upgraded.

5. Log into the controller node using ssh.

6.  Pull the docker/ucp image for the version you want to upgrade to.

    ```bash
    $ docker pull docker/ucp:$UCP_VERSION
    ```

7.  Upgrade the controller node.

    ```none
    $ docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:$UCP_VERSION upgrade -i
    ```

8. Delete your browser cache.

    Your browser can continue serving cached pages from the older version of
    UCP. To avoid errors on the UCP UI, ensure you delete the browser cache.

9. Make sure the controller node is healthy.

    In your browser, navigate to the **UCP web app**. In the **Nodes** page
    confirm that the controller is running, and the cluster is healthy.

10. Add the controller node back to the load balancing pool.

11. Upgrade other controller nodes one at a time, using steps 4-9.

12. Upgrade other nodes in the cluster, using steps 5-9.

## Where to go next

* [UCP release notes](../release_notes.md)