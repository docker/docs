---
description: OpenStack driver for machine
keywords: machine, OpenStack, driver
title: OpenStack
hide_from_sitemap: true
---

Create machines on [OpenStack](http://www.openstack.org/software/)

Mandatory:

-   `--openstack-auth-url`: Keystone service base URL.
-   `--openstack-flavor-id` or `--openstack-flavor-name`: Identify the flavor used for the machine.
-   `--openstack-image-id` or `--openstack-image-name`: Identify the image used for the machine.

## Usage

    $ docker-machine create --driver openstack vm

## Options

-   `--openstack-active-timeout`: The timeout in seconds until the OpenStack instance must be active.
-   `--openstack-availability-zone`: The availability zone in which to launch the server.
-   `--openstack-config-drive`: Whether OpenStack should mount a configuration drive for the machine.
-   `--openstack-domain-name` or `--openstack-domain-id`: Domain to use for authentication (Keystone v3 only).
-   `--openstack-endpoint-type`: Endpoint type can be `internalURL`, `adminURL`, or `publicURL`. It is a helper for the driver
    to choose the right URL in the OpenStack service catalog. If not provided the default is `publicURL`.
-   `--openstack-floatingip-pool`: The IP pool used to get a public IP can assign it to the machine. If there is an
    IP address already allocated but not assigned to any machine, this IP is chosen and assigned to the machine. If
    there is no IP address already allocated, a new IP is allocated and assigned to the machine.
-   `--openstack-keypair-name`: Specify the existing Nova keypair to use.
-   `--openstack-insecure`: Explicitly allow openstack driver to perform "insecure" SSL (https) requests. The server's certificate is not verified against any certificate authorities. This option should be used with caution.
-   `--openstack-ip-version`: If the instance has both IPv4 and IPv6 address, you can select IP version. If not provided, defaults to `4`.
-   `--openstack-net-name` or `--openstack-net-id`: Identify the private network the machine is connected to. If your OpenStack project contains only one private network it is used automatically.
-   `--openstack-password`: User password. It can be omitted if the standard environment variable `OS_PASSWORD` is set.
-   `--openstack-private-key-file`: Used with `--openstack-keypair-name`, associates the private key to the keypair.
-   `--openstack-region`: The region to work on. Can be omitted if there is only one region on the OpenStack.
-   `--openstack-sec-groups`: If security groups are available on your OpenStack you can specify a comma separated list
    to use for the machine, such as `secgrp001,secgrp002`.
-   `--openstack-ssh-port`: Customize the SSH port if the SSH server on the machine does not listen on the default port.
-   `--openstack-ssh-user`: The username to use for SSH into the machine. If not provided defaults to `root`.
-   `--openstack-tenant-name` or `--openstack-tenant-id`: Identify the tenant in which the machine is created.
-   `--openstack-user-data-file`: File containing an OpenStack userdata script.
-   `--openstack-username`: User identifier to authenticate with.

#### Environment variables and default values

| CLI option                      | Environment variable   | Default     |
| ------------------------------- | ---------------------- | ----------- |
| `--openstack-active-timeout`    | `OS_ACTIVE_TIMEOUT`    | `200`       |
| `--openstack-auth-url`          | `OS_AUTH_URL`          | -           |
| `--openstack-availability-zone` | `OS_AVAILABILITY_ZONE` | -           |
| `--openstack-config-drive`      | `OS_CONFIG_DRIVE`      | `false`     |
| `--openstack-domain-id`         | `OS_DOMAIN_ID`         | -           |
| `--openstack-domain-name`       | `OS_DOMAIN_NAME`       | -           |
| `--openstack-endpoint-type`     | `OS_ENDPOINT_TYPE`     | `publicURL` |
| `--openstack-flavor-id`         | `OS_FLAVOR_ID`         | -           |
| `--openstack-flavor-name`       | `OS_FLAVOR_NAME`       | -           |
| `--openstack-floatingip-pool`   | `OS_FLOATINGIP_POOL`   | -           |
| `--openstack-image-id`          | `OS_IMAGE_ID`          | -           |
| `--openstack-image-name`        | `OS_IMAGE_NAME`        | -           |
| `--openstack-insecure`          | `OS_INSECURE`          | `false`     |
| `--openstack-ip-version`        | `OS_IP_VERSION`        | `4`         |
| `--openstack-keypair-name`      | `OS_KEYPAIR_NAME`      | -           |
| `--openstack-net-id`            | `OS_NETWORK_ID`        | -           |
| `--openstack-net-name`          | `OS_NETWORK_NAME`      | -           |
| `--openstack-password`          | `OS_PASSWORD`          | -           |
| `--openstack-private-key-file`  | `OS_PRIVATE_KEY_FILE`  | -           |
| `--openstack-region`            | `OS_REGION_NAME`       | -           |
| `--openstack-sec-groups`        | `OS_SECURITY_GROUPS`   | -           |
| `--openstack-ssh-port`          | `OS_SSH_PORT`          | `22`        |
| `--openstack-ssh-user`          | `OS_SSH_USER`          | `root`      |
| `--openstack-tenant-id`         | `OS_TENANT_ID`         | -           |
| `--openstack-tenant-name`       | `OS_TENANT_NAME`       | -           |
| `--openstack-user-data-file`    | `OS_USER_DATA_FILE`    | -           |
| `--openstack-username`          | `OS_USERNAME`          | -           |
