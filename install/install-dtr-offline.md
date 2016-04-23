<!--[metadata]>
+++
title = "Install offline"
description = "Learn how to install Docker Trusted Registry on a machine with no internet access."
keywords = ["docker, registry, install, offline"]
[menu.main]
parent="workw_dtr_install"
weight=30
+++
<![end-metadata]-->


# Install DTR offline

The procedure to install Docker Trusted Registry on a node is the same,
whether that node has access to the internet or not.

The only difference when installing DTR on an offline node, is that instead
of pulling the DTR images from Docker Hub, you use a computer that is connected
to the internet to download a single package with all DTR images. Then you
copy that package to the node where you’ll install DTR.

1. Get the DTR package.

    Use a computer with internet access to download a single package with all DTR
    images. As an example, to download UCP 2.0, run:

    ```bash
    $ wget https://packages.docker.com/dtr/2.0/dtr-2.0.0.tar
    ```

2. Transfer the package to the offline node.

    Now that you have the DTR package file, transfer it to the node where you
    want to install Docker Trusted Registry. You can use the `scp` command
    for this.

    ```bash
    $ scp ./dtr-2.0.0.tar user@dtr-node:/tmp
    ```

3. Load the DTR images.

    Once the package is on the node where you want to install DTR, you can use
    the `docker load` command, to load the images from the .tar file.

    ```bash
    $ docker load < /tmp/dtr-2.0.0.tar
    ```

4. Install DTR.

    Now that the offline node has all the images needed to install UCP,
    you can [install DTR that machine](install-dtr.md).


## Where to go next

* [DTR architecture](../architecture.md)
* [Install DTR](install-dtr.md)
