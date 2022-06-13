---
description: Frequently asked questions
keywords: desktop, mac, faqs
title: Frequently asked questions for Mac
---

## Mac FAQs

### What is Docker.app?

`Docker.app` is Docker Desktop on Mac. It bundles the Docker client and Docker Engine. `Docker.app` uses the macOS Hypervisor.framework to run containers.

### What is HyperKit?

HyperKit is a hypervisor built on top of the Hypervisor.framework in macOS. It runs entirely in userspace and has no other dependencies.

We use HyperKit to eliminate the need for other VM products, such as Oracle
VirtualBox or VMWare Fusion.

### What is the benefit of HyperKit?

HyperKit is thinner than VirtualBox and VMWare fusion, and the version included is customized for Docker workloads on Mac.

### Why is com.docker.vmnetd still running after I quit the app?

The privileged helper process `com.docker.vmnetd` is started by `launchd` and
runs in the background. The process does not consume any resources unless
Docker.app connects to it, so it's safe to ignore.
