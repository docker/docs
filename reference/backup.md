+++
title = "backup"
description = "Stream a tar file to stdout containing all UCP data volumes."
keywords= ["docker, ucp, backup, restore"]
[menu.main]
identifier = "ucp_ref_backup"
parent = "ucp_ref"
+++

# docker/ucp id

Stream a tar file to stdout containing all UCP data volumes.

## Usage

```bash
docker run --rm -i \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           backup [command options] > backup.tar
```

## Description

This utility will dump out a tar file containing all the contents of the
volumes used by UCP on this controller. This can be used to make periodic
backups suitable for use in the 'restore' command. Only UCP infrastructure
containers are backed up by this tool.

When backing up an HA cluster, take backups of all controllers, one at a
time, in quick succession, and keep track of the exact time and sequence
when you performed each backup. You will need this timestamp/sequence
information if you restore more than one controller together.

>**WARNING**: During the backup, all UCP infrastructure containers will be
temporarily stopped on this controller to prevent data corruption.  No user
containers will be stopped during the backup.

>**WARNING**: This backup will contain private keys and other sensitive information
and should be stored securely.  You may use the `--passphrase` flag to enable
built-in PGP compatible encryption.

## Options

```nohighlight
--debug, -D         Enable debug mode
--jsonlog           Produce json formatted output for easier parsing
--interactive, -i   Enable interactive mode.  You will be prompted to enter all required information
--id value          The ID of the UCP instance to backup
--root-ca-only      Backup only the root CA certificates and keys from this controller node
--passphrase value  Encrypt the tar file with the provided passphrase [$UCP_PASSPHRASE]
```
