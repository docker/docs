---
description: Learn how to upgrade your Docker Trusted Registry
keywords: docker, dtr, upgrade, install
title: Upgrade DTR
---

The first step in upgrading to a new minor version or patch release of DTR 2.2,
is ensuring you're running DTR 2.1. If that's not the case, start by upgrading
your installation to version 2.1, and then upgrade to 2.2.

There is no downtime when upgrading a highly-available DTR cluster. If your
DTR deployment has a single replica, schedule the upgrade to take place outside
business peak hours to ensure the impact on your business is close to none.

## Step 1. Upgrade DTR to 2.1

Make sure you're running DTR 2.1. If that's not the case, [upgrade your installation to the 2.1 version](/datacenter/dtr/2.1/guides/install/upgrade/.md).

## Step 2. Upgrade DTR



To upgrade DTR, **login with ssh** into a node that's part of the UCP cluster.
Then pull the latest version of DTR:

```none
$ docker pull {{ page.docker_image }}
```

If the node you're upgrading doesn't have access to the internet, you can
use a machine with internet connection to
[pull all the DTR images](../install/install-offline.md).

Once you have the latest images on the node, run the upgrade command:

```none
$ docker run -it --rm \
  {{ page.docker_image }} upgrade \
  --ucp-insecure-tls
```

By default the upgrade command runs in interactive mode and prompts you for
any necessary information. You can also check the
[reference documentation](../../../reference/cli/index.md) for other existing flags.

## Where to go next

* [Release notes](release-notes.md)
