---
title: Install Docker Trusted Registry offline
description: Learn how to install Docker Trusted Registry on a machine with no internet
  access.
keywords:
- docker, registry, install, offline
---

The procedure to install Docker Trusted Registry on a node is the same,
whether that node has access to the internet or not.

The only difference when installing DTR on an offline node, is that instead
of pulling the DTR images from Docker Hub, you use a computer that is connected
to the internet to download a single package with all DTR images. Then you
copy that package to the nodes where you’ll install DTR.

1.  Get the DTR package.

    Use a computer with internet access to download a single package with all
    Docker Datacenter components:

    ```bash
    $ wget https://packages.docker.com/caas/ucp-{{ page.ucp_latest_version }}_dtr-{{ page.dtr_latest_version }}.tar.gz -O docker-datacenter.tar.gz
    ```

2.  Transfer the package to the offline nodes.

    Now that you have the DTR package in your machine, you can transfer it to the
    nodes that you want to install DTR. For each node run:

    ```bash
    $ scp docker-datacenter.tag.gz <user>@<host>:/tmp
    ```

3. Login into the nodes where you transferred the images.

4.  Load the images.

    Once the package is on the nodes where you want to install DTR, you can use
    the `docker load` command, to load the images from the .tar file. On each
    node, run:

    ```bash
    $ docker load < /tmp/docker-datacenter.tar.gz
    ```

5. Install DTR.

    Now that the offline node has all the images needed to install DTR,
    you can [install DTR that host](index.md).


## Where to go next

* [DTR architecture](../architecture.md)
* [Install DTR](index.md)
