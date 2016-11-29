---
description: Stream a tar file to stdout containing all UCP data volumes.
keywords: docker, ucp, backup, restore
title: docker/ucp backup
---

Stream a tar file to stdout containing all UCP data volumes.

## Usage

```none
docker run --rm -i \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  backup [command options] > backup.tar
```

## Description

This utility will dump out a tar file containing all the contents of the
volumes used by UCP on this node.  This can be used to make periodic
backups suitable for use in the 'restore' command.

When backing up an HA cluster, take backups of all controllers one at
a time, in quick succession, and keep track of the exact time when you
performed each backup.  You will need this timestamp information if you
restoring multiple controllers together.

WARNING

* During the backup, all UCP containers will be temporarily stopped
on this node to prevent data corruption.
* This backup will contain private keys and other sensitive information
and should be stored securely.  You may use the '--passphrase' flag to enable
built-in PGP compatible encryption.

## Options

| Option              | Description                                                                      |
|:--------------------|:---------------------------------------------------------------------------------|
| `--debug, -D`       | Enable debug mode                                                                |
| `--jsonlog`         | Produce json formatted output for easier parsing                                 |
| `--id`              | The ID of the UCP instance to backup                                             |
| `--root-ca-only`    | Backup only the root CA certificates and keys from this controller node          |
| `--passphrase`      | Encrypt the tar file with the provided passphrase [$UCP_PASSPHRASE]              |
| `--interactive, -i` | Enable interactive mode. You will be prompted to enter all required information. |