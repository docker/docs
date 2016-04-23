<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/install/upgrade/"]
title = "Upgrade to DTR 2.0"
description = "Learn how to upgrade your Docker Trusted Registry to the latest major release."
keywords = ["docker, dtr, upgrade, install"]
[menu.main]
parent="workw_dtr_install"
identifier="dtr_upgrade_major"
weight=40
+++
<![end-metadata]-->

# Upgrade to DTR 2.0

The first step in upgrading to Docker Trusted Registry (DTR) 2.0, is ensuring
you are running DTR 1.4.3. If that's not the case, start by upgrading your
installation to version 1.4.3, and then upgrade to DTR 2.0.

To upgrade from DTR 1.4.3 to 2.0 you first need to do a fresh installation of
DTR 2.0. Then you migrate the data from your DTR 1.4.3 installation to the 2.0
installation. Finally, you decommission your 1.4.3 by uninstalling it.

## Step 1. Install DTR 2.0

The first step in upgrading to DTR 2.0 is doing a fresh installation of DTR 2.0.
This can be done on the same node where DTR 1.4.3 is already running or on a
new node.

If you decide to install the new DTR on the same node, you'll need
to install it on a port other than 443, since DTR 1.4.3 is already using it.
Use these instructions to install DTR 2.0:

* [Release notes](../../release-notes/release-notes.md)
* [System requirements](../system-requirements.md)
* [DTR architecture](../../architecture.md)
* [Install DTR 2.0](../install-dtr.md)


## Step 2. Migrate metadata

Once you have your DTR 1.4.3 and the new DTR 2.0 running, you can migrate
configurations, accounts, and repository metadata from one installation to
another.

For this, you can use the `docker/dtr migrate` command. This command
migrates configurations, accounts, and repository metadata. It doesn't migrate
the images that are on the storage backend used by DTR 1.4.3.

To find what options are available on the migrate command, check the reference
documentation, or run:

```bash
$ docker run --rm -it docker/dtr migrate --help
```

To start the migration, on the host running DTR 1.4.3, run:

```bash
# Get the certificates used by UCP
$ curl https://$UCP_HOST/ca > ucpca.crt

# Get the certificates used by DTR 2.0
$ docker run -it --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/dtr dump-certs \
  --host $UCP_HOST --ucp-ca "$(cat ucpca.crt)" \
  --pod-id $DTR_POD_ID > dtrca.crt

# Migrate configurations, accounts, and repository metadata
docker run -it --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/dtr migrate \
  --host $UCP_HOST --ucp-ca "$(cat ucpca.crt)" \
  --dtr-host https://$DTR_HOST --dtr-ca "$(cat dtrca.crt)" \
  --pod-id $DTR_POD_ID
```

## Step 3. Test your installation

After the migration finishes, test your DTR 2.0 installation to make sure it is
properly configured.
In your browser navigate to the DTR **Settings page**, and check that DTR 2.0:

* Is correctly licensed,
* Has the correct domain name configured,
* The storage backend is correctly configured,
* User authentication is correctly configured.

You should also validate that you can now push and pull images to DTR 2.0.

## Step 4. Decommission DTR 1.4.3

Once you've fully tested your new installation, you can uninstall DTR 1.4.3.

<!-- TODO: include instructions on how to uninstall -->

## Where to go next

* [Install DTR offline](../install-dtr-offline.md)
* [Monitor DTR](../../monitor-troubleshoot/monitor.md)
