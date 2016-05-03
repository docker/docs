<!--[metadata]>
+++
title = "Install offline"
description = "Learn how to install Docker Universal Control Plane. on a machine with no internet access."
keywords = ["docker, ucp, install, offline"]
[menu.main]
parent="mn_ucp_installation"
identifier="ucp_install_offline"
weight=30
+++
<![end-metadata]-->

# Install UCP offline

The procedure to install Docker Universal Control Plane on a host is the same,
whether that host has access to the internet or not.

The only difference when installing on an offline host,
is that instead of pulling the UCP images from Docker Hub, you use a
computer that is connected to the internet to download a single package with
all the images. Then you copy that package to the host where you’ll install UCP.


## Step 1. Download UCP

Use a computer with internet access to download a single package with all UCP
images. As an example, to download UCP 1.1, run:

```bash
$ wget https://packages.docker.com/ucp/1.0/ucp-1.1.tar
```

## Step 2. Transfer UCP

Now that you have the UCP package in your machine, you can transfer it to the
host where you'll be installing Docker UCP. You can use the Secure Copy command
for this:

```bash
$ scp ./ucp-1.1.tar user@ucp-host:/tmp
```

## Step 3. Load the UCP images

Once the UCP package is transferred to the host, you can use the `docker load`
command, to load the images from the tar archive. On the host were you are
going to install UCP, run:

```bash
$ docker load < /tmp/ucp-1.1.tar
```

## Step 4. Install Docker UCP

Now that the offline host has all the images needed to install UCP,
you can [install Docker UCP that machine](install-production.md).


## Where to go next

* [Install UCP for production](install-production.md).
* [UCP system requirements](system-requirements.md)
