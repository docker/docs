---
aliases:
- /mackit/faqs/
description: Frequently asked questions
keywords:
- mac faqs
menu:
  main:
    identifier: docker-mac-faqs
    parent: pinata_mac_menu
    weight: 7
title: FAQs
---

#  Frequently Asked Questions (FAQs)

**Looking for popular FAQs on Docker for Mac?** Check out the [Docker Knowledge Hub](http://success.docker.com/) for knowledge base articles, FAQs, technical support for various subscription levels, and more.

## Stable and beta channels

**Q: How do I get the stable or beta version of Docker for Mac?**

A: Use the download links for the channels given in the topic [Download Docker for Mac](index.md#download-docker-for-mac).

This topic also has more information about the two channels.

**Q: What is the difference between the stable and beta versions of Docker for Mac?**

A: Two different download channels are available for Docker for Mac:

* The stable channel provides a general availability release-ready installer for a fully baked and tested, more reliable app. The stable version of Docker for Mac comes with the latest released version of Docker Engine. The release schedule is synched with Docker Engine releases and hotfixes.

* The beta channel provides an installer with new features we are working on, but is not necessarily fully tested. It comes with the experimental version of Docker Engine. Bugs, crashes and issues are more likely to occur with the beta app, but you get a chance to preview new functionality, experiment, and provide feedback as the apps evolve. Releases are typically more frequent than for stable, often one or more per month.


**Q: Can I switch back and forth between stable and beta versions of Docker for Mac?**

A: Yes, you can switch between versions; try out the betas to see what's new, then go back to stable for other work. However, doing so can de-stabalize your development environment, particularly in cases where you switch from a newer (beta) channel to older (stable).

For example, containers created with a newer beta version of Docker for Mac may not work after you switch back to stable because they may have been created leveraging beta features that aren't in stable yet. Just keep this in mind as you create and work with beta containers, perhaps in the spirit of a playground space where you are prepared to troubleshoot or start over.

## What is Docker.app?

`Docker.app` is Docker for Mac, a bundle of Docker client, and Docker
Engine. `Docker.app` uses the OS X
Hypervisor.framework (part of MacOS X 10.10 Yosemite and higher)
to run containers, meaning that _**no separate VirtualBox is required**_.

The Docker for Mac Beta is a work-in-progress, so some things will not work as well as expected, but we are interested in getting your feedback at this early stage.

## What kind of feedback are we looking for?

Everything is fair game. We'd like your impressions on the download-install process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about problems, what you like, or functionality you'd like to see added.

With 1.12-RC experimental features being included, we'd like to get some feedback on [Docker Swarm](https://docs.docker.com/engine/swarm/)

## What if I have problems or questions?

You can find the list of frequent issues in
[Logs and Troubleshooting](troubleshoot.md).

If your issue is not listed or solved there, please search and open a thread in the [Docker for Mac forum](https://forums.docker.com/c/docker-for-mac).

## What are system requirements for Docker for Mac?

Note that you need a Mac that supports hardware virtualization, which is most non ancient ones; i.e., use OS X `10.10.3+` or `10.11` (OS X Yosemite or OS X El Capitan). See also "What to know before you install" in [Getting Started](index.md).

<a name="faq-toolbox"></a>
## Do I need to uninstall Docker Toolbox to use Docker for Mac?

No, you can use these side by side. Docker Toolbox leverages a Docker daemon installed using `docker-machine` in a machine called `default`. Running `eval $(docker-machine env default)` in a shell sets DOCKER environment variables locally to connect to the default machine using Engine from Toolbox. To check whether Toolbox DOCKER environment variables are set, run `env | grep DOCKER`.

To make the client talk to the Docker for Mac Engine, run the command `unset ${!DOCKER_*}` to unset all DOCKER environment variables in the current shell. (Now, `env | grep DOCKER` should return no output.) You can have multiple command line shells open, some set to talk to Engine from Toolbox and others set to talk to Docker for Mac. The same applies to `docker-compose`.

## What is HyperKit?

HyperKit is a hypervisor built on top of the Hypervisor.framework in OS X 10.10 Yosemite and higher. It runs entirely in userspace and has no other dependencies.

We use HyperKit to eliminate the need for other VM products, such as Oracle Virtualbox or VMWare Fusion.

## What is the benefit of HyperKit?

It is thinner than VirtualBox and VMWare fusion, and the version we include is tailor made for Docker workloads on the Mac.

<p style="margin-bottom:300px">&nbsp;</p>
