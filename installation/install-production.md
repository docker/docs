<!--[metadata]>
+++
aliases = [ "/ucp/production-install/"]
title = "Install UCP for production"
description = "Learn how to install Docker Universal Control Plane on production"
keywords = ["Universal Control Plane, UCP, install"]
[menu.main]
parent="mn_ucp_installation"
identifier="ucp_install_production"
weight=20
+++
<![end-metadata]-->

# Install UCP for production

Docker Universal Control Plane (UCP) is a containerized application that can be
installed on-premise or on a cloud infrastructure.

If you're installing Docker Datacenter on Azure, [follow this guide](https://success.docker.com/?cid=ddc-on-azure).

## Step 1: Validate the system requirements

The first step in installing UCP, is ensuring your
infrastructure has all the [requirements UCP needs to run](system-requirements.md).


## Step 2: Install CS Docker on all nodes

UCP is a containerized application that requires CS Docker Engine 1.12.0 or
above to run.

So on each host that you want to be part of the UCP cluster,
[install CS Docker Engine 1.12.0 or above](https://docs.docker.com/docker-trusted-registry/cs-engine/install/).

## Step 3: Customize named volumes

Skip this step if you want to use the defaults provided by UCP.

Docker UCP uses named volumes to persist data. If you want
to customize the volume drivers used to manage these volumes, you can
create the volumes before installing UCP. When you install UCP, the installer
will notice that the volumes already exist, and will start using them.
[Learn about the named volumes used by UCP](../architecture.md).

If these volumes don't exist, they'll be automatically created when installing
UCP.

## Step 4: Install UCP

To install UCP you use the `docker/ucp` image, which has commands to install and
manage UCP.

To install UCP:

1. Use ssh to log in into the host where you want to install UCP.

2. Run the following command:

    ```bash
    $ docker run --rm -it --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp install \
      --host-address <node-ip-address> \
      --interactive
    ```

    This runs the install command in interactive mode, so that you're
    prompted for any necessary configuration values.
    To find what other options are available in the install command, check the
    [reference documentation](../reference/install.md).

## Step 5: License your installation

Now that UCP is installed, you need to license it. In your browser, navigate
to the UCP web UI and upload your license.

![](../images/install-production-1.png)

If you don't have a license yet, [learn how to get a free trial license](license.md).

## Step 6: Join manager nodes

Skip this step if you don't want your cluster to be highly available.

To make your UCP cluster fault-tolerant and highly available, you
can join more manager nodes to your cluster. Manager nodes are the nodes in the
cluster that perform the orchestration and cluster management tasks, and
dispatch tasks for worker nodes to execute.
[Learn more about high-availability](../high-availability/set-up-high-availability.md).

To add join manager nodes to the cluster, go to the **UCP web UI**, navigate to
the **Resources** page, and go to the **Nodes** section.

![](../images/install-production-2.png)

Click the **Add Node button** to add a new node.

![](../images/install-production-3.png)

Check the 'Add node as a manager' option to make the node a manager. Also set
the 'Use a custom listen address' option to specify the IP of the host that
you'll be joining to the cluster.

For each node manager node that you want to join to the cluster, login into the
node using ssh, and run the join command that is displayed on UCP.

![](../images/install-production-4.png)

After you run the join command in the node, the node starts being displayed
in UCP.

## Step 7: Join worker nodes

Skip this step if you don't want to add more nodes to run and scale your apps.

To add more computational resources to the cluster, you can join worker nodes.
These nodes execute tasks assigned to them by the manager nodes. For this,
use the same steps used to join manager nodes, but don't check the
'Add node as a manager' option.

## Step 8. Download a client certificate bundle

To validate that your cluster is correctly configured, you should try accessing
the cluster with the Docker CLI client. For this, you'll need to get a client
certificate bundle.
[Learn more about user bundles](../access-ucp/cli-based-access.md).

## Where to go next

* [Use externally-signed certificates](../configuration/use-externally-signed-certs.md)
* [Integrate with LDAP](../configuration/ldap-integration.md)
