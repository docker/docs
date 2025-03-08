---
description: Learn how to install Docker Trusted Registry on a machine with no internet
  access.
keywords: docker, registry, install, offline
title: Install Docker Trusted Registry offline
---

The procedure to install Docker Trusted Registry on a host is the same,
whether that host has access to the internet or not.

The only difference when installing on an offline host,
is that instead of pulling the UCP images from Docker Hub, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll install DTR.

## Versions available

{% include components/ddc_url_list.html %}

## Download the offline package

Use a computer with internet access to download a single package with all
Docker Datacenter components:

```bash
$ wget <package-url> -O docker-datacenter.tar.gz
```

Now that you have the package in your local machine, you can transfer it to
the machines where you want to install DTR.

For each machine where you want to install DTR:

1.  Copy the Docker Datacenter package to that machine.

    ```bash
    $ scp docker-datacenter.tar.gz <user>@<host>:/tmp
    ```

2.  Use ssh to login into the hosts where you transferred the package.

3.  Load the Docker Datacenter images.

    Once the package is transferred to the hosts, you can use the
    `docker load` command, to load the Docker images from the tar archive:

    ```bash
    $ docker load < docker-datacenter.tar.gz
    ```

## Install DTR

Now that the offline hosts have all the images needed to install DTR,
you can [install DTR on that host](index.md).


## Where to go next

* [DTR architecture](../architecture.md)
* [Install DTR](index.md)
