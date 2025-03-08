---
description: Storage configuration for Docker Trusted Registry
keywords: docker, documentation, about, technology, understanding, configuration,
  storage, storage drivers, Azure, S3, Swift, enterprise, hub, registry
title: Configure where images are stored
---

After installing Docker Trusted Registry, one of your first tasks is to
designate and configure the Trusted Registry storage backend. This document provides the following:

* Information describing your storage backend options.
* Configuration steps using either the Trusted Registry UI or a YAML file.

While there is a default storage backend, `filesystem`, the Trusted Registry offers other options that are cloud-based. This flexibility to configure to a different storage backend allows you to:

* Scale your Trusted Registry
* Leverage storage redundancy
* Store your images anywhere in the cloud
* Take advantage of other features that are critical to your organization

At first, you might have explored Docker Trusted Registry and Docker Engine by
 [installing](../install/index.md)
them on your system to familiarize yourself with them.
However, for various reasons such as deployment purposes or continuous
integration, it makes sense to think about your long term organization’s needs
when selecting a storage backend. The Trusted Registry natively supports TLS and
basic authentication. It also ships with a [notification system](/registry/notifications.md), calling webhooks in response to activity. The notification system also provides both extensive logging and
reporting, which is useful for organizations that want to collect metrics.

## Understand the Trusted Registry storage backend

By default, your Trusted Registry data resides as a data volume on the host
`filesystem`. This is where your repositories and images are stored. This
storage driver is the local posix `filesystem` and is configured to use a
directory tree in the local filesystem. It's suitable for development or small
deployments. The `filesystem` can be located on the same computer as the Trusted Registry, or on a separate system.

Additionally, the Trusted Registry supports these cloud-based storage drivers:

* Amazon Simple Storage Solution **S3**
* OpenStack **Swift**
* Microsoft **Azure** Blob Storage

<!--* **Rados**:  A driver storing objects in a Ceph Object Storage pool. -->

### Filesystem

If you select `filesystem`, then the Trusted Registry uses the local disk to
store registry files. This backend has a single, required `rootdirectory`
parameter which specifies a subdirectory of `/var/local/dtr/imagestorage` in
which all registry files are stored. The default value of `/local` means the
files are stored in `/var/local/dtr/image-storage/local`.

The Trusted Registry stores all its data at this location, so ensure there is
adequate space available. To do so, you can run the following commands:

* To analyze the disk usage: `docker exec -it <container_name> bash` then run `df -h`.
* To see the file size of your containers, use the `-s` argument of `docker ps -s`.

### Amazon S3

S3 stores data as objects within “buckets” where you read, write, and delete
objects in that container. It too, has a `rootdirectory` parameter. If you select this option, there will be some tasks that you need to first perform [on AWS](https://aws.amazon.com/s3/getting-started/).   

1. You must create an S3 bucket, and write down its name and the AWS zone it
runs on.
2. Determine write permissions for your bucket.
3. S3 flavor comes with DEBUG=false by default. If you need to debug, then you need to add `-e DEBUG=True`.
4. Specify an AWS region, which is dependent on your S3 location, for example, use `-e AWS_REGION=”eu-west-1”`.
5. Ensure your host time is correct. If your host clock is still running on the main computer, but not on the docker host virtual machine, then you will have
time differences. This may cause an issue if you try to authenticate with Amazon
web services.
6. You will also need your AWS access key and secret key. Learn [more about it ](http://docs.aws.amazon.com/general/latest/gr/managing-aws-access-keys.html) here.

Conversely, you can further limit what users access in the Trusted Registry when you use AW to host your Trusted Registry. Instead of using the UI to enter information, you can create an [IAM user policy](http://docs.aws.amazon.com/AmazonS3/latest/dev/example-policies-s3.html) which is a JSON description of the effects, actions, and resources available to
a user. The advantage of using this method instead of configuring through the Trusted Registry UI is that you can restrict what users can access. You apply the policy as part of the process of installing the Trusted Registry on AW. To set a policy through the AWS command line, save the policy into a file,
for example `TrustedRegistryUserPerms.json`, and pass it to the
put-user-policy AWS command:

```
$ aws iam put-user-policy --user-name MyUser --policy-name TrustedRegistryUserPerms --policy-document file://C:\Temp\TrustedRegistryUserPerms.json
```

You can also set a policy through your AWS console. For more information about
setting IAM policies using the command line or the console, review the AWS
[Overview of IAM Policies](http://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html) article or visit the console Policies page.

The following example describes the minimum permissions set which allows
Trusted Registry users to access, push, pull, and delete images.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:ListAllMyBuckets",
            "Resource": "arn:aws:s3:::*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetBucketLocation"
            ],
            "Resource": "arn:aws:s3:::<INSERT YOUR BUCKET HERE>"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject"
            ],
            "Resource": "arn:aws:s3:::<INSERT YOUR BUCKET HERE>/*"
        }
    ]
}

```

### OpenStack Swift

OpenStack Swift, also known as OpenStack Object Storage, is an open source
object storage system that is licensed under the Apache 2.0 license. Refer to [Swift documentation](http://docs.openstack.org/developer/swift/) to get started.

<!--
### Ceph Rados

**(Details are missing as this is currently being developed for DTR 1.4.3?)**

Ceph implements distributed object storage. The object storage “product”,
service or capabilities, consist of a Ceph Storage Cluster and a Ceph
Object Gateway.

For additional information see the Ceph documentation  [here](http://docs.ceph.com/docs/master/rados/) and [here](http://docs.ceph.com/docs/hammer/radosgw/).
-->

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

1.  Navigate to the Trusted Registry UI > Settings > Storage.

2.  Select Download to get the text based file. It contains a minimum amount of
    information and you're going to need additional data based on your driver
    and business requirements.

3.  Go [here](/registry/configuration.md#list-of-configuration-options") to see
    the open source YAML file. Copy the sections you need and paste into your
    `storage.yml` file. Some settings may contradict others, so ensure your
    choices make sense.

4.  Save the YAML file and return to the UI.

5. On the Storage screen, upload the file, review your changes, and click Save.

## See also

* [Use your own certificates](index.md)
