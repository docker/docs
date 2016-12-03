---
description: Rackspace driver for machine
keywords:
- machine, Rackspace, driver
menu:
  main:
    parent: smn_machine_drivers
title: Rackspace
---

# Rackspace

Create machines on [Rackspace cloud](http://www.rackspace.com/cloud)

## Usage

    $ docker-machine create --driver rackspace --rackspace-username=user --rackspace-api-key=KEY --rackspace-region=region vm

## Options

-   `--rackspace-username`: **required** Rackspace account username.
-   `--rackspace-api-key`: **required** Rackspace API key.
-   `--rackspace-region`: **required** Rackspace region name.
-   `--rackspace-endpoint-type`: Rackspace endpoint type (`adminURL`, `internalURL` or the default `publicURL`).
-   `--rackspace-image-id`: Rackspace image ID. Default: Ubuntu 16.04 LTS (Xenial Xerus) (PVHVM).
-   `--rackspace-flavor-id`: Rackspace flavor ID. Default: General Purpose 1GB.
-   `--rackspace-ssh-user`: SSH user for the newly booted machine.
-   `--rackspace-ssh-port`: SSH port for the newly booted machine.
-   `--rackspace-docker-install`: Set if Docker has to be installed on the machine.
-   `--rackspace-active-timeout`: Rackspace active timeout

The Rackspace driver will use `821ba5f4-712d-4ec8-9c65-a3fa4bc500f9` (Ubuntu 16.04 LTS) by default.

#### Environment variables and default values

| CLI option                   | Environment variable | Default                                |
| ---------------------------- | -------------------- | -------------------------------------- |
| **`--rackspace-username`**   | `OS_USERNAME`        | -                                      |
| **`--rackspace-api-key`**    | `OS_API_KEY`         | -                                      |
| **`--rackspace-region`**     | `OS_REGION_NAME`     | `IAD` (Northern Virginia)              |
| `--rackspace-endpoint-type`  | `OS_ENDPOINT_TYPE`   | `publicURL`                            |
| `--rackspace-image-id`       | `OS_IMAGE_ID`        | `821ba5f4-712d-4ec8-9c65-a3fa4bc500f9` |
| `--rackspace-flavor-id`      | `OS_FLAVOR_ID`       | `general1-1`                           |
| `--rackspace-ssh-user`       | `OS_SSH_USER`        | `root`                                 |
| `--rackspace-ssh-port`       | `OS_SSH_PORT`        | `22`                                   |
| `--rackspace-docker-install` | -                    | `true`                                 |
| `--rackspace-active-timeout` | `OS_ACTIVE_TIMEOUT`  | `300`                                  |
