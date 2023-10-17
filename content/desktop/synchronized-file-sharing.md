---
title: Synchronized file sharing
description: Get started with Synchronized file sharing on Docker Desktop.
sitemap: false
---

> **Early Access**
>
> Synchronized file sharing is an early-access feature. 
>
> If you want to test out Synchronized file sharing, you can
> [sign up for the early access program](https://www.docker.com/build-early-access-program/?utm_source=docs).
{ .restricted }

Synchronized file sharing is an alternative file sharing mechanism powered by [Mutagen](https://mutagen.io/). It provides fast and flexible host-to-VM file sharing by replacing bind mounts with synchronized filesystem caches. 

![Image of Synchronized file shares pane](images/synched-file-sharing.png)
 
## Who is it for?

For developers who: 
- Have large repositories or monorepos with 100 000 files or more totaling hundreds of megabytes or even gigabytes
- Are using virtual filesystems, such as VirtioFS, gRPC FUSE, and osxfs, which are no longer scaling well with their codebases 
- Regularly encounter performance limitations
- Don't want to worry about file ownership or spend time resolving conflicting file-ownership information when modifying multiple containers

## How does Synchronized file sharing work?

A Synchronized file share behaves just like a virtual file share, but takes advantage of Mutagen's high-performance, low-latency code synchronization engine to create a synchronized cache of the host files on an ext4 filesystem within the Docker Desktop VM. If you make filesystem changes on the host or in the VM’s containers it propagates via bidirectional synchronization.

Synchronized file sharing also adds an additional Docker Context, `desktop-linux-mutagen`, which uses Docker socket middleware to automatically replace bind mounts with in-VM caches whenever possible.

After creating a file share instance, any container using a bind mount that references a host filesystem location which matches the stated synchronized file share host location, is provided by the Synchronized file sharing feature. Bind mounts that don't satisfy this condition are passed to the normal virtual filesystem bind-mounting mechanism, for example VirtioFS or gRPC-FUSE.

> **Important**
>
> Synchronized file sharing isn't available on WSL or when using Windows containers. 
{ .important }

## Set up - TBD


- DB customer
- Signed in to DD 
- DD version 4.25

## Create a file share instance 

To create a synchronized file share:
1. Navigate to **Settings** within Docker Desktop and then within the **Resource** section, select the **File sharing** tab. 
2. In the **Synchronized file shares** section, select the **plus** icon.
3. Select a host folder to share. The synchronized file share should initialize and be usable.

File shares take a few seconds to initialize as files are copied into the Docker Desktop VM. During this time, the status indicator displays **Preparing**.

When the status indicator displays **Watching for filesystem changes**, your files are available to the VM through all the standard bind mount mechanisms, whether that's `-vm` in the command line or specified in your `compose.yml` file.

## Explore your file share instance

The **Synchronized file shares** section displays all your file share instances and provides useful information about each instance including:
- The origin of the file share content
- A Status update
- How much space each file share is using
- The number of filesystem entry counts
- The number of symbolic links
- Which container(s) is using the file share instance

Selecting a file share instance expands the dropdown and exposes this information.

## Use `.syncignore`

You can use a `.syncignore` file at the root of each file share, to exclude local files from your file share instance. It supports the same syntax as `.dockerignore` files and excludes, and/or re-includes, paths from synchronization. `.syncignore` files are ignored at any location other than the root of the file share.
 
Some example of things you might want to add to your `.syncignore` file are:
- ?

In general, the contents of your `.syncignore` file should be similar to what you have in your `.gitignore` file.

## Frequently asked questions

### What is the difference between Synchronized file sharing and the Mutagen extension?

Synchronized file sharing is essentially a replacement for the Mutagen extension. However, it provides an improved user experience and the ability to use a [`.syncignore` file](#use-syncignore). You also no longer have to set the default ownership on files as this is now handled automatically.

## Known issues

- Synchronized file sharing currently doesn't work with [Resource Saver mode](use-desktop/resource-saver.md). Make sure you turn off Resource Saver mode before using Synchronized file sharing.

- Upon launching Docker Desktop, it can take between 5-10 seconds for Synchronized file sharing to fully initialize. During this window, file share instances display as **Connecting** and any new containers created during this window won't replace bind mounts with Synchronized file sharing.

## Feedback

Thanks for trying the new Synchronized file sharing feature. Give feedback or report any bugs you may find through the issues tracker on the feedback form.
