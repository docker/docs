---
title: docker/ucp backup
description: Create a backup of a UCP manager node
keywords: docker, ucp, cli, backup
---

Create a backup of a UCP manager node

## Description

This command creates a tar file with the contents of the volumes used by
this UCP manager node, and prints it. You can then use the `restore` command to
restore the data from an existing backup.

To create backups of a multi-node swarm, you only need to backup a single manager
node. The restore operation will reconstitute a new UCP installation from the
backup of any previous manager.

Note:

  * During the backup, UCP is temporarily stopped. This does not affect your
    applications.

  * The backup contains private keys and other sensitive information. Use the
    `--passphrase` flag to encrypt the backup with PGP-compatible encryption.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--interactive, i`|Run in interactive mode and prompt for configuration values|
|`--id`|The ID of the UCP instance to backup|
|`--passphrase`|Encrypt the tar file with a passphrase|
