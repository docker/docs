---
description: Frequently asked questions
keywords: mac faqs
redirect_from:
- /mackit/faqs/
title: Frequently asked questions (FAQ)
---

**Looking for popular FAQs on Docker for Mac?** Check out the
[Docker Success Center](http://success.docker.com/){: target="_blank" class="_"}
for knowledge base articles, FAQs, technical support for subscription levels, and more.

## Questions about Docker.app
### Stable and Edge channels

**Q: How do I get the Stable or Edge version of Docker for Mac?**

A: Use the download links for the channels given in the topic
[Download Docker for Mac](install.md#download-docker-for-mac){: target="_blank" class="_"}.

This topic also has more information about the two channels.

**Q: What is the difference between the Stable and Edge versions of Docker for Mac?**

A: Two different download channels are available for Docker for Mac:

* The **Stable channel** provides a general availability release-ready installer
  for a fully baked and tested, more reliable app. The Stable version of Docker
  for Mac comes with the latest released version of Docker Engine. The release
  schedule is synched with Docker Engine releases and hotfixes. On the Stable
  channel, you can select whether to send usage statistics and other data.

* The **Edge channel** provides an installer with new features we are working on,
  but is not necessarily fully tested. It comes with the experimental version of
  Docker Engine. Bugs, crashes, and issues are more likely to occur with the Edge
  app, but you get a chance to preview new functionality, experiment, and provide
  feedback as the apps evolve. Releases are typically more frequent than for
  Stable, often one or more per month. Usage statistics and crash reports are sent
  by default. You do not have the option to disable this on the Edge channel.

**Q: Can I switch back and forth between Stable and Edge versions of Docker for Mac?**

A: Yes, you can switch between versions to try out the Edge releases to see what's new,
then go back to Stable for other work. However, **you can have only one app
installed at a time**. Switching back and forth between Stable and Edge apps can
destabilize your development environment, particularly in cases where you switch
from a newer (Edge) channel to older (Stable).

For example, containers created with a newer Edge version of Docker for Mac may
not work after you switch back to Stable because they may have been created
leveraging Edge features that aren't in Stable yet. Just keep this in mind as
you create and work with Edge containers, perhaps in the spirit of a playground
space where you are prepared to troubleshoot or start over.

<font color="#CC3366">To safely switch between Edge and Stable versions be sure
to save images and export the containers you need, then uninstall the current
version before installing another. The workflow is described in more detail
below.</font><br>

Do the following each time:

1.  Use `docker save` to save any images you want to keep. (See
    [save](/engine/reference/commandline/save) in the Docker Engine command
    line reference.)

2.  Use `docker export` to export containers you want to keep. (See
    [export](/engine/reference/commandline/export) in the Docker Engine
    command line reference.)

3.  Uninstall the current app (whether Stable or Edge).

4.  Install a different version of the app (Stable or Edge).

### What is Docker.app?

`Docker.app` is Docker for Mac, a bundle of Docker client, and Docker Engine.
`Docker.app` uses the macOS Hypervisor.framework (part of macOS 10.10 Yosemite
and higher) to run containers, meaning that _**no separate VirtualBox is
required**_.

### What are system requirements for Docker for Mac?

You need a Mac that supports hardware virtualization and can run at
least macOS `10.10.3+` or `10.11` (macOS Yosemite or macOS El Capitan). See also
[What to know before you install](install#what-to-know-before-you-install) in
the install guide.

### Do I need to reinstall Docker for Mac if I change the name of my macOS account?

Yes, if you [change the name of your macOS user account and home
folder](https://support.apple.com/en-us/HT201548){: target="_blank" class="_"} after Docker for Mac is
already installed, you need to uninstall and reinstall Docker for Mac. Otherwise, Docker for Mac fails to start. See also,
[Docker for Mac does not start if Mac user account and home folder are renamed after installing Docker for Mac](troubleshoot.md#docker-for-mac-does-not-start-if-mac-user-account-and-home-folder-are-renamed-after-installing-the-app) in Troubleshooting.

### Do I need to uninstall Docker Toolbox to use Docker for Mac?

No, you can use these side by side. Docker Toolbox leverages a Docker daemon
installed using `docker-machine` in a machine called `default`. Running `eval
$(docker-machine env default)` in a shell sets DOCKER environment variables
locally to connect to the default machine using Engine from Toolbox. To check
whether Toolbox DOCKER environment variables are set, run `env | grep DOCKER`.

To make the client talk to the Docker for Mac Engine, run the command `unset
${!DOCKER_*}` to unset all DOCKER environment variables in the current shell.
(Now, `env | grep DOCKER` should return no output.) You can have multiple
command line shells open, some set to talk to Engine from Toolbox and others set
to talk to Docker for Mac. The same applies to `docker-compose`.

### How do I uninstall Docker Toolbox?

You might decide that you do not need Toolbox now that you have Docker for Mac,
and want to uninstall it. For details on how to perform a clean uninstall of
Toolbox on the Mac, see [How to uninstall Toolbox](/toolbox/toolbox_install_mac/#how-to-uninstall-toolbox)
in the Toolbox Mac topics.

## Questions about feedback and help
### What kind of feedback are we looking for?

Everything is fair game. We'd like your impressions on the download-install
process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about problems, what you like, or
functionality you'd like to see added.

We are especially interested in getting feedback on the new swarm mode described
in [Docker Swarm](/engine/swarm/index.md). A good place to start is the
[tutorial](/engine/swarm/swarm-tutorial/index.md).

### What if I have problems or questions?

You can find the list of frequent issues in
[Logs and Troubleshooting](troubleshoot.md).

If you do not find a solution in Troubleshooting, browse issues on
[Docker for Mac issues on GitHub](https://github.com/docker/for-mac/issues){: target="_blank" class="_"} or
create a new one. You can also create new issues based on diagnostics. To learn more,
see
[Diagnose problems, send feedback, and create GitHub issues](troubleshoot.md#diagnose-problems-send-feedback-and-create-github-issues).

[Docker for Mac forum](https://forums.docker.com/c/docker-for-mac){: target="_blank" class="_"}
provides discussion threads as well, and you can create discussion topics there,
but we recommend using the GitHub issues over the forums for better tracking and
response.

### How can I opt out of sending my usage data?

If you do not want auto-send of usage data, use the Stable channel. For more
information, see [Stable and Edge channels](#stable-and-edge-channels) ("What is the difference between the Stable and Edge versions of Docker for Mac?").


## How can I...?
### Can I use Docker for Mac with swarm mode?

Yes, you can use Docker for Mac to test single-node features of [swarm
mode](/engine/swarm/index.md) introduced with Docker Engine 1.12, including
initializing a swarm with a single node, creating services, and scaling
services. Docker “Moby” on Hyperkit serves as the single swarm node. You can
also use Docker Machine, which comes with Docker for Mac, to create and
experiment a multi-node swarm. Check out the tutorial at [Get started with swarm
mode](/engine/swarm/swarm-tutorial/index.md).

### How do I connect to the remote Docker Engine API?

You might need to provide the location of the Engine API for Docker clients and
development tools.

On Docker for Mac, clients can connect to the Docker Engine through a Unix
socket: `unix:///var/run/docker.sock`.

See also [Docker Engine API](/engine/api.md) and Docker for Mac forums topic
[Using pycharm Docker plugin..](https://forums.docker.com/t/using-pycharm-docker-plugin-with-docker-beta/8617){: target="_blank" class="_"}.

If you are working with applications like [Apache Maven](https://maven.apache.org/){: target="_blank" class="_"}
that expect settings for `DOCKER_HOST` and `DOCKER_CERT_PATH` environment
variables, specify these to connect to Docker instances through Unix sockets.
For example:

```bash
export DOCKER_HOST=unix:///var/run/docker.sock
```

### How do I connect from a container to a service on the host?

The Mac has a changing IP address (or none if you have no network access). Our
current recommendation is to attach an unused IP to the `lo0` interface on the
Mac so that containers can connect to this address.

For a full explanation and examples, see
[I want to connect from a container to a service on the host](networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host)
under
[Known Limitations, Use Cases, and Workarounds](networking.md#known-limitations-use-cases-and-workarounds)
in the Networking topic.

### How do I connect to a container from the Mac?

Our current recommendation is to publish a port, or to connect from another
container. This is what you need to do even on Linux if the container
is on an overlay network, not a bridge network, as these are not routed.

For a full explanation and examples, see
[I want to connect to a container from the Mac](networking.md#i-want-to-connect-to-a-container-from-the-mac)
under
[Known Limitations, Use Cases, and Workarounds](networking.md#known-limitations-use-cases-and-workarounds)
in the Networking topic.

### How do I add custom CA certificates?

Starting with Docker for Mac Beta 27 and Stable 1.12.3, all trusted certificate
authorities (CAs) (root or intermediate) are supported.

For full information on adding server and client side certs, see
[Add TLS certificates](index.md#add-tls-certificates) in the Getting Started topic.

### How do I add client certificates?

Starting with Docker for Mac 17.06.0-ce, you do not need to push your
certificates with `git` commands anymore. You can put your client certificates
in `~/.docker/certs.d/<MyRegistry>:<Port>/client.cert` and
`~/.docker/certs.d/<MyRegistry>:<Port>/client.key`.

For full information on adding server and client side certs, see
[Add TLS certificates](index.md#add-tls-certificates) in the Getting Started topic.

### Can I pass through a USB device to a container?

Unfortunately it is not possible to pass through a USB device (or a
serial port) to a container. For use cases requiring this, we
recommend the use of [Docker Toolbox](/toolbox/overview.md).

## Disk Usage

### What is the disk image?

The containers and images are stored in a _disk image_ named
`Docker.raw` or `Docker.qcow2` depending on your settings (see below).
By default, the disk image is stored in
`~/Library/Containers/com.docker.docker/Data`, in the `vms/0`
subdirectory starting with Docker for Mac 18.01, and
`com.docker.driver.amd64-linux` before.

### Qcow2 or Raw?

Starting with High Sierra with Apple Filesystem (APFS) enabled, Docker
uses disk images in the "raw" format (`Docker.raw`), otherwise in the
Qcow2 format (`Docker.qcow2`).

### Docker.raw consumes an insane amount of disk space!

This is an illusion.  Docker uses the raw format on Macs running the Apple
Filesystem (APFS).  APFS supports _sparse files_, which compress long runs of
zeroes representing unused space.  The output of `ls` is misleading, because it
lists the logical size of the file rather than its physical size. To see the
physical size, add the `-ks` switch; to see the logical size in human readable
form, add `-lh`:

```bash
$ cd ~/Library/Containers/com.docker.docker/Data
$ cd vms/0   # or com.docker.driver.amd64-linux
$ ls -klsh Docker.raw
2333548 -rw-r--r--@ 1 akim  staff    64G Dec 13 17:42 Docker.raw
```

In this listing, the logical size is 64GB, but the physical size is
only 2.3GB.

Alternatively, you may use `du` (disk usage):

```bash
$ du -h Docker.raw
2,2G	Docker.raw
```

### How do I reduce the size of Docker.qcow2?

If your Docker for Mac uses the Qcow format, the [disk image
file](#what-is-the-disk-image) is `Docker.qcow2`.  This file grows
on-demand up to a default maximum file size of 64GiB.

In Docker 1.12 the only way to free space on the host is to delete
this file and restart the app. Unfortunately this removes all images
and containers.

In Docker 1.13 there is preliminary support for "TRIM" to non-destructively
free space on the host. First free space within the `Docker.qcow2` by
removing unneeded containers and images with the following commands:

- `docker ps -a`: list all containers
- `docker image ls`: list all images
- `docker system prune`: (new in 1.13): deletes all stopped containers, all
  volumes not used by at least one container, and all images without at least one
  referring container.

Note the `Docker.qcow2` does not shrink in size immediately.
In 1.13 a background `cron` job runs `fstrim` every 15 minutes.
If the space needs to be reclaimed sooner, run this command:

```bash
$ docker run --rm -it --privileged --pid=host walkerlee/nsenter -t 1 -m -u -i -n fstrim /var
```

Once the `fstrim` has completed, restart the app. When the app shuts down, it
compacts the file and free up space. The app
takes longer than usual to restart because it must wait for the
compaction to complete.

For background conversation thread on this, see
[Docker.qcow2 never shrinks...](https://github.com/docker/for-mac/issues/371)
on Docker for Mac GitHub issues.


## Components of Docker for Mac
### What is HyperKit?

HyperKit is a hypervisor built on top of the Hypervisor.framework in macOS 10.10
Yosemite and higher. It runs entirely in userspace and has no other
dependencies.

We use HyperKit to eliminate the need for other VM products, such as Oracle
Virtualbox or VMWare Fusion.

### What is the benefit of HyperKit?

It is thinner than VirtualBox and VMWare fusion, and the version we include is
tailor made for Docker workloads on the Mac.

### Why is com.docker.vmnetd running after I quit the app?

The privileged helper process `com.docker.vmnetd` is started by `launchd` and
runs in the background. The process does not consume any resources unless
Docker.app connects to it, so it's safe to ignore.
