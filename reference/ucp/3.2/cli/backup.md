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
  * If using the `--file` option, the path to the file must be bind mounted onto the container that is performing the backup, and the filepath must be relative to the container's file tree. For example:
  ```
  docker run <other options> --mount type=bind,src=/home/user/backup:/backup docker/ucp --file /backup/backup.tar
  ```
    
> **Note**: A bind mount with a `/backup/` target path must be added to the container performing the backup. In this case, the backup file is placed in the source directory of the bind mount. For example:
>    ```
>    docker run -v /nfs/ucp-backups:/backup docker/ucp-backup --file backup1.tar
>    ```
>    This command places the backup under the `/nfs/ucp-backups/backup1.tar` path on the host.

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--file`|Name of the file for backup contents. This is ignored when `--interactive` is specified. If not provided, backup contents are sent to stdout.|
|`--include-logs`|Only applicable is `--file` is specified. If `true`, includes logs from the backup execution in a file adjacent to the backup file, specified with the same name but with a `.log` extension. **Note**: Log files are not encrypted.|
|`--interactive, i`|Run in interactive mode and prompt for configuration values|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--passphrase`|Encrypt the tar file with a passphrase|
