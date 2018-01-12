---
description: Exoscale driver for machine
keywords: machine, exoscale, driver
title: Exoscale
---

Create machines on [Exoscale](https://www.exoscale.ch/).

Get your API key and API secret key from [API details](https://portal.exoscale.ch/account/api) and pass them to `machine create` with the `--exoscale-api-key` and `--exoscale-api-secret-key` options.

## Usage

    $ docker-machine create --driver exoscale \
        --exoscale-api-key=API \
        --exoscale-api-secret-key=SECRET \
        vm

## Options

-   `--exoscale-affinity-group`: [Anti-affinity group][anti-affinity] the machine will be started in.
-   `--exoscale-api-key`: **required** Your API key;
-   `--exoscale-api-secret-key`: **required** Your API secret key;
-   `--exoscale-availability-zone`: Exoscale [availability zone][datacenters] (CH-DK-2, AT-VIE-1, DE-FRA-1, ...);
-   `--exoscale-disk-size`: Disk size for the host in GB (10, 50, 100, 200, 400);
-   `--exoscale-image`: Image template (e.g. `Linux Ubuntu 16.04 LTS 64-bit` also known as `ubuntu-16.04`, [see below](#image-template-name));
-   `--exoscale-instance-profile`: Instance profile (Small, Medium, Large, ...);
-   `--exoscale-security-group`: Security group. _It will be created if it doesn't exist_;
-   `--exoscale-ssh-user`: SSH username (e.g. `ubuntu`, [see below](#ssh-username));
-   `--exoscale-url`: Your API endpoint;
-   `--exoscale-userdata`: Path to file containing user data for [cloud-init](https://cloud-init.io/);

### Environment variables and default values

| CLI option                      | Environment variable         | Default                           |
| ------------------------------- | ---------------------------- | --------------------------------- |
| `--exoscale-affinity-group`     | `EXOSCALE_AFFINITY_GROUP`    | -                                 |
| **`--exoscale-api-key`**        | `EXOSCALE_API_KEY`           | -                                 |
| **`--exoscale-api-secret-key`** | `EXOSCALE_API_SECRET`        | -                                 |
| `--exoscale-availability-zone`  | `EXOSCALE_AVAILABILITY_ZONE` | `ch-dk-2`                         |
| `--exoscale-disk-size`          | `EXOSCALE_DISK_SIZE`         | `50`                              |
| `--exoscale-image`              | `EXOSCALE_IMAGE`             | `Linux Ubuntu 16.04 LTS 64-bit`   |
| `--exoscale-instance-profile`   | `EXOSCALE_INSTANCE_PROFILE`  | `small`                           |
| `--exoscale-security-group`     | `EXOSCALE_SECURITY_GROUP`    | `docker-machine`                  |
| `--exoscale-ssh-user`           | `EXOSCALE_SSH_USER`          | -                                 |
| `--exoscale-url`                | `EXOSCALE_ENDPOINT`          | `https://api.exoscale.ch/compute` |
| `--exoscale-userdata`           | `EXOSCALE_USERDATA`          | -                                 |

**NB:** the _instance profile_, _image_, and _availability zone_ are case insensitive.

### Image template name

The [VM templates][templates] available at Exoscale are listed on the Portal when adding a new instance.

For any Linux template, you may use the shorter name composed only of the name and version. E.g.

| Full name                       | Short name           |
| ------------------------------- | -------------------- |
| Linux Debian 8 64-bit           | `debian-8`           |
| Linux Ubuntu 16.04 LTS 64-bit   | `ubuntu-16.04`       |
| Linux CentOS 7.3 64-bit         | `centos-7.3`         |
| Linux CoreOS stable 1298 64-bit | `coreos-stable-1298` |

**NB:** Docker won't work for non-Linux machines like OpenBSD and Windows Server.

### SSH Username

The exoscale driver does a wild guess to match the default SSH user. If left empty, it picks a suitable one:

- `centos` for Centos 7.3+;
- `core` for Linux CoreOS;
- `debian` for Debian 8+;
- `ubuntu` for Ubuntu;
- otherwise, `root`.

### Custom security group

If a custom security group is provided, you need to ensure that you allow TCP ports 22 and 2376 in an ingress rule.

Moreover, if you want to use [Docker Swarm](/engine/swarm/swarm-tutorial/), also add TCP port 2377.

### More than 8 docker machines?

There is a limit to the number of machines that an anti-affinity group can have.  This can be worked around by specifying an additional anti-affinity group using `--exoscale-affinity-group=docker-machineX`

[templates]: https://www.exoscale.ch/open-cloud/templates/
[datacenters]: https://www.exoscale.ch/infrastructure/datacenters/
[anti-affinity]: https://community.exoscale.ch/documentation/compute/anti-affinity-groups/
