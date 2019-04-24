---
description: Yandex.Cloud Compute driver for machine
keywords: machine, Yandex.Cloud, driver
title: Yandex.Cloud
---

Create machines on [Yandex.Cloud](http://cloud.yandex.com/).
You need a Yandex account and a folder ID.
See [https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy](https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy) for details on hierarchy of resources.

## Authentication

Before using the Yandex.Cloud driver, configure one of the following authentication methods.

### OAuth token

Follow the [instructions](https://cloud.yandex.com/docs/iam/concepts/authorization/oauth-token) to get OAuth access token and 
pass it to `docker-machine create` with the `--yandex-token` option:
   
    $ docker-machine create --driver yandex --yandex-token OAUTH_TOKEN dm01

### Service Account key file

For [authorization](https://cloud.yandex.com/docs/iam/concepts/authorization/#sa) with [Service Account](https://cloud.yandex.com/docs/iam/concepts/users/service-accounts)
you need to get one of [authorization keys](https://cloud.yandex.com/docs/iam/concepts/authorization/key). Follow the
step by step [instructions](https://cloud.yandex.com/docs/iam/operations/iam-token/create-for-sa#keys-create) to create
Service Account key file. Use path to that file with the `--yandex-sa-key-file` option:

    $ docker-machine create --driver yandex --yandex-sa-key-file PATH_TO_KEY_FILE dm01

> **NOTE:** Do not forget to set proper [roles](https://cloud.yandex.com/docs/iam/concepts/access-control/roles) for 
> Service Account.

## Example

To create a machine instance, specify `--driver yandex`, the folder ID and the machine name.

    $ docker-machine create --driver yandex --yandex-token OAUTH_TOKEN --yandex-folder-id FOLDER_ID dm01
    $ docker-machine create --driver yandex \
      --yandex-token OAUTH_TOKEN   \
      --yandex-folder-id FOLDER_ID \
      --yandex-zone ru-central1-a  \
      dm02

## Options

One of required:

-   `--yandex-sa-key-file`: Path to a file containing Service Account Key. 
-   `--yandex-token`: OAuth token to access Yandex.Cloud API.

Optional:
-   `--yandex-cloud-id`: Cloud ID to guess Folder ID for instance if latter is not specified.
-   `--yandex-cores`: The number of cores available to the instance.
-   `--yandex-disk-size`: The disk size of instance in gigabytes.
-   `--yandex-disk-type`: The disk type of instance, e.g. `network-hdd`.
-   `--yandex-folder-id`: Folder ID where instance should be created.
-   `--yandex-image-family`: The image family name to create the boot disk from.
-   `--yandex-image-folder-id`: Folder ID to lookup latest image ID by family name.
-   `--yandex-image-id`: ID of the image to create the boot disk from.
-   `--yandex-labels`: Instance labels in a form of `key=value` pairs.
-   `--yandex-memory`: The amount of memory available to the instance, specified in gigabytes.
-   `--yandex-platform-id`: ID of the hardware platform configuration for the instance.
-   `--yandex-preemptible`: Instance preemptibility.
-   `--yandex-ssh-port`: SSH port.
-   `--yandex-ssh-user`: SSH username.
-   `--yandex-subnet-id`: Specify subnet in which to provision VM.
-   `--yandex-use-internal-ip`:  If this option is specified, docker-machine uses internal rather than public NATed IPs. The flag is persistent in the sense that a machine created with it retains the IP. It's useful for managing docker machines from another machine on the same network, such as when deploying a swarm.
-   `--yandex-userdata`: Path to file containing User Data for the instance.
-   `--yandex-zone`: The zone to launch the instance.

Yandex.Cloud supports [image families](https://cloud.yandex.com/docs/compute/concepts/images#family).
The driver uses the `ubuntu-1604-lts` image family to create instance unless specified otherwise.
An image family is like an image alias that always points to the latest image in the family. To create an
instance from an image family, set `--yandex-image-family` to the family's name.

To get list of public available images and which family they belong to (if any) install [CLI](https://cloud.yandex.com/docs/cli/quickstart)
and run:

    $ yc compute images list --folder-id=standard-images

It is also possible to specify exact instance image ID by setting `--yandex-image-id`. 

#### Environment variables and default values

| CLI option                 | Environment variable  | Default           |
|:---------------------------|:----------------------|:------------------|
| `--yandex-cloud-id`        | `YC_CLOUD_ID`         | -                 |
| `--yandex-cores`           | `YC_CORES`            | 1                 |
| `--yandex-disk-size`       | `YC_DISK_SIZE`        | 20                |
| `--yandex-disk-type`       | `YC_DISK_TYPE`        | `network-hdd`     |
| `--yandex-folder-id`       | `YC_FOLDER_ID`        | -                 |
| `--yandex-image-family`    | `YC_IMAGE_FAMILY`     | `ubuntu-1604-lts` |
| `--yandex-image-folder-id` | `YC_IMAGE_FOLDER_ID`  | `standard-images` |
| `--yandex-image-id`        | `YC_IMAGE_ID`         | -                 |
| `--yandex-labels`          | `YC_LABELS`           | -                 |
| `--yandex-memory`          | `YC_MEMORY`           | 1                 |
| `--yandex-platform-id`     | `YC_PLATFORM_ID`      | `standard-v1`     |
| `--yandex-preemptible`     | `YC_PREEMPTIBLE`      | `false`           |
| `--yandex-sa-key-file`     | `YC_SA_KEY_FILE`      | -                 |
| `--yandex-ssh-port`        | `YC_SSH_PORT`         | `22`              |
| `--yandex-ssh-user`        | `YC_SSH_USER`         | `yc-user`         |
| `--yandex-subnet-id`       | `YC_SUBNET_ID`        | -                 |
| `--yandex-token`           | `YC_TOKEN`            | -                 |
| `--yandex-use-internal-ip` | `YC_USE_INTERNAL_IP`  | `false`           |
| `--yandex-userdata`        | `YC_USERDATA`         | -                 |
| `--yandex-zone`            | `YC_ZONE`             | `ru-central1-a`   |


## Folder ID

Default folder ID is determined automatically.
If your account or Service Account has access to more that one folder, you should specify a folder id with the `--yandex-folder-id` flag.

