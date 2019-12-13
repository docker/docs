---
description: Generic driver for machine
keywords: machine, Generic, driver
title: Generic
hide_from_sitemap: true
---

Create machines using an existing VM/Host with SSH.

This is useful if you are using a provider that Machine does not support
directly or if you would like to import an existing host to allow Docker
Machine to manage.

The driver performs a list of tasks on create:

-   If docker is not running on the host, it is installed automatically.
-   It updates the host packages (`apt-get update`, `yum update`...).
-   It generates certificates to secure the docker daemon.
-   If the host uses systemd, it creates /etc/systemd/system/docker.service.d/10-machine.conf
-   The docker daemon restarts, thus all running containers are stopped.
-   The hostname is updated to fit the machine name.


### Example

To create a machine instance, specify `--driver generic`, the IP address or DNS
name of the host and the path to the SSH private key authorized to connect
to the host.

    $ docker-machine create \
      --driver generic \
      --generic-ip-address=203.0.113.81 \
      --generic-ssh-key ~/.ssh/id_rsa \
      vm

### Sudo privileges

The user that is used to SSH into the host can be specified with
`--generic-ssh-user` flag. This user needs password-less sudo
privileges.
If it's not the case, you need to edit the `sudoers` file and configure the user
as a sudoer with `NOPASSWD`. See https://help.ubuntu.com/community/Sudoers.

### Options

-   `--generic-engine-port`: Port to use for Docker Daemon (Note: This flag does not work with boot2docker).
-   `--generic-ip-address`: **required** IP Address of host.
-   `--generic-ssh-key`: Path to the SSH user private key.
-   `--generic-ssh-user`: SSH username used to connect.
-   `--generic-ssh-port`: Port to use for SSH.

> **Note**: You must use a base operating system supported by Machine.

#### Environment variables and default values

| CLI option                 | Environment variable | Default                   |
| -------------------------- | -------------------- | ------------------------- |
| `--generic-engine-port`    | `GENERIC_ENGINE_PORT`| `2376`                    |
| **`--generic-ip-address`** | `GENERIC_IP_ADDRESS` | -                         |
| `--generic-ssh-key`        | `GENERIC_SSH_KEY`    | -                         |
| `--generic-ssh-user`       | `GENERIC_SSH_USER`   | `root`                    |
| `--generic-ssh-port`       | `GENERIC_SSH_PORT`   | `22`                      |

### Systemd settings

For systems that use systemd, if you have an existing configuration defined in
'/etc/systemd/system/docker.service.d/' this  may conflict with the settings created by
docker-machine.  Make sure you don't have any other configuration files in this location
that override the [ExecStart] setting.

Once you have confirmed any conflicting settings have been removed, run
`sudo systemctl daemon reload` followed by `sudo systemctl restart docker`


