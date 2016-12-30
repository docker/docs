---
description: Frequently asked questions
keywords: windows faqs
title: Frequently asked questions (FAQ)
---

>**Looking for popular FAQs on Docker for Windows?** Check out the [Docker
Knowledge Hub](http://success.docker.com/) for knowledge base articles, FAQs,
technical support for various subscription levels, and more.

### Questions about stable and beta channels

**Q: How do I get the stable or beta version of Docker for Windows?**

A: Use the download links for the channels given in the topic [Download Docker
for Windows](index.md#download-docker-for-windows).

This topic also has more information about the two channels.

**Q: What is the difference between the stable and beta versions of Docker for Windows?**

A: Two different download channels are available for Docker for Windows:

* The **stable channel** provides a general availability release-ready installer for a fully baked and tested, more reliable app. The stable version of Docker
for Windows comes with the latest released version of Docker Engine.  The
release schedule is synched with Docker Engine releases and hotfixes.  On the stable channel, you can select whether to send usage statistics and other data.

* The **beta channel** provides an installer with new features we are working on, but is not necessarily fully tested. It comes with the experimental version of
Docker Engine. Bugs, crashes and issues are more likely to occur with the beta
app, but you get a chance to preview new functionality, experiment, and provide
feedback as the apps evolve. Releases are typically more frequent than for
stable, often one or more per month.  Usage statistics and crash reports are sent by default. You do not have the option to disable this on the beta channel.

**Q: Can I switch back and forth between stable and beta versions of Docker for Windows?**

A: Yes, you can switch between versions to try out the betas to see what's new,
then go back to stable for other work. However, **you can have only one app
installed at a time**. Switching back and forth between stable and beta apps can
destabilize your development environment, particularly in cases where you
switch from a newer (beta) channel to older (stable).

For example, containers created with a newer beta version of Docker for Windows
may not work after you switch back to stable because they may have been created
leveraging beta features that aren't in stable yet. Just keep this in mind as
you create and work with beta containers, perhaps in the spirit of a playground
space where you are prepared to troubleshoot or start over.

<font color="#CC3366">To safely switch between beta and stable versions be sure
to save images and export the containers you need, then uninstall the current
version before installing another. The workflow is described in more detail
below.</font><br>

Do the following each time:

1. Use `docker save` to save any images you want to keep. (See
[save](/engine/reference/commandline/save.md) in the Docker Engine command line
reference.)

2. Use `docker export` to export containers you want to keep. (See
[export](/engine/reference/commandline/export.md) in the Docker Engine command
line reference.)

3. Uninstall the current app (whether stable or beta).

4. Install a different version of the app (stable or beta).

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

If you do not find a solution in Troubleshooting, browse issues on [Docker for
Windows issues on GitHub](https://github.com/docker/for-win/issues) or create a
new one. You can also create new issues based on diagnostics. To learn more
about running diagnostics and about Docker for Windows GitHub issues, see
[Diagnose and Feedback](index.md#diagnose-and-feedback).

[Docker for Windows forum](https://forums.docker.com/c/docker-for-windows)
provides discussion threads as well, and you can create discussion topics there,
but we recommend using the GitHub issues over the forums for better tracking and
response.

### Can I use Docker for Windows with new swarm mode?

Yes! You can use Docker for Windows to test single-node features of [swarm
mode](/engine/swarm/index.md) introduced with Docker Engine 1.12, including
initializing a swarm with a single node, creating services, and scaling
services. Docker “Moby” on Hyper-V will serve as the single swarm node. You can
also use Docker Machine, which comes with Docker for Windows, to create and
experiment with a multi-node swarm. Check out the tutorial at [Get started with
swarm mode](/engine/swarm/swarm-tutorial/index.md).

### How do I connect to the remote Docker Engine API?

You might need to provide the location of the remote API for Docker clients and development tools.

On Docker for Windows, clients can connect to the Docker Engine through a **named pipe**: `npipe:////./pipe/docker_engine`, or **TCP socket** at this URL: `http://localhost:2375`.

This sets `DOCKER_HOST` and `DOCKER_CERT_PATH` environment variables to the given values (for the named pipe or TCP socket, whichever you use).

See also [Docker Remote API](/engine/reference/api/docker_remote_api.md) and the Docker for Windows forums topic [How to find the remote API](https://forums.docker.com/t/how-to-find-the-remote-api/20988).

### Why doesn't `nodemon` pick up file changes in a container mounted on a shared drive?

Currently, `inotify` does not work on Docker for Windows. This is a known issue.
For more information and a temporary workaround, see [inotify on shared drives
does not work](troubleshoot.md#inotify-on-shared-drives-does-not-work) in
[Troubleshooting](troubleshoot.md).

### Are symlinks supported?

Docker for Windows supports symbolic links (symlinks) created within containers.
Symlinks will resolve within and across containers.
Symlinks created elsewhere (e.g., on the host) will not work.

To learn more about the reasons for this limitation, see the following discussions:

* GitHub issue: [Symlinks don't work as expected](https://github.com/docker/for-win/issues/109#issuecomment-251307391)

* Docker for Windows forums topic: [Symlinks on shared volumes not supported](https://forums.docker.com/t/symlinks-on-shared-volumes-not-supported/9288)

### How do I add custom CA certificates?

Starting with Docker for Windows 1.12.1, 2016-09-16 (stable) and Beta 26 (2016-09-14 1.12.1-beta26), all trusted CAs (root or intermediate) are supported. Docker recognizes certs stored under Trust Root Certification Authorities or Intermediate Certification Authorities.

Docker for Windows creates a certificate bundle of all user-trusted CAs based on the Windows certificate store, and appends it to Moby trusted certificates. So if an enterprise SSL certificate is trusted by the user on the host, it will be trusted by Docker for Windows.

To learn more, see the GitHub issue [Allow user to add custom Certificate Authorities ](https://github.com/docker/for-win/issues/48).

### Why does Docker for Windows sometimes lose network connectivity (e.g., `push`/`pull` doesn't work)?

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
see [Shared Drives](index.md#shared-drives) in the Getting Started topic.

To learn more about using Docker for Windows and Docker Machine, see [What to
know before you install](index.md#what-to-know-before-you-install) in the
Getting Started topic. For more about Docker Machine itself, see [What is Docker
Machine?](/machine/overview.md#what-is-docker-machine)

### How do I run Windows containers on Docker on Windows Server 2016?

See [About Windows containers and Windows Server
2016](index.md#about-windows-containers-and-windows-server-2016).

A full tutorial is available in [docker/labs](https://github.com/docker/labs) at
[Getting Started with Windows
Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md).

### Why is Windows 10 Home not supported?

Docker for Windows requires the Hyper-V Windows feature which is not
available on Home-edition.

### Why is Windows 10 required?

Docker for Windows uses Windows Hyper-V. While older Windows versions have
Hyper-V, their Hyper-V implementations lack features critical for Docker for
Windows to work.

### Why does Docker for Windows fail to start when firewalls or anti-virus software is installed?

Some firewalls and anti-virus software might be incompatible with Hyper-V and
some Windows 10 builds  (possibly, the Anniversary Update), which impacts Docker
for Windows. See details and workarounds in [Docker fails to start when firewall
or anti-virus software is
installed](troubleshoot.md#docker-fails-to-start-when-firewall-or-anti-virus-software-is-installed)
in [Troubleshooting](troubleshoot.md).

### How do I uninstall Docker Toolbox?

You might decide that you do not need Toolbox now that you have Docker for Windows, and want to uninstall it. For
details on how to perform a clean uninstall of Toolbox on Windows, see [How to
uninstall Toolbox](/toolbox/toolbox_install_windows.md#how-to-uninstall-toolbox)
in the Toolbox Windows topics.
