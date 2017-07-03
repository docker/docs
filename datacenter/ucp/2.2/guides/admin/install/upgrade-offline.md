---
title: Upgrade UCP offline
description: Learn how to upgrade Docker Universal Control Plane on a machine with no internet access.
keywords: ucp, upgrade, offline
---

Upgrading Universal Control Plane is the same, whether your hosts have access to
the internet or not.

The only difference when installing on an offline host is that instead of
pulling the UCP images from Docker Hub, you use a computer that's connected
to the internet to download a single package with all the images. Then you
copy this package to the host where you upgrade UCP. 

## Versions available

{% include components/ddc_url_list.html %}

## Download the offline package

Use a computer with internet access to download a single package with all
UCP components:

```bash
$ wget <package-url> -O ucp.tar.gz
$ wget <package-url> -O dtr.tar.gz
```

Now that you have the package in your local machine, you can transfer it to
the machines where you want to upgrade UCP.

For each machine that you want to manage with UCP:

1.  Copy the offline package to the machine.

    ```bash
    $ scp ucp.tar.gz <user>@<host>:/tmp
    $ scp dtr.tar.gz <user>@<host>:/tmp
    ```

2.  Use ssh to log in to the hosts where you transferred the package.

3.  Load the UCP and DTR images.

    Once the package is transferred to the hosts, you can use the
    `docker load` command, to load the Docker images from the tar archive:

    ```bash
    $ docker load < ucp.tar.gz
    $ docker load < dtr.tar.gz
    ```

## Upgrade UCP

Now that the offline hosts have all the images needed to upgrade UCP,
you can [upgrade Docker UCP](index.md).


## Where to go next

* [Upgrade UCP](index.md)
* [Release Notes](release-notes.md)
