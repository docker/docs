---
title: Install UCP offline
description: Learn how to install Docker Universal Control Plane. on a machine with
  no internet access.
keywords:
- docker, ucp, install, offline
---

The procedure to install Docker Universal Control Plane on a host is the same,
whether that host has access to the internet or not.

The only difference when installing on an offline host,
is that instead of pulling the UCP images from Docker Hub, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll install UCP.


1.  Get the UCP package.

    Use a computer with internet access to download a single package with all
    Docker Datacenter components:

    ```bash
    $ wget https://packages.docker.com/caas/ucp-2.0.0_dtr-2.1.0.tar.gz -O docker-datacenter.tar.gz
    ```

2.  Transfer the package to the offline nodes.

    Now that you have the UCP package in your machine, you can transfer it to the
    host that you want to manage with UCP. For each host:

    ```bash
    $ scp docker-datacenter.tar.gz <user>@<host>:/tmp
    ```

3. Login into the hosts where you transferred the images.

4.  Load the UCP images.

    Once the UCP package is transferred to the hosts, you can use the
    `docker load` command, to load the images from the tar archive. On each
    host, run:

    ```bash
    $ docker load < docker-datacenter.tar.gz
    ```

5.  Install Docker UCP.

    Now that the offline hosts have all the images needed to install UCP,
    you can [install Docker UCP that host](install-production.md).


## Where to go next

* [Install UCP for production](install-production.md).
* [UCP system requirements](system-requirements.md)
