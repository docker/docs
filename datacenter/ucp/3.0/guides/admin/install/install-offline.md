---
title: Install UCP offline
description: Learn how to install Docker Universal Control Plane. on a machine with
  no internet access.
keywords: UCP, install, offline, Docker EE
---

The procedure to install Docker Universal Control Plane on a host is the same,
whether the host has access to the internet or not.

The only difference when installing on an offline host is that instead of
pulling the UCP images from Docker Hub, you use a computer that's connected
to the internet to download a single package with all the images. Then you
copy this package to the host where you install UCP. The offline installation
process works only if one of the following is true:

-  All of the swarm nodes, managers and workers alike, have internet access
   to Docker Hub, and
-  None of the nodes, managers and workers alike, have internet access to
   Docker Hub.

If the managers have access to Docker Hub while the workers don't,
installation will fail.

## Versions available

Use a computer with internet access to download the UCP package from the
following links.

{% include components/ddc_url_list_2.html product="ucp" version="2.2" %}

## Download the offline package

You can also use these links to get the UCP package from the command
line:

```bash
$ wget <ucp-package-url> -O ucp.tar.gz
```

Now that you have the package in your local machine, you can transfer it to
the machines where you want to install UCP.

For each machine that you want to manage with UCP:

1.  Copy the UCP package to the machine.

    ```bash
    $ scp ucp.tar.gz <user>@<host>
    ```

2.  Use ssh to log in to the hosts where you transferred the package.

3.  Load the UCP images.

    Once the package is transferred to the hosts, you can use the
    `docker load` command, to load the Docker images from the tar archive:

    ```bash
    $ docker load < ucp.tar.gz
    ```

Follow the same steps for the DTR binaries.

## Install UCP

Now that the offline hosts have all the images needed to install UCP,
you can [install Docker UCP on one of the manager nodes](index.md).


## Where to go next

* [Install UCP](index.md).
* [System requirements](system-requirements.md)
