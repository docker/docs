---
description: Oracle VirtualBox driver for machine
keywords: machine, Oracle VirtualBox, driver
title: Oracle VirtualBox
hide_from_sitemap: true
---

Create machines locally using [VirtualBox](https://www.virtualbox.org/).
This driver requires VirtualBox 5+ to be installed on your host.
Using VirtualBox 4.3+ should work but emits a warning. Older versions
do not work.

## Usage

    $ docker-machine create --driver=virtualbox vbox-test

You can create an entirely new machine or you can convert a Boot2Docker VM into
a machine by importing the VM. To convert a Boot2Docker VM, you'd use the following
command:

    $ docker-machine create -d virtualbox --virtualbox-import-boot2docker-vm boot2docker-vm b2d

The size of the VM's disk can be configured this way:

    $ docker-machine create -d virtualbox --virtualbox-disk-size "100000" large

## Options

-   `--virtualbox-boot2docker-url`: The URL of the boot2docker image. Defaults to the latest available version.
-   `--virtualbox-cpu-count`: Number of CPUs to use to create the VM. Defaults to single CPU.
-   `--virtualbox-disk-size`: Size of disk for the host in MB.
-   `--virtualbox-host-dns-resolver`: Use the host DNS resolver. (Boolean value, defaults to false)
-   `--virtualbox-hostonly-cidr`: The CIDR of the host only adapter.
-   `--virtualbox-hostonly-nicpromisc`: Host Only Network Adapter Promiscuous Mode. Possible options are deny , allow-vms, allow-all
-   `--virtualbox-hostonly-nictype`: Host Only Network Adapter Type. Possible values are '82540EM' (Intel PRO/1000), 'Am79C973' (PCnet-FAST III), and 'virtio' Paravirtualized network adapter.
-   `--virtualbox-hostonly-no-dhcp`: Disable the Host Only DHCP Server
-   `--virtualbox-import-boot2docker-vm`: The name of a Boot2Docker VM to import.
-   `--virtualbox-memory`: Size of memory for the host in MB.
-   `--virtualbox-nat-nictype`: Specify the NAT Network Adapter Type. Possible values are '82540EM' (Intel PRO/1000), 'Am79C973' (PCnet-FAST III) and 'virtio' Paravirtualized network adapter.
-   `--virtualbox-no-dns-proxy`: Disable proxying all DNS requests to the host (Boolean value, default to false)
-   `--virtualbox-no-share`: Disable the mount of your home directory
-   `--virtualbox-no-vtx-check`: Disable checking for the availability of hardware virtualization before the vm is started
-   `--virtualbox-share-folder`: Mount the specified directory instead of the default home location.
-   `--virtualbox-ui-type`: Specify the UI Type: (gui|sdl|headless|separate)

The `--virtualbox-boot2docker-url` flag takes a few different forms. By
default, if no value is specified for this flag, Machine checks locally for
a boot2docker ISO. If one is found, it is used as the ISO for the
created machine. If one is not found, the latest ISO release available on
[boot2docker/boot2docker](https://github.com/boot2docker/boot2docker) is
downloaded and stored locally for future use. Therefore, you must run
`docker-machine upgrade` deliberately on a machine if you wish to update the "cached"
boot2docker ISO.

This is the default behavior (when `--virtualbox-boot2docker-url=""`), but the
option also supports specifying ISOs by the `http://` and `file://` protocols.
`file://` looks at the path specified locally to locate the ISO: for
instance, you could specify `--virtualbox-boot2docker-url
file://$HOME/Downloads/rc.iso` to test out a release candidate ISO that you have
downloaded already. You could also just get an ISO straight from the Internet
using the `http://` form.

To customize the host only adapter, you can use the `--virtualbox-hostonly-cidr`
flag. This specifies the host IP and Machine calculates the VirtualBox
DHCP server address (a random IP on the subnet between `.1` and `.25`) so
it does not clash with the specified host IP.
Machine specifies the DHCP lower bound to `.100` and the upper bound
to `.254`. For example, a specified CIDR of `192.168.24.1/24` would have a
DHCP server between `192.168.24.2-25`, a lower bound of `192.168.24.100` and
upper bound of `192.168.24.254`.

With the flag `--virtualbox-share-folder`, you can specify which folder the host 
shares with the created machine. The format is `local-folder:machine-folder`. 
For example, `\\?\C:\docker-share:\home\users\`. if you specify the flag with the
docker-toolbox using docker-machine from a Windows cmd, it looks like 
`C:\docker-share\\:/home/users`. The `:` sign needs to be escaped.

#### Environment variables and default values

| CLI option                           | Environment variable               | Default                  |
|:-------------------------------------|:-----------------------------------|:-------------------------|
| `--virtualbox-boot2docker-url`       | `VIRTUALBOX_BOOT2DOCKER_URL`       | _Latest boot2docker url_ |
| `--virtualbox-cpu-count`             | `VIRTUALBOX_CPU_COUNT`             | `1`                      |
| `--virtualbox-disk-size`             | `VIRTUALBOX_DISK_SIZE`             | `20000`                  |
| `--virtualbox-host-dns-resolver`     | `VIRTUALBOX_HOST_DNS_RESOLVER`     | `false`                  |
| `--virtualbox-hostonly-cidr`         | `VIRTUALBOX_HOSTONLY_CIDR`         | `192.168.99.1/24`        |
| `--virtualbox-hostonly-nicpromisc`   | `VIRTUALBOX_HOSTONLY_NIC_PROMISC`  | `deny`                   |
| `--virtualbox-hostonly-nictype`      | `VIRTUALBOX_HOSTONLY_NIC_TYPE`     | `82540EM`                |
| `--virtualbox-hostonly-no-dhcp`      | `VIRTUALBOX_HOSTONLY_NO_DHCP`      | `false`                  |
| `--virtualbox-import-boot2docker-vm` | `VIRTUALBOX_BOOT2DOCKER_IMPORT_VM` | `boot2docker-vm`         |
| `--virtualbox-memory`                | `VIRTUALBOX_MEMORY_SIZE`           | `1024`                   |
| `--virtualbox-nat-nictype`           | `VIRTUALBOX_NAT_NICTYPE`           | `82540EM`                |
| `--virtualbox-no-dns-proxy`          | `VIRTUALBOX_NO_DNS_PROXY`          | `false`                  |
| `--virtualbox-no-share`              | `VIRTUALBOX_NO_SHARE`              | `false`                  |
| `--virtualbox-no-vtx-check`          | `VIRTUALBOX_NO_VTX_CHECK`          | `false`                  |
| `--virtualbox-share-folder`          | `VIRTUALBOX_SHARE_FOLDER`          | -                        |
| `--virtualbox-ui-type`               | `VIRTUALBOX_UI_TYPE`               | `headless`               |

## Known Issues

Vboxfs suffers from a [longstanding bug](https://www.virtualbox.org/ticket/9069)
causing [sendfile(2)](http://linux.die.net/man/2/sendfile) to serve cached file
contents.

This causes problems when using a web server such as Nginx to serve
static files from a shared volume. For development environments, a good
workaround is to disable sendfile in your server configuration.
