---
description: Storage configuration for Docker Trusted Registry
keywords: docker, documentation, about, technology, understanding, configuration,
  storage, storage drivers, Azure, S3, Swift, enterprise, hub, registry
title: Configure DTR image storage
---

After installing Docker Trusted Registry, one of your first tasks is to
designate and configure the Trusted Registry storage backend.  This document
provides the following:

* Information describing your storage backend options.
* Configuration steps using either the Trusted Registry UI or a YAML file.

The default storage backend, `filesystem`, stores and serves images from the
*local* filesystem.  In a HA setup this fails, as each node can only access its
own files.

DTR allows you to confiugure your image storage via distributed stores, such as
Amazon S3, NFS, or Google Cloud Storage. This flexibility to configure to a
different storage backend allows you to:

* Scale your Trusted Registry
* Leverage storage redundancy
* Store your images anywhere in the cloud
* Take advantage of other features that are critical to your organization

At first, you might have explored Docker Trusted Registry and Docker Engine by
installing them on your system in order to familiarize yourself with them.
However, for various reasons such as deployment purposes or continuous
integration, it makes sense to think about your long term organization’s needs
when selecting a storage backend. The Trusted Registry natively supports TLS and
basic authentication.

## Understand the Trusted Registry storage backend

Your Trusted Registry data (images etc.) are stored using the configured
**storage driver** within DTR's settings.  This defaults to the local
filesystem which uses your OS' posix operations to store and serve images.

Additionally, the Trusted Registry supports these cloud-based storage drivers:

* Amazon Simple Storage Solution **S3** (and S3-compatible servers)
* OpenStack **Swift**
* Microsoft **Azure** Blob Storage
* **Google Cloud** Storage

### Filesystem

The `filesystem` driver operates on the host's local filesystem.  In HA
environments this needs to be shared via NFS, otherwise each node in your setup
will only be able to see their own local data.  For more information on
configuring NFS [see the NFS docs](/datacenter/dtr/2.2/guides/admin/
configure/external-storage/nfs/).

By default, docker creates a volume named `dtr-registry-${replica-id}` which is
used to host your data.  You can supply a different volume name or directory
when installing or reconfiguring docker to change where DTR stores your data
locally.

When using your local filesystem (or NFS) to serve images ensure there is enough
available space, otherwise pushes will begin to fail.

You can see the total space used locally by running `du -hs "path-to-volume"`.
The path to the docker volume can be found by running `docker volume ls` to list
volumes and `docker volume inspect dtr-registry-$replicaId` to show the path.

### Amazon S3

DTR supports AWS S3 plus other file servers that are S3 compatible, such as
Minio.  For more information on configuring S3 or a compatible backend see the
[S3 configuration guide](
/datacenter/dtr/2.2/guides/admin/configure/external-storage/s3/).


### OpenStack Swift

OpenStack Swift, also known as OpenStack Object Storage, is an open source
object storage system that is licensed under the Apache 2.0 license. Refer to [Swift documentation](http://docs.openstack.org/developer/swift/) to get started.


### Microsoft Azure

This storage backend uses Microsoft’s Azure Blob storage. Data is stored within
a paid Windows Azure storage account. Refer to Microsoft's Azure
[documentation](https://azure.microsoft.com/en-us/services/storage/) which
explains how to set up your Storage account.

## Configure your Trusted Registry storage backend

Once you select your driver, you need to configure it through the UI or use a
YAML file (which is discussed further in this document.)

1. From the main Trusted Registry screen, navigate to Settings > Storage.
2. Under Storage Backend, use the drop down menu to select your storage. The screen refreshes to reflect your option.
3. Enter your configuration settings. If you're not sure what a particular parameter does, then find your driver from the following headings so that you can see a detailed explanation.
4. Click Save. The Trusted Registry restarts so that your changes take effect.

>**Note**: Changing your storage backend requires you to restart the Trusted Registry.

See the [Registry configuration](/registry/configuration.md)
documentation for the full options specific to each driver. Storage drivers can
be customized through the [Docker Registry storage driver
API](/registry/storage-drivers/index.md#storage-driver-api).


### Filesystem settings

The [filesystem storage backend](/registry/configuration.md#filesystem)
has only one setting, the "Storage directory".

### S3 settings

If you select the [S3 storage backend](/registry/configuration.md#s3), then you
need to set  "AWS region", "Bucket name", "Access Key", and "Secret Key".

### Azure settings

Set the "Account name", "Account key", "Container", and "Realm" on the [Azure storage backend](/registry/configuration.md#azure) page.

### Openstack Swift settings

View the [Openstack Swift settings](/registry/configuration.md#openstack-swift)
documentation so that you can set up your storage settings: authurl, username,
password, container, tenant, tenantid, domain, domainid, insecureskipverify,
region, chunksize, and prefix.

## Configure using a YAML file

If the previous quick setup options are not sufficient to configure your
Registry options, you can upload a YAML file. The schema of this file is
identical to that used by the [Registry](/registry/configuration.md).

There are several benefits to using a YAML file as it can provide an
additional level of granularity in defining your storage backend. Advantages
include:

* Overriding specific configuration options.
* Overriding the entire configuration file.
* Selecting from the entire list of configuration options.

**To configure**:

1. Navigate to the Trusted Registry UI > Settings > Storage.
2. Select Download to get the text based file. It contains a minimum amount
of information and you're going to need additional data based on your driver and
business requirements.
3. Go [here](/registry/configuration.md#list-of-configuration-options") to see the open source YAML file. Copy the sections you need and paste into your `storage.yml` file. Note that some settings may contradict others, so
ensure your choices make sense.
4. Save the YAML file and return to the UI.
5. On the Storage screen, upload the file, review your changes, and click Save.

## Where to go next

* [Set up high availability](../set-up-high-availability.md)
