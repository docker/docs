---
description: Frequently asked questions
keywords: windows faqs
title: Frequently asked questions (FAQ)
---

**Looking for popular FAQs on Docker for Windows?** Check out the
[Docker Success Center](http://success.docker.com/){: target="_blank" class="_"}
for knowledge base articles, FAQs, technical support for subscription levels, and more.

### Questions about Stable and Edge channels

#### How do I get the Stable or Edge version of Docker for Windows?

Use the download links for the channels given in the topic
[Download Docker for Windows](install#download-docker-for-windows).

This topic also has more information about the two channels.

#### What is the difference between the Stable and Edge versions of Docker for Windows?

Two different download channels are available for Docker for Windows:

* The **Stable channel** provides a general availability release-ready installer
  for a fully baked and tested, more reliable app. The Stable version of Docker
  for Windows comes with the latest released version of Docker Engine. The
  release schedule is synched with Docker Engine releases and hotfixes. On the
  Stable channel, you can select whether to send usage statistics and other data.

* The **Edge channel** provides an installer with new features we are working on,
  but is not necessarily fully tested. It comes with the experimental version of
  Docker Engine. Bugs, crashes, and issues are more likely to occur with the Edge
  app, but you get a chance to preview new functionality, experiment, and provide
  feedback as the apps evolve. Releases are typically more frequent than for
  Stable, often one or more per month. Usage statistics and crash reports are
  sent by default. You do not have the option to disable this on the Edge channel.

#### Can I switch back and forth between Stable and Edge versions of Docker for Windows?

Yes, you can switch between versions to try out the Edge release to see what's
new, then go back to Stable for other work. However, **you can have only one app
installed at a time**. Switching back and forth between Stable and Edge apps can
destabilize your development environment, particularly in cases where you switch
from a newer (Edge) channel to older (Stable).

For example, containers created with a newer Edge version of Docker for Windows
may not work after you switch back to Stable because they may have been created
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
    [export](/engine/reference/commandline/export) in the Docker Engine command
    line reference.)

3.  Uninstall the current app (whether Stable or Edge).

4.  Install a different version of the app (Stable or Edge).


### Feeback
#### What kind of feedback are we looking for?

Everything is fair game. We'd like your impressions on the download-install
process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about problems, what you like, or
functionality you'd like to see added.

We are especially interested in getting feedback on the new swarm mode described
in [Docker Swarm](/engine/swarm/). A good place to start is the
[tutorial](/engine/swarm/swarm-tutorial/).

#### What if I have problems or questions?

You can find the list of frequent issues in
[Logs and Troubleshooting](troubleshoot).

If you do not find a solution in Troubleshooting, browse issues on
[Docker for Windows issues on GitHub](https://github.com/docker/for-win/issues){: target="_blank" class="_"}
or create a new one. You can also create new issues based on diagnostics. To
learn more about running diagnostics and about Docker for Windows GitHub issues,
see [Diagnose and Feedback](/docker-for-windows#diagnose--feedback).

[Docker for Windows forum](https://forums.docker.com/c/docker-for-windows){: target="_blank" class="_"}
provides discussion threads as well, and you can create discussion topics there,
but we recommend using the GitHub issues over the forums for better tracking and
response.

#### How can I opt out of sending my usage data?

If you do not want auto-send of usage data, use the Stable channel. For more
information, see [Stable and Edge channels](#questions-about-stable-and-edge-channels) ("What is the difference between the Stable and Edge versions of Docker for Windows?").

### Can I use Docker for Windows with new swarm mode?

Yes! You can use Docker for Windows to test single-node features of
[swarm mode](/engine/swarm/) introduced with Docker Engine 1.12, including
initializing a swarm with a single node, creating services, and scaling
services. Docker “Moby” on Hyper-V serves as the single swarm node. You can also
use Docker Machine, which comes with Docker for Windows, to create and
experiment with a multi-node swarm. Check out the tutorial at
[Get started with swarm mode](/engine/swarm/swarm-tutorial/).

### How do I connect to the remote Docker Engine API?

You might need to provide the location of the Engine API for Docker clients and development tools.

On Docker for Windows, clients can connect to the Docker Engine through a
**named pipe**: `npipe:////./pipe/docker_engine`, or **TCP socket** at this URL:
`tcp://localhost:2375`.

This sets `DOCKER_HOST` and `DOCKER_CERT_PATH` environment variables to the
given values (for the named pipe or TCP socket, whichever you use).

See also [Docker Engine API](/engine/api) and the Docker for Windows forums
topic
[How to find the remote API](https://forums.docker.com/t/how-to-find-the-remote-api/20988){: target="_blank" class="_"}.

### Volumes
#### Can I change permissions on shared volumes for container-specific deployment requirements?

No, at this point, Docker for Windows does not enable you to control (`chmod`)
the Unix-style permissions on [shared volumes](/docker-for-windows#shared-drives) for
deployed containers, but rather sets permissions to a default value of
[0755](http://permissions-calculator.org/decode/0755/){: target="_blank" class="_"}
(`read`, `write`, `execute` permissions for `user`, `read` and `execute` for
`group`) which is not configurable.

For workarounds and to learn more, see
[Permissions errors on data directories for shared volumes](troubleshoot#permissions-errors-on-data-directories-for-shared-volumes).

#### Why doesn't `nodemon` pick up file changes in a container mounted on a shared drive?

Currently, `inotify` does not work on Docker for Windows. This is a known issue.
For more information and a temporary workaround, see
[inotify on shared drives does not work](troubleshoot#inotify-on-shared-drives-does-not-work){: target="_blank" class="_"}
in [Troubleshooting](troubleshoot).

#### Are symlinks supported?

Docker for Windows supports symbolic links (symlinks) created within containers.
Symlinks resolve within and across containers.
Symlinks created outside of Docker do not work.

To learn more about the reasons for this limitation, see the following discussions:

* GitHub issue:
  [Symlinks don't work as expected](https://github.com/docker/for-win/issues/109#issuecomment-251307391){: target="_blank" class="_"}

* Docker for Windows forums topic:
  [Symlinks on shared volumes not supported](https://forums.docker.com/t/symlinks-on-shared-volumes-not-supported/9288){: target="_blank" class="_"}


### Certificates

#### How do I add custom CA certificates?

Starting with Docker for Windows 1.12.1, 2016-09-16 (Stable) and Beta 26
(2016-09-14 1.12.1-beta26), all trusted Certificate Authorities (CA) (root or
intermediate) are supported. Docker recognizes certs stored under Trust Root
Certification Authorities or Intermediate Certification Authorities.

Docker for Windows creates a certificate bundle of all user-trusted CAs based on
the Windows certificate store, and appends it to Moby trusted certificates. So
if an enterprise SSL certificate is trusted by the user on the host, it is
trusted by Docker for Windows.

To learn more about how to install a CA root certificate for the registry, see
[Verify repository client with certificates](/engine/security/certificates)
in the Docker Engine topics.

#### How do I add client certificates?

Starting with Docker for Windows 17.06.0-ce, you do not need to push your
certificates with `git` commands anymore. You can put your client certificates
in `~/.docker/certs.d/<MyRegistry>:<Port>/client.cert` and
`~/.docker/certs.d/<MyRegistry>:<Port>/client.key`.

When the Docker for Windows application starts up, it copies the
`~/.docker/certs.d` folder on your Windows system to the `/etc/docker/certs.d`
directory on Moby (the Docker for Windows virtual machine running on Hyper-V).

You need to restart Docker for Windows after making any changes to the keychain
or to the `~/.docker/certs.d` directory in order for the changes to take effect.

The registry cannot be listed as an _insecure registry_ (see
[Docker Daemon](/docker-for-windows#daemon)). Docker for Windows ignores
certificates listed under insecure registries, and does not send client
certificates. Commands like `docker run` that attempt to pull from the registry
produce error messages on the command line, as well as on the registry.

To learn more about how to set the client TLS certificate for verification, see
[Verify repository client with certificates](/engine/security/certificates)
in the Docker Engine topics.

### Why does Docker for Windows sometimes lose network connectivity, causing `push` or `pull` commands to fail?

Networking is not yet fully stable across network changes and system sleep
cycles. Exit and start Docker to restore connectivity.

### Can I use VirtualBox alongside Docker 4 Windows?

Unfortunately, VirtualBox (and other hypervisors like VMWare) cannot run when
Hyper-V is enabled on Windows.

### Can I share local drives and filesystem with my Docker Machine VMs?

No, you cannot share local drives with Docker Machine nodes when using Docker
for Windows with Hyper-V. Shared drives can be made available to containers, but
Docker for Windows does not support mounts for nodes you created with
`docker-machine`.

For more about sharing local drives with containers using Docker for Windows,
see [Shared drives](/docker-for-windows#shared-drives) in the Getting
Started topic.

To learn more about using Docker for Windows and Docker Machine, see
[What to know before you install](install#what-to-know-before-you-install) in the
Getting Started topic. For more about Docker Machine itself, see
[What is Docker Machine?](/machine/overview#what-is-docker-machine), and the
[Hyper-V driver](/machine/drivers/hyper-v) for Docker Machine.

### Windows Requirements

#### How do I run Windows containers on Docker on Windows Server 2016?

See [About Windows containers and Windows Server 2016](/install/windows/docker-ee/#about-docker-ee-containers-and-windows-server).

A full tutorial is available in [docker/labs](https://github.com/docker/labs){: target="_blank" class="_"} at
[Getting Started with Windows Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md){: target="_blank" class="_"}.

#### Why is Windows 10 Home not supported?

Docker for Windows requires the Hyper-V Windows feature which is not
available on Home-edition.

#### Why is Windows 10 required?

Docker for Windows uses Windows Hyper-V. While older Windows versions have
Hyper-V, their Hyper-V implementations lack features critical for Docker for
Windows to work.

#### Why does Docker for Windows fail to start when firewalls or anti-virus software is installed?

Some firewalls and anti-virus software might be incompatible with Hyper-V and
some Windows 10 builds (possibly, the Anniversary Update), which impacts Docker
for Windows. See details and workarounds in
[Docker fails to start when firewall or anti-virus software is installed](troubleshoot#docker-fails-to-start-when-firewall-or-anti-virus-software-is-installed)
in [Troubleshooting](troubleshoot).

### How do I uninstall Docker Toolbox?

You might decide that you do not need Toolbox now that you have Docker for
Windows, and want to uninstall it. For details on how to perform a clean
uninstall of Toolbox on Windows, see
[How to uninstall Toolbox](/toolbox/toolbox_install_windows#how-to-uninstall-toolbox) in the
Toolbox Windows topics.
