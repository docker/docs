<!--[metadata]>
+++
title ="Replicate CAs for high availability"
description="Docker Universal Control plane has support for high availability. Learn how to set up your installation to ensure it tolerates failures."
keywords= ["replica, controller, availability, high, ucp"]
[menu.main]
parent="mn_ucp_high_availability"
weight=10
+++
<![end-metadata]-->

# Replicate CAs for high availability

Internally, UCP uses two Certificate Authorities:

* `ucp-cluster-root-ca` signs the certificates that are used by new nodes
joining the cluster. This CA also generates the certificates for admin client
certificate bundles.
* `ucp-client-root-ca` signs non-admin client certificate bundles.

When configuring UCP for high-availability, you need to ensure the CAs running
on each UCP controller node are interchangeable. This is done by using the same
root certificates and keys for every CA on the cluster.


## Replicating the CAs across controller nodes

There are two ways to ensure the same root certificates and keys are used on
every controller node:

* When joining a controller node:

    When you install the first controller node, root certificates and keys are
    automatically generated.
    After installing the first controller node, you can backup the certificates
    and keys used by that controller using.


    Then, when joining new controller nodes to the cluster, you use the use
    the backup you've created to make the CAs in both nodes use the same root
    certificates and keys.

* After joining all controller nodes:

    Alternatively, you can replicate root certificates and keys through all
    controller nodes after you've joined them to the cluster.

    After the installation, you backup the root certificates and keys of the
    controller that had the valid CA material and restore all other controllers
    with that backup archive, one by one.


## Backup the certificates and keys

To create a backup of the root certificates and keys used by the CAs, use the
backup command. Notice that this command stops the UCP containers, so you
should use it outside business peak hours.

```bash
$ docker run --rm -i --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp backup --root-ca-only --passphrase "secret" > /tmp/backup.tar
```

Where:

* root-ca-only, specifies to only backup the CA certificates and keys,
* passphrase, encrypts the backup with a given passphrase

## Restore the certificate and keys

Once you have a backup archive of the certificates and keys used by the CAs
of a controller node, you can use it to make CAs in other controller nodes
use the same certificate and private key.

You can use the restore command for this:

```bash
$ docker run --rm -i --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp restore --root-ca-only --passphrase "secret" < /tmp/backup.tar
```

Where:

* root-ca-only, specifies to only restore the volumes used by the CAs,
* passphrase, specifies the passphrase to decrypt the backup archive.

## Where to go next

* [Set up high availability](set-up-high-availability.md)
* [Backups and disaster recovery](backups-and-disaster-recovery.md)
