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

Use a computer with internet access to download the UCP package from the
following links.

{% include components/ddc_url_list_2.html product="ucp" version="3.0" %}

## Download the offline package

You can also use these links to get the UCP package from the command
line:

```bash
$ wget <ucp-package-url> -O ucp.tar.gz
```

Now that you have the package in your local machine, you can transfer it to
the machines where you want to upgrade UCP.

For each machine that you want to manage with UCP:

1.  Copy the offline package to the machine.

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

## Upgrade UCP

Now that the offline hosts have all the images needed to upgrade UCP,
you can [upgrade Docker UCP](upgrade.md).

## Where to go next

- [UCP release notes](../../release-notes.md)
- [Upgrade UCP](upgrade.md)
