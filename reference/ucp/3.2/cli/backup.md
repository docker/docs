---
title: docker/ucp backup
description: Create a backup of a UCP manager node
keywords: ucp, cli, backup
---

Create a backup of a UCP manager node.

## Usage

```bash
docker container run --log-driver none --rm -i \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
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

## Options

| Option                 | Description                                                                   |
|:-----------------------|:------------------------------------------------------------------------------|
| `--debug, -D`          | Enable debug mode                                                             |
|  --file *value*        | Name of the file to write the backup contents to. Ignored in interactive mode |
| `--jsonlog`            | Produce json formatted output for easier parsing                              |
| `--interactive, -i`    | Run in interactive mode and prompt for configuration values                   |
| `--no-passphrase`      | Opt out to encrypt the tar file with a passphrase (not recommended)           |
| `--passphrase` *value* | Encrypt the tar file with a passphrase                                        |
