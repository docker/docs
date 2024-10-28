---
title: Virtual Machine Manager for Docker Desktop on Mac
linkTitle: Virtual Machine Manager 
params:
  sidebar:
    badge:
      color: green
      text: New
keywords: virtualization software, resource allocation, mac, docker desktop, vm monitoring, vm performance, apple silicon
description: Discover Docker Desktop for Mac's Virtual Machine Manager (VMM) options, including the new Docker VMM for Apple Silicon, offering enhanced performance and efficiency
weight: 150
---

The Virtual Machine Manager (VMM) in Docker Desktop for Mac is responsible for creating and managing the virtual machine used to run containers. Depending on your system architecture and performance needs, you can choose from multiple VMM options in Docker Desktop's [settings](/manuals/desktop/settings.md#general). This page provides an overview of the available options.

## Docker VMM (Beta)

Docker VMM is a new, container-optimized hypervisor introduced in Docker Desktop 4.35 and available on Apple Silicon Macs only. It has been designed to enhance speed and resource efficiency, making it an ideal choice for developers seeking to optimize their workflow efficiency. 

Docker VMM brings exciting advancements specifically tailored for Apple Silicon machines. By optimizing both the Linux kernel and hypervisor layers, Docker VMM delivers significant performance enhancements across common developer tasks. 

Some key performance enhancements provided by Docker VMM include:
 - Faster I/O operations: Iterating over a large shared filesystem with `find` is now 2x faster than in Docker Desktop 4.34 (with a cold cache).
 - Improved caching: With a warm cache, performance can improve by as much as 25x, even surpassing native Mac operations.

These improvements directly impact developers who rely on frequent file access, network interactions, and overall system responsiveness during containerized development. Docker VMM marks a significant leap in speed, enabling smoother workflows and faster iteration cycles.

### Known issues 

As Docker VMM is still in Beta, there are a few known limitations:

- Docker VMM does not currently support Rosetta, so emulation of amd64 architectures is slow. Docker is exploring potential solutions.
-  Certain databases, like MongoDB and Cassandra, may fail when using virtiofs with Docker VMM. This issue is expected to be resolved in a future release.

## Apple Virtualization Framework

The Apple Virtualization Framework is a stable and well-established option for managing virtual machines on Mac. It has been a reliable choice for many Mac users over the years. This framework is best suited for developers who prefer a proven solution with solid performance and broad compatibility.

## QEMU (Legacy) for Apple Silicon

> [!NOTE]
>
> QEMU will be deprecated in a future release.

QEMU is a legacy virtualization option for Apple Silicon Macs, primarily supported for older use cases. 

Docker recommends transitioning to newer alternatives, such as Docker VMM or the Apple Virtualization Framework, as they offer superior performance and ongoing support. Docker VMM, in particular, offers substantial speed improvements and a more efficient development environment, making it a compelling choice for developers working with Apple Silicon.

## HyperKit (Legacy) for Intel-based Macs

> [!NOTE]
>
> HyperKit will be deprecated in a future release.

HyperKit is another legacy virtualization option, specifically for Intel-based Macs. Like QEMU, it is still available but considered deprecated. Docker recommends switching to modern alternatives for better performance and to future-proof your setup.