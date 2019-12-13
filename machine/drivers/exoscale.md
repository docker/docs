---
description: Exoscale driver for machine
keywords: machine, exoscale, driver
title: Exoscale
hide_from_sitemap: true
---

Create machines on [Exoscale](https://www.exoscale.com/).

Get your API key and API secret key from [API details](https://portal.exoscale.com/account/api) and pass them to `machine create` with the `--exoscale-api-key` and `--exoscale-api-secret-key` options.

## Usage

    $ docker-machine create --driver exoscale \
        --exoscale-api-key=API \
        --exoscale-api-secret-key=SECRET \
        MY_COMPUTE_INSTANCE

If you encounter any troubles, activate the debug mode with `docker-machine --debug create ...`.

## Options

-   `--exoscale-affinity-group`: [Anti-affinity group][anti-affinity] the machine is started in.
-   `--exoscale-api-key`: **required** Your API key;
-   `--exoscale-api-secret-key`: **required** Your API secret key;
-   `--exoscale-availability-zone`: Exoscale [availability zone][datacenters] (ch-dk-2, at-vie-1, de-fra-1, ...);
-   `--exoscale-disk-size`: Disk size for the host in GiB (at least 10);
-   `--exoscale-image`: Image template, for example `ubuntu-16.04`, also known as `Linux Ubuntu 16.04 LTS 64-bit`, [see below](#image-template-name));
-   `--exoscale-instance-profile`: Instance profile (Small, Medium, Large, ...);
-   `--exoscale-security-group`: Security group. _It is created if it doesn't exist_;
-   `--exoscale-ssh-key`: Path to the SSH user private key. _A new one is created if left empty_;
-   `--exoscale-ssh-user`: SSH username to connect, such as `ubuntu`, [see below](#ssh-username));
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
| `--exoscale-ssh-key`            | `EXOSCALE_SSH_KEY`           | -                                 |
| `--exoscale-ssh-user`           | `EXOSCALE_SSH_USER`          | -                                 |
| `--exoscale-url`                | `EXOSCALE_ENDPOINT`          | `https://api.exoscale.ch/compute` |
| `--exoscale-userdata`           | `EXOSCALE_USERDATA`          | -                                 |

**NB:** the _instance profile_, _image_, and _availability zone_ are case insensitive.

### Image template name

The [VM templates][templates] available at Exoscale are listed on the Portal
when adding a new instance.

For any Linux template, you may use the shorter name composed only of the name
and version, as shown below.

| Full name                       | Short name           |
| ------------------------------- | -------------------- |
| Linux Debian 8 64-bit           | `debian-8`           |
| Linux Ubuntu 16.04 LTS 64-bit   | `ubuntu-16.04`       |
| Linux CentOS 7.3 64-bit         | `centos-7.3`         |
| Linux CoreOS stable 1298 64-bit | `coreos-stable-1298` |

**NB:** Docker doesn't work for non-Linux machines like OpenBSD or Windows Server.

### SSH Username

The Exoscale driver does an educated guess to pick the correct default SSH
user. If left empty, it picks a suitable one following those rules:

- `centos` for CentOS;
- `core` for Linux CoreOS (aka Container Linux);
- `debian` for Debian;
- `ubuntu` for Ubuntu;
- `fedora` for Fedora;
- `cloud-user` for Red Hat;
- otherwise, `root`.

### Custom security group

If a custom security group is provided, you need to ensure that you allow TCP ports 22 and 2376 in an ingress rule.

Moreover, if you want to use [Docker Swarm](/engine/swarm/swarm-tutorial/), also add TCP port 2377, UDP/TCP on 7946, and UDP on 4789.

### Debian 9

The [default storage driver][storagedriver] may fail on Debian, specifying `overlay2` should resolve this issue.

    $ docker-machine create --engine-storage-driver overlay2 ...`

### More than 8 docker machines?

There is a limit to the number of machines that an anti-affinity group can have.  This can be worked around by specifying an additional anti-affinity group using `--exoscale-affinity-group=docker-machineX`

[storagedriver]: https://docs.docker.com/storage/storagedriver/select-storage-driver/#docker-ce
[templates]: https://www.exoscale.com/templates/
[datacenters]: https://www.exoscale.com/datacenters/
[anti-affinity]: https://community.exoscale.com/documentation/compute/anti-affinity-groups/
