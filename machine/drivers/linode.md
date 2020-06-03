---
description: Linode driver for machine
keywords: machine, Linode, driver
title: Linode
hide_from_sitemap: true
---

Create machines on [Linode](https://www.linode.com).

### Credentials

You will need a Linode APIv4 Personal Access Token.  Get one here: <https://developers.linode.com/api/v4#section/Personal-Access-Token>.
Supply the token to `docker-machine create -d linode` with `--linode-token`.

### Install

`docker-machine` is required, [see the installation documentation](https://docs.docker.com/machine/install-machine/).

Then, install the latest release of the Linode machine driver for your environment from the [releases list](https://github.com/linode/docker-machine-driver-linode/releases).

### Usage

```bash
docker-machine create -d linode --linode-token=<linode-token> linode
```

See the [Linode Docker machine driver project page](https://github.com/linode/docker-machine-driver-linode) for more examples.

#### Options, Environment Variables, and Defaults

| Argument | Env | Default | Description
| --- | --- | --- | ---
| `linode-token` | `LINODE_TOKEN` | None | **required** Linode APIv4 Token (see [here](https://developers.linode.com/api/v4#section/Personal-Access-Token))
| `linode-root-pass` | `LINODE_ROOT_PASSWORD` | *generated* | The Linode Instance `root_pass` (password assigned to the `root` account)
| `linode-authorized-users` | `LINODE_AUTHORIZED_USERS` | None | Linode user accounts (separated by commas) whose Linode SSH keys will be permitted root access to the created node
| `linode-label` | `LINODE_LABEL` | *generated* | The Linode Instance `label`, unless overridden this will match the docker-machine name.  This `label` must be unique on the account.
| `linode-region` | `LINODE_REGION` | `us-east` | The Linode Instance `region` (see [here](https://api.linode.com/v4/regions))
| `linode-instance-type` | `LINODE_INSTANCE_TYPE` | `g6-standard-4` | The Linode Instance `type` (see [here](https://api.linode.com/v4/linode/types))
| `linode-image` | `LINODE_IMAGE` | `linode/ubuntu18.04` | The Linode Instance `image` which provides the Linux distribution (see [here](https://api.linode.com/v4/images)).
| `linode-ssh-port` | `LINODE_SSH_PORT` | `22` | The port that SSH is running on, needed for Docker Machine to provision the Linode.
| `linode-ssh-user` | `LINODE_SSH_USER` | `root` | The user as which docker-machine should log in to the Linode instance to install Docker.  This user must have passwordless sudo.
| `linode-docker-port` | `LINODE_DOCKER_PORT` | `2376` | The TCP port of the Linode that Docker will be listening on
| `linode-swap-size` | `LINODE_SWAP_SIZE` | `512` | The amount of swap space provisioned on the Linode Instance
| `linode-stackscript` | `LINODE_STACKSCRIPT` | None | Specifies the Linode StackScript to use to create the instance, either by numeric ID, or using the form *username*/*label*.
| `linode-stackscript-data` | `LINODE_STACKSCRIPT_DATA` | None | A JSON string specifying data that is passed (via UDF) to the selected StackScript.
| `linode-create-private-ip` | `LINODE_CREATE_PRIVATE_IP` | None | A flag specifying to create private IP for the Linode instance.
| `linode-tags` | `LINODE_TAGS` | None | A comma separated list of tags to apply to the Linode resource
| `linode-ua-prefix` | `LINODE_UA_PREFIX` | None | Prefix the User-Agent in Linode API calls with some 'product/version'

#### Notes

* When using the `linode/containerlinux` `linode-image`, the `linode-ssh-user` will default to `core`
* A `linode-root-pass` will be generated if not provided.  This password will not be shown. Rely on `docker-machine ssh` or [Linode's Rescue features](https://www.linode.com/docs/quick-answers/linode-platform/reset-the-root-password-on-your-linode/) to access the node directly.

#### Debugging

Detailed run output will be emitted when using the LinodeGo `LINODE_DEBUG=1` option along with the `docker-machine` `--debug` option.

```bash
LINODE_DEBUG=1 docker-machine --debug  create -d linode --linode-token=$LINODE_TOKEN machinename
```

