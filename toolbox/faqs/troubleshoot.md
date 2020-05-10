---
advisory: toolbox
description: Troubleshooting connectivity and certificate issues
keywords: beginner, getting started, FAQs, troubleshooting, Docker
title: Troubleshooting
---

Typically, the QuickStart works out-of-the-box, but some scenarios can cause problems.

## Example errors

You might get errors when attempting to connect to a machine (such as with `docker-machine env default`) or pull an image from Docker Hub (as with `docker run hello-world`).

The errors you get might be specific to certificates, like this:

      Error checking TLS connection: Error checking and/or regenerating the certs: There was an error validating certificates for host "192.168.99.100:2376": dial tcp 192.168.99.100:2376: i/o timeout

Others explicitly suggest regenerating certificates:

      Error checking TLS connection: Error checking and/or regenerating the certs: There was an error validating certificates for host "192.168.99.100:2376": x509: certificate is valid for 192.168.99.101, not 192.168.99.100
      You can attempt to regenerate them using 'docker-machine regenerate-certs [name]'.
      Be advised that this will trigger a Docker daemon restart which will stop running containers.

Or, indicate a network timeout, like this:

      bash-3.2$ docker run hello-world
      Unable to find image 'hello-world:latest' locally
      Pulling repository docker.io/library/hello-world
      Network timed out while trying to connect to https://index.docker.io/v1/repositories/library/hello-world/images. You may want to check your internet connection or if you are behind a proxy.
      bash-3.2$

## Solutions

Here are some quick solutions to help get back on track. These examples assume the Docker host is a machine called `default`.

#### Regenerate certificates

Some errors explicitly tell you to regenerate certificates. You might also try this for other errors that are certificate and/or connectivity related.

      $ docker-machine regenerate-certs default
        Regenerate TLS machine certs?  Warning: this is irreversible. (y/n): y
        Regenerating TLS certificates

#### Restart the Docker host

    $ docker-machine restart default

After the machine starts, set the environment variables for the command window.

    $ eval $(docker-machine env default)

Run `docker-machine ls` to verify that the machine is running and that this command window is configured to talk to it, as indicated by an asterisk for the active machine (__*__).

    $ docker-machine ls
    NAME             ACTIVE   DRIVER         STATE     URL                         SWARM   DOCKER    ERRORS
    default          *        virtualbox     Running   tcp://192.168.99.101:2376           v1.10.1

#### Stop the machine, remove it, and create a new one.

    $ docker-machine stop default
      Stopping "default"...
      Machine "default" was stopped.

    $ docker-machine rm default
      About to remove default
      Are you sure? (y/n): y
      Successfully removed default

You can use the `docker-machine create` command with the `virtualbox` driver to create a new machine called `default` (or any name you want for the machine).

    $ docker-machine create --driver virtualbox default
      Running pre-create checks...
      (default) Default Boot2Docker ISO is out-of-date, downloading the latest release...
      (default) Latest release for github.com/boot2docker/boot2docker is v1.10.1
      (default) Downloading
      ...
      Docker is up and running!
      To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env default

Set the environment variables for the command window.

    $ eval $(docker-machine env default)

Run `docker-machine ls` to verify that the new machine is running and that this command window is configured to talk to it, as indicated by an asterisk for the active machine (__*__).

<a name="machine-http-proxy-solutions"></a>

## HTTP proxies and connectivity errors

A special brand of connectivity errors can be caused by HTTP proxy. If you install Docker Toolbox on a system using a virtual private network (VPN) that uses an HTTP proxy (such as a corporate network), you might encounter errors when the client attempts to connect to the server.

Here are examples of this type of error:

      $ docker run hello-world
      An error occurred trying to connect: Post https://192.168.99.100:2376/v1.20/containers/create: Forbidden

      $ docker run ubuntu echo "hi"
      An error occurred trying to connect: Post https://192.168.99.100:2376/v1.20/containers/create: Forbidden

### Configure HTTP proxy settings on Docker machines

When Toolbox creates virtual machines (VMs) it runs `start.sh`, where it gets values for `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY`, and passes them as `create` options to create the `default machine`.

You can reconfigure HTTP proxy settings for private networks on already-created Docker machines, such as the `default` machine, then change the configuration when you are using the same system on a different network.

Alternatively, you can modify proxy settings on your machine(s) manually through the configuration file at `/var/lib/boot2docker/profile` inside the VM, or configure proxy settings as a part of a `docker-machine create` command.

Both solutions are described below.

#### Update /var/lib/boot2docker/profile on the Docker machine

One way to solve this problem is to update the file `/var/lib/boot2docker/profile` on an existing machine to specify the proxy settings you want.

This file lives on the VM itself, so you need to `ssh` into the machine, then edit and save the file there.

You can add your machine addresses as values for a `NO_PROXY` setting, and also specify proxy servers that you know about and you want to use. Typically setting your Docker machine URLs to `NO_PROXY` solves this type of connectivity problem, so that example is shown here.

1. Use `ssh` to log in to the virtual machine. This example logs in to the
   `default` machine.

        $ docker-machine ssh default
        docker@default:~$ sudo vi /var/lib/boot2docker/profile

2. Add a `NO_PROXY` setting to the end of the file similar to the example below.

        # replace with your office's proxy environment
        export "HTTP_PROXY=http://PROXY:PORT"
        export "HTTPS_PROXY=http://PROXY:PORT"
        # you can add more no_proxy with your environment.
        export "NO_PROXY=192.168.99.*,*.local,169.254/16,*.example.com,192.168.59.*"

3. Restart Docker.

    After you modify the `profile` on your VM, restart Docker and log out of the machine.

        docker@default:~$ sudo /etc/init.d/docker restart
        docker@default:~$ exit

  Re-try Docker commands. Both Docker and Kitematic should run properly now.

  When you move to a different network (for example, leave the office's corporate network and return home), remove or comment out these proxy settings in `/var/lib/boot2docker/profile` and restart Docker.

#### Create machines manually using --engine env to specify proxy settings

Rather than reconfigure automatically-created machines, you can delete them and create your `default` machine and others manually with the `docker-machine create` command, using the `--engine env` flag to specify the proxy settings you want.

Here is an example of creating a `default` machine with proxies set to `http://example.com:8080` and `https://example.com:8080`, and a `N0_PROXY` setting for the server `example2.com`.

    docker-machine create -d virtualbox \
    --engine-env HTTP_PROXY=http://example.com:8080 \
    --engine-env HTTPS_PROXY=https://example.com:8080 \
    --engine-env NO_PROXY=example2.com \
    default


To learn more about using `docker-machine create`, see the [create](../../machine/reference/create.md) command in the [Docker Machine](../../machine/overview.md) reference.

&nbsp;
