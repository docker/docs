---
description: VMware vCloud Air driver for machine
keywords: machine, VMware vCloud Air, driver
title: VMware vCloud Air
hide_from_sitemap: true
---

Creates machines on [vCloud Air](http://vcloud.vmware.com) subscription service.
You need an account within an existing subscription of vCloud Air VPC or
Dedicated Cloud.

## Usage

    $ docker-machine create --driver vmwarevcloudair --vmwarevcloudair-username=user --vmwarevcloudair-password=SECRET vm

## Options

-   `--vmwarevcloudair-catalog`: Catalog.
-   `--vmwarevcloudair-catalogitem`: Catalog Item.
-   `--vmwarevcloudair-computeid`: Compute ID (if using Dedicated Cloud).
-   `--vmwarevcloudair-cpu-count`: VM CPU Count.
-   `--vmwarevcloudair-docker-port`: Docker port.
-   `--vmwarevcloudair-edgegateway`: Organization Edge Gateway.
-   `--vmwarevcloudair-memory-size`: VM Memory Size in MB.
-   `--vmwarevcloudair-orgvdcnetwork`: Organization VDC Network to attach.
-   `--vmwarevcloudair-password`: **required** vCloud Air Password.
-   `--vmwarevcloudair-provision`: Install Docker binaries.
-   `--vmwarevcloudair-publicip`: Org Public IP to use.
-   `--vmwarevcloudair-ssh-port`: SSH port.
-   `--vmwarevcloudair-username`: **required** vCloud Air Username.
-   `--vmwarevcloudair-vdcid`: Virtual Data Center ID.

The VMware vCloud Air driver uses the `Ubuntu Server 12.04 LTS (amd64 20140927)` image by default.

#### Environment variables and default values

| CLI option                        | Environment variable      | Default                                    |
| --------------------------------- | ------------------------- | ------------------------------------------ |
| `--vmwarevcloudair-catalog`       | `VCLOUDAIR_CATALOG`       | `Public Catalog`                           |
| `--vmwarevcloudair-catalogitem`   | `VCLOUDAIR_CATALOGITEM`   | `Ubuntu Server 12.04 LTS (amd64 20140927)` |
| `--vmwarevcloudair-computeid`     | `VCLOUDAIR_COMPUTEID`     | -                                          |
| `--vmwarevcloudair-cpu-count`     | `VCLOUDAIR_CPU_COUNT`     | `1`                                        |
| `--vmwarevcloudair-docker-port`   | `VCLOUDAIR_DOCKER_PORT`   | `2376`                                     |
| `--vmwarevcloudair-edgegateway`   | `VCLOUDAIR_EDGEGATEWAY`   | `<vdcid>`                                  |
| `--vmwarevcloudair-memory-size`   | `VCLOUDAIR_MEMORY_SIZE`   | `2048`                                     |
| `--vmwarevcloudair-orgvdcnetwork` | `VCLOUDAIR_ORGVDCNETWORK` | `<vdcid>-default-routed`                   |
| **`--vmwarevcloudair-password`**  | `VCLOUDAIR_PASSWORD`      | -                                          |
| `--vmwarevcloudair-provision`     | `VCLOUDAIR_PROVISION`     | `true`                                     |
| `--vmwarevcloudair-publicip`      | `VCLOUDAIR_PUBLICIP`      | -                                          |
| `--vmwarevcloudair-ssh-port`      | `VCLOUDAIR_SSH_PORT`      | `22`                                       |
| **`--vmwarevcloudair-username`**  | `VCLOUDAIR_USERNAME`      | -                                          |
| `--vmwarevcloudair-vdcid`         | `VCLOUDAIR_VDCID`         | -                                          |