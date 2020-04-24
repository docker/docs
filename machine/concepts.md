---
description: Understand concepts for Docker Machine, including drivers, base OS, IP addresses, environment variables
keywords: docker, machine, amazonec2, azure, digitalocean, google, openstack, rackspace, softlayer, virtualbox, vmwarefusion, vmwarevcloudair, vmwarevsphere, exoscale
title: Machine concepts and getting help
---

Docker Machine allows you to provision Docker machines in a variety of environments, including virtual machines that reside on your local system, on cloud providers, or on bare metal servers (physical computers). Docker Machine creates a Docker host, and you use the Docker Engine client as needed to build images and create containers on the host.

## Drivers for creating machines

To create a virtual machine, you supply Docker Machine with the name of the driver you want to use. The driver determines where the virtual machine is created. For example, on a local Mac or Windows system, the driver is typically Oracle VirtualBox. For provisioning physical machines, a generic driver is provided. For cloud providers, Docker Machine supports drivers such as AWS, Microsoft Azure, DigitalOcean, and many more. The Docker Machine reference includes a complete [list of supported drivers](drivers/index.md).

## Default base operating systems for local and cloud hosts

Since Docker runs on Linux, each VM that Docker Machine provisions relies on a
base operating system. For convenience, there are default base operating
systems. For the Oracle Virtual Box driver, this base operating system
is [boot2docker](https://github.com/boot2docker/boot2docker). For drivers used
to connect to cloud providers, the base operating system is Ubuntu 12.04+. You
can change this default when you create a machine. The Docker Machine reference
includes a complete [list of supported operating systems](drivers/os-base.md).

## IP addresses for Docker hosts

For each machine you create, the Docker host address is the IP address of the
Linux VM. This address is assigned by the `docker-machine create` subcommand.
You use the `docker-machine ls` command to list the machines you have created.
The `docker-machine ip <machine-name>` command returns a specific host's IP
address.

## Configuring CLI environment variables for a Docker host

Before you can run a `docker` command on a machine, you need to configure your
command-line to point to that machine. The `docker-machine env <machine-name>`
subcommand outputs the configuration command you should use.

For a complete list of `docker-machine` subcommands, see the
[Docker Machine subcommand reference](reference/help.md).

## Custom root Certificate Authority for Registry

if your registry is signed by a custom root Certificate Authority and it is
not registered with Docker Engine, you may see the following error message:

```none
x509: certificate signed by unknown authority
```

As discussed in the
[Docker Engine documentation](../engine/security/certificates.md#understand-the-configuration)
place the certificates in `/etc/docker/certs.d/hostname/ca.crt`
where `hostname` is your Registry server's hostname.

```console
docker-machine scp certfile default:ca.crt
docker-machine ssh default
sudo mv ~/ca.crt /etc/docker/certs.d/hostname/ca.crt
exit
docker-machine restart
```

## Crash reporting

Provisioning a host is a complex matter that can fail for a lot of reasons. Your
workstation may have a wide variety of shell, network configuration, VPN, proxy
or firewall issues. There are also reasons from the other end of the chain:
your cloud provider or the network in between.

To help `docker-machine` be as stable as possible, we added a monitoring of
crashes whenever you try to `create` or `upgrade` a host. This sends, over
HTTPS, to Bugsnag some information about your `docker-machine` version, build,
OS, ARCH, the path to your current shell and, the history of the last command as
you could see it with a `--debug` option. This data is sent to help us pinpoint
recurring issues with `docker-machine` and is only transmitted in the case
of a crash of `docker-machine`.

To opt out of error reporting, create a `no-error-report`
file in your `$HOME/.docker/machine` directory:

    $ mkdir -p ~/.docker/machine && touch ~/.docker/machine/no-error-report

The file doesn't need to have any contents.

## Getting help

Docker Machine is still in its infancy and under active development. If you need
help, would like to contribute, or simply want to talk about the project with
like-minded individuals, we have a number of open channels for communication.

- To report bugs or file feature requests, use the
  [issue tracker on Github](https://github.com/docker/machine/issues).
- To talk about the project with people in real time,  join the
  `#docker-machine` channel on IRC.
- To contribute code or documentation changes,
  [submit a pull request on Github](https://github.com/docker/machine/pulls).

For more information and resources, visit
[our help page](../opensource/ways.md).

## Where to go next

- Create and run a Docker host on your [local system using VirtualBox](get-started.md)
- Provision multiple Docker hosts [on your cloud provider](get-started-cloud.md)
- [Docker Machine driver reference](drivers/index.md){: target="_blank" class="_"}
- [Docker Machine subcommand reference](reference/help.md){: target="_blank" class="_"}
