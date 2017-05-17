---
title: Deploy DDC on Linux servers
description: Learn how to get a trial license and install Docker Datacenter
keywords: docker, datacenter, install, orchestration
redirect_from:
- /datacenter/try/
---

The best way to try Docker Datacenter for yourself is to get the [30-day
trial available at the Docker Store](https://store.docker.com/bundles/docker-datacenter).

Once you get your trial license, you can install Docker Datacenter in your
Linux servers. Make sure all the hosts you want to manage in with Docker
Datacenter have a minimum of:

* Linux kernel version 3.10 or higher
* CS Docker Engine version 1.13.0 or higher
* 4.00 GB of RAM
* 3.00 GB of available disk space

Also make sure the hosts are running one of these operating systems:

* CentOS 7.1 or 7.2
* Red Hat Enterprise Linux 7.0, 7.1, 7.2, or 7.3
* Ubuntu 14.04 LTS or 16.04 LTS
* SUSE Linux Enterprise 12

[Learn more about the Docker Datacenter system requirements](../ucp/2.1/guides/admin/install/system-requirements.md)


### Step 1: Install Docker EE Container Engine

Install the Docker Enterprise Edition container engine on all hosts you want to manage by following [these instructions](/engine/installation/).

### Step 2: Install Universal Control Plane

Docker Universal Control Plane (UCP) allows managing from a centralized place
your images, applications, networks, and other computing resources.

Use ssh to log in to the host where you want to install UCP and run:

```none
docker run --rm -it --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_latest_image }} install \
  --host-address <node-ip-address> \
  --interactive
```

This runs the install command in interactive mode, so that you're prompted
for any necessary configuration values.

[Learn more about the UCP installation](../ucp/2.1/guides/admin/install/index.md).


### Step 3: License your installation

Now that UCP is installed, you need to license it. In your browser, navigate
to the UCP web UI, log in with your administrator credentials and upload your
license.

![](../images/try-ddc-1.png){: .with-border}

[Get a free trial license if you don't have one](https://store.docker.com/bundles/docker-datacenter).

### Step 4: Join more nodes to UCP

Join more nodes so that you can manage them from UCP.
Go to the **UCP web UI**, navigate to the **Resources** page, and go to
the **Nodes** section.

![](../images/try-ddc-2.png){: .with-border}

Click the **Add Node button** to add a new node.

![](../images/try-ddc-3.png){: .with-border}


Check the 'Add node as a manager' option to join the node as a manager
to provide replication and make UCP highly available. For a highly available
installation, make sure you have 3, 5, or 7 manager nodes.

Copy the command to your clipboard, and run it on every node that you want
to be managed by UCP. After you run the command in the node, the node
will show up in the UP web UI.

### Step 5: Install Docker Trusted Registry

Docker Trusted Registry (DTR) is a private image registry so that you can
manage who has access to your Docker images. DTR needs to be installed on
a node that is being managed by UCP.

Use ssh to log in to the host where you already installed UCP, and run:

```none
docker run -it --rm \
  {{ page.dtr_latest_image }} install \
  --ucp-node <node-hostname> \
  --ucp-insecure-tls
```

Where the `--ucp-node` is the hostname of the UCP node where you want to deploy
DTR. `--ucp-insecure-tls` tells the installer to trust the certificates used
by UCP.

## Where to go next

* [Scale your UCP cluster](../ucp/2.1/guides/admin/configure/scale-your-cluster.md)
* [Deploy an application](../ucp/2.1/guides/user/services/index.md)
