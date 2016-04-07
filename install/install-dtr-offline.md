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

The procedure to install Docker Trusted Registry on a host is the same,
whether that host has access to the internet or not.

The only difference when installing on an offline host,
is that instead of pulling the DTR images from Docker Hub, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll install DTR.

## Step 1. Get the DTR package

Use a computer with internet access to download a single package with all DTR
images. As an example, to download UCP 2.0, run:

```bash
$ wget https://packages.docker.com/dtr/2.0/dtr-2.0.0.tar
```

## Step 2. Copy the package
Now that you have the DTR package file, transfer it to the host where you want
to install Docker Trusted Registry. You can use the `scp` command for this.

```bash
$ scp ./dtr-2.0.0.tar user@dtr-host:/tmp
```

## Step 3. Load the DTR images

Once the package is on the host where you want to install DTR, you can use
the `docker load` command, to load the images from the .tar file.

```bash
$ docker load < /tmp/dtr-2.0.0.tar
```

## Step 4. Install DTR

Now that the offline host has all the images needed to install UCP,
you can [install DTR that machine](install-dtr.md).
