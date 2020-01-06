---
description: VMware vSphere driver for machine
keywords: machine, VMware vSphere, driver
title: VMware vSphere
hide_from_sitemap: true
---

Creates machines on a [VMware vSphere](http://www.vmware.com/products/vsphere)
Virtual Infrastructure. The machine must have a working vSphere ESXi
installation. You can use a paid license or free 60 day trial license. Your
installation may also include an optional VCenter server.

## Usage

    $ docker-machine create --driver vmwarevsphere --vmwarevsphere-username=user --vmwarevsphere-password=SECRET vm

## Options

-   `--vmwarevsphere-boot2docker-url`: URL for boot2docker image.
-   `--vmwarevsphere-cpu-count`: CPU number for Docker VM.
-   `--vmwarevsphere-datacenter`: Datacenter for Docker VM (must be set to `ha-datacenter` when connecting to a single host).
-   `--vmwarevsphere-datastore`: Datastore for Docker VM.
-   `--vmwarevsphere-disk-size`: Size of disk for Docker VM (in MB).
-   `--vmwarevsphere-folder`: vSphere folder for the docker VM. This folder must already exist in the datacenter.
-   `--vmwarevsphere-hostsystem`: vSphere compute resource where the docker VM is instantiated. This can be omitted if using a cluster with DRS.
-   `--vmwarevsphere-memory-size`: Size of memory for Docker VM (in MB).
-   `--vmwarevsphere-network`: Network where the Docker VM is attached.
-   `--vmwarevsphere-password`: **required** vSphere Password.
-   `--vmwarevsphere-pool`: Resource pool for Docker VM.
-   `--vmwarevsphere-username`: **required** vSphere Username.
-   `--vmwarevsphere-vcenter-port`: vSphere Port for vCenter.
-   `--vmwarevsphere-vcenter`: IP/hostname for vCenter (or ESXi if connecting directly to a single host).

The VMware vSphere driver uses the latest boot2docker image.

#### Environment variables and default values

| CLI option                        | Environment variable      | Default                  |
| --------------------------------- | ------------------------- | ------------------------ |
| `--vmwarevsphere-boot2docker-url` | `VSPHERE_BOOT2DOCKER_URL` | _Latest boot2docker url_ |
| `--vmwarevsphere-cpu-count`       | `VSPHERE_CPU_COUNT`       | `2`                      |
| `--vmwarevsphere-datacenter`      | `VSPHERE_DATACENTER`      | -                        |
| `--vmwarevsphere-datastore`       | `VSPHERE_DATASTORE`       | -                        |
| `--vmwarevsphere-disk-size`       | `VSPHERE_DISK_SIZE`       | `20000`                  |
| `--vmwarevsphere-hostsystem`      | `VSPHERE_HOSTSYSTEM`      | -                        |
| `--vmwarevsphere-memory-size`     | `VSPHERE_MEMORY_SIZE`     | `2048`                   |
| `--vmwarevsphere-network`         | `VSPHERE_NETWORK`         | -                        |
| **`--vmwarevsphere-password`**    | `VSPHERE_PASSWORD`        | -                        |
| `--vmwarevsphere-pool`            | `VSPHERE_POOL`            | -                        |
| **`--vmwarevsphere-username`**    | `VSPHERE_USERNAME`        | -                        |
| `--vmwarevsphere-vcenter-port`    | `VSPHERE_VCENTER_PORT`    | 443                      |
| `--vmwarevsphere-vcenter`         | `VSPHERE_VCENTER`         | -                        |
| `--vmwarevsphere-folder`          | `VSPHERE_FOLDER`          | -                        |
