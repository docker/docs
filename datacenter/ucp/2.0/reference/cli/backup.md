---
description: Create a backup of a UCP manager node
keywords: docker, dtr, cli, backup
title: docker/ucp backup
---

Create a backup of a UCP manager node

## Usage

```bash

docker run -i --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  backup [command options] > backup.tar

```

## Description

This command creates a tar file with the contents of the volumes used by
this UCP manager node, and prints it. You can then use the 'restore' command to
restore the data from an existing backup.

To create backups of a multi-node swarm, you should create backups for all manager
nodes, one at a time in a quick succession, and keep track of the exact time and
sequence when you performed each backup. You will need this sequence information
if you restore more than one manager node at a time.

Note:

* During the backup, UCP is temporarily stopped. This does not affect your
  applications.

* The backup contains private keys and other sensitive information. Use the
  '--passphrase' flag to encrypt the backup with PGP-compatible encryption.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--interactive, i`|Run in interactive mode and prompt for configuration values|
|`--id`|The ID of the UCP instance to backup|
|`--passphrase`|Encrypt the tar file with a passphrase|
