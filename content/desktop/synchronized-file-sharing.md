---
title: Synchronized file sharing
description: Get started with Synchronized file sharing on Docker Desktop.
sitemap: false
---

> **Early Access**
>
> Synchronized file sharing is an early-access feature. 
>
> If you would like to be considered for Synchronized file sharing testing, you can
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

After creating a file share instance, any container using a bind mount that references a host filesystem location which matches the stated synchronized file share host location, is provided by the Synchronized file sharing feature. Bind mounts that don't satisfy this condition are passed to the normal virtual filesystem bind-mounting mechanism, for example VirtioFS or gRPC-FUSE.

> **Important**
>
> Synchronized file sharing isn't available on WSL or when using Windows containers. 
{ .important }

## Set up

Before you get started with Synchronized file sharing, you need to:

- Download and install [Docker Desktop version 4.25](release-notes.md)
- Have a Pro, Team, or Business subscription
- Have a Docker ID that's part of a Docker organization participating in the Early access program

To enable Synchronized file sharing:
1. Sign in to Docker Desktop.
2. Anywhere within the Docker Dashboard, enter the Konami Code which has the key sequence of `up, up, down, down, left, right, left, b, a`.
3. From the list of experimental features, select **Mutagen File Sharing**.
4. At the bottom of the list, select **Go back to dashboard**

Synchronized file sharing is now enabled and ready for you to create your first file share instance.

## Create a file share instance 

To create a synchronized file share:
1. Navigate to **Settings** within Docker Desktop and then within the **Resource** section, select the **File sharing** tab. 
2. In the **Synchronized file shares** section, select the **plus** icon.
3. Select a host folder to share. The synchronized file share should initialize and be usable.

File shares take a few seconds to initialize as files are copied into the Docker Desktop VM. During this time, the status indicator displays **Preparing**.

When the status indicator displays **Watching for filesystem changes**, your files are available to the VM through all the standard bind mount mechanisms, whether that's `-v` in the command line or specified in your `compose.yml` file.

>**Note**
>
> When you create a new service, setting the [bind mount option consistency](../engine/reference/commandline/secret_create.md#options-for-bind-mounts) to `:consistent` bypasses synchronized file sharing. 

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
- Large dependency directories, for example `node_modules` and `composer` directories (unless you rely on accessing them via a bind mount)
- `.git directories` (again, unless you need them)

In general, the contents of your `.syncignore` file can include things anything that's going to cost you time to sync up-front and cost you storage. 

## Frequently asked questions

### What is the difference between Synchronized file sharing and the Mutagen extension?

Since Docker [acquired Mutagen](https://www.docker.com/blog/mutagen-acquisition/), Synchronized file sharing is essentially the direct integration of the Mutagen extension into Docker Desktop. However, Synchronized file sharing provides an improved user experience and the ability to use a [`.syncignore` file](#use-syncignore). You also no longer have to set the default ownership on files as this is now handled automatically.

## Known issues

- Upon launching Docker Desktop, it can take between 5-10 seconds for Synchronized file sharing to fully initialize. During this window, file share instances display as **Connecting** and any new containers created during this window won't replace bind mounts with Synchronized file sharing.

- `.syncignore` changes won't cause deletion until file share recreation. Files that become ignored due to a modification to ``.syncignore`` are left in place, but no longer update through synchronization.

- Case conflicts, due to Linux being case-sensitive and macOS/Windows only being case-preserving, display as **File exists** problems in the GUI.

- File share instances mounted into [ECI](hardened-desktop/enhanced-container-isolation/_index.md) containers are currently read-only.

- You cannot remove a file share instance during the initial synchronization. You have to wait for it to complete before **Delete** has any effect.

## Feedback

Thanks for trying the new Synchronized file sharing feature. Give feedback or report any bugs you may find through the #docker-eap-sync-file-share Slack channel.
