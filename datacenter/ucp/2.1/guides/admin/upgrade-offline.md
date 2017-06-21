---
description: Learn how to upgrade Docker Universal Control Plane. on a machine with
  no internet access.
keywords: docker, ucp, upgrade, offline
title: Upgrade UCP offline
---

Upgrading Universal Control Plane is the same, whether your hosts have access to
the internet or not.

The only difference when upgrading on an offline host,
is that instead of pulling the UCP images from Docker Store, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll upgrade UCP.

## Versions available

{% include components/ddc_url_list.html %}

## Download the offline package

Use a computer with internet access to download a single package with all
UCP components:

```bash
$ wget <package-url> -O docker-datacenter.tar.gz
```

Now that you have the package in your local machine, you can transfer it to
the machines where you want to upgrade UCP.

For each machine that you want to manage with UCP:

1.  Copy the offline package to that machine.

    ```bash
    $ scp docker-datacenter.tar.gz <user>@<host>:/tmp
    ```

2.  Use ssh to log in to the hosts where you transferred the package.

3.  Load the UCP images.

    Once the package is transferred to the hosts, you can use the
    `docker load` command, to load the Docker images from the tar archive:

    ```bash
    $ docker load < docker-datacenter.tar.gz
    ```

## Upgrade UCP

Now that the offline hosts have all the images needed to upgrade UCP,
you can [upgrade Docker UCP](upgrade.md).


## Where to go next

* [Upgrade UCP](upgrade.md)
* [Release Notes](release-notes.md)
