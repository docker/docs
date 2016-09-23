---
description: Frequently asked questions
keywords:
- windows faqs
menu:
  main:
    identifier: docker-windows-faqs
    parent: pinata_win_menu
    weight: 4
title: FAQs
---

#  Frequently Asked Questions (FAQs)

>**Looking for popular FAQs on Docker for Windows?** Check out the [Docker Knowledge Hub](http://success.docker.com/) for knowledge base articles, FAQs, technical support for various subscription levels, and more.

## Questions about stable and beta channels

**Q: How do I get the stable or beta version of Docker for Windows?**

A: Use the download links for the channels given in the topic [Download Docker for Windows](index.md#download-docker-for-windows).

This topic also has more information about the two channels.

**Q: What is the difference between the stable and beta versions of Docker for Windows?**

A: Two different download channels are available for Docker for Windows:

* The stable channel provides a general availability release-ready installer for a fully baked and tested, more reliable app. The stable version of Docker for Windows comes with the latest released version of Docker Engine.  The release schedule is synched with Docker Engine releases and hotfixes.

* The beta channel provides an installer with new features we are working on, but is not necessarily fully tested. It comes with the experimental version of Docker Engine. Bugs, crashes and issues are more likely to occur with the beta app, but you get a chance to preview new functionality, experiment, and provide feedback as the apps evolve. Releases are typically more frequent than for stable, often one or more per month.


**Q: Can I switch back and forth between stable and beta versions of Docker for Windows?**

A: Yes, you can switch between versions; try out the betas to see what's new, then go back to stable for other work. However, doing so can de-stabilize your development environment, particularly in cases where you switch from a newer (beta) channel to older (stable).

For example, containers created with a newer beta version of Docker for Windows may not work after you switch back to stable because they may have been created leveraging beta features that aren't in stable yet. Just keep this in mind as you create and work with beta containers, perhaps in the spirit of a playground space where you are prepared to troubleshoot or start over.

## Other Questions

**Q: Why does Docker for Windows sometimes lose network connectivity (for example, `push`/`pull` doesn't work)?**

A: Networking is not yet fully stable across network changes and system sleep cycles. Exit and start Docker to restore connectivity.

**Q: Can I use VirtualBox alongside Docker 4 Windows?**

A: Unfortunately, VirtualBox (and other hypervisors like VMWare) cannot run when Hyper-V is enabled on Windows.

**Q: Why is Windows 10 Home not supported?**

A: Docker for Windows requires the Hyper-V Windows feature which is not available on Home-edition.

**Q: Why is Windows 10 required?**

A: Docker for Windows uses Windows Hyper-V. While older Windows versions have Hyper-V, their Hyper-V implementations lack features critical for Docker for Windows to work.

<p style="margin-bottom:300px">&nbsp;</p>
