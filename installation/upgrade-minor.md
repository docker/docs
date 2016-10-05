<!--[metadata]>
+++
aliases = [ "/ucp/upgrade-ucp/",
            "/ucp/installation/upgrade/"]
title = "Upgrade to UCP 1.1"
description = "Learn how to upgrade Docker Universal Control Plane with minimal impact to your users."
keywords = ["Docker, UCP, upgrade, update"]
[menu.main]
parent="mn_ucp_installation"
identifier="ucp_upgrade_minor"
weight=60
+++
<![end-metadata]-->

# Upgrade to UCP 1.1

This page guides you in upgrading Docker Universal Control Plane (UCP) to
version 1.1.

Before upgrading to a new version of UCP, check the
[release notes](../release-notes/index.md) for the version you are upgrading to.
There you'll find information about the new features, breaking changes, and
other relevant information for upgrading to a particular version.

## Plan the upgrade

You should plan for the upgrade to take place outside business hours. Even
though there is no expected downtime, upgrading outside business hours ensures
there is no impact to your users.

Before starting the upgrade, make sure your cluster is healthy. If a problem
occurs, that will make it easier to find and troubleshoot any problems.

Also, don't make changes to the cluster configurations while you're upgrading
the cluster. That can lead to misconfigurations that are difficult to
troubleshoot.

## Avoid downtime

If you have multiple controller nodes and are using a load balancer to
distribute requests across nodes, drain the load balancer connections to that
node before upgrading it. While connections are being drained, all new
requests are sent to the other nodes and existing connections to that node are
given some time to complete.

This way users can continue using UCP, while that controller node is being
upgraded.

## Upgrade your cluster

Start by upgrading the controller nodes one by one. Once all controller nodes
are upgraded, upgrade existing worker nodes. To upgrade a node:

1. Drain the node from the load balancer.

    If this is a controller node that is being load-balanced, configure the
    load balancer to drain existing connections to this node.

2. Log into the node using ssh.

3. Pull the docker/ucp image for the version you want to upgrade to.

    ```bash
    $ docker pull docker/ucp:1.1.2
    ```

4. Upgrade the node.

    ```bash
    $ docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:1.1.2 \
      upgrade --interactive
    ```

    This runs the upgrade command in interactive mode, so that you are prompted
    for any necessary configuration values.

5. Delete your browser cache.

    When upgrading controller nodes, make sure you delete the browser cache
    before accessing the UCP web UI. Your browser can continue serving cached
    pages from the older UCP version, which can cause errors.

6. Make sure the node is healthy.

    In your browser, navigate to the **UCP web UI**. In the **Nodes** page
    confirm that the node is running, and the cluster is healthy.

    If this is a controller node, you can add the node back in the load
    balancing pool.

7. Upgrade other controller nodes using the same procedure.

8. Upgrade other worker nodes using the same procedure.

## Where to go next

* [UCP release notes](../release-notes/index.md)
