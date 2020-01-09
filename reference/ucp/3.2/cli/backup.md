---
title: docker/ucp backup
description: Create a backup of a UCP manager node
keywords: ucp, cli, backup
---

>{% include enterprise_label_shortform.md %}

Create a backup of a UCP manager node.

## Usage

```bash
docker container run \
    --rm \
    --interactive \
    --name ucp \
    --log-driver none \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    backup [command options] > backup.tar
```

## Description

This command creates a tar file with the contents of the volumes used by
this UCP manager node, and prints it. You can then use the `restore` command to
restore the data from an existing backup.

To create backups of a multi-node cluster, you only need to back up a single
manager node. The restore operation will reconstitute a new UCP installation
from the backup of any previous manager.

Note:

  * The backup contains private keys and other sensitive information. Use the
    `--passphrase` flag to encrypt the backup with PGP-compatible encryption
    or `--no-passphrase` to opt out (not recommended).

  * If using the `--file` option, the path to the file must be bind mounted
    onto the container that is performing the backup, and the filepath must be
    relative to the container's file tree. For example:

    ```
    docker run <other options> --mount type=bind,src=/home/user/backup:/backup docker/ucp --file /backup/backup.tar
    ```

### SELinux

If you are installing UCP on a manager node with SELinunx enabled at the daemon
and operating system level, you will need to pass `--security-opt
label=disable` in to your install command. This flag will disable SELinux
policies on the installation container.  The UCP installation container mounts
and configures the Docker Socket as part of the UCP installation container,
therefore the UCP installation will fail with a permission denied error if you
fail to pass in this flag.

```
FATA[0000] unable to get valid Docker client: unable to ping Docker daemon: Got
permission denied while trying to connect to the Docker daemon socket at
unix:///var/run/docker.sock: Get http://%2Fvar%2Frun%2Fdocker.sock/_ping: dial
unix /var/run/docker.sock: connect: permission denied - If SELinux is enabled
on the Docker daemon, make sure you run UCP with "docker run --security-opt
label=disable -v /var/run/docker.sock:/var/run/docker.sock ..."
```

An installation command for a system with SELinux enabled at the daemon level
would be:

```bash
docker container run \
    --rm \
    --interactive \
    --name ucp \
    --security-opt label=disable \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    backup [command options] > backup.tar
```

## Options

| Option                 | Description                                                                                                                                                                 |
|:-----------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--debug, -D`          | Enable debug mode                                                                                                                                                           |
| `--file *value*`       | Name of the file to write the backup contents to. Ignored in interactive mode                                                                                               |
| `--jsonlog`            | Produce json formatted output for easier parsing                                                                                                                            |
| `--include-logs`       | Only relevant if `--file` is also included. If true, an encrypted `backup.log` file will be stored alongside the `backup.tar` in the mounted directory. Default is `true`.  |
| `--interactive, -i`    | Run in interactive mode and prompt for configuration values                                                                                                                 |
| `--no-passphrase`      | Opt out to encrypt the tar file with a passphrase (not recommended)                                                                                                         |
| `--passphrase` *value* | Encrypt the tar file with a passphrase                                                                                                                                      |
