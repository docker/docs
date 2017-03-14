---
title: "Getting Started, Part 1: Orientation and Setup"
---

## Orientation

While we'll define concepts along the way, it is good for you to
understand [what Docker is](#WHATDOCKER) and [why you would use
Docker](#WHYDOCKER) before we begin.

In this tutorial, you will:

1. Get set up and oriented, on this page.
2. [Create a "Hello World" application that identifies its environment](part2.md)
3. [Hook up a visitor counter](part3.md)
4. [Configure persistent data for our application](part4.md)
5. [Set up a scalable cluster in production](part5.md)
6. [Configure load-balancing and SSL](part6.md)
7. [Learn next steps](part7.md)

The application itself is very simple so that you are not too distracted by
what the code is doing. After all, the value of Docker is in how it can build,
ship, and run applications; it's totally agnostic as to what your application
actually does.

## Setup

Before we get started, make sure your system has the latest version of Docker
installed.

[Install Docker](/engine/installation/index.md){: class="button darkblue-btn"}

## A brief history of containers

A container image is a lightweight,
stand-alone, executable package of a piece of software that includes everything
needed to run it: code, runtime, system tools, system libraries, settings. 

A container is a runtime instance of an image -- what the image becomes in
memory when actually executed. 

A container image is like the disk image of a virtual machine, but without an OS.
That's because containers run apps natively on the host machine's kernel. They
don't need to have the performance characteristics of virtual machines that only
get virtual access to host resources -- containers can get native access.

Consider this diagram comparing the two:

![Virtual machine stack example](https://www.docker.com/sites/default/files/VM%402x.png)

Virtual machines run guest operating systems -- note the OS layer in each box. This
is resource intensive, and the resulting disk image is an entangelment of OS 
settings, and system-installed dependencies, and OS security patches, all bundled
into large multi-gigabyte files.

![Container stack example](https://www.docker.com/sites/default/files/Container%402x.png)

Containers can share a single kernel, and the only information packaged in a 
container image is the executable and its package dependencies, which never need
to be installed on the host system. These processes run like native processes, and
you can manage them individually by running commands like `docker ps` -- just like
you would run `ps` on Linux to see any other active executable.

## Conclusion

The unit of scale being an individual, portable executable means that CI/CD can push
updates to one part of a distributed application, system dependencies aren't a thing
you worry about, resource density is increased, and orchestrating scaling behavior
is a matter of spinning up new executables, not new VM hosts. We'll be learning about
all of those things, but first let's learn to walk.

[On to "Getting Started, Part 2: Creating and Building Your App" >>](part2.md){: class="button outline-btn"}
